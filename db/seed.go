package db

import (
	"log"
	"math/rand"
	"slices"

	"github.com/guregu/null/v6/zero"
	"github.com/pleimann/camel-do/services/project"
	"github.com/pleimann/camel-do/services/task"
)

func Seed(count int, taskService *task.TaskService, projectService *project.ProjectService) {
	projects, err := project.GenerateRandomProjects(5)
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

	projects = slices.Collect(projectsIndex.Values())

	for t := range tasks.All() {
		randProject := projects[rand.Intn(len(projects))]

		t.ProjectID = zero.StringFrom(randProject.ID)

		if err := taskService.AddTask(&t); err != nil {
			log.Fatal(err)
		}
	}
}
