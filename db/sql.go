package db

import (
	"dsacli/config"
	"dsacli/types"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLDatabase struct {
	db *gorm.DB
}

func NewSQLDatabase(cfg config.Config) (Database, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&types.Question{}, &types.TodayQuestion{}); err != nil {
		return nil, err
	}

	return SQLDatabase{db: db}, nil
}
