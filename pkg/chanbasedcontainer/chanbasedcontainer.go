// ordered container of T indexed by Ks
package chanbasedcontainer

import (
	"container/list"
	"context"
	"errors"
	"sync"
)

const (
	CHANSIZE = 10
)

type Indexable[K comparable] interface {
	Key() K
}

type ChanBasedContainer[K comparable, T Indexable[K]] struct {
	objmap  map[K]T
	objlist list.List

	inchan  chan T
	outchan chan T
	delchan chan K

	ctx  context.Context
	can  context.CancelFunc
	once sync.Once
}

//func (r *ChanBasedContainer[K, T])

// Add will place key onto the in channel
func (r *ChanBasedContainer[K, T]) Add(thing T) {
	r.inchan <- thing
}

// next will return the first item in the container
// this will block until one is ready
func (r *ChanBasedContainer[K, T]) Next() (T, error) {
	select {
	case thing := <-r.outchan:
		return thing, nil
	case <-r.ctx.Done():
		return *new(T), errors.New("container is closed")
	}
}

// Delete will place key onto the delete channel
func (r *ChanBasedContainer[K, T]) Delete(key K) {
	r.delchan <- key
}

// InChan
func (r *ChanBasedContainer[K, T]) InChan() chan<- T {
	return r.inchan
}

// OutChan
func (r *ChanBasedContainer[K, T]) OutChan() <-chan T {
	return r.outchan
}

// DelChan
func (r *ChanBasedContainer[K, T]) DelChan() chan<- K {
	return r.delchan
}

// Close the ChanBasedContainer
func (r *ChanBasedContainer[K, T]) Close() {
	r.can()
	r.once.Do(func() {
		close(r.outchan)
	})
}

func New[K comparable, T Indexable[K]]() (*ChanBasedContainer[K, T], error) {
	objs := make(map[K]T)
	objl := *list.New()
	ichan := make(chan T, CHANSIZE)
	ochan := make(chan T, CHANSIZE)
	dchan := make(chan K, CHANSIZE)
	c, cancel := context.WithCancel(context.Background())
	r := ChanBasedContainer[K, T]{objmap: objs, objlist: objl, inchan: ichan, outchan: ochan, delchan: dchan, ctx: c, can: cancel}
	return &r, nil
}
