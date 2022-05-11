package command

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lunjon/http/style"
	"github.com/lunjon/http/util"
)

var FormatComponents = []string{"status", "statuscode", "headers", "body"}

type Formatter interface {
	Format(*http.Response) ([]byte, error)
}

type DefaultFormatter struct {
	components []string
	color      bool
	bold       style.StyleFunc
}

func NewDefaultFormatter(color bool, components []string) (*DefaultFormatter, error) {
	if len(components) > len(FormatComponents) {
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
		bold:       style.NewBuilder().Bold(true).Build(),
	}, nil
}

func (f *DefaultFormatter) Format(r *http.Response) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	for _, comp := range f.components {
		switch comp {
		case "status":
			fmt.Fprintln(buf, r.Status)
		case "statuscode":
			fmt.Fprintln(buf, r.StatusCode)
		case "headers":
			f.addHeaders(buf, r)
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

func (f *DefaultFormatter) addHeaders(w io.Writer, r *http.Response) {
	taber := util.NewTaber("")
	for name, value := range r.Header {
		n := fmt.Sprintf("%s:", name)
		v := fmt.Sprint(value)
		taber.WriteLine(f.bold(n), v)
	}
	fmt.Fprint(w, taber.String())
}

func (f *DefaultFormatter) addBody(w io.Writer, r *http.Response) error {
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	fmt.Fprintln(w, "")
	return err
}
