package panes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FloatingTextPaneModel struct {
	styles    lipgloss.Style
	width     int
	height    int
	isVisible bool
	content   string
}

func (m *FloatingTextPaneModel) updateStyles() {
	m.styles = lipgloss.NewStyle().
		Width(m.width / 3).
		Height(m.height / 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(love))
}

func NewFloatingTextPane(width, height int) *FloatingTextPaneModel {
	pane := &FloatingTextPaneModel{
		width:     width,
		height:    height,
		isVisible: false,
		content:   keyboard + "Enter Connection URL",
	}

	pane.updateStyles()

	return pane
}

func (m *FloatingTextPaneModel) Init() tea.Cmd {
	return nil
}

func (m *FloatingTextPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateStyles()

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.isVisible = false
		}
	}

	return m, nil
}

func (m *FloatingTextPaneModel) setContentString(content string) {
	m.content = content
}

func (m *FloatingTextPaneModel) View() string {
	if !m.isVisible {
		return ""
	}

	return m.styles.Render(m.content)
}
