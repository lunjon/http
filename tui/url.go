package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunjon/http/complete"
)

var exampleURLs = []string{"http://localhost", "https://golang.org"}

type urlModel struct {
	input   textinput.Model
	method  string
	urls    []string
	matches []string
}

func initialURLModel(method string) urlModel {
	input := textinput.NewModel()
	input.Prompt = ""
	input.Focus()
	input.CharLimit = 150
	input.Width = 50

	return urlModel{
		method:  method,
		input:   input,
		urls:    exampleURLs,
		matches: exampleURLs,
	}
}

func (m urlModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m urlModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			text, matches := complete.Complete(m.input.Value(), m.urls)
			m.input.SetValue(text)
			m.input.SetCursor(len(text))
			m.matches = matches

		case "enter":
			url := m.input.Value()
			return initialRequestModel(m.method, url), nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.matches = complete.Matches(m.input.Value(), m.urls)

	return m, cmd
}

func (m urlModel) View() string {
	s := fmt.Sprintf("Method: %s\n\n", m.method)
	s += fmt.Sprintf("URL: %s\n", m.input.View())

	for _, u := range m.matches {
		s += fmt.Sprintf("  %s\n", u)
	}

	return s
}
