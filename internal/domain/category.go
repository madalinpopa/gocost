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
	// GetForMonth retrieves all categories for a specific month.
	GetForMonth(monthKey string) ([]Category, error)

	// Add saves a new category for a specific month.
	Add(monthKey string, category Category) error

	// Update modifies an existing category for a specific month.
	Update(monthKey string, category Category) error

	// Delete removes a category for a specific month using its ID.
	Delete(monthKey string, categoryID string) error

	// CopyFromMonth copies all categories from a previous month to the current one,
	// returning the number of categories copied.
	CopyFromMonth(fromMonthKey, toMonthKey string) (int, error)
}
