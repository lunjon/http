package client

import "time"

type MIMEType string

const (
	MIMETypeHTML    MIMEType = "text/html"
	MIMETypeCSV     MIMEType = "text/csv"
	MIMETypeJSON    MIMEType = "application/json"
	MIMETypeXML     MIMEType = "application/xml"
	MIMETypeUnknown MIMEType = "unknown"
)

func (m MIMEType) String() string {
	return string(m)
}

var (
	KnownHTTPHeaders map[string][]string
)

func init() {
	// Build known headers and a list of common values.
	// See more: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers

	// https://www.iana.org/assignments/media-types/media-types.xhtml
	application := withPrefix("application", "json", "xml", "gzip", "zip")
	text := withPrefix("text", "html", "csv", "css", "xml", "")

	contentTypes := make([]string, len(application)+len(text))
	contentTypes = append(contentTypes, application...)
	contentTypes = append(contentTypes, text...)

	utc := time.Now().UTC()

	KnownHTTPHeaders = map[string][]string{
		"Content-Type":     contentTypes,
		"Content-Encoding": {"gzip", "compress", "deflate", "br"},
		"Accept":           contentTypes,
		"Accept-Encoding":  {},
		"Authorization":    {"Bearer ", "Basic "},
		"X-Correlation-Id": {},
		"User-Agent":       {},
		"Date":             {utc.Format(time.RFC1123), utc.Format(time.RFC1123Z)},
	}
}

func withPrefix(prefix string, variant ...string) []string {
	values := make([]string, len(variant))
	for _, v := range variant {
		values = append(values, prefix+"/"+v)
	}
	return values
}
