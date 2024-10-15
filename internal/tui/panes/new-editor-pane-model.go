package panes

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdkingsbury/americano/internal/drivers"
)

var sqlKeywords = map[string]bool{
	"SELECT": true, "FROM": true, "WHERE": true,
	"INSERT": true, "UPDATE": true, "DELETE": true,
	"CREATE": true, "TABLE": true, "JOIN": true,
	"ON": true, "INNER": true, "LEFT": true,
	"RIGHT": true, "GROUP": true, "ORDER": true,
	"BY": true, "DESC": true, "ASC": true,
}

func highlightSQL(text string) string {
	words := strings.Split(text, " ")

	for i, word := range words {
		upperWord := strings.ToUpper(word)
		if sqlKeywords[upperWord] {
			words[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(pine)).Bold(true).Render(word)
		} else if strings.HasPrefix(word, "'") && strings.HasSuffix(word, "'") {
			words[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(gold)).Bold(true).Render(word)
		} else if strings.HasPrefix(word, "--") {
			words[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(muted)).Bold(true).Render(word)
		}
	}
	return strings.Join(words, " ")
}

func isPrintable(keyMsg tea.KeyMsg) bool {
	s := keyMsg.String()
	return len(s) == 1 && s[0] >= 32 && s[0] <= 126
}

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
	focused      bool
	isActive     bool
	db           drivers.Database
	keys         editorKeyMap
}

type editorKeyMap struct {
	ExecuteQuery key.Binding
	Focus        key.Binding
	Up           key.Binding
	Down         key.Binding
	Left         key.Binding
	Right        key.Binding
	Enter        key.Binding
	Backspace    key.Binding
}

func newEditorPaneKeymap() editorKeyMap {
	return editorKeyMap{
		ExecuteQuery: key.NewBinding(
			key.WithKeys("ctrl+e"),
			key.WithHelp("ctrl+e", "execute query"),
		),
		Focus: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "toggle focus"),
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
		m.keys.ExecuteQuery,
		m.keys.Up,
		m.keys.Down,
		m.keys.Left,
		m.keys.Right,
		m.keys.Enter,
		m.keys.Backspace,
		m.keys.Focus,
	}
}

func NewEditorPane(width, height int, db drivers.Database) *EditorPaneModel {
	pane := &EditorPaneModel{
		width:     width,
		height:    height,
		buffer:    []string{""},
		cursorRow: 0,
		cursorCol: 0,
		err:       nil,
		focused:   true,
		db:        db,
		keys:      newEditorPaneKeymap(),
	}

	pane.updateStyles()

	return pane
}

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

func (m *EditorPaneModel) resizeTextArea() {
	m.width = m.width - 42
	m.height = m.height - 17
}

func (m *EditorPaneModel) Init() tea.Cmd {
	return nil
}

func (m *EditorPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateStyles()

	case InsertQueryMsg:
		m.buffer = strings.Split(msg.Query, "\n") // Convert query string to buffer
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Focus) {
			m.focused = !m.focused
			return m, tea.Batch(cmds...)
		}

		if !m.focused {
			return m, nil
		}

		switch {
		case key.Matches(msg, m.keys.ExecuteQuery):
			query := strings.Join(m.buffer, "\n")
			return m, func() tea.Msg {
				return m.db.ExecuteQuery(query)
			}
		case key.Matches(msg, m.keys.Up):
			if m.cursorRow > 0 {
				m.cursorRow--
				m.cursorCol = min(m.cursorCol, len(m.buffer[m.cursorRow]))
			}
		case key.Matches(msg, m.keys.Down):
			if m.cursorRow < len(m.buffer)-1 {
				m.cursorRow++
				m.cursorCol = min(m.cursorCol, len(m.buffer[m.cursorRow]))
			}
		case key.Matches(msg, m.keys.Left):
			if m.cursorCol > 0 {
				m.cursorCol--
			} else if m.cursorRow > 0 {
				m.cursorRow--
				m.cursorCol = len(m.buffer[m.cursorRow])
			}
		case key.Matches(msg, m.keys.Right):
			if m.cursorCol < len(m.buffer[m.cursorRow]) {
				m.cursorCol++
			} else if m.cursorRow < len(m.buffer)-1 {
				m.cursorRow++
				m.cursorCol = 0
			}
		case key.Matches(msg, m.keys.Enter):
			currentLine := m.buffer[m.cursorRow]
			m.buffer[m.cursorRow] = currentLine[:m.cursorCol]
			m.buffer = append(m.buffer[:m.cursorRow+1], append([]string{currentLine[m.cursorCol:]}, m.buffer[m.cursorRow+1:]...)...)
			m.cursorRow++
			m.cursorCol = 0
		case key.Matches(msg, m.keys.Backspace):
			if m.cursorCol > 0 {
				m.buffer[m.cursorRow] = m.buffer[m.cursorRow][:m.cursorCol-1] + m.buffer[m.cursorRow][m.cursorCol:]
				m.cursorCol--
			} else if m.cursorRow > 0 {
				prevLine := m.buffer[m.cursorRow-1]
				m.buffer = append(m.buffer[:m.cursorRow], m.buffer[m.cursorRow+1:]...)
				m.cursorRow--
				m.cursorCol = len(prevLine)
				m.buffer[m.cursorRow] = prevLine + m.buffer[m.cursorRow]
			}
		default:
			if isPrintable(msg) {
				m.buffer[m.cursorRow] = m.buffer[m.cursorRow][:m.cursorCol] + msg.String() + m.buffer[m.cursorRow][m.cursorCol:]
				m.cursorCol++
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *EditorPaneModel) View() string {
	var paneStyle lipgloss.Style
	if m.isActive {
		paneStyle = m.activeStyles
	} else {
		paneStyle = m.styles
	}

	var builder strings.Builder
	for i, line := range m.buffer {
		if i == m.cursorRow {
			cursor := ""
			if m.focused {
				cursor = lipgloss.NewStyle().Foreground(lipgloss.Color(rose)).Render("|")
			} else {
				cursor = lipgloss.NewStyle().Background(lipgloss.Color(text)).Render(" ")
			}

			highlightedLine := highlightSQL(line[:m.cursorCol]) + cursor + highlightSQL(line[m.cursorCol:])
			builder.WriteString(highlightedLine)
		} else {
			builder.WriteString(highlightSQL(line))
		}
		builder.WriteString("\n")
	}

	return paneStyle.Render(builder.String())
}
