package db

import (
	"github.com/pleimann/camel-do/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseService struct {
	*gorm.DB
}

func NewDatabaseService(dbFile string) *DatabaseService {
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.Task{})

	databaseService := DatabaseService{
		DB: db,
	}

	return &databaseService
}
