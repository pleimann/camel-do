package model

import (
	"time"
)

// Project represents a project in the task tracking application.
type Project struct {
	ID        string    `jet:"column:id"`
	CreatedAt time.Time `jet:"column:created_at"`
	UpdatedAt time.Time `jet:"column:updated_at"`

	Name  string `form:"projectName" jet:"column:name"`
	Color Color  `form:"color,default:Zinc" jet:"column:color"`  // Color of the task
	Icon  Icon   `form:"icon,default:Unknown" jet:"column:icon"` // Icon to identify project
}

type ProjectIndex = map[string]Project
