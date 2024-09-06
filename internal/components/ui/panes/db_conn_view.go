package panes

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2).Bold(true).Foreground(lipgloss.Color(text))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color(rose))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type DBConnItems struct {
	Name     string
	URL      string
	isButton bool
}

func (i DBConnItems) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(DBConnItems)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.Name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type DBConnModel struct {
	list       list.Model
	choice     DBConnItems
	focusIndex int
}

func NewDBConnModel(width int) *DBConnModel {
	items := []list.Item{
		DBConnItems{Name: "󰆺 Add Connection", URL: "", isButton: true},
	}

	li := list.New(items, itemDelegate{}, width/4, listHeight)
	li.Title = "Database Connections"
	li.SetShowStatusBar(false)
	li.SetFilteringEnabled(false)
	li.SetShowHelp(false)
	li.Styles.Title = titleStyle
	li.Styles.PaginationStyle = paginationStyle

	pane := &DBConnModel{
		list: li,
	}

	return pane
}

func (m *DBConnModel) AddConnection(name, url string) {
	m.list.InsertItem(len(m.list.Items()), DBConnItems{Name: name, URL: url, isButton: false})
}

func (m *DBConnModel) FocusedOnButton() bool {
	item, ok := m.list.SelectedItem().(DBConnItems)
	return ok && item.isButton
}

func (m *DBConnModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			item, ok := m.list.SelectedItem().(DBConnItems)
			if ok {
				if item.isButton {
					// Notify SidebarPane model that the Add Connection button was clicked
          fmt.Println("Button Clicked")
					// return m, func() tea.Msg { return SubmitFormMsg{} }
				} else if item.URL != "" {
					m.choice = item
					fmt.Println(item.Name)
				}
			}
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *DBConnModel) Init() tea.Cmd {
	return nil
}

func (m *DBConnModel) View() string {
	var content string

	content += titleStyle.Render(m.list.Title)
	content += m.list.View()

	return content
}
