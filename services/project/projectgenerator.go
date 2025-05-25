package project

import (
	"fmt"
	"math/rand/v2"
	"time"

	lorem "github.com/derektata/lorem/ipsum"
	"github.com/pleimann/camel-do/model"
)

// GenerateRandomProjects generates a slice of Project with random data.
func GenerateRandomProjects() ([]model.Project, error) {
	count := rand.IntN(5) + 1
	if count < 1 || count > 5 {
		return nil, fmt.Errorf("project count must be between 1 and 5, got %d", count)
	}

	projects := []model.Project{}
	for i := 0; i < count; i++ {
		project := GenerateRandomProject()
		projects = append(projects, project)
	}

	return projects, nil
}

var loremGen = lorem.NewGenerator()

// generateRandomProject generates a single project with random data.
func GenerateRandomProject() model.Project {
	// Seed the random number generator.
	rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))

	name := loremGen.Words[rand.IntN(len(loremGen.Words))]

	icon := model.IconValues()[rand.IntN(len(model.IconValues()))]

	color := model.ColorValues()[rand.IntN(len(model.ColorValues()))]

	createdAt := time.Now().Add(time.Duration(-rand.IntN(7*24)) * time.Hour)
	updatedAt := createdAt.Add(time.Duration(rand.IntN(72)) * time.Hour)
	if updatedAt.After(time.Now()) {
		updatedAt = time.Now()
	}

	return model.Project{
		Name:      name,
		Icon:      icon,
		Color:     color,
		CreatedAt: createdAt, // Set the creation timestamp
		UpdatedAt: updatedAt, // Set the update timestamp
	}
}
