package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	historyView viewport.Model // this will be scrollable
}

func initializeModel(hist string) model {
	m := model{
		historyView: viewport.New(0, 0),
	}

	m.historyView.SetContent(hist)

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	return m.historyView.View()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch e := msg.(type) {
	case tea.WindowSizeMsg:
		m.historyView.Width = e.Width
		m.historyView.Height = e.Height
	case tea.KeyMsg:
		switch e.String() {
		case "ctrl-c", "q":
			return m, tea.Quit
		}
	}

	m.historyView, cmd = m.historyView.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func main() {
	historyPath := filepath.Join(os.Getenv("HOME"), ".zsh_history")
	file, err := os.Open(historyPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var str strings.Builder

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if idx := strings.Index(line, ";"); idx != -1 {
			str.WriteString(line[idx+1:] + "\n")
		} else {
			str.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	m := initializeModel(str.String())

	program := tea.NewProgram(m, tea.WithAltScreen())

	_, err = program.Run()
	if err != nil {
		log.Fatal(err)
	}
}
