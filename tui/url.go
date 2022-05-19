package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type urlModel struct {
	method string
	input  textinput.Model
}

func newURLModel(method string) urlModel {
	input := textinput.NewModel()
	input.Placeholder = "https://"
	input.Focus()
	input.CharLimit = 150
	input.Width = 50

	return urlModel{
		method: method,
		input:  input,
	}
}

func (m urlModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m urlModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			url := m.input.Value()
			fmt.Println(url)
			return m, tea.Quit
		}
	case error:
		fmt.Println("Error:", msg.Error())
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m urlModel) View() string {
	s := fmt.Sprintf("Method: %s\n\n", m.method)
	s += fmt.Sprintf("Enter URL%s\n", m.input.View())
	return s + "\n\nPress q to quit.\n"
}
