package command

import (
	"net/http"

	"github.com/lunjon/httpreq/internal/rest"
	"github.com/spf13/cobra"
)

func handleDelete(cmd *cobra.Command, args []string) {
	url := args[0]
	headerString := getStringFlagValue(HeaderFlagName, cmd)
	header, err := getHeaders(headerString)
	checkError(err, 2, true, cmd)

	req, err := rest.BuildRequest(http.MethodDelete, url, nil, header)
	checkError(err, 2, true, cmd)

	signRequest := getBoolFlagValue(AWSSigV4FlagName, cmd)
	if signRequest {
		region := getStringFlagValue(AWSRegionFlagName, cmd)
		profile := getStringFlagValue(AWSProfileFlagName, cmd)
		err = rest.SignRequest(req, nil, region, profile)

		checkError(err, 2, true, cmd)
	}

	handleRequest(req, cmd)
}
