package panes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ResultPaneModel struct {
	styles       lipgloss.Style
	activeStyles lipgloss.Style
	width        int
	height       int
	err          error
	isActive     bool
}

func NewResultPane(width, height int) *ResultPaneModel {
	pane := &ResultPaneModel{
		width:  width,
		height: height,
		err:    nil,
	}

	pane.updateStyles()

	return pane
}

func (m *ResultPaneModel) updateStyles() {
	m.styles = lipgloss.NewStyle().
		Width(m.width - 3).
		Height(m.height / 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(iris))

	m.activeStyles = lipgloss.NewStyle().
		Width(m.width - 3).
		Height(m.height / 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(rose))
}

func (m *ResultPaneModel) Init() tea.Cmd {
	return nil
}

func (m *ResultPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateStyles()

	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *ResultPaneModel) View() string {
	var paneStyle lipgloss.Style
	if m.isActive {
		paneStyle = m.activeStyles
	} else {
		paneStyle = m.styles
	}

	resultPane := paneStyle.Render()

	if m.err != nil {
		resultPane += lipgloss.NewStyle().
			Foreground(lipgloss.Color("red")).
			Render(m.err.Error())
	}

	return resultPane
}
