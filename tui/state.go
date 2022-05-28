package tui

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/lunjon/http/internal/types"
)

type state struct {
	method  types.Option[string]
	url     types.Option[string]
	headers types.Option[http.Header]
	body    types.Option[[]byte]
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

func (s state) setBody(b []byte) state {
	s.body = s.body.Set(b)
	return s
}

func (s state) render() string {
	var b strings.Builder
	if s.method.IsSome() {
		method := s.method.Value()
		b.WriteString(fmt.Sprintf("Method:  %s\n", confirmedStyle.Render(method)))
	}

	if s.url.IsSome() {
		url := s.url.Value()
		b.WriteString(fmt.Sprintf("URL:     %s\n", confirmedStyle.Render(url)))
	}

	if s.headers.IsSome() {
		headers := s.headers.Value()
		if len(headers) > 0 {
			b.WriteString("Headers:\n")
			taber := types.NewTaber("  - ")
			for name, values := range headers {
				key := headerKeyStyle(name + ":")
				value := headerValueStyle(strings.Join(values, "; "))
				taber.WriteLine(key, value)
			}
			b.WriteString(taber.String())
		} else {
			b.WriteString("Headers: ")
			b.WriteString(confirmedStyle.Render("[]"))
		}
	}

	return b.String()
}
