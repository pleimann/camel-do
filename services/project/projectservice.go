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

func (s *ProjectService) GetProject(id string) (model.Project, error) {
	var project model.Project

	if err := s.db.First(&project, "id = ?", id).Error; err != nil {
		return project, err
	}

	return project, nil
}

func (s *ProjectService) UpdateProject(project *model.Project) {
	panic("unimplemented")
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

func (s *ProjectService) DeleteProject(id string) error {
	slog.Debug("deleting project", "id", id)

	if err := s.db.Delete(&model.Project{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}
