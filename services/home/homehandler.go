package home

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/pleimann/camel-do/services/cal"
	"github.com/pleimann/camel-do/services/project"
	"github.com/pleimann/camel-do/services/task"
	"github.com/pleimann/camel-do/templates"
	"github.com/pleimann/camel-do/templates/pages"
)

type HomeHandler struct {
	taskService     *task.TaskService
	calendarService *cal.CalendarService
	projectService  *project.ProjectService
}

// HomeHandler handles a view for the index page.
func NewHomeHandler(taskService *task.TaskService, calendarService *cal.CalendarService, projectService *project.ProjectService) HomeHandler {
	return HomeHandler{
		taskService:     taskService,
		calendarService: calendarService,
		projectService:  projectService,
	}
}

func (h HomeHandler) ServeHTTP(c echo.Context) error {
	// Check, if the current URL is '/'.
	if c.Request().URL.Path != "/" {
		// If not, return HTTP 404 error.
		return echo.NewHTTPError(http.StatusNotFound, "render page method %s status path %s", c.Request().Method, c.Request().URL.Path)
	}

	todaysEvents, err := h.calendarService.GetTodaysEvents()
	if err != nil {
		msg := fmt.Sprintf("get tasks for today %s", err)
		slog.Error(msg, "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, msg)
	}
	slog.Debug("Today's Calendar Events", "events", todaysEvents)

	// Get backlog and tasks scheduled for today
	backlogTasks, err := h.taskService.GetBacklogTasks()
	if err != nil {
		msg := fmt.Sprintf("get backlog tasks %s", err)
		slog.Error(msg, "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, msg)
	}

	todaysTasks, err := h.taskService.GetTodaysTasks()
	if err != nil {
		msg := fmt.Sprintf("get tasks for today %s", err)
		slog.Error(msg, "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, msg)
	}

	projectIndex, err := h.projectService.GetProjects()
	if err != nil {
		msg := fmt.Sprintf("get all projects %s", err)
		slog.Error(msg, "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, msg)
	}

	weekday := time.Now().Weekday()

	main := pages.Main(backlogTasks, weekday, todaysTasks, projectIndex)

	// Define template layout for index page.
	indexTemplate := templates.Layout(
		templates.Config{
			Title:    "Camel Do ", // define title text
			LoginUri: "http://localhost:4000/auth/google/login",
		},
		templates.MetaTags(
			"camel-do, todo, tasks", // define meta keywords
			"Welcome to Camel Do! You're here because camels are awesome and you need more of them in your life.", // define meta description
		),
		main,
	)

	return render(c, indexTemplate)
}

func render(ctx echo.Context, cmp templ.Component) error {
	return cmp.Render(ctx.Request().Context(), ctx.Response())
}
