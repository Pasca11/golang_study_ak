package main

import "testing"

func TestFactorialIterative(t *testing.T) {
	res := factorialIterative(5)

	if res != 120 {
		t.Errorf("FactorialIterative(5)=%d; want 120", res)
	}
}

func TestFactorialRecursive(t *testing.T) {
	res := factorialRecursive(5)
	if res != 120 {
		t.Errorf("FactorialRecursive(5)=%d; want 120", res)
	}
}

func TestCompareWhichFactorialIsFaster(t *testing.T) {
	res := compareWhichFactorialIsFaster()

	if res != "Iterative is faster" {
		t.Errorf("Res is %s. Expected Iteraive is faster.", res)
	}
}
