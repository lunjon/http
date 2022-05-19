package tui

import (
	"fmt"
	"log"
	"net/http"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunjon/http/client"
	"github.com/lunjon/http/logging"
)

type urlModel struct {
	method string
	input  textinput.Model
}

func newURLModel(method string) urlModel {
	input := textinput.NewModel()
	input.Placeholder = "https://"
	input.Focus()
	input.CharLimit = 150
	input.Width = 50

	return urlModel{
		method: method,
		input:  input,
	}
}

func (m urlModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m urlModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			url := m.input.Value()
			send(m.method, url)
			return m, tea.Quit
		}
	case error:
		fmt.Println("Error:", msg.Error())
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m urlModel) View() string {
	s := fmt.Sprintf("Method: %s\n\n", m.method)
	s += fmt.Sprintf("Enter URL%s\n", m.input.View())
	return s + "\n\nPress q to quit.\n"
}

func send(method, url string) {
	logger := logging.NewSilentLogger()
	httpClient := client.NewClient(&http.Client{}, logger, logger)
	u, err := client.ParseURL(url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req, err := httpClient.BuildRequest(method, u, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := httpClient.Send(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Status)
}
