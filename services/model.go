package services

import (
	"time"
)

type Color int

const (
	ColorRed Color = iota
	ColorOrange
	ColorAmber
	ColorYellow
	ColorLime
	ColorGreen
	ColorEmerald
	ColorTeal
	ColorCyan
	ColorSky
	ColorBlue
	ColorIndigo
	ColorViolet
	ColorPurple
	ColorFuchsia
	ColorPink
	ColorRose
	ColorSlate
	ColorGray
	ColorZinc
	ColorNeutral
	ColorStone
)

var colorName = map[Color]string{
	ColorRed:     "red",
	ColorOrange:  "orange",
	ColorAmber:   "amber",
	ColorYellow:  "yellow",
	ColorLime:    "lime",
	ColorGreen:   "green",
	ColorEmerald: "emerald",
	ColorTeal:    "teal",
	ColorCyan:    "cyan",
	ColorSky:     "sky",
	ColorBlue:    "blue",
	ColorIndigo:  "indigo",
	ColorViolet:  "violat",
	ColorPurple:  "purple",
	ColorFuchsia: "fuchsio",
	ColorPink:    "pink",
	ColorRose:    "rose",
	ColorSlate:   "slate",
	ColorGray:    "gray",
	ColorZinc:    "zinc",
	ColorNeutral: "neutral",
	ColorStone:   "stone",
}

func (c Color) String() string {
	return colorName[c]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex function
func (c Color) EnumIndex() int {
	return int(c)
}

// Task represents a task in the task tracking application.
type Task struct {
	ID          int           `json:"id"`          // Unique identifier for the task
	Color       Color         `json:"color"`       // Color of the task
	Title       string        `json:"title"`       // Title of the task
	Description string        `json:"description"` // Description of the task
	StartTime   time.Time     `json:"startTime"`   // Start time of the task
	Duration    time.Duration `json:"duration"`    // Duration of the task
	Completed   bool          `json:"completed"`   // Status of task completion
	CreatedAt   time.Time     `json:"createdAt"`   // Timestamp indicating when the task was created.
	UpdatedAt   time.Time     `json:"updatedAt"`   // Timestamp indicating when the task was last updated.
}

// NewTask creates a new Task instance with default values.
func NewTask(title string, description string, startTime time.Time, duration time.Duration) Task {
	return Task{
		Title:       title,
		Description: description,
		StartTime:   startTime,
		Duration:    duration,
		Color:       ColorAmber,
		Completed:   false,      // default to not completed.
		CreatedAt:   time.Now(), // Set the creation timestamp
		UpdatedAt:   time.Now(), // Set the update timestamp
	}
}
