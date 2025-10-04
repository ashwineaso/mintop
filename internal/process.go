package internal

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
