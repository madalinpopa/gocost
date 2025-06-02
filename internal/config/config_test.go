package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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
