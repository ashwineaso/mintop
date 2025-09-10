package internal

import "github.com/charmbracelet/lipgloss"

func (m Model) View() string {
	column := m.baseStyle.Width(m.width).Padding(1, 0, 0, 0).Render

	content := m.baseStyle.
		Width(m.width).
		Height(m.height).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				column(m.ViewHeader()),
				column(m.ViewProcess()),
			),
		)

	return content
}

func (m Model) ViewHeader() string {
	return m.CpuUsage.String() + "\n" + m.MemUsage.String()
}

func (m Model) ViewProcess() string {
	return "Process"
}
