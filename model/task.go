package model

import (
	"time"

	"pleimann.com/camel-do/utils"
)

// Task represents a task in the task tracking application.
type Task struct {
	ID          int            `json:"id"`          // Unique identifier for the task
	Title       string         `json:"title"`       // Title of the task
	Description string         `json:"description"` // Description of the task
	Color       Color          `json:"color"`       // Color of the task
	Icon        Icon           `json:"icon"`        // Icon to identify project
	StartTime   time.Time      `json:"startTime"`   // Start time of the task
	Duration    utils.Duration `json:"duration"`    // Duration of the task
	Completed   bool           `json:"completed"`   // Status of task completion
	CreatedAt   time.Time      `json:"createdAt"`   // Timestamp indicating when the task was created.
	UpdatedAt   time.Time      `json:"updatedAt"`   // Timestamp indicating when the task was last updated.
}

// NewTask creates a new Task instance with default values.
func NewTask(id int, title string, description string, color Color, icon Icon, startTime time.Time, duration utils.Duration, completed bool, createdAt time.Time, updatedAt time.Time) Task {
	return Task{
		ID:          id,
		Title:       title,
		Description: description,
		Color:       color,
		Icon:        icon,
		StartTime:   startTime,
		Duration:    duration,
		Completed:   completed, // default to not completed.
		CreatedAt:   createdAt, // Set the creation timestamp
		UpdatedAt:   updatedAt, // Set the update timestamp
	}
}
