package main

import (
	"fmt"
)

func main() {
	print(squareRoot(6225))
}

func squareRoot(n int) int {
	start := 1
	end := n
	mid := (start + end) / 2
	iter := 20
	for iter > 0 {
		sq := mid * mid
		diff := sq - n
		fmt.Println("Start:", start, "Mid:", mid, "End:", end, "Diff:", diff)
		if diff == 0 {
			break
		}
		if diff < 0 {
			//check for small difference
			if diff > -4 {
				break
			}
			start = mid

		}
		if diff > 0 {
			//check for small difference
			if diff < 4 {
				break
			}
			end = mid

		}
		mid = (start + end) / 2
		iter--

	}

	return mid

}
