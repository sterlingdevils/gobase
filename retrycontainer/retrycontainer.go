// ordered container of T indexed by Ks
package retrycontainer

import (
	"container/list"
	"context"
	"sync"
)

const (
	CHANSIZE = 10
)

type RetryContainer[K comparable, T any] struct {
	objmap  map[K]T
	objlist list.List

	inchan  chan T
	outchan chan T
	delchan chan K

	ctx  context.Context
	can  context.CancelFunc
	once sync.Once
}

// InChan
func (r *RetryContainer[K, T]) InChan() chan<- T {
	return r.inchan
}

// OutChan
func (r *RetryContainer[K, T]) OutChan() <-chan T {
	return r.outchan
}

// DelChan
func (r *RetryContainer[K, T]) DelChan() chan<- K {
	return r.delchan
}

// Close the RetryContainer
func (r *RetryContainer[K, T]) Close() {
	r.can()
	r.once.Do(func() {
		close(r.outchan)
	})
}

func New[K comparable, T any]() (*RetryContainer[K, T], error) {
	objs := make(map[K]T)
	objl := *list.New()
	ichan := make(chan T, CHANSIZE)
	ochan := make(chan T, CHANSIZE)
	dchan := make(chan K, CHANSIZE)
	c, cancel := context.WithCancel(context.Background())
	r := RetryContainer[K, T]{objmap: objs, objlist: objl, inchan: ichan, outchan: ochan, delchan: dchan, ctx: c, can: cancel}
	return &r, nil
}
