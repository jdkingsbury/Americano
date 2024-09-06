package panes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DBTreeModel struct{}

func NewDBTreeModel() *DBTreeModel {
	view := &DBTreeModel{}

	return view
}

func (m *DBTreeModel) Init() tea.Cmd {
	return nil
}

func (m *DBTreeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// switch msg := msg.(type) {
	// case tea.KeyMsg:
	// 	fmt.Println("Key press in DB Tree view:", msg.String())
	// }
	return m, nil
}

func (m *DBTreeModel) View() string {
	return lipgloss.NewStyle().Render("DB Tree View")
}
