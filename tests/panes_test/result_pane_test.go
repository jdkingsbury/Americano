package panes_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/jdkingsbury/americano/internal/drivers"
	"github.com/jdkingsbury/americano/internal/tui/panes"
	"github.com/jdkingsbury/americano/msgtypes"
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

	queryMsg := drivers.QueryResultMsg{Columns: columns, Rows: rows, Error: nil}
	model, err := resultPane.HandleMsg(queryMsg)

  if err != nil {
    t.Fatalf("HandleMsg returned an error: %v", err)
  }

	resultPane, ok := model.(*panes.ResultPaneModel)
	if !ok {
		t.Fatalf("expected *panes.ResultPaneModel, got %T", model)
	}

	if len(resultPane.Table().Rows()) != len(rows) {
		t.Errorf("Expected %d rows, but got %d rows", len(rows), len(resultPane.Table().Rows()))
	}

	if len(resultPane.Table().Columns()) != len(columns) {
		t.Errorf("Expected %d columns, but got %d columns", len(columns), len(resultPane.Table().Columns()))
	}
}

func TestResultPane_DisplayError(t *testing.T) {
	resultPane := panes.NewResultPaneModel(80, 20)

	expectedError := errors.New("test error")

	errMsg := msgtypes.NewErrMsg(expectedError)
	model, _ := resultPane.HandleMsg(errMsg)
	resultPane, ok := model.(*panes.ResultPaneModel)
	if !ok {
		t.Fatalf("expected *panes.ResultPaneModel, got %T", model)
	}

	output := resultPane.View()

	if !strings.Contains(output, expectedError.Error()) {
		t.Errorf("Expected error message '%s' to be displayed, but got '%s'", expectedError.Error(), output)
	}
}

func TestResultPane_DisplayNotification(t *testing.T) {
	resultPane := panes.NewResultPaneModel(80, 20)

	expectedNotification := "Operation completed successfully"
	notificationMsg := msgtypes.NewNotificationMsg(expectedNotification)

	model, _ := resultPane.HandleMsg(notificationMsg)
	resultPane, ok := model.(*panes.ResultPaneModel)

	if !ok {
		t.Fatalf("expected *panes.ResultPaneModel, got %T", model)
	}

	output := resultPane.View()

	if !strings.Contains(output, expectedNotification) {
		t.Errorf("Expected notification message '%s' to be displayed, but got '%s'", expectedNotification, output)
	}
}
