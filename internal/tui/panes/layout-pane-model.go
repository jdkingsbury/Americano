package panes

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdkingsbury/americano/internal/drivers"
	"github.com/jdkingsbury/americano/msgtypes"
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
	keys        layoutKeyMap
}

type layoutKeyMap struct {
	NextPane key.Binding
	PrevPane key.Binding
	Help     key.Binding
	Quit     key.Binding
}

func newLayoutPaneKeyMapModel() layoutKeyMap {
	return layoutKeyMap{
		NextPane: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next pane"),
		),
		PrevPane: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "previous pane"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("Q"),
			key.WithHelp("Q", "quit americano"),
		),
	}
}

func NewLayoutModel() *LayoutModel {
	sideBarPane := NewSideBarPane(0, 0)
	resultPane := NewResultPaneModel(0, 0)
	editorPane := NewEditorPane(0, 0, nil)
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
		keys:   newLayoutPaneKeyMapModel(),
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
}

func setupEditorPaneForDBConnection(dbURL string, width, height int) (*EditorPaneModel, tea.Cmd) {
	db, notificationMsg := drivers.ConnectToDatabase(dbURL)
	if db == nil {
		return nil, func() tea.Msg {
			return notificationMsg
		}
	}

	// Initialize the editor pane with connected database
	editorPane := NewEditorPane(width, height, db)
	return editorPane, func() tea.Msg {
		return notificationMsg
	}
}

func setupDBTreeForDBConnection(dbURL string) (*DBTreeModel, tea.Cmd) {
	db, notificationMsg := drivers.ConnectToDatabase(dbURL)
	if db == nil {
		return nil, func() tea.Msg {
			return notificationMsg
		}
	}

	// Initialize the db tree with connected database
	dbTree := NewDBTreeModel(db)
	return dbTree, func() tea.Msg {
		return notificationMsg
	}
}

func (m *LayoutModel) Init() tea.Cmd {
	return m.setActivePane(true)
}

func (m *LayoutModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case InsertQueryMsg:
		// Pass the query to editor pane
		editorPane := m.panes[EditorPane].(*EditorPaneModel)
		m.panes[EditorPane], cmd = editorPane.Update(msg)
		return m, cmd

	case SetupDBTreeMsg:
		dbTree, setupCmd := setupDBTreeForDBConnection(msg.dbURL)

		if dbTree != nil {
			sideBarPane := m.panes[SideBarPane].(*SideBarPaneModel)
			sideBarPane.dbTreeModel = dbTree
			sideBarPane.currentView = DBTreeView
		}

		return m, setupCmd

	case SetupEditorPaneMsg:
		editorPane, setupCmd := setupEditorPaneForDBConnection(msg.dbURL, m.width, m.height)
		if editorPane != nil {
			m.panes[EditorPane] = editorPane
		}

		return m, setupCmd

	case SetKeyMapMsg:
		m.footer.SetKeyBindings(msg.FullHelpKeys, msg.ShortHelpKeys)
		return m, nil

	case msgtypes.NotificationMsg:
		resultPane := m.panes[ResultPane].(*ResultPaneModel)
		resultPane.Update(msg)

	case msgtypes.ErrMsg:
		resultPane := m.panes[ResultPane].(*ResultPaneModel)
		resultPane.Update(msg)

	case drivers.QueryResultMsg:
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
			// Check if using the editor pane
		} else if m.currentPane == EditorPane {
			editorPane := m.panes[EditorPane].(*EditorPaneModel)
			if editorPane.focused {
				break
			}
		}

		switch {
		case key.Matches(msg, m.keys.NextPane):
			m.setActivePane(false)
			m.currentPane = pane((int(m.currentPane) + 1) % len(m.panes))
			return m, tea.Batch(cmd, m.setActivePane(true))

		case key.Matches(msg, m.keys.PrevPane):
			m.setActivePane(false)
			m.currentPane = pane((int(m.currentPane) - 1 + len(m.panes)) % len(m.panes))
			return m, tea.Batch(cmd, m.setActivePane(true))

		case key.Matches(msg, m.keys.Help):
			if m.currentPane != EditorPane {
				m.footer.showFullHelp = !m.footer.showFullHelp
			}

		case key.Matches(msg, m.keys.Quit):
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
func (m *LayoutModel) setActivePane(isActive bool) tea.Cmd {
	layoutFullHelp := [][]key.Binding{
		{m.keys.NextPane, m.keys.PrevPane, m.keys.Quit}, // Layout keybindings
	}

	switch pane := m.panes[m.currentPane].(type) {
	case *SideBarPaneModel:
		pane.isActive = isActive
		if isActive {
			return func() tea.Msg {
				return SetKeyMapMsg{
					FullHelpKeys:  append(layoutFullHelp, pane.KeyMap()),
					ShortHelpKeys: append(pane.KeyMap()),
				}
			}
		}
	case *EditorPaneModel:
		pane.isActive = isActive
		if isActive {
			return func() tea.Msg {
				return SetKeyMapMsg{
					FullHelpKeys:  append(layoutFullHelp, pane.KeyMap()),
					ShortHelpKeys: append(pane.KeyMap()),
				}
			}
		}
	case *ResultPaneModel:
		pane.isActive = isActive
		if isActive {
			return func() tea.Msg {
				return SetKeyMapMsg{
					FullHelpKeys:  append(layoutFullHelp, pane.KeyMap()),
					ShortHelpKeys: append(pane.KeyMap()),
				}
			}
		}
	}

	return nil
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
