package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/data"
	"github.com/madalinpopa/gocost/internal/ui"
)

// handlePopulateCategoriesMsg copy categories from the previous month if exists
func (m App) handlePopulateCategoriesMsg(msg ui.PopulateCategoriesMsg) (App, tea.Cmd) {
	prevRecord, exists := m.Data.MonthlyData[msg.PreviousMonthKey]
	if !exists || len(prevRecord.Categories) == 0 {
		return m.SetErrorStatus("No categories found in previous month")
	}

	var newCategories []data.Category
	for _, category := range prevRecord.Categories {
		newCategory := data.Category{
			CatID:        category.CatID,
			GroupID:      category.GroupID,
			CategoryName: category.CategoryName,
			Expense:      make(map[string]data.ExpenseRecord), // Reset expenses
		}
		newCategories = append(newCategories, newCategory)
	}

	// Create or update current month record
	currentRecord := data.MonthlyRecord{
		Incomes:    []data.IncomeRecord{}, // Start with empty incomes
		Categories: newCategories,
	}

	// If current month record exists, preserve incomes
	if existingRecord, exists := m.Data.MonthlyData[msg.CurrentMonthKey]; exists {
		currentRecord.Incomes = existingRecord.Incomes
	}

	m.Data.MonthlyData[msg.CurrentMonthKey] = currentRecord

	// Reset focus to first group in monthly model
	m.MonthlyModel = m.MonthlyModel.ResetFocus()

	return m.SetSuccessStatus("Categories populated from previous month")
}

// handleModelsWindowResize updates the width and height within views
func (m App) handleModelsWindowResize(msg tea.Msg) (tea.Model, []tea.Cmd) {
	var cmds []tea.Cmd

	updatedMonthlyModel, moCmd := m.MonthlyModel.Update(msg)
	if mo, ok := updatedMonthlyModel.(ui.MonthlyModel); ok {
		m.MonthlyModel = mo
	}
	cmds = append(cmds, moCmd)

	updatedIncomeModel, moCmd := m.IncomeModel.Update(msg)
	if mo, ok := updatedIncomeModel.(ui.IncomeModel); ok {
		m.IncomeModel = mo
	}
	cmds = append(cmds, moCmd)

	updatedCategoryGroupModel, cgCmd := m.CategoryGroupModel.Update(msg)
	if cgMo, ok := updatedCategoryGroupModel.(ui.CategoryGroupModel); ok {
		m.CategoryGroupModel = cgMo
	}
	cmds = append(cmds, cgCmd)

	updatedCategoryModel, cgCmd := m.CategoryModel.Update(msg)
	if cgMo, ok := updatedCategoryModel.(ui.CategoryModel); ok {
		m.CategoryModel = cgMo
	}
	cmds = append(cmds, cgCmd)

	updatedExpenseModel, expCmd := m.ExpenseModel.Update(msg)
	if expMo, ok := updatedExpenseModel.(ui.ExpenseModel); ok {
		m.ExpenseModel = expMo
	}
	cmds = append(cmds, expCmd)

	return m, cmds
}

// handleViewErrorMsg handles the display of error messages.
func (m App) handleViewErrorMsg(msg ui.ViewErrorMsg) (tea.Model, tea.Cmd) {
	return m.SetErrorStatus(msg.Text)
}

// handleExpenseViewMsg handles the display o expense form view
func (m App) handleExpenseViewMsg(msg ui.ExpenseViewMsg) (tea.Model, tea.Cmd) {
	m.activeView = viewExpense
	m.ExpenseModel = ui.NewExpenseModel(m.Data, msg.Category, msg.MonthKey)
	return m, m.ExpenseModel.Init()
}

// handleSaveExpenseMsg handles the saving of expense data
func (m App) handleSaveExpenseMsg(msg ui.SaveExpenseMsg) (tea.Model, tea.Cmd) {
	monthRecord, exists := m.Data.MonthlyData[msg.MonthKey]
	if !exists {
		monthRecord = data.MonthlyRecord{
			Categories: []data.Category{},
		}
	}

	// Find and update the category
	for i, category := range monthRecord.Categories {
		if category.CatID == msg.Category.CatID {
			if category.Expense == nil {
				category.Expense = make(map[string]data.ExpenseRecord)
			}
			category.Expense[category.CatID] = msg.Expense
			monthRecord.Categories[i] = category
			break
		}
	}

	m.Data.MonthlyData[msg.MonthKey] = monthRecord

	if err := data.SaveData(m.FilePath, m.Data); err != nil {
		return m.SetErrorStatus("Failed to save expense")
	}

	// Update models with new data
	m.MonthlyModel = m.MonthlyModel.UpdateData(m.Data)
	m.IncomeModel = m.IncomeModel.UpdateData(m.Data)
	m.CategoryModel = m.CategoryModel.UpdateData(m.Data)

	// Set focus to the category that was just updated
	m.MonthlyModel = m.MonthlyModel.SetFocusToCategory(msg.Category)

	// Return to monthly view
	m.activeView = viewMonthlyOverview
	return m.SetSuccessStatus("Expense saved successfully")
}

