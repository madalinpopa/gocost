package ui

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/madalinpopa/gocost/internal/config"
	"github.com/madalinpopa/gocost/internal/domain"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

type focusLevel int

const (
	focusLevelGroups focusLevel = iota
	focusLevelCategories
)

type MonthlyModel struct {
	MonthYear
	WindowSize

	Level                focusLevel
	focusedGroupIndex    int
	focusedCategoryIndex int

	categories     []domain.Category
	categoryGroups []domain.CategoryGroup
	incomes        []domain.IncomeRecord

	groupsViewport     viewport.Model
	categoriesViewport viewport.Model
	ready              bool
}

// NewMonthlyModel creates a new MonthlyModel instance.
func NewMonthlyModel(appData AppData, monthYear MonthYear) MonthlyModel {
	return MonthlyModel{
		MonthYear:          monthYear,
		categories:         appData.Categories,
		categoryGroups:     appData.CategoryGroups,
		incomes:            appData.Incomes,
		groupsViewport:     viewport.New(80, 20),
		categoriesViewport: viewport.New(80, 20),
		ready:              false,
	}
}

// Init initializes the MonthlyModel.
func (m MonthlyModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the MonthlyModel state.
func (m MonthlyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		headerHeight := lipgloss.Height(m.getHeader(decimal.Zero, ""))
		footerHeight := lipgloss.Height(m.getFooter(decimal.Zero, decimal.Zero, ""))
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			availableHeight := msg.Height - verticalMarginHeight - 4 // -4 for padding (2) and newlines (2)
			groupsViewportHeight := m.calculateGroupsViewportHeight(availableHeight)
			categoriesViewportHeight := m.calculateCategoriesViewportHeight(availableHeight)

			m.groupsViewport = viewport.New(msg.Width, groupsViewportHeight)
			m.categoriesViewport = viewport.New(msg.Width, categoriesViewportHeight)
			m.groupsViewport.YPosition = headerHeight
			m.categoriesViewport.YPosition = headerHeight
			m.groupsViewport.SetContent(m.getGroupsContent(nil, ""))
			m.categoriesViewport.SetContent(m.getCategoriesContent(nil, ""))
			m.ready = true
		} else {
			m.groupsViewport.Width = msg.Width
			m.categoriesViewport.Width = msg.Width
			availableHeight := msg.Height - verticalMarginHeight - 4 // -4 for padding (2) and newlines (2)
			groupsViewportHeight := m.calculateGroupsViewportHeight(availableHeight)
			categoriesViewportHeight := m.calculateCategoriesViewportHeight(availableHeight)
			m.groupsViewport.Height = groupsViewportHeight
			m.categoriesViewport.Height = categoriesViewportHeight
		}
		return m, nil

	case tea.KeyMsg:
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

	if m.ready {
		switch m.Level {
		case focusLevelGroups:
			m.groupsViewport, cmd = m.groupsViewport.Update(msg)
			cmds = append(cmds, cmd)
		case focusLevelCategories:
			m.categoriesViewport, cmd = m.categoriesViewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the MonthlyModel displaying the monthly overview.
func (m MonthlyModel) View() string {
	if !m.ready {
		return AppStyle.Width(m.Width).Height(m.Height).Render("\n  Initializing...")
	}

	var b strings.Builder

	defaultCurrency := viper.GetString(config.CurrencyField)

	var totalExpenses decimal.Decimal
	var totalExpensesGroup map[string]decimal.Decimal

	totalIncome := m.getMonthIncome()
	totalExpenses, totalExpensesGroup = m.getMonthExpenses()

	balance := totalIncome.Sub(totalExpenses)

	header := m.getHeader(totalIncome, defaultCurrency)
	footer := m.getFooter(totalExpenses, balance, defaultCurrency)

	if m.ready {
		switch m.Level {
		case focusLevelGroups:
			m.groupsViewport.SetContent(m.getGroupsContent(totalExpensesGroup, defaultCurrency))
		case focusLevelCategories:
			m.categoriesViewport.SetContent(m.getCategoriesContent(totalExpensesGroup, defaultCurrency))
		}
	}

	b.WriteString(header)
	b.WriteString("\n")

	switch m.Level {
	case focusLevelGroups:
		b.WriteString(m.groupsViewport.View())
	case focusLevelCategories:
		groupHeader := m.getCategoryGroupHeader(totalExpensesGroup, defaultCurrency)
		if groupHeader != "" {
			b.WriteString(groupHeader)
			b.WriteString("\n")
		}
		b.WriteString(m.categoriesViewport.View())
	}

	b.WriteString("\n")
	b.WriteString(footer)

	return AppStyle.Render(b.String())
}

// getGroupsContent generates the content for the groups viewport.
func (m MonthlyModel) getGroupsContent(totalGroupExpenses map[string]decimal.Decimal, currency string) string {
	var b strings.Builder

	if len(m.categoryGroups) == 0 {
		b.WriteString(MutedText.Render("No category groups. (g)"))
		return b.String()
	}

	if len(m.categories) == 0 {
		b.WriteString(MutedText.Render("No categories for this month."))
		return b.String()
	}

	// Group categories by their GroupID
	categoriesByGroup := make(map[string][]domain.Category)
	for _, category := range m.categories {
		categoriesByGroup[category.GroupID] = append(categoriesByGroup[category.GroupID], category)
	}

	// Get ordered list of groups and filter for visible ones
	orderedGroups := m.getOrderedGroups()
	visibleGroups := m.getVisibleGroups(orderedGroups, categoriesByGroup)

	if len(visibleGroups) == 0 {
		b.WriteString(MutedText.Render("No categories for this month"))
		return b.String()
	}

	for visibleIdx, group := range visibleGroups {
		groupStyle := NormalListItem
		groupPrefix := "  "

		if visibleIdx == m.focusedGroupIndex {
			groupStyle = FocusedListItem
			groupPrefix = "> "
		}

		var groupTotal decimal.Decimal
		if totalGroupExpenses != nil {
			groupTotal = totalGroupExpenses[group.GroupID]
		}
		groupNameRender := groupStyle.Render(fmt.Sprintf("%s%s", groupPrefix, group.GroupName))
		totalRender := MutedText.Render("Total:")
		groupTotalRender := groupStyle.Render(fmt.Sprintf("%s %s %s", totalRender, groupTotal.String(), currency))

		groupHeaderSpacerWidth := max(m.Width-lipgloss.Width(groupNameRender)-lipgloss.Width(groupTotalRender)-AppStyle.GetHorizontalPadding(), 0)
		groupHeader := lipgloss.JoinHorizontal(lipgloss.Left, groupNameRender, CreateSpacer(groupHeaderSpacerWidth).Render(""), groupTotalRender)

		b.WriteString(groupHeader)
		if visibleIdx < len(visibleGroups)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// getCategoriesContent generates the content for the categories viewport (without group header).
func (m MonthlyModel) getCategoriesContent(totalGroupExpenses map[string]decimal.Decimal, currency string) string {
	var b strings.Builder

	if len(m.categoryGroups) == 0 {
		b.WriteString(MutedText.Render("No category groups. (g)"))
		return b.String()
	}

	if len(m.categories) == 0 {
		b.WriteString(MutedText.Render("No categories for this month."))
		return b.String()
	}

	// Group categories by their GroupID
	categoriesByGroup := make(map[string][]domain.Category)
	for _, category := range m.categories {
		categoriesByGroup[category.GroupID] = append(categoriesByGroup[category.GroupID], category)
	}

	// Get ordered list of groups and filter for visible ones
	orderedGroups := m.getOrderedGroups()
	visibleGroups := m.getVisibleGroups(orderedGroups, categoriesByGroup)

	if len(visibleGroups) == 0 {
		b.WriteString(MutedText.Render("No categories for this month"))
		return b.String()
	}

	if m.focusedGroupIndex >= len(visibleGroups) {
		b.WriteString(MutedText.Render("Invalid group selection"))
		return b.String()
	}

	// Get the focused group
	selectedGroup := visibleGroups[m.focusedGroupIndex]
	categories := categoriesByGroup[selectedGroup.GroupID]

	amountColWidth := len("Amount")
	budgetColWidth := len("/Budget")
	statusColWidth := len("Status")
	notesColWidth := len("Notes")
	columnSpacing := 2

	for _, category := range categories {
		var expense domain.ExpenseRecord
		var hasExpense bool
		if len(category.Expense) > 0 {
			for _, exp := range category.Expense {
				expense = exp
				hasExpense = true
				break
			}
		}

		amountStr := "0.00"
		budgetStr := "0.00"
		statusStr := "Not Set"
		notesIndicator := ""

		if hasExpense {
			amountStr = fmt.Sprintf("%.2f", expense.Amount)
			budgetStr = fmt.Sprintf("%.2f", expense.Budget)
			statusStr = expense.Status
			if expense.Notes != "" {
				notesIndicator = " (N)"
			}
		}

		amountText := fmt.Sprintf("%s %s", amountStr, currency)
		budgetText := fmt.Sprintf("/%s %s", budgetStr, currency)
		statusText := fmt.Sprintf("[%s]", statusStr)

		if len(amountText) > amountColWidth {
			amountColWidth = len(amountText)
		}
		if len(budgetText) > budgetColWidth {
			budgetColWidth = len(budgetText)
		}
		if len(statusText) > statusColWidth {
			statusColWidth = len(statusText)
		}
		if len(notesIndicator) > notesColWidth {
			notesColWidth = len(notesIndicator)
		}
	}

	// Add column headers
	headerStyle := MutedText
	amountHeader := headerStyle.Render(CreateLeftAlignedColumn(amountColWidth).Render("Amount"))
	budgetHeader := headerStyle.Render(CreateLeftAlignedColumn(budgetColWidth).Render("/Budget"))
	statusHeader := headerStyle.Render(CreateLeftAlignedColumn(statusColWidth).Render("Status"))
	notesHeader := headerStyle.Render(CreateLeftAlignedColumn(notesColWidth).Render("Notes"))

	totalColumnsWidth := amountColWidth + budgetColWidth + statusColWidth + notesColWidth + (columnSpacing * 3)
	headerSpacerWidth := max(m.Width-totalColumnsWidth-AppStyle.GetHorizontalPadding(), 1)

	headerLine := lipgloss.JoinHorizontal(
		lipgloss.Left,
		CreateSpacer(headerSpacerWidth).Render(""),
		amountHeader,
		CreateColumnSpacer(columnSpacing).Render(""),
		budgetHeader,
		CreateColumnSpacer(columnSpacing).Render(""),
		statusHeader,
		CreateColumnSpacer(columnSpacing).Render(""),
		notesHeader,
	)
	b.WriteString(headerLine)
	b.WriteString("\n")

	// Display categories
	for catIdx, category := range categories {
		catStyle := NormalListItem
		catPrefix := "    "

		if catIdx == m.focusedCategoryIndex {
			catStyle = FocusedListItem
			catPrefix = "  > "
		}

		var expense domain.ExpenseRecord
		var hasExpense bool
		if len(category.Expense) > 0 {
			// Get the first (and should be only) expense for this category
			for _, exp := range category.Expense {
				expense = exp
				hasExpense = true
				break
			}
		}

		// Format expense data
		amountStr := "0.00"
		budgetStr := "0.00"
		statusStr := "Not Set"
		notesIndicator := ""

		if hasExpense {
			amountStr = fmt.Sprintf("%.2f", expense.Amount)
			budgetStr = fmt.Sprintf("%.2f", expense.Budget)
			statusStr = expense.Status
			if expense.Notes != "" {
				notesIndicator = " (N)"
			}
		}

		// Build category line with columns using consistent widths
		catNameRender := catStyle.Render(fmt.Sprintf("%s%s", catPrefix, category.CategoryName))
		amountRender := catStyle.Render(CreateRightAlignedColumn(amountColWidth).Render(fmt.Sprintf("%s %s", amountStr, currency)))
		budgetRender := catStyle.Render(CreateRightAlignedColumn(budgetColWidth).Render(fmt.Sprintf("/%s %s", budgetStr, currency)))
		statusRender := catStyle.Render(CreateCenterAlignedColumn(statusColWidth).Render(RenderStatusBadge(statusStr)))
		notesRender := catStyle.Render(CreateCenterAlignedColumn(notesColWidth).Render(notesIndicator))

		// Calculate spacing for category name
		nameWidth := lipgloss.Width(catNameRender)
		totalColumnsWidth := amountColWidth + budgetColWidth + statusColWidth + notesColWidth + (columnSpacing * 3)
		availableWidth := m.Width - AppStyle.GetHorizontalPadding()
		spacerWidth := max(availableWidth-nameWidth-totalColumnsWidth, 1)

		categoryLine := lipgloss.JoinHorizontal(
			lipgloss.Left,
			catNameRender,
			CreateSpacer(spacerWidth).Render(""),
			amountRender,
			CreateColumnSpacer(columnSpacing).Render(""),
			budgetRender,
			CreateColumnSpacer(columnSpacing).Render(""),
			statusRender,
			CreateColumnSpacer(columnSpacing).Render(""),
			notesRender,
		)
		b.WriteString(categoryLine)
		if catIdx < len(categories)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// getCategoryGroupHeader generates the sticky group header for category navigation.
func (m MonthlyModel) getCategoryGroupHeader(totalGroupExpenses map[string]decimal.Decimal, currency string) string {
	if len(m.categoryGroups) == 0 || len(m.categories) == 0 {
		return ""
	}

	// Group categories by their GroupID
	categoriesByGroup := make(map[string][]domain.Category)
	for _, category := range m.categories {
		categoriesByGroup[category.GroupID] = append(categoriesByGroup[category.GroupID], category)
	}

	// Get visible groups
	orderedGroups := m.getOrderedGroups()
	visibleGroups := m.getVisibleGroups(orderedGroups, categoriesByGroup)

	if len(visibleGroups) == 0 || m.focusedGroupIndex >= len(visibleGroups) {
		return ""
	}

	// Get the focused group
	selectedGroup := visibleGroups[m.focusedGroupIndex]

	// Display group header
	var groupTotal decimal.Decimal
	if totalGroupExpenses != nil {
		groupTotal = totalGroupExpenses[selectedGroup.GroupID]
	}
	groupNameRender := ActiveGroupStyle.Render(fmt.Sprintf(">> %s", selectedGroup.GroupName))
	totalRender := MutedText.Render("Total:")
	groupTotalRender := ActiveGroupStyle.Render(fmt.Sprintf("%s %s %s", totalRender, groupTotal.String(), currency))

	groupHeaderSpacerWidth := max(m.Width-lipgloss.Width(groupNameRender)-lipgloss.Width(groupTotalRender)-AppStyle.GetHorizontalPadding(), 0)
	groupHeader := lipgloss.JoinHorizontal(lipgloss.Left, groupNameRender, CreateSpacer(groupHeaderSpacerWidth).Render(""), groupTotalRender)

	return groupHeader
}

// calculateGroupsViewportHeight calculates the appropriate height for the groups viewport.
func (m MonthlyModel) calculateGroupsViewportHeight(availableHeight int) int {
	// Group categories by their GroupID
	categoriesByGroup := make(map[string][]domain.Category)
	for _, category := range m.categories {
		categoriesByGroup[category.GroupID] = append(categoriesByGroup[category.GroupID], category)
	}

	// Get visible groups
	orderedGroups := m.getOrderedGroups()
	visibleGroups := m.getVisibleGroups(orderedGroups, categoriesByGroup)

	desiredHeight := max(len(visibleGroups)+1, 1)
	return min(desiredHeight, max(1, availableHeight))
}

// calculateCategoriesViewportHeight calculates the appropriate height for the categories viewport.
func (m MonthlyModel) calculateCategoriesViewportHeight(availableHeight int) int {
	// Group categories by their GroupID
	categoriesByGroup := make(map[string][]domain.Category)
	for _, category := range m.categories {
		categoriesByGroup[category.GroupID] = append(categoriesByGroup[category.GroupID], category)
	}

	// Get visible groups
	orderedGroups := m.getOrderedGroups()
	visibleGroups := m.getVisibleGroups(orderedGroups, categoriesByGroup)

	if len(visibleGroups) == 0 || m.focusedGroupIndex >= len(visibleGroups) {
		return max(1, availableHeight)
	}

	selectedGroup := visibleGroups[m.focusedGroupIndex]
	categories := categoriesByGroup[selectedGroup.GroupID]

	// Account for sticky group header (1 line) when calculating available height
	groupHeaderHeight := 1
	adjustedAvailableHeight := max(availableHeight-groupHeaderHeight, 1)

	desiredHeight := max(len(categories)+2, 1)
	return min(desiredHeight, max(1, adjustedAvailableHeight))
}

// ensureGroupsCursorVisible ensures the cursor is visible in the groups viewport.
func (m MonthlyModel) ensureGroupsCursorVisible() MonthlyModel {
	if !m.ready {
		return m
	}

	content := m.getGroupsContent(nil, "")
	m.groupsViewport.SetContent(content)

	// Group categories by their GroupID
	categoriesByGroup := make(map[string][]domain.Category)
	for _, category := range m.categories {
		categoriesByGroup[category.GroupID] = append(categoriesByGroup[category.GroupID], category)
	}

	// Get visible groups
	orderedGroups := m.getOrderedGroups()
	visibleGroups := m.getVisibleGroups(orderedGroups, categoriesByGroup)

	if len(visibleGroups) == 0 {
		return m
	}

	viewportTop := m.groupsViewport.YOffset
	viewportBottom := viewportTop + m.groupsViewport.Height - 1

	if m.focusedGroupIndex > viewportBottom {
		newOffset := max(m.focusedGroupIndex-m.groupsViewport.Height+1, 0)
		m.groupsViewport.SetYOffset(newOffset)
	}
	if m.focusedGroupIndex < viewportTop {
		m.groupsViewport.SetYOffset(m.focusedGroupIndex)
	}
	return m
}

// ensureCategoriesCursorVisible ensures the cursor is visible in the categories viewport.
func (m MonthlyModel) ensureCategoriesCursorVisible() MonthlyModel {
	if !m.ready {
		return m
	}

	content := m.getCategoriesContent(nil, "")
	m.categoriesViewport.SetContent(content)

	// Group categories by their GroupID
	categoriesByGroup := make(map[string][]domain.Category)
	for _, category := range m.categories {
		categoriesByGroup[category.GroupID] = append(categoriesByGroup[category.GroupID], category)
	}

	// Get visible groups
	orderedGroups := m.getOrderedGroups()
	visibleGroups := m.getVisibleGroups(orderedGroups, categoriesByGroup)

	if len(visibleGroups) == 0 || m.focusedGroupIndex >= len(visibleGroups) {
		return m
	}

	selectedGroup := visibleGroups[m.focusedGroupIndex]
	categories := categoriesByGroup[selectedGroup.GroupID]

	if len(categories) == 0 {
		return m
	}

	viewportTop := m.categoriesViewport.YOffset
	viewportBottom := viewportTop + m.categoriesViewport.Height - 1

	// +1 offset for column headers (group header is now sticky and outside viewport)
	adjustedCursor := m.focusedCategoryIndex + 1

	if adjustedCursor > viewportBottom {
		newOffset := max(adjustedCursor-m.categoriesViewport.Height+1, 0)
		m.categoriesViewport.SetYOffset(newOffset)
	}
	if adjustedCursor < viewportTop {
		m.categoriesViewport.SetYOffset(adjustedCursor)
	}
	return m
}

// updateViewportHeight updates the viewport heights based on current window size.
func (m MonthlyModel) updateViewportHeight() MonthlyModel {
	if !m.ready {
		return m
	}

	headerHeight := lipgloss.Height(m.getHeader(decimal.Zero, ""))
	footerHeight := lipgloss.Height(m.getFooter(decimal.Zero, decimal.Zero, ""))
	verticalMarginHeight := headerHeight + footerHeight
	availableHeight := m.Height - verticalMarginHeight - 4 // -4 for padding (2) and newlines (2)

	groupsViewportHeight := m.calculateGroupsViewportHeight(availableHeight)
	categoriesViewportHeight := m.calculateCategoriesViewportHeight(availableHeight)

	m.groupsViewport.Height = groupsViewportHeight
	m.categoriesViewport.Height = categoriesViewportHeight
	return m
}

// getMonthIncome calculates the total income for the month.
func (m MonthlyModel) getMonthIncome() decimal.Decimal {
	var totalIncome decimal.Decimal
	for _, income := range m.incomes {
		amount := decimal.NewFromFloat(income.Amount)
		totalIncome = totalIncome.Add(amount)
	}
	return totalIncome
}

// getMonthExpenses calculates total expenses and group totals for the month.
func (m MonthlyModel) getMonthExpenses() (decimal.Decimal, map[string]decimal.Decimal) {
	var expenseTotals decimal.Decimal
	groupTotals := make(map[string]decimal.Decimal)

	for _, category := range m.categories {
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

// getHeader renders the header section with month/year and total income.
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
		CreateSpacer(headerSpacerWidth).Render(""),
		headerRight,
	)

	bottomBorder := TopBorder.
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

// getFooter renders the footer section with totals, balance, and key hints.
func (m MonthlyModel) getFooter(totalExpenses, balance decimal.Decimal, defaultCurrency string) string {
	var b bytes.Buffer

	footerStyle := CreateFooterStyle(m.Width)

	var keyHints string
	populateHint := ""
	if m.currentMonthHasNoCategories() {
		populateHint = " | p: Populate"
	}

	resetHint := ""
	if !m.isViewingCurrentMonth() {
		resetHint = " | r: Reset"
	}

	switch m.Level {
	case focusLevelGroups:
		keyHints = "j/k: Nav | Ent: Select" + populateHint + " | i: Income | c: Categories | g: Groups | h/l: Month" + resetHint
	case focusLevelCategories:
		keyHints = "j/k: Nav | Ent: Expense | t: Toggle | Esc: Back" + populateHint + " | i: Income | c: Categories | g: Groups | h/l: Month" + resetHint
	}
	totalExpensesStr := fmt.Sprintf("Total Expenses: %s %s", totalExpenses.String(), defaultCurrency)

	balanceStr := fmt.Sprintf("Balance: %s %s", balance.String(), defaultCurrency)
	footerSummarySpacerWidth := max(m.Width-lipgloss.Width(totalExpensesStr)-lipgloss.Width(balanceStr)-AppStyle.GetHorizontalPadding(), 0)

	space := CreateSpacer(footerSummarySpacerWidth).Render("")
	footerSummary := lipgloss.JoinHorizontal(lipgloss.Top, totalExpensesStr, space, balanceStr)

	b.WriteString("\n\n")
	b.WriteString(footerStyle.Render(lipgloss.JoinVertical(lipgloss.Left, footerSummary, "", MutedText.Render(keyHints))))

	return b.String()
}

// getOrderedGroups returns all category groups sorted by order.
func (m MonthlyModel) getOrderedGroups() []domain.CategoryGroup {
	var orderedGroups []domain.CategoryGroup
	for _, group := range m.categoryGroups {
		orderedGroups = append(orderedGroups, group)
	}

	sort.Slice(orderedGroups, func(i, j int) bool {
		return orderedGroups[i].Order < orderedGroups[j].Order
	})

	return orderedGroups
}

// getVisibleGroups filters groups to only those containing categories.
func (m MonthlyModel) getVisibleGroups(orderedGroups []domain.CategoryGroup, categoriesByGroup map[string][]domain.Category) []domain.CategoryGroup {
	var visibleGroups []domain.CategoryGroup
	for _, group := range orderedGroups {
		if _, hasCategories := categoriesByGroup[group.GroupID]; hasCategories {
			visibleGroups = append(visibleGroups, group)
		}
	}
	return visibleGroups
}

// handleGroupNavigation processes keyboard input when navigating groups.
func (m MonthlyModel) handleGroupNavigation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Group categories by their GroupID
	categoriesByGroup := make(map[string][]domain.Category)
	for _, category := range m.categories {
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
			m = m.ensureGroupsCursorVisible()
		}
	case "k", "up":
		if numVisibleGroups > 0 {
			m.focusedGroupIndex--
			if m.focusedGroupIndex < 0 {
				m.focusedGroupIndex = numVisibleGroups - 1
			}
			m = m.ensureGroupsCursorVisible()
		}
	case "enter":
		if numVisibleGroups > 0 && m.focusedGroupIndex >= 0 && m.focusedGroupIndex < numVisibleGroups {
			// The focused index now directly maps to visible groups
			m.Level = focusLevelCategories
			m.focusedCategoryIndex = 0
			m = m.ensureCategoriesCursorVisible()
		}

	}
	return m, nil
}

// handleCategoryNavigation processes keyboard input when navigating categories.
func (m MonthlyModel) handleCategoryNavigation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Group categories by their GroupID
	categoriesByGroup := make(map[string][]domain.Category)
	for _, category := range m.categories {
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
			m = m.ensureCategoriesCursorVisible()
		}
	case "k", "up":
		if numCategories > 0 {
			m.focusedCategoryIndex--
			if m.focusedCategoryIndex < 0 {
				m.focusedCategoryIndex = numCategories - 1
			}
			m = m.ensureCategoriesCursorVisible()
		}
	case "enter":
		if numCategories > 0 && m.focusedCategoryIndex >= 0 && m.focusedCategoryIndex < numCategories {
			selectedCategory := categoriesInGroup[m.focusedCategoryIndex]
			monthKey := GetMonthKey(m.CurrentMonth, m.CurrentYear)
			return m, func() tea.Msg {
				return ExpenseViewMsg{
					MonthKey: monthKey,
					Category: selectedCategory,
				}
			}
		}
	case "t":
		// Toggle expense status for selected category
		if numCategories > 0 && m.focusedCategoryIndex >= 0 && m.focusedCategoryIndex < numCategories {
			selectedCategory := categoriesInGroup[m.focusedCategoryIndex]
			monthKey := GetMonthKey(m.CurrentMonth, m.CurrentYear)
			return m, func() tea.Msg {
				return ToggleExpenseStatusMsg{
					MonthKey: monthKey,
					Category: selectedCategory,
				}
			}
		}
	case "esc":
		// Go back to group navigation
		m.Level = focusLevelGroups
		m.focusedCategoryIndex = 0
		m = m.ensureGroupsCursorVisible()

	}
	return m, nil
}

// handlePopulateCategories initiates populating categories from previous month.
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

// currentMonthHasNoCategories checks if the current month has any categories.
func (m MonthlyModel) currentMonthHasNoCategories() bool {
	return len(m.categories) == 0
}

// SetMonthYear updates the current month/year and resets focus.
func (m MonthlyModel) SetMonthYear(month time.Month, year int) MonthlyModel {
	m.CurrentMonth = month
	m.CurrentYear = year
	return m.ResetFocus()
}

// ResetFocus resets the navigation focus to the top level.
func (m MonthlyModel) ResetFocus() MonthlyModel {
	m.Level = focusLevelGroups
	m.focusedGroupIndex = 0
	m.focusedCategoryIndex = 0
	return m
}

// SetFocusToCategory sets focus to a specific category in the navigation.
func (m MonthlyModel) SetFocusToCategory(category domain.Category) MonthlyModel {
	// Group categories by their GroupID
	categoriesByGroup := make(map[string][]domain.Category)
	for _, cat := range m.categories {
		categoriesByGroup[cat.GroupID] = append(categoriesByGroup[cat.GroupID], cat)
	}

	// Get ordered and visible groups
	orderedGroups := m.getOrderedGroups()
	visibleGroups := m.getVisibleGroups(orderedGroups, categoriesByGroup)

	// Find the group index for this category
	groupIndex := -1
	for i, group := range visibleGroups {
		if group.GroupID == category.GroupID {
			groupIndex = i
			break
		}
	}

	if groupIndex == -1 {
		return m.ResetFocus()
	}

	// Find the category index within the group
	categoryIndex := -1
	categoriesInGroup := categoriesByGroup[category.GroupID]
	for i, cat := range categoriesInGroup {
		if cat.CatID == category.CatID {
			categoryIndex = i
			break
		}
	}

	if categoryIndex == -1 {
		return m.ResetFocus()
	}

	// Set focus to the found category
	m.Level = focusLevelCategories
	m.focusedGroupIndex = groupIndex
	m.focusedCategoryIndex = categoryIndex
	return m
}

// UpdateData refreshes the model with new data.
func (m MonthlyModel) UpdateData(appData AppData) MonthlyModel {
	m.categories = appData.Categories
	m.categoryGroups = appData.CategoryGroups
	m.incomes = appData.Incomes

	// Reset focus indices if they're out of bounds
	// Group categories by their GroupID
	categoriesByGroup := make(map[string][]domain.Category)
	for _, category := range m.categories {
		categoriesByGroup[category.GroupID] = append(categoriesByGroup[category.GroupID], category)
	}

	// Get visible groups
	orderedGroups := m.getOrderedGroups()
	visibleGroups := m.getVisibleGroups(orderedGroups, categoriesByGroup)

	if m.focusedGroupIndex >= len(visibleGroups) && len(visibleGroups) > 0 {
		m.focusedGroupIndex = len(visibleGroups) - 1
	} else if len(visibleGroups) == 0 {
		m.focusedGroupIndex = 0
	}

	// Reset category focus if needed
	if len(visibleGroups) > 0 && m.focusedGroupIndex < len(visibleGroups) {
		selectedGroup := visibleGroups[m.focusedGroupIndex]
		categoriesInGroup := categoriesByGroup[selectedGroup.GroupID]
		if m.focusedCategoryIndex >= len(categoriesInGroup) && len(categoriesInGroup) > 0 {
			m.focusedCategoryIndex = len(categoriesInGroup) - 1
		} else if len(categoriesInGroup) == 0 {
			m.focusedCategoryIndex = 0
		}
	}

	// Update viewport heights and content
	m = m.updateViewportHeight()
	if m.ready {
		switch m.Level {
		case focusLevelGroups:
			m = m.ensureGroupsCursorVisible()
		case focusLevelCategories:
			m = m.ensureCategoriesCursorVisible()
		}
	}

	return m
}

// isViewingCurrentMonth checks if the currently viewed month is the actual current month
func (m MonthlyModel) isViewingCurrentMonth() bool {
	now := time.Now()
	return m.CurrentMonth == now.Month() && m.CurrentYear == now.Year()
}
