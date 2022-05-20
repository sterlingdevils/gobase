package chanbasedcontainer_test

import (
	"github.com/sterlingdevils/gobase/pkg/chanbasedcontainer"
)

type Node struct {
	key  uint64
	data []byte
}

func (N Node) Key() uint64 {
	return N.key
}

func ExampleChanBasedContainer() {
	_, _ = chanbasedcontainer.New[uint64, Node]()
	// Output:
}

// ExampleChanBasedContainer_Close
func ExampleChanBasedContainer_Close() {
	r, _ := chanbasedcontainer.New[uint64, Node]()
	r.Close()
	// Output:
}

// ExampleChanBasedContainer_InChan
func ExampleChanBasedContainer_InChan() {
	node := &Node{key: 7, data: []byte("This is a test")}
	r, _ := chanbasedcontainer.New[uint64, *Node]()
	cin := r.InChan()
	cin <- node
	r.Close()
	// Output:
}

// ExampleChanBasedContainer_Add
func ExampleChanBasedContainer_Add() {
	node := &Node{key: 7, data: []byte("This is a test")}
	r, _ := chanbasedcontainer.New[uint64, *Node]()
	r.Add(node)
	r.Close()
	// Output:
}
