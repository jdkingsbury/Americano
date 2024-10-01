package panes

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TODO: See if we can use the help bubble tea component to help with keymaps

/* Basic Footer View */

type FooterModel struct {
	style   lipgloss.Style
	keyMaps []key.Binding
	width   int
	height  int
}

func NewFooterPane(width int) *FooterModel {
	s := lipgloss.NewStyle().
		Width(width).
		Height(1).
		Foreground(lipgloss.Color(text)).
		Padding(0, 1)

	footer := &FooterModel{
		style:  s,
		width:  width,
		height: 1,
	}
	return footer
}

func (m *FooterModel) SetKeyBindings(keyMaps []key.Binding) {
	m.keyMaps = keyMaps
}

func (m *FooterModel) Init() tea.Cmd {
	return nil
}

func (m *FooterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = 1
	}
	return m, nil
}

func (m FooterModel) View() string {
	if len(m.keyMaps) == 0 {
		return m.style.Render("Q: Quit | ?: Help")
	}

	var renderKeys []string
	for _, k := range m.keyMaps {
		renderKeys = append(renderKeys, k.Help().Key+" "+k.Help().Desc)
	}

	return m.style.Render(" " + lipgloss.JoinHorizontal(lipgloss.Left, renderKeys...))
}
