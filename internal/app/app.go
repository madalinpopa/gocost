package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/data"
	"github.com/madalinpopa/gocost/internal/ui"
)

type currentView int

const (
	viewMonthlyOverview currentView = iota
	viewCategoryGroup
)

type App struct {
	ui.AppData
	ui.MonthYear
	ui.WindowSize
	ui.AppViews

	activeView currentView
}

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
			CategoryGroupModel: ui.NewCategoryGroupModel(initialData),
		},
	}
}

func (m App) Init() tea.Cmd {
	switch m.activeView {

	case viewMonthlyOverview:
		if m.MonthlyModel != nil {
			return m.MonthlyModel.Init()
		}

	case viewCategoryGroup:
		if m.CategoryGroupModel != nil {
			return m.CategoryGroupModel.Init()
		}

	}
	return nil
}

func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch m.activeView {

		case viewMonthlyOverview:
			updatedModel, cmd := m.handleMonthlyViewKeys(msg.String())
			if cmd != nil || updatedModel != m {
				return updatedModel, cmd
			}
			if m.MonthlyModel != nil {
				updatedMonthlyModel, monthlyCmd := m.MonthlyModel.Update(msg)
				if mo, ok := updatedMonthlyModel.(ui.MonthlyModel); ok {
					m.MonthlyModel = &mo
				}
				return m, monthlyCmd
			}

		case viewCategoryGroup:
			updatedModel, cmd := m.handleCategoryGroupViewKeys(msg.String())
			if cmd != nil || updatedModel != m {
				return updatedModel, cmd
			}
			if m.CategoryGroupModel != nil {
				updatedCategoryGroupModel, categoryCmd := m.CategoryGroupModel.Update(msg)
				if cgMo, ok := updatedCategoryGroupModel.(ui.CategoryGroupModel); ok {
					m.CategoryGroupModel = &cgMo
				}
				return m, categoryCmd
			}
		}

	case tea.WindowSizeMsg:
		var cmds []tea.Cmd
		var updatedModel tea.Model
		m.Width = msg.Width
		m.Height = msg.Height
		updatedModel, cmds = m.handleModelsWindowResize(msg)
		return updatedModel, tea.Batch(cmds...)

	case ui.MonthlyViewMsg:
		return m.handleMonthlyViewMsg(msg)

	case ui.GroupAddMsg:
		return m.handleGroupAddMsg(msg)

	case ui.GroupDeleteMsg:
		return m.handleGroupDeleteMsg(msg)

	case ui.GroupUpdateMsg:
		return m.handleGroupUpdateMsg(msg)

	case ui.GroupManageCategoriesMsg:
		return m.handleGroupManageCategoriesMsg()
	}

	return m, nil
}

func (m App) View() string {

	var viewContent string

	switch m.activeView {

	case viewMonthlyOverview:
		if m.MonthlyModel != nil {
			viewContent = m.MonthlyModel.View()
		} else {
			viewContent = "Monthly overview loading..."
		}

	case viewCategoryGroup:
		if m.CategoryGroupModel != nil {
			viewContent = m.CategoryGroupModel.View()
		} else {
			viewContent = "Category groups loading..."
		}
	default:
		viewContent = "Error: View not found or not initialized"
	}

	return viewContent
}
