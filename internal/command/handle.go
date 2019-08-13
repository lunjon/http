package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/lunjon/httpreq/internal/rest"
	"github.com/spf13/cobra"
)

// This function handles all from sending the request to printing the results.
func handleRequest(req *http.Request, cmd *cobra.Command) {
	res := rest.SendRequest(req)
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
