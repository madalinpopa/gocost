package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestGetDefaultDataDir(t *testing.T) {
	t.Run("successful case - directory creation", func(t *testing.T) {
		// Get expected path
		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Skipf("Cannot get user home directory: %v", err)
		}
		expectedPath := filepath.Join(homeDir, dataDir)

		// Remove directory if it exists to test creation
		os.RemoveAll(expectedPath)

		result, err := getDefaultDataDir()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result != expectedPath {
			t.Errorf("expected %s, got %s", expectedPath, result)
		}

		// Verify directory was created
		info, err := os.Stat(result)
		if err != nil {
			t.Fatalf("expected directory to be created, got error: %v", err)
		}

		if !info.IsDir() {
			t.Errorf("expected path to be a directory")
		}

		// Verify permissions
		expectedPerm := os.FileMode(0755)
		if info.Mode().Perm() != expectedPerm {
			t.Errorf("expected permissions %v, got %v", expectedPerm, info.Mode().Perm())
		}
	})

	t.Run("directory already exists", func(t *testing.T) {
		// First call to ensure directory exists
		firstResult, err := getDefaultDataDir()
		if err != nil {
			t.Fatalf("expected no error on first call, got %v", err)
		}

		// Second call should return same path without error
		secondResult, err := getDefaultDataDir()
		if err != nil {
			t.Fatalf("expected no error on second call, got %v", err)
		}

		if firstResult != secondResult {
			t.Errorf("expected same result on both calls: %s != %s", firstResult, secondResult)
		}

		// Verify directory still exists
		if _, err := os.Stat(secondResult); os.IsNotExist(err) {
			t.Errorf("expected directory to exist after second call")
		}
	})

	t.Run("returns correct path format", func(t *testing.T) {
		result, err := getDefaultDataDir()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Skipf("Cannot get user home directory: %v", err)
		}

		expectedPath := filepath.Join(homeDir, dataDir)
		if result != expectedPath {
			t.Errorf("expected %s, got %s", expectedPath, result)
		}

		// Verify the path ends with the dataDir constant
		if filepath.Base(result) != dataDir {
			t.Errorf("expected path to end with %s, got %s", dataDir, filepath.Base(result))
		}
	})

	t.Run("returns path when file exists instead of directory", func(t *testing.T) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Skipf("Cannot get user home directory: %v", err)
		}

		expectedPath := filepath.Join(homeDir, dataDir)

		// Remove directory if it exists
		os.RemoveAll(expectedPath)

		// Create a file with the same name as the expected directory
		file, err := os.Create(expectedPath)
		if err != nil {
			t.Fatalf("Could not create test file: %v", err)
		}
		file.Close()

		// Clean up after test
		defer os.Remove(expectedPath)

		// NOTE: This documents current behavior - function doesn't validate if path is a directory
		// The function returns the path even if it's a file, not a directory
		result, err := getDefaultDataDir()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result != expectedPath {
			t.Errorf("expected %s, got %s", expectedPath, result)
		}

		// Verify that the path exists but is not a directory
		info, err := os.Stat(result)
		if err != nil {
			t.Fatalf("expected path to exist, got error: %v", err)
		}

		if info.IsDir() {
			t.Errorf("expected path to be a file, but it's a directory")
		}
	})

	t.Run("isolated test with temporary directory", func(t *testing.T) {
		// Create a temporary directory for testing
		tempDir := t.TempDir()

		// Create a mock home directory structure
		mockHomeDir := filepath.Join(tempDir, "home", "testuser")
		err := os.MkdirAll(mockHomeDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create mock home directory: %v", err)
		}

		// Temporarily change HOME environment variable
		originalHome := os.Getenv("HOME")
		defer func() {
			if originalHome != "" {
				os.Setenv("HOME", originalHome)
			} else {
				os.Unsetenv("HOME")
			}
		}()
		os.Setenv("HOME", mockHomeDir)

		// Test the function
		result, err := getDefaultDataDir()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expectedPath := filepath.Join(mockHomeDir, dataDir)
		if result != expectedPath {
			t.Errorf("expected %s, got %s", expectedPath, result)
		}

		// Verify directory was created
		if _, err := os.Stat(result); os.IsNotExist(err) {
			t.Errorf("expected directory to be created")
		}
	})

	t.Run("handles nested directory creation", func(t *testing.T) {
		// Create a temporary directory for testing
		tempDir := t.TempDir()

		// Create a mock home directory that doesn't exist yet
		mockHomeDir := filepath.Join(tempDir, "nonexistent", "home", "testuser")

		// Temporarily change HOME environment variable
		originalHome := os.Getenv("HOME")
		defer func() {
			if originalHome != "" {
				os.Setenv("HOME", originalHome)
			} else {
				os.Unsetenv("HOME")
			}
		}()
		os.Setenv("HOME", mockHomeDir)

		// First create the mock home directory
		err := os.MkdirAll(mockHomeDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create mock home directory: %v", err)
		}

		// Test the function
		result, err := getDefaultDataDir()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expectedPath := filepath.Join(mockHomeDir, dataDir)
		if result != expectedPath {
			t.Errorf("expected %s, got %s", expectedPath, result)
		}

		// Verify directory was created
		info, err := os.Stat(result)
		if err != nil {
			t.Fatalf("expected directory to be created, got error: %v", err)
		}

		if !info.IsDir() {
			t.Errorf("expected path to be a directory")
		}
	})

	t.Run("handles read-only filesystem", func(t *testing.T) {
		// Create a temporary directory for testing
		tempDir := t.TempDir()

		// Create a mock home directory
		mockHomeDir := filepath.Join(tempDir, "readonly_home")
		err := os.MkdirAll(mockHomeDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create mock home directory: %v", err)
		}

		// Make the home directory read-only
		err = os.Chmod(mockHomeDir, 0444)
		if err != nil {
			t.Fatalf("Failed to make directory read-only: %v", err)
		}

		// Restore permissions after test
		defer os.Chmod(mockHomeDir, 0755)

		// Temporarily change HOME environment variable
		originalHome := os.Getenv("HOME")
		defer func() {
			if originalHome != "" {
				os.Setenv("HOME", originalHome)
			} else {
				os.Unsetenv("HOME")
			}
		}()
		os.Setenv("HOME", mockHomeDir)

		// Test the function - should fail to create directory
		_, err = getDefaultDataDir()
		if err == nil {
			t.Errorf("expected error when creating directory in read-only filesystem, got nil")
		}

		// Verify error message contains relevant information about the failure
		if !strings.Contains(err.Error(), "could not stat data directory") && !strings.Contains(err.Error(), "could not create data directory") {
			t.Errorf("expected error message about directory operation failure, got: %v", err)
		}
	})
}

