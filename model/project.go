package model

import "gorm.io/gorm"

// Project represents a project in the task tracking application.
type Project struct {
	gorm.Model

	ID    string `gorm:"primaryKey"`
	Name  string `gorm:"unique"`
	Color Color  `gorm:"default:zinc"`    // Color of the task
	Icon  Icon   `gorm:"default:unknown"` // Icon to identify project
}
