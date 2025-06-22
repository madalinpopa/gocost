package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/domain"
)

// AppData now holds slices of data for the UI, rather than the whole data root.
// This decouples the UI from the persistence layer.
type AppData struct {
	Categories     []domain.Category
	CategoryGroups []domain.CategoryGroup
	Incomes        []domain.IncomeRecord
}

// MonthYear represents the current month and year.
type MonthYear struct {
	CurrentMonth time.Month
	CurrentYear  int
}

// WindowSize represents the size of a window.
type WindowSize struct {
	Width  int
	Height int
}

// AppViews holds all the models used by App
type AppViews struct {
	MonthlyModel       MonthlyModel
	CategoryGroupModel CategoryGroupModel
	CategoryModel      CategoryModel
	IncomeModel        IncomeModel
	IncomeFormModel    IncomeFormModel
	ExpenseModel       ExpenseModel
}

// ViewErrorMsg represents an error message and the associated model to handle the error state.
type ViewErrorMsg struct {
	Text  string
	Model tea.Model
}

// MonthlyViewMsg represents a message used to trigger updates or actions related to the monthly view in the application.
type MonthlyViewMsg struct{}

// GroupDeleteMsg represents a message to delete a specific category group.
type GroupDeleteMsg struct {
	Group domain.CategoryGroup
}

// GroupAddMsg is a message used to represent the addition of a new CategoryGroup.
type GroupAddMsg struct {
	Group domain.CategoryGroup
}

// GroupUpdateMsg represents a message used to indicate an update to a CategoryGroup.
type GroupUpdateMsg struct {
	Group domain.CategoryGroup
}

// ManageGroupsMsg is a message used to switch to the group management view.
type ManageGroupsMsg struct{}

// SelectGroupMsg is a message used to indicate or trigger the selection of a group in the application state.
type SelectGroupMsg struct{}

type SelectedGroupMsg struct {
	Group domain.CategoryGroup
}

// CategoryAddMsg represents a message containing a month key and a category to be added.
type CategoryAddMsg struct {
	MonthKey string
	Category domain.Category
}

// CategoryUpdateMsg is a message used to update a category for a specific month key.
type CategoryUpdateMsg struct {
	MonthKey string
	Category domain.Category
}

// CategoryDeleteMsg represents a message to delete a specific category for a given month.
type CategoryDeleteMsg struct {
	MonthKey string
	Category domain.Category
}

// CategoryDeleteMove represents the structure for moving a category during deletion within a specific month.
type CategoryDeleteMove struct {
	MonthKey string
	Category domain.Category
}

// FilterCategoriesMsg represents a message containing a filter text used to filter categories.
type FilterCategoriesMsg struct {
	FilterText string
}

// IncomeViewMsg is a message used to signal a view transition to the income view.
type IncomeViewMsg struct{}

// AddIncomeFormMsg represents a message to trigger displaying the Add Income form for a specific month.
type AddIncomeFormMsg struct {
	MonthKey string
}

// SaveIncomeMsg represents a message used to save an income record for a specified month.
type SaveIncomeMsg struct {
	MonthKey string
	Income   domain.IncomeRecord
}

// EditIncomeMsg represents a message for editing an income record for a specific month identified by MonthKey.
type EditIncomeMsg struct {
	MonthKey string
	Income   domain.IncomeRecord
}

// DeleteIncomeMsg represents a message for deleting an income record for a specific month key.
type DeleteIncomeMsg struct {
	MonthKey string
	Income   domain.IncomeRecord
}

// PopulateCategoriesMsg represents a message containing keys for the current and previous month's categories.
type PopulateCategoriesMsg struct {
	CurrentMonthKey  string
	PreviousMonthKey string
}

// ExpenseViewMsg represents a message used to display expenses for a specific month and category.
type ExpenseViewMsg struct {
	MonthKey string
	Category domain.Category
}

// SaveExpenseMsg represents a message for saving an expense record in a specific category and month.
type SaveExpenseMsg struct {
	MonthKey string
	Category domain.Category
	Expense  domain.ExpenseRecord
}

// EditExpenseMsg represents a message for editing an expense entry within a specific month and category context.
type EditExpenseMsg struct {
	MonthKey string
	Category domain.Category
	Expense  domain.ExpenseRecord
}

// DeleteExpenseMsg represents a message to delete an expense from a specified category for a specific month.
type DeleteExpenseMsg struct {
	MonthKey string
	Category domain.Category
}

// ToggleExpenseStatusMsg is used to toggle the status of an expense in a specific category and month.
type ToggleExpenseStatusMsg struct {
	MonthKey string
	Category domain.Category
}

// ReturnToMonthlyWithFocusMsg is sent when returning to monthly view with specific category focus
type ReturnToMonthlyWithFocusMsg struct {
	Category domain.Category
}

// CategoryViewMsg is sent when returning to the category view
type CategoryViewMsg struct{}

// CategoryViewWithMonthMsg is sent when switching to category view with specific month context
type CategoryViewWithMonthMsg struct {
	MonthYear
}
