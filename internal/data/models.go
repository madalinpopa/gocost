package data

// Category represents the monthly expenses category
type Category struct {
	CatID        string `json:"catId"`
	GroupID      string `json:"groupId"`
	CategoryName string `json:"categoryName"`
}

// CategoryGroup holds one or more categories
type CategoryGroup struct {
	GroupID   string `json:"groupId"`
	GroupName string `json:"groupName"`
}

// MonthlyRecord holds one or more income and expense records
type MonthlyRecord struct {
	Incomes    []IncomeRecord  `json:"incomes"`
	Expenses   []ExpenseRecord `json:"expenses"`
	Categories []Category      `json:"categories"`
}

// IncomeRecord represents an income record
type IncomeRecord struct {
	IncomeID    string  `json:"incomeId"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

// ExpenseRecord represents an expense record
type ExpenseRecord struct {
	CatID  string  `json:"catId"`
	Budget string  `json:"budget"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
	Notes  string  `json:"notes"`
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
