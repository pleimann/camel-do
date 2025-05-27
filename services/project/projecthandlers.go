package project

import (
	"database/sql"
	"errors"
	"log/slog"
	"maps"
	"net/http"
	"slices"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"

	"github.com/pleimann/camel-do/model"
	"github.com/pleimann/camel-do/templates/pages"
	"github.com/pleimann/camel-do/utils"
)

type ProjectHandler struct {
	*echo.Group
	projectService *ProjectService
}

func NewProjectHandler(group *echo.Group, projectService *ProjectService) *ProjectHandler {
	projectHandler := &ProjectHandler{
		Group:          group,
		projectService: projectService,
	}

	group.GET("/new", projectHandler.handleNewProject).Name = "new-project"
	group.GET("/list", projectHandler.handleListProjects).Name = "list-projects"
	group.GET("/edit/:id", projectHandler.handleEditProject).Name = "edit-project"

	group.POST("/", projectHandler.handleProjectCreate).Name = "create-project"
	group.DELETE("/:id", projectHandler.handleProjectDelete).Name = "delete-project"
	group.PUT("/:id", projectHandler.handleProjectUpdate).Name = "update-project"

	return projectHandler
}

func extractTaskId(c echo.Context) string {
	var idString string
	if c.QueryParams().Has("id") {
		idString = c.QueryParam("id")
	} else {
		idString = c.Param("id")
	}

	return idString
}

func (h *ProjectHandler) handleNewProject(c echo.Context) error {
	newProjectDialogTemplate := pages.ProjectDialog(nil)

	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, newProjectDialogTemplate); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "render template", err)
	}

	return nil
}

func (h *ProjectHandler) handleEditProject(c echo.Context) error {
	id := extractTaskId(c)

	slog.Debug("ProjectHandler.handleEditProject", "projectId", id)

	if project, err := h.projectService.GetProject(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "getting project", err)

		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "getting project", err)
		}

	} else {
		editProjectDialogTemplate := pages.ProjectDialog(project)

		if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, editProjectDialogTemplate); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "render template", err)
		}
	}

	return nil
}

func (h *ProjectHandler) handleListProjects(c echo.Context) error {
	projectsIndex, err := h.projectService.GetProjects()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "getting projects", err)
	}

	projects := slices.Collect(maps.Values(projectsIndex))

	listProjectsDialogTemplate := pages.ProjectList(projects)

	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, listProjectsDialogTemplate); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "render template", err)
	}

	return nil
}

func (h *ProjectHandler) handleProjectCreate(c echo.Context) error {
	defer c.Request().Body.Close()

	if err := c.Request().ParseForm(); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "parsing form data", err)
	}

	slog.Debug("ProjectHandler.handleProjectCreate", "form", c.Request().PostForm.Encode())

	project := model.Project{}

	if err := utils.Decoder().Decode(&project, c.Request().PostForm); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "decoding form data", err)
	}

	slog.Debug("ProjectHandler.handleProjectCreate", "project", project)

	if err := h.projectService.AddProject(project); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "adding project", err)
	}

	return htmx.NewResponse().
		AddTrigger(htmx.Trigger("close-modal")).
		Write(c.Response().Writer)
}

func (h *ProjectHandler) handleProjectDelete(c echo.Context) error {
	defer c.Request().Body.Close()

	id := extractTaskId(c)

	slog.Debug("ProjectHandler.handleProjectDelete", "projectId", id)

	// TODO: remove projectID from linked tasks

	if err := h.projectService.DeleteProject(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "deleting project", err)

		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "deleting project", err)
		}
	}

	return nil
}

func (h *ProjectHandler) handleProjectUpdate(c echo.Context) error {
	id := extractTaskId(c)

	c.Logger().Debug("ProjectHandler.handleProjectUpdate", "projectId", id)

	var project model.Project
	if err := c.Bind(&project); err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	c.Logger().Debug("ProjectHandler.handleProjectUpdate", "project", project)

	if err := h.projectService.UpdateProject(id, project); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "adding project", err)
	}

	return htmx.NewResponse().
		Refresh(true).
		Write(c.Response().Writer)
}
