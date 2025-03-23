package project

import (
	"log/slog"

	"github.com/oklog/ulid/v2"
	"github.com/pleimann/camel-do/model"
	"github.com/pleimann/camel-do/services/db"
)

type Config struct {
}

// ProjectService is a service for managing projects to which tasks belong.
type ProjectService struct {
	config *Config
	db     *db.DatabaseService
}

func NewService(config *Config, db *db.DatabaseService) (*ProjectService, error) {
	taskService := &ProjectService{
		config: config,
		db:     db,
	}

	return taskService, nil
}

func (s *ProjectService) GetProjects() ([]model.Project, error) {
	var projects []model.Project
	if err := s.db.Find(&projects).Error; err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *ProjectService) AddProject(project *model.Project) error {
	slog.Debug("adding project", "project", project)

	if project.ID == "" {
		project.ID = ulid.Make().String()
	}

	if err := s.db.Create(project).Error; err != nil {
		return err
	}

	return nil
}
