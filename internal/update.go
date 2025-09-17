package internal

import (
	"time"

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

		return m, m.tickEvery()
	}

	return m, nil
}
