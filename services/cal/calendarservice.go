package cal

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"time"

	"github.com/guregu/null/v6/zero"
	"github.com/pleimann/camel-do/model"
	"github.com/pleimann/camel-do/services/oauth"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type CalendarServiceConfig struct {
}

// TaskService is a service for managing tasks.
type CalendarService struct {
	config         *CalendarServiceConfig
	db             *sql.DB
	googleCalendar *calendar.Service
}

func NewCalendarService(config *CalendarServiceConfig, db *sql.DB) (*CalendarService, error) {
	client := oauth.NewGoogleAuth().GetClient()

	ctx := context.Background()

	service, err := calendar.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		slog.Error("error creating google tasks service", "error", err.Error())
		return nil, err
	}

	calendarService := &CalendarService{
		config:         config,
		db:             db,
		googleCalendar: service,
	}

	return calendarService, nil
}

func (t *CalendarService) GetTodaysEvents() ([]model.Event, error) {
	year, month, day := time.Now().Date()

	start := time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())

	events := t.getUpcomingEvents(start, time.Hour*24)

	return events, nil
}

func (s *CalendarService) getUpcomingEvents(
	startTime time.Time,
	duration time.Duration,
) []model.Event {
	events, err := s.googleCalendar.Events.
		List("primary").
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(startTime.Local().Format(time.RFC3339)).
		TimeMax(startTime.Add(duration).Format(time.RFC3339)).
		MaxResults(10).
		OrderBy("startTime").
		Do()

	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	modelEvents := []model.Event{}

	for _, event := range events.Items {
		modelEvents = append(modelEvents, toModelEvent(event))
	}

	return modelEvents
}

func toModelEvent(event *calendar.Event) model.Event {
	startTime, _ := time.Parse(time.RFC3339, event.Start.DateTime)

	return model.Event{
		Task: model.Task{
			GTaskID:     zero.StringFrom(event.Id),
			Title:       zero.StringFrom(event.Description),
			Description: zero.StringFrom(event.ConferenceData.Notes),
			StartTime:   zero.TimeFrom(startTime),
		},
	}
}
