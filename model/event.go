package model

import (
	"cmp"
	"iter"
	"slices"
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
	duration zero.Int32,
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
		Position:    NewTimelinePosition(startTime.Time, duration.Int32),
	}

	event := Event{
		Task: task,
	}

	return event
}

type EventList struct {
	events []Event
}

func (tl *EventList) Len() int {
	return len(tl.events)
}

func (tl *EventList) IsEmpty() bool {
	return len(tl.events) == 0
}

func (tl *EventList) All() iter.Seq[Event] {
	return slices.Values(tl.events)
}

func (tl *EventList) Push(event Event) {
	tl.events = append(tl.events, event)
}

func (tl *EventList) Sort() {
	slices.SortFunc(tl.events, func(a, b Event) int {
		if n := a.StartTime.Time.Compare(b.StartTime.Time); n != 0 {
			return n
		}

		// Times are equal order by Rank
		if n := cmp.Compare(a.Rank.Int32, b.Rank.Int32); n != 0 {
			return n
		}

		return 0
	})
}

func NewEventList() *EventList {
	return &EventList{
		events: make([]Event, 0),
	}
}
