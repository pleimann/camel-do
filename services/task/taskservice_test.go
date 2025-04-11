package task

import (
	"testing"
	"time"

	"github.com/guregu/null/v6/zero"
	"github.com/pleimann/camel-do/model"
	"github.com/stretchr/testify/assert"
)

func TestToDbTask(t *testing.T) {
	t.Run("converts task with all fields", func(t *testing.T) {
		// Arrange
		now := time.Now().UTC()

		startTime, err := time.Parse(time.RFC3339, "2025-04-07T09:30:00Z")
		if err != nil {
			t.Fatal(err)
		}

		task := &model.Task{
			ID:          "test-id",
			CreatedAt:   now,
			Title:       "Test Task",
			Description: zero.StringFrom("Test Description"),
			StartTime:   zero.TimeFrom(startTime),
			Duration:    30 * time.Minute,
			Completed:   true,
			Rank:        1.0,
			ProjectId:   zero.StringFrom("project-1"),
		}

		// Act
		dbTask, err := toDbTask(task)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, task.ID, *dbTask.ID)
		assert.Equal(t, task.CreatedAt, *dbTask.CreatedAt)
		assert.Equal(t, task.Title, *dbTask.Title)
		assert.Equal(t, task.Description, dbTask.Description)
		assert.Equal(t, task.Completed, *dbTask.Completed)
		assert.Equal(t, task.Rank, *dbTask.Rank)
		assert.Equal(t, task.ProjectId, dbTask.ProjectId)
		assert.True(t, dbTask.StartTime.Valid)
		assert.Equal(t, int32(30), *dbTask.Duration)
	})

	t.Run("handles empty start date/time", func(t *testing.T) {
		// Arrange
		task := &model.Task{
			ID:        "test-id",
			Title:     "Test Task",
			Duration:  45 * time.Minute,
			Completed: false,
			Rank:      1.0,
		}

		// Act
		dbTask, err := toDbTask(task)

		// Assert
		assert.NoError(t, err)
		assert.False(t, dbTask.StartTime.Valid)
		assert.Equal(t, int32(45), *dbTask.Duration)
	})
}
