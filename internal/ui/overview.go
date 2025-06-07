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
	"github.com/shopspring/decimal"
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

	case tea.Msg:
		switch m.level {

		case focusLevelGroups:
			return m.handleGroupNavigation(msg)
		case focusLevelCategories:
			return m.handleCategoryNavigation(msg)

		}
	}

	return m, nil
}

func (m MonthlyModel) View() string {
	var b strings.Builder

	defaultCurrency := viper.GetString(config.CurrencyField)

	var totalExpenses decimal.Decimal
	var totalExpensesGroup map[string]decimal.Decimal
	var totalIncome decimal.Decimal

	monthKey := GetMonthKey(m.CurrentMonth, m.CurrentYear)
	record, ok := m.Data.MonthlyData[monthKey]
	if ok {
		totalIncome = m.getMonthIncome(record)
		totalExpenses, totalExpensesGroup = m.getMonthExpenses(record, m.Data.CategoryGroups)
	}

	balance := totalIncome.Sub(totalExpenses)

	_, _ = balance, totalExpensesGroup

	// amountColWidth := 12        // For "1234.56 $"
	// statusColWidth := 11        // For "[Not Paid]"
	// notesIndicatorColWidth := 4 // For " (N)"

	// columnSpacer := "  " // Two spaces

	header := m.getHeader(totalIncome, defaultCurrency)
	content := m.getContent(totalExpensesGroup, defaultCurrency)
	footer := m.getFooter(totalExpenses, balance, defaultCurrency)

	b.WriteString(header)
	b.WriteString(content)
	b.WriteString(footer)

	return AppStyle.Render(b.String())

}

func (m MonthlyModel) getHeader(totalIncome decimal.Decimal, defaultCurrency string) string {
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

	income := fmt.Sprintf("Total Income: %s %s", totalIncome.String(), defaultCurrency)
	b.WriteString(MutedText.Render(income))
	b.WriteString("\n\n")

	return b.String()
}

func (m MonthlyModel) getFooter(totalExpenses, balance decimal.Decimal, defaultCurrency string) string {
	var b bytes.Buffer

	footerStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(ColorSubtleBorder).
		Width(m.Width - AppStyle.GetHorizontalPadding()).
		PaddingTop(1)

	keyHints := "j/k: Nav | Ent: Select/Edit | i: Income | g: Groups"
	totalExpensesStr := fmt.Sprintf("Total Expenses: %s %s", totalExpenses.String(), defaultCurrency)

	balanceStr := fmt.Sprintf("Balance: %s %s", balance.String(), defaultCurrency)
	footerSummarySpacerWidth := max(m.Width-lipgloss.Width(totalExpensesStr)-lipgloss.Width(balanceStr)-AppStyle.GetHorizontalPadding(), 0)

	space := lipgloss.NewStyle().Width(footerSummarySpacerWidth).Render("")
	footerSummary := lipgloss.JoinHorizontal(lipgloss.Top, totalExpensesStr, space, balanceStr)

	b.WriteString("\n\n")
	b.WriteString(footerStyle.Render(lipgloss.JoinVertical(lipgloss.Left, footerSummary, MutedText.Render(keyHints))))

	return b.String()
}

func (m MonthlyModel) getContent(totalGroupExpenses map[string]decimal.Decimal, currency string) string {
	var b strings.Builder

	if len(m.Data.CategoryGroups) == 0 {
		b.WriteString(MutedText.Render("No category groups. (g)"))
		b.WriteString("\n")
	}
	var expenseSectionContent []string
	for groupIdx, group := range m.Data.CategoryGroups {
		groupStyle := NormalListItem
		groupPrefix := "  "

		if m.level == focusLevelGroups && groupIdx == m.focusedGroupIndex {
			groupStyle = FocusedListItem
			groupPrefix = "> "
		} else if m.level == focusLevelCategories && groupIdx == m.focusedGroupIndex {
			groupStyle = HeaderText.Bold(false).Foreground(lipgloss.Color("220"))
			groupPrefix = ">>"
		}
		groupTotal := totalGroupExpenses[group.GroupID]
		groupNameRender := groupStyle.Render(fmt.Sprintf("%s %s", groupPrefix, group.GroupName))
		groupTotalRender := groupStyle.Render(fmt.Sprintf("Total: %s %s", groupTotal.String(), currency))

		groupHeaderSpacerWidth := max(m.Width - lipgloss.Width(groupNameRender) - lipgloss.Width(groupTotalRender) - AppStyle.GetHorizontalPadding())
		expenseSectionContent = append(expenseSectionContent, lipgloss.JoinHorizontal(lipgloss.Left, groupNameRender, lipgloss.NewStyle().Width(groupHeaderSpacerWidth).Render(""), groupTotalRender))
	}

	b.WriteString(strings.Join(expenseSectionContent, "\n"))
	return b.String()
}

func (m MonthlyModel) getMonthIncome(monthRecord data.MonthlyRecord) decimal.Decimal {
	var totalIncome decimal.Decimal
	for _, income := range monthRecord.Incomes {
		amount := decimal.NewFromFloat(income.Amount)
		totalIncome = totalIncome.Add(amount)
	}
	return totalIncome
}

func (m MonthlyModel) getMonthExpenses(mr data.MonthlyRecord, g []data.CategoryGroup) (decimal.Decimal, map[string]decimal.Decimal) {
	var expenseTotals decimal.Decimal
	groupTotals := make(map[string]decimal.Decimal)

	for _, expense := range mr.Expenses {
		expenseDecimal := decimal.NewFromFloat(expense.Amount)
		expenseTotals = expenseTotals.Add(expenseDecimal)

		for _, group := range g {

			for _, cat := range group.Categories {
				if cat.CatID == expense.CatID {
					amount := decimal.NewFromFloat(expense.Amount)
					groupTotals[group.GroupID] = groupTotals[group.GroupID].Add(amount)
					break
				}
			}

		}
	}
	return expenseTotals, groupTotals
}

func (m MonthlyModel) handleGroupNavigation(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m MonthlyModel) handleCategoryNavigation(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
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
