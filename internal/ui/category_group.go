package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/data"
)

type CategoryGroupModel struct {
	Data
	WindowSize

	cursor int
	groups []data.CategoryGroup

	isEditingName bool            // True if currently editing a group name or adding new one
	editInput     textinput.Model // Text input for the group name
	editingIndex  int             // Index of the group being edited, -1 for new group
}

func NewCategoryGroupModel(data *data.DataRoot) *CategoryGroupModel {
	ti := textinput.New()
	ti.Placeholder = "Group Name"
	ti.CharLimit = 30
	ti.Width = 30
	return &CategoryGroupModel{
		Data: Data{
			Root: data,
		},
		groups:       data.CategoryGroups,
		editInput:    ti,
		editingIndex: -1,
	}
}

func (m CategoryGroupModel) Init() tea.Cmd {
	return nil
}

func (m CategoryGroupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m CategoryGroupModel) View() string {
	return "Category group view"
}
