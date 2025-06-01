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

type AppViews struct {
	monthlyModel       *ui.MonthlyModel
	categoryGroupModel *ui.CategoryGroupModel
}

type App struct {
	ui.Data
	ui.MonthYear
	ui.WindowSize
	AppViews
	activeView currentView
}

func New(initialData *data.DataRoot, dataFilePath string) App {
	now := time.Now()
	currentM := now.Month()
	currentY := now.Year()

	return App{
		Data: ui.Data{
			Root:     initialData,
			FilePath: dataFilePath,
		},
		MonthYear: ui.MonthYear{
			CurrentMonth: currentM,
			CurrentYear:  currentY,
		},
		AppViews: AppViews{
			monthlyModel:       ui.NewMonthlyModel(initialData, currentM, currentY),
			categoryGroupModel: ui.NewCategoryGroupModel(initialData),
		},
	}
}

func (m App) Init() tea.Cmd {
	switch m.activeView {

	case viewMonthlyOverview:
		if m.monthlyModel != nil {
			return m.monthlyModel.Init()
		}

	case viewCategoryGroup:
		if m.categoryGroupModel != nil {
			return m.categoryGroupModel.Init()
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
			if m.monthlyModel != nil {
				updatedMonthlyModel, monthlyCmd := m.monthlyModel.Update(msg)
				if mo, ok := updatedMonthlyModel.(ui.MonthlyModel); ok {
					m.monthlyModel = &mo
				}
				return m, monthlyCmd
			}

		case viewCategoryGroup:
			updatedModel, cmd := m.handleCategoryGroupViewKeys(msg.String())
			if cmd != nil || updatedModel != m {
				return updatedModel, cmd
			}
			if m.categoryGroupModel != nil {
				updatedCategoryGroupModel, categoryCmd := m.categoryGroupModel.Update(msg)
				if cgMo, ok := updatedCategoryGroupModel.(ui.CategoryGroupModel); ok {
					m.categoryGroupModel = &cgMo
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
		return m.handleMonthlyViewMsg()

	case ui.GroupAddMsg:
		return m.handleGroupAddMsg(msg)

	case ui.GroupDeleteMsg:
		return m.handleGroupDeleteMsg()

	case ui.GroupUpdateMsg:
		return m.handleGroupUpdateMsg()

	case ui.GroupManageCategoriesMsg:
		return m.handleGroupManageCategoriesMsg()
	}

	return m, nil
}

func (m App) View() string {

	var viewContent string

	switch m.activeView {

	case viewMonthlyOverview:
		if m.monthlyModel != nil {
			viewContent = m.monthlyModel.View()
		} else {
			viewContent = "Monthly overview loading..."
		}

	case viewCategoryGroup:
		if m.categoryGroupModel != nil {
			viewContent = m.categoryGroupModel.View()
		} else {
			viewContent = "Category groups loading..."
		}
	default:
		viewContent = "Error: View not found or not initialized"
	}

	return viewContent
}

func (m App) handleModelsWindowResize(msg tea.Msg) (tea.Model, []tea.Cmd) {
	var cmds []tea.Cmd

	if m.monthlyModel != nil {
		updatedMonthlyModel, moCmd := m.monthlyModel.Update(msg)
		if mo, ok := updatedMonthlyModel.(ui.MonthlyModel); ok {
			m.monthlyModel = &mo
		}
		cmds = append(cmds, moCmd)
	}

	if m.categoryGroupModel != nil {
		updatedCategoryGroupModel, cgCmd := m.categoryGroupModel.Update(msg)
		if cgMo, ok := updatedCategoryGroupModel.(ui.CategoryGroupModel); ok {
			m.categoryGroupModel = &cgMo
		}
		cmds = append(cmds, cgCmd)
	}
	return m, cmds
}
