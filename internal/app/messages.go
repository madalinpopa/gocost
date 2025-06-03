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
	if msg, ok := msg.(ui.MonthlyViewMsg); ok {
		m.activeView = viewMonthlyOverview
		if m.MonthlyModel != nil {
			updatedModel := m.MonthlyModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
			m.MonthlyModel = &updatedModel
		}
		fmt.Println(msg)
	}
	return m, nil
}

// handleGroupAddMsg handles the addition of a new category group. It updates the data model,
func (m App) handleGroupAddMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(ui.GroupAddMsg); ok {

		m.Data.CategoryGroups = append(m.Data.CategoryGroups, msg.Group)

		if err := data.SaveData(m.FilePath, m.Data); err != nil {
			fmt.Printf("Error while saving data: %v", err)
		} else {
			updatedModel := m.CategoryGroupModel.UpdateData(m.Data)
			m.CategoryGroupModel = &updatedModel
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

				// Check if group has any category. If so, you need to delete first
				// the category before deleting the group.
				if len(group.Categories) > 0 {
					canDelete = false
					// TODO: Need to set status message here
					fmt.Printf("Cannot delete group '%s': contains categories.", group.GroupName)
					break
				}
				break
			}
		}
		if canDelete && groupIndexToDelete != -1 {
			m.Data.CategoryGroups = slices.Delete(m.Data.CategoryGroups, groupIndexToDelete, groupIndexToDelete+1)
			if err := data.SaveData(m.FilePath, m.Data); err != nil {
				// TODO: Need to set status message here
			} else {
				// TODO: Need to set status message here
				fmt.Println("Delete group: ", groupName)
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
	if msg, ok := msg.(ui.GroupUpdateMsg); ok {
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
			// Set status message
		} else {
			// Set status message
		}

		if m.CategoryGroupModel != nil {
			updatedModel := m.CategoryGroupModel.UpdateData(m.Data)
			m.CategoryGroupModel = &updatedModel
		}
	}
	return m, nil
}

// handleGroupManageCategoriesMsg handles the management of categories within a group.
func (m App) handleGroupManageCategoriesMsg() (tea.Model, tea.Cmd) {
	return m, nil
}
