package model

import (
	"time"

	"github.com/google/uuid"

	m "github.com/pleimann/camel-do/.gen/model"
)

// Task represents a task in the task tracking application.
type Task struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Title       string        `schema:"title,required" db:"title"`               // Title of the task
	Description string        `schema:"description" db:"description"`            // Description of the task
	StartTime   *time.Time    `schema:"startTime" db:"start_time"`               // Start time of the task
	Duration    time.Duration `schema:"duration,default:0s" sql:"duration"`      // Duration of the task
	Completed   bool          `schema:"completed,default:false" sql:"completed"` // Status of task completion
	GTaskId     string        `db:"gtask_id"`
	Rank        int32         `db:"rank"`                          // Sort order
	ProjectID   uuid.UUID     `schema:"projectId" db:"project_id"` // Foreign key referencing the project associated with the task.
}

func ConvertTasks(tasks []m.Tasks) []Task {
	modelTasks := make([]Task, len(tasks))
	for i, t := range tasks {
		modelTasks[i] = ConvertTask(&t)
	}

	return modelTasks
}

func ConvertTask(t *m.Tasks) Task {
	id, _ := uuid.Parse(*t.ID)
	duration := time.Duration(*t.Duration)

	var startTime *time.Time = nil
	if t.StartTime != nil && t.StartTime.Valid {
		startTime = &t.StartTime.Time
	}

	task := Task{
		ID:          id,
		CreatedAt:   *t.CreatedAt,
		UpdatedAt:   *t.UpdatedAt,
		Title:       *t.Title,
		Description: *t.Description,
		StartTime:   startTime,
		Duration:    duration,
		Completed:   *t.Completed,
		Rank:        *t.Rank,
		ProjectID:   *t.ProjectId,
	}

	return task
}
