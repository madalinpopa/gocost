package domain

// IncomeRecord represents an income record.
type IncomeRecord struct {
	IncomeID    string  `json:"incomeId"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

// IncomeRepository defines the interface for interacting with income data.
type IncomeRepository interface {
	// GetForMonth retrieves all income records for a specific month.
	GetForMonth(monthKey string) ([]IncomeRecord, error)

	// Add saves a new income record for a specific month.
	Add(monthKey string, income IncomeRecord) error

	// Update modifies an existing income record for a specific month.
	Update(monthKey string, income IncomeRecord) error

	// Delete removes an income record for a specific month using its ID.
	Delete(monthKey string, incomeID string) error
}
