package panes

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdkingsbury/americano/internal/components/drivers"
)

const listHeight = 14

var (
	listTitleStyle        = lipgloss.NewStyle().MarginLeft(2).Bold(true).Foreground(lipgloss.Color(text))
	listItemStyle         = lipgloss.NewStyle().Padding(0, 1)
	listSelectedItemStyle = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color(rose))
	listPaginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
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

	fn := listItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return listSelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type DBConnModel struct {
	list       list.Model
	choice     DBConnItems
	focusIndex int
	database   drivers.Database
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
	li.Styles.Title = listTitleStyle
	li.Styles.PaginationStyle = listPaginationStyle

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
			// Handle database connection
			item, ok := m.list.SelectedItem().(DBConnItems)
			if ok && item.URL != "" {
				// Connect to the selected database
				db, err := drivers.ConnectToDatabase(item.URL)
				if err != nil {
					fmt.Println("Error connecting to database:", err)
				} else {
					// Store the connected database instance
					m.database = db
					fmt.Printf("Connected to %s\n", item.Name)
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

	content += listTitleStyle.Render(m.list.Title)
	content += m.list.View()

	return content
}
