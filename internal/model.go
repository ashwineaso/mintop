package internal

import (
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

type Model struct {
	width  int
	height int

	lastUpdate   time.Time
	processTable table.Model
	tableStyle   lipgloss.Style
	baseStyle    lipgloss.Style
	viewStyle    lipgloss.Style

	HostInfo host.InfoStat
	CpuUsage cpu.TimesStat
	MemUsage mem.VirtualMemoryStat
}

type TickMsg time.Time

func NewModel() Model {
	tableStyle := table.DefaultStyles()
	tableStyle.Selected = lipgloss.NewStyle().Background(lipgloss.Color("62"))

	// Creates a new table with specified columns and initial empty rows.
	processTable := table.New(
		// We use this to define our table "header"
		table.WithColumns([]table.Column{
			{Title: "PID", Width: 10},
			{Title: "Name", Width: 25},
			{Title: "CPU", Width: 12},
			{Title: "MEM", Width: 12},
			{Title: "Username", Width: 12},
			{Title: "Time", Width: 12},
		}),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(20),
		table.WithStyles(tableStyle),
	)

	return Model{
		processTable: processTable,
		tableStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2),
		baseStyle: lipgloss.NewStyle(),
		viewStyle: lipgloss.NewStyle(),
	}
}
