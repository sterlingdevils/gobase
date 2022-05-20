package chanbasedcontainer_test

import (
	"fmt"

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

func ExampleChanBasedContainer_fulltest() {
	r, _ := chanbasedcontainer.New[uint64, *Node]()
	r.InChan() <- &Node{key: 1, data: []byte("I don't care what it is")}
	r.InChan() <- &Node{key: 2, data: []byte("This is a test")}
	n := <-r.OutChan()
	fmt.Printf("%v, %v", n.key, n.data)
	r.Close()
	// Output:
	// 2, [84 104 105 115 32 105 115 32 97 32 116 101 115 116]
}
