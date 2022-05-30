package tui

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lunjon/http/internal/style"
)

const (
	confirmButtonText = " [ confirm ] "
	okIcon            = ""
)

var (
	noStyle        = lipgloss.NewStyle()
	boldStyle      = lipgloss.NewStyle().Bold(true)
	errorStyle     = style.RedB
	confirmedStyle = boldStyle.Copy().Foreground(lipgloss.Color("10"))
	focusedStyle   = boldStyle.Copy().Foreground(lipgloss.Color("14"))
	blurredStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

	queryPrefix     = blurredStyle.PaddingLeft(1).PaddingRight(1)
	askFocusedStyle = focusedStyle.Copy().Bold(true).PaddingLeft(2)
	askBlurredStyle = blurredStyle.Copy().Bold(true).PaddingLeft(2)
)

func Start(urls []string) error {
	inner := initialMethodModel(state{}, urls)
	p := tea.NewProgram(root{inner})
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

type root struct {
	inner tea.Model
}

func (r root) Init() tea.Cmd {
	return nil
}

func (r root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, quitBinding) {
			return r, tea.Quit
		}
	}

	var cmd tea.Cmd
	r.inner, cmd = r.inner.Update(msg)
	return r, cmd
}

func (r root) View() string {
	return r.inner.View()
}

func renderQuery(w io.Writer, name string) {
	fmt.Fprintf(
		w,
		"%s%s",
		queryPrefix.Render("?"),
		boldStyle.Render(name),
	)
}
