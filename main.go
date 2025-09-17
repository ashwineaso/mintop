package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ashwineaso/mintop/internal"
)

func main() {
	logFile, err := os.OpenFile("mintop.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		os.Exit(1)
	}
	defer logFile.Close()

	slogHandler := slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(slogHandler)
	slog.SetDefault(logger)

	// Define and parse the refresh interval flag
	refreshInterval := flag.Duration("refresh", time.Second, "Set the refresh interval for system stats")
	flag.Parse()

	fetcher := internal.LiveStatsFetcher{}

	p := tea.NewProgram(internal.NewModel(*refreshInterval, fetcher), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
