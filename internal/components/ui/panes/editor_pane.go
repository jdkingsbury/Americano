package panes

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdkingsbury/americano/msgtypes"
)

type EditorPaneModel struct {
	styles       lipgloss.Style
	activeStyles lipgloss.Style
	width        int
	height       int
	textarea     textarea.Model
	err          error
	focused      bool
	isActive     bool
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

	case tea.KeyMsg:
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

		case tea.KeyCtrlC:
			return m, tea.Quit

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

func NewEditorPane(width, height int) *EditorPaneModel {
	ti := textarea.New()
	ti.Placeholder = "Enter SQL Code Here..."
	ti.Focus()
	ti.CharLimit = 1000
	ti.ShowLineNumbers = false

	editorPaneStyle := lipgloss.NewStyle().
		Width(width - 50).
		Height(height - 17).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(iris))

	activePaneStyle := lipgloss.NewStyle().
		Width(width - 50).
		Height(height - 17).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(love))

	return &EditorPaneModel{
		styles:       editorPaneStyle,
		activeStyles: activePaneStyle,
		width:        width,
		height:       height,
		textarea:     ti,
		err:          nil,
		focused:      true,
	}
}

func (m *EditorPaneModel) resizeTextArea() {
	m.textarea.SetWidth(m.width - 50)
	m.textarea.SetHeight(m.height - 17)
}

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
