package ui

import (
	"fmt"
	"time"

	"github.com/google/uuid"
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

// GetPreviousMonth returns the previous month and year given the current month and year.
func GetPreviousMonth(year int, month time.Month) (int, time.Month) {
	currentTime := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	prevMonthTime := currentTime.AddDate(0, -1, 0)
	return prevMonthTime.Year(), prevMonthTime.Month()
}

// GetNextMonth returns the next month and year given the current month and year.
func GetNextMonth(year int, month time.Month) (int, time.Month) {
	currentTime := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	nextMonthTime := currentTime.AddDate(0, 1, 0)
	return nextMonthTime.Year(), nextMonthTime.Month()
}

// GetMonthKey returns the current month key
func GetMonthKey(month time.Month, year int) string {
	return fmt.Sprintf("%s-%d", month.String(), year)
}

// GenerateID generates a unique ID.
func GenerateID() string {
	return uuid.NewString()
}
