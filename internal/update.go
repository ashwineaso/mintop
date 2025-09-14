package internal

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
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
		case "up", "k":
			m.processTable.MoveUp(1)
			return m, nil
		case "down", "j":
			m.processTable.MoveDown(1)
			return m, nil
		}

	// Handle the TickMsg to update system stats
	case TickMsg:
		m.lastUpdate = time.Time(msg)

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

		m.SwapUsage, err = GetSwapMemStats()
		if err != nil {
			slog.Error("Failed to get Swap Memory stats", "error", err)
		}

		m.Load1, m.Load5, m.Load15, err = GetLoadAvg()
		if err != nil {
			slog.Error("Failed to get Load Average", "error", err)
		}

		processes, err := GetProcess()
		if err != nil {
			slog.Error("Failed to get process info", "error", err)
		} else {
			var rows []table.Row
			for _, p := range processes {
				rows = append(rows, table.Row{
					fmt.Sprintf("%d", p.PID),
					fmt.Sprintf("%d", p.ParentPID),
					p.Name,
					fmt.Sprintf("%.2f%%", p.CPUPercent),
					fmt.Sprintf("%.2f%%", p.MemoryPercent),
					fmt.Sprintf("%.2fMB", p.MemoryUsage),
					p.Username,
					p.RunningTime,
				})
			}

			m.processTable.SetRows(rows)
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

func GetSwapMemStats() (mem.SwapMemoryStat, error) {
	s, err := mem.SwapMemory()
	if err != nil {
		return mem.SwapMemoryStat{}, err
	}

	return mem.SwapMemoryStat{
		Total:       s.Total,
		Used:        s.Used,
		Free:        s.Free,
		UsedPercent: s.UsedPercent,
	}, nil
}

func GetLoadAvg() (float64, float64, float64, error) {
	avg, err := load.Avg()
	if err != nil {
		return 0, 0, 0, err
	}

	return avg.Load1, avg.Load5, avg.Load15, nil
}

// convertBytes converts bytes to a human-readable format (B, KB, MB, GB)
func convertBytes(bytes uint64) (string, string) {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%6.2f", float64(bytes)/float64(GB)), "GB"
	case bytes >= MB:
		return fmt.Sprintf("%6.2f", float64(bytes)/float64(MB)), "MB"
	case bytes >= KB:
		return fmt.Sprintf("%6.2f", float64(bytes)/float64(KB)), "KB"
	default:
		return fmt.Sprintf("%d", bytes), "B"
	}
}
