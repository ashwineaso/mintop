package internal

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var err error

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

	return m, nil
}

func GetCPUStats() (cpu.TimesStat, error) {
	cpuTimes, err := cpu.Times(false)
	if err != nil || len(cpuTimes) == 0 {
		return cpu.TimesStat{}, err
	}

	return cpuTimes[0], nil
}

func GetMemStats() (mem.VirtualMemoryStat, error) {
	memStats, err := mem.VirtualMemory()
	if err != nil {
		return mem.VirtualMemoryStat{}, err
	}

	return *memStats, nil
}
