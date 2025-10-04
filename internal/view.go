package internal

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	headerView := NewHeaderView(m.config, m.baseStyle, m.viewStyle)
	processView := NewProcessView(m.viewStyle)

	column := m.baseStyle.Width(m.width).Padding(1, 0, 0, 0).Render

	content := m.baseStyle.
		Width(m.width).
		Height(m.height).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				column(headerView.Render(m)),
				column(processView.Render(m.processTable)),
			),
		)
	return content
}
