package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

func main() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	printChan := make(chan struct{})

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		printSysInfo(printChan)
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			printChan <- struct{}{}
		case <-stopChan:
			close(printChan)
			wg.Wait()
			fmt.Println("Exiting...")
			return
		}
	}
}

func printSysInfo(printChan chan struct{}) {
	for range printChan {
		fmt.Println("CPU Percentage		: ", getCPUPercentage())
		fmt.Println("Memory Percentage	: ", getMemoryPercentage())
		fmt.Println("Disk Percentage	: ", getDiskPercentage())
		fmt.Println("Running Processes	: ", getRunningProcesses())
	}
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
