package panes

import (
	"github.com/charmbracelet/bubbles/textarea"
)

func NewModel() *PaneModel {
	ti := textarea.New()
	ti.Placeholder = "Enter SQL Code Here..."
	ti.Focus()
	ti.CharLimit = 1000
	ti.ShowLineNumbers = false

	return &PaneModel{
		textarea: ti,
		err:      nil,
	}
}

func (m *PaneModel) resizeTextArea() {
	m.textarea.SetWidth(m.styles.MainPane.GetWidth())
	m.textarea.SetHeight(m.styles.MainPane.GetHeight())
}
