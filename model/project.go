package model

import (
	"iter"
	"maps"
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

type ProjectIndex struct {
	projects map[string]Project
}

func NewProjectIndex() *ProjectIndex {
	return &ProjectIndex{
		projects: make(map[string]Project),
	}
}

func (pi *ProjectIndex) All() iter.Seq2[string, Project] {
	return maps.All(pi.projects)
}

func (pi *ProjectIndex) Values() iter.Seq[Project] {
	return maps.Values(pi.projects)
}

func (pi *ProjectIndex) Add(project Project) {
	pi.projects[project.ID] = project
}

func (pi *ProjectIndex) Get(id string) *Project {
	p := pi.projects[id]
	return &p
}
