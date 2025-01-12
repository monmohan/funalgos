package main

import (
	"fmt"
	"os"
	"text/scanner"
	"time"
)

func main() {
	dict = make(map[string]bool)
	loadDict()
	/*root := createGraph("head")
	root.PrintGraph(0)
	s := &Stack{}
	fmt.Println("Path from pig to sty")

	findPath(root, "tail", s)
	for i := 0; i < len(s.elements)-1; i++ {
		fmt.Print("\nfrom ", s.elements[i].word, " --> ", s.elements[i+1].word)
	}
	fmt.Println()

	//findPathWStack(root, "sty")*/

	graph := make(map[string][]string)
	//print time taken to build graph
	start := time.Now()
	buildGraph(graph)
	fmt.Println("Time taken to build graph", time.Since(start))
	//printEdges(graph)
	findPathbfs(graph, "head", "tail")
	findPathbfs(graph, "pig", "sty")
	findPathbfs(graph, "abcd", "ten")
	findPathbfs(graph, "bear", "rise")

}

var dict map[string]bool

type Node struct {
	word     string
	Children []*Node
}

func loadDict() {
	//create scanner from file
	var s scanner.Scanner
	file := "/Users/singhmo/code/golang/src/github.com/monmohan/funalgos/algo/wordl/word-ladder/words.txt"
	//reader of file
	osFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	s.Init(osFile)
	//scan the file
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		dict[s.TokenText()] = true
	}

}

func isWordInDict(word string) bool {
	return dict[word]
}

func createGraph(from string) *Node {
	root := &Node{word: from}
	q := []*Node{root}
	seen := make(map[string]bool)

	//create all possible words with increasing level
	for len(q) > 0 {
		w := q[0]
		q = q[1:]
		seen[w.word] = true
		for i := 0; i < len(w.word); i++ {
			perms := permutate(w, i)
			for _, p := range perms {
				if _, ok := seen[p.word]; !ok {
					w.Children = append(w.Children, p)
					q = append(q, p)
					seen[p.word] = true
				}
			}

		}

	}
	return root
}
func permutate(w *Node, index int) []*Node {
	var result []*Node
	for i := 0; i <= 25; i++ {
		a := getAlphabetLetter(i)[0]
		if w.word[index] != a {
			t := w.word[:index] + string(a) + w.word[index+1:]
			//fmt.Println("Permutating", w.word, "to", t)
			if isWordInDict(t) {
				result = append(result, &Node{word: t})
			}
		}
	}
	return result
}

type Stack struct {
	elements []*Node
}

func (s *Stack) push(n *Node) {
	s.elements = append(s.elements, n)
}
func (s *Stack) pop() *Node {
	if len(s.elements) == 0 {
		return nil
	}
	n := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return n
}

func findPath(from *Node, to string, s *Stack) bool {
	s.push(from)
	if from.word == to {
		return true
	}
	for _, c := range from.Children {
		if findPath(c, to, s) {
			return true
		}
	}
	s.pop()
	return false

}

func findPathWStack(from *Node, to string) bool {

	if from.word == to {
		return true
	}
	for _, c := range from.Children {
		if findPathWStack(c, to) {
			fmt.Println(c.word)
			return true
		}
	}

	return false

}

func getAlphabetLetter(index int) string {
	// Check if index is valid (0-25 for A-Z)
	if index < 0 || index > 25 {
		return "Invalid index"
	}

	// Convert index to corresponding ASCII character
	// 'A' starts at ASCII 65
	return string(rune('a' + index))
}

func (n *Node) PrintGraph(tabs int) {
	for i := 0; i < tabs; i++ {
		fmt.Printf("  ")
	}
	fmt.Println(n.word)
	t := tabs + 1
	for _, c := range n.Children {
		c.PrintGraph(t)
	}
}

func buildGraph(graph map[string][]string) {

	processed := map[string]bool{}
	queued := map[string]bool{}
	q := []string{}

	for word := range dict {

		//fmt.Println("Processing", word)
		if processed[word] {
			continue
		}

		//fmt.Println("Processing ", word)
		q = append(q, word)
		//find all words which are one letter away
		//if the word is already processed then skip
		//add the word and its one letter away words to the graph
		for len(q) > 0 {

			word := q[0]
			q = q[1:]

			neighbours := findNeighbours(word)
			//fmt.Println("Neighbours of", word, "are", neighbours)
			for _, n := range neighbours {
				if processed[n] {
					continue
				}
				graph[word] = append(graph[word], n)
				graph[n] = append(graph[n], word)
				//fmt.Println("Adding neighbour to queue ", n)
				if !queued[n] {
					q = append(q, n)
					queued[n] = true
				}

			}
			processed[word] = true
			//fmt.Println("Q", q)
			//fmt.Println("Processed", processed)
		}

	}

}

func findNeighbours(word string) []string {
	neighbours := []string{}

	for w := range dict {
		if (len(w) != len(word)) || (word == w) {
			continue
		}

		diffC := 0
		for i := 0; i < len(word); i++ {
			if word[i] != w[i] {
				diffC++
			}
		}
		if diffC > 1 {
			continue
		}
		neighbours = append(neighbours, w)
	}
	return neighbours

}

func printEdges(graph map[string][]string) {
	for k, v := range graph {
		fmt.Println("Edges for", k, "are", v)
	}
}

func findPathbfs(graph map[string][]string, from string, to string) {
	visited := make(map[string]bool)
	parent := map[string]string{}
	q := []string{}
	q = append(q, from)

	for len(q) > 0 {
		node := q[0]
		q = q[1:]
		if node == to {
			fmt.Println("Path found from", from, "to", to)

			p := to
			for p != from {
				fmt.Println(p)
				p = parent[p]
			}
			fmt.Println(from)

			return
		}
		for _, n := range graph[node] {
			if _, ok := visited[n]; !ok {
				parent[n] = node
				visited[n] = true
				q = append(q, n)
			}
		}

	}
	fmt.Println("No path found from", from, "to", to)

}
