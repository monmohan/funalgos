package main

import "fmt"

func main() {
	print(isValid("({[]})"))
}
func isValid(s string) bool {
	//loop through the string
	stack := []rune{}
	for _, char := range s {
		if char == '(' || char == '[' || char == '{' {
			stack = append(stack, char)
		} else if (char == ')' && stack[len(stack)-1] == '(') || (char == ']' && stack[len(stack)-1] == '[') || (char == '}' && stack[len(stack)-1] == '{') {
			stack = stack[:len(stack)-1]
		} else {
			return false
		}
		//print stack with characters formatted as string
		fmt.Println(string(stack))

	}
	return len(stack) == 0

}
