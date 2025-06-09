package ui

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GetPreviousMonth returns the year and month for the month before the given month and year.
func GetPreviousMonth(year int, month time.Month) (int, time.Month) {
	currentTime := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	prevMonthTime := currentTime.AddDate(0, -1, 0)
	return prevMonthTime.Year(), prevMonthTime.Month()
}

// GetNextMonth returns the year and month for the month after the given month and year.
func GetNextMonth(year int, month time.Month) (int, time.Month) {
	currentTime := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	nextMonthTime := currentTime.AddDate(0, 1, 0)
	return nextMonthTime.Year(), nextMonthTime.Month()
}

// GetMonthKey returns a string key in the format "Month-Year" for the given month and year.
func GetMonthKey(month time.Month, year int) string {
	return fmt.Sprintf("%s-%d", month.String(), year)
}

// GenerateID generates a unique UUID string.
func GenerateID() string {
	return uuid.NewString()
}
