package task

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	lorem "github.com/derektata/lorem/ipsum"
	"pleimann.com/camel-do/model"
	"pleimann.com/camel-do/services/google"
	"pleimann.com/camel-do/services/google/oauth"
	"pleimann.com/camel-do/utils"
)

// TaskService is a service for managing tasks.
type TaskService struct {
	// Tasks is a slice of Task.
	Tasks []model.Task

	config      *Config
	googleTasks *google.GoogleTasksService
}

type Config struct {
	TokenSourceProvider *oauth.TokenSourceProvider
}

func NewTaskService(config *Config) *TaskService {
	googleTasks, err := google.NewGoogleTasksService(
		context.Background(),
		config.TokenSourceProvider,
	)

	if err != nil {
		log.Fatal("Unable to create Google Tasks service: ", err)
	}

	return &TaskService{
		googleTasks: googleTasks,
		config:      config,
	}
}

func (t *TaskService) GetTasks() []model.Task {
	tasks, err := t.generateRandomTasks(rand.Intn(50) + 1)

	if err != nil {
		return []model.Task{}
	}

	t.Tasks = tasks

	return t.Tasks
}

func (t *TaskService) CreateTask(task model.Task) (model.Task, error) {
	t.Tasks = append(t.Tasks, task)

	return task, nil
}

// GenerateRandomTasks generates a slice of Task with random data.
func (t *TaskService) generateRandomTasks(count int) ([]model.Task, error) {
	if count < 1 || count > 50 {
		return nil, fmt.Errorf("task count must be between 1 and 50, got %d", count)
	}

	tasks := make([]model.Task, count)
	for i := 0; i < count; i++ {
		tasks[i] = t.generateRandomTask(i)
	}

	return tasks, nil
}

var loremGen = lorem.NewGenerator()

// generateRandomTask generates a single task with random data.
func (t *TaskService) generateRandomTask(id int) model.Task {
	// Seed the random number generator.
	rand.New(rand.NewSource(time.Now().UnixNano()))

	color := model.Color(model.ColorNames()[rand.Intn(len(model.ColorNames()))])
	icon := model.Icon(model.IconNames()[rand.Intn(len(model.IconNames()))])

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
	duration := utils.Duration{Duration: time.Duration(rand.Intn(4*60-15)+15) * time.Minute}

	// Generate random completed status.
	completed := rand.Intn(2) == 1

	createdAt := time.Now().Add(time.Duration(-rand.Intn(7*24)) * time.Hour)
	updatedAt := createdAt.Add(time.Duration(rand.Intn(72)) * time.Hour)
	if updatedAt.After(time.Now()) {
		updatedAt = time.Now()
	}

	return model.NewTask(id, title, description, color, icon, startTime, duration, completed, createdAt, updatedAt)
}
