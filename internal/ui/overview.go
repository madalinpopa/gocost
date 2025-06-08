package ui

import (
	"bytes"
	"fmt"
	"sort"
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

func NewMonthlyModel(data *data.DataRoot, month time.Month, year int) MonthlyModel {
	return MonthlyModel{
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

	case tea.KeyMsg:
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
		totalExpenses, totalExpensesGroup = m.getMonthExpenses(record)
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

	var keyHints string
	switch m.level {
	case focusLevelGroups:
		keyHints = "j/k: Nav | Ent: Select | i: Income | c: Categories | g: Groups | h/l: Month"
	case focusLevelCategories:
		keyHints = "j/k: Nav | Ent: Select | Esc: Back | i: Income | c: Categories | g: Groups | h/l: Month"
	}
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
		return b.String()
	}

	monthKey := GetMonthKey(m.CurrentMonth, m.CurrentYear)
	record, ok := m.Data.MonthlyData[monthKey]
	if !ok {
		b.WriteString(MutedText.Render("No data for this month"))
		b.WriteString("\n")
		return b.String()
	}

	// Group categories by their GroupID
	categoriesByGroup := make(map[string][]data.Category)
	for _, category := range record.Categories {
		categoriesByGroup[category.GroupID] = append(categoriesByGroup[category.GroupID], category)
	}

	// Get ordered list of groups
	orderedGroups := m.getOrderedGroups()

	var expenseSectionContent []string
	for groupIdx, group := range orderedGroups {
		groupStyle := NormalListItem
		groupPrefix := "  "

		if m.level == focusLevelGroups && groupIdx == m.focusedGroupIndex {
			groupStyle = FocusedListItem
			groupPrefix = "> "
		} else if m.level == focusLevelCategories && groupIdx == m.focusedGroupIndex {
			groupStyle = HeaderText.Bold(false).Foreground(lipgloss.Color("220"))
			groupPrefix = ">> "
		}

		var groupTotal decimal.Decimal
		if totalGroupExpenses != nil {
			groupTotal = totalGroupExpenses[group.GroupID]
		}
		groupNameRender := groupStyle.Render(fmt.Sprintf("%s%s", groupPrefix, group.GroupName))
		groupTotalRender := groupStyle.Render(fmt.Sprintf("Total: %s %s", groupTotal.String(), currency))

		groupHeaderSpacerWidth := max(m.Width-lipgloss.Width(groupNameRender)-lipgloss.Width(groupTotalRender)-AppStyle.GetHorizontalPadding(), 0)
		groupHeader := lipgloss.JoinHorizontal(lipgloss.Left, groupNameRender, lipgloss.NewStyle().Width(groupHeaderSpacerWidth).Render(""), groupTotalRender)
		expenseSectionContent = append(expenseSectionContent, groupHeader)

		// Display categories within this group if we're in category navigation mode and this is the focused group
		if m.level == focusLevelCategories && groupIdx == m.focusedGroupIndex {
			categories := categoriesByGroup[group.GroupID]
			for catIdx, category := range categories {
				catStyle := NormalListItem
				catPrefix := "    "

				if catIdx == m.focusedCategoryIndex {
					catStyle = FocusedListItem
					catPrefix = "  > "
				}

				// Calculate category total
				var categoryTotal decimal.Decimal
				for _, expense := range category.Expense {
					amount := decimal.NewFromFloat(expense.Amount)
					categoryTotal = categoryTotal.Add(amount)
				}

				catNameRender := catStyle.Render(fmt.Sprintf("%s%s", catPrefix, category.CategoryName))
				catTotalRender := catStyle.Render(fmt.Sprintf("%s %s", categoryTotal.String(), currency))

				catSpacerWidth := max(m.Width-lipgloss.Width(catNameRender)-lipgloss.Width(catTotalRender)-AppStyle.GetHorizontalPadding(), 0)
				categoryLine := lipgloss.JoinHorizontal(lipgloss.Left, catNameRender, lipgloss.NewStyle().Width(catSpacerWidth).Render(""), catTotalRender)
				expenseSectionContent = append(expenseSectionContent, categoryLine)
			}
		}
	}

	b.WriteString(strings.Join(expenseSectionContent, "\n"))
	return b.String()
}

func (m MonthlyModel) getOrderedGroups() []data.CategoryGroup {
	var orderedGroups []data.CategoryGroup
	for _, group := range m.Data.CategoryGroups {
		orderedGroups = append(orderedGroups, group)
	}

	sort.Slice(orderedGroups, func(i, j int) bool {
		return orderedGroups[i].Order < orderedGroups[j].Order
	})

	return orderedGroups
}

func (m MonthlyModel) getMonthIncome(monthRecord data.MonthlyRecord) decimal.Decimal {
	var totalIncome decimal.Decimal
	for _, income := range monthRecord.Incomes {
		amount := decimal.NewFromFloat(income.Amount)
		totalIncome = totalIncome.Add(amount)
	}
	return totalIncome
}

func (m MonthlyModel) getMonthExpenses(mr data.MonthlyRecord) (decimal.Decimal, map[string]decimal.Decimal) {
	var expenseTotals decimal.Decimal
	groupTotals := make(map[string]decimal.Decimal)

	for _, category := range mr.Categories {
		var categoryTotal decimal.Decimal
		for _, expense := range category.Expense {
			amount := decimal.NewFromFloat(expense.Amount)
			categoryTotal = categoryTotal.Add(amount)
		}
		expenseTotals = expenseTotals.Add(categoryTotal)
		groupTotals[category.GroupID] = groupTotals[category.GroupID].Add(categoryTotal)
	}

	return expenseTotals, groupTotals
}

func (m MonthlyModel) handleGroupNavigation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {

	numGroups := len(m.Data.CategoryGroups)

	switch msg.String() {

	case "j", "down":
		if numGroups > 0 {
			m.focusedGroupIndex = (m.focusedGroupIndex + 1) % numGroups
		}
	case "k", "up":
		if numGroups > 0 {
			m.focusedGroupIndex--
			if m.focusedGroupIndex < 0 {
				m.focusedGroupIndex = numGroups - 1
			}
		}
	case "enter":
		if numGroups > 0 && m.focusedGroupIndex >= 0 && m.focusedGroupIndex < numGroups {
			// Get the selected group and check if it has categories
			monthKey := GetMonthKey(m.CurrentMonth, m.CurrentYear)
			if record, ok := m.Data.MonthlyData[monthKey]; ok {
				// Get ordered list of groups to find the correct group
				orderedGroups := m.getOrderedGroups()

				if m.focusedGroupIndex < len(orderedGroups) {
					selectedGroup := orderedGroups[m.focusedGroupIndex]

					// Count categories in this group
					var categoryCount int
					for _, category := range record.Categories {
						if category.GroupID == selectedGroup.GroupID {
							categoryCount++
						}
					}

					if categoryCount > 0 {
						m.level = focusLevelCategories
						m.focusedCategoryIndex = 0
					}
				}
			}
		}

	}
	return m, nil
}

func (m MonthlyModel) handleCategoryNavigation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {

	monthKey := GetMonthKey(m.CurrentMonth, m.CurrentYear)
	record, ok := m.Data.MonthlyData[monthKey]
	if !ok {
		return m, nil
	}

	// Get the currently focused group
	orderedGroups := m.getOrderedGroups()

	if m.focusedGroupIndex >= len(orderedGroups) {
		return m, nil
	}

	selectedGroup := orderedGroups[m.focusedGroupIndex]

	// Get categories for this group
	var categoriesInGroup []data.Category
	for _, category := range record.Categories {
		if category.GroupID == selectedGroup.GroupID {
			categoriesInGroup = append(categoriesInGroup, category)
		}
	}

	numCategories := len(categoriesInGroup)

	switch msg.String() {

	case "j", "down":
		if numCategories > 0 {
			m.focusedCategoryIndex = (m.focusedCategoryIndex + 1) % numCategories
		}
	case "k", "up":
		if numCategories > 0 {
			m.focusedCategoryIndex--
			if m.focusedCategoryIndex < 0 {
				m.focusedCategoryIndex = numCategories - 1
			}
		}
	case "enter":
		if numCategories > 0 && m.focusedCategoryIndex >= 0 && m.focusedCategoryIndex < numCategories {
			// TODO: Navigate to expense view for selected category
			// selectedCategory := categoriesInGroup[m.focusedCategoryIndex]
			// Could emit a message to switch to expense management view
		}
	case "esc":
		// Go back to group navigation
		m.level = focusLevelGroups
		m.focusedCategoryIndex = 0

	}
	return m, nil
}

func (m MonthlyModel) SetMonthYear(month time.Month, year int) MonthlyModel {
	m.CurrentMonth = month
	m.CurrentYear = year

	// Reset focus when month changes, back to group navigation
	m.level = focusLevelGroups
	if len(m.Data.CategoryGroups) > 0 {
		m.focusedGroupIndex = 0
	} else {
		m.focusedGroupIndex = -1
	}
	m.focusedCategoryIndex = -1

	return m
}
