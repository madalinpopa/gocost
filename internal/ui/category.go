package ui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/data"
)

type CategoryModel struct {
	AppData
	WindowSize
	MonthYear

	MonthKey   string
	cursor     int
	categories []data.Category
}

func NewCategoryModel(initialData *data.DataRoot, month time.Month, year int) *CategoryModel {

	monthKey := GetMonthKey(month, year)

	var categories []data.Category
	if record, ok := initialData.MonthlyData[monthKey]; ok {
		categories = record.Categories
	}

	if categories == nil {
		categories = make([]data.Category, 0)
	}

	m := CategoryModel{
		AppData: AppData{
			Data: initialData,
		},
		categories: categories,
	}

	return &m
}

func (m CategoryModel) Init() tea.Cmd {
	return nil
}

func (m CategoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "q", "esc":
			return m, func() tea.Msg { return MonthlyViewMsg{} }

		case "j", "down":
			if len(m.categories) > 0 {
				m.cursor = (m.cursor + 1) % len(m.categories)
			}
			return m, nil

		case "k", "up":
			if len(m.categories) > 0 {
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(m.categories) - 1
				}
			}
			return m, nil
		}
	}

	return m, nil
}

func (m CategoryModel) View() string {
	var b strings.Builder

	b.WriteString(HeaderText.Render("Manage Expense Categories"))
	b.WriteString("\n\n")

	return b.String()
}
