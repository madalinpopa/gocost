package ui

import (
	"github.com/shopspring/decimal"
)

func ValidAmount(v string) (decimal.Decimal, error) {
	d, err := decimal.NewFromString(v)
	if err != nil {
		return decimal.Decimal{}, err
	}
	return d, nil
}
