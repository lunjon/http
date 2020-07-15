package command

import (
	"fmt"
	"net/http"
	"time"

	"github.com/lunjon/httpreq/internal/constants"
	"github.com/lunjon/httpreq/internal/logging"
	"github.com/lunjon/httpreq/internal/parse"
	"github.com/lunjon/httpreq/internal/rest"
	"github.com/spf13/cobra"
)

const(
	defaultTimeout = time.Second * 10
)

func createHandler() *Handler {
	// Create handler and it's dependencies
	logger := logging.NewLogger()
	h := NewHeader()
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
	post := buildPost(handler)
	delete := buildDelete(handler)

	parse := buildParseURL()

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
		Short: "httpreq <method> <route> [options]",
		Long: `Execute an HTTP request. Supported HTTP methods are GET, POST and DELETE.

Routes can have any of the following formats:
  * http[s]://host[:port]/path 		(use as is)
  * host.domain.example[:port]/path	(add https:// as protocol)
  * :port/path 				(assume http://localhost:port/path)
  * /path				(assume http://localhost:80/path`,
	}

	// Persistant flags
	root.PersistentFlags().BoolP(constants.VerboseFlagName, "v", false, "Shows debug logs.")
	root.PersistentFlags().DurationP(
		constants.TimeoutFlagName,
		"T",
		defaultTimeout,
		"Request timeout in seconds.")
	root.AddCommand(get, post, delete, parse)
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

func buildPost(handler *Handler) *cobra.Command {
	post := &cobra.Command{
		Use:     `post <url> --body <body>`,
		Aliases: []string{"p"},
		Short:   "HTTP POST request with a JSON body.",
		Long: `Make an HTTP POST request to the URL with a JSON body.
This command requires the --body flag, which can be a string content or a file.`,
		Args: cobra.ExactArgs(1),
		Run:  handler.Post,
	}

	post.Flags().String("body", "", "Request body to use. Can be string content or a filename.")
	addCommonFlags(post, handler)
	return post
}

func buildDelete(handler *Handler) *cobra.Command {
	delete := &cobra.Command{
		Use:     `delete <url>`,
		Aliases: []string{"d"},
		Short:   "HTTP DELETE request.",
		Args:    cobra.ExactArgs(1),
		Run:     handler.Delete,
	}

	addCommonFlags(delete, handler)
	return delete
}

func buildParseURL() *cobra.Command {
	parse := &cobra.Command{
		Use:     `parse-url <url>`,
		Aliases: []string{"url"},
		Short:   "Parse the URL and print the results.",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url, err := parse.ParseURL(args[0])
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			detailed, _ := cmd.Flags().GetBool(constants.DetailsFlagName)
			if detailed {
				fmt.Println(url.DetailString())
			} else {
				fmt.Println(url.String())
			}

		},
	}

	parse.Flags().BoolP(
		constants.DetailsFlagName,
		"d",
		false,
		"Whether to output detailed information.",
	)

	return parse
}

func addCommonFlags(cmd *cobra.Command, handler *Handler) {
	// Headers
	cmd.Flags().VarP(handler.header, constants.HeaderFlagName, "H", "")

	cmd.Flags().StringP(constants.OutputFileFlagName, "o", "", "Output the response body to the filename.")

	cmd.Flags().BoolP(
		constants.ResponseBodyOnlyFlagName,
		"R",
		false,
		"Output only the response body")

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

	// Sandbox
	cmd.Flags().Bool(constants.SandboxFlagName, false, "Run the request to a sandbox server.")
}
