package main

import (
	"github.com/rivo/tview"
)

// Need
// - Connections Pane
// - DB Side Bar
// - SQL Input Bar
// - Table Results pane
// - Shortcut commands on the bottom of the screen
// - Help window with list of commands

// Maybe
// - Saved Queries Pane
// - Query History Pane

// TODO: Check to see how to style the list and place into a pane window
func CreateDBSideBar(app *tview.Application) *tview.List {
	dbSidebar := tview.NewList()

	dbSidebar.AddItem("New query", "", 'n', nil).
		AddItem("Saved Queries", "", 's', nil).
		AddItem("Query History", "", 'h', nil).
		AddItem("Schemas", "", 'c', nil).
    // For testing
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})

	return dbSidebar
}

func CreateLayout(app *tview.Application) *tview.Flex {
	dbSidebar := CreateDBSideBar(app)

	layout := tview.NewFlex().AddItem(dbSidebar, 20, 1, true)

	return layout
}

func main() {
	app := tview.NewApplication()
	layout := CreateLayout(app)

	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
}
