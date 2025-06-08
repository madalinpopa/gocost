package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/madalinpopa/gocost/internal/data"
)

const (
	focusAmount = iota
	focusBudget
	focusStatus
	focusNotes
	focusSave
	focusCancel
)

type ExpenseModel struct {
	AppData
	WindowSize
	MonthYear

	MonthKey string

	amountInput textinput.Model
	budgetInput textinput.Model
	notesInput  textarea.Model
	status      string

	focusIndex int // 0. amount, 1: budget, 3: status, 4: notes, 5: Save, 6: Cancel

	statusOptions []string
	statusCursor  int

	expenseCategory data.Category
	existingExpense data.ExpenseRecord
}

func NewExpenseModel(category data.Category, month time.Month, year int) ExpenseModel {

	monthKey := GetMonthKey(month, year)

	ai := textinput.New()
	ai.Placeholder = "0.00"
	ai.Focus()
	ai.CharLimit = 10
	ai.Width = 20

	bi := textinput.New()
	bi.Placeholder = "0.00"
	bi.CharLimit = 10
	bi.Width = 20

	ni := textarea.New()
	ni.Placeholder = "Optional notes.."
	ni.SetHeight(3)
	ni.SetWidth(30)

	currentStatus := "Not Paid"
	expenseRecord, existing := category.Expense[category.CatID]
	if !existing {
		expenseRecord = data.ExpenseRecord{
			Status: currentStatus,
		}
	} else {
		ai.SetValue(fmt.Sprintf("%.2f", expenseRecord.Amount))
		bi.SetValue(fmt.Sprintf("%.2f", expenseRecord.Budget))
		ni.SetValue(expenseRecord.Notes)
		currentStatus = expenseRecord.Status
	}

	statusOpts := []string{"Not Paid", "Paid"}
	statusIdx := 0
	if currentStatus == "Paid" {
		statusIdx = 1
	}

	if currentStatus == "Paid" {
		statusIdx = 1
	}
	m := ExpenseModel{
		MonthKey:        monthKey,
		expenseCategory: category,
		existingExpense: expenseRecord,
		amountInput:     ai,
		budgetInput:     bi,
		notesInput:      ni,
		statusOptions:   statusOpts,
		statusCursor:    statusIdx,
		WindowSize: WindowSize{
			Width:  50,
			Height: 15,
		},
	}
	m.amountInput.Width = m.Width - 10
	m.budgetInput.Width = m.Width - 10
	m.notesInput.SetWidth(m.Width - 6)

	return m
}

func (m ExpenseModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ExpenseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:

		switch msg.String() {

		case "esc":
			return m, func() tea.Msg { return MonthlyViewMsg{} }

		case "tab", "shift+tab", "up", "down":
			// Focus traversal
			if msg.String() == "shift+tab" || msg.String() == "up" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > focusCancel {
				m.focusIndex = focusAmount
			} else if m.focusIndex < focusAmount {
				m.focusIndex = focusCancel
			}

			// Update focus on inputs
			m.amountInput.Blur()
			m.budgetInput.Blur()
			m.notesInput.Blur()

			switch m.focusIndex {
			case focusAmount:
				m.amountInput.Focus()
				cmds = append(cmds, textinput.Blink)
			case focusBudget:
				m.budgetInput.Focus()
				cmds = append(cmds, textinput.Blink)
			case focusNotes:
				m.notesInput.Focus()
				cmds = append(cmds, textarea.Blink)
			}

		case "enter":
			if m.focusIndex == focusSave {
				return m, nil
			} else if m.focusIndex == focusCancel {
				return m, func() tea.Msg { return MonthlyViewMsg{} }
			}

			// If on status, toggle it
			if m.focusIndex == focusStatus {
				m.statusCursor = (m.statusCursor + 1) % len(m.statusOptions)
				m.status = m.statusOptions[m.statusCursor]
			}

		// Handle spacebar for status toggle when status is focused
		case " ":
			if m.focusIndex == focusStatus {
				m.statusCursor = (m.statusCursor + 1) % len(m.statusOptions)
				m.status = m.statusOptions[m.statusCursor]
			} else {
				// If not status, pass space to focused input
				if m.amountInput.Focused() {
					m.amountInput, cmd = m.amountInput.Update(msg)
					cmds = append(cmds, cmd)
				} else if m.budgetInput.Focused() {
					m.budgetInput, cmd = m.budgetInput.Update(msg)
					cmds = append(cmds, cmd)
				} else if m.notesInput.Focused() {
					m.notesInput, cmd = m.notesInput.Update(msg)
					cmds = append(cmds, cmd)
				}
			}

		default:
			if m.amountInput.Focused() {
				m.amountInput, cmd = m.amountInput.Update(msg)
				cmds = append(cmds, cmd)
			} else if m.budgetInput.Focused() {
				m.budgetInput, cmd = m.budgetInput.Update(msg)
				cmds = append(cmds, cmd)
			} else if m.notesInput.Focused() {
				m.notesInput, cmd = m.notesInput.Update(msg)
				cmds = append(cmds, cmd)
			}
		}

	}
	return m, tea.Batch(cmds...)
}

func (m ExpenseModel) View() string {
	var b strings.Builder

	title := fmt.Sprintf("Expense: %s", m.expenseCategory.CategoryName)
	b.WriteString(HeaderText.Render(title))
	b.WriteString("\n\n")

	// Amount
	b.WriteString("Amount: \n")
	b.WriteString(m.amountInput.View())
	b.WriteString("\n\n")

	// Budget
	b.WriteString("Budget: \n")
	b.WriteString(m.budgetInput.View())
	b.WriteString("\n\n")

	statusLine := "Status: "
	for i, opt := range m.statusOptions {
		cursor := " "
		style := NormalListItem
		if i == m.statusCursor {
			style = style.Bold(true)
		}
		if m.focusIndex == focusStatus && i == m.statusCursor {
			style = FocusedListItem
			cursor = ">"
		}
		statusLine += style.Render(fmt.Sprintf("%s[%s]", cursor, opt)) + "  "
	}
	b.WriteString(statusLine)
	b.WriteString("\n\n")

	// Notes
	b.WriteString("Notes: \n")
	b.WriteString(m.notesInput.View())
	b.WriteString("\n\n")

	// Buttons
	saveButton := "[ Save ]"
	cancelButton := "[ Cancel ]"

	if m.focusIndex == focusSave {
		saveButton = FocusedListItem.Render(saveButton)
	}
	if m.focusIndex == focusCancel {
		cancelButton = FocusedListItem.Render(cancelButton)
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Top, saveButton, "  ", cancelButton)
	b.WriteString(buttons)
	b.WriteString("\n\n")
	b.WriteString(MutedText.Render("(Tab/Shift+Tab to navigate, Enter to select/save, Esc to cancel)"))

	popupContent := AppStyle.Width(m.Width).Align(lipgloss.Center).Render(b.String())
	return FocusedBorder.Render(popupContent)
}
