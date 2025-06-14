package data

import (
	"github.com/madalinpopa/gocost/internal/domain"
)

// DataRoot represents the root data structure, specific to the JSON store.
type DataRoot struct {
	DefaultCurrency string                          `json:"defaultCurrency"`
	CategoryGroups  map[string]domain.CategoryGroup `json:"CategoryGroups"`
	MonthlyData     map[string]domain.MonthlyRecord `json:"monthlyData"`
}

// NewDataRoot creates a new instance of DataRoot
func NewDataRoot() *DataRoot {
	return &DataRoot{
		CategoryGroups: make(map[string]domain.CategoryGroup, 0),
		MonthlyData:    make(map[string]domain.MonthlyRecord, 0),
	}
}
