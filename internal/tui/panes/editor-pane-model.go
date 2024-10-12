package panes

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdkingsbury/americano/internal/drivers"
)

/* Handles The SQL Editor Pane*/

var sqlKeywords = map[string]bool{
	"SELECT": true, "FROM": true, "WHERE": true,
	"INSERT": true, "UPDATE": true, "DELETE": true,
	"CREATE": true, "TABLE": true, "JOIN": true,
	"ON": true,
}

func highlightSQL(text string) string {
	words := strings.Split(text, " ")

	for i, word := range words {
		if sqlKeywords[strings.ToUpper(word)] {
			words[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(pine)).Bold(true).Render(word)
		}
	}
	return strings.Join(words, " ")
}

type InsertQueryMsg struct {
	Query string
}

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
	keys         editorKeyMap
}

type editorKeyMap struct {
	ExecuteQuery key.Binding
}

func newEditorPaneKeymap() editorKeyMap {
	return editorKeyMap{
		ExecuteQuery: key.NewBinding(
			key.WithKeys("ctrl+e"),
			key.WithHelp("ctrl+e", "execute query"),
		),
	}
}

// Initialize Editor Pane
func NewEditorPane(width, height int, db drivers.Database) *EditorPaneModel {
	ti := textarea.New()
	ti.Placeholder = "Enter SQL Code Here..."
	ti.CharLimit = 1000
	ti.ShowLineNumbers = false
	ti.Prompt = " "

	pane := &EditorPaneModel{
		width:    width,
		height:   height,
		textarea: ti,
		err:      nil,
		focused:  false,
		db:       db,
		keys:     newEditorPaneKeymap(),
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

func (m *EditorPaneModel) Query() string {
	return m.textarea.Value()
}

func (m *EditorPaneModel) KeyMap() []key.Binding {
	return []key.Binding{m.keys.ExecuteQuery}
}

func (m *EditorPaneModel) Init() tea.Cmd {
	return nil
}

func (m *EditorPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resizeTextArea()

	case InsertQueryMsg:
		m.textarea.Reset()
		m.textarea.SetValue(msg.Query)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.ExecuteQuery):
			query := m.textarea.Value()
			return m, func() tea.Msg {
				return m.db.ExecuteQuery(query)
			}
		}

		switch msg.Type {
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
	}

	// Used for resizing the text area view to fit main pane
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
	var paneStyle lipgloss.Style
	if m.isActive {
		paneStyle = m.activeStyles
	} else {
		paneStyle = m.styles
	}

	// Get the entire rendered view for ensuring we have the desired dimensions
	textareaView := m.textarea.View()

	// Get the text input to check for highlighting
	rawText := m.textarea.Value()

	// Use highlight sql text function to highlight text
	highlightedText := highlightSQL(rawText)

	// Replace the raw text in textarea view with the highlighted text
	highlightedView := strings.Replace(textareaView, rawText, highlightedText, 1)

	return paneStyle.Render(highlightedView)
}
