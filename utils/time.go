package utils

import (
	"fmt"
	"time"
)

func FormatTaskDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%2dh %2dm", hours, minutes)

	} else {
		return fmt.Sprintf("%2dm", minutes)
	}
}
