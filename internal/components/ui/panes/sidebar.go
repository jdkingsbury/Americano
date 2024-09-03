package panes

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdkingsbury/americano/msgtypes"
)

/* Handles the side bar pane */

type SideBarItem struct {
	Name      string
	IsButton  bool
	IsSection bool
	SectionID int // ID to group items under each section
}

type SideBarPaneModel struct {
	styles             lipgloss.Style
	activeStyles       lipgloss.Style
	width              int
	height             int
	cursor             int
	items              []SideBarItem
	isCollapsed        map[int]bool // Track the collapsed state of each section
	err                error
	isActive           bool // Check if the pane is active
	isAddingConnection bool
	textInput          textinput.Model
}

func (m *SideBarPaneModel) updateStyles() {
	m.styles = lipgloss.NewStyle().
		Width(m.width / 4).
		Height(m.height - 17).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(iris))

	m.activeStyles = lipgloss.NewStyle().
		Width(m.width / 4).
		Height(m.height - 17).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(love))
}

func NewSideBarPane(width, height int) *SideBarPaneModel {
	ti := textinput.New()
	ti.Placeholder = keyboard + " Enter Connection URL"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = width/4 - 2

	pane := &SideBarPaneModel{
		width:  width,
		height: height,
		items: []SideBarItem{
			{Name: "Database Connections", IsSection: true, SectionID: 1},
			{Name: " 󰆺 Add Connection", IsButton: true, SectionID: 1},
		},
		isCollapsed: make(map[int]bool),
		err:         nil,
		textInput:   ti,
	}

	pane.updateStyles() // Initialize styles

	return pane
}

func (m *SideBarPaneModel) Init() tea.Cmd {
	return nil
}

func (m *SideBarPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateStyles()

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			if m.isAddingConnection {
				connectionURL := m.textInput.Value()
				fmt.Println("Adding Connection:", connectionURL)
				m.isAddingConnection = false
				m.textInput.Reset()
			} else if m.cursor < len(m.items) {
				item := m.items[m.cursor]
				if item.IsSection {
					m.isCollapsed[item.SectionID] = !m.isCollapsed[item.SectionID]
				} else if item.IsButton {
					m.isAddingConnection = true
					m.textInput.Focus()
				}
			}
		}

	case msgtypes.ErrMsg:
		m.err = msg
		return m, nil
	}

	if m.isAddingConnection {
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, cmd
}

func (m *SideBarPaneModel) View() string {
	var content string

	for i, item := range m.items {
		itemStyle := lipgloss.NewStyle()

		// Highlight item based on cursor position
		if i == m.cursor {
			itemStyle = itemStyle.Foreground(lipgloss.Color(rose)).Bold(true)
		}

		if item.IsSection {
			sectionStyle := itemStyle.Bold(true)
			if m.isCollapsed[item.SectionID] {
				content += sectionStyle.Render(fmt.Sprintf("%s %s", caretRight, item.Name)) + "\n"
			} else {
				content += sectionStyle.Render(fmt.Sprintf("%s %s", caretdown, item.Name)) + "\n"
			}
		} else {
			if !m.isCollapsed[item.SectionID] {
				content += itemStyle.Render(item.Name) + "\n"
			}
		}
	}

	if m.isAddingConnection {
		content += m.textInput.View()
	}

	var paneStyle lipgloss.Style
	if m.isActive {
		paneStyle = m.activeStyles
	} else {
		paneStyle = m.styles
	}

	sideBarPane := paneStyle.Render(content)
	return sideBarPane
}
