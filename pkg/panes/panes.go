package panes

import "github.com/charmbracelet/lipgloss"

// Render functions for each pane
func TopLeftPane(style lipgloss.Style) string {
	return style.Render("Top Left Pane")
}

func BottomLeftPane(style lipgloss.Style) string {
	return style.Render("Bottom Left Pane")
}

func MainPane(style lipgloss.Style) string {
	return style.Render("Main Pane")
}

func BottomPane(style lipgloss.Style) string {
	return style.Render("Bottom Pane")
}
