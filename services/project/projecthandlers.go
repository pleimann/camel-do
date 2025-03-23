package project

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/gorilla/mux"
	"github.com/pleimann/camel-do/model"
	"github.com/pleimann/camel-do/templates/components"
	"github.com/pleimann/camel-do/templates/pages"
)

type ProjectHandler struct {
	*mux.Router
	projectService *ProjectService
}

func NewHandler(router *mux.Router, projectService *ProjectService) *ProjectHandler {
	projectHandler := &ProjectHandler{
		Router:         router,
		projectService: projectService,
	}

	router.HandleFunc("/new", projectHandler.handleNewProject).Methods(http.MethodGet)
	router.HandleFunc("/", projectHandler.handleProjectCreate).Methods(http.MethodPost)

	return projectHandler
}

func (h *ProjectHandler) handleNewProject(w http.ResponseWriter, r *http.Request) {
	newProjectDialogTemplate := pages.ProjectDialog(nil)

	if err := htmx.NewResponse().RenderTempl(r.Context(), w, newProjectDialogTemplate); err != nil {
		h.handleError(w, r, http.StatusInternalServerError, "render template", err)
		return
	}
}

func (h *ProjectHandler) handleProjectCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := r.ParseForm(); err != nil {
		h.handleError(w, r, http.StatusUnprocessableEntity, "parsing form data", err)
		return
	}

	slog.Debug("handleProjectCreate", "form", r.PostForm.Encode())

	var project model.Project

	if err := model.Decoder().Decode(&project, r.PostForm); err != nil {
		h.handleError(w, r, http.StatusUnprocessableEntity, "decoding form data", err)
		return
	}

	slog.Debug("handleProjectCreate", "project", project)

	if err := h.projectService.AddProject(&project); err != nil {
		h.handleError(w, r, http.StatusInternalServerError, "adding project", err)
		return
	}

	htmx.NewResponse().
		AddTrigger(htmx.Trigger("close-modal")).
		Write(w)
}

func (h *ProjectHandler) handleError(w http.ResponseWriter, r *http.Request, code int, location string, err error) {
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
