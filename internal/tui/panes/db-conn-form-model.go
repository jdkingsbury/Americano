package panes

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	formTitleStyle    = lipgloss.NewStyle().MarginLeft(2).Bold(true).Foreground(lipgloss.Color(text))
	formFocusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(rose)).Bold(true).Padding(0, 1)    // Rose for focused input
	formBlurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(subtle)).Faint(true).Padding(0, 1) // Muted for unfocused input
	formSubmitStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(rose)).Bold(true).Padding(0, 1)    // Rose for the submit button
	formBlurredSubmit = lipgloss.NewStyle().Foreground(lipgloss.Color(muted)).Faint(true).Padding(0, 1)  // Muted for inactive submit button
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
	title      string
	keys       dbFormKeyMap
}

type dbFormKeyMap struct {
	CancelForm key.Binding
	NextInput  key.Binding
	PrevInput  key.Binding
	SubmitForm key.Binding
}

func newDBFormKeyMap() dbFormKeyMap {
	return dbFormKeyMap{
		CancelForm: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel form input"),
		),
		NextInput: key.NewBinding(
			key.WithKeys("tab", "down"),
			key.WithHelp("↓/tab", "next input field"),
		),
		PrevInput: key.NewBinding(
			key.WithKeys("shift+tab", "up"),
			key.WithHelp("↑/shift+tab", "previous input field"),
		),
		SubmitForm: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "submit form"),
		),
	}
}

func NewDBFormModel() *DBFormModel {
	m := DBFormModel{
		inputs: make([]textinput.Model, 2),
		submit: "[ Submit ]",
		title:  "Add Connection",
		keys:   newDBFormKeyMap(),
	}

	var ti textinput.Model
	for i := range m.inputs {
		ti = textinput.New()
		ti.CharLimit = 156
		ti.Width = 30

		switch i {
		case 0:
			ti.Placeholder = "Enter Connection Name"
			ti.Focus()
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
		m.inputs[i].SetValue("")
		if i == 0 {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
}

func (m *DBFormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *DBFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.CancelForm):
			return m, func() tea.Msg {
				return CancelFormMsg{}
			}
		case key.Matches(msg, m.keys.NextInput):
			m.focusIndex = (m.focusIndex + 1) % (len(m.inputs) + 1)

		case key.Matches(msg, m.keys.PrevInput):
			m.focusIndex = (m.focusIndex - 1 + len(m.inputs) + 1) % (len(m.inputs) + 1)

		case key.Matches(msg, m.keys.SubmitForm):
			if m.focusIndex == len(m.inputs) {
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

	output += formTitleStyle.Render(m.title) + "\n"

	// Input fields
	for i := range m.inputs {
		if i == m.focusIndex {
			output += formFocusedStyle.Render(m.inputs[i].View()) + "\n"
		} else {
			output += formBlurredStyle.Render(m.inputs[i].View()) + "\n"
		}
	}

	// Button field
	if m.focusIndex == len(m.inputs) { // Focused state for submit button
		output += formSubmitStyle.Render("\n[ Submit ]\n")
	} else {
		output += formBlurredSubmit.Render("\nSubmit\n")
	}

	return output
}
