package project

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/gorilla/mux"
	"github.com/pleimann/camel-do/model"
	"github.com/pleimann/camel-do/services/db"
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
	router.HandleFunc("/list", projectHandler.handleListProjects).Methods(http.MethodGet)
	router.HandleFunc("/edit/{id}", projectHandler.handleEditProject).Methods(http.MethodGet)

	router.HandleFunc("/", projectHandler.handleProjectCreate).Methods(http.MethodPost)
	router.HandleFunc("/{id}", projectHandler.handleProjectDelete).Methods(http.MethodDelete)
	router.HandleFunc("/{id}", projectHandler.handleProjectUpdate).Methods(http.MethodPut)

	return projectHandler
}

func (h *ProjectHandler) handleNewProject(w http.ResponseWriter, r *http.Request) {
	newProjectDialogTemplate := pages.ProjectDialog(model.Project{})

	if err := htmx.NewResponse().RenderTempl(r.Context(), w, newProjectDialogTemplate); err != nil {
		h.handleError(w, r, http.StatusInternalServerError, "render template", err)
		return
	}
}

func (h *ProjectHandler) handleEditProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	slog.Debug("handleEditProject", "projectId", id)

	if project, err := h.projectService.GetProject(id); err != nil {
		if errors.Is(err, db.NotFoundError("not found")) {
			h.handleError(w, r, http.StatusNotFound, "getting project", err)

		} else {
			h.handleError(w, r, http.StatusInternalServerError, "getting project", err)
		}

	} else {
		editProjectDialogTemplate := pages.ProjectDialog(project)

		if err := htmx.NewResponse().RenderTempl(r.Context(), w, editProjectDialogTemplate); err != nil {
			h.handleError(w, r, http.StatusInternalServerError, "render template", err)
			return
		}
	}
}

func (h *ProjectHandler) handleListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.projectService.GetProjects()

	if err != nil {
		h.handleError(w, r, http.StatusInternalServerError, "getting projects", err)
		return
	}

	listProjectsDialogTemplate := pages.ProjectList(projects)

	if err := htmx.NewResponse().RenderTempl(r.Context(), w, listProjectsDialogTemplate); err != nil {
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

func (h *ProjectHandler) handleProjectDelete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)

	id := vars["id"]
	slog.Debug("handleProjectDelete", "projectId", id)

	if err := h.projectService.DeleteProject(id); err != nil {
		if errors.Is(err, db.NotFoundError("not found")) {
			h.handleError(w, r, http.StatusNotFound, "deleting project", err)

		} else {
			h.handleError(w, r, http.StatusInternalServerError, "deleting project", err)
		}

		return
	}
}

func (h *ProjectHandler) handleProjectUpdate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)

	id := vars["id"]
	slog.Debug("handleProjectUpdate", "projectId", id)

	if err := r.ParseForm(); err != nil {
		if errors.Is(err, db.NotFoundError("not found")) {
			h.handleError(w, r, http.StatusNotFound, "deleting project", err)

		} else {
			h.handleError(w, r, http.StatusInternalServerError, "deleting project", err)
		}

		return
	}

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
