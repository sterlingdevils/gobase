package gobase

import (
	"context"
	"sync"
)

type Mux[T any] struct {
	in1  <-chan T
	in2  <-chan T
	out  chan T
	ctx  context.Context
	can  context.CancelFunc
	once sync.Once
}

func (m *Mux[T]) protectChanWrite(t T) {
	defer RecoverFromClosedChan()
	select {
	case m.out <- t:
	case <-m.ctx.Done():
	}
}

func (m *Mux[T]) doMux(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		// If both input channels are nil then we are done
		if m.in1 == nil && m.in2 == nil {
			return
		}
		select {
		case o, more := <-m.in1:
			if more {
				m.protectChanWrite(o)
			} else {
				m.in1 = nil // Set to nil so we skip it in next select
			}
		case o, more := <-m.in2:
			if more {
				m.protectChanWrite(o)
			} else {
				m.in2 = nil // Set to nil so we skip it in next select
			}
		}
	}
}

// OutputChan returns read only output channel
func (m *Mux[T]) OutputChan() <-chan T {
	return m.out
}

// Close will close the output channel
func (m *Mux[T]) Close() {
	m.in1 = nil
	m.in2 = nil
	m.can() // Close our context to free the write to output Chan
	m.once.Do(func() {
		close(m.out)
	})
}

// New returns a pointer to a new mux. Mux will create a new output Channel and put things
// from both input channels onto the output channel
func (Mux[T]) New(wg *sync.WaitGroup, in1, in2 <-chan T, chanSize int) *Mux[T] {
	c, cancel := context.WithCancel(context.Background())
	m := Mux[T]{in1: in1, in2: in2, out: make(chan T, chanSize), ctx: c, can: cancel}
	wg.Add(1)
	go m.doMux(wg)
	return &m
}
