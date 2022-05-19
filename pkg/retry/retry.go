package retry

import (
	"context"
	"sync"

	"github.com/sterlingdevils/gobase/pkg/chantools"
	"github.com/sterlingdevils/gobase/pkg/obj"
	"github.com/sterlingdevils/gobase/retrycontainer"
)

type Retry struct {
	objin  chan *obj.Obj
	objout chan *obj.Obj
	ackin  chan uint64

	wg sync.WaitGroup

	ctx  context.Context
	can  context.CancelFunc
	once sync.Once

	retrycontainer *retrycontainer.RetryContainer[uint64, *obj.Obj]
}

const (
	CHANSIZE = 10
)

// ObjIn
func (r *Retry) ObjIn() chan<- *obj.Obj {
	return r.objin
}

// ObjOut
func (r *Retry) ObjOut() <-chan *obj.Obj {
	return r.objout
}

// AckIn
func (r *Retry) AckIn() chan<- uint64 {
	return r.ackin
}

// chcecksendout do a safe write to the output channel
func (r *Retry) checksendout(o *obj.Obj) {
	defer chantools.RecoverFromClosedChan()

	// Check if context expired, if so just drop it
	if o.Ctx.Err() != nil {
		return
	}

	// Send to output channel
	select {
	case <-o.Ctx.Done():
	case r.objout <- o:
	case <-r.ctx.Done():
		return
	}

	// Send to retry channel
	select {
	case <-o.Ctx.Done():
	case r.retrycontainer.InChan() <- o:
	case <-r.ctx.Done():
		return
	}
}

// mainloop
func (r *Retry) mainloop() {
	defer r.wg.Done()

	for {
		select {
		case o := <-r.objin:
			r.checksendout(o)
		case a := <-r.ackin:
			r.retrycontainer.DelChan() <- a
		case o := <-r.retrycontainer.OutChan():
			_ = o
		case <-r.ctx.Done():
			return
		}
	}
}

// Close us
func (r *Retry) Close() {
	r.can()
	r.once.Do(func() {
		close(r.objout)
	})

	// close the retry container
	r.retrycontainer.Close()
}

// New
func New() (*Retry, error) {
	c, cancel := context.WithCancel(context.Background())
	oin := make(chan *obj.Obj, CHANSIZE)
	oout := make(chan *obj.Obj, CHANSIZE)
	ain := make(chan uint64, CHANSIZE)

	r := Retry{objin: oin, objout: oout, ackin: ain, ctx: c, can: cancel}

	// Create a retry container
	var err error
	r.retrycontainer, err = retrycontainer.New[uint64, *obj.Obj]()
	if err != nil {
		return nil, err
	}

	r.wg.Add(1)
	go r.mainloop()

	return &r, nil
}
