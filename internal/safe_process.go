package internal

import (
	"fmt"
	"log/slog"

	"github.com/shirou/gopsutil/v4/process"
)

func safeProcessString(f func() (string, error)) string {
	defer func() {
		if r := recover(); r != nil {
			slog.Error(fmt.Sprintf("Recovered from panic: %v", r))
		}
	}()
	val, err := f()
	if err != nil {
		return "Unknown"
	}
	return val
}

func safeProcessFloat64(f func() (float64, error)) float64 {
	defer func() {
		if r := recover(); r != nil {
			slog.Error(fmt.Sprintf("Recovered from panic: %v", r))
		}
	}()
	val, err := f()
	if err != nil {
		return 0.0
	}
	return val
}

func safeProcessFloat32(f func() (float32, error)) float32 {
	defer func() {
		if r := recover(); r != nil {
			slog.Error(fmt.Sprintf("Recovered from panic: %v", r))
		}
	}()
	val, err := f()
	if err != nil {
		return 0.0
	}
	return val
}

func safeProcessInt32(f func() (int32, error)) int32 {
	defer func() {
		if r := recover(); r != nil {
			slog.Error(fmt.Sprintf("Recovered from panic: %v", r))
		}
	}()
	val, err := f()
	if err != nil {
		return 0
	}
	return val
}

func safeProcessInt64(f func() (int64, error)) int64 {
	defer func() {
		if r := recover(); r != nil {
			slog.Error(fmt.Sprintf("Recovered from panic: %v", r))
		}
	}()
	val, err := f()
	if err != nil {
		return 0
	}
	return val
}

func safeMemoryInfo(p *process.Process) *process.MemoryInfoStat {
	defer func() {
		if r := recover(); r != nil {
			slog.Error(fmt.Sprintf("Recovered from panic: %v", r))
		}
	}()
	memInfo, err := p.MemoryInfo()
	if err != nil {
		return nil
	}
	return memInfo
}