// handleEditExpenseMsg handles edit expense message
func (m App) handleEditExpenseMsg(msg ui.EditExpenseMsg) (tea.Model, tea.Cmd) {
	m.activeView = viewExpense
	m.ExpenseModel = ui.NewExpenseModel(m.Data, msg.Category, msg.MonthKey)
	return m, m.ExpenseModel.Init()
}

// handleDeleteExpenseMsg find and clear the expense from the category (reset to default values)
func (m App) handleDeleteExpenseMsg(msg ui.DeleteExpenseMsg) (tea.Model, tea.Cmd) {
	monthRecord, exists := m.Data.MonthlyData[msg.MonthKey]
	if exists {
		for i, category := range monthRecord.Categories {
			if category.CatID == msg.Category.CatID {
				// Instead of deleting, reset the expense to default values
				if category.Expense == nil {
					category.Expense = make(map[string]data.ExpenseRecord)
				}
				category.Expense[category.CatID] = data.ExpenseRecord{
					Amount: 0,
					Budget: 0,
					Status: "Not Paid",
					Notes:  "",
				}
				monthRecord.Categories[i] = category
				break
			}
		}
		m.Data.MonthlyData[msg.MonthKey] = monthRecord

		if err := data.SaveData(m.FilePath, m.Data); err != nil {
			return m.SetErrorStatus("Failed to clear expense")
		}
	}

	// Update models with new data
	m.MonthlyModel = m.MonthlyModel.UpdateData(m.Data)
	m.IncomeModel = m.IncomeModel.UpdateData(m.Data)
	m.CategoryModel = m.CategoryModel.UpdateData(m.Data)

	// Set focus to the category that was just cleared
	m.MonthlyModel = m.MonthlyModel.SetFocusToCategory(msg.Category)

	// Return to monthly view
	m.activeView = viewMonthlyOverview
	return m.SetSuccessStatus("Expense cleared successfully")
}

// handleMonthlyViewMsg switches the active view to the monthly overview and updates the MonthlyModel
// with the current month and year, if it exists.
func (m App) handleMonthlyViewMsg(msg ui.MonthlyViewMsg) (tea.Model, tea.Cmd) {
	m.MonthlyModel = m.MonthlyModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
	m.CategoryModel = m.CategoryModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
	m.activeView = viewMonthlyOverview
	return m, nil
}

// handleGroupAddMsg handles the addition of a new category group. It updates the data model,
func (m App) handleGroupAddMsg(msg ui.GroupAddMsg) (tea.Model, tea.Cmd) {

	m.Data.CategoryGroups[msg.Group.GroupID] = msg.Group

	if err := data.SaveData(m.FilePath, m.Data); err != nil {
		return m.SetErrorStatus(fmt.Sprintf("Error while saving data: %v", err))
	}
	m.CategoryGroupModel = m.CategoryGroupModel.UpdateData(m.Data)
	return m.SetSuccessStatus(fmt.Sprintf("Group '%s' added successfully", msg.Group.GroupName))
}

// handleGroupDeleteMsg handles the deletion of a category group.
func (m App) handleGroupDeleteMsg(msg ui.GroupDeleteMsg) (tea.Model, tea.Cmd) {
	canDelete := true

	// Check if any categories are using this group
	for _, monthRecord := range m.Data.MonthlyData {
		for _, category := range monthRecord.Categories {
			if category.GroupID == msg.Group.GroupID {
				canDelete = false
				break
			}
		}
		if !canDelete {
			break
		}
	}

	if !canDelete {
		m.CategoryGroupModel = m.CategoryGroupModel.UpdateData(m.Data)
		return m.SetErrorStatus(fmt.Sprintf("Cannot delete group '%s': contains categories", msg.Group.GroupName))
	}

	if canDelete {
		delete(m.Data.CategoryGroups, msg.Group.GroupID)

		if err := data.SaveData(m.FilePath, m.Data); err != nil {
			return m.SetErrorStatus(fmt.Sprintf("Error while saving data: %v", err))
		}

		m.CategoryGroupModel = m.CategoryGroupModel.UpdateData(m.Data)
		return m.SetSuccessStatus(fmt.Sprintf("Group '%s' deleted successfully", msg.Group.GroupName))
	}
	return m, nil
}

