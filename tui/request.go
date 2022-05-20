package tui

import (
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

type response *http.Response

type requestModel struct {
	method string
	url    string
	res    response
}

func initialRequestModel(method, url string) requestModel {
	return requestModel{
		method: method,
		url:    url,
	}
}

func (m requestModel) Init() tea.Cmd {
	return nil
}

func (m requestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			return m, tea.Quit // TODO: send request
		}
	case response:
		m.res = msg

	}

	return m, nil
}

func (m requestModel) View() string {
	s := fmt.Sprintf("Method: %s\n", styler.WhiteB(m.method))
	s += fmt.Sprintf("URL:    %s\n\n", styler.WhiteB(m.url))

	if m.res == nil {
		s += "Send request?"
	}

	return s
}
