package udprx

import (
	"context"
	"fmt"
	"log"
	"net"
	"sdbase/pkg/chantools"
	"sync"
)

const (
	maxDSz = 65507
)

type UDPRx struct {
	port int
	out  chan []byte
	ctx  context.Context
	can  context.CancelFunc
	once sync.Once
}

func (u *UDPRx) protectChanWrite(t []byte) {
	defer chantools.RecoverFromClosedChan()
	select {
	case u.out <- t:
	case <-u.ctx.Done():
	}
}

// listen loop
func (u *UDPRx) listen(wg *sync.WaitGroup) {
	defer wg.Done()

	address := fmt.Sprintf(":%d", u.port)
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		buf := make([]byte, maxDSz)
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("readfromudp err: ", err)
			continue
		}

		u.protectChanWrite(buf[:n])

		// Check if the context is cancled
		if u.ctx.Err() != nil {
			return
		}
	}
}

// mainloop
func (u *UDPRx) mainloop(wg *sync.WaitGroup) {
	defer wg.Done()

	wg.Add(1)
	go u.listen(wg)

	for {
		select {
		case <-u.ctx.Done():
		}
	}
}

// ------- Public Methods --------

// OutputChan returns read only output channel
func (u *UDPRx) OutputChan() <-chan []byte {
	return u.out
}

func (u *UDPRx) Close() {
	u.can()
	u.once.Do(func() {
		close(u.out)
	})
}

func New(wg *sync.WaitGroup, port int) (*UDPRx, error) {
	c, cancel := context.WithCancel(context.Background())
	udp := UDPRx{out: make(chan []byte), port: port, ctx: c, can: cancel}

	wg.Add(1)
	go udp.mainloop(wg)

	return &udp, nil
}
