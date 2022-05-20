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

func readAndPrint(num int, c <-chan *Node) {
	for i := 0; i < num; i++ {
		n := <-c
		fmt.Printf("%v, %v\n", n.key, n.data)
	}
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

func ExampleChanBasedContainer_testfirst() {
	r, _ := chanbasedcontainer.New[uint64, *Node]()
	n1 := &Node{key: 1, data: []byte("I don't care what it is")}
	n2 := &Node{key: 2, data: []byte("This is a test")}
	r.InChan() <- n1
	r.InChan() <- n2

	readAndPrint(2, r.OutChan())

	r.Close()
	// Output:
	// 1, [73 32 100 111 110 39 116 32 99 97 114 101 32 119 104 97 116 32 105 116 32 105 115]
	// 2, [84 104 105 115 32 105 115 32 97 32 116 101 115 116]
}

func ExampleChanBasedContainer_testdeloffirst() {
	r, _ := chanbasedcontainer.New[uint64, *Node]()
	r.InChan() <- &Node{key: 1, data: []byte("I don't care what it is")}
	r.InChan() <- &Node{key: 2, data: []byte("This is a test")}
	r.DelChan() <- 1

	readAndPrint(1, r.OutChan())

	r.Close()
	// Output:
	// 2, [84 104 105 115 32 105 115 32 97 32 116 101 115 116]
}

func ExampleChanBasedContainer_testdelofsecond() {
	r, _ := chanbasedcontainer.New[uint64, *Node]()
	r.InChan() <- &Node{key: 1, data: []byte("I don't care what it is")}
	r.InChan() <- &Node{key: 2, data: []byte("This is a test")}
	r.InChan() <- &Node{key: 3, data: []byte("This is a test again")}
	r.DelChan() <- 2

	readAndPrint(2, r.OutChan())

	r.Close()
	// Output:
	// 1, [73 32 100 111 110 39 116 32 99 97 114 101 32 119 104 97 116 32 105 116 32 105 115]
	// 3, [84 104 105 115 32 105 115 32 97 32 116 101 115 116 32 97 103 97 105 110]
}

func ExampleChanBasedContainer_testdelonNotThere() {
	r, _ := chanbasedcontainer.New[uint64, *Node]()
	r.InChan() <- &Node{key: 1, data: []byte("I don't care what it is")}
	r.InChan() <- &Node{key: 2, data: []byte("This is a test")}
	r.DelChan() <- 3

	readAndPrint(2, r.OutChan())

	r.Close()
	// Output:
	// 1, [73 32 100 111 110 39 116 32 99 97 114 101 32 119 104 97 116 32 105 116 32 105 115]
	// 2, [84 104 105 115 32 105 115 32 97 32 116 101 115 116]
}

func ExampleChanBasedContainer_duptest() {
	r, _ := chanbasedcontainer.New[uint64, *Node]()
	r.InChan() <- &Node{key: 1, data: []byte("I don't care what it is")}
	// This should be dropped as a dup
	r.InChan() <- &Node{key: 1, data: []byte("This is a test")}

	readAndPrint(1, r.OutChan())

	r.Close()
	// Output:
	// 1, [73 32 100 111 110 39 116 32 99 97 114 101 32 119 104 97 116 32 105 116 32 105 115]
}
