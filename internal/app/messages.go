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
func (m App) handleMonthlyViewMsg() (tea.Model, tea.Cmd) {
	m.activeView = viewMonthlyOverview
	if m.MonthlyModel != nil {
		updatedModel := m.MonthlyModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
		m.MonthlyModel = &updatedModel
	}
	return m, nil
}

// handleGroupAddMsg handles the addition of a new category group. It updates the data model,
func (m App) handleGroupAddMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(ui.GroupAddMsg); ok {

		m.Data.Root.CategoryGroups = append(m.Data.Root.CategoryGroups, msg.Group)

		if err := data.SaveData(m.Data.FilePath, m.Data.Root); err != nil {
			fmt.Printf("Error while saving data: %v", err)
		} else {
			updatedModel := m.CategoryGroupModel.UpdateData(m.Data.Root)
			m.CategoryGroupModel = &updatedModel
		}
	}
	return m, nil
}

// handleGroupDeleteMsg handles the deletion of a category group.
func (m App) handleGroupDeleteMsg() (tea.Model, tea.Cmd) {
	return m, nil
}

// handleGroupUpdateMsg handles the update of a category group.
func (m App) handleGroupUpdateMsg() (tea.Model, tea.Cmd) {
	return m, nil
}

// handleGroupManageCategoriesMsg handles the management of categories within a group.
func (m App) handleGroupManageCategoriesMsg() (tea.Model, tea.Cmd) {
	return m, nil
}
