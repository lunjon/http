package client

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	portPattern      = regexp.MustCompile(`^:\d+`)
	protoPattern     = regexp.MustCompile(`^https?://`)
	localhostPattern = regexp.MustCompile(`^(localhost|127\.0\.0\.1)`)
	hostPattern      = regexp.MustCompile(`^[a-z](\.[a-z]+)*`)
	aliasPattern     = regexp.MustCompile(`\{[\w]+\}`)
)

// ParseURL parses the given URL
func ParseURL(url string, aliases map[string]string) (*url.URL, error) {
	url = strings.TrimSpace(url)
	if url == "" {
		return nil, fmt.Errorf("empty URL")
	}

	if aliases != nil && aliasPattern.MatchString(url) {
		var err error
		url, err = substitute(url, aliases)
		if err != nil {
			return nil, err
		}
	}

	// :port/path
	if portPattern.MatchString(url) {
		return parseURL("http://localhost" + url)
	}

	// localhost
	if localhostPattern.MatchString(url) {
		return parseURL("http://" + url)
	}

	// https?://...
	if protoPattern.MatchString(url) {
		return parseURL(url)
	}

	// api.com...
	if hostPattern.MatchString(url) {
		return parseURL("https://" + url)
	}

	return nil, fmt.Errorf("invalid URL format: %s", url)
}

func parseURL(s string) (*url.URL, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("empty URL")
	}
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	if u.Host == "" {
		return nil, fmt.Errorf("missing host")
	}

	if !strings.Contains(u.Scheme, "http") {
		return nil, fmt.Errorf("invalid scheme: '%s'", u.Scheme)
	}

	return u, nil
}

func substitute(url string, aliases map[string]string) (string, error) {
	matches := aliasPattern.FindAllStringSubmatch(url, -1)
	if len(matches) == 0 {
		return "", fmt.Errorf("expected aliases but found none")
	}

	for _, match := range matches[0] {
		s := strings.TrimPrefix(match, "{")
		s = strings.TrimSuffix(s, "}")
		sub, found := aliases[s]
		if !found {
			return "", fmt.Errorf("unknown alias: %s", s)
		}
		url = strings.ReplaceAll(url, match, sub)
	}
	return url, nil
}
