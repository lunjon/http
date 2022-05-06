package command

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/lunjon/http/client"
	"github.com/lunjon/http/logging"
	"github.com/lunjon/http/util"
	"github.com/spf13/cobra"
)

type FailFunc func(status int)
type runFunc func(*cobra.Command, []string)
type checkRedirectFunc func(*http.Request, []*http.Request) error

type execError struct {
	err       error
	showUsage bool
}

func newUserError(err error) *execError {
	return &execError{
		err:       err,
		showUsage: true,
	}
}

func (e *execError) Error() string {
	return e.err.Error()
}

func (e *execError) Unwrap() error {
	return e.err
}

var (
	noConfigure   = func(*cobra.Command) {}
	bodyConfigure = func(cmd *cobra.Command) {
		cmd.Flags().StringP(
			bodyFlagName,
			"B",
			"",
			"Request body to use. Can be string content or a filename.",
		)
	}
)

const (
	defaultTimeout    = time.Second * 30
	defaultAWSRegion  = "eu-west-1"
	defaultHeadersEnv = "DEFAULT_HEADERS"
)

// Build the root command for http and set version.
func Build(version string) (*cobra.Command, error) {
	cfg, err := newDefaultConfig(version)
	if err != nil {
		return nil, err
	}
	return build(version, cfg), nil
}

func build(version string, cfg *config) *cobra.Command {
	root := &cobra.Command{
		Use:   "http",
		Short: "http <method> <url> [options]",
		Long: `Executes an HTTP request. Supported HTTP methods are GET, HEAD, PUT, POST, PATCH and DELETE.
URL parameter is always required and must match something like "[https?://]host[:port][/path][?query]"

Protocol and host of the URL can be implicit if given like [host]:port/path...
Examples:
 * localhost/path	->	http://localhost/path
 * :1234/index		->	http://localhost:1234/index`,
	}

	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(*cobra.Command, []string) {
			fmt.Printf("http version: %s\n", version)
		},
	})

	// HTTP
	get := buildHTTPCommand(cfg, http.MethodGet, "", noConfigure)
	head := buildHTTPCommand(cfg, http.MethodHead, "", noConfigure)
	options := buildHTTPCommand(cfg, http.MethodOptions, "", noConfigure)
	post := buildHTTPCommand(cfg, http.MethodPost, `Make an HTTP POST request to the URL
This command requires the --body flag, which can be a string content or a file.`, bodyConfigure)
	put := buildHTTPCommand(cfg, http.MethodPut, `Make an HTTP PUT request to the URL
This command requires the --body flag, which can be a string content or a file.`, bodyConfigure)
	patch := buildHTTPCommand(cfg, http.MethodPatch, `Make an HTTP PATCH request to the URL
This command requires the --body flag, which can be a string content or a file.`, bodyConfigure)
	del := buildHTTPCommand(cfg, http.MethodDelete, "", noConfigure)
	root.AddCommand(get, head, options, post, put, patch, del)

	// URL alias
	alias := buildAlias(cfg)
	root.AddCommand(alias)

	// Persistant flags
	root.PersistentFlags().BoolP(verboseFlagName, "v", false, "Show logs.")
	return root
}

func buildHTTPClient(cmd *cobra.Command) (*http.Client, error) {
	timeout, _ := cmd.Flags().GetDuration(timeoutFlagName)

	var tlsConfig tls.Config
	certPub, _ := cmd.Flags().GetString(certpubFlagName)
	certKey, _ := cmd.Flags().GetString(certkeyFlagName)

	if certPub != "" && certKey == "" {
		return nil, errCertFlags
	} else if certPub == "" && certKey != "" {
		return nil, errCertFlags
	} else if certPub != "" && certKey != "" {
		cert, err := tls.LoadX509KeyPair(certPub, certKey)
		if err != nil {
			return nil, err
		}

		tlsConfig = tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	noFollowRedirects, _ := cmd.Flags().GetBool(noFollowRedirectsFlagName)
	var redirect checkRedirectFunc = nil
	if noFollowRedirects {
		fmt.Println("Setting custom redirect")
		redirect = func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return &http.Client{
		Timeout:       timeout,
		CheckRedirect: redirect,
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tlsConfig,
		},
	}, nil
}

// Returns a run function that handles a request for the given HTTP method
// and respects the config.
func buildRequestRun(method string, cfg *config) runFunc {
	return func(cmd *cobra.Command, args []string) {
		logger := logging.NewLogger()

		verbose, _ := cmd.Flags().GetBool(verboseFlagName)
		if verbose {
			logger.SetOutput(cfg.logs)
		} else {
			logger.SetOutput(ioutil.Discard)
		}

		httpClient, err := buildHTTPClient(cmd)
		checkInitError(err, cmd)

		traceLogger := logging.NewLogger()
		trace, _ := cmd.Flags().GetBool(traceFlagName)
		if trace {
			traceLogger.SetOutput(cfg.logs)
		} else {
			traceLogger.SetOutput(ioutil.Discard)
		}

		cl := client.NewClient(httpClient, logger, traceLogger)

		display, _ := cmd.Flags().GetString(displayFlagName)

		var formatter Formatter
		switch display {
		case "all":
			formatter, _ = NewDefaultFormatter(true, FormatComponents)
		case "", "none":
			formatter, _ = NewDefaultFormatter(true, []string{})
		default:
			components := strings.Split(strings.TrimSpace(display), ",")
			components = util.Map(components, strings.TrimSpace)
			formatter, err = NewDefaultFormatter(true, components)
			checkErr(err)
		}

		var signer client.RequestSigner
		signRequest, _ := cmd.Flags().GetBool(awsSigV4FlagName)
		if signRequest {
			logger.Print("Signing request using Sig V4")
			region, _ := cmd.Flags().GetString(awsRegionFlagName)
			creds := credentials.NewCredentials(&credentials.EnvProvider{})
			signer = client.NewAWSigner(v4.NewSigner(creds), region)
		} else {
			signer = client.DefaultSigner{}
		}

		fail, _ := cmd.Flags().GetBool(failFlagName)
		cfg.setFail(fail)
		repeat, _ := cmd.Flags().GetInt(repeatFlagName)
		cfg.setRepeat(repeat)

		aliasManager := newAliasLoader(cfg.aliasFilepath)

		handler := newHandler(
			cl,
			aliasManager,
			formatter,
			signer,
			logger,
			os.Exit,
			cfg,
		)

		url := args[0]
		bodyFlag, _ := cmd.Flags().GetString(bodyFlagName)
		err = handler.handleRequest(method, url, bodyFlag)
		if err != nil {
			fmt.Fprintf(cfg.errs, "error: %s\n", err)
			var execErr *execError
			if errors.As(err, &execErr) && execErr.showUsage {
				cmd.Usage()
			}
			os.Exit(1)
		}
	}
}

func buildHTTPCommand(
	cfg *config,
	method,
	long string,
	configure func(*cobra.Command),
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s <url>", strings.ToLower(method)),
		Short: fmt.Sprintf("HTTP %s request", strings.ToUpper(method)),
		Args:  cobra.ExactArgs(1),
		Run:   buildRequestRun(method, cfg),
	}

	addCommonFlags(cmd, cfg.headerOpt)
	configure(cmd)
	return cmd
}

