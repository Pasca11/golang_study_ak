package main

import (
	"fmt"
	"runtime"
	"time"
)

func factorialRecursive(n int) int {
	if n == 0 {
		return 1
	}
	return n * factorialRecursive(n-1)
}
func factorialIterative(n int) int {
	res := 1
	for i := 1; i <= n; i++ {
		res *= i
	}
	return res
}

// выдает true, если реализация быстрее и false, если медленнее
func compareWhichFactorialIsFaster() string {
	testN := 100000
	start := time.Now()
	factorialRecursive(testN)
	res1 := time.Since(start)

	start = time.Now()
	factorialIterative(testN)
	res2 := time.Since(start)

	if res1 > res2 {
		return "Iterative is faster"
	}
	return "Recursive is faster"
}
func main() {
	fmt.Println("Go version:", runtime.Version())
	fmt.Println("Go OS/Arch:", runtime.GOOS, "/", runtime.GOARCH)
	fmt.Println("Which factorial is faster?")
	fmt.Println(compareWhichFactorialIsFaster())
}
