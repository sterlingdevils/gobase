/*
  package chanbaedcontainer implements an ordered container with only channels as the API

  The input channel is used to add things to the container.
  The output channel will contain the head of the container when read
  The delete channel is used to remove things from the container before the are read out of
  the output channel.

  This uses Go 1.18 generics,  Things must impement the Indexable interface:
    has a method to return a comparable key

  Things come in and go out the channels in order.  Things can be removed while in the container
  by passing their key to the delete channel
*/
package chanbasedcontainer

import (
	"container/list"
	"context"
	"sync/atomic"
)

const (
	CHANSIZE = 0
)

type Keyable[K comparable] interface {
	Key() K
}

type ChanBasedContainer[K comparable, T Keyable[K]] struct {
	// We use map to hold the thing, and an ordered list of keys
	tmap  map[K]T
	tlist *list.List

	inchan  chan T
	outchan chan T
	delchan chan K

	ctx context.Context
	can context.CancelFunc

	// Holds the current thing we are trying to send
	onetosend *T

	approxSize int32
}

func (r *ChanBasedContainer[_, T]) addT(thing T) {
	k := thing.Key()
	if _, b := r.tmap[k]; b {
		return
	}

	r.tlist.PushBack(k)
	r.tmap[k] = thing
}

func (r *ChanBasedContainer[K, _]) delK(index K) {
	// If we have one to delete, check if its the one we are waiting to send
	if r.onetosend != nil {
		if (*r.onetosend).Key() == index {
			r.onetosend = nil
			return
		}
	}

	// Not onetosend so check in the continer
	for curr := r.tlist.Front(); curr != nil; curr = curr.Next() {
		val := curr.Value.(K)
		if val == index {
			r.tlist.Remove(curr)
			delete(r.tmap, index)
			break
		}
	}
}

// grab the head of the container or nil if we are empty
func (r *ChanBasedContainer[K, T]) pop() *T {
	if len(r.tmap) == 0 {
		return nil
	}

	e := r.tlist.Front()
	k := e.Value.(K)
	t := r.tmap[k]
	r.delK(k)
	return &t
}

// ApproxSize returns something close to the number of items in the container, maybe.
// Only updated at the start of each mainloop
func (r *ChanBasedContainer[_, _]) ApproxSize() int32 {
	return atomic.LoadInt32(&r.approxSize)
}

// InChan
func (r *ChanBasedContainer[_, T]) InChan() chan<- T {
	return r.inchan
}

// DelChan
func (r *ChanBasedContainer[K, _]) DelChan() chan<- K {
	return r.delchan
}

// OutChan
func (r *ChanBasedContainer[_, T]) OutChan() <-chan T {
	return r.outchan
}

// Close the ChanBasedContainer
func (r *ChanBasedContainer[_, _]) Close() {
	// Cancel our context
	r.can()
}

// RecoverFromClosedChan is used when it is OK if the channel is closed we are writing on
// This is not great using the string compare but the go runtime uses a generic error so we
// can't trap this any other way.
func recoverFromClosedChan() {
	if r := recover(); r != nil {
		if e, ok := r.(error); ok && e.Error() == "send on closed channel" {
			return
		}
		panic(r)
	}
}

// mainloop
// If the container is empty, only listen for
func (r *ChanBasedContainer[_, T]) mainloop() {
	defer close(r.outchan)
	defer recoverFromClosedChan()

	for {
		// Check if we have one ready to send
		if r.onetosend == nil {
			r.onetosend = r.pop() // pop will return nil if one is not ready
		}

		if r.onetosend == nil {
			// Save the current size
			atomic.StoreInt32(&r.approxSize, int32(len(r.tmap)))
			// None to send so don't select on output channel
			select {
			case t := <-r.inchan:
				r.addT(t)
			case k := <-r.delchan:
				r.delK(k)
			case <-r.ctx.Done():
				return
			}
		} else {
			// Save the current size
			atomic.StoreInt32(&r.approxSize, int32(len(r.tmap))+1)

			// We have one to send so select on output channel
			select {
			case r.outchan <- *r.onetosend:
				// Now that we sent it, clean onetosend so we get the next one
				r.onetosend = nil
			case t := <-r.inchan:
				r.addT(t)
			case k := <-r.delchan:
				r.delK(k)
			case <-r.ctx.Done():
				return
			}
		}
	}
}

// New returns a reference to a a container or error if there was a problem
// for performance T should be a pointer
func New[K comparable, T Keyable[K]]() (*ChanBasedContainer[K, T], error) {
	con, cancel := context.WithCancel(context.Background())
	r := ChanBasedContainer[K, T]{
		tmap:    make(map[K]T),
		tlist:   list.New(),
		inchan:  make(chan T, CHANSIZE),
		outchan: make(chan T, CHANSIZE),
		delchan: make(chan K, CHANSIZE),
		ctx:     con,
		can:     cancel}

	go r.mainloop()

	return &r, nil
}
