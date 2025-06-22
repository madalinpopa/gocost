package app

import (
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/domain"
	"github.com/madalinpopa/gocost/internal/service"
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

// App represents the main application. It now holds services instead of raw data.
type App struct {
	ui.MonthYear
	ui.WindowSize
	ui.AppViews

	filePath      string
	activeView    currentView
	statusMessage string
	isInitialized bool // Flag to track initial model creation

	// Services for business logic
	categorySvc *service.CategoryService
	groupSvc    *service.GroupService
	incomeSvc   *service.IncomeService
}

// New creates a new instance of the application.
func New(
	categoryService *service.CategoryService,
	groupService *service.GroupService,
	incomeService *service.IncomeService,
	dataFilePath string,
) App {
	now := time.Now()
	currentM := now.Month()
	currentY := now.Year()

	app := App{
		filePath: dataFilePath,
		MonthYear: ui.MonthYear{
			CurrentMonth: currentM,
			CurrentYear:  currentY,
		},
		categorySvc: categoryService,
		groupSvc:    groupService,
		incomeSvc:   incomeService,
	}

	// Initial data load and model creation
	app = app.refreshDataForModels()

	return app
}

// refreshDataForModels fetches the latest data from services and updates the UI models.
func (m App) refreshDataForModels() App {
	monthKey := ui.GetMonthKey(m.CurrentMonth, m.CurrentYear)

	groups, err := m.groupSvc.GetAllGroups()
	if err != nil {
		log.Printf("Error fetching groups: %v", err)
	}

	categories, err := m.categorySvc.GetCategoriesForMonth(monthKey)
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
	}

	incomes, err := m.incomeSvc.GetIncomesForMonth(monthKey)
	if err != nil {
		log.Printf("Error fetching incomes: %v", err)
	}

	appData := ui.AppData{
		Categories:     categories,
		CategoryGroups: groups,
		Incomes:        incomes,
	}

	// Use constructors on first run, then update methods
	if !m.isInitialized {
		m.MonthlyModel = ui.NewMonthlyModel(appData, m.CurrentMonth, m.CurrentYear)
		m.CategoryModel = ui.NewCategoryModel(appData, m.CurrentMonth, m.CurrentYear)
		monthYear := ui.MonthYear{CurrentMonth: m.CurrentMonth, CurrentYear: m.CurrentYear}
		m.CategoryGroupModel = ui.NewCategoryGroupModel(groups, m.Width, m.Height, monthYear)
		m.IncomeModel = ui.NewIncomeModel(incomes, m.CurrentMonth, m.CurrentYear)
		m.ExpenseModel = ui.NewExpenseModel(domain.Category{}, "")
		m.isInitialized = true
	} else {
		// First, update the data slices in each model
		m.MonthlyModel = m.MonthlyModel.UpdateData(appData)
		m.CategoryModel = m.CategoryModel.UpdateData(appData)
		m.CategoryGroupModel = m.CategoryGroupModel.UpdateData(groups)
		m.IncomeModel = m.IncomeModel.UpdateData(incomes)

		// Then, update the month/year for each model that tracks it
		m.MonthlyModel = m.MonthlyModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
		m.CategoryModel = m.CategoryModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
		m.CategoryGroupModel = m.CategoryGroupModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
		m.IncomeModel = m.IncomeModel.SetMonthYear(m.CurrentMonth, m.CurrentYear)
	}

	return m
}

// Init initializes the application.
func (m App) Init() tea.Cmd {
	return nil // No initial command needed
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
				return m.refreshDataForModels(), nil
			case "c":
				m.activeView = viewCategory
				return m.refreshDataForModels(), nil
			case "g":
				m.activeView = viewCategoryGroup
				return m.refreshDataForModels(), nil
			case "h":
				m.CurrentYear, m.CurrentMonth = ui.GetPreviousMonth(m.CurrentYear, m.CurrentMonth)
				return m.refreshDataForModels(), nil
			case "l":
				m.CurrentYear, m.CurrentMonth = ui.GetNextMonth(m.CurrentYear, m.CurrentMonth)
				return m.refreshDataForModels(), nil
			case "r":
				now := time.Now()
				m.CurrentMonth = now.Month()
				m.CurrentYear = now.Year()
				return m.refreshDataForModels(), nil
			}
			updatedMonthlyModel, monthlyCmd := m.MonthlyModel.Update(msg)
			if mo, ok := updatedMonthlyModel.(ui.MonthlyModel); ok {
				m.MonthlyModel = mo
			}
			return m, monthlyCmd
		case viewIncome, viewCategoryGroup, viewCategory, viewExpense, viewIncomeForm:
			// Delegate message to the active view
			var updatedModel tea.Model
			var cmd tea.Cmd
			switch m.activeView {
			case viewIncome:
				updatedModel, cmd = m.IncomeModel.Update(msg)
				if model, ok := updatedModel.(ui.IncomeModel); ok {
					m.IncomeModel = model
				}
			case viewIncomeForm:
				updatedModel, cmd = m.IncomeFormModel.Update(msg)
				if model, ok := updatedModel.(ui.IncomeFormModel); ok {
					m.IncomeFormModel = model
				}
			case viewCategoryGroup:
				updatedModel, cmd = m.CategoryGroupModel.Update(msg)
				if model, ok := updatedModel.(ui.CategoryGroupModel); ok {
					m.CategoryGroupModel = model
				}
			case viewCategory:
				updatedModel, cmd = m.CategoryModel.Update(msg)
				if model, ok := updatedModel.(ui.CategoryModel); ok {
					m.CategoryModel = model
				}
			case viewExpense:
				updatedModel, cmd = m.ExpenseModel.Update(msg)
				if model, ok := updatedModel.(ui.ExpenseModel); ok {
					m.ExpenseModel = model
				}
			}
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		updatedModel, cmds := m.handleModelsWindowResize(msg)
		return updatedModel, tea.Batch(cmds...)

	// ... other message handlers
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
	case ui.ManageGroupsMsg:
		return m.handleManageGroupsMsg()
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
	case ui.CategoryViewMsg:
		return m.handleCategoryViewMsg()
	case ui.CategoryViewWithMonthMsg:
		return m.handleCategoryViewWithMonthMsg(msg)
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
	statusLine := "\n"
	if m.HasStatus() {
		statusLine += m.GetStatusMessage()
	}
	viewContent += statusLine
	return viewContent
}