func TestLoadConfig(t *testing.T) {
	t.Run("creates config when not found", func(t *testing.T) {
		// Create a temporary directory for testing
		tempDir := t.TempDir()

		// Create a mock home directory
		mockHomeDir := filepath.Join(tempDir, "home")
		err := os.MkdirAll(mockHomeDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create mock home directory: %v", err)
		}

		// Temporarily change HOME environment variable
		originalHome := os.Getenv("HOME")
		defer func() {
			if originalHome != "" {
				os.Setenv("HOME", originalHome)
			} else {
				os.Unsetenv("HOME")
			}
		}()
		os.Setenv("HOME", mockHomeDir)

		// Clear viper state
		viper.Reset()

		// Get config file path
		configDir := filepath.Join(mockHomeDir, dataDir)
		configFilePath := filepath.Join(configDir, "config.json")

		// Test the function with default currency and testing=true to bypass prompt
		err = LoadConfig(DefaultCurrency, configFilePath, true)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify config file was created
		expectedConfigPath := filepath.Join(mockHomeDir, dataDir, "config.json")
		if _, err := os.Stat(expectedConfigPath); os.IsNotExist(err) {
			t.Errorf("expected config file to be created at %s", expectedConfigPath)
		}

		// Verify default values are set
		if viper.GetString(CurrencyField) != DefaultCurrency {
			t.Errorf("expected currency %s, got %s", DefaultCurrency, viper.GetString(CurrencyField))
		}
	})

	t.Run("loads existing config", func(t *testing.T) {
		// Create a temporary directory for testing
		tempDir := t.TempDir()

		// Create a mock home directory
		mockHomeDir := filepath.Join(tempDir, "home")
		err := os.MkdirAll(mockHomeDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create mock home directory: %v", err)
		}

		// Temporarily change HOME environment variable
		originalHome := os.Getenv("HOME")
		defer func() {
			if originalHome != "" {
				os.Setenv("HOME", originalHome)
			} else {
				os.Unsetenv("HOME")
			}
		}()
		os.Setenv("HOME", mockHomeDir)

		// Create config directory
		configDir := filepath.Join(mockHomeDir, dataDir)
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create config directory: %v", err)
		}

		// Create a simple config file
		configPath := filepath.Join(configDir, "config.json")
		configContent := `{"currency": "RON", "dataDir": "/test/path", "dataFilename": "/test/data.json"}`
		err = os.WriteFile(configPath, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create config file: %v", err)
		}

		// Clear viper state
		viper.Reset()

		// Test the function with default currency and testing=true to bypass prompt
		err = LoadConfig(DefaultCurrency, configPath, true)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify config values were loaded
		if viper.GetString(CurrencyField) != "RON" {
			t.Errorf("expected currency RON, got %s", viper.GetString(CurrencyField))
		}
	})

	t.Run("loads config with partial values - missing fields are not set to defaults", func(t *testing.T) {
		// Create a temporary directory for testing
		tempDir := t.TempDir()

		// Create a mock home directory
		mockHomeDir := filepath.Join(tempDir, "home")
		err := os.MkdirAll(mockHomeDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create mock home directory: %v", err)
		}

		// Temporarily change HOME environment variable
		originalHome := os.Getenv("HOME")
		defer func() {
			if originalHome != "" {
				os.Setenv("HOME", originalHome)
			} else {
				os.Unsetenv("HOME")
			}
		}()
		os.Setenv("HOME", mockHomeDir)

		// Create config directory
		configDir := filepath.Join(mockHomeDir, dataDir)
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create config directory: %v", err)
		}

		// Create a config file with only currency field
		configPath := filepath.Join(configDir, "config.json")
		configContent := `{"currency": "EUR"}`
		err = os.WriteFile(configPath, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create config file: %v", err)
		}

		// Clear viper state
		viper.Reset()

		// Test the function with default currency and testing=true to bypass prompt
		err = LoadConfig(DefaultCurrency, configPath, true)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify currency from file is loaded
		if viper.GetString(CurrencyField) != "EUR" {
			t.Errorf("expected currency EUR from config file, got %s", viper.GetString(CurrencyField))
		}

		// Verify missing fields are empty (no defaults set for existing config files)
		if viper.GetString(DataDirField) != "" {
			t.Errorf("expected DataDirField to be empty for existing config, got %s", viper.GetString(DataDirField))
		}

		if viper.GetString(DataFileField) != "" {
			t.Errorf("expected DataFileField to be empty for existing config, got %s", viper.GetString(DataFileField))
		}
	})

	t.Run("loads complete config file with all values", func(t *testing.T) {
		// Create a temporary directory for testing
		tempDir := t.TempDir()

		// Create a mock home directory
		mockHomeDir := filepath.Join(tempDir, "home")
		err := os.MkdirAll(mockHomeDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create mock home directory: %v", err)
		}

		// Temporarily change HOME environment variable
		originalHome := os.Getenv("HOME")
		defer func() {
			if originalHome != "" {
				os.Setenv("HOME", originalHome)
			} else {
				os.Unsetenv("HOME")
			}
		}()
		os.Setenv("HOME", mockHomeDir)

		// Create config directory
		configDir := filepath.Join(mockHomeDir, dataDir)
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create config directory: %v", err)
		}

		// Create a complete config file
		configPath := filepath.Join(configDir, "config.json")
		configContent := `{
			"currency": "USD",
			"dataDir": "/custom/path",
			"dataFilename": "/custom/data.json"
		}`
		err = os.WriteFile(configPath, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create config file: %v", err)
		}

		// Clear viper state
		viper.Reset()

		// Test the function with default currency and testing=true to bypass prompt
		err = LoadConfig(DefaultCurrency, configPath, true)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify all values from file are loaded correctly
		if viper.GetString(CurrencyField) != "USD" {
			t.Errorf("expected currency USD, got %s", viper.GetString(CurrencyField))
		}

		if viper.GetString(DataDirField) != "/custom/path" {
			t.Errorf("expected dataDir /custom/path, got %s", viper.GetString(DataDirField))
		}

		if viper.GetString(DataFileField) != "/custom/data.json" {
			t.Errorf("expected dataFilename /custom/data.json, got %s", viper.GetString(DataFileField))
		}
	})

	t.Run("handles invalid JSON config file", func(t *testing.T) {
		// Create a temporary directory for testing
		tempDir := t.TempDir()

		// Create a mock home directory
		mockHomeDir := filepath.Join(tempDir, "home")
		err := os.MkdirAll(mockHomeDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create mock home directory: %v", err)
		}

		// Temporarily change HOME environment variable
		originalHome := os.Getenv("HOME")
		defer func() {
			if originalHome != "" {
				os.Setenv("HOME", originalHome)
			} else {
				os.Unsetenv("HOME")
			}
		}()
		os.Setenv("HOME", mockHomeDir)

		// Create config directory
		configDir := filepath.Join(mockHomeDir, dataDir)
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create config directory: %v", err)
		}

		// Create an invalid JSON config file
		configPath := filepath.Join(configDir, "config.json")
		configContent := `{"currency": "USD", "invalidJson":`
		err = os.WriteFile(configPath, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create config file: %v", err)
		}

		// Clear viper state
		viper.Reset()

		// Test the function - should return an error
		err = LoadConfig(DefaultCurrency, configPath, true)
		if err == nil {
			t.Errorf("expected error when reading invalid JSON config, got nil")
		}

		// Verify error message indicates config read failure
		if !strings.Contains(err.Error(), "failed to read config file") {
			t.Errorf("expected error message about config read failure, got: %v", err)
		}
	})
}

