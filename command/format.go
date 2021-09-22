package command

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/lunjon/http/client"
)

type Formatter interface {
	Format(*client.Result) ([]byte, error)
}

type DefaultFormatter struct{}

func (f DefaultFormatter) Format(r *client.Result) ([]byte, error) {
	return r.Body()
}

type NullFormatter struct{}

func (f NullFormatter) Format(*client.Result) ([]byte, error) {
	return nil, nil
}

type BriefFormatter struct{}

func (f BriefFormatter) Format(r *client.Result) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	w := tabwriter.NewWriter(buf, 0, 0, 2, ' ', tabwriter.TabIndent)

	req := r.Request()
	fmt.Fprintf(w, "Method:\t%s\n", req.Method)
	fmt.Fprintf(w, "URL:\t%s\n", req.URL.String())
	fmt.Fprintf(w, "Status:\t%s", r.Status())
	err := w.Flush()

	return buf.Bytes(), err
}
