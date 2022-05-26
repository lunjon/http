package cli

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/lunjon/http/internal/client"
	"github.com/lunjon/http/internal/config"
	"github.com/lunjon/http/internal/format"
	"github.com/lunjon/http/internal/logging"
	"github.com/lunjon/http/internal/server"
	"github.com/lunjon/http/internal/util"
	"github.com/lunjon/http/tui"
	"github.com/spf13/cobra"
)

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

func build(
	version string,
	cfg cliConfig,
) *cobra.Command {
	root := &cobra.Command{
		Use:   "http",
		Short: "Starts an interactive session.",
		Long: `Sends HTTP requests - using either the TUI or CLI.

Supported HTTP methods are GET, HEAD, OPTIONS, PUT, POST, PATCH and DELETE.
URL parameter is always required and must be a valid URL.

Protocol and host of the URL can be implicit if given like [host]:port/path...
Examples:
 * localhost/path	->	http://localhost/path
 * :1234/index		->	http://localhost:1234/index
 * domain.com		->	https://domain.com

A request body can be specified in three ways:
 * stdin: pipe or IO redirection
 * --body '...': request body from a string
 * --body file: read content from a file`,
		Run: func(cmd *cobra.Command, args []string) {
			appConfig, err := cfg.getAppConfig()
			checkErr(err, cfg.errors)

			urls := []string{}
			for _, url := range appConfig.Aliases {
				urls = append(urls, url)
			}

			err = tui.Start(urls)
			checkErr(err, cfg.errors)
		},
	}

	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(*cobra.Command, []string) {
			styler := format.NewStyler()
			fmt.Printf("http version: %s\n", styler.WhiteB(version))
		},
	})

	// HTTP
	httpCommands := []struct {
		method string
		conf   func(*cobra.Command)
	}{
		{http.MethodGet, noConfigure},
		{http.MethodHead, noConfigure},
		{http.MethodOptions, noConfigure},
		{http.MethodPost, bodyConfigure},
		{http.MethodPut, bodyConfigure},
		{http.MethodPatch, bodyConfigure},
		{http.MethodDelete, noConfigure},
	}
	for _, cmd := range httpCommands {
		root.AddCommand(buildHTTPCommand(cfg, cmd.method, cmd.conf))
	}

	root.AddCommand(buildServer(cfg))
	root.AddCommand(buildConfig(cfg))

	// Persistant flags
	root.PersistentFlags().BoolP(verboseFlagName, "v", false, "Show logs.")
	return root
}

func buildHTTPClient(cfg config.Config, cmd *cobra.Command) (*http.Client, error) {
	flags := cmd.Flags()

	var tlsConfig tls.Config
	certPub, _ := flags.GetString(certpubFlagName)
	certKey, _ := flags.GetString(certkeyFlagName)

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

	noFollowRedirects, _ := flags.GetBool(noFollowRedirectsFlagName)
	var redirect checkRedirectFunc = nil
	if noFollowRedirects {
		redirect = func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return &http.Client{
		Timeout:       cfg.Timeout,
		CheckRedirect: redirect,
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tlsConfig,
		},
	}, nil
}

func updateConfig(cmd *cobra.Command, cfg config.Config) config.Config {
	flags := cmd.Flags()
	if flags.Changed(failFlagName) {
		v, _ := flags.GetBool(failFlagName)
		cfg = cfg.UseFail(v)
	}

	if flags.Changed(verboseFlagName) {
		v, _ := flags.GetBool(verboseFlagName)
		cfg = cfg.UseVerbose(v)
	}
	if flags.Changed(traceFlagName) {
		v, _ := flags.GetBool(traceFlagName)
		cfg = cfg.UseVerbose(v)
	}
	if flags.Changed(traceFlagName) {
		v, _ := flags.GetBool(traceFlagName)
		cfg = cfg.UseVerbose(v)
	}

	return cfg
}

