package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/data"
)

type IncomeModel struct {
	AppData
	WindowSize
	MonthYear

	cursor        int
	incomeEntries []data.IncomeRecord
}

func NewIncomeModel(initialData *data.DataRoot, month time.Month, year int) *IncomeModel {
	mKey := GetMonthKey(month, year)

	var incomeEntries []data.IncomeRecord

	if incomes, ok := initialData.MonthlyData[mKey]; ok {
		incomeEntries = incomes.Incomes
	} else {
		incomeEntries = make([]data.IncomeRecord, 0)
	}

	return &IncomeModel{
		incomeEntries: incomeEntries,
		MonthYear: MonthYear{
			CurrentMonth: month,
			CurrentYear:  year,
		},
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
