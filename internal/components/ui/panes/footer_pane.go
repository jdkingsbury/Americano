package panes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FooterModel struct {
	style  lipgloss.Style
	width  int
	height int
}

func (m *FooterModel) Init() tea.Cmd {
	return nil
}

func (m *FooterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = 1
		m.updateStyle()
	}
	return m, nil
}

func (m FooterModel) View() string {
	return m.style.Render("Q: Quit | ?: Help")
}

func NewFooterPane(width int) *FooterModel {
	footer := &FooterModel{
		width:  width,
		height: 1,
	}
	footer.updateStyle()
	return footer
}

func (m *FooterModel) updateStyle() {
  m.style = lipgloss.NewStyle().
    Width(m.width).
    Height(m.height).
    Foreground(lipgloss.Color(text)).
    Padding(0, 1) 
}
