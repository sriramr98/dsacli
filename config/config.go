package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const AppName = "dsacli"
const DefaultDBFileName = "dsacli.db"

type Config struct {
	DbPath string
}

func NewConfig(dbPath string) Config {
	return Config{
		DbPath: dbPath,
	}
}

func NewDefaultConfig() Config {
	dbPath, err := getDBPath(DefaultDBFileName)
	if err != nil {
		panic(err)
	}

	return Config{
		DbPath: dbPath,
	}
}

// SqLite3 file is created at ~/.dsacli/dsacli.db
func getDBPath(dbFileName string) (string, error) {
	appDir, err := getAppDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(appDir, dbFileName), nil
}

// Creates a folder ~/.dsacli
func getAppDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(homeDir, fmt.Sprintf(".%s", AppName))
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}
	return appDir, nil
}
