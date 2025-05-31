package app

import tea "github.com/charmbracelet/bubbletea"

func (m App) handleMonthlyView(key string) (tea.Model, tea.Cmd) {

	switch key {

	case "q":
		return m, tea.Quit
	}

	return m, nil
}
