package model

import (
	"time"

	"github.com/guregu/null/v6/zero"
)

// Task represents a task in the task tracking application.
type Task struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time

	Title       string      `schema:"title"`                   // Title of the task
	Description zero.String `schema:"description"`             // Description of the task
	StartTime   zero.Time   `schema:"startTime"`               // Start time of the task
	Duration    zero.Int32  `schema:"duration"`                // Duration of the task
	Completed   bool        `schema:"completed,default:false"` // Status of task completion
	GTaskId     string
	Rank        zero.Int32  // Sort order
	ProjectID   zero.String `schema:"projectId"` // Foreign key referencing the project associated with the task.
}
