package task

import (
	"log/slog"
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/gorilla/mux"

	"github.com/pleimann/camel-do/templates/pages"
)

func Routes(router *mux.Router, taskService *TaskService) {
	router.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		newTaskDialogTemplate := pages.TaskDialog(nil)

		if err := htmx.NewResponse().RenderTempl(r.Context(), w, newTaskDialogTemplate); err != nil {
			// If not, return HTTP 400 error.
			w.WriteHeader(http.StatusInternalServerError)
			slog.Error("render template", "method", r.Method, "status", http.StatusInternalServerError, "path", r.URL.Path)
			return
		}
	}).Methods(http.MethodGet)

}
