package internal

import (
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
)

type Model struct {
	config Config

	width        int
	height       int
	statsFetcher StatsFetcher

	lastUpdate   time.Time
	processTable table.Model
	tableStyle   table.Styles
	baseStyle    lipgloss.Style
	viewStyle    lipgloss.Style

	HostInfo  *host.InfoStat
	CpuUsage  *cpu.TimesStat
	MemUsage  *mem.VirtualMemoryStat
	SwapUsage *mem.SwapMemoryStat
	LoadAvg   *load.AvgStat

	processManager ProcessManager
	processOptions ProcessOptions

	hasLoaded bool
}

type TickMsg time.Time

func NewModel(config Config, fetcher StatsFetcher, processManager ProcessManager) Model {
	tableStyle := table.DefaultStyles()
	tableStyle.Selected = lipgloss.NewStyle().Background(config.Colors.TableSelectionBackground)

	// Creates a new table with specified columns and initial empty rows.
	processTable := table.New(
		// We use this to define our table "header"
		table.WithColumns([]table.Column{
			{Title: "PID", Width: 6},
			{Title: "PPID", Width: 6},
			{Title: "Name", Width: 30},
			{Title: "CPU%", Width: 6},
			{Title: "MEM%", Width: 6},
			{Title: "MEM(MB)", Width: 10},
			{Title: "Username", Width: 12},
			{Title: "Time", Width: 12},
		}),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(config.ProcessTableHeight),
		table.WithStyles(tableStyle),
	)

	processOptions := ProcessOptions{
		SortBy:    SortByCPU,
		Limit:     config.ProcessLimit,
		Ascending: false,
	}

	return Model{
		config: config,

		statsFetcher: fetcher,
		processTable: processTable,
		tableStyle:   tableStyle,
		baseStyle:    lipgloss.NewStyle(),
		viewStyle:    lipgloss.NewStyle(),

		processManager: processManager,
		processOptions: processOptions,
	}
}
