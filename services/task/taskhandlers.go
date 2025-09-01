package task

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/angelofallars/htmx-go"
	"github.com/guregu/null/v6/zero"
	"github.com/labstack/echo/v4"

	"github.com/pleimann/camel-do/model"
	"github.com/pleimann/camel-do/services/project"
	"github.com/pleimann/camel-do/templates/blocks/backlog"
	"github.com/pleimann/camel-do/templates/blocks/tasklist"
	"github.com/pleimann/camel-do/templates/components"
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

	group.POST("", taskHandler.handleCreateTask).Name = "create-task"

	group.GET("/new", taskHandler.handleNewTask).Name = "new-task"
	group.GET("/edit/:id", taskHandler.handleEditTask).Name = "edit-task"

	group.PUT("/:id", taskHandler.handleTaskUpdate).Name = "update-task"
	group.DELETE("/:id", taskHandler.handleTaskDelete).Name = "delete-task"
	group.PUT("/:id/complete", taskHandler.handleTaskComplete).Name = "complete-task"
	group.PUT("/:id/hide", taskHandler.handleTaskHide).Name = "hide-task"
	group.GET("/:id/schedule", taskHandler.handleScheduleDialog).Name = "schedule-dialog"
	group.PUT("/:id/schedule", taskHandler.handleScheduleTask).Name = "schedule-task"
	group.DELETE("/:id/schedule", taskHandler.handleUnScheduleTask).Name = "unschedule-task"

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

	task, err := h.taskService.GetTask(taskId)
	if err != nil {
		if utils.IsNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, "getting task", err)

		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "getting task", err)
		}

	}

	projectsIndex, err := h.projectService.GetProjects()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "getting projects", err)
	}

	taskDialogTemplate := pages.TaskDialog(projectsIndex, task)

	dialogTemplate := components.Dialog(taskDialogTemplate)

	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, dialogTemplate); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "render template", err)
	}

	return nil
}

func (h *TaskHandler) handleScheduleDialog(c echo.Context) error {
	taskId := extractTaskId(c)

	task, err := h.taskService.GetTask(taskId)
	if err != nil {
		if utils.IsNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, "getting task", err)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "getting task", err)
		}
	}

	scheduleDialogTemplate := components.ScheduleDialog(task)
	dialogTemplate := components.Dialog(scheduleDialogTemplate)

	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, dialogTemplate); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "render template", err)
	}

	return nil
}

func (h *TaskHandler) handleScheduleTask(c echo.Context) error {
	taskId := extractTaskId(c)

	// Parse form data
	if err := c.Request().ParseForm(); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "parsing form data", err)
	}

	dateStr := c.FormValue("date")
	timeStr := c.FormValue("time")

	// Parse date and time
	var scheduledTime time.Time
	var err error

	if dateStr == "" && timeStr == "" {
		// Fall back to default scheduling (next 15-minute slot)
		scheduledTime = time.Now().Truncate(15 * time.Minute).Add(15 * time.Minute)

	} else {
		var parsedDate time.Time

		if dateStr != "" {
			// Parse the date (YYYYMMDD format)
			parsedDate, err = time.Parse("20060102", dateStr)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid date format", err)
			}
		}

		var parsedTime time.Time
		if timeStr != "" {
			// Parse the time (12-hour format: "03 04 PM")
			parsedTime, err = time.Parse("03 04 PM", timeStr)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid time format", err)
			}
		} else {
			// Set to nearest 15-minute increment in the future
			parsedTime = time.Now().Truncate(15 * time.Minute).Add(15 * time.Minute)
		}

		// Combine date and time
		scheduledTime = time.Date(
			parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
			parsedTime.Hour(), parsedTime.Minute(), 0, 0,
			time.Local,
		)
	}

	if err := h.taskService.ScheduleTask(taskId, zero.TimeFrom(scheduledTime)); err != nil {
		if utils.IsNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, "scheduling task", err)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "scheduling task", err)
		}
	}

	task, err := h.taskService.GetTask(taskId)
	if err != nil {
		return fmt.Errorf("getting task: %w", err)
	}

	projectsIndex, err := h.projectService.GetProjects()
	if err != nil {
		return fmt.Errorf("getting projects: %w", err)
	}

	project := projectsIndex.Get(task.ProjectID.String)

	taskViewTemplate := components.Encapsulate("ul", "afterbegin:#tasklist", tasklist.TaskView(*task, project))

	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, taskViewTemplate); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "scheduling task", err)
	}

	return nil
}

