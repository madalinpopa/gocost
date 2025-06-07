package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/madalinpopa/gocost/internal/data"
	"github.com/shopspring/decimal"
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

	incomeId         string
	descriptionInput textinput.Model
	amountInput      textinput.Model

	focusIndex int
}

func NewIncomeFormModel(currentMonth time.Month, year int, income *data.IncomeRecord) IncomeFormModel {

	monthKey := GetMonthKey(currentMonth, year)

	descInput := textinput.New()
	descInput.Placeholder = "e.g., Salary, Freelance Project"
	descInput.Focus()
	descInput.CharLimit = 50
	descInput.Width = 30

	amountInput := textinput.New()
	amountInput.Placeholder = "0.00"
	amountInput.CharLimit = 10
	amountInput.Width = 20

	newEntry := true
	originalEntryId := ""

	var incomeAmount decimal.Decimal
	if income != nil {
		newEntry = false
		incomeAmount = decimal.NewFromFloat(income.Amount)
		originalEntryId = income.IncomeID
		descInput.SetValue(income.Description)
		amountInput.SetValue(incomeAmount.String())

	}

	m := IncomeFormModel{
		incomeId:         originalEntryId,
		NewEntry:         newEntry,
		MonthKey:         monthKey,
		descriptionInput: descInput,
		amountInput:      amountInput,
		WindowSize: WindowSize{
			Width:  50,
			Height: 10,
		},
	}

	descInput.Width = m.Width - 10
	amountInput.Width = m.Width - 10

	return m
}

func (m IncomeFormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m IncomeFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return IncomeViewMsg{} }
		}
	}

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
