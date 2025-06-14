package domain

// CategoryGroup holds one or more categories.
type CategoryGroup struct {
	GroupID   string `json:"groupId"`
	Order     int    `json:"order"`
	GroupName string `json:"groupName"`
}

// GroupRepository defines the interface for interacting with category group data.
type GroupRepository interface {
	GetAllGroups() ([]CategoryGroup, error)
	GetGroupByID(groupID string) (CategoryGroup, error)
	AddGroup(group CategoryGroup) error
	UpdateGroup(group CategoryGroup) error
	DeleteGroup(groupID string) error
}
