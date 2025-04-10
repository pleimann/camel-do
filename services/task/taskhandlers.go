package task

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/gorilla/mux"

	"github.com/pleimann/camel-do/model"
	"github.com/pleimann/camel-do/services/project"
	"github.com/pleimann/camel-do/templates/components"
	"github.com/pleimann/camel-do/templates/pages"
	"github.com/pleimann/camel-do/utils"
)

type TaskHandler struct {
	*mux.Router
	taskService    *TaskService
	projectService *project.ProjectService
}

func NewTaskHandler(
	router *mux.Router, taskService *TaskService, projectsService *project.ProjectService,
) *TaskHandler {
	taskHandler := &TaskHandler{
		Router:         router,
		taskService:    taskService,
		projectService: projectsService,
	}

	router.HandleFunc("/new", taskHandler.handleNewTask).Methods(http.MethodGet)
	router.HandleFunc("/edit/{id}", taskHandler.handleEditTask).Methods(http.MethodGet)
	router.HandleFunc("/", taskHandler.handleTaskCreate).Methods(http.MethodPost)
	router.HandleFunc("/{id}", taskHandler.handleTaskUpdate).Methods(http.MethodPut)
	router.HandleFunc("/{id}", taskHandler.handleTaskDelete).Methods(http.MethodDelete)
	router.HandleFunc("/{id}/complete", taskHandler.handleTaskComplete).Methods(http.MethodPut)

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

func (h *TaskHandler) handleNewTask(w http.ResponseWriter, r *http.Request) {
	projectsIndex, err := h.projectService.GetProjects()

	if err != nil {
		h.handleError(w, r, http.StatusInternalServerError, "getting projects", err)
		return
	}

	newTaskDialogTemplate := pages.TaskDialog(projectsIndex, nil)

	if err := htmx.NewResponse().RenderTempl(r.Context(), w, newTaskDialogTemplate); err != nil {
		h.handleError(w, r, http.StatusInternalServerError, "render template", err)
		return
	}
}

func (h *TaskHandler) handleEditTask(w http.ResponseWriter, r *http.Request) {
	taskId := extractTaskId(r)

	if task, err := h.taskService.GetTask(taskId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.handleError(w, r, http.StatusNotFound, "getting task", err)
			return
		} else {
			h.handleError(w, r, http.StatusInternalServerError, "getting task", err)
			return
		}

	} else {
		projectsIndex, err := h.projectService.GetProjects()

		if err != nil {
			h.handleError(w, r, http.StatusInternalServerError, "getting projects", err)
			return
		}

		newTaskDialogTemplate := pages.TaskDialog(projectsIndex, task)

		if err := htmx.NewResponse().RenderTempl(r.Context(), w, newTaskDialogTemplate); err != nil {
			h.handleError(w, r, http.StatusInternalServerError, "render template", err)
			return
		}
	}
}

func (h *TaskHandler) handleTaskCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := r.ParseForm(); err != nil {
		h.handleError(w, r, http.StatusUnprocessableEntity, "parsing form data", err)
		return
	}

	slog.Debug("TaskHandler.handleTaskCreate", "form", r.PostForm.Encode())

	task := &model.Task{}

	if err := utils.Decoder().Decode(task, r.PostForm); err != nil {
		h.handleError(w, r, http.StatusUnprocessableEntity, "decoding form data", err)
		return
	}

	slog.Debug("TaskHandler.handleTaskCreate", "task", task)

	if err := h.taskService.AddTask(task); err != nil {
		h.handleError(w, r, http.StatusInternalServerError, "adding task", err)
		return
	}

	slog.Debug("TaskHandler.handleTaskCreate: get project", "projectId", task.ProjectID)

	var err error
	var project *model.Project
	if task.ProjectID.Valid {
		project, err = h.projectService.GetProject(task.ProjectID.ValueOrZero())

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.handleError(w, r, http.StatusNotFound, "getting project", err)
				return

			} else {
				h.handleError(w, r, http.StatusInternalServerError, "getting project", err)
				return
			}
		}
	}

	if task.StartTime.Valid {
		// TODO Else it might belong on today's timeline but just close the dialog for now
		htmx.NewResponse().
			AddTrigger(htmx.Trigger("close-modal")).
			Reswap(htmx.SwapNone).
			RenderHTML(w, template.HTML(""))

	} else {
		slog.Debug("TaskHandler.handleTaskCreate: render AddedTaskCard", "task", task)

		addedTaskTemplate := components.TaskCard(*task, project)

		htmx.NewResponse().
			AddTrigger(htmx.Trigger("close-modal")).
			Retarget(components.BacklogSelector).
			Reswap(htmx.SwapAfterBegin).
			RenderTempl(r.Context(), w, addedTaskTemplate)
	}
}

