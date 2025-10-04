package internal

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// ProcessView handles rendering of the process table.
type ProcessView struct {
	viewStyle lipgloss.Style
}

// NewProcessView creates a new ProcessView instance.
func NewProcessView(viewStyle lipgloss.Style) *ProcessView {
	return &ProcessView{
		viewStyle: viewStyle,
	}
}

// Render renders the process table.
func (p *ProcessView) Render(processTable table.Model) string {
	return p.viewStyle.Render(processTable.View())
}
