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
	keyBindingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")) // Customize color (example: pink)
	helpTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))  // Customize color (example: green)
)

type KeyMap struct {
	shortHelp []key.Binding
	fullHelp  [][]key.Binding
}

type FooterModel struct {
	style  lipgloss.Style
	help   help.Model
	keyMap KeyMap
	width  int
	height int
}
type SetKeyMapMsg struct {
	FullHelpKeys  [][]key.Binding
	ShortHelpKeys []key.Binding
}

func NewFooterPane(width int) *FooterModel {
	s := lipgloss.NewStyle().
		Width(width).
		Height(1).
		Foreground(lipgloss.Color(text)).
		Padding(0, 1)

	footer := &FooterModel{
		style:  s,
		help:   help.New(),
		width:  width,
		height: 1,
	}
	return footer
}

func (k KeyMap) ShortHelp() []key.Binding {
	return k.shortHelp
}

// NOTE: Not using full help at the moment
func (k KeyMap) FullHelp() [][]key.Binding {
	return k.fullHelp
}

// NOTE: Will update to allow for full key maps
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
	var shortHelp []string

	// Apply custom styles to each keybinding in short help
	for _, kb := range m.keyMap.ShortHelp() {
		shortHelp = append(shortHelp, keyBindingStyle.Render(kb.Help().Key)+": "+helpTextStyle.Render(kb.Help().Desc)+" | ")
	}

	// Apply custom styles to each keybinding in full help
	// for _, section := range m.keyMap.FullHelp() {
	//     for _, kb := range section {
	//         fullHelp = append(fullHelp, keyBindingStyle.Render(kb.Help().Key)+": "+helpTextStyle.Render(kb.Help().Desc))
	//     }
	// }

	// Join the styled keybindings into one string
	return lipgloss.JoinHorizontal(lipgloss.Top, shortHelp...)
}

// func (m *FooterModel) View() string {
// 	return m.help.View(m.keyMap)
// }
