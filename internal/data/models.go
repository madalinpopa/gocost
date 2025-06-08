package data

// IncomeRecord represents an income record
type IncomeRecord struct {
	IncomeID    string  `json:"incomeId"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

// ExpenseRecord represents an expense record
type ExpenseRecord struct {
	Budget string  `json:"budget"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
	Notes  string  `json:"notes"`
}

// Category represents the monthly expenses category
type Category struct {
	CatID        string `json:"catId"`
	GroupID      string `json:"groupId"`
	CategoryName string `json:"categoryName"`
	Expense      map[string]ExpenseRecord
}

// CategoryGroup holds one or more categories
type CategoryGroup struct {
	GroupID   string `json:"groupId"`
	Order     int    `json:"order"`
	GroupName string `json:"groupName"`
}

// MonthlyRecord holds one or more income and expense records
type MonthlyRecord struct {
	Incomes    []IncomeRecord `json:"incomes"`
	Categories []Category     `json:"categories"`
}

// DataRoot represents the root data structure
type DataRoot struct {
	DefaultCurrency string                   `json:"defaultCurrency"`
	CategoryGroups  map[string]CategoryGroup `json:"CategoryGroups"`
	MonthlyData     map[string]MonthlyRecord `json:"monthlyData"`
}

// NewDataRoot creates a new instance of DataRoot
func NewDataRoot() *DataRoot {
	return &DataRoot{
		CategoryGroups: make(map[string]CategoryGroup, 0),
		MonthlyData:    make(map[string]MonthlyRecord, 0),
	}
}
