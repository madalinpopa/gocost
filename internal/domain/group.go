package domain

// CategoryGroup holds one or more categories.
type CategoryGroup struct {
	GroupID   string `json:"groupId"`
	Order     int    `json:"order"`
	GroupName string `json:"groupName"`
}

// GroupRepository defines the interface for interacting with category group data.
type GroupRepository interface {
	// GetAll retrieves all category groups.
	GetAll() ([]CategoryGroup, error)

	// GetByID retrieves a single category group by its ID.
	GetByID(groupID string) (CategoryGroup, error)

	// Add saves a new category group.
	Add(group CategoryGroup) error

	// Update modifies an existing category group.
	Update(group CategoryGroup) error

	// Delete removes a category group using its ID.
	Delete(groupID string) error
}
