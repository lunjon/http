package command

import (
	"fmt"

	"github.com/lunjon/httpreq/internal/constants"
	"github.com/lunjon/httpreq/internal/rest"
	"github.com/lunjon/httpreq/pkg/parse"
	"github.com/spf13/cobra"
)

// Build the root command for httpreq.
func Build() *cobra.Command {
	get := buildGet()
	post := buildPost()
	delete := buildDelete()
	run := buildRun()
	sandbox := buildSandbox()
	parse := buildParseURL()

	root := &cobra.Command{
		Use:   "httpreq",
		Short: "httpreq <method> <route> [options]",
		Long: `Execute an HTTP request. Supported HTTP methods are GET, POST and DELETE.

Routes can have any of the following formats:
  * http[s]://host[:port]/path 	(use as is)
  * :port/path 			(assume http://localhost:port/path)
  * /path			(assume http://localhost:80/path

Headers are specified as a comma separated list of keypairs: --header name1(:|=)value1,name2(:|=)value2 ...
or specified multiple times: --header name1(:|=)value1 --header name2(:|=)value2`,
	}

	root.AddCommand(get, post, delete, run, sandbox, parse)
	return root
}

func buildGet() *cobra.Command {
	get := &cobra.Command{
		Use:   "get <url>",
		Short: "HTTP GET request.",
		Args:  cobra.ExactArgs(1),
		Run:   handleGet,
	}

	addCommonFlags(get)
	return get
}

func buildPost() *cobra.Command {
	post := &cobra.Command{
		Use:   `post <url> --json <body>`,
		Short: "HTTP POST request with a JSON body.",
		Long: `Make an HTTP POST request to the URL with a JSON body.
This command requires the --json flag, which should be a string conforming to valid JSON.`,
		Args: cobra.ExactArgs(1),
		Run:  handlePost,
	}

	post.Flags().String("json", "", "JSON body to use")
	addCommonFlags(post)
	return post
}

func buildDelete() *cobra.Command {
	delete := &cobra.Command{
		Use:   `delete <url>`,
		Short: "HTTP DELETE request.",
		Args:  cobra.ExactArgs(1),
		Run:   handleDelete,
	}

	addCommonFlags(delete)
	return delete
}

func buildRun() *cobra.Command {
	run := &cobra.Command{
		Use:   `run <file>`,
		Short: "Run requests from a spec file.",
		Long:  "The spec file must be a valid JSON or YAML file.",
		Args:  cobra.ExactArgs(1),
		Run:   handleRun,
	}

	run.Flags().StringSliceP(
		constants.RunTargetFlagName,
		"t",
		[]string{},
		`Run the specified target(s) from the file.
Use a comma separated list for multiple targets, e.g. --target a,b
or specify the flag multiple times, e.g. --target a --target b`)
	return run
}

func buildSandbox() *cobra.Command {
	sandbox := &cobra.Command{
		Use:   `sandbox`,
		Short: "Starts a local server. Default to port 8118.",
		Run:   rest.StartSandbox,
	}

	sandbox.Flags().IntP(
		constants.SandboxPortFlagName,
		"p",
		8118,
		`The port to use.`)
	return sandbox
}

func buildParseURL() *cobra.Command {
	parse := &cobra.Command{
		Use:   `parse-url <url>`,
		Short: "Parse the URL and print the results.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url, err := parse.ParseURL(args[0])
			if err != nil {
				fmt.Printf("error: %v\n", err)
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

func addCommonFlags(cmd *cobra.Command) {
	// Headers
	cmd.Flags().StringSlice(
		constants.HeaderFlagName,
		[]string{},
		`HTTP header to use in the request.
Value should be a keypair separated by equal sign (=) or colon (:), e.q. key=value.`)

	cmd.Flags().String(constants.OutputFileFlagName, "", "Output the response body to the filename.")

	// AWS signature V4 flags
	cmd.Flags().BoolP(constants.AWSSigV4FlagName, "4", false, "Use AWS signature V4 as authentication in the request. Requires the --aws-region option.")
	cmd.Flags().String(constants.AWSRegionFlagName, constants.DefaultAWSRegion, "The AWS region to use in the AWS signature.")
	cmd.Flags().String(constants.AWSProfileFlagName, "", "The name of an AWS profile in your AWS configuration. If not specified, environment variables are used.")

	// Sandbox
	cmd.Flags().Bool(constants.SandboxFlagName, false, "Run the request to a sandbox server.")
}
