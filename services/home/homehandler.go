package home

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/angelofallars/htmx-go"
	"github.com/pleimann/camel-do/services/project"
	"github.com/pleimann/camel-do/services/task"
	"github.com/pleimann/camel-do/templates"
	"github.com/pleimann/camel-do/templates/pages"
)

type HomeHandler struct {
	taskService    *task.TaskService
	projectService *project.ProjectService
}

// HomeHandler handles a view for the index page.
func NewHomeHandler(taskService *task.TaskService, projectService *project.ProjectService) HomeHandler {
	return HomeHandler{
		taskService:    taskService,
		projectService: projectService,
	}
}

func (h HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check, if the current URL is '/'.
	if r.URL.Path != "/" {
		// If not, return HTTP 404 error.
		http.NotFound(w, r)
		slog.Error("render page", "method", r.Method, "status", http.StatusNotFound, "path", r.URL.Path)
		return
	}

	// Get backlog and tasks scheduled for today
	backlogTasks, err := h.taskService.GetBacklogTasks()
	if err != nil {
		slog.Error("get backlog tasks", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	todaysTasks, err := h.taskService.GetTodaysTasks()
	if err != nil {
		slog.Error("get tasks for today", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	slog.Info("", slog.Any("todaysTasks", todaysTasks))

	projectIndex, err := h.projectService.GetProjects()
	if err != nil {
		slog.Error("get all projects", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	weekday := time.Now().Weekday()

	// Define template layout for index page.
	indexTemplate := templates.Layout(
		templates.Config{
			Title:    "Camel Do ", // define title text
			LoginUri: "http://localhost:4000/auth/google/login",
		},
		pages.MetaTags(
			"camel-do, todo, tasks", // define meta keywords
			"Welcome to Camel Do! You're here because camels are awesome and you need more of them in your life.", // define meta description
		),
		pages.BodyContent(backlogTasks, weekday, todaysTasks, projectIndex), // define body content
	)

	// Render index page template.
	if err := htmx.NewResponse().RenderTempl(r.Context(), w, indexTemplate); err != nil {
		// If not, return HTTP 400 error.
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("render template", "method", r.Method, "status", http.StatusInternalServerError, "path", r.URL.Path)
		return
	}

	// Send log message.
	slog.Info("render page", "method", r.Method, "status", http.StatusOK, "path", r.URL.Path)
}
