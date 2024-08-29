package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

type styles struct {
	base,
	block,
	enumerator,
	toggle,
	db,
	table,
	column lipgloss.Style
}

func defaultStyles() styles {
	var s styles
	s.base = lipgloss.NewStyle().
		Background(lipgloss.Color("57")).
		Foreground(lipgloss.Color("225"))
	s.block = s.base.
		Padding(1, 3).
		Margin(1, 3).
		Width(40)
	s.enumerator = s.base.
		Foreground(lipgloss.Color("212")).
		PaddingRight(1)
	s.toggle = s.base.
		Foreground(lipgloss.Color("207")).
		PaddingRight(1)
	s.db = s.base
	s.table = s.base
	s.column = s.base
	return s
}

type db struct {
	name   string
	open   bool
	styles styles
}

func (d db) String() string {
	t := d.styles.toggle.Render
	n := d.styles.db.Render
	if d.open {
		return t("▼") + n(d.name)
	}
	return t("▶") + n(d.name)
}

type table struct {
	name   string
	open   bool
	styles styles
}

func (t table) String() string {
	toggle := t.styles.toggle.Render
	tableName := t.styles.table.Render

	if t.open {
		return toggle("▼") + tableName(t.name)
	}
	return toggle("▶") + tableName(t.name)
}

type column struct {
  name string
  styles styles
}

func (c column) String() string {
  return c.styles.column.Render(c.name)
}

func PostgresTreeView() {
  s := defaultStyles()

	t := tree.Root(db{"DB Name", true, s}).
    Enumerator(tree.RoundedEnumerator).
    EnumeratorStyle(s.enumerator).
		Child("New Query").
		Child("Saved Queries").
		Child(
			tree.New().
				Root("Schemas").
				Child("Information Schema").
				Child("PG Catalog").
				Child("PG Toast").
				Child("Public"),
		)

	fmt.Println(s.block.Render(t.String()))
}

// Tree view using lipgloss
func main() {
	PostgresTreeView()
}
