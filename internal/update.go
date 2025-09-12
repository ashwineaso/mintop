package internal

import (
	"fmt"
	"log/slog"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

func (m Model) Update(teaMsg tea.Msg) (tea.Model, tea.Cmd) {
	var err error

	switch msg := teaMsg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case TickMsg:
		m.lastUpdate = time.Time(msg)
		slog.Debug("Tick at", "time", m.lastUpdate)

		// Update Host info every minute
		m.HostInfo, err = GetHostInfo()
		if err != nil {
			slog.Error("Failed to get Host info", "error", err)
		}

		// Update CPU and Memory stats every second
		m.CpuUsage, err = GetCPUStats()
		if err != nil {
			// handle error appropriately, e.g., log it or set a default value
			slog.Error("Failed to get CPU stats", "error", err)
		}

		m.MemUsage, err = GetMemStats()
		if err != nil {
			// handle error appropriately, e.g., log it or set a default value
			slog.Error("Failed to get Memory stats", "error", err)
		}

		return m, tickEvery()
	}

	return m, nil
}

func GetHostInfo() (host.InfoStat, error) {
	info, err := host.Info()
	if err != nil {
		return host.InfoStat{}, err
	}

	return *info, nil
}

func GetCPUStats() (cpu.TimesStat, error) {
	cpuTimes, err := cpu.Times(false)
	if err != nil || len(cpuTimes) == 0 {
		slog.Error("Failed to get CPU stats", "error", err)
		return cpu.TimesStat{}, err
	}

	currStats := cpuTimes[0]

	// Calculate total time
	total := currStats.User + currStats.System + currStats.Idle + currStats.Nice +
		currStats.Iowait + currStats.Irq + currStats.Softirq + currStats.Steal +
		currStats.Guest

	if total == 0 {
		return cpu.TimesStat{}, nil
	}

	// Overwrite TimesStat fields with percentage values
	currStats.User = (currStats.User / total) * 100
	currStats.System = (currStats.System / total) * 100
	currStats.Idle = (currStats.Idle / total) * 100
	currStats.Nice = (currStats.Nice / total) * 100
	currStats.Iowait = (currStats.Iowait / total) * 100
	currStats.Irq = (currStats.Irq / total) * 100
	currStats.Softirq = (currStats.Softirq / total) * 100
	currStats.Steal = (currStats.Steal / total) * 100
	currStats.Guest = (currStats.Guest / total) * 100

	return currStats, nil
}

func GetMemStats() (mem.VirtualMemoryStat, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return mem.VirtualMemoryStat{}, err
	}

	return mem.VirtualMemoryStat{
		Total:       v.Total,
		Used:        v.Used,
		Free:        v.Free,
		UsedPercent: v.UsedPercent,
		Available:   v.Available,
	}, nil
}

func convertBytes(bytes uint64) (string, string) {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%5.2f", float64(bytes)/float64(GB)), "GB"
	case bytes >= MB:
		return fmt.Sprintf("%5.2f", float64(bytes)/float64(MB)), "MB"
	case bytes >= KB:
		return fmt.Sprintf("%5.2f", float64(bytes)/float64(KB)), "KB"
	default:
		return fmt.Sprintf("%d", bytes), "B"
	}
}
