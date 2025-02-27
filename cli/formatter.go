package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lunjon/http/internal/client"
	"github.com/lunjon/http/internal/history"
	"github.com/lunjon/http/internal/style"
	"github.com/lunjon/http/internal/types"
)

var ResponseComponents = []string{"status", "headers", "body"}

type Format string

const (
	TextFormat Format = "text"
	JSONFormat Format = "json"
)

type Formatter interface {
	FormatResponse(*http.Response) ([]byte, error)
	FormatHistory([]history.Entry) ([]byte, error)
}

type NullFormatter struct{}

func (f NullFormatter) FormatResponse(*http.Response) ([]byte, error) { return nil, nil }
func (f NullFormatter) FormatHistory([]history.Entry) ([]byte, error) { return nil, nil }

func FormatterFromString(format Format) (Formatter, error) {
	switch format {
	case TextFormat:
		return &textFormatter{}, nil
	case JSONFormat:
		return &jsonFormatter{}, nil
	}

	return nil, fmt.Errorf("unknown format: %s", format)
}

// textFormatter will only output the body, if any.
type textFormatter struct{}

func (f *textFormatter) FormatResponse(r *http.Response) ([]byte, error) {
	return readBody(r)
}

func readBody(r *http.Response) ([]byte, error) {
	if r.StatusCode == 204 && r.Header.Get("Content-Length") == "" {
		return nil, nil
	}

	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)

	// Is it application/json? Try to indent nicely.
	if strings.Contains(r.Header.Get(contentTypeHeader), string(client.MIMETypeJSON)) {
		err = json.Indent(buf, b, "", "  ")
		if err != nil {
			_, err = buf.Write(b)
		}
	} else {
		_, err = buf.Write(b)
	}

	return buf.Bytes(), err
}

func (f *textFormatter) addHeaders(w io.Writer, r *http.Response) {
	taber := types.NewTaber("")
	for name, value := range headerToMap(r.Header) {
		n := fmt.Sprintf("%s:", name)
		v := fmt.Sprint(value)
		taber.WriteLine(style.Bold.Render(n), v)
	}
	fmt.Fprint(w, taber.String())
}

func (f *textFormatter) FormatHistory([]history.Entry) ([]byte, error) {
	return nil, nil
}

type jsonFormatter struct {
	components []string
}

func (f *jsonFormatter) FormatResponse(r *http.Response) ([]byte, error) {
	output := struct {
		Status     string            `json:"status,omitempty"`
		StatusCode int               `json:"statusCode,omitempty"`
		Headers    map[string]string `json:"headers,omitempty"`
		Body       *string           `json:"body,omitempty"`
	}{
		Status:     r.Status,
		StatusCode: r.StatusCode,
		Headers:    headerToMap(r.Header),
	}

	body, err := readBody(r)
	if err != nil {
		return nil, err
	}

	if body != nil {
		b := string(body)
		output.Body = &b
	}

	return json.MarshalIndent(output, "", " ")
}

func (f *jsonFormatter) FormatHistory([]history.Entry) ([]byte, error) {
	return nil, nil
}

func headerToMap(h http.Header) map[string]string {
	headers := map[string]string{}
	for name, values := range h {
		headers[name] = strings.Join(values, "; ")
	}
	return headers
}
