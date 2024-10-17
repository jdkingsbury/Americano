package panes

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdkingsbury/americano/internal/drivers"
)

const (
	NormalMode = iota
	InsertMode
)

type InsertQueryMsg struct {
	Query string
}

type EditorPaneModel struct {
	styles       lipgloss.Style
	activeStyles lipgloss.Style
	width        int
	height       int
	buffer       []string
	cursorRow    int
	cursorCol    int
	err          error
	isActive     bool
	mode         int
	db           drivers.Database
	keys         editorKeyMap
}

type editorKeyMap struct {
	ExecuteQuery key.Binding
	Up           key.Binding
	Down         key.Binding
	Left         key.Binding
	Right        key.Binding
	Enter        key.Binding
	Backspace    key.Binding
}

func newEditorKeymap() editorKeyMap {
	return editorKeyMap{
		ExecuteQuery: key.NewBinding(
			key.WithKeys("ctrl+e"),
			key.WithHelp("ctrl+e", "execute query"),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "move right"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "new line"),
		),
		Backspace: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("backspace", "delete character"),
		),
	}
}

func (m *EditorPaneModel) KeyMap() []key.Binding {
	return []key.Binding{
		m.keys.Up,
		m.keys.Down,
		m.keys.Left,
		m.keys.Right,
		m.keys.ExecuteQuery,
	}
}

func NewEditorPane(width, height int, db drivers.Database) *EditorPaneModel {
	pane := &EditorPaneModel{
		width:     width,
		height:    height,
		buffer:    []string{""},
		cursorRow: 0,
		cursorCol: 0,
		mode:      NormalMode,
		err:       nil,
		db:        db,
		keys:      newEditorKeymap(),
	}

	pane.updateStyles()

	return pane
}

func (m *EditorPaneModel) updateStyles() {
	m.styles = lipgloss.NewStyle().
		Width(m.width - 42).
		Height(m.height - 17).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(iris)).
		Faint(true)

	m.activeStyles = lipgloss.NewStyle().
		Width(m.width - 42).
		Height(m.height - 17).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(rose))
}

func (m *EditorPaneModel) Init() tea.Cmd {
	return nil
}

