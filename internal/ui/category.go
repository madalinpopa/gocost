package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
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

	addCategory   bool               // When adding a new category
	moveCategory  bool               // When moving a category to a new group
	selectedGroup data.CategoryGroup // Selected group
	movingCategory data.Category     // Category being moved

	isEditingName bool
	editInput     textinput.Model
	editingIndex  int
}

func NewCategoryModel(initialData *data.DataRoot, month time.Month, year int) CategoryModel {

	monthKey := GetMonthKey(month, year)

	ti := textinput.New()
	ti.Placeholder = "Category name"
	ti.CharLimit = 30
	ti.Width = 30

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
	}

}

func (m CategoryModel) Init() tea.Cmd {
	return nil
}

func (m CategoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.addCategory {

		switch msg := msg.(type) {

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

	case tea.KeyMsg:

		switch msg.String() {

		case "q", "esc":
			if m.IsMovingCategory() {
				// Cancel move operation
				m = m.ResetMoveState()
				return m, nil
			}
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

		case "a", "n":
			return m, func() tea.Msg { return SelectGroupMsg{} }
		case "e":
			if len(m.categories) > 0 {
				if m.cursor >= 0 && m.cursor < len(m.categories) {
					m.editingIndex = m.cursor
					m.editInput.SetValue(m.categories[m.cursor].CategoryName)
					m.editInput.Placeholder = "Edit Category Name"
					return m.focusInput()
				}
			}

		case "d":
			if len(m.categories) > 0 {
				if m.cursor >= 0 && m.cursor < len(m.categories) {
					selectedCategory := m.categories[m.cursor]
					return m, func() tea.Msg { return CategoryDeleteMsg{MonthKey: m.MonthKey, Category: selectedCategory} }
				}
			}

		case "m":
			if len(m.categories) > 0 {
				if m.cursor >= 0 && m.cursor < len(m.categories) {
					m.movingCategory = m.categories[m.cursor]
					return m, func() tea.Msg { return SelectGroupMsg{} }
				}
			}
		}
	}

	return m, nil
}

func (m CategoryModel) View() string {
	var b strings.Builder

	title := "Manage Expense Categories"
	if m.IsMovingCategory() {
		title = fmt.Sprintf("Moving: %s", m.movingCategory.CategoryName)
	}
	b.WriteString(HeaderText.Render(title))
	b.WriteString("\n\n")

	if m.isEditingName || m.addCategory {
		b.WriteString("Enter Category Name (Enter to save, Esc to cancel):\n")
		b.WriteString(m.editInput.View())
		b.WriteString("\n")
	} else {
		if m.IsMovingCategory() {
			b.WriteString(MutedText.Render(fmt.Sprintf("Select a new group for category '%s'", m.movingCategory.CategoryName)))
			b.WriteString("\n\n")
		}
		
		if len(m.categories) == 0 {
			b.WriteString(MutedText.Render("No category defined yet."))
		} else {
			for i, item := range m.categories {
				style := NormalListItem
				prefix := " "
				if i == m.cursor {
					style = FocusedListItem
					prefix = "> "
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
				line := fmt.Sprintf("%s %s - %s", prefix, item.CategoryName, groupName)
				b.WriteString(style.Render(line))
				b.WriteString("\n")
			}
		}
		b.WriteString("\n\n")
		keyHints := "(j/k: Nav, a/n: Add, e: Edit, d: Delete, m: Move, Esc/q: Back)"
		if m.IsMovingCategory() {
			keyHints = "(Select a group to move the category, Esc/q: Cancel)"
		}
		b.WriteString(MutedText.Render(keyHints))
	}

	viewStr := AppStyle.Width(m.Width).Height(m.Height - 3).Render(b.String())
	return viewStr
}

func (m CategoryModel) AddCategory(group data.CategoryGroup) (CategoryModel, tea.Cmd) {
	m.addCategory = true
	m.selectedGroup = group
	m.editInput.Focus()
	return m, textinput.Blink
}

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

func (m CategoryModel) IsMovingCategory() bool {
	return m.movingCategory.CatID != ""
}

func (m CategoryModel) ResetMoveState() CategoryModel {
	m.movingCategory = data.Category{}
	m.moveCategory = false
	return m
}

func (m CategoryModel) focusInput() (tea.Model, tea.Cmd) {
	m.isEditingName = true
	m.editInput.Focus()
	return m, textinput.Blink
}

func (m CategoryModel) UpdateData(updatedData *data.DataRoot) CategoryModel {
	m.Data = updatedData

	var categories []data.Category
	if record, ok := m.Data.MonthlyData[m.MonthKey]; ok {
		categories = record.Categories
	}

	m.categories = categories
	if m.cursor >= len(m.categories) && len(m.categories) > 0 {
		m.cursor = len(m.categories) - 1
	} else {
		m.cursor = 0
	}
	
	// Reset move state when data is updated
	m = m.ResetMoveState()
	
	return m
}
