package panes

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ResultPaneModel struct {
	styles       lipgloss.Style
	activeStyles lipgloss.Style
	width        int
	height       int
	err          error
	isActive     bool
	table        table.Model
}

// Initialize Result Pane
func NewResultPaneModel(width, height int) *ResultPaneModel {
	columns := []table.Column{
		{Title: "Column1", Width: 10},
		{Title: "Column2", Width: 10},
	}

	rows := []table.Row{}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(height/3),
	)

	pane := &ResultPaneModel{
		width:  width,
		height: height,
		err:    nil,
		table:  t,
	}

	pane.updateStyles()

	return pane
}

func (m *ResultPaneModel) UpdateTable(columns []string, rowData [][]string) {
	// Create table columns from the column names
	tableColumns := []table.Column{}
	for _, col := range columns {
		tableColumns = append(tableColumns, table.Column{Title: col, Width: 15})
	}

  // Convert row data to the format expected by the table component
	tableRows := []table.Row{}
	for _, row := range rowData {
		tableRows = append(tableRows, table.Row(row))
	}

  // Update the table model with the new columns and rows
	m.table.SetColumns(tableColumns)
	m.table.SetRows(tableRows)
}

// Code for changing from active to inactive window
func (m *ResultPaneModel) updateStyles() {
	m.styles = lipgloss.NewStyle().
		Width(m.width - 3).
		Height(m.height / 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(iris))

	m.activeStyles = lipgloss.NewStyle().
		Width(m.width - 3).
		Height(m.height / 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(rose))
}

// Code for functionality on start
func (m *ResultPaneModel) Init() tea.Cmd {
	return nil
}

// Code for updating the state
func (m *ResultPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateStyles()
		m.table.SetHeight(m.height / 3)

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q":
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

// Result Pane View
func (m *ResultPaneModel) View() string {
	var paneStyle lipgloss.Style
	if m.isActive {
		paneStyle = m.activeStyles
	} else {
		paneStyle = m.styles
	}

	resultPane := paneStyle.Render(m.table.View())

	if m.err != nil {
		resultPane += lipgloss.NewStyle().
			Foreground(lipgloss.Color("red")).
			Render(m.err.Error())
	}

	return resultPane
}
