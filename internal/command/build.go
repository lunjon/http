package command

import (
	"net/http"
	"time"

	"github.com/lunjon/httpreq/internal/constants"
	"github.com/lunjon/httpreq/internal/logging"
	"github.com/lunjon/httpreq/internal/rest"
	"github.com/spf13/cobra"
)

const (
	defaultTimeout = time.Second * 10
)

func createHandler() *Handler {
	// Create handler and it's dependencies
	logger := logging.NewLogger()
	h := NewHeaderOption()
	httpClient := &http.Client{}
	restClient := rest.NewClient(httpClient, logger)
	handler := NewHandler(restClient, logger, h)
	return handler
}

// Build the root command for httpreq.
func Build() *cobra.Command {
	handler := createHandler()

	// HTTP
	get := buildGet(handler)
	head := buildHead(handler)
	post := buildPost(handler)
	put := buildPut(handler)
	patch := buildPatch(handler)
	delete := buildDelete(handler)

	root := &cobra.Command{
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Verbose flag
			verbose, _ := cmd.Flags().GetBool(constants.VerboseFlagName)
			handler.Verbose(verbose)

			// Timeout flag
			timeout, _ := cmd.Flags().GetDuration(constants.TimeoutFlagName)
			handler.Timeout(timeout)
		},
		Use:   "httpreq",
		Short: "httpreq <method> <url> [options]",
		Long: `Execute an HTTP request. Supported HTTP methods are GET, HEAD, PUT, POST, PATCH and DELETE.
URL parameter must be a valid HTTP URL; i.e. it must match something like "https?://host[:port][/path][?query]"`,
	}

	// Persistant flags
	root.PersistentFlags().BoolP(constants.VerboseFlagName, "v", false, "Shows debug logs.")
	root.PersistentFlags().DurationP(
		constants.TimeoutFlagName,
		"T",
		defaultTimeout,
		"Request timeout in seconds.")
	root.AddCommand(get, head, post, put, patch, delete)
	return root
}

func buildGet(handler *Handler) *cobra.Command {
	get := &cobra.Command{
		Use:     "get <url>",
		Aliases: []string{"g"},
		Short:   "HTTP GET request.",
		Args:    cobra.ExactArgs(1),
		Run:     handler.Get,
	}

	addCommonFlags(get, handler)
	return get
}

func buildHead(handler *Handler) *cobra.Command {
	head := &cobra.Command{
		Use:   "head <url>",
		Aliases: []string{"h", "hd"},
		Short: "HTTP HEAD request.",
		Args:  cobra.ExactArgs(1),
		Run:   handler.Head,
	}

	addCommonFlags(head, handler)
	return head
}

func buildPost(handler *Handler) *cobra.Command {
	post := &cobra.Command{
		Use:     `post <url> --body <body>`,
		Aliases: []string{"po"},
		Short:   "HTTP POST request with a body.",
		Long: `Make an HTTP POST request to the URL with a body.
This command requires the --body flag, which can be a string content or a file.`,
		Args: cobra.ExactArgs(1),
		Run:  handler.Post,
	}

	post.Flags().String("body", "", "Request body to use. Can be string content or a filename.")
	addCommonFlags(post, handler)
	return post
}

func buildPatch(handler *Handler) *cobra.Command {
	patch := &cobra.Command{
		Use:   `patch <url> --body <body>`,
		Aliases: []string{"pa"},
		Short: "HTTP PATCH request with a body.",
		Long: `Make an HTTP PATCH request to the URL with a body.
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
		Aliases: []string{"pu"},
		Short: "HTTP PUT request with a body.",
		Long: `Make an HTTP PUT request to the URL with a body.
This command requires the --body flag, which can be a string content or a file.`,
		Args: cobra.ExactArgs(1),
		Run:  handler.Put,
	}

	put.Flags().String("body", "", "Request body to use. Can be string content or a filename.")
	addCommonFlags(put, handler)
	return put
}

func buildDelete(handler *Handler) *cobra.Command {
	delete := &cobra.Command{
		Use:     `delete <url>`,
		Aliases: []string{"d", "de", "del"},
		Short:   "HTTP DELETE request.",
		Args:    cobra.ExactArgs(1),
		Run:     handler.Delete,
	}

	addCommonFlags(delete, handler)
	return delete
}

func addCommonFlags(cmd *cobra.Command, handler *Handler) {
	// Headers
	cmd.Flags().VarP(handler.header, constants.HeaderFlagName, "H", "")

	// AWS signature V4 flags
	cmd.Flags().BoolP(
		constants.AWSSigV4FlagName,
		"4",
		false,
		"Use AWS signature V4 as authentication in the request. Requires the --aws-region option.")

	cmd.Flags().String(
		constants.AWSRegionFlagName,
		constants.DefaultAWSRegion,
		"The AWS region to use in the AWS signature.")

	cmd.Flags().String(
		constants.AWSProfileFlagName,
		"",
		"The name of an AWS profile in your AWS configuration. If not specified, environment variables are used.")

	// Silent mode
	cmd.Flags().BoolP(constants.SilentFlagName, "s", false, "Suppress output of response body.")
}
