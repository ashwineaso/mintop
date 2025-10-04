package internal

import (
	"sort"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

type SortCriteria string

const (
	SortByCPU    SortCriteria = "cpu"
	SortByMemory SortCriteria = "memory"
	SortByPID    SortCriteria = "pid"
	SortByName   SortCriteria = "name"
)

// ProcessOptions represents options for fetching processes.
type ProcessOptions struct {
	SortBy    SortCriteria
	Limit     int
	Ascending bool
}

// ProcessManager defines the interface for fetching and managing processes.
type ProcessManager interface {
	GetProcesses(opts ProcessOptions) ([]ProcessInfo, error)
}

type DefaultProcessManager struct{}

func NewProcessManager() ProcessManager {
	return &DefaultProcessManager{}
}

func (m *DefaultProcessManager) GetProcesses(opts ProcessOptions) ([]ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var processInfos []ProcessInfo
	for _, p := range procs {

		parentPid := safeProcessInt32(p.Ppid)
		name := safeProcessString(p.Name)
		username := safeProcessString(p.Username)
		cpuPercent := safeProcessFloat64(p.CPUPercent)
		memoryPercent := safeProcessFloat32(p.MemoryPercent)
		createTime := safeProcessInt64(p.CreateTime)

		memoryUsage := 0.0
		memoryInfo := safeMemoryInfo(p)
		if memoryInfo != nil {
			memoryUsage = float64(memoryInfo.RSS) / (1024 * 1024) // Convert bytes to MB
		}

		runningTime := "Unknown"
		if createTime > 0 {
			runningTime = time.Since(time.Unix(0, createTime*int64(time.Millisecond))).Truncate(time.Second).String()
		}

		processInfos = append(processInfos, ProcessInfo{
			PID:           p.Pid,
			ParentPID:     parentPid,
			Name:          name,
			Username:      username,
			CPUPercent:    cpuPercent,
			MemoryPercent: memoryPercent,
			MemoryUsage:   memoryUsage,
			RunningTime:   runningTime,
		})
	}

	// Sort the process based on the SortBy Criteria in the options
	// and the opts.Ascending to determine the sort direction
	switch opts.SortBy {
	case SortByCPU:
		sort.Slice(processInfos, func(i, j int) bool {
			if opts.Ascending {
				return processInfos[i].CPUPercent < processInfos[j].CPUPercent
			}
			return processInfos[i].CPUPercent > processInfos[j].CPUPercent
		})
	case SortByName:
		sort.Slice(processInfos, func(i, j int) bool {
			if opts.Ascending {
				return processInfos[i].Name < processInfos[j].Name
			}
			return processInfos[i].Name > processInfos[j].Name
		})
	case SortByMemory:
		sort.Slice(processInfos, func(i, j int) bool {
			if opts.Ascending {
				return processInfos[i].MemoryUsage < processInfos[j].MemoryUsage
			}
			return processInfos[i].MemoryUsage > processInfos[j].MemoryUsage
		})
	case SortByPID:
		sort.Slice(processInfos, func(i, j int) bool {
			if opts.Ascending {
				return processInfos[i].PID < processInfos[j].PID
			}
			return processInfos[i].PID > processInfos[j].PID
		})
	}

	if len(processInfos) > opts.Limit {
		processInfos = processInfos[:opts.Limit]
	}

	return processInfos, nil
}

func DefaultProcessOptions() ProcessOptions {
	return ProcessOptions{
		SortBy:    SortByCPU,
		Limit:     25,
		Ascending: false,
	}
}
