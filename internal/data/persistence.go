package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// LoadData loads data from a file.
func LoadData(filePath string, currency string) (*DataRoot, error) {

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return NewDataRoot(), nil
		}
		return nil, err
	}

	if len(fileData) == 0 {
		return NewDataRoot(), nil
	}

	var dataRoot DataRoot
	err = json.Unmarshal(fileData, &dataRoot)
	if err != nil {
		return nil, err
	}

	if dataRoot.CategoryGroups == nil {
		dataRoot.CategoryGroups = make([]CategoryGroup, 0)
	}

	if dataRoot.MonthlyData == nil {
		dataRoot.MonthlyData = make(map[string]MonthlyRecord, 0)
	}

	dataRoot.DefaultCurrency = currency

	return &dataRoot, nil
}

// SaveData saves data to a file.
func SaveData(filePath string, data *DataRoot) error {

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}
	return nil
}