// handleGroupUpdateMsg handles the update of a category group.
func (m App) handleGroupUpdateMsg(msg ui.GroupUpdateMsg) (tea.Model, tea.Cmd) {
	groupId := msg.Group.GroupID
	groupName := msg.Group.GroupName

	_, existing := m.Data.CategoryGroups[groupId]
	if !existing {
		return m.SetErrorStatus(fmt.Sprintf("Failed to update group. Group not found: %s", groupName))
	}

	m.Data.CategoryGroups[groupId] = msg.Group

	if err := data.SaveData(m.FilePath, m.Data); err != nil {
		return m.SetErrorStatus(fmt.Sprintf("Error while saving data: %v", err))
	}
	m.CategoryGroupModel = m.CategoryGroupModel.UpdateData(m.Data)
	return m.SetSuccessStatus(fmt.Sprintf("Group '%s' updated successfully", groupName))
}

// handleAddIncomeFormMsg handles the display of the income form.
func (m App) handleAddIncomeFormMsg(msg ui.AddIncomeFormMsg) (tea.Model, tea.Cmd) {
	m.IncomeFormModel = ui.NewIncomeFormModel(m.CurrentMonth, m.CurrentYear, nil)
	m.activeView = viewIncomeForm
	return m, nil
}

// handleIncomeViewMsg handles the display of income data.
func (m App) handleIncomeViewMsg(msg ui.IncomeViewMsg) (tea.Model, tea.Cmd) {
	m.IncomeModel = ui.NewIncomeModel(m.Data, m.CurrentMonth, m.CurrentYear)
	m.activeView = viewIncome
	return m, nil
}

// handleSaveIncomeMsg handles the saving of income data.
func (m App) handleSaveIncomeMsg(msg ui.SaveIncomeMsg) (tea.Model, tea.Cmd) {

	// Get month record, if not exists, create a new one with the income record
	monthRecord, ok := m.Data.MonthlyData[msg.MonthKey]
	if !ok {
		monthRecord = data.MonthlyRecord{
			Incomes:    make([]data.IncomeRecord, 0),
			Categories: make([]data.Category, 0),
		}
	}

	// check if income exists
	found := false
	for i, income := range monthRecord.Incomes {
		if income.IncomeID == msg.Income.IncomeID {
			monthRecord.Incomes[i] = msg.Income
			found = true
			break
		}
	}

	if !found {
		monthRecord.Incomes = append(monthRecord.Incomes, msg.Income)
	}

	m.Data.MonthlyData[msg.MonthKey] = monthRecord
	m.IncomeModel = ui.NewIncomeModel(m.Data, m.CurrentMonth, m.CurrentYear)
	m.activeView = viewIncome

	err := data.SaveData(m.FilePath, m.Data)
	if err != nil {
		return m.SetErrorStatus("Failed to save income")
	}
	successMsg := "Income was added"
	if found {
		successMsg = "Income was updated"
	}
	return m.SetSuccessStatus(successMsg)
}

// handleEditIncomeMsg handles the editing of income data.
func (m App) handleEditIncomeMsg(msg ui.EditIncomeMsg) (tea.Model, tea.Cmd) {
	m.IncomeFormModel = ui.NewIncomeFormModel(m.CurrentMonth, m.CurrentYear, &msg.Income)
	m.activeView = viewIncomeForm
	return m, nil
}

// handleDeleteIncomeMsg handles the deletion of income data.
func (m App) handleDeleteIncomeMsg(msg ui.DeleteIncomeMsg) (tea.Model, tea.Cmd) {
	if monthRecord, ok := m.Data.MonthlyData[msg.MonthKey]; ok {
		var updatedIncomes []data.IncomeRecord

		for _, income := range monthRecord.Incomes {
			if income.IncomeID != msg.Income.IncomeID {
				updatedIncomes = append(updatedIncomes, income)
			}
		}

		monthRecord.Incomes = updatedIncomes
		m.Data.MonthlyData[msg.MonthKey] = monthRecord
		m.IncomeModel = ui.NewIncomeModel(m.Data, m.CurrentMonth, m.CurrentYear)

		err := data.SaveData(m.FilePath, m.Data)
		if err != nil {
			return m.SetErrorStatus("Failed to delete income")
		}
		return m.SetSuccessStatus("Income was deleted")
	}
	return m, nil
}

