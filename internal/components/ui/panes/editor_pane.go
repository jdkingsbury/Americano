package panes

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdkingsbury/americano/internal/drivers"
	"github.com/jdkingsbury/americano/msgtypes"
)

/* Handles The SQL Editor Pane*/

// TODO: Work on displaying query results

type EditorPaneModel struct {
	styles       lipgloss.Style
	activeStyles lipgloss.Style
	width        int
	height       int
	textarea     textarea.Model
	err          error
	focused      bool
	isActive     bool
	db           drivers.Database
	resultPane   *ResultPaneModel
}

// Initialize Editor Pane
func NewEditorPane(width, height int, db drivers.Database, resultPane *ResultPaneModel) *EditorPaneModel {
	ti := textarea.New()
	ti.Placeholder = "Enter SQL Code Here..."
	ti.CharLimit = 1000
	ti.ShowLineNumbers = false

	pane := &EditorPaneModel{
		width:      width,
		height:     height,
		textarea:   ti,
		err:        nil,
		focused:    true,
		db:         db,
		resultPane: resultPane,
	}

	pane.updateStyles()

	return pane
}

// Code for changing from active to inactive window
func (m *EditorPaneModel) updateStyles() {
	m.styles = lipgloss.NewStyle().
		Width(m.width - 42).
		Height(m.height - 17).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(iris))

	m.activeStyles = lipgloss.NewStyle().
		Width(m.width - 42).
		Height(m.height - 17).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(rose))
}

// Code for functionality on start
func (m *EditorPaneModel) Init() tea.Cmd {
	return m.textarea.Focus()
}

// Code for updating the state
func (m *EditorPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resizeTextArea()

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			query := m.textarea.Value()
			return m, func() tea.Msg {
				return m.db.ExecuteQuery(query)
			}
		}

		switch msg.Type {
		// Keymap to switch stop editing
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
				m.focused = false
			} else {
				cmd = m.textarea.Focus()
				m.focused = true
				cmds = append(cmds, cmd)
			}

		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				m.focused = true
				cmds = append(cmds, cmd)
			}
		}
	case msgtypes.ErrMsg:
		m.err = msg
		return m, nil
	}

	// Resizes the text area view to fit main pane
	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// Helper function to resize the text area
func (m *EditorPaneModel) resizeTextArea() {
	m.textarea.SetWidth(m.width - 42)
	m.textarea.SetHeight(m.height - 17)
}

// Editor View
func (m EditorPaneModel) View() string {
	// Render text area inside the main pane

	var paneStyle lipgloss.Style
	if m.isActive {
		paneStyle = m.activeStyles
	} else {
		paneStyle = m.styles
	}

	mainPane := paneStyle.Render(m.textarea.View())
	return mainPane
}
