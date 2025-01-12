package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Node struct {
	Value    string
	Children []*Node
}
type Process struct {
	Pid     string
	Ppid    string
	Command string
}

func main() {
	pmap := make(map[string][]*Process)
	cmd := exec.Command("ps", "-eo", "pid,ppid,command")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing ps: %v\n", err)

	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	scanner.Scan() // skip header

	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		comm := line[2]
		if len(line[2]) > 10 {
			comm = line[2][0:10]
		}
		//fmt.Printf("%s %s %s\n", line[0], line[1], comm)
		p := &Process{Pid: line[0], Ppid: line[1], Command: comm}
		pmap[p.Ppid] = append(pmap[p.Ppid], p)

	}

	root := mkTreeFromProcesses(pmap, "1")
	root.Print(0)
}

func mkTreeFromProcesses(pmap map[string][]*Process, pid string) *Node {
	node := &Node{Value: pid}
	for _, p := range pmap[pid] {
		node.Children = append(node.Children, mkTreeFromProcesses(pmap, p.Pid))
	}
	return node
}

func mkTree() *Node {
	root := &Node{Value: "root"}
	c3 := &Node{Value: "c1.1"}
	c4 := &Node{Value: "c1.2"}
	c5 := &Node{Value: "c1.3"}
	c1 := &Node{Value: "c1", Children: []*Node{c3, c4, c5}}
	c2 := &Node{Value: "c2"}
	c7 := &Node{Value: "c3"}
	root.Children = []*Node{c1, c2, c7}
	c8 := &Node{Value: "c3.1"}
	c7.Children = []*Node{c8}
	return root
}

func (n *Node) Print(tabs int) {
	for i := 0; i < tabs; i++ {
		fmt.Printf("  ")
	}
	fmt.Println(n.Value)
	t := tabs + 1
	for _, c := range n.Children {
		c.Print(t)
	}
}
