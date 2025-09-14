package internal

import (
	"sort"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

type ProcessInfo struct {
	PID           int32
	ParentPID     int32
	Name          string
	Username      string
	CPUPercent    float64
	MemoryPercent float32
	MemoryUsage   float64
	RunningTime   string
}

func GetProcess() ([]ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var processInfos []ProcessInfo
	for _, p := range procs {
		pid := p.Pid
		parentPid, err := p.Ppid()
		if err != nil {
			parentPid = 0
		}

		name, err := p.Name()
		if err != nil {
			name = "Unknown Process"
		}

		username, err := p.Username()
		if err != nil {
			username = "Unknown User"
		}

		cpuPercent, err := p.CPUPercent()
		if err != nil {
			cpuPercent = 0.0
		}

		memoryPercent, err := p.MemoryPercent()
		if err != nil {
			memoryPercent = 0.0
		}

		memoryUsage := 0.0
		memoryInfo, err := p.MemoryInfo()
		if memoryInfo != nil {
			memoryUsage = float64(memoryInfo.RSS) / (1024 * 1024) // Convert bytes to MB
		}

		createTime, err := p.CreateTime()
		if err != nil {
			createTime = 0
		}

		runningTime := time.Since(time.Unix(0, createTime*int64(time.Millisecond))).Truncate(time.Second).String()
		if createTime == 0 {
			runningTime = "Unknown"
		}

		processInfos = append(processInfos, ProcessInfo{
			PID:           pid,
			ParentPID:     parentPid,
			Name:          name,
			Username:      username,
			CPUPercent:    cpuPercent,
			MemoryPercent: memoryPercent,
			MemoryUsage:   memoryUsage,
			RunningTime:   runningTime,
		})
	}

	// Sort processes by CPU usage in descending order
	// and limit to top 20 processes
	sort.Slice(processInfos, func(i, j int) bool {
		return processInfos[i].CPUPercent > processInfos[j].CPUPercent
	})
	if len(processInfos) > 20 {
		processInfos = processInfos[:20]
	}

	return processInfos, nil
}
