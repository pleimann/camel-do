package model

import (
	"time"

	"github.com/guregu/null/v6/zero"

	m "github.com/pleimann/camel-do/db/model"
)

// Task represents a task in the task tracking application.
type Task struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time

	Title       string        `schema:"title,required"`          // Title of the task
	Description zero.String   `schema:"description"`             // Description of the task
	StartTime   zero.Time     `schema:"startTime"`               // Start time of the task
	Duration    time.Duration `schema:"duration,default:0s"`     // Duration of the task
	Completed   bool          `schema:"completed,default:false"` // Status of task completion
	GTaskId     string
	Rank        int32       // Sort order
	ProjectID   zero.String `schema:"projectId"` // Foreign key referencing the project associated with the task.
}

func ConvertTasks(tasks []m.Tasks) []Task {
	modelTasks := make([]Task, len(tasks))
	for i, t := range tasks {
		modelTasks[i] = ConvertTask(&t)
	}

	return modelTasks
}

func ConvertTask(t *m.Tasks) Task {
	id := *t.ID
	duration := time.Duration(*t.Duration).Round(time.Minute)

	task := Task{
		ID:          id,
		CreatedAt:   *t.CreatedAt,
		UpdatedAt:   *t.UpdatedAt,
		Title:       *t.Title,
		Description: t.Description,
		StartTime:   t.StartTime,
		Duration:    duration,
		Completed:   *t.Completed,
		Rank:        *t.Rank,
		ProjectID:   t.ProjectId,
	}

	return task
}
