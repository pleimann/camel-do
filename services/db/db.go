package db

import (
	"fmt"

	"github.com/pleimann/camel-do/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseService struct {
	*gorm.DB
}

func NewDatabaseService(dbFile string) (*DatabaseService, error) {
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	db.AutoMigrate(&model.Task{}, &model.Project{})

	databaseService := DatabaseService{
		DB: db,
	}

	return &databaseService, nil
}
