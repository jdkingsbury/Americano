package panes_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdkingsbury/americano/internal/drivers"
	"github.com/jdkingsbury/americano/internal/tui/panes"
	"github.com/stretchr/testify/assert"
)

type mockDatabase struct {
	executedQuery string
}

func (db *mockDatabase) Connect(url string) error {
	return nil
}

func (db *mockDatabase) CloseConnection() error {
	return nil
}

func (db *mockDatabase) ExecuteQuery(query string) drivers.QueryResultMsg {
	db.executedQuery = query
	return drivers.QueryResultMsg{
		Columns: []string{"id", "name"},
		Rows:    [][]string{{"1", "Fernando Tatis"}, {"2", "Manny Machado"}},
		Error:   nil,
	}
}

func (db *mockDatabase) GetDatabaseName() (string, error) {
	return "mockdb", nil
}

func (db *mockDatabase) GetTables() ([]string, error) {
	return []string{"users", "orders"}, nil
}

func TestNewEditorPane(t *testing.T) {
	db := &mockDatabase{}
	editor := panes.NewEditorPane(0, 0, db)

	query := "SELECT * FROM users;"
	editor.Update(panes.InsertQueryMsg{Query: query})

	keyMsg := tea.KeyMsg{Type: tea.KeyCtrlE}
	editor.Update(keyMsg)

	assert.Equal(t, query, db.executedQuery)
}

func TestGetTables(t *testing.T) {
	db := &mockDatabase{}

	tables, err := db.GetTables()

	assert.NoError(t, err)
	assert.Equal(t, []string{"users", "orders"}, tables)
}
