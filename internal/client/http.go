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
	contentTypes := concat("application", "json", "xml")
	contentTypes = append(contentTypes, concat("text", "html", "csv")...)

	utc := time.Now().UTC()

	KnownHTTPHeaders = map[string][]string{
		// https://www.iana.org/assignments/media-types/media-types.xhtml
		"Content-Type":     contentTypes,
		"Accept":           contentTypes,
		"Authorization":    {"Bearer ", "Basic "},
		"X-Correlation-Id": {},
		"Date":             {utc.Format(time.RFC1123), utc.Format(time.RFC1123Z)},
		"User-Agent":       {},
	}
}

func concat(name string, variant ...string) []string {
	values := make([]string, len(variant))
	for _, v := range variant {
		values = append(values, name+"/"+v)
	}
	return values
}
