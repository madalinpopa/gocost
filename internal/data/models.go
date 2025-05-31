package data

// Category represents the monthly expenses category
type Category struct {
	CatID        string `json:"catId"`
	CategoryName string `json:"categoryName"`
}

// CategoryGroup holds one or more categories
type CategoryGroup struct {
	GroupID    string     `json:"groupId"`
	GroupName  string     `json:"groupName"`
	Categories []Category `json:"categories"`
}

// MonthlyRecord holds one or more income and expense records
type MonthlyRecord struct {
	Incomes  []IncomeRecord  `json:"incomes"`
	Expenses []ExpenseRecord `json:"expenses"`
}

// IncomeRecord represents an income record
type IncomeRecord struct {
	IncomeID    string  `json:"incomeId"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

// ExpenseRecord represents an expense record
type ExpenseRecord struct {
	MasterCatID string  `json:"masterCatId"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"`
	Notes       string  `json:"notes"`
}

// DataRoot represents the root data structure
type DataRoot struct {
	DefaultCurrency      string                   `json:"defaultCurrency"`
	MasterCategoryGroups []CategoryGroup          `json:"masterCategoryGroups"`
	MonthlyData          map[string]MonthlyRecord `json:"monthlyData"`
}
