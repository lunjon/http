package cli

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/lunjon/http/cli/options"
	"github.com/lunjon/http/internal/client"
	"github.com/lunjon/http/internal/config"
	"github.com/lunjon/http/internal/history"
	"github.com/lunjon/http/internal/logging"
	"github.com/lunjon/http/internal/server"
	"github.com/lunjon/http/internal/style"
	"github.com/spf13/cobra"
)

const (
	verbGroupID = "verbs"
)

var (
	noConfigure   = func(*cobra.Command) {}
	bodyConfigure = func(cmd *cobra.Command) {
		flags := cmd.Flags()
		flags.String(
			options.DataStringFlagName,
			"",
			"Use string as request body.",
		)
		flags.String(
			options.DataFileFlagName,
			"",
			"Read request body from file.",
		)
		cmd.MarkFlagFilename(options.DataFileFlagName)
		flags.Bool(
			options.DataStdinFlagName,
			false,
			"Read request body from stdin.",
		)
		flags.StringArray(
			options.DataURLEncodeFlagName,
			[]string{},
			`URL encoded body. Can be called multiple times (per value).
Should be specified in format "key=value".`,
		)
		cmd.MarkFlagsMutuallyExclusive(
			options.DataStringFlagName,
			options.DataFileFlagName,
			options.DataStdinFlagName,
			options.DataURLEncodeFlagName,
		)

	}
)

func build(
	version string,
	cfg cliConfig,
) *cobra.Command {
	root := &cobra.Command{
		Use:   "http",
		Short: `http - send HTTP requests from your command-line`,
		Long: `http - send HTTP requests from your command-line

Supported HTTP methods are GET, HEAD, OPTIONS, PUT, POST, PATCH and DELETE.

Protocol and host of the URL can be implicit if given like [host]:port/path...
Examples:
 * localhost/path	->	http://localhost/path
 * :1234/index		->	http://localhost:1234/index
 * domain.com		->	https://domain.com
`,
	}

	root.AddGroup(&cobra.Group{
		ID:    verbGroupID,
		Title: "HTTP Commands:",
	})

	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(*cobra.Command, []string) {
			fmt.Printf("http version: %s\n", style.Bold.Render(version))
		},
	})

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

	root.AddCommand(buildHistory(cfg))
	root.AddCommand(buildServe(cfg))
	root.AddCommand(buildConfig(cfg))

	// Persistant flags
	root.PersistentFlags().BoolP(options.VerboseFlagName, "v", false, "Show logs.")

	root.Flags().SortFlags = true
	return root
}

func updateConfig(cmd *cobra.Command, cfg config.Config) config.Config {
	flags := cmd.Flags()
	if flags.Changed(options.FailFlagName) {
		v, _ := flags.GetBool(options.FailFlagName)
		cfg = cfg.UseFail(v)
	}

	if flags.Changed(options.VerboseFlagName) {
		v, _ := flags.GetBool(options.VerboseFlagName)
		cfg = cfg.UseVerbose(v)
	}
	if flags.Changed(options.TLSTraceFlagName) {
		v, _ := flags.GetBool(options.TLSTraceFlagName)
		cfg = cfg.UseVerbose(v)
	}
	if flags.Changed(options.TimeoutFlagName) {
		v, _ := flags.GetDuration(options.TimeoutFlagName)
		cfg = cfg.UseTimeout(v)
	}

	return cfg
}

