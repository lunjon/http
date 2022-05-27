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
