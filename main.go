package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type styles struct {
	topLeftPane    lipgloss.Style
	bottomLeftPane lipgloss.Style
	bottomPane     lipgloss.Style
	mainPane       lipgloss.Style
}

func defaultStyles(width, height int) styles {
	topLeftHeight := height / 3
	bottomLeftHeight := height / 4
	mainPaneWidth := width - 35

	s := styles{
		topLeftPane:    lipgloss.NewStyle().Width(30).Height(topLeftHeight).Border(lipgloss.RoundedBorder()).Padding(1),
		bottomLeftPane: lipgloss.NewStyle().Width(30).Height(bottomLeftHeight).Border(lipgloss.RoundedBorder()).Padding(1),
		bottomPane:     lipgloss.NewStyle().Width(width - 3).Height(height / 4).Border(lipgloss.RoundedBorder()).Padding(1),
		mainPane:       lipgloss.NewStyle().Width(mainPaneWidth).Height(topLeftHeight + bottomLeftHeight + 2).Border(lipgloss.RoundedBorder()).Padding(1),
	}
	return s
}

type model struct {
	styles     styles
	showBottom bool
	width      int
	height     int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.styles = defaultStyles(m.width, m.height)
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			m.showBottom = true // Used for testing but will be used for displaying query results
		}
	}

	return m, nil
}

func (m model) View() string {
	s := m.styles

	// Render the panes
	topLeftPane := s.topLeftPane.Render("Top Left Pane")
	bottomLeftPane := s.bottomLeftPane.Render("Bottom Left Pane")
	mainPane := s.mainPane.Render("Main Pane")

	var bottomPane string
	if m.showBottom {
		bottomPane = s.bottomPane.Render("Bottom Pane")
	} else {
		bottomPane = ""
	}

	// Arrange Panes
	leftSide := lipgloss.JoinVertical(lipgloss.Top, topLeftPane, bottomLeftPane)
	rightSide := lipgloss.JoinVertical(lipgloss.Top, mainPane)

	layout := lipgloss.JoinHorizontal(lipgloss.Top, leftSide, rightSide)

	if m.showBottom {
		layout = lipgloss.JoinVertical(lipgloss.Top, layout, bottomPane)
	}

	return layout
}

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

	// styles := defaultStyles()

	p := tea.NewProgram(model{}, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
