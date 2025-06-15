package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/domain"
	"github.com/madalinpopa/gocost/internal/ui"
)

// handlePopulateCategoriesMsg copy categories from the previous month if it exists.
func (m App) handlePopulateCategoriesMsg(msg ui.PopulateCategoriesMsg) (tea.Model, tea.Cmd) {
	count, err := m.categorySvc.CopyCategoriesFromMonth(msg.PreviousMonthKey, msg.CurrentMonthKey)
	if err != nil {
		return m.SetErrorStatus(fmt.Sprintf("Failed to copy categories: %v", err))
	}
	if count == 0 {
		return m.SetErrorStatus(fmt.Sprintf("No categories found in %s to copy from", msg.PreviousMonthKey))
	}

	app := m.refreshDataForModels()
	app.MonthlyModel = app.MonthlyModel.ResetFocus()

	return app.SetSuccessStatus(fmt.Sprintf("Successfully copied %d categories from %s", count, msg.PreviousMonthKey))
}

// handleModelsWindowResize updates the width and height within views.
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

// handleExpenseViewMsg handles the display of the expense form view.
func (m App) handleExpenseViewMsg(msg ui.ExpenseViewMsg) (tea.Model, tea.Cmd) {
	m.activeView = viewExpense
	m.ExpenseModel = ui.NewExpenseModel(msg.Category, msg.MonthKey)
	return m, m.ExpenseModel.Init()
}

// handleSaveExpenseMsg handles the saving of expense data.
func (m App) handleSaveExpenseMsg(msg ui.SaveExpenseMsg) (tea.Model, tea.Cmd) {
	category := msg.Category
	if category.Expense == nil {
		category.Expense = make(map[string]domain.ExpenseRecord)
	}
	category.Expense[category.CatID] = msg.Expense

	err := m.categorySvc.UpdateCategory(msg.MonthKey, category)
	if err != nil {
		return m.SetErrorStatus(fmt.Sprintf("Failed to save expense: %v", err))
	}

	app := m.refreshDataForModels()
	app.MonthlyModel = app.MonthlyModel.SetFocusToCategory(msg.Category)
	app.activeView = viewMonthlyOverview
	return app.SetSuccessStatus(fmt.Sprintf("Expense for '%s' saved successfully", msg.Category.CategoryName))
}

// handleEditExpenseMsg handles edit expense message.
func (m App) handleEditExpenseMsg(msg ui.EditExpenseMsg) (tea.Model, tea.Cmd) {
	m.activeView = viewExpense
	m.ExpenseModel = ui.NewExpenseModel(msg.Category, msg.MonthKey)
	return m, m.ExpenseModel.Init()
}

// handleDeleteExpenseMsg clears the expense from the category.
func (m App) handleDeleteExpenseMsg(msg ui.DeleteExpenseMsg) (tea.Model, tea.Cmd) {
	category := msg.Category
	category.Expense[category.CatID] = domain.ExpenseRecord{
		Amount: 0,
		Budget: 0,
		Status: "Not Paid",
		Notes:  "",
	}

	err := m.categorySvc.UpdateCategory(msg.MonthKey, category)
	if err != nil {
		return m.SetErrorStatus(fmt.Sprintf("Failed to clear expense: %v", err))
	}

	app := m.refreshDataForModels()
	app.MonthlyModel = app.MonthlyModel.SetFocusToCategory(msg.Category)
	app.activeView = viewMonthlyOverview
	return app.SetSuccessStatus(fmt.Sprintf("Expense for category '%s' has been cleared", msg.Category.CategoryName))
}

// handleToggleExpenseStatusMsg toggles the status of an expense.
func (m App) handleToggleExpenseStatusMsg(msg ui.ToggleExpenseStatusMsg) (tea.Model, tea.Cmd) {
	category := msg.Category
	if category.Expense == nil {
		category.Expense = make(map[string]domain.ExpenseRecord)
	}

	currentExpense := category.Expense[category.CatID]
	if currentExpense.Status == "Paid" {
		currentExpense.Status = "Not Paid"
	} else {
		currentExpense.Status = "Paid"
	}
	category.Expense[category.CatID] = currentExpense

	err := m.categorySvc.UpdateCategory(msg.MonthKey, category)
	if err != nil {
		return m.SetErrorStatus(fmt.Sprintf("Failed to toggle status: %v", err))
	}

	app := m.refreshDataForModels()
	app.MonthlyModel = app.MonthlyModel.SetFocusToCategory(category)
	return app.SetSuccessStatus(fmt.Sprintf("Status for '%s' toggled to '%s'", category.CategoryName, currentExpense.Status))
}

