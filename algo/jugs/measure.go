package main

import (
	"fmt"
)

type jug struct {
	capacity int
	current  int
}

type Node struct {
	jugs []jug
}

func (n *Node) fillJug(index int) *Node {
	if n.jugs[index].current == n.jugs[index].capacity {
		return n
	}
	newNode := &Node{jugs: make([]jug, len(n.jugs))}
	copy(newNode.jugs, n.jugs)
	newNode.jugs[index].current = newNode.jugs[index].capacity
	return newNode
}
func (n *Node) emptyJug(index int) *Node {
	if n.jugs[index].current == 0 {
		return n
	}
	newNode := &Node{jugs: make([]jug, len(n.jugs))}
	copy(newNode.jugs, n.jugs)
	newNode.jugs[index].current = 0
	return newNode
}

func (n *Node) transferJug(from, to int) *Node {
	if n.jugs[to].current == n.jugs[to].capacity {
		return n
	}
	if n.jugs[from].current == 0 {
		return n
	}
	newNode := &Node{jugs: make([]jug, len(n.jugs))}
	copy(newNode.jugs, n.jugs)
	a := n.jugs[from].current + n.jugs[to].current
	if a > n.jugs[to].capacity {
		newNode.jugs[to].current = newNode.jugs[to].capacity
		newNode.jugs[from].current = a - newNode.jugs[to].capacity
	} else {
		newNode.jugs[to].current = a
		newNode.jugs[from].current = 0
	}

	return newNode
}

func main() {
	/*j1 := jug{capacity: 5}
	j2 := jug{capacity: 3}
	n := &Node{jugs: []jug{j1, j2}}
	n1 := n.fillJug(0)
	n2 := n1.transferJug(0, 1)

	fmt.Println(n, n1, n2)*/
	start := &Node{jugs: []jug{{5, 0}, {3, 0}}}
	findPath(start, 4)

}

func findPath(start *Node, capacity int) *Node {
	q := []*Node{start}
	visited := make(map[*Node]bool)
	printQ := 2
	path := make(map[*Node]*Node)
	for len(q) > 0 {

		n := q[0]
		q = q[1:]
		visited[n] = true
		for i := 0; i < len(n.jugs); i++ {
			if n.jugs[i].current == capacity {
				//print path
				for n != nil {
					fmt.Println(n)
					n = path[n]
				}
				return n
			}
		}
		//find children
		//all fill children
		for i := 0; i < len(n.jugs); i++ {
			c1 := n.fillJug(i)
			if _, ok := visited[c1]; !ok {
				q = append(q, c1)
				path[c1] = n
				visited[c1] = true
			}
		}

		//all empty children
		for i := 0; i < len(n.jugs); i++ {
			c1 := n.emptyJug(i)
			if _, ok := visited[c1]; !ok {
				q = append(q, c1)
				path[c1] = n
				visited[c1] = true
			}
		}
		//all transfer children
		for i := 0; i < len(n.jugs); i++ {
			for j := 0; j < len(n.jugs); j++ {
				if i != j {
					c1 := n.transferJug(i, j)
					if _, ok := visited[c1]; !ok {
						q = append(q, c1)
						path[c1] = n
						visited[c1] = true
					}
				}
			}
		}
		if printQ != 0 {
			//print contents of q
			for _, n1 := range q {
				fmt.Print(n1, ",")
			}
			printQ--
			fmt.Println()
		}

	}
	return nil

}
