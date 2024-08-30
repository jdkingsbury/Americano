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
		}
	}

	return m, nil
}

func (m PaneModel) View() string {
	s := m.styles

	// Render the panes
	topLeftPane := panes.TopLeftPane(s.TopLeftPane)
	mainPane := panes.MainPane(s.MainPane)
	bottomPane := panes.BottomPane(s.BottomPane)

	// Arrange Panes
	leftSide := lipgloss.JoinVertical(lipgloss.Top, topLeftPane)
	rightSide := lipgloss.JoinVertical(lipgloss.Top, mainPane)

	layout := lipgloss.JoinHorizontal(lipgloss.Top, leftSide, rightSide)
	layout = lipgloss.JoinVertical(lipgloss.Top, layout, bottomPane)

	return layout
}
