package panes

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

/* Side Bar Types */

// Bubble Tea list types for Side Bar
const (
	listHeight = 14
)

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

// SideBar Item types
type SideBarItem struct {
	Name     string
	IsButton bool
}
type SideBarView int

// SideBar Views
const (
	ConnectionsView SideBarView = iota
	DBTreeView
)

// Database Connection types for Side Bar
type DatabaseConnection struct {
	Name string
	URL  string
}
