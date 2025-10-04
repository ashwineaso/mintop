package internal

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

type ColorConfig struct {
	TableSelectionBackground lipgloss.Color
	ProgressBarFilled        lipgloss.Color
	ProgressBarEmpty         lipgloss.Color
}

type Config struct {
	RefreshInterval    time.Duration
	ProcessLimit       int
	Colors             ColorConfig
	ProcessTableHeight int
}

func DefaultConfig() *Config {
	return &Config{
		RefreshInterval: time.Second,
		ProcessLimit:    25,
		Colors: ColorConfig{
			TableSelectionBackground: lipgloss.Color("62"),
			ProgressBarFilled:        lipgloss.Color("#aad700"),
			ProgressBarEmpty:         lipgloss.Color("#e7e3db"),
		},
		ProcessTableHeight: 25,
	}
}

func (c *Config) WithRefreshInterval(d time.Duration) Config {
	c.RefreshInterval = d
	return *c
}

func (c *Config) WithProcessLimit(limit int) Config {
	c.ProcessLimit = limit
	return *c
}

func (c *Config) WithProcessTableHeight(height int) Config {
	c.ProcessTableHeight = height
	return *c
}
