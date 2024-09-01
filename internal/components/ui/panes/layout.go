package panes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type pane int

const (
	SideBarPane pane = iota
	EditorPane
)

type LayoutModel struct {
	currentPane pane
	panes       []tea.Model
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
			m.currentPane = pane((int(m.currentPane) + 1) % len(m.panes))
		case "shift+tab":
			m.currentPane = pane((int(m.currentPane) - 1 + len(m.panes)) % len(m.panes))
		}
	}

	if int(m.currentPane) >= 0 && int(m.currentPane) < len(m.panes) {
		model := m.panes[m.currentPane]
		m.panes[m.currentPane], cmd = model.Update(msg)
	}

	return m, cmd
}

func (m *LayoutModel) View() string {
	sideBarView := m.panes[SideBarPane].View()
	editorView := m.panes[EditorPane].View()

	leftSide := lipgloss.JoinHorizontal(lipgloss.Left, sideBarView)
	rightSide := lipgloss.JoinHorizontal(lipgloss.Left, editorView)

	layout := lipgloss.JoinHorizontal(lipgloss.Left, leftSide, rightSide)

	return layout
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
		}
	}
}

func NewLayoutModel() *LayoutModel {
  sideBarPane := NewSideBarPane(0, 0)
  editorPane := NewEditorPane(0, 0)

	return &LayoutModel{
		currentPane: EditorPane,
		panes: []tea.Model{
			sideBarPane, // Index 0
			editorPane,  // Index 1
		},
		width:  0,
		height: 0,
	}
}
