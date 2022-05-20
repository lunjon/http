package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunjon/http/client"
	"github.com/lunjon/http/complete"
	"github.com/lunjon/http/format"
)

type methodModel struct {
	cursor  int
	methods []string
	matches []string
	method  Option[string]
	input   textinput.Model
}

func initialMethodModel() methodModel {
	input := textinput.NewModel()
	input.Prompt = ": "
	input.Focus()
	input.CharLimit = 150
	input.Width = 50

	return methodModel{
		methods: client.SupportedMethods,
		matches: client.SupportedMethods,
		method:  Option[string]{},
		input:   input,
	}
}

func (m methodModel) Init() tea.Cmd {
	return nil
}

func (m methodModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.matches)-1 {
				m.cursor++
			}

		case "enter", " ":
			selectedMethod := m.matches[m.cursor]
			return newURLModel(selectedMethod), nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	m.matches = complete.Matches(m.input.Value(), m.methods)
	if m.cursor > len(m.matches)-1 {
		m.cursor = len(m.matches) - 1
	}
	return m, cmd
}

func (m methodModel) View() string {
	s := fmt.Sprintf("Select method%s\n\n", m.input.View())

	for i, choice := range m.matches {
		ch := choice
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = "*" // cursor!
			ch = format.WhiteB(ch)
		}

		s += fmt.Sprintf("%s %s\n", cursor, ch)
	}

	return s
}
