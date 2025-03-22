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

	decoder.RegisterConverter(model.IconCircleHelp, func(input string) reflect.Value {
		if input == "" {
			return reflect.ValueOf(model.IconCircleHelp)
		}

		color, _ := model.ParseIcon(input)

		return reflect.ValueOf(color)
	})
}

func (t *TaskHandler) handleTaskCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := r.ParseForm(); err != nil {
		errorRespone := fmt.Sprintf("Error parsing form data: %s", err.Error())

		http.Error(w, errorRespone, http.StatusUnprocessableEntity)
		return
	}

	slog.Info("handleTaskCreate", "form", r.PostForm.Encode())

	var task model.Task

	if err := decoder.Decode(&task, r.PostForm); err != nil {
		errorRespone := fmt.Sprintf("Error decoding form data: %s", err.Error())

		http.Error(w, errorRespone, http.StatusUnprocessableEntity)
		return
	}

	slog.Info("handleTaskCreate", "task", task)

	if err := t.taskService.AddTask(&task); err != nil {
		errorRespone := fmt.Sprintf("<dir>%s</div>", err.Error())

		htmx.NewResponse().
			StatusCode(http.StatusInternalServerError).
			RenderHTML(w, template.HTML(errorRespone))
	}

	backlogTemplate := components.AddedTaskCard(task)

	htmx.NewResponse().
		AddTrigger(htmx.Trigger("close-modal")).
		RenderTempl(r.Context(), w, backlogTemplate)
}
