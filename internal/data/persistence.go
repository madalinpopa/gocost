package data

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultDataFilename = "expenses_data.json"
	dataDir             = ".gocost"
)

// GetDataFilePath returns the path to the data file.
func GetDataFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user home directory: %w", err)
	}

	appDataDir := filepath.Join(homeDir, dataDir)

	// Create directory if not exists
	if _, err := os.Stat(appDataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(appDataDir, 0755); err != nil {
			return "", fmt.Errorf("could not create data directory: %w", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("could not stat data directory: %w", err)
	}

	return filepath.Join(appDataDir, defaultDataFilename), nil
}
