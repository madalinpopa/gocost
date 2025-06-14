package service

import "github.com/madalinpopa/gocost/internal/domain"

// GroupService encapsulates business logic for category groups.
type GroupService struct {
	repo domain.GroupRepository
}

// NewGroupService creates a new GroupService.
func NewGroupService(r domain.GroupRepository) *GroupService {
	return &GroupService{repo: r}
}

// GetAllGroups retrieves all category groups.
func (s *GroupService) GetAllGroups() ([]domain.CategoryGroup, error) {
	return s.repo.GetAllGroups()
}

// GetGroupByID retrieves a single category group by its ID.
func (s *GroupService) GetGroupByID(groupID string) (domain.CategoryGroup, error) {
	return s.repo.GetGroupByID(groupID)
}

// AddGroup adds a new category group.
func (s *GroupService) AddGroup(group domain.CategoryGroup) error {
	// Future validation logic can be added here.
	return s.repo.AddGroup(group)
}

// UpdateGroup updates an existing category group.
func (s *GroupService) UpdateGroup(group domain.CategoryGroup) error {
	return s.repo.UpdateGroup(group)
}

// DeleteGroup deletes a category group by its ID.
func (s *GroupService) DeleteGroup(groupID string) error {
	return s.repo.DeleteGroup(groupID)
}
