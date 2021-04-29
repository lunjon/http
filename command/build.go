package command

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/lunjon/http/logging"
	"github.com/lunjon/http/rest"
	"github.com/spf13/cobra"
)

const (
	defaultTimeout = time.Second * 10
	description    = `Executes an HTTP request. Supported HTTP methods are GET, HEAD, PUT, POST, PATCH and DELETE.
URL parameter is always required and must match something like "[https?://]host[:port][/path][?query]"

Protocol and host of the URL can be implicit if given like [host]:port/path...
Examples:
 * localhost/path	->	http://localhost/path
 * :1234/index		->	http://localhost:1234/index
`
	DefaultAWSRegion  = "eu-west-1"
	DefaultHeadersEnv = "DEFAULT_HEADERS"
)

func createHandler() *Handler {
	logger := logging.NewLogger()
	h := NewHeaderOption()
	httpClient := &http.Client{}
	restClient := rest.NewClient(httpClient, logger)
	handler := NewHandler(restClient, logger, h)
	return handler
}

// Build the root command for http and set version.
func Build(version string) *cobra.Command {
	handler := createHandler()

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

	// Command for generating completion
	comp := buildComp(root)
	root.AddCommand(comp)

	// Persistant flags
	root.PersistentFlags().BoolP(VerboseFlagName, "v", false, "Show logs.")
	root.PersistentFlags().DurationP(TimeoutFlagName, "T", defaultTimeout, "Request timeout duration.")

	return root
}

func buildRoot(handler *Handler) *cobra.Command {
	root := &cobra.Command{
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool(VerboseFlagName)
			handler.Verbose(verbose)

			timeout, _ := cmd.Flags().GetDuration(TimeoutFlagName)
			handler.Timeout(timeout)
		},
		Use:   "http",
		Short: "http <method> <url> [options]",
		Long:  description,
	}
	return root
}

func buildComp(root *cobra.Command) *cobra.Command {
	filenameDefault := ""
	gen := &cobra.Command{
		Use:   "comp <type>",
		Short: "Generate completion for a supported shell: bash, zsh or fish",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			shell := args[0]
			filename, _ := cmd.Flags().GetString("filename")
			if filename == filenameDefault {
				filename = shell
			}

			var err error
			switch shell {
			case "bash":
				err = root.GenBashCompletionFile(filename)
			case "zsh":
				err = root.GenZshCompletionFile(filename)
			case "fish":
				err = root.GenFishCompletionFile(filename, false)
			default:
				fmt.Fprintf(os.Stderr, "invalid shell type: %s", shell)
				os.Exit(1)
			}

			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to generate %s completion file: %v", shell, err)
				os.Exit(1)
			}
		},
	}

	gen.Flags().StringP("filename", "f", filenameDefault, "Output file name.")
	return gen
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

	post.Flags().StringP(BodyFlagName, "B", "", "Request body to use. Can be string content or a filename.")
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
func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Error: failed to read alias file: %s\n", err)
	os.Exit(1)

}

func buildAlias(handler *Handler) *cobra.Command {
	return &cobra.Command{
		Use:   "alias <name> <url>",
		Short: "Create a persistant URL alias",
		Run:   handler.handleAlias,
	}
}

func addCommonFlags(cmd *cobra.Command, handler *Handler) {
	cmd.Flags().VarP(handler.header, HeaderFlagName, "H", `HTTP header, may be specified multiple times.
The value must conform to the format "name: value". "name" and "value" can
be separated by either a colon ":" or an equal sign "=", and the space
between is optional. Can be set in the same format using the env. variable
DEFAULT_HEADERS, where multiple headers must be separated by an |.`)

	cmd.Flags().IntP(RepeatFlagName, "r", 1, "Repeat the request.")

	cmd.Flags().BoolP(
		AWSSigV4FlagName,
		"4",
		false,
		"Use AWS signature V4 as authentication in the request. Requires the --aws-region option.")
	cmd.Flags().String(
		AWSRegionFlagName,
		DefaultAWSRegion,
		"The AWS region to use in the AWS signature.")
	cmd.Flags().String(
		AWSProfileFlagName,
		"",
		"The name of an AWS profile in your AWS configuration. If not specified, environment variables are used.")

	// Silent mode
	cmd.Flags().BoolP(SilentFlagName, "s", false, "Suppress output of response body.")
}
