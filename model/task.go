package model

import (
	"time"
)

// Task represents a task in the task tracking application.
type Task struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Title       string        `schema:"title,required"`                               // Title of the task
	Description string        `schema:"description"`                                  // Description of the task
	StartTime   time.Time     `schema:"startTime"`                                    // Start time of the task
	Duration    time.Duration `schema:"duration,default:0s" gorm:"default:0"`         // Duration of the task
	Completed   bool          `schema:"completed,default:false" grom:"default:false"` // Status of task completion
	Rank        int           // Sort order
	Project     Project       `gorm:"foreignKey:ID"` // Foreign key referencing the project associated with the task.
}
