package service

import "github.com/madalinpopa/gocost/internal/domain"

// CategoryService encapsulates business logic for categories.
type CategoryService struct {
	repo domain.CategoryRepository
}

// NewCategoryService creates a new CategoryService.
func NewCategoryService(r domain.CategoryRepository) *CategoryService {
	return &CategoryService{repo: r}
}

// GetCategoriesForMonth retrieves all categories for a given month.
func (s *CategoryService) GetCategoriesForMonth(monthKey string) ([]domain.Category, error) {
	return s.repo.GetCategoriesForMonth(monthKey)
}

// AddCategory adds a new category for a given month.
func (s *CategoryService) AddCategory(monthKey string, category domain.Category) error {
	// Future validation logic could go here.
	return s.repo.AddCategory(monthKey, category)
}

// UpdateCategory updates an existing category for a given month.
func (s *CategoryService) UpdateCategory(monthKey string, category domain.Category) error {
	return s.repo.UpdateCategory(monthKey, category)
}

// DeleteCategory deletes a category for a given month by its ID.
func (s *CategoryService) DeleteCategory(monthKey string, categoryID string) error {
	return s.repo.DeleteCategory(monthKey, categoryID)
}

// CopyCategoriesFromMonth copies categories from a previous month to a new one.
func (s *CategoryService) CopyCategoriesFromMonth(fromMonthKey, toMonthKey string) (int, error) {
	return s.repo.CopyCategoriesFromMonth(fromMonthKey, toMonthKey)
}
