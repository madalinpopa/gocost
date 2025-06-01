package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/data"
	"github.com/madalinpopa/gocost/internal/ui"
)

func (m App) handleMonthlyViewMsg() (tea.Model, tea.Cmd) {
	m.activeView = viewMonthlyOverview
	if m.monthlyModel != nil {
		updatedModel := m.monthlyModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
		m.monthlyModel = &updatedModel
	}
	return m, nil
}

func (m App) handleGroupAddMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(ui.GroupAddMsg); ok {

		m.Data.Root.CategoryGroups = append(m.Data.Root.CategoryGroups, msg.Group)

		if err := data.SaveData(m.Data.FilePath, m.Data.Root); err != nil {
			fmt.Printf("Error while saving data: %v", err)
		} else {
			updatedModel := m.categoryGroupModel.UpdateData(m.Data.Root)
			m.categoryGroupModel = &updatedModel
		}
	}
	return m, nil
}
func (m App) handleGroupDeleteMsg() (tea.Model, tea.Cmd) {
	return m, nil
}
func (m App) handleGroupUpdateMsg() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m App) handleGroupManageCategoriesMsg() (tea.Model, tea.Cmd) {
	return m, nil
}
