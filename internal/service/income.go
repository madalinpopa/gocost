package service

import "github.com/madalinpopa/gocost/internal/domain"

// IncomeService encapsulates business logic for income records.
type IncomeService struct {
	repo domain.IncomeRepository
}

// NewIncomeService creates a new IncomeService.
func NewIncomeService(r domain.IncomeRepository) *IncomeService {
	return &IncomeService{repo: r}
}

// GetIncomesForMonth retrieves all income records for a given month.
func (s *IncomeService) GetIncomesForMonth(monthKey string) ([]domain.IncomeRecord, error) {
	return s.repo.GetIncomesForMonth(monthKey)
}

// AddIncome adds a new income record for a given month.
func (s *IncomeService) AddIncome(monthKey string, income domain.IncomeRecord) error {
	return s.repo.AddIncome(monthKey, income)
}

// UpdateIncome updates an existing income record for a given month.
func (s *IncomeService) UpdateIncome(monthKey string, income domain.IncomeRecord) error {
	return s.repo.UpdateIncome(monthKey, income)
}

// DeleteIncome deletes an income record for a given month by its ID.
func (s *IncomeService) DeleteIncome(monthKey string, incomeID string) error {
	return s.repo.DeleteIncome(monthKey, incomeID)
}
