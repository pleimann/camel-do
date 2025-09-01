package task

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/guregu/null/v6/zero"
	"github.com/oklog/ulid/v2"
	"github.com/pleimann/camel-do/model"
	"github.com/pleimann/camel-do/utils"
	bolt "go.etcd.io/bbolt"
)

type TaskServiceConfig struct {
}

// TaskService is a service for managing tasks.
type TaskService struct {
	config *TaskServiceConfig
	db     *bolt.DB
}

func NewTaskService(config *TaskServiceConfig, db *bolt.DB) (*TaskService, error) {
	taskService := &TaskService{
		config: config,
		db:     db,
	}

	return taskService, nil
}

func (t *TaskService) AddTask(task *model.Task) error {
	task.ID = ulid.Make().String()

	slog.Debug("TaskService.AddTask", "task", task)

	err := t.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("tasks"))
		if err != nil {
			return err
		}

		if bucket == nil {
			return fmt.Errorf("tasks bucket does not exist")
		}

		taskBytes, err := task.Marshal()

		if err != nil {
			return err
		}

		bucket.Put([]byte(task.ID), taskBytes)

		return nil
	})

	if err != nil {
		return fmt.Errorf("adding task %s %w", task.Title.String, err)
	}

	return nil
}

func (t *TaskService) GetTask(id string) (*model.Task, error) {
	slog.Debug("TaskService.GetTask", "id", id)

	task := model.Task{}
	err := t.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("tasks"))

		if bucket == nil {
			return fmt.Errorf("tasks bucket does not exist")
		}

		taskBytes := bucket.Get([]byte(id))

		if taskBytes == nil {
			return utils.NewNotFoundError("task", id)
		}

		if err := task.Unmarshal(taskBytes); err != nil {
			return fmt.Errorf(" %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("TaskService.GetTask (%s): %w", id, err)
	}

	return &task, nil
}

func (t *TaskService) CompleteToggleTask(id string) error {
	slog.Debug("TaskService.CompleteToggleTask", "id", id)

	err := t.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("tasks"))

		if bucket == nil {
			return fmt.Errorf("tasks bucket does not exist")
		}

		taskBytes := bucket.Get([]byte(id))

		if taskBytes == nil {
			return utils.NewNotFoundError("task", id)
		}

		task := model.Task{}

		if err := task.Unmarshal(taskBytes); err != nil {
			return err
		}

		task.Completed.SetValid(!task.Completed.ValueOr(false))

		if taskBytes, err := task.Marshal(); err != nil {
			return err

		} else {
			bucket.Put([]byte(id), taskBytes)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("TaskService.CompleteToggleTask (%s): %w", id, err)
	}

	return nil
}

func (t *TaskService) HiddenToggleTask(id string) error {
	slog.Debug("TaskService.CompleteToggleTask", "id", id)

	err := t.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("tasks"))

		if bucket == nil {
			return fmt.Errorf("tasks bucket does not exist")
		}

		taskBytes := bucket.Get([]byte(id))

		if taskBytes == nil {
			return utils.NewNotFoundError("task", id)
		}

		task := model.Task{}

		if err := task.Unmarshal(taskBytes); err != nil {
			return err
		}

		task.Hidden.SetValid(!task.Hidden.ValueOr(false))

		if taskBytes, err := task.Marshal(); err != nil {
			return err

		} else {
			bucket.Put([]byte(id), taskBytes)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("TaskService.HiddenToggleTask (%s): %w", id, err)
	}

	return nil
}

func (t *TaskService) UpdateTask(task *model.Task) error {
	slog.Debug("TaskService.UpdateTask", "task", task)

	err := t.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("tasks"))

		if bucket == nil {
			return fmt.Errorf("tasks bucket does not exist")
		}

		taskBytes, err := task.Marshal()

		if err != nil {
			return err
		}

		bucket.Put([]byte(task.ID), taskBytes)

		return nil
	})

	if err != nil {
		return fmt.Errorf("adding task %s %w", task.Title.String, err)
	}

	return nil
}

func (t *TaskService) ScheduleTask(id string, time zero.Time) error {
	slog.Debug("TaskService.ScheduleTask", "taskId", id)

	err := t.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("tasks"))

		taskBytes := bucket.Get([]byte(id))

		if taskBytes == nil {
			return utils.NewNotFoundError("task", id)
		}

		task := model.Task{}

		if err := task.Unmarshal(taskBytes); err != nil {
			return err
		}

		task.StartTime = time

		if taskBytes, err := task.Marshal(); err != nil {
			return err

		} else {
			bucket.Put([]byte(id), taskBytes)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("TaskService.ScheduleTask(%s): %w", id, err)
	}

	return nil
}

func (t *TaskService) DeleteTask(id string) error {
	slog.Debug("TaskService.DeleteTask", "id", id)

	err := t.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("tasks"))

		taskBytes := bucket.Get([]byte(id))
		if taskBytes == nil {
			return utils.NewNotFoundError("task", id)
		}

		if err := bucket.Delete([]byte(id)); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("TaskService.DeleteTask (%s): %w", id, err)
	}

	return nil
}

func (t *TaskService) GetBacklogTasks() (*model.TaskList, error) {
	slog.Debug("TaskService.GetBacklogTasks")

	taskList := model.NewTaskList()

	err := t.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("tasks"))

		if bucket == nil {
			return nil
		}

		err := bucket.ForEach(func(taskID, taskBytes []byte) error {
			task := model.Task{}

			if err := task.Unmarshal(taskBytes); err != nil {
				return err
			}

			if task.StartTime.IsZero() {
				taskList.Push(task)
			}

			return nil
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("TaskService.GetBacklogTasks: %w", err)
	}

	taskList.Sort()

	return taskList, nil
}

func (t *TaskService) GetTodaysTasks() (*model.TaskList, error) {
	slog.Debug("TaskService.GetTodaysTasks")

	year, month, day := time.Now().Date()

	date := time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())

	return t.GetTasksScheduledOnDate(date)
}

func (t *TaskService) GetTasksScheduledOnDate(date time.Time) (*model.TaskList, error) {
	slog.Debug("TaskService.GetTasksScheduledOnDate")

	year, month, day := date.Date()

	beginningOfDay := time.Date(year, month, day, 0, 0, 0, 0, time.Local)

	endOfDay := beginningOfDay.Add(time.Hour * 24)

	taskList := model.NewTaskList()

	slog.Debug("finding tasks between", "start", beginningOfDay, "end", endOfDay)

	err := t.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("tasks"))

		if bucket == nil {
			return nil
		}

		err := bucket.ForEach(func(taskID, taskBytes []byte) error {
			task := model.Task{}

			if err := task.Unmarshal(taskBytes); err != nil {
				return err
			}

			if task.StartTime.Valid &&
				task.StartTime.Time.Equal(beginningOfDay) ||
				(task.StartTime.Time.After(beginningOfDay) &&
					task.StartTime.Time.Before(endOfDay)) {
				slog.Debug("task", "id", task.ID, "startTime", task.StartTime, "start", beginningOfDay, "end", endOfDay)
				taskList.Push(task)
			}

			return nil
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("TaskService.GetTasksScheduledOnDate (%s - %s): %w", beginningOfDay, endOfDay, err)
	}

	slog.Debug("found tasks", "count", taskList.Len(), "start", beginningOfDay, "end", endOfDay)

	return taskList, nil
}
