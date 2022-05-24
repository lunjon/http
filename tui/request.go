package tui

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunjon/http/internal/types"
)

type result struct {
	res *http.Response
	err error
}

type requestModel struct {
	method  string
	url     string
	headers http.Header
	client  *http.Client
	result  types.Option[result]
}

func initialRequestModel(method, url string, headers http.Header) requestModel {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	return requestModel{
		method:  method,
		url:     url,
		headers: headers,
		client:  client,
		result:  types.Option[result]{},
	}
}

func (m requestModel) Init() tea.Cmd {
	return nil
}

func (m requestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, m.sendRequest
		}
	case result:
		m.result = m.result.Set(msg)
		return m, tea.Quit
	}

	return m, nil
}

func (m requestModel) View() string {
	var b strings.Builder

	if m.result.IsSome() {
		result := m.result.Value()
		if result.err != nil {
			b.WriteString(styler.RedB("Error: ") + result.err.Error())
		} else {
			var status string
			if result.res.StatusCode >= 500 {
				status = styler.RedB(result.res.Status)
			} else if result.res.StatusCode >= 400 {
				status = styler.YellowB(result.res.Status)
			} else {
				status = styler.GreenB(result.res.Status)
			}
			b.WriteString("Status: " + status)
		}

		b.WriteString("\n")
	} else {
		b.WriteString(fmt.Sprintf("Method:  %s\n", styler.WhiteB(m.method)))
		b.WriteString(fmt.Sprintf("URL:     %s\n", styler.WhiteB(m.url)))

		if len(m.headers) > 0 {
			b.WriteString("Headers:\n")
			taber := types.NewTaber("  - ")
			for name, values := range m.headers {
				key := headerKeyStyle.Render(name + ":")
				value := headerValueStyle.Render(strings.Join(values, "; "))
				taber.WriteLine(key, value)
			}
			b.WriteString(taber.String())
		} else {
			b.WriteString("Headers: ")
			b.WriteString(styler.WhiteB("[]"))
		}

		b.WriteString("\n\n")
		b.WriteString(focusedStyle.Render(" Send request?\n"))
	}

	return b.String()
}

func (m requestModel) sendRequest() tea.Msg {
	req, err := http.NewRequest(m.method, m.url, nil)
	if err != nil {
		return result{err: err}
	}
	req.Header = m.headers

	res, err := m.client.Do(req)
	return result{res: res, err: err}
}
