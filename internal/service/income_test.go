package service

import (
	"errors"
	"testing"

	"github.com/madalinpopa/gocost/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockIncomeRepo is a mock implementation of the IncomeRepository.
type mockIncomeRepo struct {
	incomes []domain.IncomeRecord
	err     error
}

func (m *mockIncomeRepo) GetIncomesForMonth(monthKey string) ([]domain.IncomeRecord, error) {
	_ = monthKey
	if m.err != nil {
		return nil, m.err
	}
	return m.incomes, nil
}
func (m *mockIncomeRepo) AddIncome(monthKey string, income domain.IncomeRecord) error {
	_, _ = monthKey, income
	if m.err != nil {
		return m.err
	}
	m.incomes = append(m.incomes, income)
	return nil
}
func (m *mockIncomeRepo) UpdateIncome(monthKey string, income domain.IncomeRecord) error {
	_, _ = income, monthKey
	return m.err
}
func (m *mockIncomeRepo) DeleteIncome(monthKey string, incomeID string) error {
	_, _ = incomeID, monthKey
	return m.err
}

func TestIncomeService(t *testing.T) {
	mockInc := domain.IncomeRecord{IncomeID: "i1", Description: "Salary"}
	mockRepo := &mockIncomeRepo{
		incomes: []domain.IncomeRecord{mockInc},
	}
	service := NewIncomeService(mockRepo)

	t.Run("GetIncomesForMonth", func(t *testing.T) {
		incomes, err := service.GetIncomesForMonth("any-month")
		require.NoError(t, err)
		assert.Equal(t, []domain.IncomeRecord{mockInc}, incomes)
	})

	t.Run("AddIncome", func(t *testing.T) {
		newIncome := domain.IncomeRecord{IncomeID: "i2", Description: "Bonus"}
		err := service.AddIncome("any-month", newIncome)
		require.NoError(t, err)
		assert.Len(t, mockRepo.incomes, 2)
	})

	t.Run("Handles Repository Error", func(t *testing.T) {
		errorRepo := &mockIncomeRepo{err: errors.New("db error")}
		errorService := NewIncomeService(errorRepo)
		_, err := errorService.GetIncomesForMonth("any-month")
		require.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}
