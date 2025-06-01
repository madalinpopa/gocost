package app

import tea "github.com/charmbracelet/bubbletea"

func (m App) handleMonthlyViewMsg() (tea.Model, tea.Cmd) {
	m.activeView = viewMonthlyOverview
	if m.monthlyModel != nil {
		updatedModel := m.monthlyModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
		m.monthlyModel = &updatedModel
	}
	return m, nil
}
