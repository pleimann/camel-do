package task

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/gorilla/mux"

	"github.com/pleimann/camel-do/model"
	"github.com/pleimann/camel-do/templates/components"
	"github.com/pleimann/camel-do/templates/pages"
	"github.com/pleimann/camel-do/utils"
)

type TaskHandler struct {
	*mux.Router
	taskService *TaskService
}

func NewTaskHandler(router *mux.Router, taskService *TaskService) *TaskHandler {
	taskHandler := &TaskHandler{
		Router:      router,
		taskService: taskService,
	}

	router.HandleFunc("/new", taskHandler.handleNewTask).Methods(http.MethodGet)
	router.HandleFunc("/new", taskHandler.handleTaskCreate).Methods(http.MethodPost)

	return taskHandler
}

func (t *TaskHandler) handleNewTask(w http.ResponseWriter, r *http.Request) {
	newTaskDialogTemplate := pages.TaskDialog(nil)

	if err := htmx.NewResponse().RenderTempl(r.Context(), w, newTaskDialogTemplate); err != nil {
		// If not, return HTTP 400 error.
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("render template", "method", r.Method, "status", http.StatusInternalServerError, "path", r.URL.Path)
		return
	}
}

func (t *TaskHandler) handleTaskCreate(w http.ResponseWriter, r *http.Request) {
	task := model.Task{}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		htmx.NewResponse().StatusCode(http.StatusBadRequest).RenderHTML(w, template.HTML("error parsing request"))
		slog.Error("parse form", "method", r.Method, "status", http.StatusBadRequest, "body", r.Body)
		return
	}
	utils.SetFormValues(r.Form, &task)

	task, err := t.taskService.AddTask(task)

	if err != nil {
		htmx.NewResponse().
			StatusCode(500).
			RenderHTML(w, template.HTML(fmt.Sprintf("<dir>%s</div>", err.Error())))
	}

	backlogTemplate := components.TaskCard(task)

	htmx.NewResponse().
		AddTrigger(htmx.Trigger("close-modal")).
		RenderTempl(r.Context(), w, backlogTemplate)
}
