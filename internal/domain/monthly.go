package domain

// ExpenseRecord represents an expense record.
type ExpenseRecord struct {
	Budget float64 `json:"budget"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
	Notes  string  `json:"notes"`
}

// MonthlyRecord holds one or more income and expense records.
type MonthlyRecord struct {
	Incomes    []IncomeRecord `json:"incomes"`
	Categories []Category     `json:"categories"`
}
