package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunjon/http/internal/client"
)

type methodModel struct {
	state   state
	cursor  int
	methods []string
	urls    []string
}

func initialMethodModel(state state, urls []string) methodModel {
	return methodModel{
		state:   state,
		methods: client.SupportedMethods,
		urls:    urls,
	}
}

func (m methodModel) Init() tea.Cmd {
	return nil
}

func (m methodModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
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
			state := m.state.setMethod(selectedMethod)
			return initialURLModel(state, m.urls), nil
		}
	}
	return m, nil
}

func (m methodModel) View() string {
	s := "Method: \n"

	for i, choice := range m.methods {
		cursor := " "
		if m.cursor == i {
			cursor = focusedStyle.Render(">")
			choice = focusedStyle.Render(choice)
		}

		s += fmt.Sprintf("  %s %s\n", cursor, choice)
	}

	return s
}
