package panes_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdkingsbury/americano/internal/drivers"
	"github.com/jdkingsbury/americano/internal/tui/panes"
	"github.com/jdkingsbury/americano/tests"
	"github.com/stretchr/testify/assert"
)

func TestEditorPane_ExecuteQuery(t *testing.T) {
	mockDB := &tests.MockDatabase{
		QueryResult: drivers.QueryResultMsg{
			Columns: []string{"id", "name"},
			Rows:    [][]string{{"1", "test"}},
			Error:   nil,
		},
	}

	editor := panes.NewEditorPane(80, 20, mockDB)

	// Insert query into editor pane
	var cmd tea.Cmd
	editorModel, cmd := editor.Update(panes.InsertQueryMsg{Query: "SELECT * FROM users;"})
	editor = editorModel.(*panes.EditorPaneModel)

	// Simulate key presses
	keyMsg := tea.KeyMsg{Type: tea.KeyCtrlE}
	editorModel, cmd = editor.Update(keyMsg)
	editor = editorModel.(*panes.EditorPaneModel)

	// Execute the command (which should run the query)
	var msg tea.Msg
	if cmd != nil {
		msg = cmd() // Capture the result message from cmd
	}

	// Type assert the result message to drivers.QueryResultMsg
	queryResult, ok := msg.(drivers.QueryResultMsg)
	assert.True(t, ok, "expected QueryResultMsg, got something else")
	assert.Nil(t, queryResult.Error, "expected no error in query result")
	assert.Equal(t, []string{"id", "name"}, queryResult.Columns, "unexpected columns")
	assert.Equal(t, [][]string{{"1", "test"}}, queryResult.Rows, "unexpected rows")

	// Check if the query is executed as expected
	assert.Equal(t, "SELECT * FROM users;", mockDB.ExecutedQuery)
}
