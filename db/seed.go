package db

import (
	"log"
	"maps"
	"math/rand"
	"slices"

	"github.com/guregu/null/v6/zero"
	"github.com/pleimann/camel-do/services/project"
	"github.com/pleimann/camel-do/services/task"
)

func Seed(count int, taskService *task.TaskService, projectService *project.ProjectService) {
	projects, err := project.GenerateRandomProjects()
	if err != nil {
		log.Fatal(err)
	}

	tasks, err := task.GenerateRandomTasks(count)
	if err != nil {
		log.Fatal(err)
	}

	for i := range projects {
		err := projectService.AddProject(projects[i])
		if err != nil {
			return
		}
	}

	projectsIndex, _ := projectService.GetProjects()

	projects = slices.Collect(maps.Values(projectsIndex))

	for _, t := range tasks {
		randProject := projects[rand.Intn(len(projects))]

		t.ProjectID = zero.StringFrom(randProject.ID)

		if err := taskService.AddTask(&t); err != nil {
			log.Fatal(err)
		}
	}
}
