package model

import (
	"time"

	"github.com/google/uuid"

	m "github.com/pleimann/camel-do/.gen/model"
)

// Project represents a project in the task tracking application.
type Project struct {
	ID        uuid.UUID `jet:"column:id"`
	CreatedAt time.Time `jet:"column:created_at"`
	UpdatedAt time.Time `jet:"column:updated_at"`

	Name  string `schema:"projectName" jet:"column:name"`
	Color Color  `schema:"color,default:Zinc" jet:"column:color"`  // Color of the task
	Icon  Icon   `schema:"icon,default:Unknown" jet:"column:icon"` // Icon to identify project
}

type ProjectIndex = map[uuid.UUID]Project

func ConvertProjects(projects []m.Projects) (modelProject []Project) {
	modelProjects := make([]Project, len(projects))
	for i, p := range projects {
		modelProjects[i] = ConvertProject(&p)
	}

	return modelProjects
}

func ConvertProject(p *m.Projects) Project {
	id, _ := uuid.Parse(*p.ID)
	color, _ := ParseColorString(*p.Color)
	icon, _ := ParseIconString(*p.Icon)

	project := Project{
		ID:        id,
		CreatedAt: *p.CreatedAt,
		UpdatedAt: *p.UpdatedAt,
		Name:      *p.Name,
		Color:     color,
		Icon:      icon,
	}

	return project
}
