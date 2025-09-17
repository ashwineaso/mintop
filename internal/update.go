package internal

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(teaMsg tea.Msg) (tea.Model, tea.Cmd) {
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
		m = m.updateStats()
		m.lastUpdate = time.Time(msg)
		m.hasLoaded = true

		return m, m.tickEvery()
	}

	return m, nil
}

func (m Model) updateStats() Model {
	var err error
	// Update Host info every minute
	m.HostInfo, err = m.statsFetcher.HostInfo()
	if err != nil {
		slog.Error("Failed to get Host info", "error", err)
	}

	// Update CPU and Memory stats every second
	m.CpuUsage, err = m.statsFetcher.CpuUsage()
	if err != nil {
		// handle error appropriately, e.g., log it or set a default value
		slog.Error("Failed to get CPU stats", "error", err)
	}

	m.MemUsage, err = m.statsFetcher.MemUsage()
	if err != nil {
		// handle error appropriately, e.g., log it or set a default value
		slog.Error("Failed to get Memory stats", "error", err)
	}

	m.SwapUsage, err = m.statsFetcher.SwapUsage()
	if err != nil {
		slog.Error("Failed to get Swap Memory stats", "error", err)
	}

	m.LoadAvg, err = m.statsFetcher.LoadAvg()
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

	return m
}
