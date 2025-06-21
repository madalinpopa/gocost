package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/madalinpopa/gocost/internal/domain"
)

// CategoryModel represents the view model that displays the categories
type CategoryModel struct {
	WindowSize
	MonthYear

	MonthKey string
	cursor   int

	categories     []domain.Category
	categoryGroups []domain.CategoryGroup

	addCategory    bool                 // When adding a new category
	moveCategory   bool                 // When moving a category to a new group
	selectedGroup  domain.CategoryGroup // Selected group
	movingCategory domain.Category      // Category being moved

	isEditingName bool
	editInput     textinput.Model
	editingIndex  int

	isFiltering        bool
	filterInput        textinput.Model
	filterText         string
	filteredCategories []domain.Category
	isFiltered         bool

	viewport viewport.Model
	ready    bool
}

// NewCategoryModel creates a new CategoryModel instance.
func NewCategoryModel(appData AppData, month time.Month, year int) CategoryModel {
	monthKey := GetMonthKey(month, year)

	ti := textinput.New()
	ti.Placeholder = "Category name"
	ti.CharLimit = 30
	ti.Width = 30

	filterTi := textinput.New()
	filterTi.Placeholder = "Filter categories..."
	filterTi.CharLimit = 50
	filterTi.Width = 50

	return CategoryModel{
		MonthKey:       monthKey,
		categories:     appData.Categories,
		categoryGroups: appData.CategoryGroups,
		editInput:      ti,
		editingIndex:   -1,
		filterInput:    filterTi,
		viewport:       viewport.New(80, 20),
		ready:          false,
	}
}

// Init initializes the CategoryModel.
func (m CategoryModel) Init() tea.Cmd {
	return nil
}

