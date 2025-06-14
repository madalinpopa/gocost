package data

import (
	"errors"
	"fmt"
	"sort"

	"github.com/madalinpopa/gocost/internal/domain"
)

// JsonRepository is a concrete implementation of the repository interfaces
// that uses a JSON file for storage.
type JsonRepository struct {
	filePath string
	data     *DataRoot
}

// NewJsonRepository creates and initializes a new JsonRepository.
// It loads data from the specified file path.
func NewJsonRepository(filePath string, defaultCurrency string) (*JsonRepository, error) {
	data, err := LoadData(filePath, defaultCurrency)
	if err != nil {
		return nil, err
	}

	return &JsonRepository{
		filePath: filePath,
		data:     data,
	}, nil
}

// save is a helper to persist the current state of r.data to the JSON file.
func (r *JsonRepository) save() error {
	return SaveData(r.filePath, r.data)
}

func (r *JsonRepository) GetAllGroups() ([]domain.CategoryGroup, error) {
	var groups []domain.CategoryGroup
	for _, group := range r.data.CategoryGroups {
		groups = append(groups, group)
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Order < groups[j].Order
	})

	return groups, nil
}

func (r *JsonRepository) GetGroupByID(groupID string) (domain.CategoryGroup, error) {
	group, ok := r.data.CategoryGroups[groupID]
	if !ok {
		return domain.CategoryGroup{}, errors.New("group not found")
	}
	return group, nil
}

func (r *JsonRepository) AddGroup(group domain.CategoryGroup) error {
	if _, exists := r.data.CategoryGroups[group.GroupID]; exists {
		return errors.New("group with this ID already exists")
	}
	r.data.CategoryGroups[group.GroupID] = group
	return r.save()
}

func (r *JsonRepository) UpdateGroup(group domain.CategoryGroup) error {
	if _, exists := r.data.CategoryGroups[group.GroupID]; !exists {
		return errors.New("group not found")
	}
	r.data.CategoryGroups[group.GroupID] = group
	return r.save()
}

func (r *JsonRepository) DeleteGroup(groupID string) error {
	// Check if any categories are using this group
	for _, monthRecord := range r.data.MonthlyData {
		for _, category := range monthRecord.Categories {
			if category.GroupID == groupID {
				group, _ := r.GetGroupByID(groupID)
				return fmt.Errorf("cannot delete group '%s': group is still being used by existing categories", group.GroupName)
			}
		}
	}

	if _, exists := r.data.CategoryGroups[groupID]; !exists {
		return errors.New("group not found")
	}
	delete(r.data.CategoryGroups, groupID)
	return r.save()
}

func (r *JsonRepository) GetIncomesForMonth(monthKey string) ([]domain.IncomeRecord, error) {
	if record, ok := r.data.MonthlyData[monthKey]; ok {
		return record.Incomes, nil
	}
	return []domain.IncomeRecord{}, nil
}

func (r *JsonRepository) AddIncome(monthKey string, income domain.IncomeRecord) error {
	monthRecord, ok := r.data.MonthlyData[monthKey]
	if !ok {
		monthRecord = domain.MonthlyRecord{
			Incomes:    make([]domain.IncomeRecord, 0),
			Categories: make([]domain.Category, 0),
		}
	}

	for _, existingIncome := range monthRecord.Incomes {
		if existingIncome.IncomeID == income.IncomeID {
			return fmt.Errorf("income record with ID %s already exists", income.IncomeID)
		}
	}

	monthRecord.Incomes = append(monthRecord.Incomes, income)
	r.data.MonthlyData[monthKey] = monthRecord
	return r.save()
}

