package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/lunjon/httpreq/internal/rest"
	"github.com/spf13/cobra"
)

func handlePost(cmd *cobra.Command, args []string) {
	url := args[0]
	headerString := getStringFlagValue(HeaderFlagName, cmd)
	header, err := getHeaders(headerString)
	checkError(err, 2, true, cmd)

	json := getStringFlagValue(JSONBodyFlagName, cmd)
	if json == "" {
		fmt.Println("no or invalid JSON body specified")
		cmd.Usage()
		os.Exit(2)
	}

	body := []byte(json)

	req, err := rest.BuildRequest(http.MethodPost, url, body, header)
	checkError(err, 2, true, cmd)

	signRequest := getBoolFlagValue(AWSSigV4FlagName, cmd)
	if signRequest {
		region := getStringFlagValue(AWSRegionFlagName, cmd)
		profile := getStringFlagValue(AWSProfileFlagName, cmd)
		err = rest.SignRequest(req, body, region, profile)

		checkError(err, 2, true, cmd)
	}

	handleRequest(req, cmd)
}
