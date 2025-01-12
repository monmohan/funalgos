package main

import (
	"fmt"
	"time"
)

func main() {
	//find time taken for each method
	t1 := time.Now()
	res := waysToClimb(40)
	fmt.Println("Number of ways to climb 40 steps: ", res)
	fmt.Println("Time taken for recursive method: ", time.Since(t1))
	t1 = time.Now()
	res = waysToClimb2(40)
	fmt.Println("Number of ways to climb 40 steps: ", res)
	fmt.Println("Time taken for iterative method: ", time.Since(t1))
	//comment the below section
	//output:
	/*$ go run staircase.go
	  Number of ways to climb 50 steps:  23837527729
	  Time taken for recursive method:  27.912315s
	  Number of ways to climb 50 steps:  23837527729
	  Time taken for iterative method:  1.334µs
	*/

}

func waysToClimb(steps int) int {
	//fmt.Println(steps)
	if steps <= 2 {
		return []int{1, 1, 2}[steps]
	}

	return waysToClimb(steps-1) + waysToClimb(steps-2) + waysToClimb(steps-3)

}

func waysToClimb2(steps int) int {
	if steps <= 2 {
		return []int{1, 1, 2}[steps]
	}

	ways := make([]int, 3)
	ways[0] = 1
	ways[1] = 1
	ways[2] = 2

	for i := 3; i <= steps; i++ {
		n := ways[0] + ways[1] + ways[2]
		ways = append(ways, n)
		//remove first
		ways = ways[1:]
		//fmt.Println(ways)
	}

	return ways[2]
}
