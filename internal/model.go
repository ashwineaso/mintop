package internal

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type Model struct {
	width  int
	height int

	processTable table.Model
	tableStyle   lipgloss.Style
	baseStyle    lipgloss.Style
	viewStyle    lipgloss.Style

	CpuUsage cpu.TimesStat
	MemUsage mem.VirtualMemoryStat
}

func NewModel() Model {
	return Model{
		processTable: table.New(
			table.WithColumns([]table.Column{
				{Title: "PID", Width: 10},
				{Title: "User", Width: 20},
				{Title: "CPU%", Width: 10},
				{Title: "MEM%", Width: 10},
				{Title: "Command", Width: 50},
			}),
			table.WithRows([]table.Row{}),
			table.WithFocused(true),
		),
		tableStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2),
		baseStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")),
		viewStyle: lipgloss.NewStyle().
			Margin(1, 2),
		CpuUsage: cpu.TimesStat{},
		MemUsage: mem.VirtualMemoryStat{},
	}
}