func (h *TaskHandler) handleTaskUpdate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := r.ParseForm(); err != nil {
		h.handleError(w, r, http.StatusUnprocessableEntity, "parsing form data", err)
		return
	}

	slog.Debug("TaskHandler.handleTaskUpdate", "form", r.PostForm.Encode())

	taskId := extractTaskId(r)

	task := &model.Task{
		ID: taskId,
	}

	if err := utils.Decoder().Decode(task, r.PostForm); err != nil {
		h.handleError(w, r, http.StatusUnprocessableEntity, "decoding form data", err)
		return
	}

	slog.Debug("TaskHandler.handleTaskUpdate", "task", task)

	if err := h.taskService.UpdateTask(task); err != nil {
		h.handleError(w, r, http.StatusInternalServerError, "adding task", err)
		return
	}

	slog.Debug("TaskHandler.handleTaskUpdate: get project", "projectId", task.ProjectID)

	var err error
	var project *model.Project
	if task.ProjectID.Valid {
		project, err = h.projectService.GetProject(task.ProjectID.ValueOrZero())

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.handleError(w, r, http.StatusNotFound, "getting project", err)
				return

			} else {
				h.handleError(w, r, http.StatusInternalServerError, "getting project", err)
				return
			}
		}
	}

	if task.StartTime.Valid {
		// TODO Else it might belong on today's timeline but just close the dialog for now
		slog.Debug("TaskHandler.handleTaskUpdate: closing task dialog", "task", task)

		htmx.NewResponse().
			AddTrigger(htmx.Trigger("close-modal")).
			Reswap(htmx.SwapNone).
			RenderHTML(w, template.HTML(""))

	} else {
		slog.Debug("TaskHandler.handleTaskUpdate: render TaskCard", "task", task)

		taskTemplate := components.TaskCard(*task, project)

		htmx.NewResponse().
			AddTrigger(htmx.Trigger("close-modal")).
			Retarget(fmt.Sprintf("%s > #task-card-%s", components.BacklogSelector, task.ID)).
			Reswap(htmx.SwapOuterHTML).
			RenderTempl(r.Context(), w, taskTemplate)
	}
}

func (h *TaskHandler) handleTaskDelete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	taskId := extractTaskId(r)

	slog.Debug("TaskHandler.handleTaskDelete", "taskId", taskId)

	if err := h.taskService.DeleteTask(taskId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.handleError(w, r, http.StatusNotFound, "deleting task", err)

		} else {
			h.handleError(w, r, http.StatusInternalServerError, "deleting task", err)
		}
	}

	htmx.NewResponse().
		RenderHTML(w, template.HTML(""))
}

func (h *TaskHandler) handleTaskComplete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	taskId := extractTaskId(r)

	slog.Debug("TaskHandler.handleTaskComplete", "taskId", taskId)

	if err := h.taskService.CompleteToggleTask(taskId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.handleError(w, r, http.StatusNotFound, "updating task", err)

		} else {
			h.handleError(w, r, http.StatusUnprocessableEntity, "updating task", err)
		}

		return
	}

	if task, err := h.taskService.GetTask(taskId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.handleError(w, r, http.StatusNotFound, "getting updated task", err)

		} else {
			h.handleError(w, r, http.StatusUnprocessableEntity, "getting updated task", err)
		}

	} else {
		var project *model.Project
		if task.ProjectID.Valid {
			project, err = h.projectService.GetProject(task.ProjectID.ValueOrZero())

			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					h.handleError(w, r, http.StatusNotFound, "getting project", err)
					return

				} else {
					h.handleError(w, r, http.StatusInternalServerError, "getting project", err)
					return
				}
			}
		}

		taskTemplate := components.TaskCard(*task, project)

		htmx.NewResponse().
			RenderTempl(r.Context(), w, taskTemplate)
	}
}

func (t *TaskHandler) handleError(w http.ResponseWriter, r *http.Request, code int, location string, err error) {
	slog.Error(location, "error", err.Error())

	errorMessage := fmt.Sprintf("Error %s: %s", location, err.Error())

	messageTemplate := components.ErrorMessage(errorMessage)

	htmx.NewResponse().
		Retarget("#messages").
		StatusCode(code).
		RenderTempl(r.Context(), w, messageTemplate)
}
