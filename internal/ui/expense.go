package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/madalinpopa/gocost/internal/data"
)

const (
	focusAmount = iota
	focusBudget
	focusNotes
	focusSave
	focusCancel
	focusClear
)

type ExpenseModel struct {
	AppData
	WindowSize
	MonthYear

	amountInput textinput.Model
	budgetInput textinput.Model
	notesInput  textarea.Model

	focusIndex int // 0. amount, 1: budget, 2: notes, 3: Save, 4: Cancel


	expenseCategory data.Category
	existingExpense data.ExpenseRecord
	monthKey        string
	hasExistingExpense bool
}

// NewExpenseModel creates a new ExpenseModel instance for managing expense data.
func NewExpenseModel(initialData *data.DataRoot, category data.Category, monthKey string) ExpenseModel {

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

	var expenseRecord data.ExpenseRecord
	var existing bool
	
	if category.Expense != nil {
		expenseRecord, existing = category.Expense[category.CatID]
	}
	
	if !existing {
		expenseRecord = data.ExpenseRecord{}
	} else {
		ai.SetValue(fmt.Sprintf("%.2f", expenseRecord.Amount))
		bi.SetValue(fmt.Sprintf("%.2f", expenseRecord.Budget))
		ni.SetValue(expenseRecord.Notes)
	}

	m := ExpenseModel{
		AppData: AppData{
			Data: initialData,
		},
		amountInput:        ai,
		budgetInput:        bi,
		notesInput:         ni,
		expenseCategory:    category,
		existingExpense:    expenseRecord,
		monthKey:           monthKey,
		hasExistingExpense: existing,
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

// Init initializes the ExpenseModel.
func (m ExpenseModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the ExpenseModel state.
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
			return m, func() tea.Msg { 
				return ReturnToMonthlyWithFocusMsg{
					Category: m.expenseCategory,
				}
			}

		case "tab", "shift+tab", "up", "down":
			// Focus traversal
			if msg.String() == "shift+tab" || msg.String() == "up" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			maxFocus := focusCancel
			if m.hasExistingExpense {
				maxFocus = focusClear
			}

			if m.focusIndex > maxFocus {
				m.focusIndex = focusAmount
			} else if m.focusIndex < focusAmount {
				m.focusIndex = maxFocus
			}

			// Skip focusClear if no existing expense
			if !m.hasExistingExpense && m.focusIndex == focusClear {
				if msg.String() == "shift+tab" || msg.String() == "up" {
					m.focusIndex = focusCancel
				} else {
					m.focusIndex = focusAmount
				}
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
				// Validate and save expense
				amount, err := ValidAmount(m.amountInput.Value())
				if err != nil {
					return m, func() tea.Msg {
						return ViewErrorMsg{
							Text:  "Please provide a valid amount",
							Model: m,
						}
					}
				}

				budget, err := ValidAmount(m.budgetInput.Value())
				if err != nil {
					return m, func() tea.Msg {
						return ViewErrorMsg{
							Text:  "Please provide a valid budget",
							Model: m,
						}
					}
				}

				status := m.existingExpense.Status
				if status == "" {
					status = "Not Paid" // Default status for new expenses
				}

				expense := data.ExpenseRecord{
					Amount: amount,
					Budget: budget,
					Status: status,
					Notes:  m.notesInput.Value(),
				}

				return m, func() tea.Msg {
					return SaveExpenseMsg{
						MonthKey: m.monthKey,
						Category: m.expenseCategory,
						Expense:  expense,
					}
				}
			} else if m.focusIndex == focusCancel {
				return m, func() tea.Msg { 
					return ReturnToMonthlyWithFocusMsg{
						Category: m.expenseCategory,
					}
				}
			} else if m.focusIndex == focusClear && m.hasExistingExpense {
				return m, func() tea.Msg {
					return DeleteExpenseMsg{
						MonthKey: m.monthKey,
						Category: m.expenseCategory,
					}
				}
			}


		// Handle spacebar for focused inputs
		case " ":
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

// View renders the ExpenseModel as a form for editing expense details.
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

	// Notes
	b.WriteString("Notes: \n")
	b.WriteString(m.notesInput.View())
	b.WriteString("\n\n")

	// Buttons
	saveButton := RenderButton("Save", m.focusIndex == focusSave)
	cancelButton := RenderButton("Cancel", m.focusIndex == focusCancel)

	var buttons string
	if m.hasExistingExpense {
		clearButton := RenderButton("Clear", m.focusIndex == focusClear)
		buttons = lipgloss.JoinHorizontal(lipgloss.Top, saveButton, "  ", cancelButton, "  ", clearButton)
	} else {
		buttons = lipgloss.JoinHorizontal(lipgloss.Top, saveButton, "  ", cancelButton)
	}
	b.WriteString(buttons)
	b.WriteString("\n\n")
	
	helpText := "(Tab/Shift+Tab to navigate, Enter to select/save, Esc to cancel"
	if m.hasExistingExpense {
		helpText += ", Clear to reset"
	}
	helpText += ", Status can be toggled from monthly view with 't')"
	b.WriteString(MutedText.Render(helpText))

	popupContent := AppStyle.Width(m.Width).Align(lipgloss.Center).Render(b.String())
	return FocusedBorder.Render(popupContent)
}
