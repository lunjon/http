package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunjon/http/complete"
)

const (
	listLimit = 5
)

type urlModel struct {
	input   textinput.Model
	method  string
	urls    []string
	matches []string
}

func initialURLModel(method string, urls []string) urlModel {
	input := textinput.NewModel()
	input.Prompt = ""
	input.Focus()
	input.CharLimit = 150
	input.Width = 50

	return urlModel{
		method:  method,
		input:   input,
		urls:    urls,
		matches: urls,
	}
}

func (m urlModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m urlModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	matches := complete.Matches(m.input.Value(), m.urls)

	// If we have one match and it is equal to the input,
	// render an empty suggestion list
	if len(matches) == 1 && m.input.Value() == matches[0] {
		m.matches = []string{}
	} else {
		m.matches = matches
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			var text string
			matches := []string{}

			if len(m.matches) == 1 {
				// We have a single match, use that
				text = m.matches[0]
			} else {
				text, matches = complete.Complete(m.input.Value(), m.urls)
			}

			m.input.SetValue(text)
			m.input.SetCursor(len(text))
			m.matches = matches
		case "enter":
			url := m.input.Value()
			return initialRequestModel(m.method, url), nil
		}
	}

	return m, cmd
}

func (m urlModel) View() string {
	s := fmt.Sprintf("Method: %s\n\n", m.method)
	s += fmt.Sprintf("URL: %s\n", m.input.View())

	// Only render top matches
	limit := listLimit
	if len(m.matches) < limit {
		limit = len(m.matches)
	}

	for _, u := range m.matches[:limit] {
		s += fmt.Sprintf("  %s\n", u)
	}

	return s
}
