package task

import (
	"database/sql"
	"fmt"
	"log/slog"
	"slices"
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

	tableTask := toTableTask(task)

	insertStmt := Tasks.INSERT(Tasks.AllColumns).
		MODEL(tableTask)

	if _, err := insertStmt.Exec(t.db); err != nil {
		return fmt.Errorf("insert new task: %w", err)
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

	modelTask := toModelTask(&tasks[0])

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

func (t *TaskService) UpdateTask(task *model.Task) error {
	slog.Debug("TaskService.UpdateTask", "task", task)

	columns := detectUpdatedColumns(task)

	tableTask := toTableTask(task)

	updateStmt := Tasks.
		UPDATE(columns).
		MODEL(tableTask).
		WHERE(Tasks.ID.EQ(String(task.ID)))

	if _, err := updateStmt.Exec(t.db); err != nil {
		return fmt.Errorf("TaskService.UpdateTask(%s): %w", task.ID, err)

	} else {
		return nil
	}
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
		ORDER_BY(Tasks.Rank, Tasks.UpdatedAt).
		LIMIT(100)

	var tasks []m.Tasks
	if err := stmt.Query(t.db, &tasks); err != nil {
		return nil, fmt.Errorf("TaskService.GetBacklogTasks: %w", err)
	}

	modelTasks := toModelTasks(tasks)

	sortTasksByStartTime(&modelTasks)

	return modelTasks, nil
}

func (t *TaskService) GetTodaysTasks() ([]model.Task, error) {
	slog.Debug("TaskService.GetTodaysTasks")

	year, month, day := time.Now().Date()

	start := time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())
	end := start.Add(time.Hour * 24)

	stmt := SELECT(Tasks.AllColumns).
		FROM(Tasks).
		WHERE(Tasks.StartTime.BETWEEN(DATETIME(start), DATETIME(end))).
		ORDER_BY(Tasks.UpdatedAt.DESC()).
		LIMIT(100)

	var tasks []m.Tasks
	if err := stmt.Query(t.db, &tasks); err != nil {
		return nil, fmt.Errorf("TaskService.GetTodaysTasks (%s - %s): %w", start, end, err)
	}

	modelTasks := toModelTasks(tasks)

	sortTasksByStartTime(&modelTasks)

	return modelTasks, nil
}

func sortTasksByStartTime(tasks *[]model.Task) {
	slices.SortFunc(*tasks, func(a, b model.Task) int {
		if a.StartTime.Valid && b.StartTime.Valid {
			timeCmp := a.StartTime.Time.Compare(b.StartTime.Time)

			if timeCmp == 0 {
				if a.Rank.Int32 < b.Rank.Int32 {
					return -1
				} else if a.Rank.Int32 > b.Rank.Int32 {
					return 1
				} else {
					if a.Duration.Int32 < b.Duration.Int32 {
						return -1
					} else if a.Duration.Int32 > b.Duration.Int32 {
						return 1
					} else {
						return 0
					}
				}

			} else {
				return timeCmp
			}
		}

		return 0
	})
}

func toTableTask(task *model.Task) m.Tasks {
	tasks := m.Tasks{
		ID:          task.ID,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		Title:       task.Title,
		Description: task.Description,
		StartTime:   task.StartTime,
		Duration:    task.Duration,
		Completed:   task.Completed,
		Rank:        task.Rank,
		ProjectID:   task.ProjectID,
		GTaskID:     task.GTaskID,
	}

	return tasks
}

func toModelTasks(tasks []m.Tasks) []model.Task {
	modelTasks := make([]model.Task, len(tasks))
	for i, t := range tasks {
		modelTasks[i] = toModelTask(&t)
	}

	return modelTasks
}

func toModelTask(t *m.Tasks) model.Task {
	id := t.ID

	task := model.NewTask(
		id,
		t.Title,
		t.Description,
		t.CreatedAt,
		t.UpdatedAt,
		t.StartTime,
		t.Duration,
		t.Completed,
		t.Rank,
		t.ProjectID,
		t.GTaskID,
	)

	return task
}

func detectUpdatedColumns(task *model.Task) ColumnList {
	columns := ColumnList{}

	if !task.CreatedAt.IsZero() {
		columns = append(columns, Tasks.CreatedAt)
	}

	if !task.UpdatedAt.IsZero() {
		columns = append(columns, Tasks.UpdatedAt)
	}

	if !task.Title.IsZero() {
		columns = append(columns, Tasks.Title)
	}

	if !task.Description.IsZero() {
		columns = append(columns, Tasks.Description)
	}

	if !task.StartTime.IsZero() {
		columns = append(columns, Tasks.StartTime)
	}

	if !task.Duration.IsZero() {
		columns = append(columns, Tasks.Duration)
	}

	if !task.Completed.IsZero() {
		columns = append(columns, Tasks.Completed)
	}

	if !task.GTaskID.IsZero() {
		columns = append(columns, Tasks.GTaskID)
	}

	if !task.Rank.IsZero() {
		columns = append(columns, Tasks.Rank)
	}

	if !task.ProjectID.IsZero() {
		columns = append(columns, Tasks.ProjectID)
	}

	return columns
}
