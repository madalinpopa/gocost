package ui

import tea "github.com/charmbracelet/bubbletea"

type MonthlyModel struct {
	Data
	MonthYear
	WindowSize
}

func (m MonthlyModel) Init() tea.Cmd {
	return nil
}

func (m MonthlyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m MonthlyModel) View() string {
	return "Monthly view"
}
