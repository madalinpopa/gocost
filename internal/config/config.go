package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	CurrencyField = "currency"
	DataDirField  = "dataDir"
	DataFileField = "dataFilename"

	DefaultCurrency     = "USD"
	dataDir             = ".gocost"
	defaultDataFilename = "expenses_data.json"
	defaultConfigName   = "config"
	defaultConfigType   = "json"
)

// PromptForCurrency asks the user to enter a default currency
func PromptForCurrency() string {
	fmt.Println("Welcome to gocost! Please enter a default currency:")
	fmt.Println("You can change this later in the config file (~/.gocost/config.json)")
	fmt.Println("\nEnter any currency code you want to use.")

	fmt.Printf("\nEnter currency code [%s]: ", DefaultCurrency)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v. Using default: %s\n", err, DefaultCurrency)
		return DefaultCurrency
	}

	currency := strings.TrimSpace(input)
	currency = strings.ToUpper(currency)

	if currency == "" {
		return DefaultCurrency
	}

	return currency
}

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

// CheckConfigFile checks if the config file exists
func CheckConfigFile() (bool, string, error) {
	dataDirPath, err := getDefaultDataDir()
	if err != nil {
		return false, "", fmt.Errorf("failed to get data directory: %w", err)
	}

	configFilename := filepath.Join(
		dataDirPath, fmt.Sprintf("%s.%s", defaultConfigName, defaultConfigType),
	)

	_, err = os.Stat(configFilename)
	if os.IsNotExist(err) {
		return false, configFilename, nil
	}
	if err != nil {
		return false, "", fmt.Errorf("failed to check config file: %w", err)
	}

	return true, configFilename, nil
}

// LoadConfig loads the configuration from the default location.
func LoadConfig(defaultCurrency string, configFilePath string) error {
	// Extract the directory path from the config file path
	dataDirPath := filepath.Dir(configFilePath)

	viper.AddConfigPath(dataDirPath)
	viper.SetConfigName(defaultConfigName)
	viper.SetConfigType(defaultConfigType)

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			// First time running the app
			dataFilename := filepath.Join(dataDirPath, defaultDataFilename)

			viper.SetDefault(CurrencyField, defaultCurrency)
			viper.SetDefault(DataDirField, dataDirPath)
			viper.SetDefault(DataFileField, dataFilename)

			if err := viper.SafeWriteConfigAs(configFilePath); err != nil {
				var configFileAlreadyExistsError viper.ConfigFileAlreadyExistsError
				if errors.As(err, &configFileAlreadyExistsError) {
					return nil
				}
				return fmt.Errorf("failed to write config file: %w", err)
			}
		} else {
			// Handle other errors from ReadInConfig
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}
	return nil
}
