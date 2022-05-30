package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunjon/http/internal/client"
	"github.com/lunjon/http/internal/complete"
	"github.com/lunjon/http/internal/style"
	"github.com/lunjon/http/internal/util"
)

const (
	listLimit = 10
)

type urlModel struct {
	help    tea.Model
	state   state
	input   textinput.Model
	urls    []string
	matches []string
}

func initialURLModel(state state, urls []string) urlModel {
	help := newHelp(keyMap{
		short: []key.Binding{},
		full: [][]key.Binding{
			{autocompleteBinding},
			{upBinding, downBinding},
		},
	})

	input := textinput.NewModel()
	input.Prompt = ""
	input.Focus()
	input.CharLimit = 150
	input.Width = 50

	return urlModel{
		help:    help,
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
		switch {
		case key.Matches(msg, autocompleteBinding):
			var text string
			matches := []string{}

			if len(m.matches) == 1 {
				text = m.matches[0]
			} else {
				text, matches = complete.Complete(m.input.Value(), m.urls)
			}

			m.input.SetValue(text)
			m.input.SetCursor(len(text))
			m.matches = matches
		case key.Matches(msg, confirmBinding):
			url, err := client.ParseURL(m.input.Value(), nil)
			if err == nil {
				state := m.state.setURL(url.String())
				method := m.state.method.Value()
				if !util.Contains([]string{"post", "put", "patch"}, strings.ToLower(method)) {
					return initialHeadersModel(state), nil
				}
				return initialBodyModel(state), nil
			}
		}
	}

	m.help, _ = m.help.Update(msg)
	return m, cmd
}

func (m urlModel) View() string {
	b := strings.Builder{}
	b.WriteString(m.state.render())

	renderQuery(&b, "URL: ")
	b.WriteString(m.input.View())
	b.WriteString("\n")

	// Display parsed value
	value := strings.TrimSpace(m.input.Value())
	if value != "" {
		line := "   "

		url, err := client.ParseURL(value, nil)
		if err != nil {
			line += style.RedB.Render("X")
			line += " "
			line += blurredStyle.Render(value)
		} else {
			line += confirmedStyle.Render(okIcon)
			line += " "
			line += blurredStyle.Render(url.String())
		}
		b.WriteString(line)
	}

	b.WriteString("\n")

	// Only render top matches
	limit := listLimit
	if len(m.matches) < limit {
		limit = len(m.matches)
	}

	for _, u := range m.matches[:limit] {
		b.WriteString(fmt.Sprintf("   %s\n", u))
	}

	b.WriteString(m.help.View())
	return b.String()
}
