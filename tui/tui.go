/*
tui -- module responsible for everything about the text user interface.

Vision
======
The user is taken through a number of steps when creating a request:

	Method  --->  URL  ---> [Headers] ---> [Body]

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
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunjon/http/internal/format"
)

func init() {
	styler = format.NewStyler()
}

var (
	styler         *format.Styler
	quitKeyBinding = key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	)
)

type Option[T any] struct {
	some  bool
	value T
}

func (o Option[T]) IsSome() bool {
	return o.some
}

func (o Option[T]) IsNone() bool {
	return !o.IsSome()
}

func (o Option[T]) Value() T {
	if !o.some {
		panic("No value")
	}
	return o.value
}

func (o Option[T]) Set(value T) Option[T] {
	return Option[T]{
		some:  true,
		value: value,
	}
}

func Start(urls []string) error {
	p := tea.NewProgram(initialModel(urls))
	return p.Start()
}

func initialModel(urls []string) root {
	inner := initialMethodModel(urls)
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
