package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type historyItem string

func (h historyItem) Title() string       { return string(h) }
func (h historyItem) Description() string { return "" }
func (h historyItem) FilterValue() string { return string(h) }

type model struct {
	list list.Model
}

func initializeModel(items []list.Item) model {
	const defaultWidth = 20
	l := list.New(items, list.NewDefaultDelegate(), defaultWidth, 10)
	l.Title = "Command History"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	return model{list: l}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch e := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(e.Width, e.Height)
	case tea.KeyMsg:
		switch e.String() {
		case "ctrl-c", "q":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.list.View()
}

func main() {
	historyPath := filepath.Join(os.Getenv("HOME"), ".zsh_history")
	file, err := os.Open(historyPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var items []list.Item

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if idx := strings.Index(line, ";"); idx != -1 {
			items = append(items, historyItem(line[idx+1:]))
		} else {
			items = append(items, historyItem(line))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	m := initializeModel(items)

	program := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}