// Update updates the CategoryModel.
func (m CategoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.isFiltering {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				filterText := strings.TrimSpace(m.filterInput.Value())
				m.isFiltering = false
				m.filterInput.Blur()

				if filterText != "" {
					return m, func() tea.Msg {
						return FilterCategoriesMsg{FilterText: filterText}
					}
				} else {
					m.filterText = ""
					m.filteredCategories = nil
					m.isFiltered = false
					m.cursor = 0
					m.filterInput.SetValue("")
					return m, nil
				}
			case "esc":
				m.isFiltering = false
				m.filterInput.Blur()
				m.filterInput.SetValue("")
				return m, nil
			}
		}
		m.filterInput, cmd = m.filterInput.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	if m.addCategory {
		switch msg := msg.(type) {
		case FilterCategoriesMsg:
			m = m.handleFilterCategories(msg)
			return m, nil
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				categoryName := strings.TrimSpace(m.editInput.Value())
				if categoryName != "" {
					newCategoryId := GenerateID()
					newCategory := domain.Category{
						CatID:        newCategoryId,
						GroupID:      m.selectedGroup.GroupID,
						CategoryName: categoryName,
						Expense:      make(map[string]domain.ExpenseRecord, 0),
					}
					m.addCategory = false
					m.editInput.SetValue("")
					m.editInput.Blur()
					return m, func() tea.Msg {
						return CategoryAddMsg{MonthKey: m.MonthKey, Category: newCategory}
					}
				}
			case "esc":
				m.addCategory = false
				m.editInput.Blur()
				m.editInput.SetValue("")
				return m, nil
			}
		}

		m.editInput, cmd = m.editInput.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	if m.isEditingName {
		switch msg := msg.(type) {
		case FilterCategoriesMsg:
			m = m.handleFilterCategories(msg)
			return m, nil
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				categoryName := strings.TrimSpace(m.editInput.Value())
				if categoryName != "" {
					updatedCategory := m.categories[m.editingIndex]
					updatedCategory.CategoryName = categoryName
					m.isEditingName = false
					m.editInput.Blur()
					return m, func() tea.Msg {
						return CategoryUpdateMsg{MonthKey: m.MonthKey, Category: updatedCategory}
					}
				}
			case "esc":
				m.isEditingName = false
				m.editInput.Blur()
				m.editInput.SetValue("")
				return m, nil
			}
		}
		m.editInput, cmd = m.editInput.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case FilterCategoriesMsg:
		m = m.handleFilterCategories(msg)
		return m, nil
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			viewportHeight := max(1, msg.Height-verticalMarginHeight-6) // -6 for padding
			m.viewport = viewport.New(msg.Width, viewportHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.getCategoriesContent())
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			viewportHeight := max(1, msg.Height-verticalMarginHeight-6) // -6 for padding
			m.viewport.Height = viewportHeight
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			if m.IsMovingCategory() {
				m = m.ResetMoveState()
				return m, nil
			}
			m = m.resetEditingState()
			return m, func() tea.Msg { return MonthlyViewMsg{} }

		case "j", "down":
			displayCategories := m.getDisplayCategories()
			if len(displayCategories) > 0 {
				m.cursor = (m.cursor + 1) % len(displayCategories)
				(&m).ensureCursorVisible()
			}
			return m, nil

		case "k", "up":
			displayCategories := m.getDisplayCategories()
			if len(displayCategories) > 0 {
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(displayCategories) - 1
				}
				(&m).ensureCursorVisible()
			}
			return m, nil
		case "/":
			m.isFiltering = true
			m.filterInput.Focus()
			return m, textinput.Blink
		case "c":
			if m.isFiltered {
				m.filterText = ""
				m.filteredCategories = nil
				m.isFiltered = false
				m.cursor = 0
				(&m).ensureCursorVisible()
				return m, nil
			}
		case "a", "n":
			return m, func() tea.Msg { return SelectGroupMsg{} }
		case "e":
			displayCategories := m.getDisplayCategories()
			if len(displayCategories) > 0 {
				if m.cursor >= 0 && m.cursor < len(displayCategories) {
					selectedCategory := displayCategories[m.cursor]
					for i, category := range m.categories {
						if category.CatID == selectedCategory.CatID {
							m.editingIndex = i
							break
						}
					}
					m.editInput.SetValue(selectedCategory.CategoryName)
					m.editInput.Placeholder = "Edit Category Name"
					return m.focusInput()
				}
			}
		case "d":
			displayCategories := m.getDisplayCategories()
			if len(displayCategories) > 0 {
				if m.cursor >= 0 && m.cursor < len(displayCategories) {
					selectedCategory := displayCategories[m.cursor]
					return m, func() tea.Msg { return CategoryDeleteMsg{MonthKey: m.MonthKey, Category: selectedCategory} }
				}
			}
		case "m":
			displayCategories := m.getDisplayCategories()
			if len(displayCategories) > 0 {
				if m.cursor >= 0 && m.cursor < len(displayCategories) {
					m.movingCategory = displayCategories[m.cursor]
					return m, func() tea.Msg { return SelectGroupMsg{} }
				}
			}
		}
	}

	if m.ready && !m.isFiltering && !m.addCategory && !m.isEditingName {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the CategoryModel.
func (m CategoryModel) View() string {
	if !m.ready {
		return AppStyle.Width(m.Width).Height(m.Height).Render("\n  Initializing...")
	}

	if m.ready && !m.isFiltering && !m.addCategory && !m.isEditingName {
		m.viewport.SetContent(m.getCategoriesContent())
	}

	if m.isFiltering || m.addCategory || m.isEditingName {
		var b strings.Builder
		b.WriteString(m.headerView())
		b.WriteString("\n")
		b.WriteString(m.footerView())
		return AppStyle.Render(b.String())
	}

	var b strings.Builder
	b.WriteString(m.headerView())
	b.WriteString("\n")
	b.WriteString(m.viewport.View())
	b.WriteString("\n")
	b.WriteString(m.footerView())
	return AppStyle.Render(b.String())
}

// getGroupName finds the group name for a given group ID.
func (m CategoryModel) getGroupName(groupID string) string {
	for _, group := range m.categoryGroups {
		if group.GroupID == groupID {
			return group.GroupName
		}
	}
	return "Unknown"
}

// AddCategory enables input field to create a new category in the specified group.
func (m CategoryModel) AddCategory(group domain.CategoryGroup) (CategoryModel, tea.Cmd) {
	m.addCategory = true
	m.selectedGroup = group
	m.editInput.Focus()
	return m, textinput.Blink
}

// MoveCategory moves the selected category to the specified group.
func (m CategoryModel) MoveCategory(group domain.CategoryGroup) (CategoryModel, tea.Cmd) {
	m.moveCategory = true
	m.selectedGroup = group
	m.movingCategory.GroupID = group.GroupID
	m.moveCategory = false
	return m, func() tea.Msg {
		return CategoryUpdateMsg{MonthKey: m.MonthKey, Category: m.movingCategory}
	}
}

// IsMovingCategory returns true if a category is currently being moved.
func (m CategoryModel) IsMovingCategory() bool {
	return m.movingCategory.CatID != ""
}

// ResetMoveState clears the moving category state.
func (m CategoryModel) ResetMoveState() CategoryModel {
	m.movingCategory = domain.Category{}
	m.moveCategory = false
	return m
}

// getDisplayCategories returns either filtered categories or all categories.
func (m CategoryModel) getDisplayCategories() []domain.Category {
	if m.isFiltered {
		return m.filteredCategories
	}
	return m.categories
}

// ensureCursorVisible updates viewport content and scrolls to keep cursor visible
func (m *CategoryModel) ensureCursorVisible() {
	if !m.ready {
		return
	}

	content := m.getCategoriesContent()
	m.viewport.SetContent(content)

	displayCategories := m.getDisplayCategories()
	if len(displayCategories) == 0 {
		return
	}

	// Calculate the current viewport bounds
	viewportTop := m.viewport.YOffset
	viewportBottom := viewportTop + m.viewport.Height - 1

	// Adjust for filter info line if present
	contentOffset := 0
	if m.isFiltered {
		contentOffset = 1 // Account for the filter info line
	}

	// Adjust cursor position for content offset
	adjustedCursor := m.cursor + contentOffset

	// If cursor is below viewport, scroll down
	if adjustedCursor > viewportBottom {
		newOffset := max(adjustedCursor-m.viewport.Height+1, 0)
		m.viewport.SetYOffset(newOffset)
	}
	// If cursor is above viewport, scroll up
	if adjustedCursor < viewportTop {
		m.viewport.SetYOffset(adjustedCursor)
	}
}

// handleFilterCategories processes the filter message and applies filtering.
func (m CategoryModel) handleFilterCategories(msg FilterCategoriesMsg) CategoryModel {
	m.filterText = msg.FilterText
	var filtered []domain.Category
	filterLower := strings.ToLower(msg.FilterText)
	for _, category := range m.categories {
		if strings.Contains(strings.ToLower(category.CategoryName), filterLower) {
			filtered = append(filtered, category)
			continue
		}
		groupName := m.getGroupName(category.GroupID)
		if strings.Contains(strings.ToLower(groupName), filterLower) {
			filtered = append(filtered, category)
		}
	}
	m.filteredCategories = filtered
	m.isFiltered = true
	m.cursor = 0
	m.filterInput.SetValue("")
	(&m).ensureCursorVisible()
	return m
}

// clearFilter clears the filter state.
func (m CategoryModel) clearFilter() CategoryModel {
	m.isFiltering = false
	m.filterText = ""
	m.filteredCategories = nil
	m.isFiltered = false
	m.filterInput.SetValue("")
	m.filterInput.Blur()
	return m
}

// resetEditingState resets all editing-related state flags and inputs.
func (m CategoryModel) resetEditingState() CategoryModel {
	m.addCategory = false
	m.moveCategory = false
	m.selectedGroup = domain.CategoryGroup{}
	m.movingCategory = domain.Category{}
	m.isEditingName = false
	m.editInput.Blur()
	m.editInput.SetValue("")
	m.editingIndex = -1
	m = m.clearFilter()
	return m
}

// focusInput activates the text input for category name editing.
func (m CategoryModel) focusInput() (tea.Model, tea.Cmd) {
	m.isEditingName = true
	m.editInput.Focus()
	return m, textinput.Blink
}

// headerView renders the header section
func (m CategoryModel) headerView() string {
	var b strings.Builder

	title := "Manage Expense Categories"
	if m.isFiltering {
		title = "Filter Categories"
	} else if m.isEditingName {
		title = "Edit Category Name"
	} else if m.addCategory {
		title = fmt.Sprintf("Add Category to %s", m.selectedGroup.GroupName)
	} else if m.IsMovingCategory() {
		title = fmt.Sprintf("Moving: %s", m.movingCategory.CategoryName)
	}

	b.WriteString(HeaderText.Render(title))
	b.WriteString("\n")

	if m.isFiltering {
		b.WriteString("\n")
		b.WriteString("Filter Categories (Enter to apply filter, Esc to cancel):\n")
		b.WriteString(m.filterInput.View())
	} else if m.isEditingName {
		b.WriteString("\n")
		b.WriteString("Enter Category Name (Enter to save, Esc to cancel):\n")
		b.WriteString(m.editInput.View())
	} else if m.addCategory {
		b.WriteString("\n")
		b.WriteString("Enter Category Name (Enter to save, Esc to cancel):\n")
		b.WriteString(m.editInput.View())
	} else if m.IsMovingCategory() {
		b.WriteString("\n")
		b.WriteString(MutedText.Render(fmt.Sprintf("Select a new group for category '%s'", m.movingCategory.CategoryName)))
	}

	return b.String()
}

// footerView renders the footer section
func (m CategoryModel) footerView() string {
	var b strings.Builder
	b.WriteString("\n")
	keyHints := "(j/k: Nav, /: Filter, a/n: Add, e: Edit, d: Delete, m: Move, Esc/q: Back)"
	if m.IsMovingCategory() {
		keyHints = "(Select a group to move the category, Esc/q: Cancel)"
	}
	if m.isFiltered {
		keyHints = "(j/k: Nav, /: Filter, a/n: Add, e: Edit, d: Delete, m: Move, c: Clear filter, Esc/q: Back)"
	}
	b.WriteString(MutedText.Render(keyHints))
	return b.String()
}

// getCategoriesContent generates the content for the viewport
func (m CategoryModel) getCategoriesContent() string {
	var b strings.Builder
	displayCategories := m.getDisplayCategories()

	if len(displayCategories) == 0 {
		if m.isFiltered {
			b.WriteString(MutedText.Render(fmt.Sprintf("No categories found matching '%s'.", m.filterText)))
		} else {
			b.WriteString(MutedText.Render("No category defined yet."))
		}
	} else {
		if m.isFiltered {
			b.WriteString(MutedText.Render(fmt.Sprintf("Showing %d of %d categories matching '%s'", len(displayCategories), len(m.categories), m.filterText)))
			b.WriteString("\n")
		}
		maxCategoryWidth := 0
		maxGroupWidth := 0
		for _, item := range displayCategories {
			if len(item.CategoryName) > maxCategoryWidth {
				maxCategoryWidth = len(item.CategoryName)
			}
			groupName := m.getGroupName(item.GroupID)
			if len(groupName) > maxGroupWidth {
				maxGroupWidth = len(groupName)
			}
		}
		categoryColWidth := maxCategoryWidth + 2
		groupColWidth := maxGroupWidth + 2
		for i, item := range displayCategories {
			style := NormalListItem
			prefix := " "
			if i == m.cursor {
				style = FocusedListItem
				prefix = ">"
			}
			if m.IsMovingCategory() && item.CatID == m.movingCategory.CatID {
				style = FocusedListItem
				prefix = "→ "
			}
			groupName := m.getGroupName(item.GroupID)
			categoryFormatted := fmt.Sprintf("%-*s", categoryColWidth, item.CategoryName)
			groupFormatted := MutedText.Render(fmt.Sprintf("%-*s", groupColWidth, groupName))
			line := fmt.Sprintf("%s %s %s", prefix, categoryFormatted, groupFormatted)
			b.WriteString(style.Render(line))
			b.WriteString("\n")
		}
	}
	return b.String()
}

// UpdateData refreshes the model with new data and resets state.
func (m CategoryModel) UpdateData(appData AppData) CategoryModel {
	m.categories = appData.Categories
	m.categoryGroups = appData.CategoryGroups
	m.cursor = 0
	m = m.resetEditingState()
	(&m).ensureCursorVisible()
	return m
}

// SetMonthYear updates the current month/year and loads corresponding categories.
func (m CategoryModel) SetMonthYear(month time.Month, year int) CategoryModel {
	m.CurrentMonth = month
	m.CurrentYear = year
	m.MonthKey = GetMonthKey(month, year)
	m = m.resetEditingState()
	(&m).ensureCursorVisible()
	return m
}
