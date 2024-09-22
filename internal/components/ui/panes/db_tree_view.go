package panes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

type styles struct {
	base,
	block,
	enumerator,
	dbTable,
	toggle,
	dbNode lipgloss.Style
}

func defaultStyles() styles {
	var s styles
	s.base = lipgloss.NewStyle().
		Foreground(lipgloss.Color("225"))
	s.block = s.base.
		Padding(1, 3).
		Margin(1, 3).
		Width(40)
	s.enumerator = s.base.
		Foreground(lipgloss.Color("212")).
		PaddingRight(1)
	s.dbTable = s.base.
		Inline(true)
	s.toggle = s.base.
		Foreground(lipgloss.Color("207")).
		PaddingRight(1)
	s.dbNode = s.base
	return s
}

type dbTable struct {
	name   string
	open   bool
	styles styles
}

func (d dbTable) String() string {
	t := d.styles.toggle.Render
	n := d.styles.dbTable.Render
	if d.open {
		return t("▼") + n(d.name)
	}
	return t("▶") + n(d.name)
}

type dbNode struct {
	name   string
	styles styles
}

func (s dbNode) String() string {
	return s.styles.dbNode.Render(s.name)
}

type DBTreeModel struct {
	tree *tree.Tree
}

func NewDBTreeModel() *DBTreeModel {
	s := defaultStyles()

	t := tree.Root(dbTable{"Activities", true, s}).
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(s.enumerator).
		Child(
			tree.Root(dbTable{"Sports", true, s}).
				Child(
					dbNode{"Basketball", s},
					dbNode{"Football", s},
					dbNode{"Baseball", s},
				),
		)

	return &DBTreeModel{
		tree: t,
	}
}

func (m *DBTreeModel) Init() tea.Cmd {
	return nil
}

func (m *DBTreeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *DBTreeModel) View() string {
	s := defaultStyles()
	return s.base.Render(m.tree.String())
}
