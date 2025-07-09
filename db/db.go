package db

import (
	"database/sql"
	"os"
	"path/filepath"
)

const (
	DBFilename string = "dsacli"
	AppName
)

func GetDB() (*sql.DB, error) {
	dbPath, err := getDBPath()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create questions table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS questions (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		url TEXT NOT NULL,
		difficulty TEXT NOT NULL,
		last_reviewed TEXT,
		sr_score INTEGER DEFAULT 0,
		attempted BOOLEAN DEFAULT 0
	);`

	if _, err := db.Exec(createTableSQL); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// SqLite3 file is created at ~/.dsacli/dsacli.db
func getDBPath() (string, error) {
	appDir, err := getAppDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(appDir, DBFilename), nil
}

// Createa a folder ~/.dsacli
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
