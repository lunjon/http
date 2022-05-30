package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunjon/http/internal/client"
)

type methodModel struct {
	help    tea.Model
	state   state
	cursor  int
	methods []string
	urls    []string
}

func initialMethodModel(state state, urls []string) methodModel {
	help := newHelp(keyMap{
		short: []key.Binding{},
		full: [][]key.Binding{
			{},
			{upBindingV, downBindingV},
		},
	})
	return methodModel{
		help:    help,
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
		switch {
		case key.Matches(msg, upBindingV):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, downBindingV):
			if m.cursor < len(m.methods)-1 {
				m.cursor++
			}
		case key.Matches(msg, confirmBinding):
			selectedMethod := m.methods[m.cursor]
			state := m.state.setMethod(selectedMethod)
			return initialURLModel(state, m.urls), nil
		}
	}

	m.help, _ = m.help.Update(msg)
	return m, nil
}

func (m methodModel) View() string {
	b := strings.Builder{}
	renderQuery(&b, "Method:")
	b.WriteString("\n")

	for i, choice := range m.methods {
		cursor := " "
		if m.cursor == i {
			cursor = focusedStyle.Render(">")
			choice = focusedStyle.Render(choice)
		}

		b.WriteString(fmt.Sprintf("   %s %s\n", cursor, choice))

	}

	b.WriteString("\n")
	b.WriteString(m.help.View())
	return b.String()
}
