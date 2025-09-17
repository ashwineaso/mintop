package internal

import "fmt"

// convertBytes converts bytes to a human-readable format (B, KB, MB, GB)
func convertBytes(bytes uint64) (string, string) {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%6.2f", float64(bytes)/float64(GB)), "GB"
	case bytes >= MB:
		return fmt.Sprintf("%6.2f", float64(bytes)/float64(MB)), "MB"
	case bytes >= KB:
		return fmt.Sprintf("%6.2f", float64(bytes)/float64(KB)), "KB"
	default:
		return fmt.Sprintf("%d", bytes), "B"
	}
}
