package components

import "github.com/charmbracelet/lipgloss"

var (
	rose = lipgloss.Color("#ebbcba")
	iris = lipgloss.Color("#c4a7e7")
	pine = lipgloss.Color("#31748f")
	foam = lipgloss.Color("#9ccfd8")
)

type PaneStyles struct {
	TopLeftPane    lipgloss.Style
	BottomLeftPane lipgloss.Style
	BottomPane     lipgloss.Style
	MainPane       lipgloss.Style
}

func CreateStyles(width, height int) PaneStyles {
	topLeftHeight := height / 3
	bottomLeftHeight := height / 4
	mainPaneWidth := width - 45

	s := PaneStyles{
		TopLeftPane:    lipgloss.NewStyle().Width(40).Height(topLeftHeight).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(rose)).Padding(1),
		BottomLeftPane: lipgloss.NewStyle().Width(40).Height(bottomLeftHeight).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(iris)).Padding(1),
		BottomPane:     lipgloss.NewStyle().Width(width - 3).Height(height / 4).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(pine)).Padding(1),
		MainPane:       lipgloss.NewStyle().Width(mainPaneWidth).Height(topLeftHeight + bottomLeftHeight + 3).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(foam)).Padding(1),
	}
	return s
}
