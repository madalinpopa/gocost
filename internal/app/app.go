package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/data"
	"github.com/madalinpopa/gocost/internal/ui"
)

// currentView represents the current view of the application.
type currentView int

const (
	viewMonthlyOverview currentView = iota
	viewIncome
	viewIncomeForm
	viewCategoryGroup
	viewCategory
	viewExpense
)

// App represents the main application.
type App struct {
	ui.AppData
	ui.MonthYear
	ui.WindowSize
	ui.AppViews

	activeView    currentView
	statusMessage string
}

// New creates a new instance of the application.
func New(initialData *data.DataRoot, dataFilePath string) App {
	now := time.Now()
	currentM := now.Month()
	currentY := now.Year()

	return App{
		AppData: ui.AppData{
			Data:     initialData,
			FilePath: dataFilePath,
		},
		MonthYear: ui.MonthYear{
			CurrentMonth: currentM,
			CurrentYear:  currentY,
		},
		AppViews: ui.AppViews{
			MonthlyModel:       ui.NewMonthlyModel(initialData, currentM, currentY),
			IncomeModel:        ui.NewIncomeModel(initialData, currentM, currentY),
			CategoryModel:      ui.NewCategoryModel(initialData, currentM, currentY),
			CategoryGroupModel: ui.NewCategoryGroupModel(initialData),
			ExpenseModel:       ui.NewExpenseModel(initialData, data.Category{}, ""),
		},
	}
}

// Init initializes the application.
func (m App) Init() tea.Cmd {
	switch m.activeView {

	case viewMonthlyOverview:
		return m.MonthlyModel.Init()

	case viewIncome:
		return m.IncomeModel.Init()

	case viewCategoryGroup:
		return m.CategoryGroupModel.Init()

	case viewCategory:
		return m.CategoryModel.Init()

	case viewExpense:
		return m.ExpenseModel.Init()
	case viewIncomeForm:
		return m.IncomeFormModel.Init()
	}
	return nil
}

