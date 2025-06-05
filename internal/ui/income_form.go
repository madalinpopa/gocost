package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/madalinpopa/gocost/internal/data"
)

const (
	editFocusDescription = iota
	editFocusAmount
	editFocusSave
	editFocusCancel
	numEditFormFields // For cycling focus
)

type IncomeFormModel struct {
	WindowSize
	NewEntry     bool
	MonthKey     string
	IncomeRecord data.IncomeRecord

	descriptionInput textinput.Model
	amountInput      textinput.Model

	focusIndex int
}

func (m IncomeFormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m IncomeFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m IncomeFormModel) View() string {
	var b strings.Builder
	title := "Add New Income"
	if !m.NewEntry {
		title = "Edit Income"
	}
	b.WriteString(HeaderText.Render(fmt.Sprintf("%s - %s", title, m.MonthKey)))
	b.WriteString("\n\n")

	b.WriteString("Description:\n")
	b.WriteString(m.descriptionInput.View())
	b.WriteString("\n\n")

	b.WriteString("Amount:\n")
	b.WriteString(m.amountInput.View())
	b.WriteString("\n\n")

	saveButton := "[ Save ]"
	cancelButton := "[ Cancel ]"
	if m.focusIndex == editFocusSave {
		saveButton = FocusedListItem.Render(saveButton)
	}
	if m.focusIndex == editFocusCancel {
		cancelButton = FocusedListItem.Render(cancelButton)
	}
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, saveButton, "  ", cancelButton)
	b.WriteString(buttons)
	b.WriteString("\n\n")
	b.WriteString(MutedText.Render("(Tab/Shift+Tab to navigate, Enter to save, Esc to cancel)"))

	popupContent := AppStyle.Width(m.Width).Align(lipgloss.Center).Render(b.String())
	return FocusedBorder.Render(popupContent)
}
