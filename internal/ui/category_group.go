package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/data"
)

type GroupDeleteMsg struct {
	GroupID string
}

type GroupAddMsg struct {
	Group data.CategoryGroup
}

type GroupUpdateMsg struct {
	Group data.CategoryGroup
}

type GroupManageCategoriesMsg struct {
	Group data.CategoryGroup
}

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

	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.isEditingName {

		// Handle actions when editing
		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch msg.String() {

			case "enter":
				groupName := strings.TrimSpace(m.editInput.Value())
				if groupName != "" {
					newGroupId := GenerateID()
					newGroup := data.CategoryGroup{
						GroupID:    newGroupId,
						GroupName:  groupName,
						Categories: make([]data.Category, 0),
					}
					return m.blurInput(), func() tea.Msg {
						return GroupAddMsg{Group: newGroup}
					}
				}

			case "esc":
				m.blurInput()
				m.editInput.SetValue("")
			}
		}
		m.editInput, cmd = m.editInput.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:

		// Handle navigation and select actions
		switch msg.String() {

		case "q", "esc":
			return m, func() tea.Msg { return SwitchToMonthlyOverviewMsg{} }

		case "j", "down":
			if len(m.groups) > 0 {
				m.cursor = (m.cursor + 1) % len(m.groups)
			}
			return m, nil

		case "k", "up":
			if len(m.groups) > 0 {
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(m.groups) - 1
				}
			}
			return m, nil

		case "a", "n": // Add new category group name
			m.editingIndex = -1
			m.editInput.SetValue("")
			m.editInput.Placeholder = "New Group Name"
			return m.focusInput()

		case "e": // Edit selected category group name
			if len(m.groups) > 0 {
				if m.cursor >= 0 && m.cursor < len(m.groups) {
					m.editingIndex = m.cursor
					m.editInput.SetValue(m.groups[m.cursor].GroupName)
					m.editInput.Placeholder = "Edit Group Name"
					return m.focusInput()
				}

			}

		case "d": // Delete selected category group
			if len(m.groups) > 0 {
				if m.cursor >= 0 && m.cursor < len(m.groups) {
					groupIDToDelete := m.groups[m.cursor].GroupID
					return m, func() tea.Msg {
						return GroupDeleteMsg{
							GroupID: groupIDToDelete,
						}
					}
				}
			}

		case "enter": // Manage categories for selected group
			if len(m.groups) > 0 {
				if m.cursor >= 0 && m.cursor < len(m.groups) {
					selectedGroup := m.groups[m.cursor]
					return m, func() tea.Msg {
						return GroupManageCategoriesMsg{
							Group: selectedGroup,
						}
					}
				}
			}

		}
		return m, nil
	}

	return m, nil
}

func (m CategoryGroupModel) View() string {
	return "Category group view"
}

func (m CategoryGroupModel) focusInput() (tea.Model, tea.Cmd) {
	m.isEditingName = true
	m.editInput.Focus()
	return m, textinput.Blink
}

func (m CategoryGroupModel) blurInput() tea.Model {
	m.isEditingName = false
	m.editInput.Blur()
	m.editingIndex = -1
	return m
}
