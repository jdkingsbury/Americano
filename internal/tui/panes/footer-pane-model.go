package panes

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TODO: See if we can use the help bubble tea component to help with keymaps

/* Basic Footer View */

var (
	keyBindingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(pine)).Padding(0, 1).Bold(true)
	helpTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(foam)).Padding(0, 1).Bold(true)
)

type KeyMap struct {
	shortHelp []key.Binding
	fullHelp  [][]key.Binding
}

type FooterModel struct {
	style        lipgloss.Style
	help         help.Model
	keyMap       KeyMap
	width        int
	height       int
	showFullHelp bool
}
type SetKeyMapMsg struct {
	FullHelpKeys  [][]key.Binding
	ShortHelpKeys []key.Binding
}

func NewFooterPane(width int) *FooterModel {
	s := lipgloss.NewStyle().
		Width(width).
		Height(1).
		Padding(0, 1)

	footer := &FooterModel{
		style:        s,
		help:         help.New(),
		width:        width,
		height:       1,
		showFullHelp: false,
	}
	return footer
}

func (k KeyMap) ShortHelp() []key.Binding {
	return k.shortHelp
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return k.fullHelp
}

func (m *FooterModel) SetKeyBindings(fullHelpKeys [][]key.Binding, shortHelpKeys []key.Binding) {
	m.keyMap = KeyMap{
		shortHelp: shortHelpKeys,
		fullHelp:  fullHelpKeys,
	}
}

func (m *FooterModel) Init() tea.Cmd {
	return nil
}

func (m *FooterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SetKeyMapMsg:
		m.SetKeyBindings(msg.FullHelpKeys, msg.ShortHelpKeys)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = 1
	}
	return m, nil
}

func (m *FooterModel) View() string {
	var helpView []string

	if m.showFullHelp {
		for _, section := range m.keyMap.FullHelp() {
			for idx, kb := range section {
				helpView = append(helpView, keyBindingStyle.Render(kb.Help().Key)+": "+helpTextStyle.Render(kb.Help().Desc))
				if idx != len(m.keyMap.FullHelp())-1 {
					helpView = append(helpView, " ")
				}
			}
		}
	} else {
		for idx, kb := range m.keyMap.ShortHelp() {
			helpView = append(helpView, keyBindingStyle.Render(kb.Help().Key)+": "+helpTextStyle.Render(kb.Help().Desc))
			if idx != len(m.keyMap.ShortHelp())-1 {
				helpView = append(helpView, " ")
			}
		}
	}

	// Join the styled keybindings into one string
	return lipgloss.JoinHorizontal(lipgloss.Top, helpView...)
}