func TestGetDefaultDataDirConstants(t *testing.T) {
	// Test that constants have expected values
	if dataDir != ".gocost" {
		t.Errorf("expected dataDir to be '.gocost', got %s", dataDir)
	}
}

func TestGetDefaultDataDirPathValidation(t *testing.T) {
	t.Run("path contains expected components", func(t *testing.T) {
		result, err := getDefaultDataDir()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify path is absolute
		if !filepath.IsAbs(result) {
			t.Errorf("expected absolute path, got %s", result)
		}

		// Verify path contains the dataDir component
		if !strings.Contains(result, dataDir) {
			t.Errorf("expected path to contain %s, got %s", dataDir, result)
		}

		// Verify path doesn't contain any unexpected characters
		if strings.Contains(result, "..") {
			t.Errorf("path should not contain '..', got %s", result)
		}
	})

	t.Run("consistent results across multiple calls", func(t *testing.T) {
		results := make([]string, 5)
		for i := 0; i < len(results); i++ {
			result, err := getDefaultDataDir()
			if err != nil {
				t.Fatalf("call %d failed: %v", i, err)
			}
			results[i] = result
		}

		// All results should be identical
		for i := 1; i < len(results); i++ {
			if results[i] != results[0] {
				t.Errorf("inconsistent results: %s != %s", results[0], results[i])
			}
		}
	})
}
