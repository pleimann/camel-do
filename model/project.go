package model

import (
	"time"
)

// Project represents a project in the task tracking application.
type Project struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name  string `schema:"projectName" gorm:"unique"`
	Color Color  `schema:"color,default:Zinc" gorm:"default:Zinc"`      // Color of the task
	Icon  Icon   `schema:"icon,default:Unknown" gorm:"default:Unknown"` // Icon to identify project
}
