package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Adaptive colors that work well in both light and dark themes
	ColorSubtleBorder  = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	ColorFocusedBorder = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	ColorHeaderText    = lipgloss.AdaptiveColor{Light: "#0E6BA8", Dark: "#04B575"}
	ColorMutedText     = lipgloss.AdaptiveColor{Light: "#6C6C6C", Dark: "#7D7D7D"}
	ColorSuccess       = lipgloss.AdaptiveColor{Light: "#047857", Dark: "#10B981"}
	ColorWarning       = lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#EF4444"}
	ColorAccent        = lipgloss.AdaptiveColor{Light: "#7C3AED", Dark: "#A855F7"}
	
	// Status colors
	ColorStatusPaid    = lipgloss.AdaptiveColor{Light: "#047857", Dark: "#10B981"}
	ColorStatusNotPaid = lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#EF4444"}
	
	// List item colors
	ColorFocusedListBg = lipgloss.AdaptiveColor{Light: "#E0E7FF", Dark: "#374151"}
	ColorFocusedListFg = lipgloss.AdaptiveColor{Light: "#1E293B", Dark: "#F3F4F6"}
	
	// Semantic background colors
	ColorGroupHeaderBg = lipgloss.AdaptiveColor{Light: "#F8FAFC", Dark: "#1F2937"}
	ColorActiveGroupBg = lipgloss.AdaptiveColor{Light: "#FEF3C7", Dark: "#454311"}
	ColorErrorBg      = lipgloss.AdaptiveColor{Light: "#FEF2F2", Dark: "#7F1D1D"}
	ColorSuccessBg    = lipgloss.AdaptiveColor{Light: "#F0FDF4", Dark: "#14532D"}
	ColorInfoBg       = lipgloss.AdaptiveColor{Light: "#EFF6FF", Dark: "#1E3A8A"}
	ColorInputBg      = lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#374151"}
	ColorInputBorder  = lipgloss.AdaptiveColor{Light: "#D1D5DB", Dark: "#4B5563"}

	// General styles
	AppStyle = lipgloss.NewStyle().Padding(1, 2)

	// Border styles
	FocusedBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorFocusedBorder)

	NormalBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSubtleBorder)
			
	TopBorder = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(ColorSubtleBorder)
			
	BottomBorder = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(ColorSubtleBorder)

	// Text styles
	HeaderText = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorHeaderText)

	MutedText = lipgloss.NewStyle().
			Foreground(ColorMutedText)
			
	AccentText = lipgloss.NewStyle().
			Foreground(ColorAccent)
			
	BoldText = lipgloss.NewStyle().
			Bold(true)

	// Status styles
	StatusPaid = lipgloss.NewStyle().
			Foreground(ColorStatusPaid).
			Bold(true)
			
	StatusNotPaid = lipgloss.NewStyle().
			Foreground(ColorStatusNotPaid).
			Bold(true)

	// List item styles
	FocusedListItem = lipgloss.NewStyle().
			Background(ColorFocusedListBg).
			Foreground(ColorFocusedListFg).
			Bold(true)
			
	NormalListItem = lipgloss.NewStyle()
	
	SelectedListItem = lipgloss.NewStyle().
			Background(ColorAccent).
			Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"}).
			Bold(true)

	// Group styles
	GroupHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorHeaderText)
			
	ActiveGroupStyle = lipgloss.NewStyle().
			Bold(false).
			Foreground(lipgloss.AdaptiveColor{Light: "#B45309", Dark: "#FBBF24"})

	// Button styles
	ButtonStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSubtleBorder)
			
	FocusedButtonStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorFocusedBorder).
			Background(ColorFocusedListBg).
			Foreground(ColorFocusedListFg).
			Bold(true)

	// Layout utility styles
	SpacerStyle = lipgloss.NewStyle()
	
	// Column alignment styles
	LeftAlign   = lipgloss.NewStyle().Align(lipgloss.Left)
	RightAlign  = lipgloss.NewStyle().Align(lipgloss.Right)
	CenterAlign = lipgloss.NewStyle().Align(lipgloss.Center)

	// Message styles
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorWarning).
			Background(ColorErrorBg).
			Padding(0, 1).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Background(ColorSuccessBg).
			Padding(0, 1).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Background(ColorInfoBg).
			Padding(0, 1)

	// Input field styles
	InputStyle = lipgloss.NewStyle().
			Background(ColorInputBg).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorInputBorder).
			Padding(0, 1)

	FocusedInputStyle = lipgloss.NewStyle().
			Background(ColorInputBg).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorFocusedBorder).
			Padding(0, 1)

	// Table styles
	TableHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorHeaderText).
			Background(ColorGroupHeaderBg).
			Padding(0, 1)

	TableCellStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// Highlight and emphasis styles
	HighlightStyle = lipgloss.NewStyle().
			Background(ColorFocusedListBg).
			Foreground(ColorFocusedListFg).
			Bold(true)

	EmphasisStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	// Progress and status indicators
	ProgressBarStyle = lipgloss.NewStyle().
			Background(ColorSubtleBorder)

	ProgressFillStyle = lipgloss.NewStyle().
			Background(ColorSuccess)
)

