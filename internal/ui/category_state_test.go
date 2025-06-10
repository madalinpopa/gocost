package ui

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/madalinpopa/gocost/internal/data"
)

func createTestCategoryModel() CategoryModel {
	initialData := &data.DataRoot{
		CategoryGroups: map[string]data.CategoryGroup{
			"group1": {
				GroupID:   "group1",
				GroupName: "Test Group",
				Order:     1,
			},
		},
		MonthlyData: map[string]data.MonthlyRecord{
			"January-2024": {
				Categories: []data.Category{
					{
						CatID:        "cat1",
						GroupID:      "group1",
						CategoryName: "Test Category",
						Expense:      make(map[string]data.ExpenseRecord),
					},
				},
			},
		},
	}
	
	return NewCategoryModel(initialData, time.January, 2024)
}

func createTestCategoryGroupModel() CategoryGroupModel {
	initialData := &data.DataRoot{
		CategoryGroups: map[string]data.CategoryGroup{
			"group1": {
				GroupID:   "group1",
				GroupName: "Test Group",
				Order:     1,
			},
		},
		MonthlyData: map[string]data.MonthlyRecord{},
	}
	
	return NewCategoryGroupModel(initialData)
}

func TestCategoryModel_SetMonthYear_ResetsEditingState(t *testing.T) {
	model := createTestCategoryModel()
	
	// Simulate being in add category mode
	model.addCategory = true
	model.selectedGroup = data.CategoryGroup{GroupID: "test", GroupName: "Test"}
	model.isEditingName = true
	model.editingIndex = 2
	model.editInput.Focus()
	model.editInput.SetValue("Test Value")
	
	// Also set some filter state
	model.isFiltered = true
	model.filterText = "test filter"
	model.filteredCategories = []data.Category{{CatID: "test"}}
	
	// Set some move state
	model.moveCategory = true
	model.movingCategory = data.Category{CatID: "moving"}
	
	// Call SetMonthYear - this should reset all editing state
	updatedModel := model.SetMonthYear(time.February, 2024)
	
	// Verify all editing flags are reset
	if updatedModel.addCategory {
		t.Error("Expected addCategory to be false after SetMonthYear")
	}
	
	if updatedModel.isEditingName {
		t.Error("Expected isEditingName to be false after SetMonthYear")
	}
	
	if updatedModel.editingIndex != -1 {
		t.Errorf("Expected editingIndex to be -1, got %d", updatedModel.editingIndex)
	}
	
	if updatedModel.editInput.Value() != "" {
		t.Errorf("Expected editInput value to be empty, got '%s'", updatedModel.editInput.Value())
	}
	
	if updatedModel.editInput.Focused() {
		t.Error("Expected editInput to be blurred after SetMonthYear")
	}
	
	// Verify filter state is reset
	if updatedModel.isFiltered {
		t.Error("Expected isFiltered to be false after SetMonthYear")
	}
	
	if updatedModel.filterText != "" {
		t.Errorf("Expected filterText to be empty, got '%s'", updatedModel.filterText)
	}
	
	if len(updatedModel.filteredCategories) != 0 {
		t.Errorf("Expected filteredCategories to be empty, got %d items", len(updatedModel.filteredCategories))
	}
	
	// Verify move state is reset
	if updatedModel.moveCategory {
		t.Error("Expected moveCategory to be false after SetMonthYear")
	}
	
	if updatedModel.movingCategory.CatID != "" {
		t.Error("Expected movingCategory to be empty after SetMonthYear")
	}
	
	if updatedModel.selectedGroup.GroupID != "" {
		t.Error("Expected selectedGroup to be empty after SetMonthYear")
	}
}

func TestCategoryModel_UpdateData_ResetsEditingState(t *testing.T) {
	model := createTestCategoryModel()
	
	// Set up editing state
	model.addCategory = true
	model.isEditingName = true
	model.editingIndex = 1
	model.editInput.SetValue("Test Value")
	model.isFiltered = true
	model.filterText = "filter"
	
	// Create new data
	newData := &data.DataRoot{
		CategoryGroups: map[string]data.CategoryGroup{},
		MonthlyData:    map[string]data.MonthlyRecord{},
	}
	
	// Call UpdateData
	updatedModel := model.UpdateData(newData)
	
	// Verify editing state is reset
	if updatedModel.addCategory {
		t.Error("Expected addCategory to be false after UpdateData")
	}
	
	if updatedModel.isEditingName {
		t.Error("Expected isEditingName to be false after UpdateData")
	}
	
	if updatedModel.editingIndex != -1 {
		t.Errorf("Expected editingIndex to be -1, got %d", updatedModel.editingIndex)
	}
	
	if updatedModel.isFiltered {
		t.Error("Expected isFiltered to be false after UpdateData")
	}
	
	if updatedModel.filterText != "" {
		t.Error("Expected filterText to be empty after UpdateData")
	}
}