// Update updates the application state.
func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch m.activeView {

		case viewMonthlyOverview:
			switch msg.String() {

			case "ctrl+c", "q":
				return m, tea.Quit

			case "i":
				m.activeView = viewIncome
				return m, m.IncomeModel.Init()

			case "c":
				m.CategoryModel = m.CategoryModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
				m.activeView = viewCategory
				return m, m.CategoryModel.Init()

			case "g":
				m.activeView = viewCategoryGroup
				return m, m.CategoryGroupModel.Init()

			case "h":
				m.CurrentYear, m.CurrentMonth = ui.GetPreviousMonth(m.CurrentYear, m.CurrentMonth)
				m.MonthlyModel = m.MonthlyModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
				m.IncomeModel = m.IncomeModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
				m.CategoryModel = m.CategoryModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
				return m, nil

			case "l":
				m.CurrentYear, m.CurrentMonth = ui.GetNextMonth(m.CurrentYear, m.CurrentMonth)
				m.MonthlyModel = m.MonthlyModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
				m.IncomeModel = m.IncomeModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
				m.CategoryModel = m.CategoryModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
				return m, nil

			case "r":
				now := time.Now()
				m.CurrentMonth = now.Month()
				m.CurrentYear = now.Year()
				m.MonthlyModel = m.MonthlyModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
				m.IncomeModel = m.IncomeModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
				m.CategoryModel = m.CategoryModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
				return m, nil
			}
			updatedMonthlyModel, monthlyCmd := m.MonthlyModel.Update(msg)
			if mo, ok := updatedMonthlyModel.(ui.MonthlyModel); ok {
				m.MonthlyModel = mo
			}
			return m, monthlyCmd

		case viewIncome:
			updatedIncomeModel, incomeCmd := m.IncomeModel.Update(msg)
			if inMo, ok := updatedIncomeModel.(ui.IncomeModel); ok {
				m.IncomeModel = inMo
			}
			return m, incomeCmd

		case viewIncomeForm:
			updatedIncomeModelForm, incomeCmd := m.IncomeFormModel.Update(msg)
			if inFoMo, ok := updatedIncomeModelForm.(ui.IncomeFormModel); ok {
				m.IncomeFormModel = inFoMo
			}
			return m, incomeCmd

		case viewCategoryGroup:
			updatedCategoryGroupModel, categoryCmd := m.CategoryGroupModel.Update(msg)
			if cgMo, ok := updatedCategoryGroupModel.(ui.CategoryGroupModel); ok {
				m.CategoryGroupModel = cgMo
			}
			return m, categoryCmd

		case viewCategory:
			updatedCategoryModel, categoryCmd := m.CategoryModel.Update(msg)
			if cgMo, ok := updatedCategoryModel.(ui.CategoryModel); ok {
				m.CategoryModel = cgMo
			}
			return m, categoryCmd

		case viewExpense:
			updatedExpenseModel, expenseCmd := m.ExpenseModel.Update(msg)
			if expMo, ok := updatedExpenseModel.(ui.ExpenseModel); ok {
				m.ExpenseModel = expMo
			}
			return m, expenseCmd
		}

	case tea.WindowSizeMsg:
		var cmds []tea.Cmd
		var updatedModel tea.Model
		m.Width = msg.Width
		m.Height = msg.Height
		updatedModel, cmds = m.handleModelsWindowResize(msg)
		return updatedModel, tea.Batch(cmds...)

	case ui.MonthlyViewMsg:
		return m.handleMonthlyViewMsg()

	case ui.IncomeViewMsg:
		return m.handleIncomeViewMsg()

	case ui.AddIncomeFormMsg:
		return m.handleAddIncomeFormMsg()

	case ui.SaveIncomeMsg:
		return m.handleSaveIncomeMsg(msg)

	case ui.EditIncomeMsg:
		return m.handleEditIncomeMsg(msg)

	case ui.DeleteIncomeMsg:
		return m.handleDeleteIncomeMsg(msg)

	case ui.GroupAddMsg:
		return m.handleGroupAddMsg(msg)

	case ui.GroupDeleteMsg:
		return m.handleGroupDeleteMsg(msg)

	case ui.GroupUpdateMsg:
		return m.handleGroupUpdateMsg(msg)

	case ui.SelectGroupMsg:
		return m.handleSelectGroupMsg()

	case ui.SelectedGroupMsg:
		return m.handleSelectedGroupMsg(msg)

	case ui.CategoryAddMsg:
		return m.handleCategoryAddMsg(msg)

	case ui.CategoryUpdateMsg:
		return m.handleCategoryUpdateMsg(msg)

	case ui.CategoryDeleteMsg:
		return m.handleCategoryDeleteMsg(msg)

	case ui.FilterCategoriesMsg:
		return m.handleFilterCategoriesMsg(msg)

	case ui.ViewErrorMsg:
		return m.handleViewErrorMsg(msg)

	case ui.PopulateCategoriesMsg:
		return m.handlePopulateCategoriesMsg(msg)

	case ui.ExpenseViewMsg:
		return m.handleExpenseViewMsg(msg)

	case ui.SaveExpenseMsg:
		return m.handleSaveExpenseMsg(msg)

	case ui.EditExpenseMsg:
		return m.handleEditExpenseMsg(msg)

	case ui.DeleteExpenseMsg:
		return m.handleDeleteExpenseMsg(msg)

	case ui.ToggleExpenseStatusMsg:
		return m.handleToggleExpenseStatusMsg(msg)

	case ui.ReturnToMonthlyWithFocusMsg:
		return m.handleReturnToMonthlyWithFocusMsg(msg)

	case StatusClearMsg:
		return m.ClearStatus(), nil
	}

	return m, nil
}

// View returns the current view of the application.
func (m App) View() string {

	var viewContent string

	switch m.activeView {

	case viewMonthlyOverview:
		viewContent = m.MonthlyModel.View()

	case viewIncome:
		viewContent = m.IncomeModel.View()

	case viewIncomeForm:
		viewContent = m.IncomeFormModel.View()

	case viewCategoryGroup:
		viewContent = m.CategoryGroupModel.View()

	case viewCategory:
		viewContent = m.CategoryModel.View()

	case viewExpense:
		viewContent = m.ExpenseModel.View()

	default:
		viewContent = "Error: View not found or not initialized"
	}

	// Add status message at the bottom if present
	statusLine := "\n\n"
	if m.HasStatus() {
		statusLine += m.GetStatusMessage()
	}
	viewContent += statusLine
	return viewContent
}
