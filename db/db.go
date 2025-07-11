package db

import (
	"dsacli/types"
	_ "embed"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	DBFilename string = "dsacli.db"
	AppName    string = "dsacli"
)

var gormDB *gorm.DB

func init() {
	dbPath, err := getDBPath()
	if err != nil {
		panic("Failed to get database path: " + err.Error())
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	gormDB = db

	if err := gormDB.AutoMigrate(&types.Question{}, &types.TodayQuestion{}); err != nil {
		log.Println("Failed to migrate database schema:", err)
	}
}

// SqLite3 file is created at ~/.dsacli/dsacli.db
func getDBPath() (string, error) {
	appDir, err := getAppDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(appDir, DBFilename), nil
}

// Creates a folder ~/.dsacli
func getAppDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(homeDir, "."+AppName)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}
	return appDir, nil
}
