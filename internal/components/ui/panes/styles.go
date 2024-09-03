package panes

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// Font Icons
var (
	caretRight = ""
	caretDown  = ""
	dbIcon     = ""
	dbAdd      = "󰆺"
	dbConn     = "󱘩"
	dbNotConn  = "󰴀"
	keyboard   = "󰥻"
)

// Colors are from the Rose-Pine Colorscheme
var (
	base          = lipgloss.Color("#191724")
	surface       = lipgloss.Color("#1f1d2e")
	overlay       = lipgloss.Color("#26233a")
	muted         = lipgloss.Color("#6e6a86")
	subtle        = lipgloss.Color("#908caa")
	text          = lipgloss.Color("#e0def4")
	love          = lipgloss.Color("#eb6f92")
	gold          = lipgloss.Color("#f6c177")
	rose          = lipgloss.Color("#ebbcba")
	pine          = lipgloss.Color("#31748f")
	foam          = lipgloss.Color("#9ccfd8")
	iris          = lipgloss.Color("#c4a7e7")
	highlightLow  = lipgloss.Color("#2a283e")
	highlightMed  = lipgloss.Color("#dfdad9")
	highlightHigh = lipgloss.Color("#cecacd")
)

/* Styles For Sidebar */
var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(1).Bold(true).Foreground(lipgloss.Color(text))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(1)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color(rose))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(1)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(1).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type SideBarConfig struct {
	TitleStyle        lipgloss.Style
	ItemStyle         lipgloss.Style
	SelectedItemStyle lipgloss.Style
	PaginationStyle   lipgloss.Style
	HelpStyle         lipgloss.Style
	QuitTextStyle     lipgloss.Style
}

func NewSideBarConfig() *SideBarConfig {
	return &SideBarConfig{
		TitleStyle:        titleStyle,
		ItemStyle:         itemStyle,
		SelectedItemStyle: selectedItemStyle,
		PaginationStyle:   paginationStyle,
		HelpStyle:         helpStyle,
		QuitTextStyle:     quitTextStyle,
	}
}
