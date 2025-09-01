package project

import (
	"encoding/gob"
	"fmt"
	"log/slog"

	"github.com/oklog/ulid/v2"
	"github.com/pleimann/camel-do/model"
	"github.com/pleimann/camel-do/utils"
	bolt "go.etcd.io/bbolt"
)

type ProjectServiceConfig struct {
}

// ProjectService is a service for managing projects to which tasks belong.
type ProjectService struct {
	config *ProjectServiceConfig
	db     *bolt.DB
}

func NewProjectService(config *ProjectServiceConfig, db *bolt.DB) (*ProjectService, error) {
	taskService := &ProjectService{
		config: config,
		db:     db,
	}

	gob.Register(model.Project{})

	return taskService, nil
}

func (s *ProjectService) GetProject(id string) (*model.Project, error) {
	slog.Debug("ProjectService.GetProject", "id", id)

	project := model.Project{}
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("projects"))

		projectBytes := bucket.Get([]byte(id))

		if projectBytes == nil {
			return utils.NewNotFoundError("project", id)
		}

		if err := project.Unmarshal(projectBytes); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("fetching project %s %w", id, err)
	}

	return &project, nil
}

func (s *ProjectService) GetProjects() (*model.ProjectIndex, error) {
	slog.Debug("ProjectService.GetProjects")

	var projectsIndex = model.NewProjectIndex()
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("projects"))

		err := bucket.ForEach(func(k, projectBytes []byte) error {
			project := model.Project{}

			if err := project.Unmarshal(projectBytes); err != nil {
				return err
			}

			projectsIndex.Add(project)

			return nil
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("fetching all projects %w", err)
	}

	return projectsIndex, nil
}

func (s *ProjectService) AddProject(project model.Project) error {
	project.ID = ulid.Make().String()

	slog.Debug("ProjectService.AddProject", "project", project)

	s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("projects"))

		if err != nil {
			return err
		}

		projectBytes, err := project.Marshal()
		if err != nil {
			return err
		}

		bucket.Put([]byte(project.ID), projectBytes)

		return nil
	})

	return nil
}

func (s *ProjectService) UpdateProject(id string, project model.Project) error {
	slog.Debug("ProjectService.UpdateProject", "project", project)

	s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("projects"))

		if err != nil {
			return err
		}

		projectBytes, err := project.Marshal()
		if err != nil {
			return err
		}

		bucket.Put([]byte(project.ID), projectBytes)

		return nil
	})

	return nil
}

func (s *ProjectService) DeleteProject(id string) error {
	slog.Debug("ProjectService.DeleteProject", "id", id)

	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("projects"))

		projectBytes := bucket.Get([]byte(id))
		if projectBytes == nil {
			return utils.NewNotFoundError("project", id)
		}

		bucket.Delete([]byte(id))

		return nil
	})

	if err != nil {
		return fmt.Errorf("ProjectService.DeleteProject (%s): %w", id, err)
	}

	return nil
}
