package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
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

		case "enter":
			if selected, ok := m.list.SelectedItem().(historyItem); ok {
				return m, func() tea.Msg {
					runCommand(string(selected))
					return tea.Quit()
				}
			}
		case "y":
			if selected, ok := m.list.SelectedItem().(historyItem); ok {
				copyToClipboard(string(selected))
			}

		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.list.View()
}

func runCommand(cmdStr string) {
	cmd := exec.Command("zsh", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	_ = cmd.Run()
}

func copyToClipboard(text string) {
	var cmd *exec.Cmd

	if _, err := exec.LookPath("pbcopy"); err == nil {
		cmd = exec.Command("pbcopy")
	} else if _, err := exec.LookPath("xclip"); err == nil {
		cmd = exec.Command("xclip", "-selection", "clipboard")
	} else {
		return // no supported clipboard tool
	}

	in, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	if err := cmd.Start(); err != nil {
		return
	}

	_, _ = in.Write([]byte(text))
	_ = in.Close()
	_ = cmd.Wait()
}

func main() {
	historyPath := filepath.Join(os.Getenv("HOME"), ".zsh_history")
	file, err := os.Open(historyPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var rawLines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if idx := strings.Index(line, ";"); idx != -1 {
			rawLines = append(rawLines, line[idx+1:])
		} else {
			rawLines = append(rawLines, line)
		}
	}

	// Reverse the slice
	for i, j := 0, len(rawLines)-1; i < j; i, j = i+1, j-1 {
		rawLines[i], rawLines[j] = rawLines[j], rawLines[i]
	}

	var items []list.Item
	for _, line := range rawLines {
		items = append(items, historyItem(line))
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
