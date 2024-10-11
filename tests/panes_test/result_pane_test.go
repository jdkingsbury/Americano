package panes_test

import (
	"testing"

	"github.com/jdkingsbury/americano/internal/tui/panes"
)

func TestResultPane_UpdateTable(t *testing.T) {
	resultPane := panes.NewResultPaneModel(80, 20)

	columns := []string{"ID", "Name", "Age", "Occupation", "Country"}
	rows := [][]string{
		{"1", "Alice", "29", "Engineer", "USA"},
		{"2", "Bob", "34", "Designer", "UK"},
		{"3", "Charlie", "22", "Student", "Canada"},
		{"4", "David", "40", "Manager", "Australia"},
		{"5", "Eve", "35", "Scientist", "Germany"},
	}

	resultPane.UpdateTable(columns, rows)

	if len(resultPane.Table().Rows()) != len(rows) {
		t.Errorf("Expected %d rows, but got %d rows", len(rows), len(resultPane.Table().Rows()))
	}

	if len(resultPane.Table().Columns()) != len(columns) {
		t.Errorf("Expected %d columns, but got %d columns", len(columns), len(resultPane.Table().Columns()))
	}
}
