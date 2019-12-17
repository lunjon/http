package command

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/lunjon/httpreq/internal/rest"
	"github.com/spf13/cobra"
)

func getHeaders(arr []string) (http.Header, error) {
	if len(arr) == 0 {
		return nil, nil
	}

	headers := http.Header{}
	re := regexp.MustCompile(`([a-zA-Z0-9\-_]+)[:=](.*)`)

	for _, h := range arr {
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

func createClient(c *http.Client, timeout int) *rest.Client {
	if c == nil {
		c = &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}
	}

	return rest.NewClient(c, timeout)
}