// handleMonthlyViewMsg switches the active view to the monthly overview.
func (m App) handleMonthlyViewMsg() (tea.Model, tea.Cmd) {
	m.activeView = viewMonthlyOverview
	app := m.refreshDataForModels()
	return app, nil
}

// handleGroupAddMsg handles the addition of a new category group.
func (m App) handleGroupAddMsg(msg ui.GroupAddMsg) (tea.Model, tea.Cmd) {
	err := m.groupSvc.AddGroup(msg.Group)
	if err != nil {
		return m.SetErrorStatus(fmt.Sprintf("Failed to add group: %v", err))
	}
	app := m.refreshDataForModels()
	return app.SetSuccessStatus(fmt.Sprintf("Group '%s' added successfully", msg.Group.GroupName))
}

// handleGroupDeleteMsg handles the deletion of a category group.
func (m App) handleGroupDeleteMsg(msg ui.GroupDeleteMsg) (tea.Model, tea.Cmd) {
	err := m.groupSvc.DeleteGroup(msg.Group.GroupID)
	if err != nil {
		return m.SetErrorStatus(fmt.Sprintf("Failed to delete group: %v", err))
	}
	app := m.refreshDataForModels()
	return app.SetSuccessStatus(fmt.Sprintf("Group '%s' deleted successfully", msg.Group.GroupName))
}

// handleGroupUpdateMsg handles the update of a category group.
func (m App) handleGroupUpdateMsg(msg ui.GroupUpdateMsg) (tea.Model, tea.Cmd) {
	err := m.groupSvc.UpdateGroup(msg.Group)
	if err != nil {
		return m.SetErrorStatus(fmt.Sprintf("Failed to update group: %v", err))
	}
	app := m.refreshDataForModels()
	return app.SetSuccessStatus(fmt.Sprintf("Group '%s' updated successfully", msg.Group.GroupName))
}

// handleAddIncomeFormMsg handles the display of the income form.
func (m App) handleAddIncomeFormMsg() (tea.Model, tea.Cmd) {
	m.IncomeFormModel = ui.NewIncomeFormModel(m.CurrentMonth, m.CurrentYear, nil)
	m.activeView = viewIncomeForm
	return m, nil
}

// handleIncomeViewMsg handles the display of income data.
func (m App) handleIncomeViewMsg() (tea.Model, tea.Cmd) {
	app := m.refreshDataForModels()
	app.activeView = viewIncome
	return app, nil
}

// handleSaveIncomeMsg handles the saving of income data.
func (m App) handleSaveIncomeMsg(msg ui.SaveIncomeMsg) (tea.Model, tea.Cmd) {
	existingIncomes, _ := m.incomeSvc.GetIncomesForMonth(msg.MonthKey)
	isUpdate := false
	for _, income := range existingIncomes {
		if income.IncomeID == msg.Income.IncomeID {
			isUpdate = true
			break
		}
	}

	var err error
	var successMsg string
	if isUpdate {
		err = m.incomeSvc.UpdateIncome(msg.MonthKey, msg.Income)
		successMsg = fmt.Sprintf("Income '%s' updated successfully", msg.Income.Description)
	} else {
		err = m.incomeSvc.AddIncome(msg.MonthKey, msg.Income)
		successMsg = fmt.Sprintf("Income '%s' added successfully", msg.Income.Description)
	}

	if err != nil {
		return m.SetErrorStatus(fmt.Sprintf("Failed to save income: %v", err))
	}

	app := m.refreshDataForModels()
	app.activeView = viewIncome
	return app.SetSuccessStatus(successMsg)
}

// handleEditIncomeMsg handles the editing of income data.
func (m App) handleEditIncomeMsg(msg ui.EditIncomeMsg) (tea.Model, tea.Cmd) {
	m.IncomeFormModel = ui.NewIncomeFormModel(m.CurrentMonth, m.CurrentYear, &msg.Income)
	m.activeView = viewIncomeForm
	return m, nil
}

