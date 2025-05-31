package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/data"
)

type AppData struct {
	data         *data.DataRoot
	dataFilePath string
}

type AppViews struct{}

type App struct {
	width  int
	height int

	currentMonth time.Month
	currentYear  int

	AppData
	AppViews
}

func New(initialData *data.DataRoot, dataFilePath string) App {
	now := time.Now()
	currentM := now.Month()
	currentY := now.Year()

	return App{
		currentMonth: currentM,
		currentYear:  currentY,
		AppData: AppData{
			data:         initialData,
			dataFilePath: dataFilePath,
		},
	}
}

func (m App) Init() tea.Cmd {
	return nil
}

func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m App) View() string {
	return "Hello World"
}
