package main

import (
	"fmt"
	"strconv"

	"github.com/shirou/gopsutil/v4/cpu"
)

func main() {
	fmt.Println("CPU Percentage		: ", getCPUPercentage())
	fmt.Println("Memory Percentage	: ", getMemoryPercentage())
	fmt.Println("Disk Percentage	: ", getDiskPercentage())
	fmt.Println("Running Processes	: ", getRunningProcesses())
}

func getCPUPercentage() string {
	cpuTimes, err := cpu.Times(false)
	if err != nil || len(cpuTimes) == 0 {
		return "Error fetching CPU times"
	}

	c := cpuTimes[0]
	user := strconv.FormatFloat(c.User, 'f', 1, 64)
	system := strconv.FormatFloat(c.System, 'f', 1, 64)
	idle := strconv.FormatFloat(c.Idle, 'f', 1, 64)
	return fmt.Sprintf("User: %s, System: %s, Idle: %s", user, system, idle)
}

func getMemoryPercentage() float64 {
	return 0.0
}

func getDiskPercentage() float64 {
	return 0.0
}

func getRunningProcesses() int {
	return 0
}
