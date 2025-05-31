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
)

type AppViews struct {
	monthlyModel *ui.MonthlyModel
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
			monthlyModel: ui.NewMonthlyModel(initialData, currentM, currentY),
		},
	}
}

func (m App) Init() tea.Cmd {
	switch m.activeView {

	case viewMonthlyOverview:
		if m.monthlyModel != nil {
			return m.monthlyModel.Init()
		}

	}
	return nil
}

func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch m.activeView {

		case viewMonthlyOverview:
			return m.handleMonthlyView(msg.String())
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		if m.monthlyModel != nil {
			updatedMonthlyModel, moCmd := m.monthlyModel.Update(msg)
			if mo, ok := updatedMonthlyModel.(*ui.MonthlyModel); ok {
				m.monthlyModel = mo
			}
			cmds = append(cmds, moCmd)
		}
		return m, tea.Batch(cmds...)
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
	default:
		viewContent = "Error: View not found or not initialized"
	}

	return viewContent
}
