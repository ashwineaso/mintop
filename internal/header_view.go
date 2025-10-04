package internal

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// HeaderView handles rendering of the header section with system stats.
type HeaderView struct {
	baseStyle lipgloss.Style
	viewStyle lipgloss.Style
}

// NewHeaderView creates a new HeaderView instance.
func NewHeaderView(config Config, baseStyle, viewStyle lipgloss.Style) *HeaderView {
	return &HeaderView{
		baseStyle: baseStyle,
		viewStyle: viewStyle,
	}
}

// Render renders the complete header section with CPU, Memory, and Load stats.
func (h *HeaderView) Render(m Model) string {
	// If there is no data to load, then return Loading text
	if !m.hasLoaded {
		return h.viewStyle.Render(lipgloss.JoinVertical(lipgloss.Top, "Loading..."))
	}

	return h.viewStyle.Render(
		lipgloss.JoinVertical(lipgloss.Top,
			h.renderHostDetails(m),
			h.renderStatsSection(m),
		),
	)
}

// renderHostDetails renders the host information line.
func (h *HeaderView) renderHostDetails(m Model) string {
	hostDetails := h.baseStyle.
		Height(1).
		Padding(1, 1).Render

	return hostDetails(fmt.Sprintf("Host: %s | OS: %s | Arch: %s | Uptime: %s",
		m.HostInfo.Hostname, m.HostInfo.OS, m.HostInfo.KernelArch, timeToHuman(m.HostInfo.Uptime)))
}

// renderStatsSection renders all the stats columns (Usage, CPU, Memory, Load Avg).
func (h *HeaderView) renderStatsSection(m Model) string {
	return lipgloss.JoinHorizontal(lipgloss.Top,
		h.renderUsageColumn(m),
		h.renderCPUColumn(m),
		h.renderMemoryColumn(m),
		h.renderLoadAvgColumn(m),
	)
}

// renderUsageColumn renders the usage progress bars column.
func (h *HeaderView) renderUsageColumn(m Model) string {
	list := h.createListStyle()
	listHeader := h.baseStyle.Bold(true).Render

	return list.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			listHeader("% Usage"),
			listItem(h.baseStyle, "CPU", fmt.Sprintf("%s %.1f", ProgressBar(100-m.CpuUsage.Idle, h.baseStyle, m.config.Colors.ProgressBarFilled, m.config.Colors.ProgressBarEmpty), 100-m.CpuUsage.Idle), "%"),
			listItem(h.baseStyle, "MEM", fmt.Sprintf("%s %.1f", ProgressBar(m.MemUsage.UsedPercent, h.baseStyle, m.config.Colors.ProgressBarFilled, m.config.Colors.ProgressBarEmpty), m.MemUsage.UsedPercent), "%"),
			listItem(h.baseStyle, "SWAP", fmt.Sprintf("%s %.1f", ProgressBar(m.SwapUsage.UsedPercent, h.baseStyle, m.config.Colors.ProgressBarFilled, m.config.Colors.ProgressBarEmpty), m.SwapUsage.UsedPercent), "%"),
		),
	)
}

// renderCPUColumn renders the CPU stats column.
func (h *HeaderView) renderCPUColumn(m Model) string {
	list := h.createListStyle().Border(lipgloss.NormalBorder(), false, true, false, false)
	listHeader := h.baseStyle.Bold(true).Render

	return list.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			listHeader("CPU"),
			listItem(h.baseStyle, "User", fmt.Sprintf("%5.2f", m.CpuUsage.User), "%"),
			listItem(h.baseStyle, "Sys ", fmt.Sprintf("%5.2f", m.CpuUsage.System), "%"),
			listItem(h.baseStyle, "Idle", fmt.Sprintf("%5.2f", m.CpuUsage.Idle), "%"),
		),
	)
}

// renderMemoryColumn renders the Memory stats column.
func (h *HeaderView) renderMemoryColumn(m Model) string {
	list := h.createListStyle().Border(lipgloss.NormalBorder(), false, true, false, false)
	listHeader := h.baseStyle.Bold(true).Render

	return list.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			listHeader("Memory"),
			h.formatMemoryItem("total", m.MemUsage.Total),
			h.formatMemoryItem("used", m.MemUsage.Used),
			h.formatMemoryItem("free", m.MemUsage.Available),
		),
	)
}

// renderLoadAvgColumn renders the Load Average column.
func (h *HeaderView) renderLoadAvgColumn(m Model) string {
	list := h.createListStyle().Border(lipgloss.NormalBorder(), false, true, false, false)
	listHeader := h.baseStyle.Bold(true).Render

	return list.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			listHeader("Load Avg"),
			listItem(h.baseStyle, "1 min", fmt.Sprintf("%5.2f", m.LoadAvg.Load1), ""),
			listItem(h.baseStyle, "5 min", fmt.Sprintf("%5.2f", m.LoadAvg.Load5), ""),
			listItem(h.baseStyle, "15 min", fmt.Sprintf("%5.2f", m.LoadAvg.Load15), ""),
		),
	)
}

// createListStyle creates the base style for list containers.
func (h *HeaderView) createListStyle() lipgloss.Style {
	return h.baseStyle.
		Border(lipgloss.NormalBorder(), false, true, false, false).
		Height(4).
		Padding(0, 1)
}

// formatMemoryItem formats a memory value with appropriate units.
func (h *HeaderView) formatMemoryItem(label string, bytes uint64) string {
	value, unit := convertBytes(bytes)
	return listItem(h.baseStyle, label, value, unit)
}
