package services

import (
	"fmt"
	"math/rand"
	"time"
)

// TaskService is a service for managing tasks.
type TaskService struct {
	// Tasks is a slice of Task.
	Tasks []Task
}

func (t *TaskService) GetTasks() []Task {
	if t.Tasks == nil {
		tasks, err := t.generateRandomTasks(rand.Intn(50) + 1)

		if err != nil {
			return []Task{}
		}

		t.Tasks = tasks
	}

	return t.Tasks
}

// GenerateRandomTasks generates a slice of Task with random data.
func (t *TaskService) generateRandomTasks(count int) ([]Task, error) {
	if count < 1 || count > 50 {
		return nil, fmt.Errorf("task count must be between 1 and 50, got %d", count)
	}

	tasks := make([]Task, count)
	for i := 0; i < count; i++ {
		tasks[i] = t.generateRandomTask(i)
	}

	return tasks, nil
}

// generateRandomTask generates a single task with random data.
func (t *TaskService) generateRandomTask(id int) Task {
	// Seed the random number generator.
	rand.New(rand.NewSource(time.Now().UnixNano()))

	color := Color(rand.Intn(len(colorName)))

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
	descriptions := []string{
		"Implement feature X",
		"Fix issue Y in module Z",
		"Write tests for the new functionality",
		"Update documentation",
		"Discuss requirements with team",
		"Prepare presentation",
		"Optimize performance",
		"Refactor for better readability",
		"Investigate unexpected behavior",
		"Setup new environment",
		"Research possible solutions",
	}
	description := descriptions[rand.Intn(len(descriptions))]

	// Generate random start time within the past week.
	startTime := time.Now().Add(time.Duration(-rand.Intn(7*24)) * time.Hour)

	// Generate random duration between 15 minutes and 4 hours.
	duration := time.Duration(rand.Intn(4*60-15)+15) * time.Minute

	// Generate random completed status.
	completed := rand.Intn(2) == 1

	createdAt := time.Now().Add(time.Duration(-rand.Intn(7*24)) * time.Hour)
	updatedAt := createdAt.Add(time.Duration(rand.Intn(72)) * time.Hour)
	if updatedAt.After(time.Now()) {
		updatedAt = time.Now()
	}

	return Task{
		ID:          id,
		Color:       color,
		Title:       title,
		Description: description,
		StartTime:   startTime,
		Duration:    duration,
		Completed:   completed,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
