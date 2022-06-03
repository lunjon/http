package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lunjon/http/internal/complete"
)

const ()

type inputModel struct {
	input     textinput.Model
	cursor    int
	items     []string
	displayed []string

	Cursor      string
	CursorStyle lipgloss.Style
}

func newInputModel(items []string, prompt string) inputModel {
	input := textinput.NewModel()
	input.Prompt = prompt
	input.Focus()
	input.CharLimit = 150
	input.Width = 50

	return inputModel{
		cursor:      -1,
		input:       input,
		items:       items,
		displayed:   []string{},
		Cursor:      ">",
		CursorStyle: focusedStyle,
	}
}

func (m inputModel) Init() tea.Cmd {
	return nil
}

func (m inputModel) Value() string {
	return m.input.Value()
}

func (m *inputModel) focusList() {
	if !m.input.Focused() {
		return
	}
	m.cursor = 0
	m.input.Blur()
}

func (m *inputModel) focusInput() tea.Cmd {
	if m.input.Focused() {
		return nil
	}
	m.cursor = -1
	return m.input.Focus()
}

func (m *inputModel) setInputValue(value string) {
	m.input.SetValue(value)
	m.input.SetCursor(len(value))
}

func (m inputModel) inputFocused() bool {
	return m.cursor < 0 && m.input.Focused()
}

func (m inputModel) listFocused() bool {
	return m.cursor > -1 && !m.input.Focused()
}

func (m inputModel) Update(msg tea.Msg) (inputModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, upBinding):
			if m.listFocused() && m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, downBinding):
			if m.inputFocused() {
				m.focusList()
			} else if m.cursor < len(m.displayed)-1 {
				m.cursor++
			}

		case key.Matches(msg, autocompleteBinding):
			if m.listFocused() {
				value := m.displayed[m.cursor]
				m.setInputValue(value)
			} else {
				var text string
				matches := []string{}

				if len(m.displayed) == 1 {
					text = m.displayed[0]
				} else {
					text, matches = complete.Complete(m.input.Value(), m.items)
				}

				m.setInputValue(text)
				m.displayed = matches
			}
		case key.Matches(msg, confirmBinding):
			if m.listFocused() {
				value := m.displayed[m.cursor]
				m.setInputValue(value)
				m.focusInput()
				m.displayed = []string{}
			}
		default:
			m.focusInput()
			m.input, cmd = m.input.Update(msg)
			input := strings.TrimSpace(m.input.Value())
			if input != "" {
				_, matches := complete.Complete(input, m.items)
				// If we have one match and it is equal to the input,
				// render an empty suggestion list
				if len(matches) == 1 && m.input.Value() == matches[0] {
					m.displayed = []string{}
				} else {
					m.displayed = matches
				}
			} else {
				m.displayed = []string{}
			}
		}
	}

	return m, cmd
}

func (m inputModel) View() string {
	b := strings.Builder{}

	b.WriteString(m.input.View())
	b.WriteString("\n")

	limit := listLimit
	if len(m.displayed) < limit {
		limit = len(m.displayed)
	}

	for i, u := range m.displayed[:limit] {
		cursor := " "
		if i == m.cursor {
			cursor = m.CursorStyle.Render(m.Cursor)
		}
		b.WriteString(fmt.Sprintf(" %s %s\n", cursor, u))
	}

	return b.String()
}
