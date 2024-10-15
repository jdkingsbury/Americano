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

var sqlKeywords = map[string]bool{
	"SELECT": true, "FROM": true, "WHERE": true,
	"INSERT": true, "UPDATE": true, "DELETE": true,
	"CREATE": true, "TABLE": true, "JOIN": true,
	"ON": true, "INNER": true, "LEFT": true,
	"RIGHT": true, "GROUP": true, "ORDER": true,
	"BY": true, "DESC": true, "ASC": true,
}

func highlightSQL(text string) string {
	var builder strings.Builder
	words := strings.Fields(text)

	for i, word := range words {
		upperWord := strings.ToUpper(word)
		if i > 0 {
			builder.WriteString(" ")
		}

		switch {
		case sqlKeywords[upperWord]:
			builder.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(pine)).Bold(true).Render(word))
		case strings.HasPrefix(word, "'") && strings.HasSuffix(word, "'"):
			builder.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(gold)).Bold(true).Render(word))
		case strings.HasPrefix(word, "--"):
			builder.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(muted)).Bold(true).Render(word))
		default:
			builder.WriteString(word)
		}
	}
	return builder.String()
}

func isPrintable(keyMsg tea.KeyMsg) bool {
	s := keyMsg.String()
	return len(s) == 1 && s[0] >= 32 && s[0] <= 126
}

func isWhiteSpace(ch byte) bool {
  return ch == ' ' || ch == '\t' || ch == '\n'
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
	isActive     bool
	mode         int
	db           drivers.Database
	keys         editorKeyMap
}

type editorKeyMap struct {
	ExecuteQuery key.Binding
	NormalMode   key.Binding
	InsertMode   key.Binding
	Up           key.Binding
	Down         key.Binding
	Left         key.Binding
	Right        key.Binding
	Enter        key.Binding
	Backspace    key.Binding
	ArrowUp      key.Binding
	ArrowDown    key.Binding
	ArrowLeft    key.Binding
	ArrowRight   key.Binding
}

func newEditorPaneKeymap() editorKeyMap {
	return editorKeyMap{
		ExecuteQuery: key.NewBinding(
			key.WithKeys("ctrl+e"),
			key.WithHelp("ctrl+e", "execute query"),
		),
		NormalMode: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "enter normal mode"),
		),
		InsertMode: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "enter insert mode"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
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
		ArrowUp: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "move up"),
		),
		ArrowDown: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "move down"),
		),
		ArrowLeft: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "move left"),
		),
		ArrowRight: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "move right"),
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
		m.keys.NormalMode,
		m.keys.InsertMode,
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
		db:        db,
		keys:      newEditorPaneKeymap(),
		mode:      NormalMode,
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

func (m *EditorPaneModel) moveCursorForwardWord() {
  line := m.buffer[m.cursorRow]

  for m.cursorCol < len(line) && !isWhiteSpace(line[m.cursorCol]) {
    m.cursorCol++
  }

  for m.cursorCol < len(line) && isWhiteSpace(line[m.cursorCol]) {
    m.cursorCol++
  }
}

func (m *EditorPaneModel) moveCursorBackwardWord() {
  line := m.buffer[m.cursorRow]

  if m.cursorCol > 0 {
    m.cursorCol--
  }

  for m.cursorCol > 0 && isWhiteSpace(line[m.cursorCol]) {
    m.cursorCol--
  }

  for m.cursorCol > 0 && !isWhiteSpace(line[m.cursorCol-1]) {
    m.cursorCol--
  }
}


func (m *EditorPaneModel) moveCursorVertically(direction int) {
	m.cursorRow += direction
	if m.cursorRow < 0 {
		m.cursorRow = 0
	} else if m.cursorRow >= len(m.buffer) {
		m.cursorRow = len(m.buffer) - 1
	}
	m.ensureCursorInBounds()
}

func (m *EditorPaneModel) ensureCursorInBounds() {
	if m.cursorRow >= len(m.buffer) {
		m.cursorRow = len(m.buffer) - 1
	}
	if m.cursorCol > len(m.buffer[m.cursorRow]) {
		m.cursorCol = len(m.buffer[m.cursorRow])
	}
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
		switch m.mode {

		case NormalMode:
			switch {
			case key.Matches(msg, m.keys.InsertMode):
				m.mode = InsertMode
				return m, nil
      // Move forward a word
      case msg.String() == "w":
        m.moveCursorForwardWord()
      // Move backward a word
      case msg.String() == "b":
        m.moveCursorBackwardWord()
			case key.Matches(msg, m.keys.Up):
				m.moveCursorVertically(-1)
			case key.Matches(msg, m.keys.Down):
				m.moveCursorVertically(1)
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
			}

		case InsertMode:
			switch {
			case key.Matches(msg, m.keys.NormalMode):
				m.mode = NormalMode
				return m, nil
			case key.Matches(msg, m.keys.Enter):
				currentLine := m.buffer[m.cursorRow]
				m.buffer[m.cursorRow] = currentLine[:m.cursorCol]
				m.buffer = append(m.buffer[:m.cursorRow+1], append([]string{currentLine[m.cursorCol:]}, m.buffer[m.cursorRow+1:]...)...)
				m.cursorRow++
				m.cursorCol = 0
			case key.Matches(msg, m.keys.Backspace):
				if m.cursorCol > 0 {
					// Deleting character in the middle of the line
					m.buffer[m.cursorRow] = m.buffer[m.cursorRow][:m.cursorCol-1] + m.buffer[m.cursorRow][m.cursorCol:]
					m.cursorCol--
				} else if m.cursorRow > 0 {
					// Save the current line before removing it
					currentLine := m.buffer[m.cursorRow]
					prevLine := m.buffer[m.cursorRow-1]

					// Remove the current line from the buffer
					m.buffer = append(m.buffer[:m.cursorRow], m.buffer[m.cursorRow+1:]...)

					// Move the cursor to the end of the previous line
					m.cursorRow--
					m.cursorCol = len(prevLine)

					// Concatenate the previous line with the current line
					m.buffer[m.cursorRow] = prevLine + currentLine
				}
			case key.Matches(msg, m.keys.ArrowUp):
				m.moveCursorVertically(-1)
			case key.Matches(msg, m.keys.ArrowDown):
				m.moveCursorVertically(1)
			case key.Matches(msg, m.keys.ArrowLeft):
				if m.cursorCol > 0 {
					m.cursorCol--
				} else if m.cursorRow > 0 {
					m.cursorRow--
					m.cursorCol = len(m.buffer[m.cursorRow])
				}
			case key.Matches(msg, m.keys.ArrowRight):
				if m.cursorCol < len(m.buffer[m.cursorRow]) {
					m.cursorCol++
				} else if m.cursorRow < len(m.buffer)-1 {
					m.cursorRow++
					m.cursorCol = 0
				}
			default:
				if isPrintable(msg) {
					m.buffer[m.cursorRow] = m.buffer[m.cursorRow][:m.cursorCol] + msg.String() + m.buffer[m.cursorRow][m.cursorCol:]
					m.cursorCol++
				}
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
			if m.mode == InsertMode {
				cursor = lipgloss.NewStyle().Foreground(lipgloss.Color(rose)).Render("|")
			} else {
				cursor = lipgloss.NewStyle().Background(lipgloss.Color(rose)).Foreground(lipgloss.Color(rose)).Render(" ")
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
