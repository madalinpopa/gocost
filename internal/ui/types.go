package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
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
	CategoryModel      *CategoryModel
	IncomeModel        *IncomeModel
	IncomeFormModel    *IncomeFormModel
}

type ViewErrorMsg struct {
	Text  string
	Model tea.Model
}

type MonthlyViewMsg struct{}

type GroupDeleteMsg struct {
	Group data.CategoryGroup
}

type GroupAddMsg struct {
	Group data.CategoryGroup
}

type GroupUpdateMsg struct {
	Group data.CategoryGroup
}

type SelectGroupMsg struct{}

type SelectedGroupMsg struct {
	Group data.CategoryGroup
}

type CategoryAddMsg struct {
	MonthKey string
	Category data.Category
}

type CategoryUpdateMsg struct {
	Category data.Category
}

type CategoryDeleteMsg struct {
	Category data.Category
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
	MonthKey string
	Income   data.IncomeRecord
}

type DeleteIncomeMsg struct {
	MonthKey string
	Income   data.IncomeRecord
}