// Returns a run function that handles a request for the given HTTP method
// and respects the config.
func buildRequestRun(
	method string,
	cfg cliConfig,
	headerOpt *HeaderOption,
) runFunc {
	return func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		appConfig, err := cfg.getAppConfig()
		checkErr(err, cfg.errors)

		appConfig = updateConfig(cmd, appConfig)

		logger := logging.New(io.Discard)
		if appConfig.Verbose {
			logger.SetOutput(cfg.logs)
		}

		traceLogger := logging.New(io.Discard)
		trace, _ := flags.GetBool(traceFlagName)
		if trace {
			traceLogger.SetOutput(cfg.logs)
		}

		httpClient, err := buildHTTPClient(appConfig, cmd)
		checkErr(err, cfg.errors)

		cl := client.NewClient(httpClient, logger, traceLogger)
		display, _ := flags.GetString(displayFlagName)

		var formatter format.ResponseFormatter
		switch display {
		case "all":
			formatter, _ = format.NewResponseFormatter(format.NewStyler(), format.ResponseComponents)
		case "", "none":
			formatter, _ = format.NewResponseFormatter(format.NewStyler(), []string{})
		default:
			components := strings.Split(strings.TrimSpace(display), ",")
			components = util.Map(components, strings.TrimSpace)
			formatter, err = format.NewResponseFormatter(format.NewStyler(), components)
			checkErr(err, cfg.errors)
		}

		var signer client.RequestSigner
		signRequest, _ := flags.GetBool(awsSigV4FlagName)
		if signRequest {
			logger.Print("Signing request using Sig V4")
			region, _ := flags.GetString(awsRegionFlagName)
			creds := credentials.NewCredentials(&credentials.EnvProvider{})
			signer = client.NewAWSigner(v4.NewSigner(creds), region)
		} else {
			signer = client.DefaultSigner{}
		}

		output := cfg.infos
		outputFile, _ := flags.GetString(outputFlagName)
		if outputFile != "" {
			file, err := os.Create(outputFile)
			checkErr(err, cfg.errors)

			defer func() {
				file.Close()
			}()
			output = file
		}

		failFunc := defaultFailFunc
		if appConfig.Fail {
			failFunc = os.Exit
		}

		handler := newRequestHandler(
			cl,
			formatter,
			signer,
			logger,
			appConfig,
			headerOpt.values,
			output,
			outputFile,
			failFunc,
		)

		url := args[0]
		bodyFlag, _ := flags.GetString(bodyFlagName)
		err = handler.handleRequest(method, url, bodyFlag)
		checkErr(err, cfg.errors)
	}
}

func buildHTTPCommand(
	cfg cliConfig,
	method string,
	configure func(*cobra.Command),
) *cobra.Command {
	headerOption := newHeaderOption()

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s <url>", strings.ToLower(method)),
		Short: fmt.Sprintf("HTTP %s request", strings.ToUpper(method)),
		Args:  cobra.ExactArgs(1),
		Run:   buildRequestRun(method, cfg, headerOption),
	}

	addCommonFlags(cmd, headerOption)
	configure(cmd)
	return cmd
}

func buildServer(cfg cliConfig) *cobra.Command {
	c := &cobra.Command{
		Use:   "server",
		Short: "Starts an HTTP server on localhost.",
		Long: `Starts an HTTP server on localhost.
Useful for local testing and debugging.`,
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetUint("port")

			// TODO: trap exit signal for graceful shutdown

			styler := format.NewStyler()
			fmt.Printf("Starting server on :%s.\n", styler.WhiteB(fmt.Sprint(port)))
			fmt.Printf("Press %s to exit.\n", styler.WhiteB("CTRL-C"))

			server := server.New(port, styler)
			err := server.Serve()
			checkErr(err, cfg.errors)
		},
	}

	c.Flags().UintP("port", "p", 8080, "Port to listen on.")
	return c
}

func buildConfig(cfg cliConfig) *cobra.Command {
	handler := newConfigHandler(cfg.configPath, cfg.infos)

	root := &cobra.Command{
		Use:   "config",
		Short: "Configuration commands.",
		Run: func(cmd *cobra.Command, _ []string) {
			err := handler.list()
			checkErr(err, cfg.errors)
		},
	}

	edit := &cobra.Command{
		Use:   "edit",
		Short: "Edit the configuration file.",
		Run: func(cmd *cobra.Command, _ []string) {
			err := handler.edit()
			checkErr(err, cfg.errors)
		},
	}

	init := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new configuration file.",
		Run: func(cmd *cobra.Command, _ []string) {
			err := handler.init()
			checkErr(err, cfg.errors)
		},
	}

	root.AddCommand(edit, init)
	root.Flags().BoolP(aliasHeadingFlagName, "n", false, "Do not display heading when listing aliases. Useful for e.g. scripting.")
	return root
}

func addCommonFlags(cmd *cobra.Command, h *HeaderOption) {
	cmd.Flags().VarP(h, headerFlagName, "H", `HTTP header, may be specified multiple times.
The value must conform to the format "name: value". "name" and "value" can
be separated by either a colon ":" or an equal sign "=", and the space
between is optional.`)

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
  none:       no output
  all:        all information
  status:     response status code text
  statuscode: response status code number
  headers:    response headers
  body:       response body
`)
	cmd.Flags().BoolP(failFlagName, "f", false, "Exit with status code > 0 if HTTP status is 400 or greater.")
	cmd.Flags().Bool(traceFlagName, false, "Output detailed TLS trace information.")
	cmd.Flags().DurationP(timeoutFlagName, "T", defaultTimeout, "Request timeout duration.")
	cmd.Flags().String(certpubFlagName, "", "Use as client certificate. Requires the --key flag.")
	cmd.Flags().String(certkeyFlagName, "", "Use as private key. Requires the --cert flag.")
	cmd.Flags().StringP(outputFlagName, "o", "", "Write output to file instead of stdout.")
	cmd.Flags().Bool(noFollowRedirectsFlagName, false, "Do not follow redirects. Default allows a maximum of 10 consecutive requests.")
}

func getAliasFilepath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homedir, ".gohttp", "aliases.json"), nil
}
