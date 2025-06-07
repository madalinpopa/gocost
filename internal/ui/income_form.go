package ui

import (
	"fmt"
	"strings"
	"time"

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

	if income != nil {
		newEntry = false
		originalEntryId = income.IncomeID
		descInput.SetValue(income.Description)
		amountInput.SetValue(fmt.Sprintf("%.2f", income.Amount))
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

	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "esc":
			return m, func() tea.Msg { return IncomeViewMsg{} }

		case "tab", "shift+tab", "up", "down":
			if msg.String() == "shift+tab" || msg.String() == "up" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > editFocusCancel {
				m.focusIndex = editFocusDescription
			} else if m.focusIndex < editFocusDescription {
				m.focusIndex = editFocusCancel
			}

			m.descriptionInput.Blur()
			m.amountInput.Blur()

			switch m.focusIndex {
			case editFocusDescription:
				m.descriptionInput.Focus()
				cmds = append(cmds, textinput.Blink)
			case editFocusAmount:
				m.amountInput.Focus()
				cmds = append(cmds, textinput.Blink)
			}

		case "enter":
			if m.focusIndex == editFocusSave {
				if m.NewEntry {

					// amount cannot be 0 or invalid
					amount, err := ValidAmount(m.amountInput.Value())
					if err != nil {
						return m, func() tea.Msg {
							return ErrorStatusMsg{
								Text:  "Please provide a valid amount",
								Model: m,
							}
						}
					}

					newIncome := data.IncomeRecord{
						IncomeID:    GenerateID(),
						Description: m.descriptionInput.Value(),
						Amount:      amount,
					}

					return m, func() tea.Msg {
						return SaveIncomeMsg{
							MonthKey: m.MonthKey,
							Income:   newIncome,
						}
					}
				} else {

					// amount cannot be 0 or invalid
					amount, err := ValidAmount(m.amountInput.Value())
					if err != nil {
						return m, func() tea.Msg {
							return ErrorStatusMsg{
								Text:  "Please provide a valid amount",
								Model: m,
							}
						}
					}

					m.IncomeRecord.Amount = amount
					m.IncomeRecord.IncomeID = m.incomeId
					m.IncomeRecord.Description = m.descriptionInput.Value()
					return m, func() tea.Msg {
						return SaveIncomeMsg{
							MonthKey: m.MonthKey,
							Income:   m.IncomeRecord,
						}
					}
				}
			} else if m.focusIndex == editFocusCancel {
				return m, func() tea.Msg { return IncomeViewMsg{} }
			}

		default:
			// Handle regular typing
			if m.descriptionInput.Focused() {
				m.descriptionInput, cmd = m.descriptionInput.Update(msg)
				cmds = append(cmds, cmd)
			} else if m.amountInput.Focused() {
				m.amountInput, cmd = m.amountInput.Update(msg)
				cmds = append(cmds, cmd)
			}
		}

	}

	// Manage blink for focused input
	isBlinking := false
	for _, c := range cmds {
		if c != nil && fmt.Sprintf("%p", c) == fmt.Sprintf("%p", textinput.Blink) {
			isBlinking = true
			break
		}
	}

	if (m.descriptionInput.Focused() || m.amountInput.Focused()) && !isBlinking {
		cmds = append(cmds, textinput.Blink)
	}

	return m, tea.Batch(cmds...)
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
