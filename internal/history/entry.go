package history

import (
	"io"
	"net/http"
	"time"

	"github.com/lunjon/http/internal/client"
)

type Entry struct {
	Timestamp      time.Time       `json:"timestamp"`
	Method         string          `json:"method"`
	URL            string          `json:"url"`
	Header         http.Header     `json:"headers"`
	Body           []byte          `json:"body"`
	ClientSettings client.Settings `json:"clientSettings"`
}

func NewEntry(req *http.Request, settings client.Settings) (Entry, error) {
	var body []byte
	if req.Body != nil {
		reader, err := req.GetBody()
		if err != nil {
			return Entry{}, err
		}

		body, err = io.ReadAll(reader)
		if err != nil {
			return Entry{}, err
		}
	}

	return Entry{
		Timestamp:      time.Now(),
		Method:         req.Method,
		URL:            req.URL.String(),
		Header:         req.Header,
		Body:           body,
		ClientSettings: settings,
	}, nil
}
