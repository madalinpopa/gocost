package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

	existingExpense data.ExpenseRecord
}

func NewExpenseModel(initialData *data.DataRoot, month time.Month, year int) ExpenseModel {

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

	// currentStatus := "Not Paid"

	return ExpenseModel{
		AppData: AppData{
			Data: initialData,
		},
		MonthKey: monthKey,

		amountInput: ai,
		budgetInput: bi,
		notesInput:  ni,
	}
}

func (m ExpenseModel) Init() tea.Cmd {
	return nil
}

func (m ExpenseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "esc":
			return m, func() tea.Msg { return MonthlyViewMsg{} }
		}

	}
	return m, nil
}

func (m ExpenseModel) View() string {
	var b strings.Builder

	// title := "Add Expense"

	return b.String()
}
