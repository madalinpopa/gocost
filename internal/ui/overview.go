package ui

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/madalinpopa/gocost/internal/data"
)

type focusLevel int

const (
	focusLevelGroups focusLevel = iota
	focusLevelCategories
)

type MonthlyModel struct {
	Data
	MonthYear
	WindowSize

	level                focusLevel
	focusedGroupIndex    int
	focusedCategoryIndex int
}

func NewMonthlyModel(data *data.DataRoot, month time.Month, year int) *MonthlyModel {
	return &MonthlyModel{
		Data: Data{
			Root: data,
		},
		MonthYear: MonthYear{
			CurrentMonth: month,
			CurrentYear:  year,
		},
	}
}

func (m MonthlyModel) Init() tea.Cmd {
	return nil
}

func (m MonthlyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	}
	return m, nil
}

func (m MonthlyModel) View() string {
	var b strings.Builder

	var totalExpenses float64
	var totalExpensesGroup map[string]float64
	var totalIncome float64

	monthKey := m.getCurrentMonth()
	record, ok := m.Data.Root.MonthlyData[monthKey]

	if ok {
		totalIncome = m.getMonthIncome(record)
		totalExpenses, totalExpensesGroup = m.getMonthExpenses(record, m.Data.Root.CategoryGroups)
	}

	balance := totalIncome - totalExpenses

	_, _ = balance, totalExpensesGroup

	header := m.getHeader()

	b.WriteString(header)

	return AppStyle.Render(b.String())

}

func (m MonthlyModel) getHeader() string {
	var b bytes.Buffer

	headerLeft := fmt.Sprintf("Month: %s %d", m.CurrentMonth.String(), m.CurrentYear)
	headerRight := MutedText.Render("(h/l Month)")

	frameSize := AppStyle.GetHorizontalFrameSize()
	headerLeftSize := lipgloss.Width(headerLeft)
	headerRightSize := lipgloss.Width(headerRight)

	// Calculate available width for the spacer in the header
	headerSpacerWidth := max(m.Width-headerLeftSize-headerRightSize-frameSize, 0)

	headerStr := lipgloss.JoinHorizontal(
		lipgloss.Top,
		HeaderText.Render(headerLeft),
		lipgloss.NewStyle().Width(headerSpacerWidth).Render(""),
		headerRight,
	)

	b.WriteString(headerStr)
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(ColorSubtleBorder).Width(m.Width - AppStyle.GetHorizontalFrameSize()).Render(""))
	b.WriteString("\n")

	return b.String()
}

func (m MonthlyModel) getCurrentMonth() string {
	monthKey := fmt.Sprintf("%s-%d", m.CurrentMonth.String(), m.CurrentYear)
	return monthKey
}

func (m MonthlyModel) getMonthIncome(monthRecord data.MonthlyRecord) float64 {
	var totalIncome float64
	for _, income := range monthRecord.Incomes {
		totalIncome += income.Amount
	}
	return totalIncome
}

func (m MonthlyModel) getMonthExpenses(mr data.MonthlyRecord, g []data.CategoryGroup) (float64, map[string]float64) {
	var expenseTotals float64
	groupTotals := make(map[string]float64)

	for _, expense := range mr.Expenses {
		expenseTotals += expense.Amount

		for _, group := range g {

			for _, cat := range group.Categories {
				if cat.CatID == expense.CatID {
					groupTotals[group.GroupID] += expense.Amount
					break
				}
			}

		}
	}
	return expenseTotals, groupTotals
}

func (m MonthlyModel) SetMonthYear(month time.Month, year int) {
	m.CurrentMonth = month
	m.CurrentYear = year

	// Reset focus when month changes, back to group navigation
	m.level = focusLevelGroups
	if len(m.Data.Root.CategoryGroups) > 0 {
		m.focusedCategoryIndex = 0
	} else {
		m.focusedGroupIndex = -1
	}
	m.focusedCategoryIndex = -1
}
