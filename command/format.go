package command

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lunjon/http/util"
)

var FormatComponents = []string{"status", "header", "body"}

type Formatter interface {
	Format(*http.Response) ([]byte, error)
}

type DefaultFormatter struct {
	components []string
	color      bool
}

func NewDefaultFormatter(color bool, components []string) (*DefaultFormatter, error) {
	if len(components) > 3 {
		return nil, fmt.Errorf("invalid format specifiers: too many")
	}

	parsed := util.Map(components, strings.ToLower)
	for _, c := range parsed {
		if !util.Contains(FormatComponents, c) {
			return nil, fmt.Errorf("invalid format specifier: %s", c)
		}
	}

	return &DefaultFormatter{
		color:      color,
		components: parsed,
	}, nil
}

func (f *DefaultFormatter) Format(r *http.Response) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	for _, comp := range f.components {
		switch comp {
		case "status":
			f.addStatus(buf, r)
		case "header":
			f.addHeader(buf, r)
		case "body":
			if err := f.addBody(buf, r); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("invalid format specifier: %s", comp)
		}
	}
	return buf.Bytes(), nil
}

func (f *DefaultFormatter) addStatus(w io.Writer, r *http.Response) {
	fmt.Fprintln(w, r.Status)
}

func (f *DefaultFormatter) addHeader(w io.Writer, r *http.Response) {
	for name, value := range r.Header {
		fmt.Fprintf(w, "%s: %s\n", name, value)
	}
}

func (f *DefaultFormatter) addBody(w io.Writer, r *http.Response) error {
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
