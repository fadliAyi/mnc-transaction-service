package main

import (
	"fmt"
)

func calculateChange(total, paid int) interface{} {
	if paid < total {
		return false // If paid amount is less than total, return false
	}

	// Calculate the raw change
	change := paid - total
	// Round down the change to the nearest 100
	change = (change / 100) * 100

	// Available denominations
	denominations := []int{100000, 50000, 20000, 10000, 5000, 2000, 1000, 500, 200, 100}
	denominationCount := make(map[int]int)

	// Calculate the number of each denomination to be returned
	for _, denom := range denominations {
		if change >= denom {
			denominationCount[denom] = change / denom
			change %= denom
		}
	}

	// Format the output
	result := fmt.Sprintf("Kembalian: %d\nPecahan uang:", (paid-total)/100*100)
	for _, denom := range denominations {
		if count, exists := denominationCount[denom]; exists && count > 0 {
			result += fmt.Sprintf("\n%d lembar/koin %d", count, denom)
		}
	}

	return result
}

func main() {
	// Example 1
	total := 700649
	paid := 800000
	result := calculateChange(total, paid)
	fmt.Println(result)

	// Example 2
	total = 575650
	paid = 580000
	result = calculateChange(total, paid)
	fmt.Println(result)

	// Example 3
	total = 657650
	paid = 600000
	result = calculateChange(total, paid)
	fmt.Println(result) // Output: False
}
