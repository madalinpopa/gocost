package ui

import (
	"testing"
)

func TestValidAmount(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      float64
		expectErr bool
	}{
		{name: "valid positive float", input: "123.45", want: 123.45, expectErr: false},
		{name: "valid negative float", input: "-789.01", want: -789.01, expectErr: false},
		{name: "valid positive integer as float", input: "100", want: 100.0, expectErr: false},
		{name: "zero value", input: "0", want: 0.0, expectErr: true},
		{name: "empty input", input: "", want: 0.0, expectErr: true},
		{name: "non-numeric input", input: "abc", want: 0.0, expectErr: true},
		{name: "input with spaces", input: " 42.0 ", want: 42.0, expectErr: false},
		{name: "input with special characters", input: "$123", want: 0.0, expectErr: true},
		{name: "large valid float", input: "123456789.123456", want: 123456789.123456, expectErr: false},
		{name: "valid scientific notation", input: "1e6", want: 1000000.0, expectErr: false},
		{name: "invalid scientific notation", input: "5e", want: 0.0, expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidAmount(tt.input)
			if (err != nil) != tt.expectErr {
				t.Fatalf("ValidAmount(%q) error = %v, expectErr %v", tt.input, err, tt.expectErr)
			}
			if !tt.expectErr && got != tt.want {
				t.Errorf("ValidAmount(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
