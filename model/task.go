package model

import (
	"time"

	"gorm.io/gorm"
)

// Task represents a task in the task tracking application.
type Task struct {
	gorm.Model

	ID          string        `schema:"id" gorm:"primaryKey,autoIncrement"`           // Unique identifier for the task
	Title       string        `schema:"title,required"`                               // Title of the task
	Description string        `schema:"description"`                                  // Description of the task
	Color       Color         `schema:"color" gorm:"default:Zinc"`                    // Color of the task
	Icon        Icon          `schema:"-" gorm:"default:unknown"`                     // Icon to identify project
	StartTime   time.Time     `schema:"startTime"`                                    // Start time of the task
	Duration    time.Duration `schema:"duration"`                                     // Duration of the task
	Completed   bool          `schema:"completed,default:false" grom:"default:false"` // Status of task completion
	CreatedAt   time.Time     // Timestamp indicating when the task was created.
	UpdatedAt   time.Time     // Timestamp indicating when the task was last updated.
	Order       int
}

// NewTask creates a new Task instance with default values.
func NewTask(id string, title string, description string, color Color, icon Icon, startTime time.Time, duration time.Duration, completed bool, createdAt time.Time, updatedAt time.Time, order int) Task {
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
		Order:       order,
	}
}
