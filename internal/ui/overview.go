package ui

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/madalinpopa/gocost/internal/config"
	"github.com/madalinpopa/gocost/internal/data"
	"github.com/spf13/viper"
)

type focusLevel int

const (
	focusLevelGroups focusLevel = iota
	focusLevelCategories
)

type MonthlyModel struct {
	AppData
	MonthYear
	WindowSize

	level                focusLevel
	focusedGroupIndex    int
	focusedCategoryIndex int
}

func NewMonthlyModel(data *data.DataRoot, month time.Month, year int) *MonthlyModel {
	return &MonthlyModel{
		AppData: AppData{
			Data: data,
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
	record, ok := m.Data.MonthlyData[monthKey]

	if ok {
		totalIncome = m.getMonthIncome(record)
		totalExpenses, totalExpensesGroup = m.getMonthExpenses(record, m.Data.CategoryGroups)
	}

	balance := totalIncome - totalExpenses

	_, _ = balance, totalExpensesGroup

	// amountColWidth := 12        // For "1234.56 $"
	// statusColWidth := 11        // For "[Not Paid]"
	// notesIndicatorColWidth := 4 // For " (N)"

	// columnSpacer := "  " // Two spaces

	header := m.getHeader(totalIncome)
	footer := m.getFooter()

	b.WriteString(header)
	b.WriteString(footer)

	return AppStyle.Render(b.String())

}

func (m MonthlyModel) getHeader(totalIncome float64) string {
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

	bottomBorder := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(ColorSubtleBorder).
		Width(m.Width - AppStyle.GetHorizontalFrameSize()).
		Render("")

	b.WriteString(headerStr)
	b.WriteString("\n")
	b.WriteString(bottomBorder)
	b.WriteString("\n")

	defaultCurrency := viper.GetString(config.CurrencyField)
	income := fmt.Sprintf("Total Income: %.2f %s", totalIncome, defaultCurrency)

	b.WriteString(MutedText.Render(income))
	return b.String()
}

func (m MonthlyModel) getCurrentMonth() string {
	monthKey := fmt.Sprintf("%s-%d", m.CurrentMonth.String(), m.CurrentYear)
	return monthKey
}

func (m MonthlyModel) getFooter() string {
	var b bytes.Buffer

	footerStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(ColorSubtleBorder).
		Width(m.Width - AppStyle.GetHorizontalPadding()).
		PaddingTop(1)

	keyHints := "j/k: Nav | Ent: Select/Edit | Ctrl+i: Income | Ctrl+g: Groups"

	b.WriteString("\n\n")
	b.WriteString(footerStyle.Render(lipgloss.JoinVertical(lipgloss.Left, MutedText.Render(keyHints))))

	return b.String()
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

func (m MonthlyModel) SetMonthYear(month time.Month, year int) MonthlyModel {
	m.CurrentMonth = month
	m.CurrentYear = year

	// Reset focus when month changes, back to group navigation
	m.level = focusLevelGroups
	if len(m.Data.CategoryGroups) > 0 {
		m.focusedCategoryIndex = 0
	} else {
		m.focusedGroupIndex = -1
	}
	m.focusedCategoryIndex = -1

	return m
}
