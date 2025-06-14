package data

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/madalinpopa/gocost/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRepo is a helper function to create a new repository in a temporary directory for testing.
func setupTestRepo(t *testing.T) *JsonRepository {
	t.Helper()
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_data.json")

	repo, err := NewJsonRepository(filePath, "USD")
	require.NoError(t, err)
	return repo
}

func TestJsonRepository_GroupOperations(t *testing.T) {
	repo := setupTestRepo(t)
	group1 := domain.CategoryGroup{GroupID: "g1", GroupName: "Group 1", Order: 1}
	group2 := domain.CategoryGroup{GroupID: "g2", GroupName: "Group 2", Order: 2}

	t.Run("Add and Get Group", func(t *testing.T) {
		// Add
		err := repo.AddGroup(group1)
		require.NoError(t, err)
		err = repo.AddGroup(group2)
		require.NoError(t, err)

		// GetGroupByID
		retrieved, err := repo.GetGroupByID("g1")
		require.NoError(t, err)
		assert.Equal(t, group1, retrieved)

		// GetAllGroups
		all, err := repo.GetAllGroups()
		require.NoError(t, err)
		assert.Len(t, all, 2)
		assert.Equal(t, "Group 1", all[0].GroupName) // Check order
	})

	t.Run("Update Group", func(t *testing.T) {
		updatedGroup := domain.CategoryGroup{GroupID: "g1", GroupName: "Group 1 Updated", Order: 1}
		err := repo.UpdateGroup(updatedGroup)
		require.NoError(t, err)

		retrieved, err := repo.GetGroupByID("g1")
		require.NoError(t, err)
		assert.Equal(t, "Group 1 Updated", retrieved.GroupName)
	})

	t.Run("Delete Group", func(t *testing.T) {
		// Test delete failure when in use
		cat := domain.Category{CatID: "c1", GroupID: "g2", CategoryName: "Test Cat"}
		err := repo.AddCategory("May-2024", cat)
		require.NoError(t, err)

		err = repo.DeleteGroup("g2")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "group is still being used")

		// Test successful delete
		err = repo.DeleteGroup("g1")
		require.NoError(t, err)

		_, err = repo.GetGroupByID("g1")
		assert.Error(t, err)

		all, err := repo.GetAllGroups()
		require.NoError(t, err)
		assert.Len(t, all, 1)
	})
}

func TestJsonRepository_IncomeOperations(t *testing.T) {
	repo := setupTestRepo(t)
	monthKey := "June-2024"
	income1 := domain.IncomeRecord{IncomeID: "i1", Description: "Salary", Amount: 5000}
	income2 := domain.IncomeRecord{IncomeID: "i2", Description: "Freelance", Amount: 1000}

	t.Run("Add and Get Income", func(t *testing.T) {
		err := repo.AddIncome(monthKey, income1)
		require.NoError(t, err)
		err = repo.AddIncome(monthKey, income2)
		require.NoError(t, err)

		incomes, err := repo.GetIncomesForMonth(monthKey)
		require.NoError(t, err)
		assert.Len(t, incomes, 2)
	})

	t.Run("Update Income", func(t *testing.T) {
		updatedIncome := domain.IncomeRecord{IncomeID: "i1", Description: "Salary", Amount: 5500}
		err := repo.UpdateIncome(monthKey, updatedIncome)
		require.NoError(t, err)

		incomes, err := repo.GetIncomesForMonth(monthKey)
		require.NoError(t, err)
		assert.Equal(t, 5500.0, incomes[0].Amount)
	})

	t.Run("Delete Income", func(t *testing.T) {
		err := repo.DeleteIncome(monthKey, "i2")
		require.NoError(t, err)

		incomes, err := repo.GetIncomesForMonth(monthKey)
		require.NoError(t, err)
		assert.Len(t, incomes, 1)
		assert.Equal(t, "i1", incomes[0].IncomeID)
	})
}

func TestJsonRepository_CategoryOperations(t *testing.T) {
	repo := setupTestRepo(t)
	monthKey := "July-2024"
	cat1 := domain.Category{CatID: "c1", GroupID: "g1", CategoryName: "Rent"}

	t.Run("Add and Get Category", func(t *testing.T) {
		err := repo.AddCategory(monthKey, cat1)
		require.NoError(t, err)

		cats, err := repo.GetCategoriesForMonth(monthKey)
		require.NoError(t, err)
		assert.Len(t, cats, 1)
		assert.Equal(t, "Rent", cats[0].CategoryName)
	})

	t.Run("Update Category", func(t *testing.T) {
		updatedCat := domain.Category{CatID: "c1", GroupID: "g1", CategoryName: "Mortgage"}
		err := repo.UpdateCategory(monthKey, updatedCat)
		require.NoError(t, err)

		cats, err := repo.GetCategoriesForMonth(monthKey)
		require.NoError(t, err)
		assert.Equal(t, "Mortgage", cats[0].CategoryName)
	})

	t.Run("Delete Category", func(t *testing.T) {
		err := repo.DeleteCategory(monthKey, "c1")
		require.NoError(t, err)

		cats, err := repo.GetCategoriesForMonth(monthKey)
		require.NoError(t, err)
		assert.Empty(t, cats)
	})
}

func TestJsonRepository_CopyFromMonth(t *testing.T) {
	repo := setupTestRepo(t)
	fromMonth := "August-2024"
	toMonth := "September-2024"
	cat1 := domain.Category{CatID: "c1", GroupID: "g1", CategoryName: "Utilities", Expense: map[string]domain.ExpenseRecord{"c1": {Amount: 100}}}
	cat2 := domain.Category{CatID: "c2", GroupID: "g1", CategoryName: "Groceries"}

	err := repo.AddCategory(fromMonth, cat1)
	require.NoError(t, err)
	err = repo.AddCategory(fromMonth, cat2)
	require.NoError(t, err)

	count, err := repo.CopyCategoriesFromMonth(fromMonth, toMonth)
	require.NoError(t, err)
	assert.Equal(t, 2, count)

	newCats, err := repo.GetCategoriesForMonth(toMonth)
	require.NoError(t, err)
	assert.Len(t, newCats, 2)
	// Verify that expenses are reset
	assert.Empty(t, newCats[0].Expense)
}

func TestJsonRepository_Persistence(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "persistent_data.json")

	// Create and modify repo1
	repo1, err := NewJsonRepository(filePath, "EUR")
	require.NoError(t, err)
	err = repo1.AddGroup(domain.CategoryGroup{GroupID: "p1", GroupName: "Persistent Group"})
	require.NoError(t, err)

	// Create repo2 from the same file
	repo2, err := NewJsonRepository(filePath, "EUR")
	require.NoError(t, err)

	// Check if repo2 loaded the data saved by repo1
	groups, err := repo2.GetAllGroups()
	require.NoError(t, err)
	assert.Len(t, groups, 1)
	assert.Equal(t, "Persistent Group", groups[0].GroupName)

	// Verify the file was actually created and has content
	fileInfo, err := os.Stat(filePath)
	require.NoError(t, err)
	assert.Greater(t, fileInfo.Size(), int64(0))
}
