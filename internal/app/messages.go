package app

import (
	"fmt"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/data"
	"github.com/madalinpopa/gocost/internal/ui"
)

// handleModelsWindowResize updates the width and height within views
func (m App) handleModelsWindowResize(msg tea.Msg) (tea.Model, []tea.Cmd) {
	var cmds []tea.Cmd

	if m.MonthlyModel != nil {
		updatedMonthlyModel, moCmd := m.MonthlyModel.Update(msg)
		if mo, ok := updatedMonthlyModel.(ui.MonthlyModel); ok {
			m.MonthlyModel = &mo
		}
		cmds = append(cmds, moCmd)
	}

	if m.IncomeModel != nil {
		updatedIncomeModel, moCmd := m.IncomeModel.Update(msg)
		if mo, ok := updatedIncomeModel.(ui.IncomeModel); ok {
			m.IncomeModel = &mo
		}
		cmds = append(cmds, moCmd)
	}

	if m.CategoryGroupModel != nil {
		updatedCategoryGroupModel, cgCmd := m.CategoryGroupModel.Update(msg)
		if cgMo, ok := updatedCategoryGroupModel.(ui.CategoryGroupModel); ok {
			m.CategoryGroupModel = &cgMo
		}
		cmds = append(cmds, cgCmd)
	}
	return m, cmds
}

// handleMonthlyViewMsg switches the active view to the monthly overview and updates the MonthlyModel
// with the current month and year, if it exists.
func (m App) handleMonthlyViewMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.activeView = viewMonthlyOverview
	if m.MonthlyModel != nil {
		updatedModel := m.MonthlyModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
		m.MonthlyModel = &updatedModel
	}
	fmt.Println(msg)
	return m, nil
}

// handleGroupAddMsg handles the addition of a new category group. It updates the data model,
func (m App) handleGroupAddMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(ui.GroupAddMsg); ok {

		m.Data.CategoryGroups = append(m.Data.CategoryGroups, msg.Group)

		if err := data.SaveData(m.FilePath, m.Data); err != nil {
			return m.SetErrorStatus(fmt.Sprintf("Error while saving data: %v", err))
		} else {
			updatedModel := m.CategoryGroupModel.UpdateData(m.Data)
			m.CategoryGroupModel = &updatedModel
			return m.SetSuccessStatus(fmt.Sprintf("Group '%s' added successfully", msg.Group.GroupName))
		}
	}
	return m, nil
}

// handleGroupDeleteMsg handles the deletion of a category group.
func (m App) handleGroupDeleteMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(ui.GroupDeleteMsg); ok {
		canDelete := true
		groupIndexToDelete := -1
		var groupName string

		for i, group := range m.Data.CategoryGroups {
			if group.GroupID == msg.GroupID {
				groupIndexToDelete = i
				groupName = group.GroupName

				// Check if group has any category. If so, you need to delete first
				// the category before deleting the group.
				if len(group.Categories) > 0 {
					canDelete = false
					break
				}
				break
			}
		}

		if !canDelete {
			if m.CategoryGroupModel != nil {
				updatedModel := m.CategoryGroupModel.UpdateData(m.Data)
				m.CategoryGroupModel = &updatedModel
			}
			return m.SetErrorStatus(fmt.Sprintf("Cannot delete group '%s': contains categories", groupName))
		}

		if canDelete && groupIndexToDelete != -1 {
			m.Data.CategoryGroups = slices.Delete(m.Data.CategoryGroups, groupIndexToDelete, groupIndexToDelete+1)
			if err := data.SaveData(m.FilePath, m.Data); err != nil {
				if m.CategoryGroupModel != nil {
					updatedModel := m.CategoryGroupModel.UpdateData(m.Data)
					m.CategoryGroupModel = &updatedModel
				}
				return m.SetErrorStatus(fmt.Sprintf("Error while saving data: %v", err))
			} else {
				if m.CategoryGroupModel != nil {
					updatedModel := m.CategoryGroupModel.UpdateData(m.Data)
					m.CategoryGroupModel = &updatedModel
				}
				return m.SetSuccessStatus(fmt.Sprintf("Group '%s' deleted successfully", groupName))
			}
		}

		if m.CategoryGroupModel != nil {
			updatedModel := m.CategoryGroupModel.UpdateData(m.Data)
			m.CategoryGroupModel = &updatedModel
		}
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
			if m.CategoryGroupModel != nil {
				updatedModel := m.CategoryGroupModel.UpdateData(m.Data)
				m.CategoryGroupModel = &updatedModel
			}
			return m.SetErrorStatus(fmt.Sprintf("Error while saving data: %v", err))
		} else {
			if m.CategoryGroupModel != nil {
				updatedModel := m.CategoryGroupModel.UpdateData(m.Data)
				m.CategoryGroupModel = &updatedModel
			}
			return m.SetSuccessStatus(fmt.Sprintf("Group '%s' updated successfully", groupName))
		}
	}
	return m, nil
}

// handleGroupManageCategoriesMsg handles the management of categories within a group.
func (m App) handleGroupManageCategoriesMsg() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m App) handleAddIncomeFormMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(ui.AddIncomeFormMsg); ok {
		fmt.Println("Heiiiiii")
		fmt.Println(msg)
	}
	return m, nil
}
