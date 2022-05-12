package udp

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/sterlingdevils/gobase/pkg/chantools"
)

type Packet struct {
	Addr *net.UDPAddr
	Data []byte
}

type ConnType int

const (
	maxDSz          = 65507
	SERVER ConnType = 1
	CLIENT ConnType = 2
)

type UDP struct {
	addr string
	in   <-chan Packet
	out  chan Packet

	conn *net.UDPConn

	ctx  context.Context
	can  context.CancelFunc
	once sync.Once

	ct ConnType
}

func (u *UDP) protectChanWrite(t Packet) {
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
		n, a, err := u.conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("readfromudp err: ", err)
			continue
		}

		p := Packet{Addr: a, Data: buf[:n]}
		u.protectChanWrite(p)

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
		switch u.ct {
		case SERVER:
			_, err := u.conn.WriteToUDP(b.Data, b.Addr)
			if err != nil {
				log.Println("udp write failed")
			}
		case CLIENT:
			_, err := u.conn.Write(b.Data)
			if err != nil {
				log.Println("udp write failed")
			}
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
			return
		}
	}
}

// ------- Public Methods --------

// OutputChan returns read only output channel
func (u *UDP) OuputChan() <-chan Packet {
	return u.out
}

func (u *UDP) Close() {
	u.can()
	u.once.Do(func() {
		close(u.out)
	})
}

// New will return a UDP connection component,  it can be setup with as a Server to listen
// for incomming connections, or a client to connect out to a server.  After that client and
// server mode work the same.
// Either way it will read from in channel and then send the packet, and it will listen
// for incomming packets on the socket and put them onto the output channel
func New(wg *sync.WaitGroup, in1 <-chan Packet, addr string, ct ConnType, outChanSize int) (*UDP, error) {
	c, cancel := context.WithCancel(context.Background())
	udp := UDP{out: make(chan Packet, outChanSize), addr: addr, ctx: c, can: cancel, in: in1, ct: ct}

	var err error = nil
	switch ct {
	case SERVER:
		if err = udp.serverConn(); err != nil {
			return nil, err
		}
	case CLIENT:
		if err := udp.clientConn(); err != nil {
			return nil, err
		}
	}

	wg.Add(1)
	go udp.mainloop(wg)

	return &udp, nil
}
