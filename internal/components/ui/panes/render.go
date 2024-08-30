package panes

import "github.com/charmbracelet/lipgloss"

// Pane Title Text and Background Color
var titleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#e0def4")).
	Background(lipgloss.Color("#26233a")).
	Padding(0, 1)

// Render functions for each pane
func TopLeftPane(style lipgloss.Style) string {
	title := titleStyle.Render("Top Left Pane")
	return style.Render(title)
}

func MainPane(style lipgloss.Style) string {
	title := titleStyle.Render("Main Pane")
	return style.Render(title)
}

func BottomPane(style lipgloss.Style) string {
	title := titleStyle.Render("Bottom Pane")
	return style.Render(title)
}
