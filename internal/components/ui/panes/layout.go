package panes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PaneModel struct {
	styles     PaneStyles
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
		m.styles = CreatePaneStyles(m.width, m.height)
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
	topLeftPane := TopLeftPane(s.TopLeftPane)
	mainPane := MainPane(s.MainPane)
	bottomPane := BottomPane(s.BottomPane)

	// Arrange Panes
	leftSide := lipgloss.JoinVertical(lipgloss.Top, topLeftPane)
	rightSide := lipgloss.JoinVertical(lipgloss.Top, mainPane)

	layout := lipgloss.JoinHorizontal(lipgloss.Top, leftSide, rightSide)
	layout = lipgloss.JoinVertical(lipgloss.Top, layout, bottomPane)

	return layout
}
