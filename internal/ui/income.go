package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/data"
)

type IncomeModel struct {
	AppData
	WindowSize
}

func NewIncomeModel(initialData *data.DataRoot) *IncomeModel {
	return &IncomeModel{
		AppData: AppData{
			Data: initialData,
		},
	}
}

func (m IncomeModel) Init() tea.Cmd {
	return nil
}

func (m IncomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	}

	return m, nil
}

func (m IncomeModel) View() string {
	return "Hello from income"
}