func (h *TaskHandler) handleUnScheduleTask(c echo.Context) error {
	taskId := extractTaskId(c)

	if err := h.taskService.ScheduleTask(taskId, zero.TimeFromPtr(nil)); err != nil {
		if utils.IsNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, "unscheduling task", err)

		} else {
			return fmt.Errorf("unscheduling task: %w", err)
		}
	}

	task, err := h.taskService.GetTask(taskId)
	if err != nil {
		return fmt.Errorf("getting task: %w", err)
	}

	projectsIndex, err := h.projectService.GetProjects()
	if err != nil {
		return fmt.Errorf("getting projects: %w", err)
	}

	project := projectsIndex.Get(task.ProjectID.String)

	taskCardTemplate := components.Encapsulate("ul", "afterbegin:#backlog", backlog.TaskCard(*task, project))

	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, taskCardTemplate); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "scheduling task", err)
	}

	return nil
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
			if utils.IsNotFoundError(err) {
				return echo.NewHTTPError(http.StatusNotFound, "getting project", err)

			} else {
				return echo.NewHTTPError(http.StatusInternalServerError, "getting project", err)
			}
		}
	}

	if task.StartTime.Valid {
		// TODO Else it might belong on today's timeline but just close the dialog for now
		if _, err := htmx.NewResponse().
			AddTrigger(htmx.Trigger("close-modal")).
			Reswap(htmx.SwapNone).
			RenderHTML(c.Response().Writer, template.HTML("")); err != nil {
			return fmt.Errorf("task start time is invalid: %w", err)
		}

	} else {
		c.Logger().Debug("TaskHandler.handleTaskCreate: render AddedTaskCard", "task", task)

		addedTaskTemplate := backlog.TaskCard(*task, project)

		if err := htmx.NewResponse().
			AddTrigger(htmx.Trigger("close-modal")).
			Retarget(backlog.Selector).
			Reswap(htmx.SwapAfterBegin).
			RenderTempl(c.Request().Context(), c.Response().Writer, addedTaskTemplate); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "getting project", err)
		}
	}

	return nil
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
			if utils.IsNotFoundError(err) {
				return echo.NewHTTPError(http.StatusNotFound, "getting project", err)

			} else {
				return fmt.Errorf("getting project: %w", err)
			}
		}
	}

	if task.StartTime.Valid {
		// TODO Else it might belong on today's timeline but just close the dialog for now
		c.Logger().Debug("TaskHandler.handleTaskUpdate: closing task dialog", "task", task)

		if _, err := htmx.NewResponse().
			AddTrigger(htmx.Trigger("close-modal")).
			Reswap(htmx.SwapNone).
			RenderHTML(c.Response().Writer, template.HTML("")); err != nil {
			return fmt.Errorf("start time is invalid: %w", err)
		}

	} else {
		c.Logger().Debug("TaskHandler.handleTaskUpdate: render TaskCard", "task", task)

		taskTemplate := backlog.TaskCard(*task, project)

		if err := htmx.NewResponse().
			AddTrigger(htmx.Trigger("close-modal")).
			Retarget(fmt.Sprintf("#%s > #%s-%s", backlog.Selector, backlog.TaskSelector, task.ID)).
			Reswap(htmx.SwapOuterHTML).
			RenderTempl(c.Request().Context(), c.Response().Writer, taskTemplate); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "getting project", err)
		}
	}

	return nil
}

func (h *TaskHandler) handleTaskDelete(c echo.Context) error {
	defer c.Request().Body.Close()

	taskId := extractTaskId(c)

	c.Logger().Debug("TaskHandler.handleTaskDelete", "taskId", taskId)

	if err := h.taskService.DeleteTask(taskId); err != nil {
		if utils.IsNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, "deleting task", err)

		} else {
			return fmt.Errorf("deleting task: %w", err)
		}
	}

	if err := c.NoContent(http.StatusNoContent); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "rendering template", err)
	}

	return nil
}

func (h *TaskHandler) handleTaskComplete(c echo.Context) error {
	defer c.Request().Body.Close()

	taskId := extractTaskId(c)

	c.Logger().Debug("TaskHandler.handleTaskComplete", "taskId", taskId)

	if err := h.taskService.CompleteToggleTask(taskId); err != nil {
		if utils.IsNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, "updating task", err)

		} else {
			return fmt.Errorf("updating task: %w", err)
		}
	}

	if task, err := h.taskService.GetTask(taskId); err != nil {
		if utils.IsNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, "fetching updated task", err)

		} else {
			return fmt.Errorf("fetching updated task: %w", err)
		}

	} else {
		var project *model.Project
		if task.ProjectID.Valid {
			project, err = h.projectService.GetProject(task.ProjectID.ValueOrZero())

			if err != nil {
				if utils.IsNotFoundError(err) {
					return echo.NewHTTPError(http.StatusNotFound, "getting project", err)

				} else {
					return fmt.Errorf("getting project: %w", err)
				}
			}
		}

		target := c.Request().Header.Get(htmx.HeaderTarget)

		var taskTemplate templ.Component
		if strings.HasPrefix(target, backlog.TaskSelector) {
			taskTemplate = backlog.TaskCard(*task, project)

		} else if strings.HasPrefix(target, tasklist.TaskSelector) {
			taskTemplate = tasklist.TaskView(*task, project)
		}

		if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, taskTemplate); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "rendering template", err)
		}

		return nil
	}
}

func (h *TaskHandler) handleTaskHide(c echo.Context) error {
	defer c.Request().Body.Close()

	taskId := extractTaskId(c)

	c.Logger().Debug("TaskHandler.handleTaskHide", "taskId", taskId)

	if err := h.taskService.HiddenToggleTask(taskId); err != nil {
		if utils.IsNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, "updating task", err)

		} else {
			return fmt.Errorf("updating task: %w", err)
		}
	}

	return nil
}
