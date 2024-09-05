package panes

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type DBConnInputModel struct {
	focusIndex   int
	inputs       []textinput.Model
	cursorMode   cursor.Mode
	isAddingConn bool
}

func NewConnInputModel(width int) *DBConnInputModel {
	inputs := make([]textinput.Model, 2)

	for i := range inputs {
		ti := textinput.New()
		ti.Cursor.Style = cursorStyle
		ti.CharLimit = 150
		ti.Width = width/4 - 2

		switch i {
		case 0:
			ti.Placeholder = "Enter Connection Name"
			ti.Focus()
			ti.PromptStyle = focusedStyle
			ti.TextStyle = focusedStyle
		case 1:
			ti.Placeholder = "Enter Connection URL"
			ti.Focus()
			ti.PromptStyle = focusedStyle
			ti.TextStyle = focusedStyle
		}

		inputs[i] = ti
	}

	m := &DBConnInputModel{
		inputs: inputs,
	}

	return m
}

func (m *DBConnInputModel) addConnection(name, url string) {
	if name == "" || url == "" {
		fmt.Println("Name and url cannot be empty")
	} else {
		fmt.Printf("Added connection: %s (%s\n", name, url)
	}
}

func (m *DBConnInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *DBConnInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()
			if s == "enter" && m.focusIndex == len(m.inputs) {
				name := m.inputs[0].Value()
				url := m.inputs[1].Value()
				m.addConnection(name, url)
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *DBConnInputModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *DBConnInputModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}
