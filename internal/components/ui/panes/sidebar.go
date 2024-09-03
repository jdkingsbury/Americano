package panes

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdkingsbury/americano/msgtypes"
)

// TODO: Look into tabs to see if it will be a good option of switching
// between connections and the database tree so we don't need to create another pane in the layout

/* Handles the side bar pane */

const (
	listHeight = 14
)

// Define styles
var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(1).Bold(true).Foreground(lipgloss.Color(text))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(1)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color(rose))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(1)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(1).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type SideBarItem struct {
	Name     string
	IsButton bool
}

func (i SideBarItem) FilterValue() string { return i.Name }

// Define custom delegate
type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(SideBarItem)
	if !ok {
		return
	}

	str := i.Name

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type DatabaseConnection struct {
	Name string
	URL  string
}

type SideBarPaneModel struct {
	styles             lipgloss.Style
	activeStyles       lipgloss.Style
	width              int
	height             int
	isActive           bool // Check if the pane is active
	isAddingConnection bool
	list               list.Model
	inputs             []textinput.Model
	focusedIndex       int
	err                error
	connections        []DatabaseConnection
}

func (m *SideBarPaneModel) addConnection(name, url string) {
	connection := DatabaseConnection{Name: name, URL: url}
	m.connections = append(m.connections, connection)

	// Append the new connection to the list
	newItem := SideBarItem{Name: " 󰇯 " + name}
	m.list.InsertItem(len(m.list.Items()), newItem)

	// Set the new item as selected
	m.list.Select(len(m.list.Items()) - 1)

	// Reset after adding the connection
	m.isAddingConnection = false
	m.inputs[0].SetValue("")
	m.inputs[1].SetValue("")
	m.focusedIndex = 0
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
	inputs := make([]textinput.Model, 2)

	placeholders := []string{"Enter Name", "Enter Connection URL"}

	for i := range inputs {
		ti := textinput.New()
		ti.Placeholder = placeholders[i]
		ti.CharLimit = 256
		ti.Width = width/4 - 2
		if i == 0 {
			ti.Focus()
		}
		inputs[i] = ti
	}

	items := []list.Item{
		SideBarItem{Name: "󰆺 Add Connection", IsButton: true},
	}

	li := list.New(items, itemDelegate{}, width/4, listHeight)
	li.Title = "Database Connections"
	li.SetShowStatusBar(false)
	li.SetFilteringEnabled(false)
	li.SetShowHelp(false) // Disable help text
	li.Styles.Title = titleStyle
	li.Styles.PaginationStyle = paginationStyle

	pane := &SideBarPaneModel{
		width:        width,
		height:       height,
		list:         li,
		inputs:       inputs,
		focusedIndex: 0,
		err:          nil,
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
		case "enter":
			if m.isAddingConnection {
				name := m.inputs[0].Value()
				url := m.inputs[1].Value()
				m.addConnection(name, url)
				m.inputs[0].SetValue("")
				m.inputs[1].SetValue("")
				m.focusedIndex = 0
				m.isAddingConnection = false
			} else {
				selectedItem := m.list.SelectedItem().(SideBarItem)
				if selectedItem.IsButton {
					m.isAddingConnection = true
				}
			}

		case "tab":
			if m.isAddingConnection && m.focusedIndex > 0 {
				m.focusedIndex--
			} else if m.isAddingConnection && m.focusedIndex < len(m.inputs)-1 {
				m.focusedIndex++
			}
		}

	case msgtypes.ErrMsg:
		m.err = msg
		return m, nil
	}

	if m.isAddingConnection {
		for i := range m.inputs {
			if i == m.focusedIndex {
				m.inputs[i].Focus()
			} else {
				m.inputs[i].Blur()
			}
			m.inputs[i], cmd = m.inputs[i].Update(msg)
		}
		return m, cmd
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *SideBarPaneModel) View() string {
	var content string

	if m.isAddingConnection {
		for _, input := range m.inputs {
			content += input.View() + "\n"
		}
	} else {
		content += titleStyle.Render(m.list.Title) + "\n"
		content += m.list.View()
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
