package main

import "testing"

func TestFibonacci(t *testing.T) {
	table := []struct {
		n        int
		expected int
	}{
		{5, 5},
		{6, 8},
		{7, 13},
		{8, 21},
		{9, 34},
		{10, 55},
	}
	for _, tt := range table {
		res := Fibonacci(tt.n)
		if res != tt.expected {
			t.Errorf("Fibonacci(%d), expected %d, got %d", tt.n, tt.expected, res)
		}
	}
}
