package main

import "testing"

func TestFactorial(t *testing.T) {
	res := Factorial(0)
	if res != 1 {
		t.Errorf("Factorial(0) failed.")
	}
	res = Factorial(1)
	if res != 1 {
		t.Errorf("Factorial(1) failed.")
	}
	res = Factorial(5)
	if res != 120 {
		t.Errorf("Factorial(5) failed.")
	}
}
