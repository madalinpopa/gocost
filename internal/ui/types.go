package ui

import (
	"time"

	"github.com/madalinpopa/gocost/internal/data"
)

// AppData represents the application data.
type AppData struct {
	Data     *data.DataRoot
	FilePath string
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
	MonthlyModel       *MonthlyModel
	CategoryGroupModel *CategoryGroupModel
	IncomeModel        *IncomeModel
	IncomeFormModel    *IncomeFormModel
}
