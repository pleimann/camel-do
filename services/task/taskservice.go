package task

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	lorem "github.com/derektata/lorem/ipsum"
	"github.com/pleimann/camel-do/model"
	"github.com/pleimann/camel-do/utils"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

type Config struct {
}

// TaskService is a service for managing tasks.
type TaskService struct {
	// Tasks is a slice of Task.
	tasks []model.Task

	config       *Config
	tasksService *tasks.Service
}

func NewTaskService(config *Config, http *http.Client) (*TaskService, error) {
	service, err := tasks.NewService(context.Background(), option.WithHTTPClient(http))

	if err != nil {
		slog.Error("error creating google tasks service", "error", err.Error())
		return nil, err
	}

	taskService := &TaskService{
		config:       config,
		tasksService: service,
	}

	return taskService, nil
}

func (t *TaskService) AddTask(task model.Task) (model.Task, error) {
	slog.Info("adding task", "task", task)

	slog.Info("NewTask", "color", task.Color, "icon", task.Icon)

	if task.Color == "" {
		task.Color = model.ColorZinc
	}

	if task.Icon == "" {
		task.Icon = model.IconCircleHelp
	}

	slog.Info("NewTask", "color", task.Color, "icon", task.Icon)

	t.tasks = append([]model.Task{task}, t.tasks...)

	return task, nil
}

func (t *TaskService) GetTasks() []model.Task {
	r, err := t.tasksService.Tasklists.List().MaxResults(10).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve task lists. %v", err)
	}

	fmt.Println("Task Lists:")
	if len(r.Items) > 0 {
		for _, i := range r.Items {
			fmt.Printf("%s (%s)\n", i.Title, i.Id)
		}
	} else {
		fmt.Print("No task lists found.")
	}

	tasks, err := t.generateRandomTasks(rand.Intn(50) + 1)

	if err != nil {
		return []model.Task{}
	}

	t.tasks = tasks

	return t.tasks
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
