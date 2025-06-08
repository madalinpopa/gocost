package ui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/data"
)

type ExpenseModel struct {
	AppData
	WindowSize
	MonthYear

	MonthKey string
}

func NewExpenseModel(initialData *data.DataRoot, month time.Month, year int) ExpenseModel {

	monthKey := GetMonthKey(month, year)

	return ExpenseModel{
		AppData: AppData{
			Data: initialData,
		},
		MonthKey: monthKey,
	}
}

func (m ExpenseModel) Init() tea.Cmd {
	return nil
}

func (m ExpenseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ExpenseModel) View() string {
	var b strings.Builder

	return b.String()
}
