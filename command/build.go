package command

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/lunjon/httpreq/logging"
	"github.com/lunjon/httpreq/rest"
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
	alias, err := readAliasFile()
	checkErr(err)
	handler := NewHandler(restClient, logger, h, alias)
	return handler
}

// Build the root command for httpreq and set version.
func Build(version string) *cobra.Command {
	handler := createHandler()

	root := &cobra.Command{
		Version: version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Verbose flag
			verbose, _ := cmd.Flags().GetBool(VerboseFlagName)
			handler.Verbose(verbose)

			// Timeout flag
			timeout, _ := cmd.Flags().GetDuration(TimeoutFlagName)
			handler.Timeout(timeout)
		},
		Use:   "httpreq",
		Short: "httpreq <method> <url> [options]",
		Long:  description,
	}

	// HTTP
	get := buildGet(handler)
	head := buildHead(handler)
	post := buildPost(handler)
	put := buildPut(handler)
	patch := buildPatch(handler)
	delete := buildDelete(handler)
	root.AddCommand(get, head, post, put, patch, delete)

	// URL alias
	url := buildURL()
	root.AddCommand(url)

	// Command for generating completion
	gen := buildGen(root)
	root.AddCommand(gen)

	// Persistant flags
	root.PersistentFlags().BoolP(VerboseFlagName, "V", false, "Shows debug logs.")
	root.PersistentFlags().DurationP(
		TimeoutFlagName,
		"T",
		defaultTimeout,
		"Request timeout in seconds.")

	return root
}

func buildGen(root *cobra.Command) *cobra.Command {
	filenameDefault := ""
	gen := &cobra.Command{
		Use:     "gen <type>",
		Aliases: []string{"g"},
		Short:   "Generate completion for shell <type>: bash, zsh, fish",
		Args:    cobra.ExactArgs(1),
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
		Use:     "head <url>",
		Aliases: []string{"h", "hd"},
		Short:   "HTTP HEAD request.",
		Args:    cobra.ExactArgs(1),
		Run:     handler.Head,
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
		Use:     `patch <url> --body <body>`,
		Aliases: []string{"pa"},
		Short:   "HTTP PATCH request with a body.",
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
		Use:     `put <url> --body <body>`,
		Aliases: []string{"pu"},
		Short:   "HTTP PUT request with a body.",
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
	del := &cobra.Command{
		Use:     `delete <url>`,
		Aliases: []string{"d", "de", "del"},
		Short:   "HTTP DELETE request.",
		Args:    cobra.ExactArgs(1),
		Run:     handler.Delete,
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

func buildURL() *cobra.Command {
	listAlias := func() {
		alias, err := readAliasFile()
		checkErr(err)
		for a, url := range alias {
			fmt.Printf("%s  ->  %s\n", a, url)
		}
	}

	setAlias := func(alias, url string) {
		filepath := getAliasFilepath()
		file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		checkErr(err)
		defer file.Close()
		_, err = file.WriteString(fmt.Sprintf("%s %s\n", alias, url))
		checkErr(err)
	}

	return &cobra.Command{
		Use:   "url <alias> <url>",
		Short: "Create a persistant URL alias.",
		Run: func(_ *cobra.Command, args []string) {
			switch len(args) {
			case 0:
				listAlias()
			case 2:
				setAlias(args[0], args[1])
			default:
				fmt.Fprintln(os.Stderr, "unknown number of arguments")
			}
		},
	}
}

func addCommonFlags(cmd *cobra.Command, handler *Handler) {
	// Headers
	cmd.Flags().VarP(handler.header, HeaderFlagName, "H", `HTTP header, may be specified multiple times.
The value must conform to the format "name: value". "name" and "value" can
be separated by either a colon ":" or an equal sign "=", and the space
between is optional. Can be set in the same format using the env. variable
DEFAULT_HEADERS, where multiple headers must be separated by an |.`)

	// AWS signature V4 flags
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

func readAliasFile() (map[string]string, error) {
	alias := make(map[string]string)
	filepath := getAliasFilepath()
	file, err := os.Open(filepath)
	if os.IsNotExist(err) {
		return alias, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		s := strings.Split(line, " ")
		if len(s) != 2 {
			continue
		}
		alias[s[0]] = s[1]
	}
	return alias, nil
}

func getAliasFilepath() string {
	filepath, err := os.UserHomeDir()
	checkErr(err)
	return path.Join(filepath, ".httpreq", "alias")
}