func (r *JsonRepository) UpdateIncome(monthKey string, income domain.IncomeRecord) error {
	monthRecord, ok := r.data.MonthlyData[monthKey]
	if !ok {
		return fmt.Errorf("no data found for month %s", monthKey)
	}

	found := false
	for i, existingIncome := range monthRecord.Incomes {
		if existingIncome.IncomeID == income.IncomeID {
			monthRecord.Incomes[i] = income
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("income record with ID %s not found for update", income.IncomeID)
	}

	r.data.MonthlyData[monthKey] = monthRecord
	return r.save()
}

func (r *JsonRepository) DeleteIncome(monthKey string, incomeID string) error {
	monthRecord, ok := r.data.MonthlyData[monthKey]
	if !ok {
		return fmt.Errorf("no data found for month %s", monthKey)
	}

	found := false
	var updatedIncomes []domain.IncomeRecord
	for _, income := range monthRecord.Incomes {
		if income.IncomeID == incomeID {
			found = true
		} else {
			updatedIncomes = append(updatedIncomes, income)
		}
	}

	if !found {
		return fmt.Errorf("income record with ID %s not found for deletion", incomeID)
	}

	monthRecord.Incomes = updatedIncomes
	r.data.MonthlyData[monthKey] = monthRecord
	return r.save()
}

func (r *JsonRepository) GetCategoriesForMonth(monthKey string) ([]domain.Category, error) {
	if record, ok := r.data.MonthlyData[monthKey]; ok {
		return record.Categories, nil
	}
	return []domain.Category{}, nil
}

func (r *JsonRepository) AddCategory(monthKey string, category domain.Category) error {
	monthRecord, ok := r.data.MonthlyData[monthKey]
	if !ok {
		monthRecord = domain.MonthlyRecord{
			Incomes:    make([]domain.IncomeRecord, 0),
			Categories: make([]domain.Category, 0),
		}
	}

	monthRecord.Categories = append(monthRecord.Categories, category)
	r.data.MonthlyData[monthKey] = monthRecord
	return r.save()
}

func (r *JsonRepository) UpdateCategory(monthKey string, category domain.Category) error {
	monthRecord, ok := r.data.MonthlyData[monthKey]
	if !ok {
		return fmt.Errorf("no data found for month %s", monthKey)
	}

	found := false
	for i, existingCategory := range monthRecord.Categories {
		if existingCategory.CatID == category.CatID {
			monthRecord.Categories[i] = category
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("category with ID %s not found for update", category.CatID)
	}

	r.data.MonthlyData[monthKey] = monthRecord
	return r.save()
}

func (r *JsonRepository) DeleteCategory(monthKey string, categoryID string) error {
	monthRecord, ok := r.data.MonthlyData[monthKey]
	if !ok {
		return fmt.Errorf("no data found for month %s", monthKey)
	}

	found := false
	var updatedCategories []domain.Category
	for _, category := range monthRecord.Categories {
		if category.CatID == categoryID {
			found = true
		} else {
			updatedCategories = append(updatedCategories, category)
		}
	}

	if !found {
		return fmt.Errorf("category with ID %s not found for deletion", categoryID)
	}

	monthRecord.Categories = updatedCategories
	r.data.MonthlyData[monthKey] = monthRecord
	return r.save()
}

func (r *JsonRepository) CopyCategoriesFromMonth(fromMonthKey, toMonthKey string) (int, error) {
	prevRecord, exists := r.data.MonthlyData[fromMonthKey]
	if !exists || len(prevRecord.Categories) == 0 {
		return 0, fmt.Errorf("no categories found in %s to copy from", fromMonthKey)
	}

	var newCategories []domain.Category
	for _, category := range prevRecord.Categories {
		newCategory := domain.Category{
			CatID:        category.CatID,
			GroupID:      category.GroupID,
			CategoryName: category.CategoryName,
			Expense:      make(map[string]domain.ExpenseRecord), // Reset expenses
		}
		newCategories = append(newCategories, newCategory)
	}

	currentRecord := domain.MonthlyRecord{
		Incomes:    []domain.IncomeRecord{}, // Start with empty incomes
		Categories: newCategories,
	}

	if existingRecord, exists := r.data.MonthlyData[toMonthKey]; exists {
		currentRecord.Incomes = existingRecord.Incomes
	}

	r.data.MonthlyData[toMonthKey] = currentRecord

	err := r.save()
	if err != nil {
		return 0, err
	}

	return len(newCategories), nil
}
