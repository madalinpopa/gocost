package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/madalinpopa/gocost/internal/domain"
)

// jsonStore represents the root data structure, specific to the JSON file.
// It is an unexported implementation detail of the JsonRepository.
type jsonStore struct {
	DefaultCurrency string                          `json:"defaultCurrency"`
	CategoryGroups  map[string]domain.CategoryGroup `json:"CategoryGroups"`
	MonthlyData     map[string]domain.MonthlyRecord `json:"monthlyData"`
}

// newJsonStore creates a new instance of jsonStore.
func newJsonStore() *jsonStore {
	return &jsonStore{
		CategoryGroups: make(map[string]domain.CategoryGroup, 0),
		MonthlyData:    make(map[string]domain.MonthlyRecord, 0),
	}
}

// JsonRepository is a concrete implementation of the repository interfaces
// that uses a JSON file for storage.
type JsonRepository struct {
	filePath string
	store    *jsonStore
}

// NewJsonRepository creates and initializes a new JsonRepository.
// It loads data from the specified file path.
func NewJsonRepository(filePath string, defaultCurrency string) (*JsonRepository, error) {
	store, err := loadData(filePath, defaultCurrency)
	if err != nil {
		return nil, err
	}

	return &JsonRepository{
		filePath: filePath,
		store:    store,
	}, nil
}

// save is a helper to persist the current state of r.store to the JSON file.
func (r *JsonRepository) save() error {
	return saveData(r.filePath, r.store)
}

// --- GroupRepository Implementation ---

func (r *JsonRepository) GetAllGroups() ([]domain.CategoryGroup, error) {
	var groups []domain.CategoryGroup
	for _, group := range r.store.CategoryGroups {
		groups = append(groups, group)
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Order < groups[j].Order
	})
	return groups, nil
}

func (r *JsonRepository) GetGroupByID(groupID string) (domain.CategoryGroup, error) {
	group, ok := r.store.CategoryGroups[groupID]
	if !ok {
		return domain.CategoryGroup{}, errors.New("group not found")
	}
	return group, nil
}

func (r *JsonRepository) AddGroup(group domain.CategoryGroup) error {
	if _, exists := r.store.CategoryGroups[group.GroupID]; exists {
		return errors.New("group with this ID already exists")
	}
	r.store.CategoryGroups[group.GroupID] = group
	return r.save()
}

func (r *JsonRepository) UpdateGroup(group domain.CategoryGroup) error {
	if _, exists := r.store.CategoryGroups[group.GroupID]; !exists {
		return errors.New("group not found")
	}
	r.store.CategoryGroups[group.GroupID] = group
	return r.save()
}

func (r *JsonRepository) DeleteGroup(groupID string) error {
	for _, monthRecord := range r.store.MonthlyData {
		for _, category := range monthRecord.Categories {
			if category.GroupID == groupID {
				group, _ := r.GetGroupByID(groupID)
				return fmt.Errorf("cannot delete group '%s': group is still being used by existing categories", group.GroupName)
			}
		}
	}
	if _, exists := r.store.CategoryGroups[groupID]; !exists {
		return errors.New("group not found")
	}
	delete(r.store.CategoryGroups, groupID)
	return r.save()
}

// --- IncomeRepository Implementation ---

func (r *JsonRepository) GetIncomesForMonth(monthKey string) ([]domain.IncomeRecord, error) {
	if record, ok := r.store.MonthlyData[monthKey]; ok {
		return record.Incomes, nil
	}
	return []domain.IncomeRecord{}, nil
}

func (r *JsonRepository) AddIncome(monthKey string, income domain.IncomeRecord) error {
	monthRecord, ok := r.store.MonthlyData[monthKey]
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
	r.store.MonthlyData[monthKey] = monthRecord
	return r.save()
}

func (r *JsonRepository) UpdateIncome(monthKey string, income domain.IncomeRecord) error {
	monthRecord, ok := r.store.MonthlyData[monthKey]
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
	r.store.MonthlyData[monthKey] = monthRecord
	return r.save()
}

func (r *JsonRepository) DeleteIncome(monthKey string, incomeID string) error {
	monthRecord, ok := r.store.MonthlyData[monthKey]
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
	r.store.MonthlyData[monthKey] = monthRecord
	return r.save()
}

// --- CategoryRepository Implementation ---

func (r *JsonRepository) GetCategoriesForMonth(monthKey string) ([]domain.Category, error) {
	if record, ok := r.store.MonthlyData[monthKey]; ok {
		return record.Categories, nil
	}
	return []domain.Category{}, nil
}

func (r *JsonRepository) AddCategory(monthKey string, category domain.Category) error {
	monthRecord, ok := r.store.MonthlyData[monthKey]
	if !ok {
		monthRecord = domain.MonthlyRecord{
			Incomes:    make([]domain.IncomeRecord, 0),
			Categories: make([]domain.Category, 0),
		}
	}
	monthRecord.Categories = append(monthRecord.Categories, category)
	r.store.MonthlyData[monthKey] = monthRecord
	return r.save()
}

func (r *JsonRepository) UpdateCategory(monthKey string, category domain.Category) error {
	monthRecord, ok := r.store.MonthlyData[monthKey]
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
	r.store.MonthlyData[monthKey] = monthRecord
	return r.save()
}

func (r *JsonRepository) DeleteCategory(monthKey string, categoryID string) error {
	monthRecord, ok := r.store.MonthlyData[monthKey]
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
	r.store.MonthlyData[monthKey] = monthRecord
	return r.save()
}

func (r *JsonRepository) CopyCategoriesFromMonth(fromMonthKey, toMonthKey string) (int, error) {
	prevRecord, exists := r.store.MonthlyData[fromMonthKey]
	if !exists || len(prevRecord.Categories) == 0 {
		return 0, fmt.Errorf("no categories found in %s to copy from", fromMonthKey)
	}
	var newCategories []domain.Category
	for _, category := range prevRecord.Categories {
		newCategory := domain.Category{
			CatID:        category.CatID,
			GroupID:      category.GroupID,
			CategoryName: category.CategoryName,
			Expense:      make(map[string]domain.ExpenseRecord),
		}
		newCategories = append(newCategories, newCategory)
	}
	currentRecord := domain.MonthlyRecord{
		Incomes:    []domain.IncomeRecord{},
		Categories: newCategories,
	}
	if existingRecord, exists := r.store.MonthlyData[toMonthKey]; exists {
		currentRecord.Incomes = existingRecord.Incomes
	}
	r.store.MonthlyData[toMonthKey] = currentRecord
	err := r.save()
	if err != nil {
		return 0, err
	}
	return len(newCategories), nil
}

// --- Unexported persistence helpers ---

func loadData(filePath string, currency string) (*jsonStore, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return newJsonStore(), nil
		}
		return nil, err
	}
	if len(fileData) == 0 {
		return newJsonStore(), nil
	}
	var store jsonStore
	err = json.Unmarshal(fileData, &store)
	if err != nil {
		return nil, err
	}
	if store.CategoryGroups == nil {
		store.CategoryGroups = make(map[string]domain.CategoryGroup, 0)
	}
	if store.MonthlyData == nil {
		store.MonthlyData = make(map[string]domain.MonthlyRecord, 0)
	}
	store.DefaultCurrency = currency
	return &store, nil
}

func saveData(filePath string, store *jsonStore) error {
	jsonData, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}
	return nil
}
