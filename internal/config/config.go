package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	CurrencyField = "currency"
	DataDirField  = "dataDir"
	DataFileField = "dataFilename"

	defaultCurrency     = "RON"
	dataDir             = ".gocost"
	defaultDataFilename = "expenses_data.json"
	defaultConfigName   = "config"
	defaultConfigType   = "json"
)

// getDataDir returns a path like "/home/user/.local/share/gocost" or whatever
// the XDG_DATA_HOME var is set to
func getDataDir() (string, error) {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", errors.New("could not get user home directory: " + err.Error())
		}
		dataHome = filepath.Join(homeDir, ".local", "share")
	}

	appDataDir := filepath.Join(dataHome, "gocost")

	if err := ensureDir(appDataDir); err != nil {
		return "", err
	}
	return appDataDir, nil
}

// getConfigDir returns a path like "/home/user/.config/gocost" or whatever
// the XDG_CONFIG_HOME var is set to
func getConfigDir() (string, error) {
    configHome := os.Getenv("XDG_CONFIG_HOME")
    if configHome == "" {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            return "", errors.New("could not get user home directory: " + err.Error())
        }
        configHome = filepath.Join(homeDir, ".config")
    }

    appConfigDir := filepath.Join(configHome, "gocost")

    if err := ensureDir(appConfigDir); err != nil {
        return "", err
    }
    return appConfigDir, nil
}

// ensureDir just creates a directory if it doesn't exist
func ensureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if mkErr := os.MkdirAll(path, 0755); mkErr != nil {
			return fmt.Errorf("could not create directory %s: %w", path, mkErr)
		}
	} else if err != nil {
		return fmt.Errorf("could not stat directory %s: %w", path, err)
	}
	return nil
}

// LoadConfig loads the configuration from the default location.
func LoadConfig() error {

	dataDirPath, err := getDataDir()
	if err != nil {
		return fmt.Errorf("failed to get data directory: %w", err)
	}

	configDirPath, err := getConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get data directory: %w", err)
	}

	viper.AddConfigPath(configDirPath)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {

		if _, ok := err.(viper.ConfigFileNotFoundError); ok {

			dataFilename := filepath.Join(dataDirPath, defaultDataFilename)

			viper.SetDefault(CurrencyField, defaultCurrency)
			viper.SetDefault(DataDirField, dataDirPath)
			viper.SetDefault(DataFileField, dataFilename)

			configFilename := filepath.Join(
				configDirPath, fmt.Sprintf("%s.%s", defaultConfigName, defaultConfigType),
			)

			if err := viper.SafeWriteConfigAs(configFilename); err != nil {
				if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
					return nil
				}
				return fmt.Errorf("failed to write config file: %w", err)
			}
		} else {
			return fmt.Errorf("failed to read config file: %w", err)
		}

	}
	return nil

}
