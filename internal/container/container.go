package container

import (
	"github.com/alirezaghasemi/user-manager/internal/config"
	"github.com/alirezaghasemi/user-manager/internal/config/database"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Container struct {
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewContainer(cfg config.Config) *Container {
	// Database
	db, err := database.NewDatabaseConnection(cfg)
	if err != nil {
		panic(err)
	}

	// Validator
	validate := validator.New()

	return &Container{
		DB:       db.Connection(),
		Validate: validate,
	}
}
