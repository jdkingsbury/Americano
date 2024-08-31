package panes

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type PaneModel struct {
	styles   PaneStyles
	width    int
	height   int
	textarea textarea.Model
	err      error
}

func (m *PaneModel) Init() tea.Cmd {
	return nil
}

func (m *PaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.styles = CreatePaneStyles(m.width, m.height)
		m.resizeTextArea()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// Pane Title Text and Background Color
var titleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#e0def4")).
	Background(lipgloss.Color("#26233a")).
	Align(lipgloss.Center).
	Padding(0, 1)

func (m PaneModel) View() string {
	s := m.styles

	titleLeftPane := titleStyle.Render("Top Left Pane")
	titleBottomPane := titleStyle.Render("Bottom Pane")

	// Render the panes
	topLeftPane := s.TopLeftPane.Render(titleLeftPane)
	mainPane := s.MainPane.Render(m.textarea.View())
	bottomPane := s.BottomPane.Render(titleBottomPane)

	// Arrange Panes
	leftSide := lipgloss.JoinVertical(lipgloss.Top, topLeftPane)
	rightSide := lipgloss.JoinVertical(lipgloss.Top, mainPane)

	layout := lipgloss.JoinHorizontal(lipgloss.Top, leftSide, rightSide)
	layout = lipgloss.JoinVertical(lipgloss.Top, layout, bottomPane)

	return layout
}
