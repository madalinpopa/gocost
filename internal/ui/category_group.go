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

	addCategory   bool
	isEditingName bool            // True if currently editing a group name or adding new one
	editInput     textinput.Model // Text input for the group name
	editingIndex  int             // Index of the group being edited, -1 for new group
}

func NewCategoryGroupModel(initialData *data.DataRoot) *CategoryGroupModel {
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

	return &CategoryGroupModel{
		AppData: AppData{
			Data: initialData,
		},
		groups:       groups,
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

		// Handle navigation and select actions
		switch msg.String() {

		case "q", "esc":
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
					groupToDelete := m.groups[m.cursor]
					return m, func() tea.Msg {
						return GroupDeleteMsg{
							Group: groupToDelete,
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

	var b strings.Builder

	b.WriteString(HeaderText.Render("Manage Category Groups"))
	if m.addCategory {
		b.WriteString(HeaderText.Render("Select group"))
	}
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
				line := fmt.Sprintf("%s %d. %s (ID: %s)", prefix, item.Order, item.GroupName, item.GroupID)
				b.WriteString(style.Render(line))
				b.WriteString("\n")
			}
		}
		b.WriteString("\n\n")
		keyHints := "(j/k: Nav, a/n: Add, e: Edit, d: Delete, Esc/q: Back)"
		if m.addCategory {
			keyHints = "(j/k: Nav, Enter: Select, Esc/q: Back)"
		}
		b.WriteString(MutedText.Render(keyHints))
	}

	viewStr := AppStyle.Width(m.Width).Height(m.Height - 3).Render(b.String())
	return viewStr
}

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
	return m
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
