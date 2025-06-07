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

type StatusClearMsg struct{}

type MonthlyViewMsg struct{}

type GroupDeleteMsg struct {
	GroupID string
}

type GroupAddMsg struct {
	Group data.CategoryGroup
}

type GroupUpdateMsg struct {
	Group data.CategoryGroup
}

type GroupManageCategoriesMsg struct {
	Group data.CategoryGroup
}

type IncomeViewMsg struct{}

type AddIncomeFormMsg struct {
	MonthKey string
}

type SaveIncomeMsg struct {
	MonthKey string
	Income   data.IncomeRecord
}

type EditIncomeMsg struct {
	MonthKey     string
	IncomeRecord data.IncomeRecord
}

type DeleteIncomeMsg struct {
	MonthKey     string
	IncomeRecord data.IncomeRecord
}
