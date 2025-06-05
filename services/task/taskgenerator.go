package task

import (
	"fmt"
	"math"
	"math/rand/v2"
	"time"

	lorem "github.com/derektata/lorem/ipsum"
	"github.com/guregu/null/v6/zero"
	"github.com/pleimann/camel-do/model"
)

// GenerateRandomTasks generates a slice of Task with random data.
func GenerateRandomTasks(count int) ([]model.Task, error) {
	if count < 1 || count > 50 {
		return nil, fmt.Errorf("task count must be between 1 and 50, got %d", count)
	}

	tasks := make([]model.Task, count)
	for i := 0; i < count; i++ {
		tasks[i] = GenerateRandomTask()
	}

	return tasks, nil
}

var loremGen = lorem.NewGenerator()

// generateRandomTask generates a single task with random data.
func GenerateRandomTask() model.Task {
	// Seed the random number generator.
	rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))

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
	title := titles[rand.IntN(len(titles))]

	// Generate random description.
	description := zero.StringFrom(loremGen.Generate(rand.IntN(20) + 5))

	// Generate random start time within the past week.
	var startTime zero.Time
	if rand.IntN(5) == 1 {
		startTime = zero.TimeFrom(time.Now().Add(time.Duration(-rand.IntN(7*24)) * time.Hour))
	}

	// Generate random duration between 15 minutes and 4 hours.
	duration := int32(math.Round(rand.Float64()*4.0) * 15)

	// Generate random completed status.
	completed := rand.IntN(3) == 1

	createdAt := time.Now().Add(time.Duration(-rand.IntN(7*24)) * time.Hour)
	updatedAt := createdAt.Add(time.Duration(rand.IntN(72)) * time.Hour)
	if updatedAt.After(time.Now()) {
		updatedAt = time.Now()
	}

	return model.Task{
		Title:       zero.StringFrom(title),
		Description: description,
		ProjectID:   zero.StringFromPtr(nil),
		StartTime:   startTime,
		Duration:    zero.Int32From(duration),
		Completed:   zero.BoolFrom(completed), // default to not completed.
		CreatedAt:   createdAt,                // Set the creation timestamp
		UpdatedAt:   updatedAt,                // Set the update timestamp
	}
}
