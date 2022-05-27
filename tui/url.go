package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunjon/http/internal/complete"
	"github.com/lunjon/http/internal/util"
)

const (
	listLimit = 10
)

type urlModel struct {
	state   state
	input   textinput.Model
	urls    []string
	matches []string
}

func initialURLModel(state state, urls []string) urlModel {
	input := textinput.NewModel()
	input.Prompt = ""
	input.Focus()
	input.CharLimit = 150
	input.Width = 50

	return urlModel{
		state:   state,
		input:   input,
		urls:    urls,
		matches: []string{},
	}
}

func (m urlModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m urlModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.input, cmd = m.input.Update(msg)

	input := strings.TrimSpace(m.input.Value())
	if input != "" {
		_, matches := complete.Complete(m.input.Value(), m.urls)
		// If we have one match and it is equal to the input,
		// render an empty suggestion list
		if len(matches) == 1 && m.input.Value() == matches[0] {
			m.matches = []string{}
		} else {
			m.matches = matches
		}
	} else {
		m.matches = []string{}
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
			method := m.state.method.Value()
			state := m.state.setURL(m.input.Value())
			if !util.Contains([]string{"post", "put", "patch"}, strings.ToLower(method)) {
				return initialHeadersModel(state), nil
			}
			return initialBodyModel(state), nil
		}
	}

	return m, cmd
}

func (m urlModel) View() string {
	s := m.state.render()
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
