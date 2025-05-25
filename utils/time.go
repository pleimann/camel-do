package utils

import (
	"fmt"
)

func FormatTaskDuration(d int32) string {
	hours := int(d / 60)
	minutes := d % 60

	if hours > 0 {
		return fmt.Sprintf("%2dh %2dm", hours, minutes)

	} else {
		return fmt.Sprintf("%2dm", minutes)
	}
}
