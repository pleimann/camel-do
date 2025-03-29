package task

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/google/uuid"
	. "github.com/pleimann/camel-do/.gen/table"

	m "github.com/pleimann/camel-do/.gen/model"

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
	task.ID = uuid.New()

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

func (t *TaskService) GetTask(id uuid.UUID) (*model.Task, error) {
	slog.Debug("TaskService.GetTask", "id", id)

	stmt := SELECT(Tasks.AllColumns).
		FROM(Tasks).
		WHERE(Tasks.ID.EQ(UUID(id))).
		LIMIT(1)

	var tasks []m.Tasks
	if err := stmt.Query(t.db, &tasks); err != nil {
		return nil, fmt.Errorf("get tasks: %w", err)
	}

	modelTask := model.ConvertTask(&tasks[0])

	return &modelTask, nil
}

func (t *TaskService) CompleteTask(id uuid.UUID, completed bool) error {
	slog.Debug("completing task", "id", id, "completed", completed)

	updateStmt := Tasks.UPDATE(Tasks.Completed).
		SET(Tasks.Completed.SET(Bool(true))).
		WHERE(Tasks.ID.EQ(UUID(id)))

	if res, err := updateStmt.Exec(t.db); err != nil {
		return fmt.Errorf("complete task (%s): %w", id, err)

	} else {
		rows, _ := res.RowsAffected()

		slog.Debug("TaskService.CompleteTask: records updated", "count", rows)
	}

	return nil
}

func (t *TaskService) UpdateTask(task *model.Task) (*model.Task, error) {
	slog.Debug("updating task", "task", task)

	updateStmt := Tasks.
		UPDATE(Tasks.MutableColumns).
		MODEL(task).
		WHERE(Tasks.ID.EQ(UUID(task.ID))).
		RETURNING(Tasks.AllColumns)

	var updatedTasks *m.Tasks

	if err := updateStmt.Query(t.db, updatedTasks); err != nil {
		return nil, err
	}

	updatedTask := model.ConvertTask(updatedTasks)

	return &updatedTask, nil
}

func (t *TaskService) DeleteTask(id uuid.UUID) error {
	slog.Debug("deleting task", "id", id)

	deleteStmt := Tasks.DELETE().
		WHERE(Tasks.ID.EQ(UUID(id)))

	if res, err := deleteStmt.Exec(t.db); err != nil {
		return fmt.Errorf("delete task (%s): %w", id, err)

	} else {
		rows, _ := res.RowsAffected()

		slog.Debug("TaskService.DeleteTask: records deleted", "count", rows)
	}

	return nil
}

func (t *TaskService) GetBacklogTasks() ([]model.Task, error) {
	stmt := SELECT(Tasks.AllColumns).
		FROM(Tasks).
		WHERE(Tasks.StartTime.IS_NULL()).
		ORDER_BY(Tasks.UpdatedAt.DESC()).
		LIMIT(100)

	var tasks []m.Tasks
	if err := stmt.Query(t.db, &tasks); err != nil {
		return nil, fmt.Errorf("get backlog tasks: %w", err)
	}

	modelTasks := model.ConvertTasks(tasks)

	return modelTasks, nil
}

func (t *TaskService) GetTodaysTasks() ([]model.Task, error) {
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
		return nil, fmt.Errorf("get todays tasks (%s - %s): %w", start, end, err)
	}

	modelTasks := model.ConvertTasks(tasks)

	return modelTasks, nil
}
