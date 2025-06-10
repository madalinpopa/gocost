package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/data"
)

type CategoryGroupModel struct {
	AppData
	WindowSize

	cursor int
	groups []data.CategoryGroup

	selectGroup bool

	isEditingName bool            // True if currently editing a group name or adding new one
	editInput     textinput.Model // Text input for the group name
	editingIndex  int             // Index of the group being edited, -1 for new group
}

// NewCategoryGroupModel creates a new CategoryGroupModel instance.
func NewCategoryGroupModel(initialData *data.DataRoot) CategoryGroupModel {
	ti := textinput.New()
	ti.Placeholder = "Group Name"
	ti.CharLimit = 30
	ti.Width = 30

	var groups []data.CategoryGroup
	for _, value := range initialData.CategoryGroups {
		groups = append(groups, value)
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Order < groups[j].Order
	})

	return CategoryGroupModel{
		AppData: AppData{
			Data: initialData,
		},
		groups:       groups,
		editInput:    ti,
		editingIndex: -1,
	}
}

// Init initializes the CategoryGroupModel.
func (m CategoryGroupModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the CategoryGroupModel state.
func (m CategoryGroupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd



	if m.isEditingName {

		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch msg.String() {

			case "enter":
				groupName := strings.TrimSpace(m.editInput.Value())
				if groupName != "" {
					if m.editingIndex == -1 {
						newGroupId := GenerateID()

						maxOrder := 0
						for _, group := range m.groups {
							if group.Order > maxOrder {
								maxOrder = group.Order
							}
						}

						newGroup := data.CategoryGroup{
							GroupID:   newGroupId,
							GroupName: groupName,
							Order:     maxOrder + 1,
						}

						return m.blurInput(), func() tea.Msg {
							return GroupAddMsg{Group: newGroup}
						}
					} else {
						updatedGroup := m.groups[m.editingIndex]
						updatedGroup.GroupName = groupName
						return m.blurInput(), func() tea.Msg {
							return GroupUpdateMsg{Group: updatedGroup}
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

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:

		switch msg.String() {

		case "q", "esc":
			m = m.ResetSelection()
			return m, func() tea.Msg { return MonthlyViewMsg{} }

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

		case "enter":
			// Handle selection when in select group mode
			if m.selectGroup {
				if len(m.groups) > 0 {
					if m.cursor >= 0 && m.cursor < len(m.groups) {
						selectedGroup := m.groups[m.cursor]
						return m, func() tea.Msg { return SelectedGroupMsg{Group: selectedGroup} }
					}
				}
			}

		case "a", "n": // Add new category group name (only when not selecting)
			if !m.selectGroup {
				m.editingIndex = -1
				m.editInput.SetValue("")
				m.editInput.Placeholder = "New Group Name"
				return m.focusInput()
			}

		case "e": // Edit selected category group name (only when not selecting)
			if !m.selectGroup {
				if len(m.groups) > 0 {
					if m.cursor >= 0 && m.cursor < len(m.groups) {
						m.editingIndex = m.cursor
						m.editInput.SetValue(m.groups[m.cursor].GroupName)
						m.editInput.Placeholder = "Edit Group Name"
						return m.focusInput()
					}
				}
			}

		case "d": // Delete selected category group (only when not selecting)
			if !m.selectGroup {
				if len(m.groups) > 0 {
					if m.cursor >= 0 && m.cursor < len(m.groups) {
						groupToDelete := m.groups[m.cursor]
						return m, func() tea.Msg {
							return GroupDeleteMsg{
								Group: groupToDelete,
							}
						}
					}
				}
			}
		}
		return m, nil
	}

	return m, nil
}

// View renders the CategoryGroupModel.
func (m CategoryGroupModel) View() string {

	var b strings.Builder

	title := "Manage Category Groups"
	if m.selectGroup {
		title = "Select group"
	}
	b.WriteString(HeaderText.Render(title))
	b.WriteString("\n\n")

	if m.isEditingName {
		b.WriteString("Enter Category Group Name (Enter to save, Esc to cancel):\n")
		b.WriteString(m.editInput.View())
		b.WriteString("\n")
	} else {
		if len(m.groups) == 0 {
			b.WriteString(MutedText.Render("No category groups defined yet."))
		} else {
			for i, item := range m.groups {
				style := NormalListItem
				prefix := "  "
				if i == m.cursor {
					style = FocusedListItem
					prefix = "> "
				}
				groupId := MutedText.Render(item.GroupID)
				line := fmt.Sprintf("%s %d. %s (ID: %s)", prefix, item.Order, item.GroupName, groupId)
				b.WriteString(style.Render(line))
				b.WriteString("\n")
			}
		}
		b.WriteString("\n\n")
		keyHints := "(j/k: Nav, a/n: Add, e: Edit, d: Delete, Esc/q: Back)"
		if m.selectGroup {
			keyHints = "(j/k: Nav, Enter: Select, Esc/q: Back)"
		}
		b.WriteString(MutedText.Render(keyHints))
	}

	viewStr := AppStyle.Width(m.Width).Height(m.Height - 3).Render(b.String())
	return viewStr
}

// UpdateData refreshes the model with new data and resets cursor if needed.
func (m CategoryGroupModel) UpdateData(updatedData *data.DataRoot) CategoryGroupModel {
	var groups []data.CategoryGroup
	for _, value := range m.Data.CategoryGroups {
		groups = append(groups, value)
	}

	sort.Slice(groups, func(i int, j int) bool {
		return groups[i].Order < groups[j].Order
	})

	m.Data = updatedData
	m.groups = groups
	if m.cursor >= len(m.groups) && len(m.groups) > 0 {
		m.cursor = len(m.groups) - 1
	} else {
		m.cursor = 0
	}

	m = m.resetEditingState()
	return m
}

// focusInput activates the text input for group name editing.
func (m CategoryGroupModel) focusInput() (tea.Model, tea.Cmd) {
	m.isEditingName = true
	m.editInput.Focus()
	return m, textinput.Blink
}

// blurInput deactivates the text input and resets editing state.
func (m CategoryGroupModel) blurInput() tea.Model {
	m.isEditingName = false
	m.editInput.Blur()
	m.editInput.SetValue("")
	m.editingIndex = -1
	return m
}

// SelectGroup enables group selection mode.
func (m CategoryGroupModel) SelectGroup() CategoryGroupModel {
	m.selectGroup = true
	return m
}

// resetEditingState resets all editing-related state flags and inputs
func (m CategoryGroupModel) resetEditingState() CategoryGroupModel {
	m.selectGroup = false
	m.isEditingName = false
	m.editInput.Blur()
	m.editInput.SetValue("")
	m.editingIndex = -1
	return m
}

// ResetSelection disables group selection mode.
func (m CategoryGroupModel) ResetSelection() CategoryGroupModel {
	return m.resetEditingState()
}
