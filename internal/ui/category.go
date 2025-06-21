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
	// Data is now passed in as slices
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
			return m.handleFilterCategories(msg)
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
			return m.handleFilterCategories(msg)
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
		return m.handleFilterCategories(msg)
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
			}
			return m, nil

		case "k", "up":
			displayCategories := m.getDisplayCategories()
			if len(displayCategories) > 0 {
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(displayCategories) - 1
				}
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
	return m, nil
}

// View renders the CategoryModel.
func (m CategoryModel) View() string {
	var b strings.Builder
	title := "Manage Expense Categories"
	if m.IsMovingCategory() {
		title = fmt.Sprintf("Moving: %s", m.movingCategory.CategoryName)
	}
	b.WriteString(HeaderText.Render(title))
	b.WriteString("\n\n")

	if m.isFiltering {
		b.WriteString("Filter Categories (Enter to apply filter, Esc to cancel):\n")
		b.WriteString(m.filterInput.View())
		b.WriteString("\n\n")
	} else if m.isEditingName || m.addCategory {
		b.WriteString("Enter Category Name (Enter to save, Esc to cancel):\n")
		b.WriteString(m.editInput.View())
		b.WriteString("\n")
	} else {
		if m.IsMovingCategory() {
			b.WriteString(MutedText.Render(fmt.Sprintf("Select a new group for category '%s'", m.movingCategory.CategoryName)))
			b.WriteString("\n\n")
		}
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
				b.WriteString("\n\n")
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
		b.WriteString("\n\n")
		keyHints := "(j/k: Nav, /: Filter, a/n: Add, e: Edit, d: Delete, m: Move, Esc/q: Back)"
		if m.IsMovingCategory() {
			keyHints = "(Select a group to move the category, Esc/q: Cancel)"
		}
		if m.isFiltered {
			keyHints = "(j/k: Nav, /: Filter, a/n: Add, e: Edit, d: Delete, m: Move, c: Clear filter, Esc/q: Back)"
		}
		b.WriteString(MutedText.Render(keyHints))
	}
	viewStr := AppStyle.Width(m.Width).Height(m.Height - 3).Render(b.String())
	return viewStr
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

// handleFilterCategories processes the filter message and applies filtering.
func (m CategoryModel) handleFilterCategories(msg FilterCategoriesMsg) (CategoryModel, tea.Cmd) {
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
	return m, nil
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

// UpdateData refreshes the model with new data and resets state.
func (m CategoryModel) UpdateData(appData AppData) CategoryModel {
	m.categories = appData.Categories
	m.categoryGroups = appData.CategoryGroups
	m.cursor = 0
	m = m.resetEditingState()
	return m
}

// SetMonthYear updates the current month/year and loads corresponding categories.
func (m CategoryModel) SetMonthYear(month time.Month, year int) CategoryModel {
	m.CurrentMonth = month
	m.CurrentYear = year
	m.MonthKey = GetMonthKey(month, year)
	m = m.resetEditingState()
	return m
}
