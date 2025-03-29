package project

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
	"github.com/pleimann/camel-do/templates/components"
	"github.com/pleimann/camel-do/templates/pages"
	"github.com/pleimann/camel-do/utils"
)

type ProjectHandler struct {
	*mux.Router
	projectService *ProjectService
}

func NewProjectHandler(router *mux.Router, projectService *ProjectService) *ProjectHandler {
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

func (h *ProjectHandler) extractTaskId(r *http.Request, w http.ResponseWriter) *uuid.UUID {
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

func (h *ProjectHandler) handleNewProject(w http.ResponseWriter, r *http.Request) {
	newProjectDialogTemplate := pages.ProjectDialog(&model.Project{})

	if err := htmx.NewResponse().RenderTempl(r.Context(), w, newProjectDialogTemplate); err != nil {
		h.handleError(w, r, http.StatusInternalServerError, "render template", err)
		return
	}
}

func (h *ProjectHandler) handleEditProject(w http.ResponseWriter, r *http.Request) {
	id := h.extractTaskId(r, w)

	slog.Debug("ProjectHandler.handleEditProject", "projectId", id)

	if project, err := h.projectService.GetProject(*id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
	projectsIndex, err := h.projectService.GetProjects()

	if err != nil {
		h.handleError(w, r, http.StatusInternalServerError, "getting projects", err)
		return
	}

	projects := slices.Collect(maps.Values(projectsIndex))

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

	slog.Debug("ProjectHandler.handleProjectCreate", "form", r.PostForm.Encode())

	var project model.Project

	if err := utils.Decoder().Decode(&project, r.PostForm); err != nil {
		h.handleError(w, r, http.StatusUnprocessableEntity, "decoding form data", err)
		return
	}

	slog.Debug("ProjectHandler.handleProjectCreate", "project", project)

	if err := h.projectService.AddProject(project); err != nil {
		h.handleError(w, r, http.StatusInternalServerError, "adding project", err)
		return
	}

	htmx.NewResponse().
		AddTrigger(htmx.Trigger("close-modal")).
		Write(w)
}

func (h *ProjectHandler) handleProjectDelete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id := h.extractTaskId(r, w)

	slog.Debug("ProjectHandler.handleProjectDelete", "projectId", id)

	// TODO: remove projectID from linked tasks

	if err := h.projectService.DeleteProject(*id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.handleError(w, r, http.StatusNotFound, "deleting project", err)

		} else {
			h.handleError(w, r, http.StatusInternalServerError, "deleting project", err)
		}

		return
	}
}

func (h *ProjectHandler) handleProjectUpdate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id := h.extractTaskId(r, w)

	slog.Debug("ProjectHandler.handleProjectUpdate", "projectId", id)

	if err := r.ParseForm(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.handleError(w, r, http.StatusNotFound, "updating project", err)

		} else {
			h.handleError(w, r, http.StatusInternalServerError, "updating project", err)
		}

		return
	}

	slog.Debug("ProjectHandler.handleProjectUpdate", "form", r.PostForm.Encode())

	var project model.Project

	if err := utils.Decoder().Decode(&project, r.PostForm); err != nil {
		h.handleError(w, r, http.StatusUnprocessableEntity, "decoding form data", err)
		return
	}

	slog.Debug("ProjectHandler.handleProjectUpdate", "project", project)

	if err := h.projectService.UpdateProject(*id, project); err != nil {
		h.handleError(w, r, http.StatusInternalServerError, "adding project", err)
		return
	}

	htmx.NewResponse().
		Refresh(true).
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
