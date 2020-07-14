package parse

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"
)

const (
	HTTP      = "http"
	HTTPS     = "https"
	localhost = "localhost"
)

type URL struct {
	Scheme string
	Port   int
	Host   string
	Path   string
	Query  string
}

func (url *URL) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s://%s", url.Scheme, url.Host))

	if url.Port != 80 && url.Port != 443 {
		builder.WriteString(fmt.Sprintf(":%d", url.Port))
	}

	builder.WriteString(url.Path)
	if url.Query != "" {
		builder.WriteString("?")
		builder.WriteString(url.Query)
	}

	return builder.String()
}

func (url *URL) BaseURL() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s://%s", url.Scheme, url.Host))

	if url.Port != 80 && url.Port != 443 {
		builder.WriteString(fmt.Sprintf(":%d", url.Port))
	}

	return builder.String()
}

func (url *URL) DetailString() string {
	var builder strings.Builder
	const padding = 3
	w := tabwriter.NewWriter(&builder, 0, 0, padding, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "URL:\t%s\n", url.String())
	fmt.Fprintf(w, "Scheme:\t%s\n", url.Scheme)
	fmt.Fprintf(w, "Port:\t%d\n", url.Port)
	fmt.Fprintf(w, "Host:\t%s\n", url.Host)
	fmt.Fprintf(w, "Path:\t%s", url.Path)
	if url.Query != "" {
		fmt.Fprintf(w, "\nQuery:\t%s", url.Query)
	}
	w.Flush()

	return builder.String()
}

/*
ParseURL tries to parse the url and match it to a format according to:
	- /path -> http://localhost/path
	- :port/path -> http://localhost:port/path
	- [https?://]localhost[:port]/path -> https?://localhost[:port]/path
	- [https?://]host.domain[:port]/path -> [https?://]host.domain[:port]/path
	- [https?://]ip[:port]/path -> [https?://]ip[:port]/path
*/
func ParseURL(url string) (*URL, error) {
	url = strings.TrimRight(url, "/")
	local := regexp.MustCompile(`^(/[0-9a-zA-Z\-?&_%=\.])+`)
	if local.MatchString(url) {
		return &URL{
			Scheme: HTTP,
			Port:   80,
			Host:   localhost,
			Path:   url,
		}, nil
	}

	localPort := regexp.MustCompile(`^:(\d+)(/[0-9a-zA-Z\-?&_%=\.]+)*`)
	matches := localPort.FindAllStringSubmatch(url, -1)
	if matches != nil {
		portStr := matches[0][1]
		port, _ := strconv.Atoi(portStr)
		path := matches[0][2]
		return &URL{
			Scheme: HTTP,
			Port:   port,
			Host:   localhost,
			Path:   path,
		}, nil
	}

	// Scheme
	s := "((https?)(://))?"
	// Host
	h := `(([a-z0-9\-]+)(\.[a-z0-9\-]+)*)`
	// Port
	p := `(:(\d+))?`
	// Path
	r := `((/[0-9a-zA-Z\-&_%=\.]+)*)`
	// Query
	q := `(\?([0-9a-zA-Z\-&_%=\.]+))?`

	def := regexp.MustCompile(s + h + p + r + q)
	matches = def.FindAllStringSubmatch(url, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("invalid url: %s", url)
	}

	scheme := matches[0][2]
	host := matches[0][4]
	port := 443

	switch {
	case scheme == "" && regexp.MustCompile(`(localhost|127\.0\.0\.1)`).MatchString(host):
		scheme = HTTP
		port = 80
	case scheme == "":
		scheme = HTTPS
	case scheme == HTTP:
		port = 80
	}

	if scheme == host {
		return nil, fmt.Errorf("invalid url: %s", url)
	}

	portStr := matches[0][8]
	if portStr != "" {
		port, _ = strconv.Atoi(portStr)
	}

	path := matches[0][9]
	query := matches[0][12]

	return &URL{
		Scheme: scheme,
		Host:   host,
		Port:   port,
		Path:   path,
		Query:  query,
	}, nil
}
