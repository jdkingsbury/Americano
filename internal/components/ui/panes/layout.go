package panes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type pane int

const (
	SideBarPane pane = iota
	EditorPane
	ResultPane
)

type LayoutModel struct {
	currentPane pane
	panes       []tea.Model
	footer      *FooterModel
	width       int
	height      int
}

func (m *LayoutModel) Init() tea.Cmd {
	return nil
}

func (m *LayoutModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updatePaneSizes()
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			// Deactivate current pane
			m.setActivePane(false)

			// Switch to the next pane
			m.currentPane = pane((int(m.currentPane) + 1) % len(m.panes))

			// Activate the new current pane
			m.setActivePane(true)

		case "shift+tab":
			// Deactivate current pane
			m.setActivePane(false)

			// Switch to the previous pane
			m.currentPane = pane((int(m.currentPane) - 1 + len(m.panes)) % len(m.panes))

			// Activate the new current pane
			m.setActivePane(true)
		case "Q":
			return m, tea.Quit
		}
	}

	if int(m.currentPane) >= 0 && int(m.currentPane) < len(m.panes) {
		model := m.panes[m.currentPane]
		m.panes[m.currentPane], cmd = model.Update(msg)
	}

	return m, cmd
}

// Helper function to set the active status of the current pane
func (m *LayoutModel) setActivePane(isActive bool) {
	switch pane := m.panes[m.currentPane].(type) {
	case *SideBarPaneModel:
		pane.isActive = isActive
	case *EditorPaneModel:
		pane.isActive = isActive
	case *ResultPaneModel:
		pane.isActive = isActive
	}
}

func (m *LayoutModel) View() string {
	sideBarView := m.panes[SideBarPane].View()
	editorView := m.panes[EditorPane].View()
	resultView := m.panes[ResultPane].View()

	leftSide := lipgloss.JoinHorizontal(lipgloss.Left, sideBarView)
	rightSide := lipgloss.JoinHorizontal(lipgloss.Left, editorView)

	layout := lipgloss.JoinHorizontal(lipgloss.Left, leftSide, rightSide)
	layout = lipgloss.JoinVertical(lipgloss.Top, layout, resultView)

	footerView := m.footer.View()

	return lipgloss.JoinVertical(lipgloss.Top, layout, footerView)
}

func (m *LayoutModel) updatePaneSizes() {
	for _, pane := range m.panes {
		switch pane := pane.(type) {
		case *SideBarPaneModel:
			pane.width = m.width
			pane.height = m.height
			pane.updateStyles()
		case *EditorPaneModel:
			pane.width = m.width
			pane.height = m.height
			pane.resizeTextArea()
		case *ResultPaneModel:
			pane.width = m.width
			pane.height = m.height
			pane.updateStyles()
		}
	}

	m.footer.width = m.width
	m.footer.updateStyle()
}

func NewLayoutModel() *LayoutModel {
	sideBarPane := NewSideBarPane(0, 0)
	editorPane := NewEditorPane(0, 0)
	resultPane := NewResultPane(0, 0)
	footerPane := NewFooterPane(0)

	return &LayoutModel{
		currentPane: EditorPane,
		panes: []tea.Model{
			sideBarPane, // Index 0
			editorPane,  // Index 1
			resultPane,  // Index 2
		},
		footer: footerPane,
		width:  0,
		height: 0,
	}
}
