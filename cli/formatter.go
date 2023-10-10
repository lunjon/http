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
	"github.com/lunjon/http/internal/util"
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

func FormatterFromString(format Format, comps string) (Formatter, error) {
	var components []string
	switch comps {
	case "", "none":
		return NullFormatter{}, nil
	case "all":
		components = ResponseComponents
	default:
		components = strings.Split(strings.ToLower(comps), ",")
		if len(components) > len(ResponseComponents) {
			return nil, fmt.Errorf("invalid format specifiers: too many")
		}
		components = util.Map(components, strings.TrimSpace)
	}

	for _, c := range components {
		if !util.Contains(ResponseComponents, c) {
			return nil, fmt.Errorf("invalid format specifier: %s", c)
		}
	}

	switch format {
	case TextFormat:
		return &textFormatter{components}, nil
	case JSONFormat:
		return &jsonFormatter{components}, nil
	}

	return nil, fmt.Errorf("unknown format: %s", format)
}

type textFormatter struct {
	components []string
}

func (f *textFormatter) FormatResponse(r *http.Response) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	for _, comp := range f.components {
		switch comp {
		case "status":
			fmt.Fprintln(buf, r.Status)
		case "headers":
			f.addHeaders(buf, r)
		case "body":
			defer r.Body.Close()
			b, err := io.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}

			// Is it application/json?
			//   => try to indent nicely
			if r.Header.Get(contentTypeHeader) == string(client.MIMETypeJSON) {
				err = json.Indent(buf, b, "", "  ")
				if err != nil {
					_, err = buf.Write(b)
				}
			} else {
				_, err = buf.Write(b)
			}

			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("invalid format specifier: %s", comp)
		}
	}
	return buf.Bytes(), nil
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
	type output struct {
		Status  *string           `json:"status,omitempty"`
		Headers map[string]string `json:"headers,omitempty"`
		Body    *string           `json:"body,omitempty"`
	}

	body := new(output)
	for _, comp := range f.components {
		switch comp {
		case "status":
			body.Status = &r.Status
		case "headers":
			body.Headers = headerToMap(r.Header)
		case "body":
			defer r.Body.Close()
			b, err := io.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}

			buf := bytes.NewBuffer(nil)

			// Is it application/json?
			//   => try to indent nicely
			if r.Header.Get(contentTypeHeader) == string(client.MIMETypeJSON) {
				err = json.Indent(buf, b, "", "  ")
				if err != nil {
					buf.WriteString(string(b))
				}
			} else {
				buf.WriteString(string(b))
			}

			bs := buf.String()
			body.Body = &bs
		default:
			return nil, fmt.Errorf("invalid format specifier: %s", comp)
		}
	}

	return json.MarshalIndent(body, "", " ")
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