func buildAlias(cfg *config) *cobra.Command {
	c := &cobra.Command{
		Use:   "alias [<name> <url>]",
		Short: "List, create or remove persistant URL aliases",
		Long: `List, create or remove persistant URL aliases.
Valid alias commands:
  - alias: list all aliases
  - alias name url: create a persistant alias
  - alias --remove name: remove alias by name

The name must match the pattern: ^[a-zA-Z_]\w*$, in other words
it must begin with _, a small or capital letter followed by zero
or more _, letters or numbers (max size of name is 20).`,
		Run: func(cmd *cobra.Command, args []string) {
			handler := AliasHandler{
				manager: newAliasLoader(cfg.aliasFilepath),
				infos:   cfg.infos,
				errors:  cfg.errs,
			}

			var err error
			switch len(args) {
			case 0:
				if r, _ := cmd.Flags().GetString("remove"); r != "" {
					err = handler.removeAlias(r)
				} else {
					err = handler.listAlias()
				}
			case 2:
				err = handler.setAlias(args[0], args[1])
			default:
				err = fmt.Errorf("unknown number of arguments")
			}

			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	c.Flags().StringP("remove", "r", "", "Remove alias with this name.")
	return c
}

func addCommonFlags(cmd *cobra.Command, h *HeaderOption) {
	cmd.Flags().VarP(h, headerFlagName, "H", `HTTP header, may be specified multiple times.
The value must conform to the format "name: value". "name" and "value" can
be separated by either a colon ":" or an equal sign "=", and the space
between is optional. Can be set in the same format using the env. variable
DEFAULT_HEADERS, where multiple headers must be separated by an |.`)

	cmd.Flags().IntP(repeatFlagName, "r", 1, "Repeat the request.")

	cmd.Flags().BoolP(
		awsSigV4FlagName,
		"4",
		false,
		`Use AWS signature V4 as authentication in the request. AWS region can be
set with the --aws-region option. Credentials are expected to be set
in environment variables.
`)
	cmd.Flags().String(
		awsRegionFlagName,
		defaultAWSRegion,
		"The AWS region to use in the AWS signature.")

	cmd.Flags().String(displayFlagName, "body", `Comma (,) separated list of response items to display.
Possible values:
  all:    all response information
  none:   no output
  status: response status code
  header: response headers
  body:   response body`)
	cmd.Flags().Bool(noColorFlagName, false, "Do not use colored output.")
	cmd.Flags().BoolP(failFlagName, "f", false, "Exit with status code > 0 if HTTP status is 400 or greater.")
	cmd.Flags().Bool(traceFlagName, false, "Output detailed TLS trace information.")
	cmd.Flags().DurationP(timeoutFlagName, "T", defaultTimeout, "Request timeout duration.")
	cmd.Flags().String(certpubFlagName, "", "Use as client certificate. Requires the --key flag.")
	cmd.Flags().String(certkeyFlagName, "", "Use as private key. Requires the --cert flag.")
	cmd.Flags().Bool(noFollowRedirectsFlagName, false, "Do not follow redirects. Default allows a maximum of 10 consecutive requests.")
}

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func checkInitError(err error, cmd *cobra.Command) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "error: %v\n\n", err)
	cmd.Usage()
	os.Exit(1)
}

func getAliasFilepath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homedir, ".gohttp", "aliases.json"), nil
}
