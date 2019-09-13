package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"github.com/lunjon/httpreq/internal/rest"
	"github.com/lunjon/httpreq/internal/runner"
	"github.com/spf13/cobra"
)

func handleGet(cmd *cobra.Command, args []string) {
	handle(http.MethodGet, cmd, args)
}

func handlePost(cmd *cobra.Command, args []string) {
	handle(http.MethodPost, cmd, args)
}

func handleDelete(cmd *cobra.Command, args []string) {
	handle(http.MethodDelete, cmd, args)
}

func handleRun(cmd *cobra.Command, args []string) {
	spec := args[0]
	runner, err := runner.Load(spec)
	checkError(err, 2, false, cmd)

	targetString := cmd.Flag(RunTargetFlagName).Value.String()
	targetString = targetString[1 : len(targetString)-1]

	targets := []string{}
	if targetString != "" {
		targets = strings.Split(targetString, ",")
	}

	client := rest.NewClient(nil)

	results, err := runner.Run(client, targets...)
	checkError(err, 1, false, cmd)

	for _, r := range results {
		fmt.Println(r)
	}
}

func handle(method string, cmd *cobra.Command, args []string) {
	url := args[0]
	headerString := getStringFlagValue(HeaderFlagName, cmd)
	header, err := getHeaders(headerString)
	checkError(err, 2, true, cmd)

	var body []byte

	if method == http.MethodPost {
		json := getStringFlagValue(JSONBodyFlagName, cmd)
		if json == "" {
			fmt.Println("no or invalid JSON body specified")
			cmd.Usage()
			os.Exit(2)
		}

		body = []byte(json)
	}

	var client *rest.Client

	// Sandbox should send request to a local test server
	sandbox := getBoolFlagValue(SandboxFlagName, cmd)
	if sandbox {
		server := httptest.NewServer(&SandboxHandler{})
		defer server.Close()
		client = createClient(server.Client())
		url = server.URL
	} else {
		client = createClient(nil)
	}

	req, err := client.BuildRequest(method, url, body, header)
	checkError(err, 2, true, cmd)

	signRequest := getBoolFlagValue(AWSSigV4FlagName, cmd)
	if signRequest {
		region := getStringFlagValue(AWSRegionFlagName, cmd)
		profile := getStringFlagValue(AWSProfileFlagName, cmd)
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

	filename := getStringFlagValue(OutputFileFlagName, cmd)
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

func createClient(c *http.Client) *rest.Client {
	if c == nil {
		c = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	return rest.NewClient(c)
}