// Returns a function that handles a request for the given HTTP method
// and respects the config.
func buildRequestRun(
	method string,
	cfg cliConfig,
	headerOpt *options.HeaderOption,
	tlsMinVersion *options.TLSVersionOption,
	tlsMaxVersion *options.TLSVersionOption,
	certFile *options.FileOption,
	keyFile *options.FileOption,
	certKind *options.CertKindOption,
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

		// HTTP CLIENT
		traceLogger := logging.New(io.Discard)
		if flags.Changed(options.TLSTraceFlagName) {
			traceLogger.SetOutput(cfg.logs)
		}

		settings := client.NewSettings().
			WithTimeout(appConfig.Timeout).
			WithNoFollowRedirects(flags.Changed(options.NoFollowRedirectsFlagName))

		tlsOpts := client.NewTLSOptions().
			WithVersions(tlsMinVersion.Value(), tlsMaxVersion.Value())

		certFile, certFileSet := certFile.Value()
		if certFileSet {
			certKind.Update(certFile)
			keyFile, keyFileSet := keyFile.Value()
			switch certKind.Value() {
			case options.CertKindX509:
				if !keyFileSet {
					checkErr(fmt.Errorf("%s option required but not set", options.CertkeyFlagName), cfg.errors)
				}
				tlsOpts = tlsOpts.WithX509Cert(certFile, keyFile)
			case options.CertKindPKCS12:
				certPass, _ := flags.GetString(options.CertPassFlagName)
				if keyFileSet {
					checkErr(fmt.Errorf("%s option should not be specified with type %s", options.CertkeyFlagName, certKind.Value()), cfg.errors)
				}
				pfx, err := os.ReadFile(certFile)
				checkErr(err, cfg.errors)
				tlsOpts = tlsOpts.WithPKCS12Cert(pfx, certPass)
			}
		}

		settings = settings.WithTLSOptions(tlsOpts)
		cl, err := client.NewClient(settings, logger, traceLogger)
		checkErr(err, cfg.errors)

		// OUTPUT
		outputFormat, _ := flags.GetString(options.FormatFlagName)
		formatter, err := FormatterFromString(Format(outputFormat))
		checkErr(err, cfg.errors)

		var signer client.RequestSigner
		signRequest, _ := flags.GetBool(options.AWSSigV4FlagName)
		if signRequest {
			logger.Print("Signing request using Sig V4")
			region, _ := flags.GetString(options.AWSRegionFlagName)
			creds := credentials.NewCredentials(&credentials.EnvProvider{})
			signer = client.NewAWSigner(v4.NewSigner(creds), region)
		} else {
			signer = client.DefaultSigner{}
		}

		output := cfg.infos
		outputFile, _ := flags.GetString(options.OutfileFlagName)
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

		header := headerOpt.Header()
		if bearerToken, _ := flags.GetString(options.BearerFlagName); bearerToken != "" {
			header.Set("Authorization", fmt.Sprintf("Bearer %s", strings.TrimSpace(bearerToken)))
		}

		handler := newRequestHandler(
			cl,
			formatter,
			signer,
			history.NewHandler(cfg.historyPath),
			logger,
			appConfig,
			header,
			output,
			outputFile,
			failFunc,
		)

		url := args[0]
		dataOpts, err := options.DataOptionsFromFlags(cmd)
		checkErr(err, cfg.errors)

		err = handler.handleRequest(method, url, dataOpts)
		checkErr(err, cfg.errors)
	}
}

func buildHTTPCommand(
	cfg cliConfig,
	method string,
	configure func(*cobra.Command),
) *cobra.Command {
	headerOption := options.NewHeaderOption()
	tlsMinVersion := options.NewTLSVersionOption(tls.VersionTLS12)
	tlsMaxVersion := options.NewTLSVersionOption(tls.VersionTLS13)
	certFile := &options.FileOption{}
	keyFile := &options.FileOption{}
	certKind := &options.CertKindOption{}

	cmd := &cobra.Command{
		GroupID: verbGroupID,
		Use:     fmt.Sprintf("%s <url>", strings.ToLower(method)),
		Short:   fmt.Sprintf("HTTP %s request", strings.ToUpper(method)),
		Args:    cobra.ExactArgs(1),
		Run: buildRequestRun(
			method,
			cfg,
			headerOption,
			tlsMinVersion,
			tlsMaxVersion,
			certFile,
			keyFile,
			certKind,
		),
	}

	addCommonFlags(cmd, headerOption, tlsMinVersion, tlsMaxVersion, certFile, keyFile, certKind)
	configure(cmd)
	return cmd
}

func buildHistory(cfg cliConfig) *cobra.Command {
	hst := &cobra.Command{
		Use:     "history",
		Aliases: []string{"hist"},
		Short:   "Command for managing request history",
		Run: func(cmd *cobra.Command, args []string) {
			handler := history.NewHandler(cfg.historyPath)
			entries, err := handler.GetAll()
			checkErr(err, cfg.errors)

			for _, entry := range entries {
				fmt.Fprintf(cfg.infos, "%s %s\n", entry.Method, entry.URL)
			}
		},
	}

	clear := &cobra.Command{
		Use:   "clear",
		Short: "Clears request history",
		Run: func(cmd *cobra.Command, args []string) {
			handler := history.NewHandler(cfg.historyPath)
			err := handler.Clear()
			checkErr(err, cfg.errors)
		},
	}

	hst.AddCommand(clear)
	return hst
}

