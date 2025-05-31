package ui

import tea "github.com/charmbracelet/bubbletea"

type CategoryGroupModel struct {
	Data
	WindowSize
}

func (m CategoryGroupModel) Init() tea.Cmd {
	return nil
}

func (m CategoryGroupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m CategoryGroupModel) View() string {
	return "Category group view"
}
