package ui

import (
	"errors"
	"strconv"
)

// ValidAmount validates and converts a string to a float64 amount, ensuring it's not zero.
func ValidAmount(v string) (float64, error) {

	floatValue, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, err
	}

	if floatValue == 0 {
		return 0.0, errors.New("amount cannot be zero")
	}

	return floatValue, nil
}
