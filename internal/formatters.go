package internal

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ProgressBar creates a visual representation of a percentage as a progress bar.
func ProgressBar(percentage float64, baseStyle lipgloss.Style, filledColor, emptyColor lipgloss.Color) string {
	totalBars := 25
	fillBars := int(percentage / 100 * float64(totalBars))

	// renders the filled part of the progress bar with a green color.
	filled := baseStyle.
		Foreground(lipgloss.Color(filledColor)).
		Render(strings.Repeat("|", fillBars))

	// renders the empty part of the progress bar with a secondary color.
	empty := baseStyle.
		Foreground(lipgloss.Color(emptyColor)).
		Render(strings.Repeat("|", totalBars-fillBars))

	return baseStyle.Render(fmt.Sprintf("%s%s%s%s", "[", filled, empty, "]"))
}

// timeToHuman converts seconds to a human-readable format.
func timeToHuman(seconds uint64) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	days := hours / 24
	hours = hours % 24

	if days > 0 {
		return fmt.Sprintf("%d days, %02d hrs, %02d mins", days, hours, minutes)
	}
	return fmt.Sprintf("%02d hrs, %02d mins", hours, minutes)
}

// listItem formats a key-value pair with an optional suffix.
// It aligns the value to the right and renders it with the specified style.
func listItem(baseStyle lipgloss.Style, key string, value string, suffix ...string) string {
	finalSuffix := ""
	if len(suffix) > 0 {
		finalSuffix = suffix[0]
	}

	listItemValue := baseStyle.Align(lipgloss.Right).Render(fmt.Sprintf("%s %s", value, finalSuffix))
	listItemKey := func(key string) string {
		return baseStyle.Render(key + ":")
	}

	return fmt.Sprintf("%6s %s", listItemKey(key), listItemValue)
}
