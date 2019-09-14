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
	"github.com/lunjon/httpreq/pkg/parse"
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
	headerString, _ := cmd.Flags().GetString(HeaderFlagName)
	header, err := getHeaders(headerString)
	checkError(err, 2, true, cmd)

	var body []byte

	if method == http.MethodPost {
		json, _ := cmd.Flags().GetString(JSONBodyFlagName)
		if json == "" {
			fmt.Println("no or invalid JSON body specified")
			cmd.Usage()
			os.Exit(2)
		}

		body = []byte(json)
	}

	var client *rest.Client

	// Sandbox should send request to a local test server
	sandbox, _ := cmd.Flags().GetBool(SandboxFlagName)
	if sandbox {
		server := httptest.NewServer(&SandboxHandler{})
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

	signRequest, _ := cmd.Flags().GetBool(AWSSigV4FlagName)
	if signRequest {
		region, _ := cmd.Flags().GetString(AWSRegionFlagName)
		profile, _ := cmd.Flags().GetString(AWSProfileFlagName)
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

	filename, _ := cmd.Flags().GetString(OutputFileFlagName)
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
