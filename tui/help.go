package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	defaultToggleBinding = key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	)
	inputToggleBinding = key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("ctrl+h", "toggle help"),
	)
	quitBinding = key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc/ctrl+c", "quit"),
	)
	configBinding = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("/enter", "confirm"),
	)
	upBinding = key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑", "move up"),
	)
	upBindingV = key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	)
	downBinding = key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓", "move down"),
	)
	downBindingV = key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	)
	leftBinding = key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←h", "move left"),
	)
	leftBindingV = key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	)
	rightbinding = key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
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
	keys   keyMap
	help   help.Model
	toggle key.Binding
}

func newHelp(toggle key.Binding, keys keyMap) helpModel {
	return helpModel{
		toggle: toggle,
		keys:   keys,
		help:   help.New(),
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
		case key.Matches(msg, m.toggle):
			m.help.ShowAll = !m.help.ShowAll
		}
	}

	return m, nil
}

func (m helpModel) View() string {
	return "\n" + m.help.View(m.keys)
}
