package task

import (
	"fmt"
	"math/rand"
	"time"

	lorem "github.com/derektata/lorem/ipsum"
	"github.com/pleimann/camel-do/model"
)

// GenerateRandomTasks generates a slice of Task with random data.
func (t *TaskService) generateRandomTasks(count int) ([]model.Task, error) {
	if count < 1 || count > 50 {
		return nil, fmt.Errorf("task count must be between 1 and 50, got %d", count)
	}

	tasks := make([]model.Task, count)
	for i := 0; i < count; i++ {
		tasks[i] = generateRandomTask()
	}

	return tasks, nil
}

var loremGen = lorem.NewGenerator()

// generateRandomTask generates a single task with random data.
func generateRandomTask() model.Task {
	// Seed the random number generator.
	rand.New(rand.NewSource(time.Now().UnixNano()))

	color := model.Color(model.ColorNames()[rand.Intn(len(model.ColorNames())-1)])
	icon := model.Icon(model.IconNames()[rand.Intn(len(model.IconNames())-1)])

	// Generate random title.
	titles := []string{
		"Write code",
		"Read documentation",
		"Attend meeting",
		"Fix bug",
		"Test features",
		"Plan project",
		"Refactor code",
		"Deploy application",
		"Learn new technology",
		"Debug issue",
		"Review code",
	}
	title := titles[rand.Intn(len(titles))]

	// Generate random description.
	description := loremGen.Generate(rand.Intn(20) + 5)

	// Generate random start time within the past week.
	startTime := time.Now().Add(time.Duration(-rand.Intn(7*24)) * time.Hour)

	// Generate random duration between 15 minutes and 4 hours.
	duration := time.Duration(rand.Intn(4*60-15) + 15)

	// Generate random completed status.
	completed := rand.Intn(2) == 1

	createdAt := time.Now().Add(time.Duration(-rand.Intn(7*24)) * time.Hour)
	updatedAt := createdAt.Add(time.Duration(rand.Intn(72)) * time.Hour)
	if updatedAt.After(time.Now()) {
		updatedAt = time.Now()
	}

	return model.Task{
		Title:       title,
		Description: description,
		Project: model.Project{
			Color: color,
			Icon:  icon,
		},
		StartTime: startTime,
		Duration:  duration,
		Completed: completed, // default to not completed.
		CreatedAt: createdAt, // Set the creation timestamp
		UpdatedAt: updatedAt, // Set the update timestamp
	}
}
