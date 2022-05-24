package types

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"strings"
)

// Taber is used to output lines in a tabular format.
type Taber struct {
	linePrefix string
	buf        *bytes.Buffer
	tw         *tabwriter.Writer
}

// NewTaber return a Taber.
// linePrefix sets a prefix for each line.
func NewTaber(linePrefix string) *Taber {
	buf := bytes.NewBuffer(nil)
	w := tabwriter.NewWriter(buf, 0, 0, 4, ' ', 0)
	return &Taber{
		linePrefix: linePrefix,
		buf:        buf,
		tw:         w,
	}
}

// Writef bypasses the tabwriter and writes to the
// underlying buffer. Useful for setting e.g. a heading.
func (t *Taber) Writef(fmts string, args ...any) {
	fmt.Fprintf(t.buf, fmts, args...)
}

func (t *Taber) WriteLine(values ...string) {
	s := strings.Join(values, "\t")
	fmt.Fprintf(t.tw, "%s%s\n", t.linePrefix, s)
}

func (t *Taber) String() string {
	t.tw.Flush()
	return t.buf.String()
}
