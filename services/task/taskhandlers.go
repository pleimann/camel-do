package task

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"slices"

	"github.com/angelofallars/htmx-go"
	"github.com/google/uuid"
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

func NewTaskHandler(router *mux.Router, taskService *TaskService, projectsService *project.ProjectService) *TaskHandler {
	taskHandler := &TaskHandler{
		Router:         router,
		taskService:    taskService,
		projectService: projectsService,
	}

	router.HandleFunc("/new", taskHandler.handleNewTask).Methods(http.MethodGet)
	router.HandleFunc("/edit/{id}", taskHandler.handleEditTask).Methods(http.MethodGet)
	router.HandleFunc("/", taskHandler.handleTaskCreate).Methods(http.MethodPost)
	router.HandleFunc("/{id}", taskHandler.handleTaskDelete).Methods(http.MethodDelete)
	router.HandleFunc("/{id}/complete", taskHandler.handleTaskComplete).Methods(http.MethodPut)
	// router.HandleFunc("/{id}", taskHandler.handleTaskUpdate).Methods(http.MethodPut)

	return taskHandler
}

func (h *TaskHandler) extractTaskId(r *http.Request, w http.ResponseWriter) *uuid.UUID {
	var taskIdString string
	if r.URL.Query().Has("id") {
		taskIdString = r.URL.Query().Get("id")
	} else {
		taskIdString = mux.Vars(r)["id"]
	}

	if uuid, err := uuid.Parse(taskIdString); err != nil {
		return nil

	} else {
		return &uuid
	}
}

func (h *TaskHandler) handleNewTask(w http.ResponseWriter, r *http.Request) {
	projectsIndex, err := h.projectService.GetProjects()

	if err != nil {
		h.handleError(w, r, http.StatusInternalServerError, "getting projects", err)
		return
	}

	projectValues := slices.Collect(maps.Values(projectsIndex))

	newTaskDialogTemplate := pages.TaskDialog(projectValues, model.Task{})

	if err := htmx.NewResponse().RenderTempl(r.Context(), w, newTaskDialogTemplate); err != nil {
		h.handleError(w, r, http.StatusInternalServerError, "render template", err)
		return
	}
}

func (h *TaskHandler) handleEditTask(w http.ResponseWriter, r *http.Request) {
	taskId := h.extractTaskId(r, w)

	if task, err := h.taskService.GetTask(*taskId); err != nil {
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

		projectValues := slices.Collect(maps.Values(projectsIndex))

		newTaskDialogTemplate := pages.TaskDialog(projectValues, *task)

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

	slog.Debug("TaskHandler.handleTaskCreate: get created task", "task", task)

	if task, err := h.taskService.GetTask(task.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.handleError(w, r, http.StatusNotFound, "getting updated task", err)

		} else {
			h.handleError(w, r, http.StatusUnprocessableEntity, "getting updated task", err)
		}

	} else {
		slog.Debug("TaskHandler.handleTaskCreate: get project", "projectId", task.ProjectID)

		if project, err := h.projectService.GetProject(task.ProjectID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.handleError(w, r, http.StatusNotFound, "getting project", err)

			} else {
				h.handleError(w, r, http.StatusUnprocessableEntity, "getting project", err)
			}

		} else {
			slog.Debug("TaskHandler.handleTaskCreate: render AddedTaskCard", "task", task)

			addedTaskTemplate := components.AddedTaskCard(task, project)

			htmx.NewResponse().
				AddTrigger(htmx.Trigger("close-modal")).
				RenderTempl(r.Context(), w, addedTaskTemplate)
		}
	}
}

func (h *TaskHandler) handleTaskDelete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	taskId := h.extractTaskId(r, w)

	slog.Debug("TaskHandler.handleTaskDelete", "taskId", taskId)

	if err := h.taskService.DeleteTask(*taskId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.handleError(w, r, http.StatusNotFound, "deleting task", err)

		} else {
			h.handleError(w, r, http.StatusInternalServerError, "deleting task", err)
		}
	}
}

func (h *TaskHandler) handleTaskComplete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	taskId := h.extractTaskId(r, w)

	slog.Debug("TaskHandler.handleTaskComplete", "taskId", taskId)

	if err := h.taskService.CompleteTask(*taskId, true); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.handleError(w, r, http.StatusNotFound, "updating task", err)

		} else {
			h.handleError(w, r, http.StatusUnprocessableEntity, "updating task", err)
		}

		return
	}

	if task, err := h.taskService.GetTask(*taskId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.handleError(w, r, http.StatusNotFound, "getting updated task", err)

		} else {
			h.handleError(w, r, http.StatusUnprocessableEntity, "getting updated task", err)
		}

	} else {
		if project, err := h.projectService.GetProject(task.ProjectID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.handleError(w, r, http.StatusNotFound, "getting project", err)

			} else {
				h.handleError(w, r, http.StatusUnprocessableEntity, "getting project", err)
			}

		} else {
			taskTemplate := components.TaskCard(task, project)

			htmx.NewResponse().
				RenderTempl(r.Context(), w, taskTemplate)
		}
	}
}

func (t *TaskHandler) handleError(w http.ResponseWriter, r *http.Request, code int, location string, err error) {
	slog.Error(location, "error", err.Error())

	var errorMessage string

	if err != nil {
		errorMessage = fmt.Sprintf("Error %s: %s", location, err.Error())
	} else {
		errorMessage = location
	}

	messageTemplate := components.ErrorMessage(errorMessage)

	htmx.NewResponse().
		Retarget("#messages").
		StatusCode(code).
		RenderTempl(r.Context(), w, messageTemplate)
}
