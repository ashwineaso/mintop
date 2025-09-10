package internal

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return tea.Every(time.Second,
		func(t time.Time) tea.Msg {
			return time.Second
		})
}
