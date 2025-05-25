package task

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/angelofallars/htmx-go"
	"github.com/gorilla/mux"
	"github.com/guregu/null/v6/zero"
	"github.com/labstack/echo/v4"

	"github.com/pleimann/camel-do/model"
	"github.com/pleimann/camel-do/services/project"
	"github.com/pleimann/camel-do/templates/blocks/backlog"
	"github.com/pleimann/camel-do/templates/blocks/timeline"
	"github.com/pleimann/camel-do/templates/pages"
	"github.com/pleimann/camel-do/utils"
)

type TaskHandler struct {
	*echo.Group
	taskService    *TaskService
	projectService *project.ProjectService
}

func NewTaskHandler(
	group *echo.Group, taskService *TaskService, projectsService *project.ProjectService,
) *TaskHandler {
	taskHandler := &TaskHandler{
		Group:          group,
		taskService:    taskService,
		projectService: projectsService,
	}

	group.GET("/new", taskHandler.handleNewTask)
	group.GET("/edit/{id}", taskHandler.handleEditTask)

	group.PUT("/schedule/{id}", taskHandler.handleScheduleTask)
	group.POST("/", taskHandler.handleTaskCreate)
	group.PUT("/{id}", taskHandler.handleTaskUpdate)
	group.DELETE("/{id}", taskHandler.handleTaskDelete)
	group.PUT("/{id}/complete", taskHandler.handleTaskComplete)

	return taskHandler
}

func extractTaskId(r *http.Request) string {
	var idString string
	if r.URL.Query().Has("id") {
		idString = r.URL.Query().Get("id")
	} else {
		idString = mux.Vars(r)["id"]
	}

	return idString
}

func (h *TaskHandler) handleNewTask(c echo.Context) error {
	projectsIndex, err := h.projectService.GetProjects()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "getting projects", err)
	}

	newTaskDialogTemplate := pages.TaskDialog(projectsIndex, nil)

	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, newTaskDialogTemplate); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "render template", err)
	}

	return nil
}

func (h *TaskHandler) handleEditTask(c echo.Context) error {
	taskId := extractTaskId(c.Request())

	if task, err := h.taskService.GetTask(taskId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "getting task", err)

		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "getting task", err)
		}

	} else {
		projectsIndex, err := h.projectService.GetProjects()

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "getting projects", err)
		}

		newTaskDialogTemplate := pages.TaskDialog(projectsIndex, task)

		if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, newTaskDialogTemplate); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "render template", err)
		}
	}

	return nil
}

func (h *TaskHandler) handleScheduleTask(c echo.Context) error {
	taskId := extractTaskId(c.Request())

	// TODO Ensure duration is at least 15
	task := &model.Task{
		ID:        taskId,
		StartTime: zero.TimeFrom(time.Now().Truncate(15 * time.Minute).Add(15 * time.Minute)),
	}

	// TODO figure out when next open slot is
	if err := h.taskService.UpdateTask(task); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "scheduling task", err)

		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "scheduling task", err)
		}
	}

	projectsIndex, err := h.projectService.GetProjects()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "getting projects", err)
	}

	todaysTasks, err := h.taskService.GetTodaysTasks()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "getting todays tasks", err)
	}

	timelineViewTemplate := timeline.TimelineView(time.Now().Weekday(), todaysTasks, projectsIndex)

	return htmx.NewResponse().
		RenderTempl(c.Request().Context(), c.Response().Writer, timelineViewTemplate)
}

func (h *TaskHandler) handleTaskCreate(c echo.Context) error {
	defer c.Request().Body.Close()

	if err := c.Request().ParseForm(); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "parsing form data", err)
	}

	slog.Debug("TaskHandler.handleTaskCreate", "form", c.Request().PostForm.Encode())

	task := &model.Task{}

	if err := utils.Decoder().Decode(task, c.Request().PostForm); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "decoding form data", err)
	}

	slog.Debug("TaskHandler.handleTaskCreate", "task", task)

	if err := h.taskService.AddTask(task); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "adding task", err)
	}

	slog.Debug("TaskHandler.handleTaskCreate: get project", "projectId", task.ProjectID)

	var err error
	var project *model.Project
	if task.ProjectID.Valid {
		project, err = h.projectService.GetProject(task.ProjectID.ValueOrZero())

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return echo.NewHTTPError(http.StatusNotFound, "getting project", err)

			} else {
				return echo.NewHTTPError(http.StatusInternalServerError, "getting project", err)
			}
		}
	}

	if task.StartTime.Valid {
		// TODO Else it might belong on today's timeline but just close the dialog for now
		_, err := htmx.NewResponse().
			AddTrigger(htmx.Trigger("close-modal")).
			Reswap(htmx.SwapNone).
			RenderHTML(c.Response().Writer, template.HTML(""))

		return err

	} else {
		slog.Debug("TaskHandler.handleTaskCreate: render AddedTaskCard", "task", task)

		addedTaskTemplate := backlog.TaskCard(*task, project)

		return htmx.NewResponse().
			AddTrigger(htmx.Trigger("close-modal")).
			Retarget(backlog.BacklogSelector).
			Reswap(htmx.SwapAfterBegin).
			RenderTempl(c.Request().Context(), c.Response().Writer, addedTaskTemplate)
	}
}

