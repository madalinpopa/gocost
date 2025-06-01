package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/ui"
)

func (m App) handleMonthlyView(key string) (tea.Model, tea.Cmd) {

	switch key {

	case "ctrl+c", "q":
		return m, tea.Quit

	case "h":
		m.CurrentYear, m.CurrentMonth = ui.GetPreviousMonth(m.CurrentYear, m.CurrentMonth)
		if m.monthlyModel != nil {
			m.monthlyModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
		}

		return m, nil

	case "l":
		m.CurrentYear, m.CurrentMonth = ui.GetNextMonth(m.CurrentYear, m.CurrentMonth)
		if m.monthlyModel != nil {
			m.monthlyModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
		}

		return m, nil
	}

	return m, nil
}
