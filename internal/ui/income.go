package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/config"
	"github.com/madalinpopa/gocost/internal/data"
	"github.com/spf13/viper"
)

type IncomeModel struct {
	AppData
	WindowSize
	MonthYear

	cursor        int
	monthKey      string
	incomeEntries []data.IncomeRecord
}

func NewIncomeModel(initialData *data.DataRoot, month time.Month, year int) IncomeModel {
	monthKey := GetMonthKey(month, year)
	var incomeEntries []data.IncomeRecord

	if incomes, ok := initialData.MonthlyData[monthKey]; ok {
		incomeEntries = incomes.Incomes
	} else {
		incomeEntries = make([]data.IncomeRecord, 0)
	}

	return IncomeModel{
		incomeEntries: incomeEntries,
		monthKey:      monthKey,
		MonthYear: MonthYear{
			CurrentMonth: month,
			CurrentYear:  year,
		},
		AppData: AppData{
			Data: initialData,
		},
	}
}

func (m IncomeModel) Init() tea.Cmd {
	return nil
}

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
			if len(m.incomeEntries) > 0 {
				m.cursor++
				if m.cursor >= len(m.incomeEntries) {
					m.cursor = 0
				}
			}

		case "k", "up":
			if len(m.incomeEntries) > 0 {
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(m.incomeEntries) - 1
				}
			}

		case "a", "n":
			return m, func() tea.Msg {
				return AddIncomeFormMsg{MonthKey: m.monthKey}
			}

		case "e", "enter":
			if len(m.incomeEntries) > 0 && m.cursor >= 0 && m.cursor < len(m.incomeEntries) {
				incomeRecord := m.incomeEntries[m.cursor]
				return m, func() tea.Msg {
					return EditIncomeMsg{
						MonthKey: m.monthKey,
						Income:   incomeRecord,
					}
				}
			}

		case "d":
			if len(m.incomeEntries) > 0 && m.cursor >= 0 && m.cursor < len(m.incomeEntries) {
				incomeRecord := m.incomeEntries[m.cursor]
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

func (m IncomeModel) View() string {
	var b strings.Builder

	title := fmt.Sprintf("Manage Income - %s", m.monthKey)
	b.WriteString(HeaderText.Render(title))
	b.WriteString("\n\n")

	if len(m.incomeEntries) == 0 {
		b.WriteString(MutedText.Render("No income entries for this month."))
	} else {
		for i, entry := range m.incomeEntries {
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

func (m IncomeModel) SetMonthYear(month time.Month, year int) IncomeModel {

	m.CurrentMonth = month
	m.CurrentYear = year
	m.monthKey = GetMonthKey(month, year)

	if monthRecord, ok := m.Data.MonthlyData[m.monthKey]; ok {
		m.incomeEntries = monthRecord.Incomes
	} else {
		m.incomeEntries = make([]data.IncomeRecord, 0)
	}
	m.cursor = 0 // Reset cursor

	return m
}
