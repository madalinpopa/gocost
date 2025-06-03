package app

import (
	"testing"
	"time"

	"github.com/madalinpopa/gocost/internal/data"
	"github.com/madalinpopa/gocost/internal/ui"
)

func createTestApp() App {
	testData := &data.DataRoot{
		CategoryGroups: []data.CategoryGroup{},
	}

	currentMonth := time.January
	currentYear := 2024

	return App{
		AppData: ui.AppData{
			Data:     testData,
			FilePath: "test.json",
		},
		MonthYear: ui.MonthYear{
			CurrentMonth: currentMonth,
			CurrentYear:  currentYear,
		},
		WindowSize: ui.WindowSize{
			Width:  80,
			Height: 24,
		},
		AppViews: ui.AppViews{
			MonthlyModel:       ui.NewMonthlyModel(testData, currentMonth, currentYear),
			CategoryGroupModel: ui.NewCategoryGroupModel(testData),
		},
		activeView: viewMonthlyOverview,
	}
}

func TestHandleMonthlyViewKeys_QuitKeys(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{"ctrl+c should quit", "ctrl+c"},
		{"q should quit", "q"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := createTestApp()

			model, cmd := app.handleMonthlyViewKeys(tt.key)

			if model != app {
				t.Errorf("Expected model to be unchanged, got different model")
			}

			if cmd == nil {
				t.Errorf("Expected quit command, got nil")
			}
		})
	}
}

func TestHandleMonthlyViewKeys_SwitchToCategoryGroup(t *testing.T) {
	app := createTestApp()
	app.activeView = viewMonthlyOverview

	model, cmd := app.handleMonthlyViewKeys("ctrl+g")

	resultApp, ok := model.(App)
	if !ok {
		t.Fatalf("Expected App model, got %T", model)
	}

	if resultApp.activeView != viewCategoryGroup {
		t.Errorf("Expected activeView to be viewCategoryGroup, got %v", resultApp.activeView)
	}

	if cmd != nil {
		t.Errorf("Expected nil command from CategoryGroupModel.Init(), got %v", cmd)
	}
}

func TestHandleMonthlyViewKeys_SwitchToCategoryGroup_WrongView(t *testing.T) {
	app := createTestApp()
	app.activeView = viewCategoryGroup // Not in monthly overview

	model, cmd := app.handleMonthlyViewKeys("ctrl+g")

	if model != app {
		t.Errorf("Expected model to be unchanged when not in monthly overview")
	}

	if cmd != nil {
		t.Errorf("Expected nil command when not in monthly overview, got %v", cmd)
	}
}

func TestHandleMonthlyViewKeys_PreviousMonth(t *testing.T) {
	app := createTestApp()
	app.CurrentMonth = time.March
	app.CurrentYear = 2024

	model, cmd := app.handleMonthlyViewKeys("h")

	resultApp, ok := model.(App)
	if !ok {
		t.Fatalf("Expected App model, got %T", model)
	}

	if resultApp.CurrentMonth != time.February {
		t.Errorf("Expected CurrentMonth to be February, got %v", resultApp.CurrentMonth)
	}

	if resultApp.CurrentYear != 2024 {
		t.Errorf("Expected CurrentYear to be 2024, got %v", resultApp.CurrentYear)
	}

	if cmd != nil {
		t.Errorf("Expected nil command, got %v", cmd)
	}
}

func TestHandleMonthlyViewKeys_PreviousMonth_YearBoundary(t *testing.T) {
	app := createTestApp()
	app.CurrentMonth = time.January
	app.CurrentYear = 2024

	model, _ := app.handleMonthlyViewKeys("h")

	resultApp, ok := model.(App)
	if !ok {
		t.Fatalf("Expected App model, got %T", model)
	}

	if resultApp.CurrentMonth != time.December {
		t.Errorf("Expected CurrentMonth to be December, got %v", resultApp.CurrentMonth)
	}

	if resultApp.CurrentYear != 2023 {
		t.Errorf("Expected CurrentYear to be 2023, got %v", resultApp.CurrentYear)
	}
}

func TestHandleMonthlyViewKeys_NextMonth(t *testing.T) {
	app := createTestApp()
	app.CurrentMonth = time.March
	app.CurrentYear = 2024

	model, cmd := app.handleMonthlyViewKeys("l")

	resultApp, ok := model.(App)
	if !ok {
		t.Fatalf("Expected App model, got %T", model)
	}

	if resultApp.CurrentMonth != time.April {
		t.Errorf("Expected CurrentMonth to be April, got %v", resultApp.CurrentMonth)
	}

	if resultApp.CurrentYear != 2024 {
		t.Errorf("Expected CurrentYear to be 2024, got %v", resultApp.CurrentYear)
	}

	if cmd != nil {
		t.Errorf("Expected nil command, got %v", cmd)
	}
}

func TestHandleMonthlyViewKeys_NextMonth_YearBoundary(t *testing.T) {
	app := createTestApp()
	app.CurrentMonth = time.December
	app.CurrentYear = 2024

	model, _ := app.handleMonthlyViewKeys("l")

	resultApp, ok := model.(App)
	if !ok {
		t.Fatalf("Expected App model, got %T", model)
	}

	if resultApp.CurrentMonth != time.January {
		t.Errorf("Expected CurrentMonth to be January, got %v", resultApp.CurrentMonth)
	}

	if resultApp.CurrentYear != 2025 {
		t.Errorf("Expected CurrentYear to be 2025, got %v", resultApp.CurrentYear)
	}
}

func TestHandleMonthlyViewKeys_MonthNavigation_WithNilMonthlyModel(t *testing.T) {
	app := createTestApp()
	app.MonthlyModel = nil

	// Test previous month
	model, cmd := app.handleMonthlyViewKeys("h")
	resultApp, ok := model.(App)
	if !ok {
		t.Fatalf("Expected App model, got %T", model)
	}

	if resultApp.CurrentMonth != time.December {
		t.Errorf("Expected CurrentMonth to be December, got %v", resultApp.CurrentMonth)
	}

	if cmd != nil {
		t.Errorf("Expected nil command when MonthlyModel is nil, got %v", cmd)
	}

	// Test next month from the updated app
	model, cmd = resultApp.handleMonthlyViewKeys("l")
	resultApp, ok = model.(App)
	if !ok {
		t.Fatalf("Expected App model, got %T", model)
	}

	if resultApp.CurrentMonth != time.January {
		t.Errorf("Expected CurrentMonth to be January, got %v", resultApp.CurrentMonth)
	}
}

func TestHandleMonthlyViewKeys_UnknownKey(t *testing.T) {
	app := createTestApp()

	model, cmd := app.handleMonthlyViewKeys("unknown")

	if model != app {
		t.Errorf("Expected model to be unchanged for unknown key")
	}

	if cmd != nil {
		t.Errorf("Expected nil command for unknown key, got %v", cmd)
	}
}

func TestHandleMonthlyViewKeys_EmptyKey(t *testing.T) {
	app := createTestApp()

	model, cmd := app.handleMonthlyViewKeys("")

	if model != app {
		t.Errorf("Expected model to be unchanged for empty key")
	}

	if cmd != nil {
		t.Errorf("Expected nil command for empty key, got %v", cmd)
	}
}
