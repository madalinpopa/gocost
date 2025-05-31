package data

import (
	"encoding/json"
	"errors"
	"os"
)

// LoadData loads data from a file.
func LoadData(filePath string) (*DataRoot, error) {

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

	return &dataRoot, nil
}

func SaveData(filePath string, data *DataRoot) error {
	return nil
}
