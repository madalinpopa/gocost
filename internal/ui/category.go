package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/data"
)

// CategoryModel represents the view model that displays the categories
type CategoryModel struct {
	AppData
	WindowSize
	MonthYear

	MonthKey   string
	cursor     int
	categories []data.Category

	addCategory    bool               // When adding a new category
	moveCategory   bool               // When moving a category to a new group
	selectedGroup  data.CategoryGroup // Selected group
	movingCategory data.Category      // Category being moved

	isEditingName bool
	editInput     textinput.Model
	editingIndex  int

	isFiltering        bool
	filterInput        textinput.Model
	filterText         string
	filteredCategories []data.Category
	isFiltered         bool
}

// NewCategoryModel creates a new CategoryModel instance.
func NewCategoryModel(initialData *data.DataRoot, month time.Month, year int) CategoryModel {

	monthKey := GetMonthKey(month, year)

	ti := textinput.New()
	ti.Placeholder = "Category name"
	ti.CharLimit = 30
	ti.Width = 30

	filterTi := textinput.New()
	filterTi.Placeholder = "Filter categories..."
	filterTi.CharLimit = 50
	filterTi.Width = 50

	var categories []data.Category
	if record, ok := initialData.MonthlyData[monthKey]; ok {
		categories = record.Categories
	}

	if categories == nil {
		categories = make([]data.Category, 0)
	}

	return CategoryModel{
		AppData: AppData{
			Data: initialData,
		},
		MonthKey:     monthKey,
		categories:   categories,
		editInput:    ti,
		editingIndex: -1,
		filterInput:  filterTi,
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
					// Clear filter if empty
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

					newCategory := data.Category{
						CatID:        newCategoryId,
						GroupID:      m.selectedGroup.GroupID,
						CategoryName: categoryName,
						Expense:      make(map[string]data.ExpenseRecord, 0),
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
		cmds := append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	if m.isEditingName {

		// Handle actions when editing
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
				// Cancel move operation
				m = m.ResetMoveState()
				return m, nil
			}
			// Clear filter when exiting
			m = m.clearFilter()
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
				// Clear filter
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
					// Find the original index in the full categories list
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
			// Show filter info if filtering is active
			if m.isFiltered {
				b.WriteString(MutedText.Render(fmt.Sprintf("Showing %d of %d categories matching '%s'", len(displayCategories), len(m.categories), m.filterText)))
				b.WriteString("\n\n")
			}
			
			// Calculate maximum widths for dynamic column sizing
			maxCategoryWidth := 0
			maxGroupWidth := 0
			
			for _, item := range displayCategories {
				if len(item.CategoryName) > maxCategoryWidth {
					maxCategoryWidth = len(item.CategoryName)
				}
				
				var groupName string
				if group, ok := m.Data.CategoryGroups[item.GroupID]; ok {
					groupName = group.GroupName
				}
				if len(groupName) > maxGroupWidth {
					maxGroupWidth = len(groupName)
				}
			}
			
			// Add padding to column widths
			categoryColWidth := maxCategoryWidth + 2
			groupColWidth := maxGroupWidth + 2
			
			for i, item := range displayCategories {
				style := NormalListItem
				prefix := " "
				if i == m.cursor {
					style = FocusedListItem
					prefix = ">"
				}

				// Highlight the category being moved
				if m.IsMovingCategory() && item.CatID == m.movingCategory.CatID {
					style = FocusedListItem
					prefix = "→ "
				}

				var groupName string
				group, ok := m.Data.CategoryGroups[item.GroupID]
				if ok {
					groupName = group.GroupName
				}
				
				// Format category name with left alignment
				categoryFormatted := fmt.Sprintf("%-*s", categoryColWidth, item.CategoryName)
				
				// Format group name with muted style and column alignment
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

// AddCategory enables input field to create a new category in the specified group.
func (m CategoryModel) AddCategory(group data.CategoryGroup) (CategoryModel, tea.Cmd) {
	m.addCategory = true
	m.selectedGroup = group
	m.editInput.Focus()
	return m, textinput.Blink
}

// MoveCategory moves the selected category to the specified group.
func (m CategoryModel) MoveCategory(group data.CategoryGroup) (CategoryModel, tea.Cmd) {
	m.moveCategory = true
	m.selectedGroup = group

	// Update the moving category's group ID
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
	m.movingCategory = data.Category{}
	m.moveCategory = false
	return m
}

// getDisplayCategories returns either filtered categories or all categories
func (m CategoryModel) getDisplayCategories() []data.Category {
	if m.isFiltered {
		return m.filteredCategories
	}
	return m.categories
}

// handleFilterCategories processes the filter message and applies filtering
func (m CategoryModel) handleFilterCategories(msg FilterCategoriesMsg) (CategoryModel, tea.Cmd) {
	m.filterText = msg.FilterText
	
	var filtered []data.Category
	filterLower := strings.ToLower(msg.FilterText)

	for _, category := range m.categories {
		// Check if category name contains filter text
		if strings.Contains(strings.ToLower(category.CategoryName), filterLower) {
			filtered = append(filtered, category)
			continue
		}

		// Check if group name contains filter text
		if group, ok := m.Data.CategoryGroups[category.GroupID]; ok {
			if strings.Contains(strings.ToLower(group.GroupName), filterLower) {
				filtered = append(filtered, category)
			}
		}
	}

	m.filteredCategories = filtered
	m.isFiltered = true
	m.cursor = 0 // Reset cursor to first item
	m.filterInput.SetValue("")
	
	return m, nil
}

// clearFilter clears the filter state
func (m CategoryModel) clearFilter() CategoryModel {
	m.isFiltering = false
	m.filterText = ""
	m.filteredCategories = nil
	m.isFiltered = false
	m.filterInput.SetValue("")
	m.filterInput.Blur()
	return m
}

// focusInput activates the text input for category name editing.
func (m CategoryModel) focusInput() (tea.Model, tea.Cmd) {
	m.isEditingName = true
	m.editInput.Focus()
	return m, textinput.Blink
}

// UpdateData refreshes the model with new data and resets state.
func (m CategoryModel) UpdateData(updatedData *data.DataRoot) CategoryModel {
	m.Data = updatedData

	var categories []data.Category
	if record, ok := m.Data.MonthlyData[m.MonthKey]; ok {
		categories = record.Categories
	}

	m.categories = categories
	
	// Clear filter when data is updated
	m.filterText = ""
	m.filteredCategories = nil
	m.isFiltered = false
	m.cursor = 0

	// Reset move state when data is updated
	m = m.ResetMoveState()

	return m
}

// SetMonthYear updates the current month/year and loads corresponding categories.
func (m CategoryModel) SetMonthYear(month time.Month, year int) CategoryModel {
	m.CurrentMonth = month
	m.CurrentYear = year
	m.MonthKey = GetMonthKey(month, year)

	var categories []data.Category
	if record, ok := m.Data.MonthlyData[m.MonthKey]; ok {
		categories = record.Categories
	}

	if categories == nil {
		categories = make([]data.Category, 0)
	}

	m.categories = categories
	
	// Reset cursor
	if len(m.categories) > 0 {
		m.cursor = 0
	} else {
		m.cursor = 0
	}

	// Reset move state and filter when month/year changes
	m = m.ResetMoveState()
	m = m.clearFilter()

	return m
}
