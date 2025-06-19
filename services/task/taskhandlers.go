package task

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/angelofallars/htmx-go"
	"github.com/guregu/null/v6/zero"
	"github.com/labstack/echo/v4"

	"github.com/pleimann/camel-do/model"
	"github.com/pleimann/camel-do/services/project"
	"github.com/pleimann/camel-do/templates/blocks/backlog"
	"github.com/pleimann/camel-do/templates/blocks/tasklist"
	"github.com/pleimann/camel-do/templates/components"
	"github.com/pleimann/camel-do/templates/pages"
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

	group.GET("/new", taskHandler.handleNewTask).Name = "new-task"
	group.GET("/edit/:id", taskHandler.handleEditTask).Name = "edit-task"

	group.PUT("/schedule/:id", taskHandler.handleScheduleTask).Name = "schedule-task"
	group.DELETE("/schedule/:id", taskHandler.handleUnScheduleTask).Name = "unschedule-task"
	group.POST("/", taskHandler.handleCreateTask).Name = "create-task"
	group.PUT("/:id", taskHandler.handleTaskUpdate).Name = "update-task"
	group.DELETE("/:id", taskHandler.handleTaskDelete).Name = "delete-task"
	group.PUT("/:id/complete", taskHandler.handleTaskComplete).Name = "complete-task"

	return taskHandler
}

func extractTaskId(c echo.Context) string {
	var idString string
	if c.QueryParams().Has("id") {
		idString = c.QueryParam("id")
	} else {
		idString = c.Param("id")
	}

	return idString
}

func (h *TaskHandler) handleNewTask(c echo.Context) error {
	projectsIndex, err := h.projectService.GetProjects()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "getting projects", err)
	}

	newTaskDialogTemplate := pages.TaskDialog(projectsIndex, nil)

	dialogTemplate := components.Dialog(newTaskDialogTemplate)

	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, dialogTemplate); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "render template", err)
	}

	return nil
}

func (h *TaskHandler) handleEditTask(c echo.Context) error {
	taskId := extractTaskId(c)

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

		taskDialogTemplate := pages.TaskDialog(projectsIndex, task)

		dialogTemplate := components.Dialog(taskDialogTemplate)

		if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, dialogTemplate); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "render template", err)
		}
	}

	return nil
}

func (h *TaskHandler) handleScheduleTask(c echo.Context) error {
	taskId := extractTaskId(c)

	// TODO Ensure duration is at least 15
	task := model.Task{
		ID:        taskId,
		StartTime: zero.TimeFrom(time.Now().Truncate(15 * time.Minute).Add(15 * time.Minute)),
	}

	// TODO figure out when next open slot is
	if err := h.taskService.UpdateTask(&task); err != nil {
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

	taskViewTemplate := tasklist.TaskView(task, projectsIndex)

	return htmx.NewResponse().
		RenderTempl(c.Request().Context(), c.Response().Writer, taskViewTemplate)
}

func (h *TaskHandler) handleUnScheduleTask(c echo.Context) error {
	taskId := extractTaskId(c)

	// TODO Ensure duration is at least 15
	task := model.Task{
		ID:        taskId,
		StartTime: zero.Time{},
	}

	// TODO figure out when next open slot is
	if err := h.taskService.UpdateTask(&task); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "unscheduling task", err)

		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "unscheduling task", err)
		}
	}

	projectsIndex, err := h.projectService.GetProjects()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "getting projects", err)
	}

	taskViewTemplate := tasklist.TaskView(task, projectsIndex)

	return htmx.NewResponse().
		RenderTempl(c.Request().Context(), c.Response().Writer, taskViewTemplate)
}

func (h *TaskHandler) handleCreateTask(c echo.Context) error {
	task := &model.Task{}
	if err := c.Bind(task); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "decoding form data", err)
	}

	c.Logger().Debug("TaskHandler.handleTaskCreate", "task", task)

	if err := h.taskService.AddTask(task); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "adding task", err)
	}

	c.Logger().Debug("TaskHandler.handleTaskCreate: get project", "projectId", task.ProjectID)

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
		c.Logger().Debug("TaskHandler.handleTaskCreate: render AddedTaskCard", "task", task)

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

	c.Logger().Debug("TaskHandler.handleTaskUpdate", "form", c.Request().PostForm.Encode())

	taskId := extractTaskId(c)

	task := &model.Task{
		ID: taskId,
	}

	if err := c.Bind(task); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "decoding form data", err)
	}

	c.Logger().Debug("TaskHandler.handleTaskUpdate", "task", task)

	if err := h.taskService.UpdateTask(task); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "adding task", err)
	}

	c.Logger().Debug("TaskHandler.handleTaskUpdate: get project", "projectId", task.ProjectID)

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
		c.Logger().Debug("TaskHandler.handleTaskUpdate: closing task dialog", "task", task)

		_, err := htmx.NewResponse().
			AddTrigger(htmx.Trigger("close-modal")).
			Reswap(htmx.SwapNone).
			RenderHTML(c.Response().Writer, template.HTML(""))

		return err

	} else {
		c.Logger().Debug("TaskHandler.handleTaskUpdate: render TaskCard", "task", task)

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

	taskId := extractTaskId(c)

	c.Logger().Debug("TaskHandler.handleTaskDelete", "taskId", taskId)

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

	taskId := extractTaskId(c)

	c.Logger().Debug("TaskHandler.handleTaskComplete", "taskId", taskId)

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
