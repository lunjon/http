package command

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/lunjon/httpreq/internal/rest"
	"github.com/spf13/cobra"
)

/* getHeaders assumes strArray is a string that consist of
comma separated keypairs in the format key(:|=)value,
wrapped inside [].
*/
func getHeaders(strArray string) (http.Header, error) {
	strArray = strings.TrimLeft(strArray, "[")
	strArray = strings.TrimRight(strArray, "]")

	if strArray == "" {
		return nil, nil
	}

	headers := http.Header{}
	re := regexp.MustCompile(`([a-zA-Z0-9\-_]+)[:=](.*)`)

	for _, h := range strings.Split(strArray, ",") {
		matches := re.FindAllStringSubmatch(h, -1)
		if matches == nil {
			return nil, fmt.Errorf("invalid header format: %s", h)
		}
		for _, match := range matches {
			key := match[1]
			value := match[2]
			headers.Add(key, value)
		}
	}

	return headers, nil
}

// Check if err != nil. If so, print the error, command usage (if printUsage is true)
// and exit the program with the given status code.
func checkError(err error, exitStatus int, printUsage bool, cmd *cobra.Command) {
	if err == nil {
		return
	}

	fmt.Printf("error: %v\n", err)
	if printUsage {
		cmd.Usage()
	}
	os.Exit(exitStatus)
}

func createClient(c *http.Client) *rest.Client {
	if c == nil {
		c = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	return rest.NewClient(c)
}
