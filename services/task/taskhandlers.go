package task

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/angelofallars/htmx-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

	"github.com/pleimann/camel-do/model"
	"github.com/pleimann/camel-do/templates/components"
	"github.com/pleimann/camel-do/templates/pages"
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
	router.HandleFunc("/", taskHandler.handleTaskCreate).Methods(http.MethodPost)
	router.HandleFunc("/{id}", taskHandler.handleTaskDelete).Methods(http.MethodDelete)
	// router.HandleFunc("/{id}", taskHandler.handleTaskUpdate).Methods(http.MethodPut)

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

var decoder *schema.Decoder

func init() {
	decoder = schema.NewDecoder()

	decoder.RegisterConverter(model.ColorZinc, func(input string) reflect.Value {
		if input == "" {
			return reflect.ValueOf(model.ColorZinc)
		}

		color, _ := model.ParseColor(input)

		return reflect.ValueOf(color)
	})

	decoder.RegisterConverter(model.IconUnknown, func(input string) reflect.Value {
		if input == "" {
			return reflect.ValueOf(model.IconUnknown)
		}

		color, _ := model.ParseIcon(input)

		return reflect.ValueOf(color)
	})
}

func (t *TaskHandler) handleTaskCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := r.ParseForm(); err != nil {
		t.handleError(w, r, http.StatusUnprocessableEntity, "parsing form data", err)
		return
	}

	slog.Debug("handleTaskCreate", "form", r.PostForm.Encode())

	var task model.Task

	if err := decoder.Decode(&task, r.PostForm); err != nil {
		t.handleError(w, r, http.StatusUnprocessableEntity, "decoding form data", err)
		return
	}

	slog.Debug("handleTaskCreate", "task", task)

	if err := t.taskService.AddTask(&task); err != nil {
		t.handleError(w, r, http.StatusInternalServerError, "adding task", err)
		return
	}

	backlogTemplate := components.AddedTaskCard(task)

	htmx.NewResponse().
		AddTrigger(htmx.Trigger("close-modal")).
		RenderTempl(r.Context(), w, backlogTemplate)
}

func (t *TaskHandler) handleTaskDelete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)

	id := vars["id"]
	slog.Debug("handleTaskDelete", "taskId", id)

	if err := t.taskService.DeleteTask(id); err != nil {
		t.handleError(w, r, http.StatusInternalServerError, "deleting task", err)
		return
	}
}

func (t *TaskHandler) handleError(w http.ResponseWriter, r *http.Request, code int, location string, err error) {
	var errorMessage string

	if err != nil {
		errorMessage = fmt.Sprintf("Error %s: %s", location, err.Error())
	} else {
		errorMessage = location
	}

	htmx.NewResponse().
		StatusCode(code).
		RenderHTML(w, template.HTML(errorMessage))
}
