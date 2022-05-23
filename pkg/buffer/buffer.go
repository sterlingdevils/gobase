package buffer

import (
	"context"
	"errors"

	"github.com/sterlingdevils/gobase/pkg/pipeline"
)

const (
	CHANSIZE = 0
)

type Buffer[T any] struct {
	ctx context.Context
	can context.CancelFunc

	inchan  chan T
	outchan chan T

	pl pipeline.Pipelineable[T]
}

// InChan
func (r *Buffer[T]) InChan() chan<- T {
	return r.inchan
}

// OutChan
func (r *Buffer[T]) OutChan() <-chan T {
	return r.outchan
}

// PipelineChan returns a R/W channel that is used for pipelining
func (r *Buffer[T]) PipelineChan() chan T {
	return r.outchan
}

// Close
func (r *Buffer[_]) Close() {
	// If we pipelined then call Close the input pipeline
	if r.pl != nil {
		r.pl.Close()
	}

	// Cancel our context
	r.can()
}

// mainloop, read from in channel and write to out channel safely
// exit when our context is closed
func (r *Buffer[_]) mainloop() {
	defer close(r.outchan)

	for {
		select {
		case t := <-r.inchan:
			select {
			case r.outchan <- t:
			case <-r.ctx.Done():
				return
			}
		case <-r.ctx.Done():
			return
		}
	}
}

func NewWithChannel[T any](size int, in chan T) (*Buffer[T], error) {
	con, cancel := context.WithCancel(context.Background())

	r := Buffer[T]{
		ctx:     con,
		can:     cancel,
		inchan:  in,
		outchan: make(chan T, size)}

	go r.mainloop()

	return &r, nil
}

func NewWithPipeline[T any](size int, p pipeline.Pipelineable[T]) (*Buffer[T], error) {
	if p == nil {
		return nil, errors.New("bad pipeline passed in to New")
	}

	r, err := NewWithChannel(size, p.PipelineChan())
	if err != nil {
		return nil, err
	}

	r.pl = p

	return r, nil
}

func New[T any](size int) (*Buffer[T], error) {
	con, cancel := context.WithCancel(context.Background())

	r := Buffer[T]{
		ctx:     con,
		can:     cancel,
		inchan:  make(chan T, size),
		outchan: make(chan T, CHANSIZE)}

	go r.mainloop()

	return &r, nil
}
