package tui

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lunjon/http/internal/util"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = focusedStyle.Copy().PaddingLeft(2)
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingBottom(1)

	filterKeyBinding = key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	)
)

type item string

func (i item) FilterValue() string { return string(i) }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i := listItem.(item)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}
	fmt.Fprintf(w, fn(string(i)))
}

type fileSearchModel struct {
	help   tea.Model
	state  state
	list   list.Model
	items  []item
	choice string
}

func initialFileSearchModel(state state) fileSearchModel {
	paths, err := util.WalkDir(".")
	if err != nil {
		panic(err)
	}

	items := []list.Item{}
	for _, p := range paths {
		items = append(items, item(p))
	}

	l := list.New(items, itemDelegate{}, 20, 15)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.PaginationStyle = paginationStyle
	l.FilterInput.PromptStyle = boldStyle
	l.FilterInput.CursorStyle = noStyle

	help := newHelp(keyMap{
		short: []key.Binding{upBindingV, downBindingV, filterKeyBinding},
		full: [][]key.Binding{
			{},
			{upBindingV, downBindingV, filterKeyBinding},
		},
	})

	return fileSearchModel{
		help:  help,
		state: state,
		list:  l,
	}

}

func (m fileSearchModel) Init() tea.Cmd {
	return nil
}

func (m fileSearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, confirmBinding):
			i := m.list.SelectedItem().(item)
			filepath := string(i)
			b, err := os.ReadFile(filepath)
			checkError(err)

			state := m.state.setBody(filepath, b)
			return initialHeadersModel(state), nil
		}
	}

	m.help, _ = m.help.Update(msg)

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m fileSearchModel) View() string {
	s := m.state.render()
	s += "Select body file:\n"
	return s + m.list.View() + m.help.View()
}
