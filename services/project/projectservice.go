package project

import (
	"database/sql"
	"fmt"
	"log/slog"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/oklog/ulid/v2"
	m "github.com/pleimann/camel-do/db/model"
	. "github.com/pleimann/camel-do/db/table"
	"github.com/pleimann/camel-do/model"
)

type ProjectServiceConfig struct {
}

// ProjectService is a service for managing projects to which tasks belong.
type ProjectService struct {
	config *ProjectServiceConfig
	db     *sql.DB
}

func NewService(config *ProjectServiceConfig, db *sql.DB) (*ProjectService, error) {
	taskService := &ProjectService{
		config: config,
		db:     db,
	}

	return taskService, nil
}

func (s *ProjectService) GetProject(id string) (*model.Project, error) {
	slog.Debug("ProjectService.GetProject", "id", id)

	stmt := SELECT(Projects.AllColumns).
		FROM(Projects).
		WHERE(Projects.ID.EQ(String(id)))

	var projects []m.Projects
	if err := stmt.Query(s.db, &projects); err != nil {
		return nil, fmt.Errorf("get project (%s): %w", id, err)
	}

	if len(projects) == 0 {
		return nil, sql.ErrNoRows
	}

	modelProject := model.ConvertProject(&projects[0])

	return &modelProject, nil
}

func (s *ProjectService) GetProjects() (model.ProjectIndex, error) {
	slog.Debug("ProjectService.GetProjects")

	stmt := SELECT(Projects.AllColumns).
		FROM(Projects)

	var projects []m.Projects
	err := stmt.Query(s.db, &projects) // Query directly into a slice
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}

	modelProjects := model.ConvertProjects(projects)

	projectsMap := make(model.ProjectIndex)
	for _, project := range modelProjects {
		projectsMap[project.ID] = project
	}

	return projectsMap, nil
}

func (s *ProjectService) AddProject(project model.Project) error {
	project.ID = ulid.Make().String()

	slog.Debug("ProjectService.AddProject", "project", project)

	insertStmt := Projects.
		INSERT(Projects.AllColumns).
		MODEL(project)

	if res, err := insertStmt.Exec(s.db); err != nil {
		return err
	} else {
		rows, _ := res.RowsAffected()

		slog.Debug("ProjectService.AddProject: project added", "count", rows)
	}

	return nil
}

func (s *ProjectService) UpdateProject(id string, project model.Project) error {
	slog.Debug("ProjectService.UpdateProject", "project", project)

	updateStmt := Projects.
		UPDATE(Projects.MutableColumns).
		WHERE(Projects.ID.EQ(String(id))).
		MODEL(project)

	if res, err := updateStmt.Exec(s.db); err != nil {
		return err

	} else {
		rows, _ := res.RowsAffected()

		slog.Debug("ProjectService.UpdateProject: project updated", "count", rows)
	}

	return nil
}

func (s *ProjectService) DeleteProject(id string) error {
	slog.Debug("ProjectService.DeleteProject", "id", id)

	deleteStmt := Projects.DELETE().
		WHERE(Projects.ID.EQ(String(id)))

	if res, err := deleteStmt.Exec(s.db); err != nil {
		return err

	} else {
		rows, _ := res.RowsAffected()

		slog.Debug("ProjectService.DeleteProject: project deleted", "count", rows)
	}

	return nil
}
