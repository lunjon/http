package command

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/lunjon/http/client"
	"github.com/lunjon/http/logging"
	"github.com/spf13/cobra"
)

type FailFunc func(status int)

type RunFunc func(*cobra.Command, []string)

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
	get := buildGet(cfg)
	head := buildHead(cfg)
	post := buildPost(cfg)
	put := buildPut(cfg)
	patch := buildPatch(cfg)
	del := buildDelete(cfg)
	root.AddCommand(get, head, post, put, patch, del)

	// URL alias
	alias := buildAlias(cfg)
	root.AddCommand(alias)

	// Persistant flags
	root.PersistentFlags().BoolP(verboseFlagName, "v", false, "Show logs.")
	return root
}

// Returns a run function that handles a request for the given HTTP method
// and respects the config.
func buildRequestRun(method string, cfg *config) RunFunc {
	return func(cmd *cobra.Command, args []string) {
		logger := logging.NewLogger()

		verbose, _ := cmd.Flags().GetBool(verboseFlagName)
		if verbose {
			logger.SetOutput(cfg.logs)
		} else {
			logger.SetOutput(ioutil.Discard)
		}

		timeout, _ := cmd.Flags().GetDuration(timeoutFlagName)

		var tlsConfig tls.Config
		certPub, _ := cmd.Flags().GetString(certpubFlagName)
		certKey, _ := cmd.Flags().GetString(certkeyFlagName)
		if certPub != "" && certKey == "" {
			checkInitError(errCertFlags, cmd)
		} else if certPub == "" && certKey != "" {
			checkInitError(errCertFlags, cmd)
		} else if certPub != "" && certKey != "" {
			cert, err := tls.LoadX509KeyPair(certPub, certKey)
			checkInitError(err, cmd)

			tlsConfig = tls.Config{
				Certificates: []tls.Certificate{cert},
			}
		}

		httpClient := &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tlsConfig,
			},
		}

		traceLogger := logging.NewLogger()
		trace, _ := cmd.Flags().GetBool(traceFlagName)
		if trace {
			traceLogger.SetOutput(cfg.logs)
		} else {
			traceLogger.SetOutput(ioutil.Discard)
		}

		cl := client.NewClient(httpClient, logger, traceLogger)

		brief, _ := cmd.Flags().GetBool(briefFlagName)
		silent, _ := cmd.Flags().GetBool(silentFlagName)

		if brief && silent {
			checkInitError(errBriefAndSilent, cmd)
		}

		var formatter Formatter = DefaultFormatter{}
		if silent {
			formatter = NullFormatter{}
		}
		if brief {
			formatter = BriefFormatter{}
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
		err := handler.handleRequest(method, url, bodyFlag)
		if err != nil {
			var execErr *execError
			if errors.As(err, &execErr) && execErr.showUsage {
				cmd.Usage()
			}
			os.Exit(1)
		}
	}
}

func buildGet(cfg *config) *cobra.Command {
	get := &cobra.Command{
		Use:   "get <url>",
		Short: "HTTP GET request",
		Args:  cobra.ExactArgs(1),
		Run:   buildRequestRun(http.MethodGet, cfg),
	}

	addCommonFlags(get, cfg.headerOpt)
	return get
}

func buildHead(cfg *config) *cobra.Command {
	head := &cobra.Command{
		Use:   "head <url>",
		Short: "HTTP HEAD request",
		Args:  cobra.ExactArgs(1),
		Run:   buildRequestRun(http.MethodHead, cfg),
	}

	addCommonFlags(head, cfg.headerOpt)
	return head
}

func buildPost(cfg *config) *cobra.Command {
	post := &cobra.Command{
		Use:   `post <url> --body <body>`,
		Short: "HTTP POST request",
		Long: `Make an HTTP POST request to the URL
This command requires the --body flag, which can be a string content or a file.`,
		Args: cobra.ExactArgs(1),
		Run:  buildRequestRun(http.MethodPost, cfg),
	}

	post.Flags().StringP(bodyFlagName, "B", "", "Request body to use. Can be string content or a filename.")
	addCommonFlags(post, cfg.headerOpt)
	return post
}

func buildPatch(cfg *config) *cobra.Command {
	patch := &cobra.Command{
		Use:   `patch <url> --body <body>`,
		Short: "HTTP PATCH request",
		Long: `Make an HTTP PATCH request to the URL
This command requires the --body flag, which can be a string content or a file.`,
		Args: cobra.ExactArgs(1),
		Run:  buildRequestRun(http.MethodPatch, cfg),
	}

	patch.Flags().String("body", "", "Request body to use. Can be string content or a filename.")
	addCommonFlags(patch, cfg.headerOpt)
	return patch
}

func buildPut(cfg *config) *cobra.Command {
	put := &cobra.Command{
		Use:   `put <url> --body <body>`,
		Short: "HTTP PUT request",
		Long: `Make an HTTP PUT request to the URL
This command requires the --body flag, which can be a string content or a file.`,
		Args: cobra.ExactArgs(1),
		Run:  buildRequestRun(http.MethodPut, cfg),
	}

	put.Flags().String("body", "", "Request body to use. Can be string content or a filename.")
	addCommonFlags(put, cfg.headerOpt)
	return put
}

func buildDelete(cfg *config) *cobra.Command {
	del := &cobra.Command{
		Use:   `delete <url>`,
		Short: "HTTP DELETE request",
		Args:  cobra.ExactArgs(1),
		Run:   buildRequestRun(http.MethodDelete, cfg),
	}

	addCommonFlags(del, cfg.headerOpt)
	return del
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

	cmd.Flags().BoolP(silentFlagName, "s", false, "Suppress output of response body.")
	cmd.Flags().BoolP(failFlagName, "f", false, "Exit with status code > 0 if HTTP status is 400 or greater.")
	cmd.Flags().Bool(briefFlagName, false, "Output brief summary of the request.")
	cmd.Flags().Bool(traceFlagName, false, "Output detailed TLS trace information.")
	cmd.Flags().DurationP(timeoutFlagName, "T", defaultTimeout, "Request timeout duration.")
	cmd.Flags().String(certpubFlagName, "", "Use as client certificate public key  (requires --cert-key-file flag).")
	cmd.Flags().String(certkeyFlagName, "", "Use as client certificate private key (requires --cert-pub-file flag).")
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
