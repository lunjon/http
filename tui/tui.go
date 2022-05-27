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
	m := initialMethodModel(state{}, urls)
	p := tea.NewProgram(m)
	return p.Start()
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
