package tui

import (
	"bytes"
	"io"
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
	state  state
	client *http.Client
	result types.Option[result]
}

func initialRequestModel(state state) requestModel {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	return requestModel{
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
		b.WriteString(focusedStyle.Render(confirmButtonText))
	}

	b.WriteString("\n")
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