// Utility functions for common styling patterns

// CreateSpacer creates a spacer with the specified width
func CreateSpacer(width int) lipgloss.Style {
	return SpacerStyle.Width(width)
}

// CreateColumnSpacer creates a spacer for column layouts
func CreateColumnSpacer(spacing int) lipgloss.Style {
	return SpacerStyle.Width(spacing)
}

// CreateLeftAlignedColumn creates a left-aligned column with specified width
func CreateLeftAlignedColumn(width int) lipgloss.Style {
	return LeftAlign.Width(width)
}

// CreateRightAlignedColumn creates a right-aligned column with specified width
func CreateRightAlignedColumn(width int) lipgloss.Style {
	return RightAlign.Width(width)
}

// CreateCenterAlignedColumn creates a center-aligned column with specified width
func CreateCenterAlignedColumn(width int) lipgloss.Style {
	return CenterAlign.Width(width)
}

// CreateFooterStyle creates a footer style with specified width and padding
func CreateFooterStyle(width int) lipgloss.Style {
	return BottomBorder.
		Width(width - AppStyle.GetHorizontalPadding()).
		PaddingTop(1)
}

// RenderButton renders a button with appropriate styling based on focus state
func RenderButton(text string, focused bool) string {
	buttonText := "[ " + text + " ]"
	if focused {
		return FocusedButtonStyle.Render(buttonText)
	}
	return ButtonStyle.Render(buttonText)
}

// RenderStatusBadge renders a status badge with appropriate colors
func RenderStatusBadge(status string) string {
	badge := "[" + status + "]"
	switch status {
	case "Paid":
		return StatusPaid.Render(badge)
	case "Not Paid":
		return StatusNotPaid.Render(badge)
	default:
		return MutedText.Render(badge)
	}
}

// RenderErrorMessage renders an error message with appropriate styling
func RenderErrorMessage(message string) string {
	return ErrorStyle.Render("✗ " + message)
}

// RenderSuccessMessage renders a success message with appropriate styling
func RenderSuccessMessage(message string) string {
	return SuccessStyle.Render("✓ " + message)
}

// RenderInfoMessage renders an info message with appropriate styling
func RenderInfoMessage(message string) string {
	return InfoStyle.Render("ℹ " + message)
}

// RenderInput renders an input field with appropriate styling based on focus state
func RenderInput(content string, focused bool) string {
	if focused {
		return FocusedInputStyle.Render(content)
	}
	return InputStyle.Render(content)
}

// RenderTableHeader renders a table header cell
func RenderTableHeader(content string, width int) string {
	return TableHeaderStyle.Width(width).Render(content)
}

// RenderTableCell renders a table cell with specified alignment and width
func RenderTableCell(content string, width int, align lipgloss.Position) string {
	return TableCellStyle.Width(width).Align(align).Render(content)
}

// RenderProgressBar renders a progress bar with the given percentage (0-100)
func RenderProgressBar(percentage float64, width int) string {
	if width <= 0 {
		return ""
	}
	
	fillWidth := int(float64(width) * percentage / 100)
	if fillWidth > width {
		fillWidth = width
	}
	
	fill := ProgressFillStyle.Width(fillWidth).Render("")
	empty := ProgressBarStyle.Width(width - fillWidth).Render("")
	
	return lipgloss.JoinHorizontal(lipgloss.Left, fill, empty)
}

// RenderHighlight renders text with highlight styling
func RenderHighlight(text string) string {
	return HighlightStyle.Render(text)
}

// RenderEmphasis renders text with emphasis styling  
func RenderEmphasis(text string) string {
	return EmphasisStyle.Render(text)
}