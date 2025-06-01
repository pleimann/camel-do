package model

import (
	"time"

	"github.com/guregu/null/v6/zero"
)

type Event struct {
	Task

	ConferenceData string
}

func NewEvent(
	id string,
	title zero.String,
	description zero.String,
	createdAt time.Time,
	updatedAt time.Time,
	startTime zero.Time,
	duration Duration,
	projectID zero.String,
	gTaskID zero.String,
) Event {
	task := Task{
		ID:          id,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Title:       title,
		Description: description,
		StartTime:   startTime,
		Duration:    duration,
		ProjectID:   projectID,
		GTaskID:     gTaskID,
		Position:    NewTimelinePosition(startTime.Time, duration.V),
	}

	event := Event{
		Task: task,
	}

	return event
}
