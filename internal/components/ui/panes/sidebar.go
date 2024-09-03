package panes

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdkingsbury/americano/msgtypes"
)

/* Handles the side bar pane */

// TODO: Look to see if we will need to create different instances of list
// for the db tree or if bubble tea has another way of creating the tree.
// list is used for using the Bubble Tea list.

type SideBarPaneModel struct {
	listConfig         *SideBarConfig       // Bubble Tea List Syle
	styles             lipgloss.Style       // Normal Pane Style
	activeStyles       lipgloss.Style       // Active Pane Style
	width              int
	height             int
	isActive           bool                 // Check if the pane is active
	isAddingConnection bool                 // Check if adding a connection
	list               list.Model           // For storing list items
	inputs             []textinput.Model
	focusedIndex       int
	err                error
	connections        []DatabaseConnection // List of Database Connections
	currentView        SideBarView          // Display SideBar Views
}

// Initialize Side Bar Pane
func NewSideBarPane(width, height int) *SideBarPaneModel {
	listConfig := NewSideBarConfig()
	inputs := make([]textinput.Model, 2)

	placeholders := []string{"Enter Name", "Enter Connection URL"}

	for i := range inputs {
		ti := textinput.New()
		ti.Placeholder = placeholders[i]
		ti.CharLimit = 256
		ti.Width = width/4 - 2
		if i == 0 {
			ti.Focus()
		}
		inputs[i] = ti
	}

	items := []list.Item{
		SideBarItem{Name: "󰆺 Add Connection", IsButton: true},
	}

	li := list.New(items, itemDelegate{}, width/4, listHeight)
	li.Title = "Database Connections"
	li.SetShowStatusBar(false)
	li.SetFilteringEnabled(false)
	li.SetShowHelp(false) // Disable help text
	li.Styles.Title = listConfig.TitleStyle
	li.Styles.PaginationStyle = listConfig.PaginationStyle

	pane := &SideBarPaneModel{
		listConfig:   listConfig,
		width:        width,
		height:       height,
		list:         li,
		inputs:       inputs,
		focusedIndex: 0,
		err:          nil,
	}

	pane.updateStyles() // Initialize styles

	return pane
}

// Code for changing from active to inactive window
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

// TODO: Create Checks to ensure both fields are filled and to check if the connection is valid

// Code for adding a DB Connection
func (m *SideBarPaneModel) addConnection(name, url string) {
	connection := DatabaseConnection{Name: name, URL: url}
	m.connections = append(m.connections, connection)

	// Append the new connection to the list
	newItem := SideBarItem{Name: " 󰇯 " + name}
	m.list.InsertItem(len(m.list.Items()), newItem)

	// Set the new item as selected
	m.list.Select(len(m.list.Items()) - 1)

	// Reset after adding the connection
	m.isAddingConnection = false
	m.inputs[0].SetValue("")
	m.inputs[1].SetValue("")
	m.focusedIndex = 0
}


// Code for functionality on start
func (m *SideBarPaneModel) Init() tea.Cmd {
	return nil
}

// Code for updating the state
func (m *SideBarPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
  // Fetch Window Size
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateStyles()

	case tea.KeyMsg:
		switch msg.String() {

		// Keymap to switch Side Bar views
		case "v":
			if m.currentView == ConnectionsView {
				m.currentView = DBTreeView
			} else {
				m.currentView = ConnectionsView
			}

			// Keymap to exit and add the db connection
		case "enter":
			if m.isAddingConnection && m.currentView == ConnectionsView {
				name := m.inputs[0].Value()
				url := m.inputs[1].Value()
				m.addConnection(name, url)
				m.inputs[0].SetValue("")
				m.inputs[1].SetValue("")
				m.focusedIndex = 0
				m.isAddingConnection = false
			} else if m.currentView == ConnectionsView {
				selectedItem := m.list.SelectedItem().(SideBarItem)
				if selectedItem.IsButton {
					m.isAddingConnection = true
				}
			}

			// Keymap to switch between text fields
		case "up", "down":
			if m.isAddingConnection && m.currentView == ConnectionsView {
				if m.isAddingConnection && m.focusedIndex > 0 {
					m.focusedIndex--
				} else if m.isAddingConnection && m.focusedIndex < len(m.inputs)-1 {
					m.focusedIndex++
				}
			}
		}

	case msgtypes.ErrMsg:
		m.err = msg
		return m, nil
	}

	// Checks to see if we are adding a connection and which text field we are in
	if m.isAddingConnection {
		for i := range m.inputs {
			if i == m.focusedIndex {
				m.inputs[i].Focus()
			} else {
				m.inputs[i].Blur()
			}
			m.inputs[i], cmd = m.inputs[i].Update(msg)
		}
		return m, cmd
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// Side Bar Views
func (m *SideBarPaneModel) View() string {
	var content string

	// Connection Views
	switch m.currentView {
	case ConnectionsView:
		if m.isAddingConnection {
			for _, input := range m.inputs {
				content += input.View() + "\n"
			}
		} else {
			content += titleStyle.Render(m.list.Title) + "\n"
			content += m.list.View()
		}

	case DBTreeView:
		content += "Database Tree"
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
