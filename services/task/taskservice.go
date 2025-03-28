package task

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/pleimann/camel-do/model"
	"github.com/pleimann/camel-do/services/db"
	"gorm.io/gorm"
)

type Config struct {
}

// TaskService is a service for managing tasks.
type TaskService struct {
	config *Config
	db     *db.DatabaseService
}

func NewTaskService(config *Config, db *db.DatabaseService) (*TaskService, error) {
	taskService := &TaskService{
		config: config,
		db:     db,
	}

	return taskService, nil
}

func (t *TaskService) AddTask(task *model.Task) error {
	slog.Debug("adding task", "task", task)

	if task.ID == "" {
		task.ID = ulid.Make().String()
	}

	if err := t.db.Create(task).Error; err != nil {
		return err
	}

	return nil
}

func (t *TaskService) GetTask(id string) (*model.Task, error) {
	slog.Debug("getting task", "id", id)

	task := model.Task{
		ID: id,
	}

	if err := t.db.Preload("Project").First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("not found: %w", err)
		} else {
			return nil, err
		}

	}

	return &task, nil
}

func (t *TaskService) CompleteTask(id string, completed bool) (*model.Task, error) {
	slog.Debug("completing task", "id", id, "completed", completed)

	if task, err := t.GetTask(id); err != nil {
		return nil, err

	} else {
		task.Completed = completed

		if err := t.db.Save(&task).Error; err != nil {
			return nil, err
		}

		return task, nil
	}
}

func (t *TaskService) UpsertTask(task *model.Task) error {
	slog.Debug("updating task", "task", task)

	if err := t.db.Save(task).Error; err != nil {
		return err
	}

	return nil
}

func (t *TaskService) DeleteTask(id string) error {
	slog.Debug("deleting task", "id", id)
	if err := t.db.Delete(&model.Task{}, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("not found: %w", err)
		}

		return err
	}

	return nil
}

func (t *TaskService) GetBacklogTasks() ([]model.Task, error) {
	tasks := []model.Task{}
	if err := t.db.Limit(100).Order("updated_at desc").Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *TaskService) GetTodaysTasks() ([]model.Task, error) {
	tasks := []model.Task{}
	end := time.Now().UTC()
	start := end.Add(-time.Hour * 24)
	if err := t.db.Where("start_time BETWEEN ? AND ?", start, end).Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}
