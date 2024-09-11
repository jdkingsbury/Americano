package panes

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdkingsbury/americano/internal/drivers"
)

type pane int

// List Of Panes To Cycle Through
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

func NewLayoutModel() *LayoutModel {
	sideBarPane := NewSideBarPane(0, 0)
	resultPane := NewResultPaneModel(0, 0)
	editorPane := NewEditorPane(0, 0, nil, resultPane)
	footerPane := NewFooterPane(0)

	layout := &LayoutModel{
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

	// Set the initial active pane
	layout.setActivePane(true)

	return layout
}

// Updates pane sizes
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

func setupEditorPaneForDBConnection(dbURL string, width, height int, resultPane *ResultPaneModel) (*EditorPaneModel, error) {
	// Connect to database
	db, err := drivers.ConnectToDatabase(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize the editor pane with connected database and resultPane
	editorPane := NewEditorPane(width, height, db, resultPane)
	return editorPane, nil
}

// Code for functionality on start
func (m *LayoutModel) Init() tea.Cmd {
	return nil
}

// Code for updating the state
func (m *LayoutModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case SetupEditorPaneMsg:
		resultPane := m.panes[ResultPane].(*ResultPaneModel) // Retrieve the resultPane
		editorPane, err := setupEditorPaneForDBConnection(msg.dbURL, m.width, m.height, resultPane)
		if err != nil {
			fmt.Println("Error setting up editor pane:", err)
		} else {
			m.panes[EditorPane] = editorPane
			m.currentPane = EditorPane
		}

	case QueryResultsExecuted:
		m.currentPane = ResultPane

	// Fetch Window Size
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updatePaneSizes()

	case tea.KeyMsg:

		// Check if Adding Connection to disable layout commands temporarily
		if m.currentPane == SideBarPane {
			sideBarPane := m.panes[SideBarPane].(*SideBarPaneModel)
			if sideBarPane.showInputForm {
				break
			}
		}
		//     // Check if using the editor pane
		// } else if m.currentPane == EditorPane {
		// 	editorPane := m.panes[EditorPane].(*EditorPaneModel)
		// 	if editorPane.focused {
		// 		break
		// 	}
		// }

		switch msg.String() {
		// Keymap For Switching To Next Pane. Also Changes The Active Pane
		case "tab":
			// Deactivate current pane
			m.setActivePane(false)

			// Switch to the next pane
			m.currentPane = pane((int(m.currentPane) + 1) % len(m.panes))

			// Activate the new current pane
			m.setActivePane(true)

			// Keymap For Switching To Previous Pane. Also Changes The Active Pane
		case "shift+tab":
			// Deactivate current pane
			m.setActivePane(false)

			// Switch to the previous pane
			m.currentPane = pane((int(m.currentPane) - 1 + len(m.panes)) % len(m.panes))

			// Activate the new current pane
			m.setActivePane(true)

			// Keymap for Quitting the Application. Used for testing.
		case "Q":
			return m, tea.Quit
		}
	}

	// Retrieves the model for the current Pane. Ensures current pane is a valid index.
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

// Application Layout View
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
