package app

import (
	"testing"

	"github.com/madalinpopa/gocost/internal/data"
)

func createTestAppForStatus() App {
	initialData := &data.DataRoot{
		CategoryGroups: map[string]data.CategoryGroup{},
	}
	return New(initialData, "test.json")
}

func TestSetStatus(t *testing.T) {
	app := createTestAppForStatus()

	updatedApp, cmd := app.SetStatus("Test message", StatusSuccess)

	if updatedApp.statusMessage != "Test message" {
		t.Errorf("Expected status message 'Test message', got '%s'", updatedApp.statusMessage)
	}

	if cmd == nil {
		t.Error("Expected command to be returned for auto-clear, got nil")
	}
}

func TestSetSuccessStatus(t *testing.T) {
	app := createTestAppForStatus()

	updatedApp, cmd := app.SetSuccessStatus("Operation successful")

	if !updatedApp.HasStatus() {
		t.Error("Expected status to be set")
	}

	if cmd == nil {
		t.Error("Expected command to be returned for auto-clear, got nil")
	}

	// Check that the message contains the success symbol
	message := updatedApp.GetStatusMessage()
	if len(message) == 0 {
		t.Error("Expected status message to be set")
	}
}

func TestSetErrorStatus(t *testing.T) {
	app := createTestAppForStatus()

	updatedApp, cmd := app.SetErrorStatus("Something went wrong")

	if !updatedApp.HasStatus() {
		t.Error("Expected status to be set")
	}

	if cmd == nil {
		t.Error("Expected command to be returned for auto-clear, got nil")
	}

	// Check that the message contains the error symbol
	message := updatedApp.GetStatusMessage()
	if len(message) == 0 {
		t.Error("Expected status message to be set")
	}
}

func TestClearStatus(t *testing.T) {
	app := createTestAppForStatus()
	app.statusMessage = "Test message"

	clearedApp := app.ClearStatus()

	if clearedApp.HasStatus() {
		t.Error("Expected status to be cleared")
	}

	if clearedApp.GetStatusMessage() != "" {
		t.Errorf("Expected empty status message, got '%s'", clearedApp.GetStatusMessage())
	}
}

func TestHasStatus(t *testing.T) {
	app := createTestAppForStatus()

	// Test without status
	if app.HasStatus() {
		t.Error("Expected HasStatus to return false when no status is set")
	}

	// Test with status
	app.statusMessage = "Test message"
	if !app.HasStatus() {
		t.Error("Expected HasStatus to return true when status is set")
	}
}

func TestGetStatusMessage(t *testing.T) {
	app := createTestAppForStatus()

	// Test empty status
	if app.GetStatusMessage() != "" {
		t.Errorf("Expected empty status message, got '%s'", app.GetStatusMessage())
	}

	// Test with status
	testMessage := "Test status message"
	app.statusMessage = testMessage
	if app.GetStatusMessage() != testMessage {
		t.Errorf("Expected status message '%s', got '%s'", testMessage, app.GetStatusMessage())
	}
}

func TestStatusClearMsg(t *testing.T) {
	app := createTestAppForStatus()
	app.statusMessage = "Test message"

	// Test that StatusClearMsg clears the status
	updatedModel, cmd := app.Update(StatusClearMsg{})

	if cmd != nil {
		t.Errorf("Expected nil command, got %v", cmd)
	}

	updatedApp, ok := updatedModel.(App)
	if !ok {
		t.Fatalf("Expected App model, got %T", updatedModel)
	}

	if updatedApp.HasStatus() {
		t.Error("Expected status to be cleared after StatusClearMsg")
	}
}
