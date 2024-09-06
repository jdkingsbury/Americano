package panes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SideBarView int

// SideBar Views
const (
	ConnectionsView SideBarView = iota
	DBTreeView
)

type SideBarPaneModel struct {
	styles        lipgloss.Style
	activeStyles  lipgloss.Style
	width         int
	height        int
	isActive      bool
	currentView   SideBarView
	dbConnModel   *DBConnModel
	dbTreeModel   *DBTreeModel
	dbFormModel   *DBFormModel
	showInputForm bool
}

func NewSideBarPane(width, height int) *SideBarPaneModel {
	dbConnModel := NewDBConnModel(width)
	dbTreeModel := NewDBTreeModel()
	dbFormModel := NewDBFormModel()

	pane := &SideBarPaneModel{
		width:       width,
		height:      height,
		dbConnModel: dbConnModel,
		dbTreeModel: dbTreeModel,
		dbFormModel: dbFormModel,
		currentView: ConnectionsView,
	}

	pane.updateStyles()

	return pane
}

func (m *SideBarPaneModel) updateStyles() {
	m.styles = lipgloss.NewStyle().
		Width(m.width / 4).
		Height(m.height - 17).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(iris))

	m.activeStyles = lipgloss.NewStyle().
		Width(m.width / 4).
		Height(m.height - 17).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(rose))
}

func (m *SideBarPaneModel) Init() tea.Cmd {
	return nil
}

func (m *SideBarPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateStyles()

	case tea.KeyMsg:
		switch msg.String() {
		case "v":
			if m.currentView == ConnectionsView {
				m.currentView = DBTreeView
			} else {
				m.currentView = ConnectionsView
			}

		case "enter":
			if m.currentView == ConnectionsView && m.dbConnModel.FocusedOnButton() {
				m.showInputForm = true
			}
		}

	case CancelFormMsg:
		// Hide form after submission
		m.showInputForm = false
		// Reset Form
		m.dbFormModel.Reset()
		// Switch back to connections view
		m.currentView = ConnectionsView

	case SubmitFormMsg:
		// Add New Connection
		m.dbConnModel.AddConnection(msg.Name, msg.URL)
		// Hide form after submission
		m.showInputForm = false
		// Reset Form
		m.dbFormModel.Reset()
		// Switch back to connections view
		m.currentView = ConnectionsView
	}

	// Input Form for adding a connection
	if m.showInputForm {
		updatedForm, formCmd := m.dbFormModel.Update(msg)
		m.dbFormModel = updatedForm.(*DBFormModel)
		cmd = tea.Batch(cmd, formCmd)
	} else if m.currentView == ConnectionsView {
		updatedModel, modelCmd := m.dbConnModel.Update(msg)
		m.dbConnModel = updatedModel.(*DBConnModel)
		cmd = tea.Batch(cmd, modelCmd)
	} else if m.currentView == DBTreeView {
		updateModel, modelCmd := m.dbTreeModel.Update(msg)
		m.dbTreeModel = updateModel.(*DBTreeModel)
		cmd = tea.Batch(cmd, modelCmd)
	}

	return m, cmd
}

func (m *SideBarPaneModel) View() string {
	var content string

	// Connection Views
	if m.showInputForm {
		content = m.dbFormModel.View()
	} else if m.currentView == ConnectionsView {
		content = m.dbConnModel.View()
	} else if m.currentView == DBTreeView {
		content = m.dbTreeModel.View()
	}

	var paneStyle lipgloss.Style
	if m.isActive {
		paneStyle = m.activeStyles
	} else {
		paneStyle = m.styles
	}

	sideBarPane := paneStyle.Render(content)
	return sideBarPane
}
