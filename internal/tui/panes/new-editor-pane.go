package panes

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdkingsbury/americano/internal/drivers"
)

// TODO:
// Add the functionality to ensure code works on multiline
// Work on adding cursor blinking when in inset mode
// Work on move forward and backward a word to ensure that we always end up on the first character of a word

type TokenType int

const (
	TokenKeyword TokenType = iota
	TokenIdentifier
	TokenNumber
	TokenSymbol
	TokenString
	TokenComment
)

type token struct {
	Type  TokenType
	Value string
}

var sqlKeywords = map[string]struct{}{
	"SELECT": {}, "FROM": {}, "WHERE": {}, "INSERT": {}, "UPDATE": {},
	"DELETE": {}, "CREATE": {}, "DROP": {}, "ALTER": {}, "JOIN": {},
	"LEFT": {}, "RIGHT": {}, "INNER": {}, "OUTER": {}, "GROUP": {},
	"ORDER": {}, "BY": {}, "LIMIT": {}, "DISTINCT": {}, "AND": {},
	"OR": {}, "NOT": {}, "IN": {}, "LIKE": {}, "AS": {},
}

var tokenStyles = map[TokenType]lipgloss.Style{
	TokenKeyword:    lipgloss.NewStyle().Foreground(lipgloss.Color(pine)),
	TokenIdentifier: lipgloss.NewStyle().Foreground(lipgloss.Color(text)),
	TokenNumber:     lipgloss.NewStyle().Foreground(lipgloss.Color(gold)),
	TokenSymbol:     lipgloss.NewStyle().Foreground(lipgloss.Color(subtle)),
	TokenString:     lipgloss.NewStyle().Foreground(lipgloss.Color(gold)),
	TokenComment:    lipgloss.NewStyle().Foreground(lipgloss.Color(subtle)).Italic(true),
}

func isKeyword(word string) bool {
	_, exists := sqlKeywords[strings.ToUpper(word)]
	return exists
}

func isSymbol(word string) bool {
	symbolSet := "{}[](),.;+-/*=&<>"
	for _, ch := range word {
		if strings.ContainsRune(symbolSet, ch) {
			return true
		}
	}
	return false
}

func tokenize(line string) []token {
	var tokens []token
	var currentToken strings.Builder
	var inString bool
	tokenType := TokenIdentifier // default token type

	for _, char := range line {
		switch {
		case char == ' ' && !inString:
			// Complete the current token before processing the space
			if currentToken.Len() > 0 {
				word := currentToken.String()
				if isKeyword(word) {
					tokenType = TokenKeyword
				} else if isNumber(word) {
					tokenType = TokenNumber
				} else {
					tokenType = TokenIdentifier
				}
				tokens = append(tokens, token{Type: tokenType, Value: word})
				currentToken.Reset()
			}
			// Add space as a symbol token
			tokens = append(tokens, token{Type: TokenSymbol, Value: " "})

		case isSymbol(string(char)) && !inString:
			// Close the current token before processing a symbol
			if currentToken.Len() > 0 {
				word := currentToken.String()
				if isKeyword(word) {
					tokenType = TokenKeyword
				} else if isNumber(word) {
					tokenType = TokenNumber
				} else {
					tokenType = TokenIdentifier
				}
				tokens = append(tokens, token{Type: tokenType, Value: word})
				currentToken.Reset()
			}
			tokens = append(tokens, token{Type: TokenSymbol, Value: string(char)}) // Symbol as its own token

		case char == '"' || char == '\'':
			// Handle string literals
			if inString {
				// Close token string
				currentToken.WriteRune(char)
				tokens = append(tokens, token{Type: TokenString, Value: currentToken.String()})
				currentToken.Reset()
				inString = false
			} else {
				// Start a new token string
				if currentToken.Len() > 0 {
					// Flush any non-string tokens before starting a new string
					word := currentToken.String()
					if isKeyword(word) {
						tokenType = TokenKeyword
					} else if isNumber(word) {
						tokenType = TokenNumber
					} else {
						tokenType = TokenIdentifier
					}
					tokens = append(tokens, token{Type: tokenType, Value: word})
					currentToken.Reset()
				}
				inString = true
				currentToken.WriteRune(char)
			}

		default:
			currentToken.WriteRune(char) // Continue building the word
		}
	}

	// Add the last token if there's any remaining text
	if currentToken.Len() > 0 {
		word := currentToken.String()
		if isKeyword(word) {
			tokenType = TokenKeyword
		} else if isNumber(word) {
			tokenType = TokenNumber
		} else {
			tokenType = TokenIdentifier
		}
		tokens = append(tokens, token{Type: tokenType, Value: word})
	}

	return tokens
}

