package panes

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdkingsbury/americano/internal/drivers"
)

var (
	treeTitleStyle        = lipgloss.NewStyle().MarginLeft(2).Bold(true).Foreground(lipgloss.Color(text))
	treeItemStyle         = lipgloss.NewStyle().Padding(0, 1)
	treeSelectedItemStyle = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color(rose)).Background(lipgloss.Color(highlightLow))
)

const (
	openCaret   = "▾" // Downward caret for open state
	closedCaret = "▸" // Rightward caret for closed state
)

type DBTreeMsg struct {
	Notification string
	Error        error
}

type ListItem struct {
	Title    string
	SubItems []ListItem
	IsOpen   bool
	Query    string
}

// FlatListItem is used for the rendering the list items
type FlatListItem struct {
	Title     string
	Level     int
	IsOpen    bool
	IsSubItem bool
}

type DBTreeModel struct {
	originalList []ListItem
	flatList     []FlatListItem
	cursor       int
}

func NewDBTreeModel(db drivers.Database) *DBTreeModel {
	var originalList []ListItem

	if db == nil {
		originalList = []ListItem{
			{Title: "No connection"},
		}
	} else {
		dbName, err := db.GetDatabaseName()
		if err != nil {
			dbName = "Unknown db"
		}

		originalList = []ListItem{
			{
				Title:    dbName,
				SubItems: buildDBTree(db),
				IsOpen:   false,
			},
		}
	}

	flatList := flattenList(originalList, 0)

	return &DBTreeModel{
		originalList: originalList,
		flatList:     flatList,
		cursor:       0,
	}
}

// TODO: Work on adding saved queries table when save query feature is added
func buildDBTree(db drivers.Database) []ListItem {
	tables, err := db.GetTables()
	if err != nil {
		return []ListItem{
			{Title: "No connection"},
		}
	}

	tablesItem := ListItem{
		Title:    "Tables",
		IsOpen:   false,
		SubItems: buildTableList(tables),
	}

	// savedQueriesItem := ListItem {
	//   Title: "Saved Queries",
	//   IsOpen: false,
	// }

	return []ListItem{
		tablesItem,
		// savedQueriesItem,
	}
}

func buildTableList(tables []string) []ListItem {
	var tableItems []ListItem
	for _, table := range tables {
		tableItems = append(tableItems, ListItem{
			Title:    table,
			SubItems: buildTableSubItems(table),
			IsOpen:   false,
		})
	}

	return tableItems
}

// Sub items containing queries for tables
func buildTableSubItems(table string) []ListItem {
	return []ListItem{
		{Title: " list", Query: fmt.Sprintf("SELECT * FROM %s;", table)},
		{Title: " column", Query: fmt.Sprintf("PRAGMA table_info(%s);", table)},
		{Title: " foreign key", Query: fmt.Sprintf("PRAGMA foreign_key_list(%s);", table)},
	}
}

func flattenList(items []ListItem, level int) []FlatListItem {
	var flatList []FlatListItem
	for _, item := range items {
		flatItem := FlatListItem{
			Title:     item.Title,
			Level:     level,
			IsOpen:    item.IsOpen,
			IsSubItem: len(item.SubItems) > 0,
		}
		flatList = append(flatList, flatItem)

		// If the item is open and has subitems, recursively flatten the subitems
		if item.IsOpen && len(item.SubItems) > 0 {
			flatList = append(flatList, flattenList(item.SubItems, level+1)...)
		}
	}
	return flatList
}

// Toggles and item's open/collapse state and rebuilds the flat list
func (m *DBTreeModel) toggleItemOpen() {
	// Find the item in the original list
	m.updateOriginalListState(m.originalList, m.flatList[m.cursor].Title, 0, m.flatList[m.cursor].Level)

	// Rebuild the flat list based on the updated original list
	m.flatList = flattenList(m.originalList, 0)
}

func getQueryForItem(title string, items []ListItem) string {
	for _, item := range items {
		if item.Title == title && item.Query != "" {
			return item.Query
		}
		if len(item.SubItems) > 0 {
			query := getQueryForItem(title, item.SubItems)
			if query != "" {
				return query
			}
		}
	}
	return ""
}

// NOTE: Function is returning a bool to ensure recursion stops when found and also for testing when tests are added.

// Update the collapsible state and return true or false if found.
func (m *DBTreeModel) updateOriginalListState(items []ListItem, title string, currentLevel, targetLevel int) bool {
	for i := range items {
		// Check if this is the correct item based on the title and level
		if items[i].Title == title && currentLevel == targetLevel {
			// Toggle open state
			items[i].IsOpen = !items[i].IsOpen
			return true // Exit after toggling
		}

		// Recursively check subitems if they exist
		if len(items[i].SubItems) > 0 {
			// If the item is found in the sublist, return true to stop recursion
			if m.updateOriginalListState(items[i].SubItems, title, currentLevel+1, targetLevel) {
				return true
			}
		}
	}

	return false // Return false if the item was not found in the branch
}

func renderFlatList(flatList []FlatListItem, cursor int) string {
	var b strings.Builder

	title := treeTitleStyle.Render("Database Connection Tree")
	b.WriteString(fmt.Sprintf("%s\n", title))

	for i, item := range flatList {
		indent := strings.Repeat("  ", item.Level)
		caret := ""

		if item.IsOpen && item.IsSubItem {
			caret = openCaret
		} else if !item.IsOpen && item.IsSubItem {
			caret = closedCaret
		}

		if i == cursor {
			b.WriteString(fmt.Sprintf("%s %s %s\n", indent, caret, treeSelectedItemStyle.Render("> "+item.Title))) // Highlight the selected item
		} else {
			b.WriteString(fmt.Sprintf("%s %s %s\n", indent, caret, treeItemStyle.Render(item.Title))) // Normal item
		}
	}

	return b.String()
}

func (m *DBTreeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.flatList)-1 {
				m.cursor++
			}

		case "enter", " ":
			// Check if the selected item has an associated query
			query := getQueryForItem(m.flatList[m.cursor].Title, m.originalList)
			if query != "" {
				// Send the message with the query for the editor
				return m, func() tea.Msg {
					return InsertQueryMsg{Query: query}
				}
			} else {
				m.toggleItemOpen()
			}
		}
	}
	return m, nil
}

func (m *DBTreeModel) Init() tea.Cmd {
	return nil
}

func (m *DBTreeModel) View() string {
	return renderFlatList(m.flatList, m.cursor)
}
