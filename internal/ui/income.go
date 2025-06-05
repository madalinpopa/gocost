package ui

import tea "github.com/charmbracelet/bubbletea"

type IncomeModel struct{}

func (m IncomeModel) Init() tea.Cmd {
	return nil
}

func (m IncomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m IncomeModel) View() string {
	return "Hello from income"
}
