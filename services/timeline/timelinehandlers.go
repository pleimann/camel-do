package timeline

import (
	"net/http"
	"time"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"

	"github.com/pleimann/camel-do/services/cal"
	"github.com/pleimann/camel-do/services/project"
	"github.com/pleimann/camel-do/services/task"
	"github.com/pleimann/camel-do/templates/blocks/timeline"
)

type TimelineHandler struct {
	*echo.Group
	taskService     *task.TaskService
	calendarService *cal.CalendarService
	projectService  *project.ProjectService
}

func NewTaskHandler(
	group *echo.Group,
	taskService *task.TaskService,
	calendarService *cal.CalendarService,
	projectsService *project.ProjectService,
) *TimelineHandler {
	timlineHandler := &TimelineHandler{
		Group:           group,
		taskService:     taskService,
		calendarService: calendarService,
		projectService:  projectsService,
	}

	group.GET("", timlineHandler.handleGetTimeline).Name = "get-timeline"

	return timlineHandler
}

func (h *TimelineHandler) handleGetTimeline(c echo.Context) error {
	var date time.Time = time.Now()
	if c.QueryParams().Has("date") {
		dateString := c.QueryParam("date")

		if d, err := time.ParseInLocation("20060102", dateString, time.Local); err != nil {
			return c.String(http.StatusBadRequest, "invalid date `"+dateString+"`")

		} else {
			date = d
		}
	}

	tasks, err := h.taskService.GetTasksScheduledOnDate(date)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "getting tasks", err)
	}

	events, err := h.calendarService.GetTodaysEvents()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "getting events", err)
	}

	projectsIndex, err := h.projectService.GetProjects()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "getting projects", err)
	}

	timelineViewTemplate := timeline.TimelineView(date, tasks, events, projectsIndex, nil)

	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, timelineViewTemplate); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "render template", err)
	}

	return nil
}
