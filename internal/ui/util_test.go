package ui

import (
	"testing"
	"time"
)

func TestGetPreviousMonth(t *testing.T) {
	tests := []struct {
		year      int
		month     time.Month
		wantYear  int
		wantMonth time.Month
	}{
		{2024, time.January, 2023, time.December},
		{2024, time.February, 2024, time.January},
		{2024, time.March, 2024, time.February},
		{2000, time.January, 1999, time.December},
		{2020, time.December, 2020, time.November},
	}

	for _, tt := range tests {
		gotYear, gotMonth := GetPreviousMonth(tt.year, tt.month)
		if gotYear != tt.wantYear || gotMonth != tt.wantMonth {
			t.Errorf("GetPreviousMonth(%d, %v) = (%d, %v); want (%d, %v)",
				tt.year, tt.month, gotYear, gotMonth, tt.wantYear, tt.wantMonth)
		}
	}
}

func TestGetNextMonth(t *testing.T) {
	tests := []struct {
		year      int
		month     time.Month
		wantYear  int
		wantMonth time.Month
	}{
		{2024, time.December, 2025, time.January},
		{2024, time.January, 2024, time.February},
		{2024, time.February, 2024, time.March},
		{1999, time.December, 2000, time.January},
		{2020, time.November, 2020, time.December},
	}

	for _, tt := range tests {
		gotYear, gotMonth := GetNextMonth(tt.year, tt.month)
		if gotYear != tt.wantYear || gotMonth != tt.wantMonth {
			t.Errorf("GetNextMonth(%d, %v) = (%d, %v); want (%d, %v)",
				tt.year, tt.month, gotYear, gotMonth, tt.wantYear, tt.wantMonth)
		}
	}
}
