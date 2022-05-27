package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunjon/http/internal/config"
	"github.com/lunjon/http/internal/util"
)

type selection int

const (
	choiceEditor selection = iota
	choiceFile   selection = iota
	choiceSkip   selection = iota
)

var (
	editorBinding = key.NewBinding(
		key.WithKeys("e", "E"),
		key.WithHelp("e", "Open editor"),
	)
	fileBinding = key.NewBinding(
		key.WithKeys("f", "F"),
		key.WithHelp("f", "Search files"),
	)
	skipBinding = key.NewBinding(
		key.WithKeys("s", "S"),
		key.WithHelp("s", "Skip body"),
	)
)

type choice struct {
	key  string
	text string
}

func (c choice) render(focused bool) string {
	cursor := " "
	var key string
	if focused {
		key = focusedStyle.Render(c.key)
		cursor = focusedStyle.Render(">")
	} else {
		key = boldStyle.Render(c.key)
	}
	return fmt.Sprintf("%s %s  %s", cursor, key, blurredStyle.Render(c.text))
}

type bodyModel struct {
	help    tea.Model
	state   state
	cursor  int
	choices []choice
}

func initialBodyModel(state state) bodyModel {
	keys := keyMap{
		short: []key.Binding{helpToggleBinding},
		full: [][]key.Binding{
			{upBindingV, downBindingV},
			{quitBinding, helpToggleBinding},
		},
	}
	help := newHelp(keys)

	choices := []choice{
		{"e", "Open editor (" + config.Editor + ")"},
		{"f", "Search files"},
		{"s", "Skip"},
	}

	return bodyModel{
		help:    help,
		state:   state,
		choices: choices,
	}
}

func (m bodyModel) Init() tea.Cmd {
	return nil
}

func (m bodyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, upBindingV):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, downBindingV):
			if m.cursor < 2 {
				m.cursor++
			}
		case key.Matches(msg, editorBinding):
			return m, choiceCmd(choiceEditor)
		case key.Matches(msg, fileBinding):
			return m, choiceCmd(choiceFile)
		case key.Matches(msg, skipBinding):
			return m, choiceCmd(choiceSkip)
		case key.Matches(msg, confirmBinding):
			return m, choiceCmd(selection(m.cursor))
		}
	case selection:
		switch msg {
		case choiceEditor:
			content, err := util.OpenEditor(config.Editor)
			checkError(err)
			state := m.state.setBody(content)
			return initialHeadersModel(state), nil
		case choiceFile:
			return initialFileSearchModel(m.state), nil
		case choiceSkip:
			return initialHeadersModel(m.state), nil
		}
	}

	m.help, _ = m.help.Update(msg)
	return m, nil
}

func (m bodyModel) View() string {
	b := &strings.Builder{}
	b.WriteString(m.state.render())

	b.WriteString("Body:\n")
	for index, c := range m.choices {
		b.WriteString("  ")
		b.WriteString(c.render(m.cursor == index) + "\n")
	}

	b.WriteString(m.help.View())
	return b.String()
}

func choiceCmd(s selection) func() tea.Msg {
	return func() tea.Msg {
		return s
	}
}
