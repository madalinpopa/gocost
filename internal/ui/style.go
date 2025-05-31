package ui

import "github.com/charmbracelet/lipgloss"

var (
	ColorSubtleBorder  = lipgloss.Color("240") // Example: Grey for normal borders
	ColorFocusedBorder = lipgloss.Color("63")  // Example: Purple for focused borders
	ColorHeaderText    = lipgloss.Color("212")
	ColorMutedText     = lipgloss.Color("240")
	ColorStatusPaid    = lipgloss.Color("78")
	ColorStatusNotPaid = lipgloss.Color("203")
	ColorFocusedListBg = lipgloss.Color("236")
	ColorFocusedListFg = lipgloss.Color("252")

	// General
	AppStyle = lipgloss.NewStyle().Padding(1, 2)

	// Borders
	FocusedBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorFocusedBorder)

	NormalBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSubtleBorder)

	// Text
	HeaderText = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorHeaderText)

	MutedText = lipgloss.NewStyle().Foreground(ColorMutedText)

	// Status Colors
	StatusPaid    = lipgloss.NewStyle().Foreground(ColorStatusPaid)
	StatusNotPaid = lipgloss.NewStyle().Foreground(ColorStatusNotPaid)

	FocusedListItem = lipgloss.NewStyle().Background(ColorFocusedListBg).Foreground(ColorFocusedListFg)
	NormalListItem  = lipgloss.NewStyle()
)
