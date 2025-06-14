package app

import (
	"path/filepath"
	"testing"

	"github.com/madalinpopa/gocost/internal/data"
	"github.com/madalinpopa/gocost/internal/domain"
	"github.com/madalinpopa/gocost/internal/service"
	"github.com/madalinpopa/gocost/internal/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRepo is a helper to create a test repository in a temporary directory.
func mockRepo(t *testing.T) *data.JsonRepository {
	t.Helper()
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_data.json")
	repo, err := data.NewJsonRepository(filePath, "USD")
	require.NoError(t, err)
	return repo
}

// createTestAppWithMocks creates a new App instance with mock services for testing.
func createTestAppWithMocks(t *testing.T) App {
	t.Helper()
	repo := mockRepo(t)
	categorySvc := service.NewCategoryService(repo)
	groupSvc := service.NewGroupService(repo)
	incomeSvc := service.NewIncomeService(repo)
	return New(categorySvc, groupSvc, incomeSvc, repo.FilePath())
}

func TestSetStatus(t *testing.T) {
	app := createTestAppWithMocks(t)

	updatedApp, cmd := app.SetStatus("Test message", StatusSuccess)

	if updatedApp.statusMessage != "Test message" {
		t.Errorf("Expected status message 'Test message', got '%s'", updatedApp.statusMessage)
	}

	if cmd == nil {
		t.Error("Expected command to be returned for auto-clear, got nil")
	}
}

func TestSetSuccessStatus(t *testing.T) {
	app := createTestAppWithMocks(t)

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
	app := createTestAppWithMocks(t)

	updatedApp, cmd := app.SetErrorStatus("Something went wrong")

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

func TestClearStatus(t *testing.T) {
	app := createTestAppWithMocks(t)
	app.statusMessage = "Test message"

	clearedApp := app.ClearStatus()

	if clearedApp.HasStatus() {
		t.Error("Expected status to be cleared")
	}

	if clearedApp.GetStatusMessage() != "" {
		t.Errorf("Expected empty status message, got '%s'", clearedApp.GetStatusMessage())
	}
}

func TestStatusClearMsg(t *testing.T) {
	app := createTestAppWithMocks(t)
	app.statusMessage = "Test message"

	updatedModel, cmd := app.Update(StatusClearMsg{})

	if cmd != nil {
		t.Errorf("Expected nil command, got %T", cmd)
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
	// Setup
	repo := mockRepo(t)
	categorySvc := service.NewCategoryService(repo)
	groupSvc := service.NewGroupService(repo)
	incomeSvc := service.NewIncomeService(repo)
	app := New(categorySvc, groupSvc, incomeSvc, repo.FilePath())
	monthKey := ui.GetMonthKey(app.CurrentMonth, app.CurrentYear)

	// Create test data
	testCategory := domain.Category{
		CatID:        "cat1",
		GroupID:      "group1",
		CategoryName: "Test Category",
		Expense: map[string]domain.ExpenseRecord{
			"cat1": {
				Amount: 100.0,
				Budget: 120.0,
				Status: "Not Paid",
				Notes:  "Test expense",
			},
		},
	}
	err := repo.AddCategory(monthKey, testCategory)
	require.NoError(t, err)
	app = app.refreshDataForModels()

	// Test toggling from "Not Paid" to "Paid"
	updatedModel, _ := app.handleToggleExpenseStatusMsg(ui.ToggleExpenseStatusMsg{
		MonthKey: monthKey,
		Category: testCategory,
	})

	updatedApp, ok := updatedModel.(App)
	require.True(t, ok)

	// Verify the underlying data was changed by the service
	cats, err := updatedApp.categorySvc.GetCategoriesForMonth(monthKey)
	require.NoError(t, err)
	require.Len(t, cats, 1)
	updatedExpense := cats[0].Expense["cat1"]
	assert.Equal(t, "Paid", updatedExpense.Status)

	// Test toggling back from "Paid" to "Not Paid"
	updatedModel2, _ := updatedApp.handleToggleExpenseStatusMsg(ui.ToggleExpenseStatusMsg{
		MonthKey: monthKey,
		Category: cats[0], // Use the updated category
	})

	updatedApp2, ok := updatedModel2.(App)
	require.True(t, ok)

	cats2, err := updatedApp2.categorySvc.GetCategoriesForMonth(monthKey)
	require.NoError(t, err)
	require.Len(t, cats2, 1)
	updatedExpense2 := cats2[0].Expense["cat1"]
	assert.Equal(t, "Not Paid", updatedExpense2.Status)
}
