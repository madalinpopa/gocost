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
	AppViews
	activeView currentView
}

func New(initialData *data.DataRoot, dataFilePath string) App {
	now := time.Now()

	return App{
		Data: ui.Data{
			Data:         initialData,
			DataFilePath: dataFilePath,
		},
		MonthYear: ui.MonthYear{
			CurrentMonth: now.Month(),
			CurrentYear:  now.Year(),
		},
		AppViews: AppViews{
			monthlyModel: new(ui.MonthlyModel),
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

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch m.activeView {

		case viewMonthlyOverview:
			switch msg.String() {

			case "q":
				return m, tea.Quit
			}
		}
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
	}

	return viewContent
}
