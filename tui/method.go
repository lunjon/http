package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunjon/http/internal/client"
	"github.com/lunjon/http/internal/types"
)

type methodModel struct {
	cursor  int
	methods []string
	method  types.Option[string]
	urls    []string
}

func initialMethodModel(urls []string) methodModel {
	return methodModel{
		methods: client.SupportedMethods,
		method:  types.Option[string]{},
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
			return initialURLModel(selectedMethod, m.urls), nil
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
