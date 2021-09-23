package command

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/lunjon/http/client"
	"github.com/lunjon/http/logging"
	"github.com/spf13/cobra"
)

const (
	defaultTimeout = time.Second * 30
	description    = `Executes an HTTP request. Supported HTTP methods are GET, HEAD, PUT, POST, PATCH and DELETE.
URL parameter is always required and must match something like "[https?://]host[:port][/path][?query]"

Protocol and host of the URL can be implicit if given like [host]:port/path...
Examples:
 * localhost/path	->	http://localhost/path
 * :1234/index		->	http://localhost:1234/index
`
	defaultAWSRegion  = "eu-west-1"
	defaultHeadersEnv = "DEFAULT_HEADERS"
)

func newDefaultHandler() *Handler {
	logger := logging.NewLogger()
	traceLogger := logging.NewLogger()

	cl := client.NewClient(&http.Client{}, logger, traceLogger)

	homedir, err := os.UserHomeDir()
	checkErr(err)
	dir := path.Join(homedir, ".gohttp")

	return NewHandler(
		cl,
		logger,
		traceLogger,
		os.Stdout,
		os.Stderr,
		dir,
		func() {
			os.Exit(1)
		},
	)
}

func Build(version string) *cobra.Command {
	h := newDefaultHandler()
	return build(version, h)
}

// Build the root command for http and set version.
func build(version string, handler *Handler) *cobra.Command {
	root := buildRoot(handler)
	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(*cobra.Command, []string) {
			fmt.Printf("http version: %s\n", version)
		},
	})

	// HTTP
	get := buildGet(handler)
	head := buildHead(handler)
	post := buildPost(handler)
	put := buildPut(handler)
	patch := buildPatch(handler)
	del := buildDelete(handler)
	root.AddCommand(get, head, post, put, patch, del)

	// URL alias
	alias := buildAlias(handler)
	root.AddCommand(alias)

	// Persistant flags
	root.PersistentFlags().BoolP(verboseFlagName, "v", false, "Show logs.")

	return root
}

func buildRoot(handler *Handler) *cobra.Command {
	root := &cobra.Command{
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			handler.init(cmd)
		},
		Use:   "http",
		Short: "http <method> <url> [options]",
		Long:  description,
	}
	return root
}

func buildGet(handler *Handler) *cobra.Command {
	get := &cobra.Command{
		Use:   "get <url>",
		Short: "HTTP GET request",
		Args:  cobra.ExactArgs(1),
		Run:   handler.Get,
	}

	addCommonFlags(get, handler)
	return get
}

func buildHead(handler *Handler) *cobra.Command {
	head := &cobra.Command{
		Use:   "head <url>",
		Short: "HTTP HEAD request",
		Args:  cobra.ExactArgs(1),
		Run:   handler.Head,
	}

	addCommonFlags(head, handler)
	return head
}

func buildPost(handler *Handler) *cobra.Command {
	post := &cobra.Command{
		Use:   `post <url> --body <body>`,
		Short: "HTTP POST request",
		Long: `Make an HTTP POST request to the URL
This command requires the --body flag, which can be a string content or a file.`,
		Args: cobra.ExactArgs(1),
		Run:  handler.Post,
	}

	post.Flags().StringP(bodyFlagName, "B", "", "Request body to use. Can be string content or a filename.")
	addCommonFlags(post, handler)
	return post
}

func buildPatch(handler *Handler) *cobra.Command {
	patch := &cobra.Command{
		Use:   `patch <url> --body <body>`,
		Short: "HTTP PATCH request",
		Long: `Make an HTTP PATCH request to the URL
This command requires the --body flag, which can be a string content or a file.`,
		Args: cobra.ExactArgs(1),
		Run:  handler.Patch,
	}

	patch.Flags().String("body", "", "Request body to use. Can be string content or a filename.")
	addCommonFlags(patch, handler)
	return patch
}

func buildPut(handler *Handler) *cobra.Command {
	put := &cobra.Command{
		Use:   `put <url> --body <body>`,
		Short: "HTTP PUT request",
		Long: `Make an HTTP PUT request to the URL
This command requires the --body flag, which can be a string content or a file.`,
		Args: cobra.ExactArgs(1),
		Run:  handler.Put,
	}

	put.Flags().String("body", "", "Request body to use. Can be string content or a filename.")
	addCommonFlags(put, handler)
	return put
}

func buildDelete(handler *Handler) *cobra.Command {
	del := &cobra.Command{
		Use:   `delete <url>`,
		Short: "HTTP DELETE request",
		Args:  cobra.ExactArgs(1),
		Run:   handler.Delete,
	}

	addCommonFlags(del, handler)
	return del
}

func buildAlias(handler *Handler) *cobra.Command {
	return &cobra.Command{
		Use:   "alias <name> <url>",
		Short: "List and create persistant URL aliases",
		Run:   handler.handleAlias,
	}
}

func addCommonFlags(cmd *cobra.Command, handler *Handler) {
	cmd.Flags().VarP(handler.header, headerFlagName, "H", `HTTP header, may be specified multiple times.
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
