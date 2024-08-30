package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdkingsbury/americano/pkg/components"
	"github.com/jdkingsbury/americano/pkg/panes"
)

type PaneModel struct {
	styles     components.PaneStyles
	showBottom bool
	width      int
	height     int
}

func NewModel() PaneModel {
	return PaneModel{}
}

func (m PaneModel) Init() tea.Cmd {
	return nil
}

func (m PaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.styles = components.CreateStyles(m.width, m.height)
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

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("206")) // Change color as needed

func (m PaneModel) View() string {
	s := m.styles

	// Define the titles for each pane
	topLeftTitle := titleStyle.Render("Top Left Pane")
	bottomLeftTitle := titleStyle.Render("Bottom Left Pane")
	mainPaneTitle := titleStyle.Render("Main Pane")
	bottomPaneTitle := titleStyle.Render("Bottom Pane")

	// Render the panes
	topLeftPane := panes.TopLeftPane(s.TopLeftPane)
	bottomLeftPane := panes.BottomLeftPane(s.BottomLeftPane)
	mainPane := panes.MainPane(s.MainPane)

	var bottomPane string
	if m.showBottom {
		bottomPane = panes.BottomPane(s.BottomPane)
	} else {
		bottomPane = ""
	}

	// Arrange Panes
	leftSide := lipgloss.JoinVertical(lipgloss.Top, topLeftTitle, topLeftPane, bottomLeftTitle, bottomLeftPane)
	rightSide := lipgloss.JoinVertical(lipgloss.Top, mainPaneTitle, mainPane)

	layout := lipgloss.JoinHorizontal(lipgloss.Top, leftSide, rightSide)

	if m.showBottom {
		layout = lipgloss.JoinVertical(lipgloss.Top, layout, bottomPaneTitle, bottomPane)
	}

	return layout
}
