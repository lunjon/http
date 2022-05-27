package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type selection int

const (
	choiceEditor selection = iota
	choiceFile   selection = iota
	choiceSkip   selection = iota
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
	method string
	url    string

	cursor  int
	choices []choice
}

func initialBodyModel(method, url string) bodyModel {
	editor, ok := os.LookupEnv("EDITOR")
	if !ok {
		editor = "vim"
	}

	choices := []choice{
		{"e", "Open editor (" + editor + ")"},
		{"f", "Search files"},
		{"s", "Skip"},
	}

	return bodyModel{
		method:  method,
		url:     url,
		choices: choices,
	}
}

func (m bodyModel) Init() tea.Cmd {
	return nil
}

func (m bodyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < 2 {
				m.cursor++
			}
		case "e", "E":
			return m, choiceCmd(choiceEditor)
		case "f", "F":
			return m, choiceCmd(choiceFile)
		case "s", "S":
			return m, choiceCmd(choiceSkip)
		case "enter":
			return m, choiceCmd(selection(m.cursor))
		}
	case selection:
		switch msg {
		case choiceEditor:
			return initialHeadersModel(m.method, m.url), nil
		case choiceFile:
			return initialHeadersModel(m.method, m.url), nil
		case choiceSkip:
			return initialHeadersModel(m.method, m.url), nil
		}
	}
	return m, nil
}

func (m bodyModel) View() string {
	b := &strings.Builder{}

	for index, c := range m.choices {
		b.WriteString("  ")
		b.WriteString(c.render(m.cursor == index) + "\n")
	}

	return b.String()
}

func choiceCmd(s selection) func() tea.Msg {
	return func() tea.Msg {
		return s
	}
}
