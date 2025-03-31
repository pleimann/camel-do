package task

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/oklog/ulid/v2"
	. "github.com/pleimann/camel-do/db/table"

	m "github.com/pleimann/camel-do/db/model"

	"github.com/pleimann/camel-do/model"
)

type TaskServiceConfig struct {
}

// TaskService is a service for managing tasks.
type TaskService struct {
	config *TaskServiceConfig
	db     *sql.DB
}

func NewTaskService(config *TaskServiceConfig, db *sql.DB) (*TaskService, error) {
	taskService := &TaskService{
		config: config,
		db:     db,
	}

	return taskService, nil
}

func (t *TaskService) AddTask(task *model.Task) error {
	task.ID = ulid.Make().String()

	slog.Debug("TaskService.AddTask", "task", task)

	insertStmt := Tasks.INSERT(Tasks.AllColumns).
		MODEL(task)

	if res, err := insertStmt.Exec(t.db); err != nil {
		return fmt.Errorf("insert new task: %w", err)

	} else {
		rows, _ := res.RowsAffected()

		slog.Debug("TaskService.AddTask: record inserted", "count", rows)
	}

	return nil
}

func (t *TaskService) GetTask(id string) (*model.Task, error) {
	slog.Debug("TaskService.GetTask", "id", id)

	stmt := SELECT(Tasks.AllColumns).
		FROM(Tasks).
		WHERE(Tasks.ID.EQ(String(id))).
		LIMIT(1)

	var tasks []m.Tasks
	if err := stmt.Query(t.db, &tasks); err != nil {
		return nil, fmt.Errorf("TaskService.GetTask (%s): %w", id, err)
	}

	modelTask := model.ConvertTask(&tasks[0])

	return &modelTask, nil
}

func (t *TaskService) CompleteToggleTask(id string) error {
	slog.Debug("TaskService.CompleteToggleTask", "id", id)

	updateStmt := Tasks.
		UPDATE(Tasks.Completed).
		SET(CASE().
			WHEN(Tasks.Completed.IS_TRUE()).
			THEN(Bool(false)).
			ELSE(Bool(true))).
		WHERE(Tasks.ID.EQ(String(id)))

	if res, err := updateStmt.Exec(t.db); err != nil {
		return fmt.Errorf("TaskService.CompleteToggleTask (%s): %w", id, err)

	} else {
		rows, _ := res.RowsAffected()

		slog.Debug("TaskService.CompleteTask: records updated", "count", rows)
	}

	return nil
}

func (t *TaskService) UpdateTask(task *model.Task) (*model.Task, error) {
	slog.Debug("TaskService.UpdateTask", "task", task)

	updateStmt := Tasks.
		UPDATE(Tasks.MutableColumns).
		MODEL(task).
		WHERE(Tasks.ID.EQ(String(task.ID))).
		RETURNING(Tasks.AllColumns)

	updatedTasks := m.Tasks{}

	if err := updateStmt.Query(t.db, &updatedTasks); err != nil {
		return nil, fmt.Errorf("TaskService.UpdateTask(%s): %w", task.ID, err)
	}

	updatedTask := model.ConvertTask(&updatedTasks)

	return &updatedTask, nil
}

func (t *TaskService) DeleteTask(id string) error {
	slog.Debug("TaskService.DeleteTask", "id", id)

	deleteStmt := Tasks.DELETE().
		WHERE(Tasks.ID.EQ(String(id)))

	if res, err := deleteStmt.Exec(t.db); err != nil {
		return fmt.Errorf("TaskService.DeleteTask (%s): %w", id, err)

	} else {
		rows, _ := res.RowsAffected()

		slog.Debug("TaskService.DeleteTask: records deleted", "count", rows)
	}

	return nil
}

func (t *TaskService) GetBacklogTasks() ([]model.Task, error) {
	slog.Debug("TaskService.GetBacklogTasks")

	stmt := SELECT(Tasks.AllColumns).
		FROM(Tasks).
		WHERE(Tasks.StartTime.IS_NULL()).
		ORDER_BY(Tasks.UpdatedAt.DESC()).
		LIMIT(100)

	var tasks []m.Tasks
	if err := stmt.Query(t.db, &tasks); err != nil {
		return nil, fmt.Errorf("TaskService.GetBacklogTasks: %w", err)
	}

	modelTasks := model.ConvertTasks(tasks)

	return modelTasks, nil
}

func (t *TaskService) GetTodaysTasks() ([]model.Task, error) {
	slog.Debug("TaskService.GetTodaysTasks")

	year, month, day := time.Now().Date()

	start := time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())
	end := start.Add(time.Hour * 24)

	stmt := SELECT(Tasks.AllColumns.As("")).
		FROM(Tasks).
		WHERE(Tasks.StartTime.BETWEEN(DATETIME(start), DATETIME(end))).
		ORDER_BY(Tasks.UpdatedAt.DESC()).
		LIMIT(100)

	tasks := []m.Tasks{}
	if err := stmt.Query(t.db, &tasks); err != nil {
		return nil, fmt.Errorf("TaskService.GetTodaysTasks (%s - %s): %w", start, end, err)
	}

	modelTasks := model.ConvertTasks(tasks)

	return modelTasks, nil
}
