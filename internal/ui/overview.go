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

	Level                focusLevel
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
		// Handle populate command at any level - only if current month has no categories
		if msg.String() == "p" && m.currentMonthHasNoCategories() {
			return m.handlePopulateCategories()
		}

		switch m.Level {

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
	populateHint := ""
	if m.currentMonthHasNoCategories() {
		populateHint = " | p: Populate"
	}

	switch m.Level {
	case focusLevelGroups:
		keyHints = "j/k: Nav | Ent: Select" + populateHint + " | i: Income | c: Categories | g: Groups | h/l: Month"
	case focusLevelCategories:
		keyHints = "j/k: Nav | Ent: Select | Esc: Back" + populateHint + " | i: Income | c: Categories | g: Groups | h/l: Month"
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

	// Get ordered list of groups and filter for visible ones
	orderedGroups := m.getOrderedGroups()
	visibleGroups := m.getVisibleGroups(orderedGroups, categoriesByGroup)

	if len(visibleGroups) == 0 {
		b.WriteString(MutedText.Render("No categories for this month"))
		b.WriteString("\n")
		return b.String()
	}

	var expenseSectionContent []string
	for visibleIdx, group := range visibleGroups {
		groupStyle := NormalListItem
		groupPrefix := "  "

		if m.Level == focusLevelGroups && visibleIdx == m.focusedGroupIndex {
			groupStyle = FocusedListItem
			groupPrefix = "> "
		} else if m.Level == focusLevelCategories && visibleIdx == m.focusedGroupIndex {
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
		if m.Level == focusLevelCategories && visibleIdx == m.focusedGroupIndex {
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

func (m MonthlyModel) getVisibleGroups(orderedGroups []data.CategoryGroup, categoriesByGroup map[string][]data.Category) []data.CategoryGroup {
	var visibleGroups []data.CategoryGroup
	for _, group := range orderedGroups {
		if _, hasCategories := categoriesByGroup[group.GroupID]; hasCategories {
			visibleGroups = append(visibleGroups, group)
		}
	}
	return visibleGroups
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
	// Get visible groups for navigation
	monthKey := GetMonthKey(m.CurrentMonth, m.CurrentYear)
	record, ok := m.Data.MonthlyData[monthKey]
	if !ok {
		return m, nil
	}

	// Group categories by their GroupID
	categoriesByGroup := make(map[string][]data.Category)
	for _, category := range record.Categories {
		categoriesByGroup[category.GroupID] = append(categoriesByGroup[category.GroupID], category)
	}

	// Get visible groups
	orderedGroups := m.getOrderedGroups()
	visibleGroups := m.getVisibleGroups(orderedGroups, categoriesByGroup)
	numVisibleGroups := len(visibleGroups)

	switch msg.String() {

	case "j", "down":
		if numVisibleGroups > 0 {
			m.focusedGroupIndex = (m.focusedGroupIndex + 1) % numVisibleGroups
		}
	case "k", "up":
		if numVisibleGroups > 0 {
			m.focusedGroupIndex--
			if m.focusedGroupIndex < 0 {
				m.focusedGroupIndex = numVisibleGroups - 1
			}
		}
	case "enter":
		if numVisibleGroups > 0 && m.focusedGroupIndex >= 0 && m.focusedGroupIndex < numVisibleGroups {
			// The focused index now directly maps to visible groups
			m.Level = focusLevelCategories
			m.focusedCategoryIndex = 0
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

	// Group categories by their GroupID
	categoriesByGroup := make(map[string][]data.Category)
	for _, category := range record.Categories {
		categoriesByGroup[category.GroupID] = append(categoriesByGroup[category.GroupID], category)
	}

	// Get visible groups
	orderedGroups := m.getOrderedGroups()
	visibleGroups := m.getVisibleGroups(orderedGroups, categoriesByGroup)

	if m.focusedGroupIndex >= len(visibleGroups) {
		return m, nil
	}

	selectedGroup := visibleGroups[m.focusedGroupIndex]

	// Get categories for this group (they already exist since this is a visible group)
	categoriesInGroup := categoriesByGroup[selectedGroup.GroupID]

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
		m.Level = focusLevelGroups
		m.focusedCategoryIndex = 0

	}
	return m, nil
}

func (m MonthlyModel) handlePopulateCategories() (tea.Model, tea.Cmd) {
	currentMonthKey := GetMonthKey(m.CurrentMonth, m.CurrentYear)
	prevYear, prevMonth := GetPreviousMonth(m.CurrentYear, m.CurrentMonth)
	prevMonthKey := GetMonthKey(prevMonth, prevYear)

	return m, func() tea.Msg {
		return PopulateCategoriesMsg{
			CurrentMonthKey:  currentMonthKey,
			PreviousMonthKey: prevMonthKey,
		}
	}
}

func (m MonthlyModel) currentMonthHasNoCategories() bool {
	currentMonthKey := GetMonthKey(m.CurrentMonth, m.CurrentYear)
	if record, exists := m.Data.MonthlyData[currentMonthKey]; exists {
		return len(record.Categories) == 0
	}
	return true
}

func (m MonthlyModel) SetMonthYear(month time.Month, year int) MonthlyModel {
	m.CurrentMonth = month
	m.CurrentYear = year
	return m.ResetFocus()
}

func (m MonthlyModel) ResetFocus() MonthlyModel {
	m.Level = focusLevelGroups
	m.focusedGroupIndex = 0
	m.focusedCategoryIndex = 0
	return m
}

func (m MonthlyModel) UpdateData(updatedData *data.DataRoot) MonthlyModel {
	m.Data = updatedData
	return m.ResetFocus()
}
