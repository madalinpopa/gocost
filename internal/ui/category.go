package ui

import (
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
				// style := NormalListItem
				// prefix := " "
				// if i  == m.cursor {
				// 	style = FocusedListItem
				// 	prefix = "> "
				// }

				// var groupName string
				// line := fmt.Sprintf("%s%s - %s", prefix, item.CategoryName, item.GroupID)
			}
		}
	}

	keyHints := "(j/k: Nav, a/n: Add, e: Edit, d: Delete, Esc/q: Back)"
	b.WriteString(MutedText.Render(keyHints))

	viewStr := AppStyle.Width(m.Width).Height(m.Height - 3).Render(b.String())
	return viewStr
}
