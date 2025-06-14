package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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
	}
}

// Init initializes the IncomeModel.
func (m IncomeModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the IncomeModel state.
func (m IncomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
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
			}

		case "k", "up":
			if len(m.incomes) > 0 {
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(m.incomes) - 1
				}
			}

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

	}

	return m, nil
}

// View renders the IncomeModel.
func (m IncomeModel) View() string {
	var b strings.Builder

	title := fmt.Sprintf("Manage Income - %s", m.monthKey)
	b.WriteString(HeaderText.Render(title))
	b.WriteString("\n\n")

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

	b.WriteString("\n\n")
	keyHints := "(j/k: Nav, a/n: Add, e/Enter: Edit, d: Delete, Esc/q: Back)"
	b.WriteString(MutedText.Render(keyHints))

	viewStr := AppStyle.Width(m.Width).Height(m.Height - 3).Render(b.String())
	return viewStr
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
	return m
}
