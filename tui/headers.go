package tui

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lunjon/http/internal/client"
	"github.com/lunjon/http/internal/complete"
	"github.com/lunjon/http/internal/types"
)

var (
	headerKeyStyle   = lipgloss.NewStyle().Bold(true)
	headerValueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
)

type suggestions struct {
	indent string
	values []string
}

type headersModel struct {
	method string
	url    string

	buttonFocus bool
	nameInput   textinput.Model
	valueInput  textinput.Model
	headers     http.Header
	headerNames []string // Used to preserve order

	knownHeaderNames   []string
	suggest            types.Option[suggestions]
	nameSuggestIndent  string
	valueSuggestIndent string
}

func newInput(prompt string, limit int) textinput.Model {
	input := textinput.New()
	input.Prompt = prompt
	input.CharLimit = limit
	return input
}

func initialHeadersModel(method, url string) headersModel {
	name := newInput("Name: ", 64)
	name.Focus()
	name.PromptStyle = focusedStyle
	name.Width = 24

	value := newInput("Value: ", 128)

	knownHeaderNames := []string{}
	for n := range client.KnownHTTPHeaders {
		knownHeaderNames = append(knownHeaderNames, n)
	}

	return headersModel{
		method:             method,
		url:                url,
		nameInput:          name,
		valueInput:         value,
		headers:            http.Header{},
		knownHeaderNames:   knownHeaderNames,
		suggest:            types.Option[suggestions]{},
		nameSuggestIndent:  strings.Repeat(" ", len(" Name: ")),
		valueSuggestIndent: strings.Repeat(" ", len(" Value: ")+len(" Name: ")+name.Width),
	}
}

func (m headersModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *headersModel) addHeader(name, value string) {
	if v := m.headers.Get(name); v == "" {
		k := http.CanonicalHeaderKey(name)
		m.headerNames = append(m.headerNames, k)
	}
	m.headers.Add(name, value)
}

func (m *headersModel) setFocus(name, value, button bool) tea.Cmd {
	var cmd tea.Cmd
	if name {
		cmd = m.nameInput.Focus()
		m.nameInput.PromptStyle = focusedStyle

		m.valueInput.Blur()
		m.valueInput.PromptStyle = noStyle

		m.buttonFocus = false

	} else if value {
		m.nameInput.Blur()
		m.nameInput.PromptStyle = noStyle

		cmd = m.valueInput.Focus()
		m.valueInput.PromptStyle = focusedStyle

		m.buttonFocus = false
	} else if button {
		m.nameInput.Blur()
		m.nameInput.PromptStyle = noStyle

		m.valueInput.Blur()
		m.valueInput.PromptStyle = noStyle

		m.buttonFocus = true
	}

	m.unsetSuggestions()
	return cmd
}

func (m *headersModel) setSuggestions(indent string, values []string) {
	opt := types.Option[suggestions]{}
	m.suggest = opt.Set(suggestions{indent, values})
}

func (m *headersModel) unsetSuggestions() {
	m.suggest = types.Option[suggestions]{}
}

func (m headersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.nameInput.Focused() {
				prefix, matches := complete.Complete(m.nameInput.Value(), m.knownHeaderNames)
				if prefix != "" && len(matches) > 0 {
					if len(matches) > 1 {
						m.setSuggestions(m.nameSuggestIndent, matches)
					} else {
						m.unsetSuggestions()
					}
					m.nameInput.SetValue(prefix)
					m.nameInput.SetCursor(len(prefix))
				}
			} else if m.valueInput.Focused() {
				headerName := strings.TrimSpace(m.nameInput.Value())
				if knownValues, ok := client.KnownHTTPHeaders[headerName]; ok {
					prefix, matches := complete.Complete(m.valueInput.Value(), knownValues)
					if prefix != "" && len(matches) > 0 {
						if len(matches) > 1 {
							m.setSuggestions(m.valueSuggestIndent, matches)
						} else {
							m.unsetSuggestions()
						}
						m.valueInput.SetValue(prefix)
						m.valueInput.SetCursor(len(prefix))
					}
				}
			}
		case "enter":
			if m.nameInput.Focused() || m.valueInput.Focused() {
				name := m.nameInput.Value()
				value := m.valueInput.Value()
				name = strings.TrimSpace(name)
				value = strings.TrimSpace(value)
				if name != "" && value != "" {
					m.unsetSuggestions()
					m.addHeader(name, value)
					m.nameInput.SetValue("")
					m.valueInput.SetValue("")

					cmd := m.setFocus(true, false, false)
					return m, cmd
				}
			} else if m.buttonFocus {
				// Switch to next model
				return initialRequestModel(m.method, m.url, m.headers), nil
			}

		case "up":
			if m.buttonFocus {
				cmd := m.setFocus(true, false, false)
				m.buttonFocus = false
				return m, cmd
			}
		case "down":
			if !m.buttonFocus {
				m.setFocus(false, false, true)
				return m, nil
			}
		case "right":
			if m.nameInput.Focused() {
				cmd := m.setFocus(false, true, false)
				return m, cmd
			}
		case "left":
			if m.valueInput.Focused() {
				cmd := m.setFocus(true, false, false)
				return m, cmd
			}
		}
	}

	var cmd1, cmd2 tea.Cmd

	// Handle character input and blinking
	m.nameInput, cmd1 = m.nameInput.Update(msg)
	m.valueInput, cmd2 = m.valueInput.Update(msg)
	return m, tea.Batch(cmd1, cmd2)
}

func (m headersModel) View() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Method:  %s\n", confirmedStyle.Render(m.method)))
	b.WriteString(fmt.Sprintf("URL:     %s\n", confirmedStyle.Render(m.url)))
	b.WriteString("Headers:\n")

	b.WriteString(" ")
	b.WriteString(m.nameInput.View())
	b.WriteString(m.valueInput.View())
	b.WriteString("\n")

	if m.suggest.IsSome() {
		sugg := m.suggest.Value()
		for _, v := range sugg.values {
			fmt.Fprintf(&b, "%s%s\n", sugg.indent, v)
		}
	}

	if len(m.headerNames) > 0 {
		b.WriteString("\n")
		taber := types.NewTaber("  - ")
		for _, name := range m.headerNames {
			key := headerKeyStyle.Render(name + ":")
			values := strings.Join(m.headers.Values(name), "; ")
			values = headerValueStyle.Render(values)

			taber.WriteLine(key, values)
		}
		b.WriteString(taber.String())
	}

	b.WriteString("\n")

	confirm := confirmButtonText
	if m.buttonFocus {
		confirm = focusedStyle.Render(confirm)
	} else {
		confirm = blurredStyle.Render(confirm)
	}

	b.WriteString(confirm)
	b.WriteString("\n")
	return b.String()
}
