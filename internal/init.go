package internal

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return m.tickEvery()
}

func (m Model) tickEvery() tea.Cmd {
	return tea.Every(m.config.RefreshInterval, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
