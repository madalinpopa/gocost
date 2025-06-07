package ui

import (
	"errors"
	"strconv"
)

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
