package main

func main() {
	print(breakChocolate(2, 3))
}

func breakChocolate(n, m int) int {
	if n == 0 || m == 0 {
		return 0
	}
	if n == 1 && m == 1 {
		return 0
	}
	if n == 1 && m > 1 {
		return m - 1
	}
	if m == 1 && n > 1 {
		return n - 1
	}
	return 1 + breakChocolate(n-1, m) + breakChocolate(1, m)
}
