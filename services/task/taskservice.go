package task

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"slices"
	"strconv"

	"github.com/pleimann/camel-do/model"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

type Config struct {
}

// TaskService is a service for managing tasks.
type TaskService struct {
	// Tasks is a slice of Task.
	tasks []model.Task

	config      *Config
	googleTasks *tasks.Service
}

func NewTaskService(config *Config, http *http.Client) (*TaskService, error) {
	service, err := tasks.NewService(context.Background(), option.WithHTTPClient(http))

	if err != nil {
		slog.Error("error creating google tasks service", "error", err.Error())
		return nil, err
	}

	taskService := &TaskService{
		config:      config,
		googleTasks: service,
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
	r, err := t.googleTasks.Tasklists.List().MaxResults(10).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve task lists. %v", err)
	}

	tasks := []model.Task{}

	fmt.Println("Task Lists:")
	if len(r.Items) > 0 {
		taskList := r.Items[0]
		gtasks, err := t.googleTasks.Tasks.List(taskList.Id).Do()
		if err != nil {
			return []model.Task{}
		}

		for i, gtask := range gtasks.Items {
			b, _ := gtask.MarshalJSON()

			fmt.Printf("%d gtask: %v\n", i, string(b))

			order, _ := strconv.Atoi(gtask.Position)

			tasks = append(tasks, model.Task{
				ID:          gtask.Id,
				Title:       gtask.Title,
				Description: gtask.Notes,
				Completed:   gtask.Completed != nil,
				Order:       order,
			})
		}

		slices.SortStableFunc(tasks, func(a, b model.Task) int {
			return cmp.Compare(a.Order, b.Order)
		})

		return tasks

	} else {
		fmt.Print("No task lists found.")
	}

	t.tasks = tasks

	return t.tasks
}
