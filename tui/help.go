package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	helpToggleBinding = key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("ctrl+h", "toggle help"),
	)
	quitBinding = key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	)
	confirmBinding = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("/enter", "confirm"),
	)
	upBinding = key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	)
	upBindingV = key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	)
	downBinding = key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	)
	downBindingV = key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	)
	leftBinding = key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "move left"),
	)
	leftBindingV = key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	)
	rightbinding = key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "move right"),
	)
	rightbindingV = key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	)
	autocompleteBinding = key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "autocomplete"),
	)
)

type keyMap struct {
	short []key.Binding
	full  [][]key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return k.short
}

func (k keyMap) FullHelp() [][]key.Binding {
	return k.full
}

type helpModel struct {
	keys keyMap
	help help.Model
}

func newHelp(keys keyMap) helpModel {
	return helpModel{
		keys: keys,
		help: help.New(),
	}
}

func (m helpModel) Init() tea.Cmd {
	return nil
}

func (m helpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, helpToggleBinding):
			m.help.ShowAll = !m.help.ShowAll
		}
	}
	return m, nil
}

func (m helpModel) View() string {
	return "\n" + m.help.View(m.keys)
}
