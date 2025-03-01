package services

import "time"

// Task represents a task in the task tracking application.
type Task struct {
	ID          int           `json:"id"`          // Unique identifier for the task
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
		Completed:   false,      // default to not completed.
		CreatedAt:   time.Now(), // Set the creation timestamp
		UpdatedAt:   time.Now(), // Set the update timestamp
	}
}
