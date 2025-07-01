package main

import (
	"log"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	historyView viewport.Model // this will be scrollable
}

func initializeModel() model {
	return model{
		historyView: viewport.New(0, 0),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	return ""
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func main() {
	m := initializeModel()

	program := tea.NewProgram(m, tea.WithAltScreen())

	_, err := program.Run()
	if err != nil {
		log.Fatal(err)
	}
}
