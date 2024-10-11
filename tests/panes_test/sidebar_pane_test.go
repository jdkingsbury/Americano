package panes_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdkingsbury/americano/internal/tui/panes"
	"github.com/stretchr/testify/assert"
)

func TestSideBarPane_SwitchView(t *testing.T) {
	sidebar := panes.NewSideBarPane(80, 20)

	assert.Equal(t, panes.ConnectionsView, sidebar.CurrentView())

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("v")}
	sidebarModel, _ := sidebar.Update(keyMsg)
	sidebar = sidebarModel.(*panes.SideBarPaneModel)

	assert.Equal(t, panes.DBTreeView, sidebar.CurrentView())
}

func TestSideBarPane_SelectItem(t *testing.T) {
	sidebar := panes.NewSideBarPane(80, 20)

	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	sidebarModel, _ := sidebar.Update(keyMsg)
	sidebar = sidebarModel.(*panes.SideBarPaneModel)

	assert.True(t, sidebar.ShowInputForm())
}
