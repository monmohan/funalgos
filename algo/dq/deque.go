package main

type node struct {
	v    interface{}
	prev *node
	next *node
}

type deque struct {
	head *node
	tail *node
}

func (dq *deque) push(v interface{}) {
	node := &node{v, dq.tail, nil}
	dq.tail = node
}

func (dq *deque) undo() interface{} {
	node := dq.tail
	dq.tail = dq.tail.prev
	dq.tail.next = nil
	node.prev = nil
	node.next = nil
	return node.v
}
