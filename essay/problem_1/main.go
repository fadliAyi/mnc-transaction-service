package main

import (
	"fmt"
	"strings"
)

func findMatchingStrings(n int, stringsList []string) interface{} {
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			if strings.ToLower(stringsList[i]) == strings.ToLower(stringsList[j]) {
				return []int{i + 1, j + 1} // Return 1-based index
			}
		}
	}
	return false
}

func main() {
	// Example 1
	n := 4
	stringsList := []string{"abcd", "acbd", "aaab", "acbd"}
	result := findMatchingStrings(n, stringsList)
	fmt.Println(result) // Output: [2 4]

	// Example 2
	n = 11
	stringsList = []string{"Satu", "Sate", "Tujuh", "Tusuk", "Tujuh", "Sate", "Bonus", "Tiga", "Puluh", "Tujuh", "Tusuk"}
	result = findMatchingStrings(n, stringsList)
	fmt.Println(result) // Output: [3 5 10]

	// Example 3
	n = 5
	stringsList = []string{"pisang", "goreng", "enak", "sekali", "rasanya"}
	result = findMatchingStrings(n, stringsList)
	fmt.Println(result) // Output: false
}
