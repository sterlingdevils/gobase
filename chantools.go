package gobase

// Requires GO >= 1.18 as we use generics

import (
	"fmt"
	"log"
	"sync"
)

// RecoverFromClosedChan is used when it is OK if the channel is closed we are writing on
// This is not great using the string compare but the go runtime uses a generic error so we
// can't trap this any other way.
func RecoverFromClosedChan() {
	if r := recover(); r != nil {
		if e, ok := r.(error); ok && e.Error() == "send on closed channel" {
			log.Println("might be recovery from closed channel. not a problem: ", e)
		} else {
			panic(r)
		}
	}
}

// Print out things from the in channel until the in channel is closed
func DisplayChan[T any](wg *sync.WaitGroup, in <-chan T) {
	defer wg.Done()
	for s := range in {
		fmt.Println(s)
	}
}

// AsyncSkip takes things from in channel and puts on the out channel if out channel
// can accept the write,  if not it drops the thing on the floor
// out channel will close when the in channel is closed
func AsyncSkip[T any](wg *sync.WaitGroup, in <-chan T, chanSize int) <-chan T {
	out := make(chan T, chanSize)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(out)
		for s := range in {
			select {
			case out <- s:
			default:
			}
		}
	}()
	return out
}
