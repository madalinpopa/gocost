package domain

// IncomeRecord represents an income record.
type IncomeRecord struct {
	IncomeID    string  `json:"incomeId"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

// IncomeRepository defines the interface for interacting with income data.
type IncomeRepository interface {
	GetIncomesForMonth(monthKey string) ([]IncomeRecord, error)
	AddIncome(monthKey string, income IncomeRecord) error
	UpdateIncome(monthKey string, income IncomeRecord) error
	DeleteIncome(monthKey string, incomeID string) error
}