func TestCategoryModel_ResetEditingState(t *testing.T) {
	model := createTestCategoryModel()
	
	// Set up all possible editing states
	model.addCategory = true
	model.moveCategory = true
	model.selectedGroup = data.CategoryGroup{GroupID: "test", GroupName: "Test"}
	model.movingCategory = data.Category{CatID: "moving", CategoryName: "Moving"}
	model.isEditingName = true
	model.editingIndex = 5
	model.editInput.Focus()
	model.editInput.SetValue("Test Value")
	model.isFiltered = true
	model.filterText = "test"
	model.filteredCategories = []data.Category{{CatID: "test"}}
	
	// Call resetEditingState
	resetModel := model.resetEditingState()
	
	// Verify everything is reset
	if resetModel.addCategory {
		t.Error("Expected addCategory to be false")
	}
	
	if resetModel.moveCategory {
		t.Error("Expected moveCategory to be false")
	}
	
	if resetModel.selectedGroup.GroupID != "" {
		t.Error("Expected selectedGroup to be empty")
	}
	
	if resetModel.movingCategory.CatID != "" {
		t.Error("Expected movingCategory to be empty")
	}
	
	if resetModel.isEditingName {
		t.Error("Expected isEditingName to be false")
	}
	
	if resetModel.editingIndex != -1 {
		t.Errorf("Expected editingIndex to be -1, got %d", resetModel.editingIndex)
	}
	
	if resetModel.editInput.Value() != "" {
		t.Errorf("Expected editInput value to be empty, got '%s'", resetModel.editInput.Value())
	}
	
	if resetModel.editInput.Focused() {
		t.Error("Expected editInput to be blurred")
	}
	
	if resetModel.isFiltered {
		t.Error("Expected isFiltered to be false")
	}
	
	if resetModel.filterText != "" {
		t.Error("Expected filterText to be empty")
	}
	
	if len(resetModel.filteredCategories) != 0 {
		t.Error("Expected filteredCategories to be empty")
	}
}

func TestCategoryGroupModel_ResetSelection_ResetsEditingState(t *testing.T) {
	model := createTestCategoryGroupModel()
	
	// Set up editing state
	model.selectGroup = true
	model.isEditingName = true
	model.editingIndex = 3
	model.editInput = textinput.New()
	model.editInput.Focus()
	model.editInput.SetValue("Test Group Name")
	
	// Call ResetSelection
	resetModel := model.ResetSelection()
	
	// Verify all state is reset
	if resetModel.selectGroup {
		t.Error("Expected selectGroup to be false after ResetSelection")
	}
	
	if resetModel.isEditingName {
		t.Error("Expected isEditingName to be false after ResetSelection")
	}
	
	if resetModel.editingIndex != -1 {
		t.Errorf("Expected editingIndex to be -1, got %d", resetModel.editingIndex)
	}
	
	if resetModel.editInput.Value() != "" {
		t.Errorf("Expected editInput value to be empty, got '%s'", resetModel.editInput.Value())
	}
	
	if resetModel.editInput.Focused() {
		t.Error("Expected editInput to be blurred after ResetSelection")
	}
}

func TestCategoryGroupModel_UpdateData_ResetsEditingState(t *testing.T) {
	model := createTestCategoryGroupModel()
	
	// Set up editing state
	model.isEditingName = true
	model.editingIndex = 2
	model.editInput.SetValue("Edit Value")
	model.selectGroup = true
	
	// Create new data
	newData := &data.DataRoot{
		CategoryGroups: map[string]data.CategoryGroup{
			"newgroup": {
				GroupID:   "newgroup",
				GroupName: "New Group",
				Order:     1,
			},
		},
		MonthlyData: map[string]data.MonthlyRecord{},
	}
	
	// Call UpdateData
	updatedModel := model.UpdateData(newData)
	
	// Verify editing state is reset
	if updatedModel.isEditingName {
		t.Error("Expected isEditingName to be false after UpdateData")
	}
	
	if updatedModel.editingIndex != -1 {
		t.Errorf("Expected editingIndex to be -1, got %d", updatedModel.editingIndex)
	}
	
	if updatedModel.editInput.Value() != "" {
		t.Error("Expected editInput value to be empty after UpdateData")
	}
	
	if updatedModel.selectGroup {
		t.Error("Expected selectGroup to be false after UpdateData")
	}
}

func TestCategoryGroupModel_ResetEditingState(t *testing.T) {
	model := createTestCategoryGroupModel()
	
	// Set up all editing states
	model.selectGroup = true
	model.isEditingName = true
	model.editingIndex = 7
	model.editInput.Focus()
	model.editInput.SetValue("Group Name")
	
	// Call resetEditingState
	resetModel := model.resetEditingState()
	
	// Verify everything is reset
	if resetModel.selectGroup {
		t.Error("Expected selectGroup to be false")
	}
	
	if resetModel.isEditingName {
		t.Error("Expected isEditingName to be false")
	}
	
	if resetModel.editingIndex != -1 {
		t.Errorf("Expected editingIndex to be -1, got %d", resetModel.editingIndex)
	}
	
	if resetModel.editInput.Value() != "" {
		t.Errorf("Expected editInput value to be empty, got '%s'", resetModel.editInput.Value())
	}
	
	if resetModel.editInput.Focused() {
		t.Error("Expected editInput to be blurred")
	}
}