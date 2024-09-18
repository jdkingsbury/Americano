package panes

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdkingsbury/americano/internal/drivers"
)

// NOTE: May need to change how we update the width and height of the table
type ClearNotificationMsg struct{}

type ResultPaneModel struct {
	styles       lipgloss.Style
	activeStyles lipgloss.Style
	width        int
	height       int
	err          error
	isActive     bool
	table        table.Model
	notification string
}

// Initialize Result Pane
func NewResultPaneModel(width, height int) *ResultPaneModel {
	columns := []table.Column{}
	rows := []table.Row{}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Apply table styles
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(iris)).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(rose)).
		Background(lipgloss.Color(highlightLow)).
		Bold(false)
	t.SetStyles(s)

	pane := &ResultPaneModel{
		width:  width,
		height: height,
		err:    nil,
		table:  t,
	}

	pane.updateStyles()

	return pane
}

// NOTE: Function is for testing the table
func (m *ResultPaneModel) TestResultPaneTable() {
	columns := []string{"ID", "Name", "Age", "Occupation", "Country"}

	rows := [][]string{
		{"1", "Alice", "29", "Engineer", "USA"},
		{"2", "Bob", "34", "Designer", "UK"},
		{"3", "Charlie", "22", "Student", "Canada"},
		{"4", "David", "40", "Manager", "Australia"},
		{"5", "Eve", "35", "Scientist", "Germany"},
	}
	m.UpdateTable(columns, rows)
}

func (m *ResultPaneModel) UpdateTable(columns []string, rowData [][]string) {
	if len(columns) == 0 {
		fmt.Println("No columns to display")
		return
	}

	// Calculate the available width for the table
	availableWidth := m.width - 16
	if availableWidth < 0 {
		availableWidth = 0
	}

	// Calculate column width dynamically
	columnWidth := availableWidth / len(columns)
	if columnWidth < 1 {
		columnWidth = 1
	}

	// fmt.Printf("Available width: %d, Column width: %d", availableWidth, columnWidth)

	// Create table columns from the column names
	tableColumns := []table.Column{}
	for _, col := range columns {
		tableColumns = append(tableColumns, table.Column{Title: col, Width: columnWidth})
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
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case drivers.DBConnMsg:

    // Appends the cmd to clear notification and error before displaying new result
		cmds = append(cmds, func() tea.Msg {
			return ClearNotificationMsg{}
		})

		if msg.Error != nil {
			m.err = msg.Error
			m.notification = ""
		} else {
			m.notification = msg.Notification
			m.err = nil
		}

	case ClearNotificationMsg:
		m.notification = ""
		m.err = nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateStyles()
		m.table.SetHeight((m.height / 3) - 3)

		// Recalculate column widths
		availableWidth := m.width - 16 // Account for borders and padding
		columns := m.table.Columns()   // Get the existing columns
		columnWidth := availableWidth / len(columns)

		for i := range columns {
			columns[i].Width = columnWidth
		}
		m.table.SetColumns(columns)

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "enter":
			m.TestResultPaneTable()
		case "q":
			return m, tea.Quit
		}
	}
	var tableCmd tea.Cmd
	m.table, tableCmd = m.table.Update(msg)
	cmds = append(cmds, tableCmd)

	return m, tea.Batch(cmds...)
}

// Result Pane View
func (m *ResultPaneModel) View() string {
	var paneStyle lipgloss.Style
	if m.isActive {
		paneStyle = m.activeStyles
	} else {
		paneStyle = m.styles
	}

	if m.err != nil {
		return paneStyle.Render(lipgloss.NewStyle().
			Foreground(lipgloss.Color(rose)).
			Render(m.err.Error()),
		)
	}

	if m.notification != "" {
		return paneStyle.Render(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(text)).
				Render(m.notification),
		)
	}

	return paneStyle.Render(m.table.View())
}
