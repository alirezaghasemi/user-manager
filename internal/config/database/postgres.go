package database

import (
	"fmt"

	"github.com/alirezaghasemi/user-manager/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB  *gorm.DB
	cfg config.Config
}

func NewDatabaseConnection(cfg config.Config) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.Name)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("cannot connect to database: %w", err)
	}

	return &Database{
		DB:  db,
		cfg: cfg,
	}, nil
}

func (d *Database) Connection() *gorm.DB {
	return d.DB
}
