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

TODO
====
 - Use lipgloss for styling

*/
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunjon/http/format"
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

func Start() error {
	p := tea.NewProgram(initialModel())
	return p.Start()
}

func initialModel() root {
	inner := initialMethodModel(format.NewStyler())
	return root{
		inner: inner,
	}
}

type root struct {
	inner tea.Model
}

func (r root) Init() tea.Cmd {
	return nil
}

func (m root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.inner, cmd = m.inner.Update(msg)
	return m, cmd
}
func (m root) View() string {
	s := "<header>\n\n"
	s += m.inner.View()
	return s + "\n\nPress ctrl+c to quit."
}
