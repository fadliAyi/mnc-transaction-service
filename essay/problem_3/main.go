package main

import (
	"fmt"
)

// Function to check if a given string of brackets is valid
func isValidBrackets(input string) bool {
	// Map of matching pairs
	matchingBrackets := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
		'>': '<',
	}

	// Stack to track opening brackets
	stack := []rune{}

	for _, char := range input {
		// If it's an opening bracket, push to stack
		if char == '(' || char == '[' || char == '{' || char == '<' {
			stack = append(stack, char)
		} else if char == ')' || char == ']' || char == '}' || char == '>' {
			// If it's a closing bracket, check if it matches the last opened one
			if len(stack) == 0 || stack[len(stack)-1] != matchingBrackets[char] {
				return false
			}
			// If it's a valid match, pop the last element from the stack
			stack = stack[:len(stack)-1]
		}
	}

	// If the stack is empty, all brackets were matched properly
	return len(stack) == 0
}

func main() {
	// Example 1 (True)
	input := "{{[<>[{{}}]]}}"
	result := isValidBrackets(input)
	fmt.Println(result) // Output: true

	// Example 2 (False)
	input = "[>]"
	result = isValidBrackets(input)
	fmt.Println(result) // Output: false

	// Example 3 (True)
	input = "{<{[[{{[]<{{[{[]<>}]}}<>>}}]]}>}"
	result = isValidBrackets(input)
	fmt.Println(result) // Output: true
}
