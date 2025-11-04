package main

import "testing"

func TestSum(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{"normal case", []int{1, 2, 3, 4, 5}, 15},
		{"empty slice", []int{}, 0},
		{"single element", []int{10}, 10},
		{"negative numbers", []int{-1, -2, -3}, -6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sum(tt.input)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
