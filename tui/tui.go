/*
tui -- module responsible for everything about the text user interface.

Vision
======
The user is taken through a number of steps when creating a request:

	Method  --->  URL  ---> [Headers] ---> [Body] ---> Request

Steps:
	Method: select HTTP method from a list
		todo:	* add text input field with search
	URL: enter the URL
		todo:	* allow auto completion from aliases
	Headers: add headers
		todo:	* auto completion for common/known headers
				* auto-completion for the values as well

At any point the user can pres <BUTTON> to select client options, such as:
 - timeout
 - certificate
*/
package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lunjon/http/internal/format"
)

const (
	confirmButtonText = " [ confirm ] "
)

var (
	noStyle        = lipgloss.NewStyle()
	boldStyle      = lipgloss.NewStyle().Bold(true)
	errorStyle     = boldStyle.Copy().Foreground(lipgloss.Color("1"))
	confirmedStyle = boldStyle.Copy().Foreground(lipgloss.Color("10"))
	focusedStyle   = boldStyle.Copy().Foreground(lipgloss.Color("14"))
	blurredStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

	styler = format.NewStyler()

	quitKeyBinding = key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	)
)

func Start(urls []string) error {
	p := tea.NewProgram(initialModel(urls))
	return p.Start()
}

func initialModel(urls []string) root {
	inner := initialMethodModel(state{}, urls)
	return root{
		inner: inner,
		help:  newHelpModel(),
	}
}

type root struct {
	inner tea.Model
	help  tea.Model
}

func (r root) Init() tea.Cmd {
	return nil
}

func (m root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, quitKeyBinding) {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.inner, cmd = m.inner.Update(msg)
	cmds := []tea.Cmd{cmd}

	m.help, cmd = m.help.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m root) View() string {
	s := m.inner.View()
	s += "\n"
	return s + m.help.View()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"%s: %s\n",
			errorStyle.Render("error"),
			err,
		)
		os.Exit(1)
	}
}
