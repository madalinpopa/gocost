package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/ui"
)

// StatusType represents the type of status message
type StatusType int

const (
	StatusSuccess StatusType = iota
	StatusError
)

// StatusClearMsg represents a message to clear the status message
type StatusClearMsg struct{}

// Status represents a status message
type Status struct {
	Message string
	Type    StatusType
}

// SetStatus updates the app's status message
func (m App) SetStatus(message string, statusType StatusType) (App, tea.Cmd) {
	m.statusMessage = message
	return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return StatusClearMsg{}
	})
}

// SetSuccessStatus sets a success status message
func (m App) SetSuccessStatus(message string) (App, tea.Cmd) {
	styledMessage := ui.StatusPaid.Render("✓ " + message)
	return m.SetStatus(styledMessage, StatusSuccess)
}

// SetErrorStatus sets an error status message
func (m App) SetErrorStatus(message string) (App, tea.Cmd) {
	styledMessage := ui.StatusNotPaid.Render("✗ " + message)
	return m.SetStatus(styledMessage, StatusError)
}

// ClearStatus clears the current status message
func (m App) ClearStatus() App {
	m.statusMessage = ""
	return m
}

// HasStatus returns true if there's a current status message
func (m App) HasStatus() bool {
	return m.statusMessage != ""
}

// GetStatusMessage returns the current status message
func (m App) GetStatusMessage() string {
	return m.statusMessage
}
