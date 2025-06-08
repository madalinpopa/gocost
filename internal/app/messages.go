package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/data"
	"github.com/madalinpopa/gocost/internal/ui"
)

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

	return m, cmds
}

// handleViewErrorMsg handles the display of error messages.
func (m App) handleViewErrorMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(ui.ViewErrorMsg); ok {
		return m.SetErrorStatus(msg.Text)
	}
	return m, nil
}

// handleMonthlyViewMsg switches the active view to the monthly overview and updates the MonthlyModel
// with the current month and year, if it exists.
func (m App) handleMonthlyViewMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(ui.MonthlyViewMsg); ok {
		m.MonthlyModel = m.MonthlyModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
		m.activeView = viewMonthlyOverview
	}
	return m, nil
}

// handleGroupAddMsg handles the addition of a new category group. It updates the data model,
func (m App) handleGroupAddMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(ui.GroupAddMsg); ok {

		m.Data.CategoryGroups[msg.Group.GroupID] = msg.Group

		if err := data.SaveData(m.FilePath, m.Data); err != nil {
			return m.SetErrorStatus(fmt.Sprintf("Error while saving data: %v", err))
		} else {
			m.CategoryGroupModel = m.CategoryGroupModel.UpdateData(m.Data)
			return m.SetSuccessStatus(fmt.Sprintf("Group '%s' added successfully", msg.Group.GroupName))
		}
	}
	return m, nil
}

// handleGroupDeleteMsg handles the deletion of a category group.
func (m App) handleGroupDeleteMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(ui.GroupDeleteMsg); ok {
		canDelete := true

		if !canDelete {
			m.CategoryGroupModel = m.CategoryGroupModel.UpdateData(m.Data)
			return m.SetErrorStatus(fmt.Sprintf("Cannot delete group '%s': contains categories", msg.Group.GroupName))
		}

		if canDelete {
			delete(m.Data.CategoryGroups, msg.Group.GroupID)

			if err := data.SaveData(m.FilePath, m.Data); err != nil {
				return m.SetErrorStatus(fmt.Sprintf("Error while saving data: %v", err))
			} else {
				m.CategoryGroupModel = m.CategoryGroupModel.UpdateData(m.Data)
				return m.SetSuccessStatus(fmt.Sprintf("Group '%s' deleted successfully", msg.Group.GroupName))
			}
		}
		m.CategoryGroupModel = m.CategoryGroupModel.UpdateData(m.Data)
	}
	return m, nil
}

// handleGroupUpdateMsg handles the update of a category group.
func (m App) handleGroupUpdateMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	found := false
	var groupName string
	if msg, ok := msg.(ui.GroupUpdateMsg); ok {
		groupName = msg.Group.GroupName
		for i, group := range m.Data.CategoryGroups {

			if group.GroupID == msg.Group.GroupID {
				found = true
				m.Data.CategoryGroups[i] = msg.Group
				break
			}

		}
	}
	if found {
		if err := data.SaveData(m.FilePath, m.Data); err != nil {
			return m.SetErrorStatus(fmt.Sprintf("Error while saving data: %v", err))
		} else {
			m.CategoryGroupModel = m.CategoryGroupModel.UpdateData(m.Data)
			return m.SetSuccessStatus(fmt.Sprintf("Group '%s' updated successfully", groupName))
		}
	}
	return m, nil
}

// handleAddIncomeFormMsg handles the display of the income form.
func (m App) handleAddIncomeFormMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(ui.AddIncomeFormMsg); ok {
		m.IncomeFormModel = ui.NewIncomeFormModel(m.CurrentMonth, m.CurrentYear, nil)
		m.activeView = viewIncomeForm
	}
	return m, nil
}

// handleIncomeViewMsg handles the display of income data.
func (m App) handleIncomeViewMsg(msg tea.Msg) (tea.Model, tea.Cmd) {

	if _, ok := msg.(ui.IncomeViewMsg); ok {
		m.activeView = viewIncome
	}

	return m, nil
}

// handleSaveIncomeMsg handles the saving of income data.
func (m App) handleSaveIncomeMsg(msg tea.Msg) (tea.Model, tea.Cmd) {

	if msg, ok := msg.(ui.SaveIncomeMsg); ok {

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

		// save data to file
		err := data.SaveData(m.FilePath, m.Data)
		if err != nil {
			return m.SetErrorStatus("Failed to save income")
		} else {
			successMsg := "Income was added"
			if found {
				successMsg = "Income was updated"
			}
			return m.SetSuccessStatus(successMsg)
		}

	}
	return m, nil
}

// handleEditIncomeMsg handles the editing of income data.
func (m App) handleEditIncomeMsg(msg tea.Msg) (tea.Model, tea.Cmd) {

	if msg, ok := msg.(ui.EditIncomeMsg); ok {
		m.IncomeFormModel = ui.NewIncomeFormModel(m.CurrentMonth, m.CurrentYear, &msg.Income)
		m.activeView = viewIncomeForm
	}
	return m, nil
}

// handleDeleteIncomeMsg handles the deletion of income data.
func (m App) handleDeleteIncomeMsg(msg tea.Msg) (tea.Model, tea.Cmd) {

	if msg, ok := msg.(ui.DeleteIncomeMsg); ok {

		monthRecord, ok := m.Data.MonthlyData[msg.MonthKey]
		if ok {

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
			} else {
				return m.SetSuccessStatus("Income was deleted")
			}
		}

	}

	return m, nil
}

func (m App) handleSelectGroupMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(ui.SelectGroupMsg); ok {
		m.CategoryGroupModel = m.CategoryGroupModel.SelectGroup()
		m.activeView = viewCategoryGroup
	}
	return m, nil
}

func (m App) handleSelectedGroupMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if msg, ok := msg.(ui.SelectedGroupMsg); ok {
		m.CategoryModel, cmd = m.CategoryModel.AddCategory(msg.Group)
		m.activeView = viewCategory
	}
	return m, cmd
}

func (m App) handleCategoryAddMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(ui.CategoryAddMsg); ok {

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

		// Update model
		m.CategoryModel = m.CategoryModel.UpdateData(m.Data)

		return m.SetSuccessStatus("Category name was saved")
	}

	return m, nil
}

func (m App) handleCategoryUpdateMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(ui.CategoryUpdateMsg); ok {
		_ = msg

	}
	return m, nil
}
