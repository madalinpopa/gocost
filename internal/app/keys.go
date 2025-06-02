package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/ui"
)

func (m App) handleMonthlyViewKeys(key string) (tea.Model, tea.Cmd) {

	switch key {

	case "ctrl+c", "q":
		return m, tea.Quit

	case "ctrl+g":
		if m.activeView == viewMonthlyOverview {
			m.activeView = viewCategoryGroup
			return m, m.CategoryGroupModel.Init()
		}

	case "h":
		m.CurrentYear, m.CurrentMonth = ui.GetPreviousMonth(m.CurrentYear, m.CurrentMonth)
		if m.MonthlyModel != nil {
			updatedModel := m.MonthlyModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
			m.MonthlyModel = &updatedModel
		}

		return m, nil

	case "l":
		m.CurrentYear, m.CurrentMonth = ui.GetNextMonth(m.CurrentYear, m.CurrentMonth)
		if m.MonthlyModel != nil {
			updatedModel := m.MonthlyModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
			m.MonthlyModel = &updatedModel
		}

		return m, nil
	}

	return m, nil
}

func (m App) handleCategoryGroupViewKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc", "q":
		m.activeView = viewMonthlyOverview
		return m, nil
	}
	return m, nil
}
