// ordered container of T indexed by Ks
package chanbasedcontainer

import (
	"container/list"
	"context"
	"errors"
	"sync"
)

const (
	CHANSIZE = 0
)

type Indexable[K comparable] interface {
	Key() K
}

type ChanBasedContainer[K comparable, T Indexable[K]] struct {
	objmap  map[K]T
	objlist *list.List

	inchan  chan T
	outchan chan T
	delchan chan K

	ctx  context.Context
	can  context.CancelFunc
	once sync.Once
}

//func (r *ChanBasedContainer[K, T])

func (r *ChanBasedContainer[_, T]) addT(thing T) {
	k := thing.Key()
	if _, b := r.objmap[k]; b {
		return
	}

	r.objlist.PushBack(k)
	r.objmap[k] = thing

	return
}

func (r *ChanBasedContainer[K, _]) delK(index K) {
	for curr := r.objlist.Front(); curr != nil; curr = curr.Next() {
		val := curr.Value.(K)
		if val == index {
			r.objlist.Remove(curr)
			delete(r.objmap, index)
			break
		}
	}
}

func (r *ChanBasedContainer[K, T]) pop() T {
	e := r.objlist.Front()
	k := e.Value.(K)
	t := r.objmap[k]
	r.delK(k)
	return t
}

func (r *ChanBasedContainer[_, _]) IsEmpty() bool {
	return len(r.objmap) == 0
}

func (r *ChanBasedContainer[_, _]) QueueSize() int {
	return len(r.objmap)
}

// Add will place key onto the in channel
func (r *ChanBasedContainer[_, T]) Add(thing T) {
	r.inchan <- thing
}

// next will return the first item in the container
// this will block until one is ready
func (r *ChanBasedContainer[_, T]) Next() (T, error) {
	select {
	case thing := <-r.outchan:
		return thing, nil
	case <-r.ctx.Done():
		return *new(T), errors.New("container is closed")
	}
}

// Delete will place key onto the delete channel
func (r *ChanBasedContainer[K, _]) Delete(key K) {
	r.delchan <- key
}

// InChan
func (r *ChanBasedContainer[_, T]) InChan() chan<- T {
	return r.inchan
}

// OutChan
func (r *ChanBasedContainer[_, T]) OutChan() <-chan T {
	return r.outchan
}

// DelChan
func (r *ChanBasedContainer[K, _]) DelChan() chan<- K {
	return r.delchan
}

// Close the ChanBasedContainer
func (r *ChanBasedContainer[_, _]) Close() {
	r.can()
	r.once.Do(func() {
		close(r.outchan)
	})
}

func (r *ChanBasedContainer[_, _]) mainloop() {
	for {
		if r.IsEmpty() {
			select {
			case t := <-r.inchan:
				r.addT(t)
			case k := <-r.delchan:
				r.delK(k)
			case <-r.ctx.Done():
				return
			}
		} else {
			select {
			case t := <-r.inchan:
				r.addT(t)
			case k := <-r.delchan:
				r.delK(k)
			case r.outchan <- r.pop():
			case <-r.ctx.Done():
				return
			}
		}
	}
}

func New[K comparable, T Indexable[K]]() (*ChanBasedContainer[K, T], error) {
	objs := make(map[K]T)
	objl := list.New()
	ichan := make(chan T, CHANSIZE)
	ochan := make(chan T, CHANSIZE)
	dchan := make(chan K, CHANSIZE)
	c, cancel := context.WithCancel(context.Background())
	r := ChanBasedContainer[K, T]{objmap: objs, objlist: objl, inchan: ichan, outchan: ochan, delchan: dchan, ctx: c, can: cancel}

	go r.mainloop()

	return &r, nil
}