func (m *EditorPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateStyles()

	case tea.KeyMsg:
		switch {
		// Execute Query
		case key.Matches(msg, m.keys.ExecuteQuery):
			// Join all lines in the buffer to get the full sql query code
			query := strings.Join(m.buffer, "\n")
			return m, func() tea.Msg {
				m.isActive = false
				return m.db.ExecuteQuery(query)
			}

			// Switch to Normal Mode
		case msg.String() == "i" && m.mode == NormalMode:
			m.mode = InsertMode

			// Switch to Insert Mode
		case msg.String() == "esc" && m.mode == InsertMode:
			m.mode = NormalMode

			// Normal Mode Commands
		case m.mode == NormalMode:
			switch {

			// Up
			case key.Matches(msg, m.keys.Up) || msg.String() == "k":
				if m.cursorRow > 0 {
					m.cursorRow--
					m.cursorCol = min(m.cursorCol, len(m.buffer[m.cursorRow]))
				}

				// Down
			case key.Matches(msg, m.keys.Down) || msg.String() == "j":
				if m.cursorRow < len(m.buffer)-1 {
					m.cursorRow++
					m.cursorCol = min(m.cursorCol, len(m.buffer[m.cursorRow]))
				}

				// Left
			case key.Matches(msg, m.keys.Left) || msg.String() == "h":
				if m.cursorCol > 0 {
					m.cursorCol--
				} else if m.cursorRow > 0 {
					m.cursorRow--
					m.cursorCol = len(m.buffer[m.cursorRow])
				}

				// Right
			case key.Matches(msg, m.keys.Right) || msg.String() == "l":
				if m.cursorCol < len(m.buffer[m.cursorRow]) {
					m.cursorCol++
				} else if m.cursorRow < len(m.buffer)-1 {
					m.cursorRow++
					m.cursorCol = 0
				}
			}

			// Insert Mode Commands
		case m.mode == InsertMode:
			switch {

			// Enter
			case key.Matches(msg, m.keys.Enter):
				// Split the current line at the cursor position
				newLine := m.buffer[m.cursorRow][m.cursorCol:]
				m.buffer[m.cursorRow] = m.buffer[m.cursorRow][:m.cursorCol]
				m.buffer = append(m.buffer[:m.cursorRow+1], append([]string{newLine}, m.buffer[m.cursorRow+1:]...)...)
				m.cursorRow++
				m.cursorCol = 0

				// Backspace
			case key.Matches(msg, m.keys.Backspace):
				if m.cursorCol > 0 {
					// Delete character before the cursor
					m.buffer[m.cursorRow] = m.buffer[m.cursorRow][:m.cursorCol-1] + m.buffer[m.cursorRow][m.cursorCol:]
					m.cursorCol--
				} else if m.cursorRow > 0 {
					// Merge the previous line
					prevLineLen := len(m.buffer[m.cursorRow-1])
					m.buffer[m.cursorRow-1] += m.buffer[m.cursorRow]
					m.buffer = append(m.buffer[:m.cursorRow], m.buffer[m.cursorRow+1:]...)
					m.cursorRow--
					m.cursorCol = prevLineLen
				}

				// Up
			case key.Matches(msg, m.keys.Up):
				if m.cursorRow > 0 {
					m.cursorRow--
					m.cursorCol = min(m.cursorCol, len(m.buffer[m.cursorRow]))
				}

				// Down
			case key.Matches(msg, m.keys.Down):
				if m.cursorRow < len(m.buffer)-1 {
					m.cursorRow++
					m.cursorCol = min(m.cursorCol, len(m.buffer[m.cursorRow]))
				}

				// Left
			case key.Matches(msg, m.keys.Left):
				if m.cursorCol > 0 {
					m.cursorCol--
				} else if m.cursorRow > 0 {
					m.cursorRow--
					m.cursorCol = len(m.buffer[m.cursorRow])
				}

				// Right
			case key.Matches(msg, m.keys.Right):
				if m.cursorCol < len(m.buffer[m.cursorRow]) {
					m.cursorCol++
				} else if m.cursorRow < len(m.buffer)-1 {
					m.cursorRow++
					m.cursorCol = 0
				}

				// Typing Characters into the Editor Pane
			default:
				if msg.Type == tea.KeyRunes {
					runes := msg.Runes
					// Insert character at cursor position
					m.buffer[m.cursorRow] = m.buffer[m.cursorRow][:m.cursorCol] + string(runes) + m.buffer[m.cursorRow][m.cursorCol:]
					m.cursorCol += len(runes)
				} else if msg.String() == " " {
					// Insert space character
					m.buffer[m.cursorRow] = m.buffer[m.cursorRow][:m.cursorCol] + " " + m.buffer[m.cursorRow][m.cursorCol:]
					m.cursorCol++
				}
			}
		}

		// Insert Query into Editor Pane
	case InsertQueryMsg:
		m.buffer = strings.Split(msg.Query, "\n")
		return m, nil
	}

	return m, nil
}

func (m *EditorPaneModel) View() string {
	var paneStyle lipgloss.Style
	if m.isActive {
		paneStyle = m.activeStyles
	} else {
		paneStyle = m.styles
	}

	// Render buffer lines and add the cursor at the correct position
	var output strings.Builder
	for i, line := range m.buffer {
		if i == m.cursorRow {
			// Render if the pane is active
			if m.isActive {

				// Render for Normal Mode
				if m.mode == NormalMode {
					cursor := "█" // Normal mode cursor
					// Insert cursor into the line at the correct column
					if m.cursorCol < len(line) {
						output.WriteString(line[:m.cursorCol] + lipgloss.NewStyle().Foreground(lipgloss.Color(rose)).Render(cursor) + line[m.cursorCol+1:])
					} else {
						output.WriteString(line + lipgloss.NewStyle().Foreground(lipgloss.Color(rose)).Render(cursor))
					}
				}

				// Render for Insert Mode
				if m.mode == InsertMode {
					cursor := "|" // Insert mode cursor
					// Insert cursor into the line at the correct column
					if m.cursorCol < len(line) {
						output.WriteString(line[:m.cursorCol] + lipgloss.NewStyle().Foreground(lipgloss.Color(rose)).Render(cursor) + line[m.cursorCol:])
					} else {
						output.WriteString(line + lipgloss.NewStyle().Foreground(lipgloss.Color(rose)).Render(cursor))
					}
				}
				// Render if the pane is inactive
			} else {
				output.WriteString(line)
			}
		} else {
			output.WriteString(line)
		}

		// Add new line unless it's the last line
		if i < len(m.buffer)-1 {
			output.WriteString("\n")
		}
	}

	return paneStyle.Render(output.String())
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
