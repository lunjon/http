package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunjon/http/internal/client"
	"github.com/lunjon/http/internal/util"
)

const (
	listLimit = 10
)

type urlModel struct {
	help  tea.Model
	state state
	input inputModel
	urls  []string
}

func initialURLModel(state state, urls []string) urlModel {
	help := newHelp(keyMap{
		short: []key.Binding{},
		full: [][]key.Binding{
			{autocompleteBinding},
			{upBinding, downBinding},
		},
	})

	return urlModel{
		help:  help,
		state: state,
		input: newInputModel(urls, ""),
		urls:  urls,
	}
}

func (m urlModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m urlModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, confirmBinding):
			url, err := client.ParseURL(m.input.Value(), nil)
			if err == nil {
				state := m.state.setURL(url.String())
				method := m.state.method.Value()
				if !util.Contains([]string{"post", "put", "patch"}, strings.ToLower(method)) {
					return initialHeadersModel(state), nil
				}
				return initialBodyModel(state), nil
			}
		}
	}

	m.help, _ = m.help.Update(msg)
	return m, cmd
}

func (m urlModel) View() string {
	b := strings.Builder{}
	b.WriteString(m.state.render())

	renderQuery(&b, "URL: ")
	b.WriteString(m.input.View())
	b.WriteString("\n")

	b.WriteString(m.help.View())
	return b.String()
}
