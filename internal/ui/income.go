package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/madalinpopa/gocost/internal/config"
	"github.com/madalinpopa/gocost/internal/domain"
	"github.com/spf13/viper"
)

type IncomeModel struct {
	WindowSize
	MonthYear

	cursor   int
	monthKey string
	incomes  []domain.IncomeRecord

	viewport viewport.Model
	ready    bool
}

// NewIncomeModel creates a new IncomeModel instance.
func NewIncomeModel(incomes []domain.IncomeRecord, month time.Month, year int) IncomeModel {
	monthKey := GetMonthKey(month, year)

	return IncomeModel{
		incomes:  incomes,
		monthKey: monthKey,
		MonthYear: MonthYear{
			CurrentMonth: month,
			CurrentYear:  year,
		},
		viewport: viewport.New(70, 20),
		ready:    false,
	}
}

// Init initializes the IncomeModel.
func (m IncomeModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the IncomeModel state.
func (m IncomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			availableHeight := msg.Height - verticalMarginHeight - 4 // -4 for padding (2) and newlines (2)
			viewportHeight := m.calculateViewportHeight(availableHeight)
			m.viewport = viewport.New(msg.Width, viewportHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.getIncomesContent())
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			availableHeight := msg.Height - verticalMarginHeight - 4 // -4 for padding (2) and newlines (2)
			viewportHeight := m.calculateViewportHeight(availableHeight)
			m.viewport.Height = viewportHeight
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "esc":
			return m, func() tea.Msg { return MonthlyViewMsg{} }

		case "j", "down":
			if len(m.incomes) > 0 {
				m.cursor++
				if m.cursor >= len(m.incomes) {
					m.cursor = 0
				}
				m = m.ensureCursorVisible()
			}
			return m, nil

		case "k", "up":
			if len(m.incomes) > 0 {
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(m.incomes) - 1
				}
				m = m.ensureCursorVisible()
			}
			return m, nil

		case "a", "n":
			return m, func() tea.Msg {
				return AddIncomeFormMsg{MonthKey: m.monthKey}
			}

		case "e", "enter":
			if len(m.incomes) > 0 && m.cursor >= 0 && m.cursor < len(m.incomes) {
				incomeRecord := m.incomes[m.cursor]
				return m, func() tea.Msg {
					return EditIncomeMsg{
						MonthKey: m.monthKey,
						Income:   incomeRecord,
					}
				}
			}

		case "d":
			if len(m.incomes) > 0 && m.cursor >= 0 && m.cursor < len(m.incomes) {
				incomeRecord := m.incomes[m.cursor]
				return m, func() tea.Msg {
					return DeleteIncomeMsg{
						MonthKey: m.monthKey,
						Income:   incomeRecord,
					}
				}
			}
		}
		return m, nil
	}

	if m.ready {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the IncomeModel.
func (m IncomeModel) View() string {
	if !m.ready {
		return AppStyle.Width(m.Width).Height(m.Height).Render("\n  Initializing...")
	}

	if m.ready {
		m.viewport.SetContent(m.getIncomesContent())
	}

	var b strings.Builder
	b.WriteString(m.headerView())
	b.WriteString("\n")
	b.WriteString(m.viewport.View())
	b.WriteString("\n")
	b.WriteString(m.footerView())
	return AppStyle.Render(b.String())
}

// headerView renders the header section of the view.
func (m IncomeModel) headerView() string {
	title := fmt.Sprintf("Manage Income - %s %d", m.CurrentMonth.String(), m.CurrentYear)
	var b strings.Builder
	b.WriteString(HeaderText.Render(title))
	b.WriteString("\n")
	return b.String()
}

// footerView renders the footer section with key hints.
func (m IncomeModel) footerView() string {
	var b strings.Builder
	b.WriteString("\n")
	keyHints := "(j/k: Nav, a/n: Add, e/Enter: Edit, d: Delete, Esc/q: Back)"
	b.WriteString(MutedText.Render(keyHints))
	return b.String()
}

// getIncomesContent generates the content for the viewport.
func (m IncomeModel) getIncomesContent() string {
	var b strings.Builder

	if len(m.incomes) == 0 {
		b.WriteString(MutedText.Render("No income entries for this month."))
	} else {
		for i, entry := range m.incomes {
			lineStyle := NormalListItem
			prefix := "  "
			if i == m.cursor {
				lineStyle = FocusedListItem
				prefix = "> "
			}
			line := fmt.Sprintf("%s%s: %.2f %s",
				prefix,
				entry.Description,
				entry.Amount,
				viper.GetString(config.CurrencyField),
			)
			b.WriteString(lineStyle.Render(line))
			b.WriteString("\n")
		}
	}
	return b.String()
}

// calculateViewportHeight calculates the appropriate height for the viewport.
func (m IncomeModel) calculateViewportHeight(availableHeight int) int {
	desiredHeight := max(len(m.incomes)+1, 1)
	return min(desiredHeight, max(1, availableHeight))
}

// ensureCursorVisible ensures the cursor is visible in the viewport.
func (m IncomeModel) ensureCursorVisible() IncomeModel {
	if !m.ready {
		return m
	}

	content := m.getIncomesContent()
	m.viewport.SetContent(content)

	if len(m.incomes) == 0 {
		return m
	}

	viewportTop := m.viewport.YOffset
	viewportBottom := viewportTop + m.viewport.Height - 1

	if m.cursor > viewportBottom {
		newOffset := max(m.cursor-m.viewport.Height+1, 0)
		m.viewport.SetYOffset(newOffset)
	}
	if m.cursor < viewportTop {
		m.viewport.SetYOffset(m.cursor)
	}
	return m
}

// updateViewportHeight updates the viewport height based on current window size.
func (m IncomeModel) updateViewportHeight() IncomeModel {
	if !m.ready {
		return m
	}

	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())

	verticalMarginHeight := headerHeight + footerHeight
	availableHeight := m.Height - verticalMarginHeight - 4 // -4 for padding (2) and newlines (2)
	viewportHeight := m.calculateViewportHeight(availableHeight)
	m.viewport.Height = viewportHeight
	return m
}

// SetMonthYear updates the current month/year and loads corresponding income entries.
func (m IncomeModel) SetMonthYear(month time.Month, year int) IncomeModel {
	m.CurrentMonth = month
	m.CurrentYear = year
	m.monthKey = GetMonthKey(month, year)
	m.cursor = 0 // Reset cursor
	return m
}

// UpdateData refreshes the model with new data and resets state.
func (m IncomeModel) UpdateData(incomes []domain.IncomeRecord) IncomeModel {
	m.incomes = incomes
	if m.cursor >= len(m.incomes) && len(m.incomes) > 0 {
		m.cursor = len(m.incomes) - 1
	} else if len(m.incomes) == 0 {
		m.cursor = 0
	}

	// Update viewport height when income data changes
	m = m.updateViewportHeight()
	if m.ready {
		m.viewport.SetContent(m.getIncomesContent())
		m.viewport.GotoTop()
	}

	return m
}
