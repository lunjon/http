package command

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"text/tabwriter"
)

type Formatter interface {
	Format(*http.Response) ([]byte, error)
}

type DefaultFormatter struct{}

func (f DefaultFormatter) Format(r *http.Response) ([]byte, error) {
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

type NullFormatter struct{}

func (f NullFormatter) Format(*http.Response) ([]byte, error) {
	return nil, nil
}

type BriefFormatter struct{}

func (f BriefFormatter) Format(r *http.Response) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	w := tabwriter.NewWriter(buf, 0, 0, 2, ' ', tabwriter.TabIndent)

	fmt.Fprintf(w, "Method:\t%s\n", r.Request.Method)
	fmt.Fprintf(w, "URL:\t%s\n", r.Request.URL.String())
	fmt.Fprintf(w, "Status:\t%s", r.Status)
	err := w.Flush()

	return buf.Bytes(), err
}
