package panes

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jdkingsbury/americano/internal/drivers"
	"github.com/jdkingsbury/americano/msgtypes"
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
	keys         resultKeyMaps
}

type resultKeyMaps struct {
	Focus key.Binding
}

func newResultKeyMaps() resultKeyMaps {
	return resultKeyMaps{
		Focus: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "toggle table focus"),
		),
	}
}

func (m ResultPaneModel) KeyMap() []key.Binding {
	return []key.Binding{m.keys.Focus}
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
		keys:   newResultKeyMaps(),
	}

	pane.updateStyles()

	return pane
}

// Used for testing the tables in the result pane
func (m *ResultPaneModel) Table() table.Model {
	return m.table
}

func (m *ResultPaneModel) UpdateTable(columns []string, rowData [][]string) {
	if len(columns) == 0 {
		msgtypes.NewNotificationMsg("No columns to display")
		return
	}

	// Constants for border, padding, and column spacing
	borderWidth := 6
	padding := 6
	columnSpacing := 1

	// Calculate the available width for the table
	availableWidth := m.width - borderWidth - padding - (len(columns)-1)*columnSpacing
	if availableWidth < 0 {
		availableWidth = 0
	}

	minColumnWidth := 10

	// Calculate column width dynamically
	columnWidth := availableWidth / len(columns)
	if columnWidth < minColumnWidth {
		columnWidth = minColumnWidth
	}

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

// Styles for result pane
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

func (m *ResultPaneModel) Init() tea.Cmd {
	return nil
}

func (m *ResultPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case drivers.QueryResultMsg:
		cmds = append(cmds, func() tea.Msg {
			return ClearNotificationMsg{}
		})

		if msg.Error != nil {
			m.err = msg.Error
			return m, nil
		}

		m.UpdateTable(msg.Columns, msg.Rows)

	case msgtypes.NotificationMsg:
		cmds = append(cmds, func() tea.Msg {
			return ClearNotificationMsg{}
		})

		m.notification = msg.Notification
		m.err = nil

	case msgtypes.ErrMsg:
		cmds = append(cmds, func() tea.Msg {
			return ClearNotificationMsg{}
		})

		m.notification = ""
		m.err = msg.Err

	// Msg for clearing notifications and errors in result pane
	case ClearNotificationMsg:
		m.notification = ""
		m.err = nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateStyles()
		m.table.SetHeight((m.height / 3) - 3)

		// Calculate column widths
		availableWidth := m.width - 16 // Account for borders and padding
		columns := m.table.Columns()   // Get the existing columns
		columnWidth := availableWidth / len(columns)

		for i := range columns {
			columns[i].Width = columnWidth
		}
		m.table.SetColumns(columns)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Focus):
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
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
		// For displaying errors
		return paneStyle.Render(lipgloss.NewStyle().
			Foreground(lipgloss.Color(rose)).
			Render(m.err.Error()),
		)
	}

	// For displaying notifications
	if m.notification != "" {
		return paneStyle.Render(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(text)).
				Render(m.notification),
		)
	}

	tableView := m.table.View()

	centeredTable := lipgloss.NewStyle().
		Padding(0, 2).
		Render(tableView)

	return paneStyle.Render(centeredTable)
}
