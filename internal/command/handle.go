package command

import (
	"fmt"
	"github.com/lunjon/httpreq/internal/constants"
	"github.com/lunjon/httpreq/internal/rest"
	"github.com/lunjon/httpreq/internal/runner"
	"github.com/lunjon/httpreq/pkg/parse"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
)

var (
	// Output is where the handler will write the results.
	// It is initialized to os.Stdout as default
	Output io.Writer = os.Stdout
)

func handleGet(cmd *cobra.Command, args []string) {
	handleRequest(http.MethodGet, nil, cmd, args)
}

func handlePost(cmd *cobra.Command, args []string) {
	bodyFlag, _ := cmd.Flags().GetString(constants.BodyFlagName)
	if bodyFlag == "" {
		fmt.Println("No or invalid JSON body specified")
		cmd.Usage()
		os.Exit(2)
	}

	// We first try to read as a file
	body, err := ioutil.ReadFile(bodyFlag)

	if os.IsNotExist(err) {
		log.Print("Failed to open input file: ", err)
		checkError(err, 1, false, cmd)
	}

	if body == nil {
		// Assume that the content was given as string
		log.Print("Assuming body was given as content string")
		body = []byte(bodyFlag)
	}

	handleRequest(http.MethodPost, body,  cmd, args)
}

func handleDelete(cmd *cobra.Command, args []string) {
	handleRequest(http.MethodDelete, nil, cmd, args)
}

func handleRequest(method string, body []byte, cmd *cobra.Command, args []string) {
	url := args[0]
	headerString, _ := cmd.Flags().GetStringSlice(constants.HeaderFlagName)
	header, err := getHeaders(headerString)
	checkError(err, 2, true, cmd)

	timeout, _ := cmd.Flags().GetInt(constants.TimeoutFlagName)
	log.Printf("Using timeout of %d seconds\n", timeout)

	var client *rest.Client

	// Sandbox should send request to a local test server
	sandbox, _ := cmd.Flags().GetBool(constants.SandboxFlagName)
	if sandbox {
		log.Println("Using sandbox mode")
		server := httptest.NewServer(&rest.SandboxHandler{})
		defer server.Close()
		client = createClient(server.Client(), timeout)

		// Re-write the URL to get correct path
		u, err := parse.ParseURL(url)
		checkError(err, 2, false, cmd)
		url = server.URL + u.Path
	} else {
		client = createClient(nil, timeout)
	}

	req, err := client.BuildRequest(method, url, body, header)
	checkError(err, 2, false, cmd)

	signRequest, _ := cmd.Flags().GetBool(constants.AWSSigV4FlagName)
	if signRequest {
		log.Println("Adding AWS Sig V4 to the request")
		region, _ := cmd.Flags().GetString(constants.AWSRegionFlagName)
		profile, _ := cmd.Flags().GetString(constants.AWSProfileFlagName)
		err = client.SignRequest(req, nil, region, profile)

		checkError(err, 2, true, cmd)
	}

	res := client.SendRequest(req)
	checkError(res.Error(), 1, false, cmd)
	outputResults(cmd, res)
}

func handleRun(cmd *cobra.Command, args []string) {
	spec, err := runner.LoadSpec(args[0])
	checkError(err, 2, false, cmd)

	err = spec.Validate()
	checkError(err, 2, false, cmd)

	targets, _ := cmd.Flags().GetStringSlice(constants.RunTargetFlagName)
	timeout, _ := cmd.Flags().GetInt(constants.TimeoutFlagName)

	var rr *runner.Runner
	var client *rest.Client

	sandbox, _ := cmd.Flags().GetBool(constants.SandboxFlagName)
	if sandbox {
		server := httptest.NewServer(&rest.SandboxHandler{})
		defer server.Close()
		client = createClient(server.Client(), timeout)

		rr = runner.NewRunner(spec, client)
		err := rr.SetBaseURL(server.URL)
		checkError(err, 2, false, cmd)
	} else {
		client = createClient(nil, timeout)
		rr = runner.NewRunner(spec, client)
	}

	results, err := rr.Run(targets...)
	checkError(err, 1, false, cmd)

	outputResults(cmd, results...)
}

func outputResults(cmd *cobra.Command, results ...*rest.Result) {
	log.Printf("Output %d results", len(results))

	printResponseBodyOnly, _ := cmd.Flags().GetBool(constants.ResponseBodyOnlyFlagName)
	filename, _ := cmd.Flags().GetString(constants.OutputFileFlagName)

	var output io.Writer
	var writeToFile bool
	output = Output

	if filename != "" {
		log.Print("Writing results to file")
		writeToFile = true
		f, err := os.Create(filename)
		checkError(err, 1, false, cmd)
		output = f
	}

	var body string
	for _, r := range results {
		if !printResponseBodyOnly {
			log.Print("Printing response info")
			fmt.Fprintln(output, r.Info())
		}

		if writeToFile {
			b, err := r.Body()
			checkError(err, 1, false, cmd)
			body = string(b)
		} else {
			body, _ = r.BodyFormatString()
		}

		_, err := fmt.Fprintln(output, body)
		checkError(err, 1, false, cmd)
	}
}
