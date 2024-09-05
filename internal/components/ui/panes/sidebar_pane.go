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
	styles         lipgloss.Style
	activeStyles   lipgloss.Style
	width          int
	height         int
	isActive       bool
	currentView    SideBarView
	dbConnModel    *DBConnModel
	connInputModel *DBConnInputModel
	showInputForm  bool
}

func NewSideBarPane(width, height int) *SideBarPaneModel {
	dbConnModel := NewDBConnModel(width)
	connInputModel := NewConnInputModel(width)

	pane := &SideBarPaneModel{
		width:          width,
		height:         height,
		dbConnModel:    dbConnModel,
		connInputModel: connInputModel,
		currentView:    ConnectionsView,
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
			if m.showInputForm {
				name := m.connInputModel.inputs[0].Value()
				url := m.connInputModel.inputs[1].Value()
				m.connInputModel.addConnection(name, url)
				m.showInputForm = false
        return m, tea.Quit
			} else {
				item, ok := m.dbConnModel.list.SelectedItem().(DBConnItems)
				if ok && item.isButton {
					m.showInputForm = true
				}
			}
		}

		// Input Form for adding a connection
		if m.showInputForm {
			updatedModel, modelCmd := m.connInputModel.Update(msg)
			m.connInputModel = updatedModel.(*DBConnInputModel)
			cmd = tea.Batch(cmd, modelCmd)
			// DB Connections view
		} else {
			updatedModel, modelCmd := m.dbConnModel.Update(msg)
			m.dbConnModel = updatedModel.(*DBConnModel)
			cmd = tea.Batch(cmd, modelCmd)
		}
		// TODO: Add DB Tree

	// switch m.currentView {
	// case ConnectionsView:
	// 	updatedModel, modelCmd := m.dbConnModel.Update(msg)
	// 	m.dbConnModel = updatedModel.(*DBConnModel)
	// 	cmd = tea.Batch(cmd, modelCmd)
	// }

	case tea.Msg:
		if msg == tea.Msg("add-connection-clicked") {
			m.showInputForm = true
		}
	}

	return m, cmd
}

func (m *SideBarPaneModel) View() string {
	var content string

	// Connection Views
	if m.showInputForm {
		content += m.connInputModel.View()
	} else {
		switch m.currentView {
		case ConnectionsView:
			// content += "Connections View"
			content += m.dbConnModel.View()
		case DBTreeView:
			content += "Database Tree View"
		}
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
