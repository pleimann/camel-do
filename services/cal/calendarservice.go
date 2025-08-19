package cal

import (
	"context"
	"database/sql"
	"fmt"
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

func NewCalendarService(config *CalendarServiceConfig, googleAuth *oauth.GoogleAuth, db *sql.DB) (*CalendarService, error) {
	client := googleAuth.GetClient()

	ctx := context.Background()

	service, err := calendar.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		return nil, fmt.Errorf("error creating google calendar service: %w", err)
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

	return t.getUpcomingEvents(start, time.Hour*24)
}

func (s *CalendarService) getUpcomingEvents(
	startTime time.Time,
	duration time.Duration,
) ([]model.Event, error) {
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
		return nil, err
	}

	modelEvents := []model.Event{}

	bytes, _ := events.MarshalJSON()
	fmt.Printf("calendar events: %v", string(bytes))

	for _, event := range events.Items {
		modelEvents = append(modelEvents, toModelEvent(event))
	}

	return modelEvents, nil
}

func toModelEvent(event *calendar.Event) model.Event {
	startTime, _ := time.Parse(time.RFC3339, event.Start.DateTime)
	endTime, _ := time.Parse(time.RFC3339, event.End.DateTime)
	createdTime, _ := time.Parse(time.RFC3339, event.Created)
	updatedTime, _ := time.Parse(time.RFC3339, event.Updated)

	duration := endTime.Sub(startTime)

	return model.Event{
		Task: model.Task{
			CreatedAt:   createdTime,
			UpdatedAt:   updatedTime,
			Title:       zero.StringFrom(event.Summary),
			Description: zero.StringFrom(event.Description),
			StartTime:   zero.TimeFrom(startTime),
			Duration:    zero.Int32From(int32(duration.Minutes())),
			GTaskID:     zero.StringFrom(event.Id),
		},
	}
}
