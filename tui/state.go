package tui

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lunjon/http/internal/style"
	"github.com/lunjon/http/internal/types"
)

var (
	confirmedPrefix = confirmedStyle.PaddingLeft(1).Render(okIcon + " ")
	stateFieldStyle = boldStyle.Copy().Width(8).Align(lipgloss.Left)
	stateValueStyle = style.Cyan.Copy()
)

type state struct {
	method          types.Option[string]
	url             types.Option[string]
	headers         types.Option[http.Header]
	body            types.Option[[]byte]
	bodyDescription string
}

func (s state) setMethod(m string) state {
	s.method = s.method.Set(m)
	return s
}

func (s state) setURL(u string) state {
	s.url = s.method.Set(u)
	return s
}

func (s state) setHeaders(h http.Header) state {
	s.headers = s.headers.Set(h)
	return s
}

func (s state) setBody(desc string, b []byte) state {
	s.body = s.body.Set(b)
	s.bodyDescription = desc
	return s
}

func (s state) render() string {
	var b strings.Builder
	if s.method.IsSome() {
		method := s.method.Value()
		renderStateValue(&b, "Method", method)
	}

	if s.url.IsSome() {
		url := s.url.Value()
		renderStateValue(&b, "URL", url)
	}

	if s.body.IsSome() {
		renderStateValue(&b, "Body", s.bodyDescription)
	}

	if s.headers.IsSome() {
		headers := s.headers.Value()
		if len(headers) > 0 {
			renderStateValue(&b, "Headers", "")

			taber := types.NewTaber("   - ")
			for name, values := range headers {
				key := headerKeyStyle.Render(name + ":")
				value := headerValueStyle.Render(strings.Join(values, "; "))
				taber.WriteLine(key, value)
			}
			b.WriteString(taber.String())
		} else {
			renderStateValue(&b, "Headers", "[]")
		}
	}

	return b.String()
}

func renderStateValue(w io.Writer, name, value string) {
	fmt.Fprintf(
		w,
		"%s%s %s\n",
		confirmedPrefix,
		stateFieldStyle.Render(name+":"),
		stateValueStyle.Render(value),
	)
}
