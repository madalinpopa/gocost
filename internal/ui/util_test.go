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

func TestIsViewingCurrentMonth(t *testing.T) {
	now := time.Now()
	
	// Get a different month to ensure test reliability
	differentMonth := time.January
	if now.Month() == time.January {
		differentMonth = time.February
	}
	
	tests := []struct {
		name      string
		month     time.Month
		year      int
		want      bool
	}{
		{
			name:  "current month and year",
			month: now.Month(),
			year:  now.Year(),
			want:  true,
		},
		{
			name:  "different month same year",
			month: differentMonth,
			year:  now.Year(),
			want:  false,
		},
		{
			name:  "same month different year",
			month: now.Month(),
			year:  now.Year() - 1,
			want:  false,
		},
		{
			name:  "different month and year",
			month: differentMonth,
			year:  2020,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := MonthlyModel{
				MonthYear: MonthYear{
					CurrentMonth: tt.month,
					CurrentYear:  tt.year,
				},
			}
			
			got := model.isViewingCurrentMonth()
			if got != tt.want {
				t.Errorf("isViewingCurrentMonth() = %v; want %v", got, tt.want)
			}
		})
	}
}
