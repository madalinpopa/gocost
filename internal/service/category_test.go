package service

import (
	"errors"
	"testing"

	"github.com/madalinpopa/gocost/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockCategoryRepo is a mock implementation of the CategoryRepository.
type mockCategoryRepo struct {
	categories []domain.Category
	err        error
}

func (m *mockCategoryRepo) GetCategoriesForMonth(monthKey string) ([]domain.Category, error) {
	_ = monthKey
	if m.err != nil {
		return nil, m.err
	}
	return m.categories, nil
}

func (m *mockCategoryRepo) AddCategory(monthKey string, category domain.Category) error {
	_ = monthKey
	if m.err != nil {
		return m.err
	}
	m.categories = append(m.categories, category)
	return nil
}

func (m *mockCategoryRepo) UpdateCategory(monthKey string, category domain.Category) error {
	_, _ = monthKey, category
	return m.err
}
func (m *mockCategoryRepo) DeleteCategory(monthKey string, categoryID string) error {
	_, _ = categoryID, monthKey
	return m.err
}
func (m *mockCategoryRepo) CopyCategoriesFromMonth(fromMonthKey, toMonthKey string) (int, error) {
	_ = fromMonthKey
	if m.err != nil {
		return 0, m.err
	}
	return len(m.categories), nil
}

func TestCategoryService(t *testing.T) {
	mockCat := domain.Category{CatID: "c1", CategoryName: "Test"}
	mockRepo := &mockCategoryRepo{
		categories: []domain.Category{mockCat},
	}
	service := NewCategoryService(mockRepo)

	t.Run("GetCategoriesForMonth", func(t *testing.T) {
		cats, err := service.GetCategoriesForMonth("any-month")
		require.NoError(t, err)
		assert.Equal(t, []domain.Category{mockCat}, cats)
	})

	t.Run("AddCategory", func(t *testing.T) {
		newCat := domain.Category{CatID: "c2", CategoryName: "New Test"}
		err := service.AddCategory("any-month", newCat)
		require.NoError(t, err)
		// Check if it was added to the mock's slice
		assert.Len(t, mockRepo.categories, 2)
		assert.Equal(t, "c2", mockRepo.categories[1].CatID)
	})

	t.Run("Handles Repository Error", func(t *testing.T) {
		errorRepo := &mockCategoryRepo{err: errors.New("db error")}
		errorService := NewCategoryService(errorRepo)
		_, err := errorService.GetCategoriesForMonth("any-month")
		require.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}