// handleSelectGroupMsg handles the selection of a category group.
func (m App) handleSelectGroupMsg(msg ui.SelectGroupMsg) (tea.Model, tea.Cmd) {
	m.CategoryGroupModel = m.CategoryGroupModel.SelectGroup()
	m.activeView = viewCategoryGroup
	return m, nil
}

// handleSelectedGroupMsg handles the selected category group and returns to Category view.
func (m App) handleSelectedGroupMsg(msg ui.SelectedGroupMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.CategoryModel.IsMovingCategory() {
		// Moving an existing category
		m.CategoryModel, cmd = m.CategoryModel.MoveCategory(msg.Group)
		m.CategoryModel = m.CategoryModel.ResetMoveState()
	} else {
		// Adding a new category
		m.CategoryModel, cmd = m.CategoryModel.AddCategory(msg.Group)
	}
	m.activeView = viewCategory
	return m, cmd
}

// handleCategoryAddMsg handles the addition of a new category.
func (m App) handleCategoryAddMsg(msg ui.CategoryAddMsg) (tea.Model, tea.Cmd) {
	monthRecord, exists := m.Data.MonthlyData[msg.MonthKey]
	if !exists {
		monthRecord = data.MonthlyRecord{
			Incomes:    make([]data.IncomeRecord, 0),
			Categories: make([]data.Category, 0),
		}
	}

	monthRecord.Categories = append(monthRecord.Categories, msg.Category)

	m.Data.MonthlyData[msg.MonthKey] = monthRecord

	if err := data.SaveData(m.FilePath, m.Data); err != nil {
		return m.SetErrorStatus("Failed to save category")
	}

	m.CategoryModel = m.CategoryModel.UpdateData(m.Data)
	m.MonthlyModel = m.MonthlyModel.UpdateData(m.Data)

	// Set focus to the newly added category
	m.MonthlyModel = m.MonthlyModel.SetFocusToCategory(msg.Category)

	return m.SetSuccessStatus("Category name was saved")
}

// handleCategoryUpdateMsg handles the update of a category.
func (m App) handleCategoryUpdateMsg(msg ui.CategoryUpdateMsg) (tea.Model, tea.Cmd) {
	monthRecord, exists := m.Data.MonthlyData[msg.MonthKey]
	if !exists {
		return m.SetErrorStatus("Failed to update category: month record not found")
	}

	// Find and update the category
	found := false
	for i, category := range monthRecord.Categories {
		if category.CatID == msg.Category.CatID {
			monthRecord.Categories[i] = msg.Category
			found = true
			break
		}
	}

	if !found {
		return m.SetErrorStatus("Failed to update category: category not found")
	}

	m.Data.MonthlyData[msg.MonthKey] = monthRecord

	if err := data.SaveData(m.FilePath, m.Data); err != nil {
		return m.SetErrorStatus("Failed to save category")
	}

	m.CategoryModel = m.CategoryModel.UpdateData(m.Data)
	m.MonthlyModel = m.MonthlyModel.UpdateData(m.Data)

	// Set focus to the updated category
	m.MonthlyModel = m.MonthlyModel.SetFocusToCategory(msg.Category)

	return m.SetSuccessStatus("Category was updated successfully")
}

// handleCategoryDeleteMsg handles the deletion of a category.
func (m App) handleCategoryDeleteMsg(msg ui.CategoryDeleteMsg) (tea.Model, tea.Cmd) {
	if monthRecord, ok := m.Data.MonthlyData[msg.MonthKey]; ok {
		var updatedCategories []data.Category

		for _, category := range monthRecord.Categories {
			if category.CatID != msg.Category.CatID {
				updatedCategories = append(updatedCategories, category)
			}
		}

		monthRecord.Categories = updatedCategories
		m.Data.MonthlyData[msg.MonthKey] = monthRecord
		m.CategoryModel = m.CategoryModel.UpdateData(m.Data)
		m.MonthlyModel = m.MonthlyModel.UpdateData(m.Data)

		// Reset focus since category was deleted
		m.MonthlyModel = m.MonthlyModel.ResetFocus()

		err := data.SaveData(m.FilePath, m.Data)
		if err != nil {
			return m.SetErrorStatus("Failed to delete category")
		}
		return m.SetSuccessStatus("Category was deleted")
	}

	return m, nil
}

// handleReturnToMonthlyWithFocusMsg handles the return to monthly view with focus on a specific category.
func (m App) handleReturnToMonthlyWithFocusMsg(msg ui.ReturnToMonthlyWithFocusMsg) (tea.Model, tea.Cmd) {
	m.MonthlyModel = m.MonthlyModel.SetFocusToCategory(msg.Category)
	m.activeView = viewMonthlyOverview
	return m, nil
}