func isNumber(word string) bool {
	for _, ch := range word {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

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

// Helper function for determining the min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func isWordChar(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsNumber(ch) || unicode.IsPunct(ch) || ch == '_' || ch == '-' || ch == '*' || ch == '+' || ch == '=' || ch == '>' || ch == '<'
}

func isDelimeter(ch rune) bool {
	return unicode.IsSpace(ch) || unicode.IsPunct(ch) && ch != '_' && ch != '-' && ch != '*' && ch != '+'
}

// Function for moving forward by a word
func (m *EditorPaneModel) moveCursorForwardByWord(line string, col int) int {
	// Skip over non word characters
	for col < len(line) && isDelimeter(rune(line[col])) {
		col++
	}

	// Skip over word characters
	for col < len(line) && isWordChar(rune(line[col])) {
		col++
	}

	for col < len(line) && isDelimeter(rune(line[col])) {
		col++
	}

	return col
}

// Function for moving backward by a word
func (m *EditorPaneModel) moveCursorBackwardByWord(line string, col int) int {
	// Skip over non word characters
	for col > 0 && isDelimeter(rune(line[col-1])) {
		col--
	}

	// Skip over word characters
	for col > 0 && isWordChar(rune(line[col-1])) {
		col--
	}

	for col > 0 && isDelimeter(rune(line[col])) {
		col--
	}

	return col
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

		// Insert Query into Editor Pane
	case InsertQueryMsg:
		m.buffer = strings.Split(msg.Query, "\n")
		return m, nil

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
			return m, nil

			// Switch to Insert Mode
		case msg.String() == "esc" && m.mode == InsertMode:
			m.mode = NormalMode
			return m, nil

			// Normal Mode Commands
		case m.mode == NormalMode:
			switch {

			// Move forward by a word
			case msg.String() == "w":
				m.cursorCol = m.moveCursorForwardByWord(m.buffer[m.cursorRow], m.cursorCol)

			// Move backward by a word
			case msg.String() == "b":
				m.cursorCol = m.moveCursorBackwardByWord(m.buffer[m.cursorRow], m.cursorCol)

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
	}

	return m, nil
}

// Function for rendering the line
func (m *EditorPaneModel) renderLine(tokens []token) string {
	var renderedLine strings.Builder

	for _, token := range tokens {
		rawTokenValue := token.Value
		tokenStyle := tokenStyles[token.Type]

		renderedLine.WriteString(tokenStyle.Render(rawTokenValue))
	}

	return renderedLine.String()
}

// Function for rendering the cursor placement
func (m *EditorPaneModel) renderCursor(tokens []token, row int) string {
	var renderedLine strings.Builder
	cursor := "█"
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(rose))

	// Calculate the total character count across tokens to track cursor placement
	charCount := 0

	for _, token := range tokens {
		rawTokenValue := token.Value
		tokenLength := len(rawTokenValue)
		tokenStyle := tokenStyles[token.Type]

		// Check if the cursor is in the token and render it in the correct spot
		if row == m.cursorRow && m.cursorCol >= charCount && m.cursorCol < charCount+tokenLength {
			cursorPos := m.cursorCol - charCount

			if m.mode == NormalMode {
				// Normal Mode: Insert block cursor at the appropriate position
				renderedLine.WriteString(tokenStyle.Render(rawTokenValue[:cursorPos]))
				renderedLine.WriteString(cursorStyle.Render(cursor))
				if cursorPos+1 < tokenLength {
					renderedLine.WriteString(tokenStyle.Render(rawTokenValue[cursorPos+1:]))
				}
			} else if m.mode == InsertMode {
				// Insert Mode: Highlight the character under the cursor
				renderedLine.WriteString(tokenStyle.Render(rawTokenValue[:cursorPos]))

				charUnderCursor := string(rawTokenValue[cursorPos])
				highlightedCharStyle := lipgloss.NewStyle().
					Background(lipgloss.Color(rose)).
					Foreground(lipgloss.Color(overlay))

				renderedLine.WriteString(highlightedCharStyle.Render(charUnderCursor))

				if cursorPos+1 < tokenLength {
					renderedLine.WriteString(tokenStyle.Render(rawTokenValue[cursorPos+1:]))
				}
			}
		} else {
			// Pane is inactive or the cursor is not on this token: render token normally
			renderedLine.WriteString(tokenStyle.Render(rawTokenValue))
		}

		charCount += tokenLength
	}

	// Case for when the cursor is beyond the last token in the line
	if row == m.cursorRow && m.cursorCol >= charCount {
		renderedLine.WriteString(cursorStyle.Render(cursor))
	}

	return renderedLine.String()
}

func (m *EditorPaneModel) View() string {
	var paneStyle lipgloss.Style
	if m.isActive {
		paneStyle = m.activeStyles
	} else {
		paneStyle = m.styles
	}

	var output strings.Builder

	for i, line := range m.buffer {
		// Tokenize line for syntax highlighting
		tokens := tokenize(line)

		if m.isActive && i == m.cursorRow {
			// Render the line with cursor
			output.WriteString(m.renderCursor(tokens, i))
		} else {
			// Render the line normally
			output.WriteString(m.renderLine(tokens))
		}

		// Add new line unless it's the last line
		if i < len(m.buffer)-1 {
			output.WriteString("\n")
		}
	}

	return paneStyle.Render(output.String())
}
