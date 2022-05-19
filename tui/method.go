package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunjon/http/client"
)

type methodModel struct {
	cursor  int
	methods []string
	method  Option[string]
}

func initialModel() methodModel {
	return methodModel{
		methods: client.SupportedMethods,
		method:  Option[string]{},
	}
}

func (m methodModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
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
			if m.cursor < len(m.methods)-1 {
				m.cursor++
			}

		case "enter", " ":
			selectedMethod := m.methods[m.cursor]
			return newURLModel(selectedMethod), nil
		}
	}
	return m, nil
}

func (m methodModel) View() string {
	s := "Select method:\n\n"

	for i, choice := range m.methods {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	return s + "\nPress q to quit.\n"
}
