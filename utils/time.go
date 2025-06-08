package utils

import (
	"fmt"
	"time"
)

func FormatDuration(d int32) string {
	duration := time.Duration(d) * time.Minute

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%2dh %2dm", hours, minutes)

	} else {
		return fmt.Sprintf("%2dm", minutes)
	}
}

func FormatTime(t time.Time) string {
	return t.Local().Format("03:04 PM")
}