func (h *TaskHandler) handleTaskUpdate(c echo.Context) error {
	defer c.Request().Body.Close()

	if err := c.Request().ParseForm(); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "parsing form data", err)
	}

	slog.Debug("TaskHandler.handleTaskUpdate", "form", c.Request().PostForm.Encode())

	taskId := extractTaskId(c.Request())

	task := &model.Task{
		ID: taskId,
	}

	if err := utils.Decoder().Decode(task, c.Request().PostForm); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "decoding form data", err)
	}

	slog.Debug("TaskHandler.handleTaskUpdate", "task", task)

	if err := h.taskService.UpdateTask(task); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "adding task", err)
	}

	slog.Debug("TaskHandler.handleTaskUpdate: get project", "projectId", task.ProjectID)

	var err error
	var project *model.Project
	if task.ProjectID.Valid {
		project, err = h.projectService.GetProject(task.ProjectID.ValueOrZero())

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return echo.NewHTTPError(http.StatusNotFound, "getting project", err)

			} else {
				return echo.NewHTTPError(http.StatusInternalServerError, "getting project", err)
			}
		}
	}

	if task.StartTime.Valid {
		// TODO Else it might belong on today's timeline but just close the dialog for now
		slog.Debug("TaskHandler.handleTaskUpdate: closing task dialog", "task", task)

		_, err := htmx.NewResponse().
			AddTrigger(htmx.Trigger("close-modal")).
			Reswap(htmx.SwapNone).
			RenderHTML(c.Response().Writer, template.HTML(""))

		return err

	} else {
		slog.Debug("TaskHandler.handleTaskUpdate: render TaskCard", "task", task)

		taskTemplate := backlog.TaskCard(*task, project)

		return htmx.NewResponse().
			AddTrigger(htmx.Trigger("close-modal")).
			Retarget(fmt.Sprintf("%s > #task-card-%s", backlog.BacklogSelector, task.ID)).
			Reswap(htmx.SwapOuterHTML).
			RenderTempl(c.Request().Context(), c.Response().Writer, taskTemplate)
	}
}

func (h *TaskHandler) handleTaskDelete(c echo.Context) error {
	defer c.Request().Body.Close()

	taskId := extractTaskId(c.Request())

	slog.Debug("TaskHandler.handleTaskDelete", "taskId", taskId)

	if err := h.taskService.DeleteTask(taskId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "deleting task", err)

		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "deleting task", err)
		}
	}

	_, err := htmx.NewResponse().
		RenderHTML(c.Response().Writer, template.HTML(""))

	return err
}

func (h *TaskHandler) handleTaskComplete(c echo.Context) error {
	defer c.Request().Body.Close()

	taskId := extractTaskId(c.Request())

	slog.Debug("TaskHandler.handleTaskComplete", "taskId", taskId)

	if err := h.taskService.CompleteToggleTask(taskId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "updating task", err)

		} else {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "updating task", err)
		}
	}

	if task, err := h.taskService.GetTask(taskId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "getting updated task", err)

		} else {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "getting updated task", err)
		}

	} else {
		var project *model.Project
		if task.ProjectID.Valid {
			project, err = h.projectService.GetProject(task.ProjectID.ValueOrZero())

			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return echo.NewHTTPError(http.StatusNotFound, "getting project", err)

				} else {
					return echo.NewHTTPError(http.StatusInternalServerError, "getting project", err)
				}
			}
		}

		taskTemplate := backlog.TaskCard(*task, project)

		htmx.NewResponse().
			RenderTempl(c.Request().Context(), c.Response().Writer, taskTemplate)

		return nil
	}
}
