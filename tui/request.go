package tui

import (
	"fmt"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunjon/http/internal/util"
)

type response *http.Response

type requestModel struct {
	method  string
	url     string
	headers http.Header
	res     response
}

func initialRequestModel(method, url string, headers http.Header) requestModel {
	return requestModel{
		method:  method,
		url:     url,
		headers: headers,
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
			return m, tea.Quit // TODO: send request
		}
	case response:
		m.res = msg

	}

	return m, nil
}

func (m requestModel) View() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Method:  %s\n", styler.WhiteB(m.method)))
	b.WriteString(fmt.Sprintf("URL:     %s\n", styler.WhiteB(m.url)))

	if len(m.headers) > 0 {
		b.WriteString("Headers:\n")
		taber := util.NewTaber("  - ")
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

	b.WriteString("\nSend request?")
	return b.String()
}
