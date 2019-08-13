package rest

import (
	"fmt"
	"regexp"
	"strings"
)

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
