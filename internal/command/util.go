package command

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func getStringFlagValue(name string, cmd *cobra.Command) string {
	flag := cmd.Flag(name)
	if flag == nil {
		panic("invalid flag: " + name)
	}

	return flag.Value.String()
}

func getBoolFlagValue(name string, cmd *cobra.Command) bool {
	flag := cmd.Flag(name)
	if flag == nil {
		panic("invalid flag: " + name)
	}

	val, _ := strconv.ParseBool(flag.Value.String())
	return val
}

func parseURL(route string) (string, error) {
	route = strings.TrimRight(route, "/")
	local := regexp.MustCompile(`^(/[0-9a-zA-Z\-?&_%])+`)
	if local.MatchString(route) {
		return "http://localhost/" + strings.TrimLeft(route, "/"), nil
	}

	localPort := regexp.MustCompile(`^:(\d+)(/[0-9a-zA-Z\-?&_%]+)*`)
	if localPort.MatchString(route) {
		return "http://localhost" + route, nil
	}

	localhost := regexp.MustCompile(`(https?)?(localhost|127\.0\.0\.1)(:\d+)?(/[0-9a-zA-Z\-?&_%]+)*`)
	if localhost.MatchString(route) {
		if strings.HasPrefix(route, "http") {
			return route, nil
		}
		return "http://" + route, nil
	}

	proto := regexp.MustCompile(`^https?://([a-z0-9\-]+)(\.[a-z0-9\-]+)+(:\d+)?(/[0-9a-zA-Z\-?&_%]+)*`)
	if proto.MatchString(route) {
		return route, nil
	}

	missingProto := regexp.MustCompile(`^([a-z0-9\-]+)(\.[a-z0-9\-]+)+(:\d+)?(/[0-9a-zA-Z\-?&_%]+)*`)
	if missingProto.MatchString(route) {
		return "https://" + route, nil
	}

	return "", fmt.Errorf("Invalid route: %s", route)
}

// getHeaders assumes strArray is a string that consist of
// comma separated keypairs in the format key(:|=)value.
func getHeaders(strArray string) (http.Header, error) {
	strArray = strings.TrimLeft(strArray, "[")
	strArray = strings.TrimRight(strArray, "]")
	strArray = strings.TrimSpace(strArray)

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
