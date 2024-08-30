package panes

import "github.com/charmbracelet/lipgloss"

// Colors are from the Rose-Pine Colorscheme
var (
	rose = lipgloss.Color("#ebbcba")
	iris = lipgloss.Color("#c4a7e7")
	pine = lipgloss.Color("#31748f")
	foam = lipgloss.Color("#9ccfd8")
  subtle = lipgloss.Color("#908caa")

)

type PaneStyles struct {
	TopLeftPane    lipgloss.Style
	BottomPane     lipgloss.Style
	MainPane       lipgloss.Style
}

func CreatePaneStyles(width, height int) PaneStyles {
	topLeftHeight := height - 17
	mainPaneWidth := width - 45

	s := PaneStyles{
		TopLeftPane:    lipgloss.NewStyle().Width(45).Height(topLeftHeight).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(iris)),
		BottomPane:     lipgloss.NewStyle().Width(width - 3).Height(height / 3).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(iris)),
		MainPane:       lipgloss.NewStyle().Width(mainPaneWidth - 5).Height(height - 17).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(iris)),
	}
	return s
}
