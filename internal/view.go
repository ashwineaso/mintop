package internal

import (
	"fmt"
	"strings"

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

	hostDetails := m.baseStyle.
		Height(1).
		Padding(1, 1).Render

	// applies bold styling to the text.
	listHeader := m.baseStyle.Bold(true).Padding(0, 1).Render

	// helper function that formats a key-value pair with an optional suffix. It aligns the value to the right and renders it with the specified style.
	listItem := func(key string, value string, suffix ...string) string {
		finalSuffix := ""
		if len(suffix) > 0 {
			finalSuffix = suffix[0]
		}

		listItemValue := m.baseStyle.Align(lipgloss.Right).Render(fmt.Sprintf("%s %s", value, finalSuffix))

		listItemKey := func(key string) string {
			return m.baseStyle.Render(key + ":")
		}

		return fmt.Sprintf("%6s %s", listItemKey(key), listItemValue)
	}
	return m.viewStyle.Render(
		lipgloss.JoinVertical(lipgloss.Top,
			hostDetails(fmt.Sprintf("Host: %s | OS: %s | Arch: %s | Uptime: %s",
				m.HostInfo.Hostname, m.HostInfo.OS, m.HostInfo.KernelArch, timeToHuman(m.HostInfo.Uptime))),
			lipgloss.JoinHorizontal(lipgloss.Top,
				// Progress Bars
				list.Render(
					lipgloss.JoinVertical(lipgloss.Left,
						listHeader("% Usage"),
						listItem("CPU", fmt.Sprintf("%s %.1f", ProgressBar(100-m.CpuUsage.Idle, m.baseStyle), 100-m.CpuUsage.Idle), "%"),
						listItem("MEM", fmt.Sprintf("%s %.1f", ProgressBar(m.MemUsage.UsedPercent, m.baseStyle), m.MemUsage.UsedPercent), "%"),
						listItem("SWAP", fmt.Sprintf("%s %.1f", ProgressBar(m.SwapUsage.UsedPercent, m.baseStyle), m.SwapUsage.UsedPercent), "%"),
					),
				),

				// CPU
				list.Border(lipgloss.NormalBorder(), false, true, false, false).Render(
					lipgloss.JoinVertical(lipgloss.Left,
						listHeader("CPU"),
						listItem("User", fmt.Sprintf("%5.2f", m.CpuUsage.User), "%"),
						listItem("Sys", fmt.Sprintf("%5.2f", m.CpuUsage.System), "%"),
						listItem("Idle", fmt.Sprintf("%5.2f", m.CpuUsage.Idle), "%"),
					),
				),

				// MEM
				list.Border(lipgloss.NormalBorder(), false, true, false, false).Render(
					lipgloss.JoinVertical(lipgloss.Left,
						listHeader("Memory"),
						func() string {
							value, unit := convertBytes(m.MemUsage.Total)
							return listItem("total", value, unit)
						}(),
						func() string {
							value, unit := convertBytes(m.MemUsage.Used)
							return listItem("used", value, unit)
						}(),
						func() string {
							value, unit := convertBytes(m.MemUsage.Available)
							return listItem("free", value, unit)
						}(),
					),
				),

				// SWAP MEM
				list.Border(lipgloss.NormalBorder(), false, true, false, false).Render(
					lipgloss.JoinVertical(lipgloss.Left,
						listHeader("Swap"),
						func() string {
							value, unit := convertBytes(m.SwapUsage.Total)
							return listItem("total", value, unit)
						}(),
						func() string {
							value, unit := convertBytes(m.SwapUsage.Used)
							return listItem("used", value, unit)
						}(),
						func() string {
							value, unit := convertBytes(m.SwapUsage.Free)
							return listItem("free", value, unit)
						}(),
					),
				),
			),
		),
	)

}

func (m Model) ViewProcess() string {
	return "Process"
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

func timeToHuman(seconds uint64) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	days := hours / 24
	hours = hours % 24

	if days > 0 {
		return fmt.Sprintf("%d days, %02d hrs, %02d mins", days, hours, minutes)
	}

	return fmt.Sprintf("%02d hrs, %02d mins", hours, minutes)
}
