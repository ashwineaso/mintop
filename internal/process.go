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
