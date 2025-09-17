package internal

import (
	"log/slog"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
)

// StatsFetcher defines the interface for fetching system statistics.
// This allows us to use a real implementation in production and a mock for testing.
type StatsFetcher interface {
	HostInfo() (*host.InfoStat, error)
	CpuUsage() (*cpu.TimesStat, error)
	MemUsage() (*mem.VirtualMemoryStat, error)
	SwapUsage() (*mem.SwapMemoryStat, error)
	LoadAvg() (*load.AvgStat, error)
}

// LiveStatsFetcher is the production implementation of StatsFetcher that uses gopsutil.
type LiveStatsFetcher struct{}

func (l LiveStatsFetcher) HostInfo() (*host.InfoStat, error) {
	info, err := host.Info()
	if err != nil {
		return &host.InfoStat{}, err
	}

	return info, nil
}

func (l LiveStatsFetcher) CpuUsage() (*cpu.TimesStat, error) {
	cpuTimes, err := cpu.Times(false)
	if err != nil || len(cpuTimes) == 0 {
		slog.Error("Failed to get CPU stats", "error", err)
		return &cpu.TimesStat{}, err
	}

	currStats := cpuTimes[0]

	// Calculate total time
	total := currStats.User + currStats.System + currStats.Idle + currStats.Nice +
		currStats.Iowait + currStats.Irq + currStats.Softirq + currStats.Steal +
		currStats.Guest

	if total == 0 {
		return &cpu.TimesStat{}, nil
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

	return &currStats, nil
}

func (l LiveStatsFetcher) MemUsage() (*mem.VirtualMemoryStat, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return &mem.VirtualMemoryStat{}, err
	}

	return &mem.VirtualMemoryStat{
		Total:       v.Total,
		Used:        v.Used,
		Free:        v.Free,
		UsedPercent: v.UsedPercent,
		Available:   v.Available,
	}, nil
}

func (l LiveStatsFetcher) SwapUsage() (*mem.SwapMemoryStat, error) {
	return mem.SwapMemory()
}

func (l LiveStatsFetcher) LoadAvg() (*load.AvgStat, error) {
	return load.Avg()
}
