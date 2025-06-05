package model

import (
	"encoding/json"
	"time"

	"github.com/guregu/null/v6/zero"
)

// Task represents a task in the task tracking application.
type Task struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time

	Title       zero.String `form:"title"`                   // Title of the task
	Description zero.String `form:"description"`             // Description of the task
	StartTime   zero.Time   `form:"startTime"`               // Start time of the task
	Duration    zero.Int32  `form:"duration"`                // Duration of the task
	Completed   zero.Bool   `form:"completed,default:false"` // Status of task completion
	Rank        zero.Int32  // Sort order
	ProjectID   zero.String `form:"projectId"` // Foreign key referencing the project associated with the task.
	GTaskID     zero.String
	Position    TimelinePosition
}

func NewTask(
	id string,
	title zero.String,
	description zero.String,
	createdAt time.Time,
	updatedAt time.Time,
	startTime zero.Time,
	duration zero.Int32,
	completed zero.Bool,
	rank zero.Int32,
	projectID zero.String,
	gTaskID zero.String,
) Task {
	task := Task{
		ID:          id,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Title:       title,
		Description: description,
		StartTime:   startTime,
		Duration:    duration,
		Completed:   completed,
		Rank:        rank,
		ProjectID:   projectID,
		GTaskID:     gTaskID,
	}

	task.Position = NewTimelinePosition(
		task.StartTime.Time,
		duration.Int32,
	)

	return task
}

func (t Task) MarshalJSONString() string {
	json, err := t.MarshalJSON()
	if err != nil {
		return ""
	}

	return string(json)
}

func (t Task) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"id":          t.ID,
		"title":       t.Title.String,
		"description": t.Description.String,
		"startTime":   t.StartTime.Time,
		"duration":    t.Duration,
		"completed":   t.Completed.Bool,
		"rank":        t.Rank.Int32,
		"projectId":   t.ProjectID.String,
		"gTaskId":     t.GTaskID.String,
		"position":    t.Position,
	})
}

var startHours = 6
var endHours = 18
var slotMinutes = 15

type TimelinePosition struct {
	Slot int
	Size int
}

func (t TimelinePosition) MarshallJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"slot": t.Slot,
		"size": t.Size,
	})
}

func NewTimelinePosition(startTime time.Time, duration int32) TimelinePosition {
	return TimelinePosition{
		Slot: slot(startTime),
		Size: size(int(duration)),
	}
}

func slot(t time.Time) int {
	return (t.Hour()-startHours)*int(60/slotMinutes) + (t.Minute() / slotMinutes) + 1
}

func size(durationMinutes int) int {
	return max(durationMinutes, slotMinutes) / slotMinutes
}