// handleDeleteIncomeMsg handles the deletion of income data.
func (m App) handleDeleteIncomeMsg(msg ui.DeleteIncomeMsg) (tea.Model, tea.Cmd) {
	err := m.incomeSvc.DeleteIncome(msg.MonthKey, msg.Income.IncomeID)
	if err != nil {
		return m.SetErrorStatus(fmt.Sprintf("Failed to delete income: %v", err))
	}
	app := m.refreshDataForModels()
	return app.SetSuccessStatus(fmt.Sprintf("Income '%s' has been deleted", msg.Income.Description))
}

// handleSelectGroupMsg handles the selection of a category group.
func (m App) handleSelectGroupMsg() (tea.Model, tea.Cmd) {
	m.CategoryGroupModel = m.CategoryGroupModel.SelectGroup()
	m.activeView = viewCategoryGroup
	return m, nil
}

// handleSelectedGroupMsg handles the selected category group and returns to Category view.
func (m App) handleSelectedGroupMsg(msg ui.SelectedGroupMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var updatedCatModel tea.Model

	if m.CategoryModel.IsMovingCategory() {
		updatedCatModel, cmd = m.CategoryModel.MoveCategory(msg.Group)
		m.CategoryModel, _ = updatedCatModel.(ui.CategoryModel)
	} else {
		updatedCatModel, cmd = m.CategoryModel.AddCategory(msg.Group)
		m.CategoryModel, _ = updatedCatModel.(ui.CategoryModel)
	}

	m.activeView = viewCategory
	return m, cmd
}

// handleCategoryAddMsg handles the addition of a new category.
func (m App) handleCategoryAddMsg(msg ui.CategoryAddMsg) (tea.Model, tea.Cmd) {
	err := m.categorySvc.AddCategory(msg.MonthKey, msg.Category)
	if err != nil {
		return m.SetErrorStatus(fmt.Sprintf("Failed to add category: %v", err))
	}
	app := m.refreshDataForModels()
	return app.SetSuccessStatus(fmt.Sprintf("Category '%s' has been created successfully", msg.Category.CategoryName))
}

// handleCategoryUpdateMsg handles the update of a category.
func (m App) handleCategoryUpdateMsg(msg ui.CategoryUpdateMsg) (tea.Model, tea.Cmd) {
	err := m.categorySvc.UpdateCategory(msg.MonthKey, msg.Category)
	if err != nil {
		return m.SetErrorStatus(fmt.Sprintf("Failed to update category: %v", err))
	}
	app := m.refreshDataForModels()
	// After moving a category, reset the state in the UI model
	if app.CategoryModel.IsMovingCategory() {
		app.CategoryModel = app.CategoryModel.ResetMoveState()
	}
	return app.SetSuccessStatus(fmt.Sprintf("Category '%s' has been updated successfully", msg.Category.CategoryName))
}

// handleCategoryDeleteMsg handles the deletion of a category.
func (m App) handleCategoryDeleteMsg(msg ui.CategoryDeleteMsg) (tea.Model, tea.Cmd) {
	err := m.categorySvc.DeleteCategory(msg.MonthKey, msg.Category.CatID)
	if err != nil {
		return m.SetErrorStatus(fmt.Sprintf("Failed to delete category: %v", err))
	}
	app := m.refreshDataForModels()
	return app.SetSuccessStatus(fmt.Sprintf("Category '%s' has been deleted", msg.Category.CategoryName))
}

// handleFilterCategoriesMsg handles filtering categories by the provided filter text.
func (m App) handleFilterCategoriesMsg(msg ui.FilterCategoriesMsg) (tea.Model, tea.Cmd) {
	updatedCategoryModel, cmd := m.CategoryModel.Update(msg)
	if cgMo, ok := updatedCategoryModel.(ui.CategoryModel); ok {
		m.CategoryModel = cgMo
	}
	return m, cmd
}

// handleReturnToMonthlyWithFocusMsg handles the return to monthly view with focus on a specific category.
func (m App) handleReturnToMonthlyWithFocusMsg(msg ui.ReturnToMonthlyWithFocusMsg) (tea.Model, tea.Cmd) {
	m = m.refreshDataForModels()
	m.MonthlyModel = m.MonthlyModel.SetFocusToCategory(msg.Category)
	m.activeView = viewMonthlyOverview
	return m, nil
}

// handleCategoryViewMsg handles the return to category view.
func (m App) handleCategoryViewMsg() (tea.Model, tea.Cmd) {
	app := m.refreshDataForModels()
	app.activeView = viewCategory
	return app, nil
}
