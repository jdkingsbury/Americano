package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication() // Initialize the application

	textView := tview.NewTextView().SetText("Hello World! Press 'q' to quit.")

	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			app.Stop()
		}
		return event
	})

	if err := app.SetRoot(textView, true).Run(); err != nil {
		panic(err)
	}
}
