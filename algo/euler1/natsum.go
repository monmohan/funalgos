package main

func main() {
	print(natSum(100))
}
func natSum(n int) int {
	limit3 := (n - 1) / 3
	limit5 := (n - 1) / 5
	limit15 := (n - 1) / 15
	sum3 := 3 * limit3 * (limit3 + 1) / 2
	sum5 := 5 * limit5 * (limit5 + 1) / 2
	sum15 := 15 * limit15 * (limit15 + 1) / 2
	return sum3 + sum5 - sum15
}
