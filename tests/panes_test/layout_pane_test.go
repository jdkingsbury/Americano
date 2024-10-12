package panes_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdkingsbury/americano/internal/tui/panes"
)

func TestLayoutModel_PaneSwitching(t *testing.T) {
	layout := panes.NewLayoutModel()

	nextPaneMsg := tea.KeyMsg{Type: tea.KeyTab}
	model, _ := layout.Update(nextPaneMsg)
	layout = model.(*panes.LayoutModel)

	if layout.CurrentPane() != panes.ResultPane {
		t.Errorf("Expected current pane to be ResultPane, got %d", layout.CurrentPane())
	}

	prevPaneMsg := tea.KeyMsg{Type: tea.KeyShiftTab}
	model, _ = layout.Update(prevPaneMsg)
	layout = model.(*panes.LayoutModel)

	if layout.CurrentPane() != panes.EditorPane {
		t.Errorf("Expected current pane to be EditorPane, got %d", layout.CurrentPane())
	}
}

func TestLayoutModel_HandleInsertQueryMsg(t *testing.T) {
	layout := panes.NewLayoutModel()
	queryMsg := panes.InsertQueryMsg{Query: "SELECT * FROM users"}

	model, _ := layout.Update(queryMsg)
	layout = model.(*panes.LayoutModel)

	editorPane := layout.Panes()[panes.EditorPane].(*panes.EditorPaneModel)
	if editorPane.Query() != "SELECT * FROM users" {
		t.Errorf("expected query to be 'SELECT * FROM users', got %s", editorPane.Query())
	}
}
