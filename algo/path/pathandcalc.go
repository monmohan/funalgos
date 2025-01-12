package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	//pretty print the output
	fmt.Println(path("/a/./b/../../c/d"))
	fmt.Println(path("/home/"))
	fmt.Println(path("/../"))
	fmt.Println(evalExpression("((((7*3)+5)-3)*6)") == 138)
	fmt.Println(evalExpression("(((2*3))+1)") == 7)

}

func path(p string) string {
	stack := []string{}
	parts := strings.Split(p, "/")
	//print parts and length
	fmt.Println(parts, len(parts))
	for _, v := range parts {
		//fmt.Println(v)
		if v == ".." {
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		} else if v == "." || v == "" {
			continue
		} else {
			stack = append(stack, v)
		}
	}

	return "/" + strings.Join(stack, "/")
}

func evalExpression(exp string) int {
	parts := []string{}
	buf := []byte{}
	for _, v := range exp {
		{
			if unicode.Is(unicode.White_Space, v) {
				if len(buf) > 0 {
					parts = append(parts, string(buf))
					buf = []byte{}
				}
				continue
			}
			if v == '(' || v == ')' || v == '+' || v == '-' || v == '*' || v == '/' {
				if len(buf) > 0 {
					parts = append(parts, string(buf))
					buf = []byte{}
				}
				parts = append(parts, string(v))
			} else {
				buf = append(buf, byte(v))
			}
		}
	}
	if len(buf) > 0 {
		parts = append(parts, string(buf))
	}
	fmt.Println(len(parts), parts)
	stack := []string{}

	for _, v := range parts {
		if v == "(" {
			stack = append(stack, v)
		} else if v == ")" {
			fmt.Println(v, stack)
			for i := len(stack) - 1; i >= 0; i-- {
				if stack[i] == "(" {
					fmt.Println("i:", i)
					v := eval(stack[i+1:])
					stack = stack[:i]
					stack = append(stack, strconv.Itoa(v))
					break
				}
			}

		} else {
			stack = append(stack, v)
		}

	}
	//print with message
	fmt.Println("Stack:", stack)
	val := ""
	if len(stack) == 1 {
		val = stack[0]
	}

	result, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("Error:", err)
		return 0

	}
	return result

}

func eval(parts []string) int {
	if len(parts) == 1 {
		result, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Println("Error:", err)
			return 0
		}
		return result
	}

	opand1, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}
	opand2, err := strconv.Atoi(parts[2])
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}
	if parts[1] == "/" {
		return opand1 / opand2
	}
	if parts[1] == "+" {
		return opand1 + opand2
	}
	if parts[1] == "-" {
		return opand1 - opand2
	}
	if parts[1] == "*" {
		return opand1 * opand2
	}
	return 0
}
