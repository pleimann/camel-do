package task

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/gorilla/mux"

	"github.com/pleimann/camel-do/model"
	"github.com/pleimann/camel-do/services/db"
	"github.com/pleimann/camel-do/services/project"
	"github.com/pleimann/camel-do/templates/components"
	"github.com/pleimann/camel-do/templates/pages"
)

type TaskHandler struct {
	*mux.Router
	taskService    *TaskService
	projectService *project.ProjectService
}

func NewHandler(router *mux.Router, taskService *TaskService, projectsService *project.ProjectService) *TaskHandler {
	taskHandler := &TaskHandler{
		Router:         router,
		taskService:    taskService,
		projectService: projectsService,
	}

	router.HandleFunc("/new", taskHandler.handleNewTask).Methods(http.MethodGet)
	router.HandleFunc("/edit", taskHandler.handleEditTask).Methods(http.MethodGet)
	router.HandleFunc("/", taskHandler.handleTaskCreate).Methods(http.MethodPost)
	router.HandleFunc("/{id}", taskHandler.handleTaskDelete).Methods(http.MethodDelete)
	router.HandleFunc("/{id}/complete", taskHandler.handleTaskComplete).Methods(http.MethodPut)
	// router.HandleFunc("/{id}", taskHandler.handleTaskUpdate).Methods(http.MethodPut)

	return taskHandler
}

func (h *TaskHandler) handleNewTask(w http.ResponseWriter, r *http.Request) {
	projects, err := h.projectService.GetProjects()

	if err != nil {
		h.handleError(w, r, http.StatusInternalServerError, "getting projects", err)
		return
	}

	newTaskDialogTemplate := pages.TaskDialog(projects, model.Task{})

	if err := htmx.NewResponse().RenderTempl(r.Context(), w, newTaskDialogTemplate); err != nil {
		h.handleError(w, r, http.StatusInternalServerError, "render template", err)
		return
	}
}

func (h *TaskHandler) handleEditTask(w http.ResponseWriter, r *http.Request) {
	if task, err := h.taskService.GetTask(r.URL.Query().Get("id")); err != nil {
		if errors.Is(err, db.NotFoundError("not found")) {
			h.handleError(w, r, http.StatusNotFound, "getting task", err)
			return
		} else {
			h.handleError(w, r, http.StatusInternalServerError, "getting task", err)
			return
		}

	} else {
		projects, err := h.projectService.GetProjects()

		if err != nil {
			h.handleError(w, r, http.StatusInternalServerError, "getting projects", err)
			return
		}

		newTaskDialogTemplate := pages.TaskDialog(projects, *task)

		if err := htmx.NewResponse().RenderTempl(r.Context(), w, newTaskDialogTemplate); err != nil {
			h.handleError(w, r, http.StatusInternalServerError, "render template", err)
			return
		}
	}
}

func (t *TaskHandler) handleTaskCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := r.ParseForm(); err != nil {
		t.handleError(w, r, http.StatusUnprocessableEntity, "parsing form data", err)
		return
	}

	slog.Debug("handleTaskCreate", "form", r.PostForm.Encode())

	var task *model.Task = &model.Task{}

	if err := model.Decoder().Decode(task, r.PostForm); err != nil {
		t.handleError(w, r, http.StatusUnprocessableEntity, "decoding form data", err)
		return
	}

	slog.Debug("handleTaskCreate", "task", task)

	if err := t.taskService.AddTask(task); err != nil {
		t.handleError(w, r, http.StatusInternalServerError, "adding task", err)
		return
	}

	if createdTask, err := t.taskService.GetTask(task.ID); err != nil {
		t.handleError(w, r, http.StatusInternalServerError, "getting task after add", err)

	} else {
		addedTaskTemplate := components.AddedTaskCard(createdTask)

		htmx.NewResponse().
			AddTrigger(htmx.Trigger("close-modal")).
			RenderTempl(r.Context(), w, addedTaskTemplate)
	}
}

func (h *TaskHandler) handleTaskDelete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)

	id := vars["id"]
	slog.Debug("handleTaskDelete", "taskId", id)

	if err := h.taskService.DeleteTask(id); err != nil {
		if errors.Is(err, db.NotFoundError("not found")) {
			h.handleError(w, r, http.StatusNotFound, "deleting task", err)

		} else {
			h.handleError(w, r, http.StatusInternalServerError, "deleting task", err)
		}

		return
	}
}

func (h *TaskHandler) handleTaskComplete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)

	id := vars["id"]
	slog.Debug("handleTaskComplete", "taskId", id)

	if task, err := h.taskService.CompleteTask(id, true); err != nil {
		if errors.Is(err, db.NotFoundError("not found")) {
			h.handleError(w, r, http.StatusNotFound, "updating task", err)

		} else {
			h.handleError(w, r, http.StatusUnprocessableEntity, "updating task", err)
		}

		return

	} else {
		taskTemplate := components.TaskCard(task)

		htmx.NewResponse().
			RenderTempl(r.Context(), w, taskTemplate)
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
