package home

import (
	"strconv"
	"time"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/pleimann/camel-do/templates/components"
)

type ComponentsService struct {
	*echo.Group
}

// NewComponentsService creates a new ComponentsService with the given Echo group
func NewComponentsService(group *echo.Group) *ComponentsService {
	service := &ComponentsService{
		Group: group,
	}

	// Register routes
	group.GET("/datepicker-calendar", service.handleDatePickerCalendar)

	return service
}

// handleDatePickerCalendar handles the HTMX request for updating the datepicker calendar
func (cs *ComponentsService) handleDatePickerCalendar(c echo.Context) error {
	monthStr := c.QueryParam("month")
	yearStr := c.QueryParam("year")
	
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		month = int(time.Now().Month())
	}
	
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		year = time.Now().Year()
	}
	
	calendar := components.DatePickerCalendar(year, time.Month(month))
	return htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, calendar)
}