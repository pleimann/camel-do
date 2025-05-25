package task

import (
	"testing"
	"time"

	"github.com/guregu/null/v6/zero"
	"github.com/pleimann/camel-do/model"
	"github.com/stretchr/testify/assert"
)

func TestToTableTask(t *testing.T) {
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
			Title:       zero.StringFrom("Test Task"),
			Description: zero.StringFrom("Test Description"),
			StartTime:   zero.TimeFrom(startTime),
			Duration:    zero.Int32From(30),
			Completed:   zero.BoolFrom(true),
			Rank:        zero.Int32From(1),
			ProjectID:   zero.StringFrom("project-1"),
		}

		// Act
		tableTask := toTableTask(task)

		// Assert
		assert.Equal(t, task.ID, tableTask.ID)
		assert.Equal(t, task.CreatedAt, tableTask.CreatedAt)
		assert.Equal(t, task.Title, tableTask.Title)
		assert.Equal(t, task.Description, tableTask.Description)
		assert.Equal(t, task.Completed, tableTask.Completed)
		assert.Equal(t, task.Rank, tableTask.Rank)
		assert.Equal(t, task.ProjectID, tableTask.ProjectID)
		assert.True(t, tableTask.StartTime.Valid)
		assert.Equal(t, zero.Int32From(30), tableTask.Duration)
	})

	t.Run("handles empty start date/time", func(t *testing.T) {
		// Arrange
		task := &model.Task{
			ID:        "test-id",
			Title:     zero.StringFrom("Test Task"),
			Duration:  zero.Int32From(45),
			Completed: zero.BoolFrom(false),
			Rank:      zero.Int32From(1),
		}

		// Act
		dbTask := toTableTask(task)

		// Assert
		assert.False(t, dbTask.StartTime.Valid)
		assert.Equal(t, zero.Int32From(45), dbTask.Duration)
	})
}
