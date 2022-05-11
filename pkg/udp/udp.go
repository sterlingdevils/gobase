package udp

import (
	"context"
	"log"
	"net"
	"sdbase/pkg/chantools"
	"sync"
)

type chantype []byte

type ConnType int

const (
	maxDSz          = 65507
	SERVER ConnType = 1
	CLIENT ConnType = 2
)

type UDP struct {
	addr string
	in   <-chan chantype
	out  chan chantype

	conn *net.UDPConn

	ctx  context.Context
	can  context.CancelFunc
	once sync.Once
}

func (u *UDP) protectChanWrite(t chantype) {
	defer chantools.RecoverFromClosedChan()
	select {
	case u.out <- t:
	case <-u.ctx.Done():
	}
}

func (u *UDP) serverConn() error {
	addr, err := net.ResolveUDPAddr("udp4", u.addr)
	if err != nil {
		return err
	}

	u.conn, err = net.ListenUDP("udp4", addr)
	if err != nil {
		return err
	}

	return nil
}

func (u *UDP) clientConn() error {
	a, err := net.ResolveUDPAddr("udp4", u.addr)
	if err != nil {
		return err
	}

	u.conn, err = net.DialUDP("udp4", nil, a)
	if err != nil {
		return err
	}

	return nil
}

// listen loop to incomming UDP
func (u *UDP) processInUDP(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		buf := make([]byte, maxDSz)
		n, _, err := u.conn.ReadFromUDP(buf)
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

// Process Incomming Channel Data
func (u *UDP) processInChan(wg *sync.WaitGroup) {
	defer wg.Done()

	for b := range u.in {
		_, err := u.conn.Write(b)
		if err != nil {
			log.Println("udp write failed")
		}
	}
}

// mainloop
func (u *UDP) mainloop(wg *sync.WaitGroup) {
	defer wg.Done()

	wg.Add(1)
	go u.processInUDP(wg)

	wg.Add(1)
	go u.processInChan(wg)

	for {
		select {
		case <-u.ctx.Done():
		}
	}
}

// ------- Public Methods --------

// OutputChan returns read only output channel
func (u *UDP) OuputChan() <-chan chantype {
	return u.out
}

func (u *UDP) Close() {
	u.can()
	u.once.Do(func() {
		close(u.out)
	})
}

func New(wg *sync.WaitGroup, in1 <-chan chantype, addr string, ct ConnType) (*UDP, error) {
	c, cancel := context.WithCancel(context.Background())
	udp := UDP{out: make(chan chantype), addr: addr, ctx: c, can: cancel, in: in1}

	switch ct {
	case SERVER:
		udp.serverConn()
	case CLIENT:
		udp.clientConn()
	}

	wg.Add(1)
	go udp.mainloop(wg)

	return &udp, nil
}
