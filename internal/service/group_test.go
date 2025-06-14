package service

import (
	"errors"
	"testing"

	"github.com/madalinpopa/gocost/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockGroupRepo is a mock implementation of the GroupRepository.
type mockGroupRepo struct {
	groups []domain.CategoryGroup
	err    error
}

func (m *mockGroupRepo) GetAllGroups() ([]domain.CategoryGroup, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.groups, nil
}
func (m *mockGroupRepo) GetGroupByID(groupID string) (domain.CategoryGroup, error) {
	return domain.CategoryGroup{}, m.err
}
func (m *mockGroupRepo) AddGroup(group domain.CategoryGroup) error {
	if m.err != nil {
		return m.err
	}
	m.groups = append(m.groups, group)
	return nil
}
func (m *mockGroupRepo) UpdateGroup(group domain.CategoryGroup) error {
	return m.err
}
func (m *mockGroupRepo) DeleteGroup(groupID string) error {
	return m.err
}

func TestGroupService(t *testing.T) {
	mockGrp := domain.CategoryGroup{GroupID: "g1", GroupName: "Utilities"}
	mockRepo := &mockGroupRepo{
		groups: []domain.CategoryGroup{mockGrp},
	}
	service := NewGroupService(mockRepo)

	t.Run("GetAllGroups", func(t *testing.T) {
		groups, err := service.GetAllGroups()
		require.NoError(t, err)
		assert.Equal(t, []domain.CategoryGroup{mockGrp}, groups)
	})

	t.Run("AddGroup", func(t *testing.T) {
		newGroup := domain.CategoryGroup{GroupID: "g2", GroupName: "Housing"}
		err := service.AddGroup(newGroup)
		require.NoError(t, err)
		assert.Len(t, mockRepo.groups, 2)
	})

	t.Run("Handles Repository Error", func(t *testing.T) {
		errorRepo := &mockGroupRepo{err: errors.New("db error")}
		errorService := NewGroupService(errorRepo)
		_, err := errorService.GetAllGroups()
		require.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}
