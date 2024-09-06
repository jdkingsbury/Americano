package panes

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(text)).Bold(true).Padding(0, 1)   // Rose for focused input
	blurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(muted)).Faint(true).Padding(0, 1) // Muted for unfocused input
	submitStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(rose)).Bold(true).Padding(0, 1)   // Rose for the submit button
	blurredSubmit = lipgloss.NewStyle().Foreground(lipgloss.Color(muted)).Faint(true).Padding(0, 1) // Muted for inactive submit button
)

type CancelFormMsg struct{}

type SubmitFormMsg struct {
	Name string
	URL  string
}

type DBFormModel struct {
	focusIndex int
	inputs     []textinput.Model
	submit     string
}

func NewDBFormModel() *DBFormModel {
	m := DBFormModel{
		inputs: make([]textinput.Model, 2),
		submit: "[ Submit ]", // Initialize the submit button label
	}

	// Initialize text inputs
	var ti textinput.Model
	for i := range m.inputs {
		ti = textinput.New()
		ti.CharLimit = 156
		ti.Width = 20

		switch i {
		case 0:
			ti.Placeholder = "Enter Connection Name"
			ti.Focus() // Focus on the first input initially
		case 1:
			ti.Placeholder = "Enter Connection URL"
		}

		m.inputs[i] = ti // Assign the initialized textinput.Model back to the slice
	}

	return &m
}

func (m *DBFormModel) Reset() {
	m.focusIndex = 0
	for i := range m.inputs {
		m.inputs[i].SetValue("") // Clear input
		if i == 0 {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
}

func (m *DBFormModel) Init() tea.Cmd {
	return textinput.Blink // Enable blinking cursor for focused input
}

func (m *DBFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg {
				return CancelFormMsg{}
			}
		case "tab": // Tab moves forward through inputs and submit button
			m.focusIndex = (m.focusIndex + 1) % (len(m.inputs) + 1) // Include submit button
		case "shift+tab": // Shift+Tab moves backward through inputs and submit button
			m.focusIndex = (m.focusIndex - 1 + len(m.inputs) + 1) % (len(m.inputs) + 1)
		case "down": // Down arrow moves forward through inputs and submit button
			m.focusIndex = (m.focusIndex + 1) % (len(m.inputs) + 1)
		case "up": // Up arrow moves backward through inputs and submit button
			m.focusIndex = (m.focusIndex - 1 + len(m.inputs) + 1) % (len(m.inputs) + 1)
		case "enter":
			if m.focusIndex == len(m.inputs) {
				// Submit button is focused, handle form submission
				return m, func() tea.Msg {
					return SubmitFormMsg{
						Name: m.inputs[0].Value(),
						URL:  m.inputs[1].Value(),
					}
				}
			}
		}

		// Update focus for inputs
		for i := range m.inputs {
			if i == m.focusIndex {
				m.inputs[i].Focus()
			} else {
				m.inputs[i].Blur()
			}
		}
	}

	// Update all inputs
	for i := range m.inputs {
		var cmd tea.Cmd
		m.inputs[i], cmd = m.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *DBFormModel) View() string {
	var output string

	// Render all input fields
	for i := range m.inputs {
		if i == m.focusIndex {
			output += focusedStyle.Render(m.inputs[i].View()) + "\n"
		} else {
			output += blurredStyle.Render(m.inputs[i].View()) + "\n"
		}
	}

	// Render submit button
	if m.focusIndex == len(m.inputs) { // Focused state for submit button
		output += submitStyle.Render("\n[ Submit ]\n")
	} else {
		output += blurredStyle.Render("\nSubmit\n")
	}

	return output
}
