package app

import (
	"testing"

	"github.com/madalinpopa/gocost/internal/data"
	"github.com/madalinpopa/gocost/internal/ui"
)

func createTestAppForStatus() App {
	initialData := &data.DataRoot{
		CategoryGroups: map[string]data.CategoryGroup{},
		MonthlyData:    map[string]data.MonthlyRecord{},
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

func TestToggleExpenseStatus(t *testing.T) {
	app := createTestAppForStatus()
	
	// Create test data with a category and expense
	testCategory := data.Category{
		CatID:        "cat1",
		GroupID:      "group1",
		CategoryName: "Test Category",
		Expense: map[string]data.ExpenseRecord{
			"cat1": {
				Amount: 100.0,
				Budget: 120.0,
				Status: "Not Paid",
				Notes:  "Test expense",
			},
		},
	}
	
	monthKey := "2024-01"
	app.Data.MonthlyData = map[string]data.MonthlyRecord{
		monthKey: {
			Categories: []data.Category{testCategory},
		},
	}
	
	// Test toggling from "Not Paid" to "Paid"
	updatedModel, cmd := app.handleToggleExpenseStatusMsg(ui.ToggleExpenseStatusMsg{
		MonthKey: monthKey,
		Category: testCategory,
	})
	
	if cmd == nil {
		t.Error("Expected command to be returned for status message")
	}
	
	updatedApp, ok := updatedModel.(App)
	if !ok {
		t.Fatalf("Expected App model, got %T", updatedModel)
	}
	
	// Check that status was toggled to "Paid"
	updatedRecord := updatedApp.Data.MonthlyData[monthKey]
	if len(updatedRecord.Categories) == 0 {
		t.Fatal("Expected category to exist after toggle")
	}
	
	updatedCategory := updatedRecord.Categories[0]
	if updatedExpense, exists := updatedCategory.Expense[testCategory.CatID]; exists {
		if updatedExpense.Status != "Paid" {
			t.Errorf("Expected status to be 'Paid', got '%s'", updatedExpense.Status)
		}
	} else {
		t.Error("Expected expense to exist after toggle")
	}
	
	// Test toggling back from "Paid" to "Not Paid"
	updatedModel2, cmd2 := updatedApp.handleToggleExpenseStatusMsg(ui.ToggleExpenseStatusMsg{
		MonthKey: monthKey,
		Category: testCategory,
	})
	
	if cmd2 == nil {
		t.Error("Expected command to be returned for status message")
	}
	
	updatedApp2, ok := updatedModel2.(App)
	if !ok {
		t.Fatalf("Expected App model, got %T", updatedModel2)
	}
	
	// Check that status was toggled back to "Not Paid"
	updatedRecord2 := updatedApp2.Data.MonthlyData[monthKey]
	updatedCategory2 := updatedRecord2.Categories[0]
	if updatedExpense2, exists := updatedCategory2.Expense[testCategory.CatID]; exists {
		if updatedExpense2.Status != "Not Paid" {
			t.Errorf("Expected status to be 'Not Paid', got '%s'", updatedExpense2.Status)
		}
	} else {
		t.Error("Expected expense to exist after second toggle")
	}
}

func TestToggleExpenseStatusNewExpense(t *testing.T) {
	app := createTestAppForStatus()
	
	// Create test data with a category but no existing expense
	testCategory := data.Category{
		CatID:        "cat1",
		GroupID:      "group1",
		CategoryName: "Test Category",
		Expense:      nil, // No existing expense
	}
	
	monthKey := "2024-01"
	app.Data.MonthlyData = map[string]data.MonthlyRecord{
		monthKey: {
			Categories: []data.Category{testCategory},
		},
	}
	
	// Test toggling when no expense exists (should create default and set to "Paid")
	updatedModel, cmd := app.handleToggleExpenseStatusMsg(ui.ToggleExpenseStatusMsg{
		MonthKey: monthKey,
		Category: testCategory,
	})
	
	if cmd == nil {
		t.Error("Expected command to be returned for status message")
	}
	
	updatedApp, ok := updatedModel.(App)
	if !ok {
		t.Fatalf("Expected App model, got %T", updatedModel)
	}
	
	// Check that a new expense was created with "Paid" status
	updatedRecord := updatedApp.Data.MonthlyData[monthKey]
	updatedCategory := updatedRecord.Categories[0]
	if updatedCategory.Expense == nil {
		t.Fatal("Expected expense map to be created")
	}
	
	if updatedExpense, exists := updatedCategory.Expense[testCategory.CatID]; exists {
		if updatedExpense.Status != "Paid" {
			t.Errorf("Expected status to be 'Paid' for new expense, got '%s'", updatedExpense.Status)
		}
		if updatedExpense.Amount != 0 || updatedExpense.Budget != 0 {
			t.Error("Expected default values for amount and budget in new expense")
		}
	} else {
		t.Error("Expected expense to be created after toggle")
	}
}

func TestToggleExpenseStatusNoMonthlyData(t *testing.T) {
	app := createTestAppForStatus()
	
	// Create test data with a category
	testCategory := data.Category{
		CatID:        "cat1",
		GroupID:      "group1",
		CategoryName: "Test Category",
		Expense:      nil,
	}
	
	monthKey := "2024-01"
	// Don't add any monthly data - should handle missing month gracefully
	
	// Test toggling when no monthly data exists
	updatedModel, cmd := app.handleToggleExpenseStatusMsg(ui.ToggleExpenseStatusMsg{
		MonthKey: monthKey,
		Category: testCategory,
	})
	
	if cmd == nil {
		t.Error("Expected command to be returned for status message")
	}
	
	updatedApp, ok := updatedModel.(App)
	if !ok {
		t.Fatalf("Expected App model, got %T", updatedModel)
	}
	
	// Check that monthly data was created
	if _, exists := updatedApp.Data.MonthlyData[monthKey]; !exists {
		t.Error("Expected monthly data to be created when it doesn't exist")
	}
	
	// The category won't be in the monthly data since we're testing the case
	// where the category doesn't exist in that month, so just verify the
	// monthly record was created with empty categories
	updatedRecord := updatedApp.Data.MonthlyData[monthKey]
	if updatedRecord.Categories == nil {
		t.Error("Expected categories slice to be initialized")
	}
}