func buildServe(cfg cliConfig) *cobra.Command {
	port := options.NewPortOption()
	statusFlagName := "show-status"
	listFlagName := "list"
	staticFlagName := "static"

	c := &cobra.Command{
		Use:   "serve",
		Short: "Starts an HTTP server on localhost",
		Long: `Starts an HTTP server on localhost. Useful for local testing.

If --static DIR is provided the server hosts the files
in that directory. Otherwise the default API is started (use --list to see endoints).`,
		Run: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()

			listRoutes, _ := flags.GetBool(listFlagName)
			if listRoutes {
				server.ListRoutes()
				return
			}

			showStatus, _ := flags.GetBool(statusFlagName)
			staticRoot, _ := flags.GetString(staticFlagName)

			opts := server.Options{
				Port:       port.Value(),
				ShowStatus: showStatus,
				StaticRoot: staticRoot,
			}

			server := server.New(opts)

			errs := make(chan error)
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

			go func() {
				err := server.Serve()
				if err != nil {
					errs <- err
				}
			}()

			select {
			case <-sig:
				server.Close()
			case err := <-errs:
				server.Close()
				checkErr(err, cfg.errors)
			}
		},
	}

	c.Flags().VarP(port, "port", "p", "Port to listen on. Must be valid number in the range 1024-65535.")
	c.Flags().String(staticFlagName, "", "Serve static files from this directory.")
	c.Flags().Bool(statusFlagName, false, "Shows current status instead of showing each request.")
	c.Flags().Bool(listFlagName, false, "List predefined routes.")
	return c
}

func buildConfig(cfg cliConfig) *cobra.Command {
	handler := newConfigHandler(cfg.configPath, cfg.infos)

	root := &cobra.Command{
		Use:   "config",
		Short: "Configuration commands",
		Run: func(cmd *cobra.Command, _ []string) {
			err := handler.list()
			checkErr(err, cfg.errors)
		},
	}

	edit := &cobra.Command{
		Use:   "edit",
		Short: "Edit the configuration file.",
		Run: func(cmd *cobra.Command, _ []string) {
			editor, _ := cmd.Flags().GetString("editor")
			err := handler.edit(editor)
			checkErr(err, cfg.errors)
		},
	}
	editor, ok := os.LookupEnv("EDITOR")
	if !ok {
		editor = "vim"
	}
	edit.Flags().StringP("editor", "e", editor, "Use as editor.")

	init := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new configuration file.",
		Run: func(cmd *cobra.Command, _ []string) {
			err := handler.init()
			checkErr(err, cfg.errors)
		},
	}

	root.AddCommand(edit, init)
	root.Flags().BoolP(options.AliasHeadingFlagName, "n", false, "Do not display heading when listing aliases. Useful for e.g. scripting.")
	return root
}

func addCommonFlags(
	cmd *cobra.Command,
	header *options.HeaderOption,
	tlsMinVersion *options.TLSVersionOption,
	tlsMaxVersion *options.TLSVersionOption,
	certFile *options.FileOption,
	keyFile *options.FileOption,
	certKind *options.CertKindOption,
) {
	cmd.Flags().VarP(header, options.HeaderFlagName, "H", `HTTP header, may be specified multiple times.
The value must conform to the format "name: value".`)

	cmd.Flags().BoolP(
		options.AWSSigV4FlagName,
		"4",
		false,
		`Use AWS signature V4 as authentication in the request. AWS region can be
set with the --aws-region option. Credentials are expected to be set
in environment variables.
`)
	flags := cmd.Flags()
	flags.String(
		options.AWSRegionFlagName,
		defaultAWSRegion,
		"The AWS region to use in the AWS signature.")

	flags.String(options.BearerFlagName, "", "Set Authorization header as OAuth2 bearer token.")
	flags.String(options.FormatFlagName, "text", `Output format of response. Possible values: text, json.`)
	flags.BoolP(options.FailFlagName, "f", false, "Exit with status code > 0 if HTTP status is 400 or greater.")
	flags.DurationP(options.TimeoutFlagName, "T", defaultTimeout, "Request timeout duration.")
	flags.StringP(options.OutfileFlagName, "o", "", "Write output to file instead of stdout.")
	flags.Bool(options.NoFollowRedirectsFlagName, false, "Do not follow redirects. Default allows a maximum of 10 consecutive requests.")

	flags.Var(certFile, options.CertfileFlagName, "Use as client certificate. Requires the --key flag.")
	cmd.MarkFlagFilename(options.CertfileFlagName)
	flags.Var(keyFile, options.CertkeyFlagName, "Use as private key. Requires the --cert flag.")
	cmd.MarkFlagFilename(options.CertkeyFlagName)
	flags.Var(certKind, options.CertKindFlagName, "Specifies certificate type.")
	flags.String(options.CertPassFlagName, "", "Use as password for certificate.")

	flags.Bool(options.TLSTraceFlagName, false, "Output detailed TLS trace information.")
	flags.Var(tlsMinVersion, options.TLSMinVersionFlagName, "Set minimum TLS version to use. Allowed values are 1.0-3.")
	flags.Var(tlsMaxVersion, options.TLSMaxVersionFlagName, "Set maximum TLS version to use. Allowed values are 1.0-3.")
}
