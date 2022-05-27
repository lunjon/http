package tui

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunjon/http/internal/types"
)

type result struct {
	res *http.Response
	err error
}

type requestModel struct {
	help   tea.Model
	state  state
	client *http.Client
	result types.Option[result]
}

func initialRequestModel(state state) requestModel {
	keys := keyMap{
		short: []key.Binding{configBinding, helpToggleBinding},
		full: [][]key.Binding{
			{configBinding},
			{helpToggleBinding, quitBinding},
		},
	}
	help := newHelp(keys)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	return requestModel{
		help:   help,
		state:  state,
		client: client,
		result: types.Option[result]{},
	}
}

func (m requestModel) Init() tea.Cmd {
	return nil
}

func (m requestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, quitBinding):
			return m, tea.Quit
		case key.Matches(msg, configBinding):
			return m, m.sendRequest
		}
	case result:
		m.result = m.result.Set(msg)
		return m, tea.Quit
	}

	m.help, _ = m.help.Update(msg)
	return m, nil
}

func (m requestModel) View() string {
	var b strings.Builder

	if m.result.IsSome() {
		// Render response

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
	} else {
		b.WriteString(m.state.render())
		b.WriteString("\n\n")
		b.WriteString(focusedStyle.Render("[ Send request? ]"))
	}

	b.WriteString("\n")
	b.WriteString(m.help.View())
	return b.String()
}

func (m requestModel) sendRequest() tea.Msg {
	method := m.state.method.Value()
	url := m.state.url.Value()

	var body io.Reader
	if m.state.body.IsSome() {
		body = bytes.NewReader(m.state.body.Value())
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return result{err: err}
	}

	if m.state.headers.IsSome() {
		req.Header = m.state.headers.Value()
	}

	res, err := m.client.Do(req)
	return result{res: res, err: err}
}
