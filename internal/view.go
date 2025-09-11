package internal

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

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

// ViewHeader renders the header section with CPU and Memory usage
func (m Model) ViewHeader() string {
	// defines the style for list items, including borders, border color, height, and padding.
	list := m.baseStyle.
		Border(lipgloss.NormalBorder(), false, true, false, false).
		Height(4).
		Padding(0, 1)

	// applies bold styling to the text.
	listHeader := m.baseStyle.Bold(true).Render

	slog.Debug("Rendering CPU", "cpu_stats", m.CpuUsage)

	return m.viewStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			fmt.Sprintf("Last update: %d milliseconds ago\n", time.Now().Sub(m.lastUpdate).Milliseconds()),
			lipgloss.JoinHorizontal(lipgloss.Top,
				// Progress Bars
				list.Render(
					lipgloss.JoinVertical(lipgloss.Left,
						listHeader("% Usage"),
						m.formatListItem("CPU", fmt.Sprintf("%s %.1f", ProgressBar(100-m.CpuUsage.Idle, m.baseStyle), 100-m.CpuUsage.Idle), "%"),
						m.formatListItem("MEM", fmt.Sprintf("%s %.1f", ProgressBar(m.MemUsage.UsedPercent, m.baseStyle), m.MemUsage.UsedPercent), "%"),
					),
				),
			),
		),
	)

}

func (m Model) ViewProcess() string {
	return "Process"
}

func (m Model) formatListItem(key, value string, suffix ...string) string {
	finalSuffix := ""
	if len(suffix) > 0 {
		finalSuffix = suffix[0]
	}

	listItemValue := m.baseStyle.
		Align(lipgloss.Right).
		Render(fmt.Sprintf("%s%s", value, finalSuffix))

	listItemKey := func(key string) string {
		return m.baseStyle.Render(key + ":")
	}

	return fmt.Sprintf("%s %s", listItemKey(key), listItemValue)
}

// ProgressBar creates a visual representation of a percentage as a progress bar.
func ProgressBar(percentage float64, baseStyle lipgloss.Style) string {
	totalBars := 20
	fillBars := int(percentage / 100 * float64(totalBars))
	// renders the filled part of the progress bar with a green color.
	filled := baseStyle.
		Foreground(lipgloss.Color("#aad700")).
		Render(strings.Repeat("|", fillBars))
	// renders the empty part of the progress bar with a secondary color.
	empty := baseStyle.
		Foreground(lipgloss.Color("#e7e3db")).
		Render(strings.Repeat("|", totalBars-fillBars))

	return baseStyle.Render(fmt.Sprintf("%s%s%s%s", "[", filled, empty, "]"))
}
