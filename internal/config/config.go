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

// getDefaultDataDir returns the default data directory path.
func getDefaultDataDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("could not get user home directory: " + err.Error())
	}
	appDataDir := filepath.Join(homeDir, dataDir)

	if _, err := os.Stat(appDataDir); os.IsNotExist(err) {
		if mkErr := os.MkdirAll(appDataDir, 0755); mkErr != nil {
			return "", fmt.Errorf("could not create data directory: %w", mkErr)
		}
	} else if err != nil {
		return "", fmt.Errorf("could not stat data directory: %w", err)
	}

	return appDataDir, nil
}

// LoadConfig loads the configuration from the default location.
func LoadConfig() error {

	dataDirPath, err := getDefaultDataDir()
	if err != nil {
		return fmt.Errorf("failed to get data directory: %w", err)
	}

	viper.AddConfigPath(dataDirPath)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {

		if _, ok := err.(viper.ConfigFileNotFoundError); ok {

			dataFilename := filepath.Join(dataDirPath, defaultDataFilename)

			viper.SetDefault(CurrencyField, defaultCurrency)
			viper.SetDefault(DataDirField, dataDirPath)
			viper.SetDefault(DataFileField, dataFilename)

			configFilename := filepath.Join(
				dataDirPath, fmt.Sprintf("%s.%s", defaultConfigName, defaultConfigType),
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