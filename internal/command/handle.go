package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/lunjon/httpreq/internal/constants"
	"github.com/lunjon/httpreq/internal/rest"
	"github.com/lunjon/httpreq/internal/runner"
	"github.com/lunjon/httpreq/pkg/parse"
	"github.com/spf13/cobra"
)

func handleGet(cmd *cobra.Command, args []string) {
	handleRequest(http.MethodGet, cmd, args)
}

func handlePost(cmd *cobra.Command, args []string) {
	handleRequest(http.MethodPost, cmd, args)
}

func handleDelete(cmd *cobra.Command, args []string) {
	handleRequest(http.MethodDelete, cmd, args)
}

func handleRequest(method string, cmd *cobra.Command, args []string) {
	url := args[0]
	headerString, _ := cmd.Flags().GetString(constants.HeaderFlagName)
	header, err := getHeaders(headerString)
	checkError(err, 2, true, cmd)

	var body []byte

	if method == http.MethodPost {
		json, _ := cmd.Flags().GetString(constants.JSONBodyFlagName)
		if json == "" {
			fmt.Println("no or invalid JSON body specified")
			cmd.Usage()
			os.Exit(2)
		}

		body = []byte(json)
	}

	var client *rest.Client

	// Sandbox should send request to a local test server
	sandbox, _ := cmd.Flags().GetBool(constants.SandboxFlagName)
	if sandbox {
		server := httptest.NewServer(&rest.SandboxHandler{})
		defer server.Close()
		client = createClient(server.Client())

		// Re-write the URL to get correct path
		u, err := parse.ParseURL(url)
		checkError(err, 2, false, cmd)
		url = server.URL + u.Path
	} else {
		client = createClient(nil)
	}

	req, err := client.BuildRequest(method, url, body, header)
	checkError(err, 2, false, cmd)

	signRequest, _ := cmd.Flags().GetBool(constants.AWSSigV4FlagName)
	if signRequest {
		region, _ := cmd.Flags().GetString(constants.AWSRegionFlagName)
		profile, _ := cmd.Flags().GetString(constants.AWSProfileFlagName)
		err = client.SignRequest(req, nil, region, profile)

		checkError(err, 2, true, cmd)
	}

	sendRequest(client, req, cmd)
}

func sendRequest(client *rest.Client, req *http.Request, cmd *cobra.Command) {
	res := client.SendRequest(req)
	checkError(res.Error(), 1, false, cmd)

	// Write request result to stdout
	fmt.Println(res)

	body, err := res.Body()
	checkError(err, 1, false, cmd)

	filename, _ := cmd.Flags().GetString(constants.OutputFileFlagName)
	if len(body) == 0 {
		if filename != "" {
			fmt.Println("no response body to write to file")
		}
		return
	}

	dst := &bytes.Buffer{}
	err = json.Indent(dst, body, "", "  ")
	if err == nil {
		body = dst.Bytes()
	}

	if filename == "" {
		fmt.Println(string(body))
	} else {
		err = ioutil.WriteFile(filename, body, 0644)
		checkError(err, 1, false, cmd)
	}
}

func handleRun(cmd *cobra.Command, args []string) {
	spec, err := runner.LoadSpec(args[0])
	checkError(err, 2, false, cmd)
	err = spec.Validate()
	checkError(err, 2, false, cmd)

	targets, _ := cmd.Flags().GetStringSlice(constants.RunTargetFlagName)

	var rr *runner.Runner
	var client *rest.Client

	sandbox, _ := cmd.Flags().GetBool(constants.SandboxFlagName)

	if sandbox {
		server := httptest.NewServer(&rest.SandboxHandler{})
		defer server.Close()
		client = createClient(server.Client())

		rr = runner.NewRunner(spec, client)
		err := rr.SetBaseURL(server.URL)
		checkError(err, 2, false, cmd)

	} else {
		client = createClient(nil)
		rr = runner.NewRunner(spec, client)
	}

	results, err := rr.Run(targets...)
	checkError(err, 1, false, cmd)

	for _, r := range results {
		fmt.Println(r)
	}
}
