package data

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// mockDataRoot is a helper for test data
var mockDataRoot = &DataRoot{}

func writeTempFile(t *testing.T, content []byte) string {
	t.Helper()
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "testdata.json")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return tmpFile
}

func TestLoadData_FileDoesNotExist(t *testing.T) {
	currency := "USD"
	nonExistent := filepath.Join(t.TempDir(), "doesnotexist.json")
	data, err := LoadData(nonExistent, currency)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if data == nil {
		t.Fatalf("expected new DataRoot, got nil")
	}
}

func TestLoadData_EmptyFile(t *testing.T) {
	currency := "USD"
	tmpFile := writeTempFile(t, []byte{})
	data, err := LoadData(tmpFile, currency)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if data == nil {
		t.Fatalf("expected new DataRoot, got nil")
	}
}

func TestLoadData_ValidJSON(t *testing.T) {
	currency := "USD"
	// Prepare a valid DataRoot JSON
	jsonBytes, err := json.Marshal(mockDataRoot)
	if err != nil {
		t.Fatalf("failed to marshal mock data: %v", err)
	}
	tmpFile := writeTempFile(t, jsonBytes)
	data, err := LoadData(tmpFile, currency)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if data == nil {
		t.Fatalf("expected DataRoot, got nil")
	}
}

func TestLoadData_InvalidJSON(t *testing.T) {
	currency := "USD"
	tmpFile := writeTempFile(t, []byte("{invalid json"))
	_, err := LoadData(tmpFile, currency)
	if err == nil {
		t.Fatalf("expected error for invalid JSON, got nil")
	}
	var syntaxErr *json.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Errorf("expected json.SyntaxError, got %T", err)
	}
}
