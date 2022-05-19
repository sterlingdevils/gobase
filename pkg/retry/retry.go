package retry

import (
	"context"
	"log"
	"sync"

	"github.com/sterlingdevils/gobase/pkg/obj"
)

type Retry struct {
	objin  chan *obj.Obj
	objout chan *obj.Obj
	ackin  chan uint64

	wg sync.WaitGroup

	ctx  context.Context
	can  context.CancelFunc
	once sync.Once
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

// mainloop
func (r *Retry) mainloop() {
	defer r.wg.Done()

	for {
		select {
		case o := <-r.objin:
			log.Println("Hello Obj: ", o)
			r.objout <- o
		case a := <-r.ackin:
			log.Println("Hello Ack: ", a)
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
}

// New
func New() (*Retry, error) {
	c, cancel := context.WithCancel(context.Background())
	oin := make(chan *obj.Obj, CHANSIZE)
	oout := make(chan *obj.Obj, CHANSIZE)
	ain := make(chan uint64, CHANSIZE)
	r := Retry{objin: oin, objout: oout, ackin: ain, ctx: c, can: cancel}

	r.wg.Add(1)
	go r.mainloop()

	return &r, nil
}
