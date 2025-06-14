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

	DefaultCurrency     = "RON"
	dataDir             = ".gocost"
	defaultDataFilename = "expenses_data.json"
	defaultConfigName   = "config"
	defaultConfigType   = "json"
)

// promptForCurrency asks the user to enter a default currency
// If testing is true, it returns the default currency without prompting
func promptForCurrency(testing bool) string {
	if testing {
		return DefaultCurrency
	}

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

	// Trim whitespace and convert to uppercase
	currency := strings.TrimSpace(input)
	currency = strings.ToUpper(currency)

	// If empty, use default
	if currency == "" {
		return DefaultCurrency
	}

	// Accept any currency code entered by the user
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

// ConfigFileExists checks if the config file exists
func ConfigFileExists() (bool, string, error) {
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
// If testing is true, it will use the default currency without prompting.
// defaultCurrency is the currency to use if the config file doesn't exist.
// configFilePath is the path to the config file.
func LoadConfig(defaultCurrency string, configFilePath string, testing bool) error {
	// Extract the directory path from the config file path
	dataDirPath := filepath.Dir(configFilePath)

	viper.AddConfigPath(dataDirPath)
	viper.SetConfigName(defaultConfigName)
	viper.SetConfigType(defaultConfigType)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// First time running the app
			dataFilename := filepath.Join(dataDirPath, defaultDataFilename)

			// Prompt for currency
			selectedCurrency := promptForCurrency(testing)

			viper.SetDefault(CurrencyField, selectedCurrency)
			viper.SetDefault(DataDirField, dataDirPath)
			viper.SetDefault(DataFileField, dataFilename)

			if err := viper.SafeWriteConfigAs(configFilePath); err != nil {
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
