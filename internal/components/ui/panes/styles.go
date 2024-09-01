package panes

import (
	"github.com/charmbracelet/lipgloss"
)

// Font Icons
var (
	caretRight = ""
	caretdown  = ""
	dbIcon     = ""
	dbAdd      = "󰆺"
	dbConn     = "󱘩"
	dbNotConn  = "󰴀"
)

// Colors are from the Rose-Pine Colorscheme
var (
	rose    = lipgloss.Color("#ebbcba")
	gold    = lipgloss.Color("#f6c177")
	iris    = lipgloss.Color("#c4a7e7")
	pine    = lipgloss.Color("#31748f")
	foam    = lipgloss.Color("#9ccfd8")
	subtle  = lipgloss.Color("#908caa")
	love    = lipgloss.Color("#eb6f92")
	overlay = lipgloss.Color("#26233a")
	surface = lipgloss.Color("#1f1d2e")
)

type PaneStyles struct {
	StatusPane  lipgloss.Style
	TopLeftPane lipgloss.Style
	BottomPane  lipgloss.Style
	MainPane    lipgloss.Style
	// footer      lipgloss.Style
}

func CreatePaneStyles(width, height int) PaneStyles {
	statusHeight := height / 4
	topLeftHeight := height - 17

	s := PaneStyles{
		StatusPane:  lipgloss.NewStyle().Width(45).Height(statusHeight).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(iris)),
		TopLeftPane: lipgloss.NewStyle().Width(45).Height(topLeftHeight).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(iris)),
		BottomPane:  lipgloss.NewStyle().Width(width - 3).Height(height / 3).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(iris)),
		MainPane:    lipgloss.NewStyle().Width(width - 50).Height(height - 17).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(iris)),
		// footer:      lipgloss.NewStyle().Background(lipgloss.Color("#26233a")).Foreground(lipgloss.Color("#e0def4")).Padding(0, 1).Align(lipgloss.Center).Width(width - 3),
	}
	return s
}
