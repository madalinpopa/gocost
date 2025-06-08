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
	selectedGroup data.CategoryGroup // Selected group

	isEditingName bool
	editInput     textinput.Model
	editingIndex  int
}

func NewCategoryModel(initialData *data.DataRoot, month time.Month, year int) *CategoryModel {

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

	m := CategoryModel{
		AppData: AppData{
			Data: initialData,
		},
		categories:   categories,
		editInput:    ti,
		editingIndex: -1,
	}

	return &m
}

func (m CategoryModel) Init() tea.Cmd {
	return nil
}

func (m CategoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.isEditingName {

		// Handle actions when editing
		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch msg.String() {

			case "enter":
				categoryName := strings.TrimSpace(m.editInput.Value())
				if categoryName != "" {
					if m.editingIndex == -1 {
						newCategoryId := GenerateID()

						newCategory := data.Category{
							CatID:        newCategoryId,
							CategoryName: categoryName,
						}

						return m.blurInput(), func() tea.Msg {
							return CategoryAddMsg{MonthKey: m.MonthKey, Category: newCategory}
						}
					} else {
						updatedCategory := m.categories[m.editingIndex]
						updatedCategory.CategoryName = categoryName
						return m.blurInput(), func() tea.Msg {
							return CategoryUpdateMsg{Category: updatedCategory}
						}
					}

				}

			case "esc":
				updatedModel := m.blurInput()
				m.editInput.SetValue("")
				return updatedModel, nil
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
		case "d":

		}
	}

	return m, nil
}

func (m CategoryModel) View() string {
	var b strings.Builder

	b.WriteString(HeaderText.Render("Manage Expense Categories"))
	b.WriteString("\n\n")

	if m.isEditingName {
		b.WriteString("Enter Category Name (Enter to save, Esc to cancel):\n")
		b.WriteString(m.editInput.View())
		b.WriteString("\n")
	} else {
		if len(m.categories) == 0 {
			b.WriteString(MutedText.Render("No category defined yet."))
		} else {
			for i, item := range m.categories {
				_, _ = i, item
				style := NormalListItem
				prefix := " "
				if i == m.cursor {
					style = FocusedListItem
					prefix = "> "
				}

				var groupName string
				group, ok := m.Data.CategoryGroups[item.GroupID]
				if ok {
					groupName = group.GroupName
				}
				line := fmt.Sprintf("%s%s - %s", prefix, item.CategoryName, groupName)
				b.WriteString(style.Render(line))
				b.WriteString("\n")
			}
		}
		b.WriteString("\n\n")
		keyHints := "(j/k: Nav, a/n: Add, e: Edit, d: Delete, Esc/q: Back)"
		b.WriteString(MutedText.Render(keyHints))
	}

	viewStr := AppStyle.Width(m.Width).Height(m.Height - 3).Render(b.String())
	return viewStr
}

func (m CategoryModel) AddCategory(group data.CategoryGroup) CategoryModel {
	m.addCategory = true
	m.selectedGroup = group
	return m
}

func (m CategoryModel) focusInput() (tea.Model, tea.Cmd) {
	m.isEditingName = true
	m.editInput.Focus()
	return m, textinput.Blink
}

func (m CategoryModel) blurInput() tea.Model {
	m.isEditingName = false
	m.editInput.Blur()
	m.editingIndex = -1
	return m
}
