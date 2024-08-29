package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	width  int
	height int
	ready  bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

	}

	return m, nil
}

func (m model) View() string {
	style := lipgloss.NewStyle().Width(m.width).Height(m.height).Background(lipgloss.Color("63"))
	return style.Render("This is a fullscreen application")
}

// func (m model) View() string {
// 	if !m.ready {
// 		return "loading..."
// 	}
//
// 	windowStyle := lipgloss.NewStyle().
// 		Width(m.width/2).
// 		Height(m.height/2).
// 		Padding(1, 2).
// 		Border(lipgloss.RoundedBorder()).
// 		BorderForeground(lipgloss.Color("63")).
// 		Align(lipgloss.Center)
//
// 	content := "Simple window"
//
// 	return windowStyle.Render(content)
// }

func main() {
	// save terminal state
	saveState := exec.Command("tput", "smcup")
	saveState.Stdout = os.Stdout
	saveState.Run()

	// Ensure the terminal state is restored when the program exits
	defer func() {
		restoreState := exec.Command("tput", "rmcup")
		restoreState.Stdout = os.Stdout
		restoreState.Run()
	}()

	p := tea.NewProgram(model{})

	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
