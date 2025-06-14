package domain

// Category represents the monthly expenses category.
type Category struct {
	CatID        string                   `json:"catId"`
	GroupID      string                   `json:"groupId"`
	CategoryName string                   `json:"categoryName"`
	Expense      map[string]ExpenseRecord `json:"expense"`
}

// CategoryRepository defines the interface for interacting with category data.
type CategoryRepository interface {
	GetCategoriesForMonth(monthKey string) ([]Category, error)
	AddCategory(monthKey string, category Category) error
	UpdateCategory(monthKey string, category Category) error
	DeleteCategory(monthKey string, categoryID string) error
	CopyCategoriesFromMonth(fromMonthKey, toMonthKey string) (int, error)
}
