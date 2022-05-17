// Package udp implements a UDP socket component that uses
// go channels to for sending and receiving UDP packets
package udp

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/sterlingdevils/gobase/pkg/chantools"
)

// Packet holds a UDP address and Data from the UDP
type Packet struct {
	// Addr holds a UDP address (with port) for the packet
	Addr *net.UDPAddr
	// Data contains the data
	Data []byte
}

const maxDSz = 65507

// ConnType are constants for UDP socket type
type ConnType int

// Socket Connection type
const (
	// SERVER used to create a listen socket
	SERVER ConnType = 1
	// CLIENT used to create a connect to socket
	CLIENT ConnType = 2
)

// UDP holds our private data for the component
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

// protectChanWrite sends to a channel with a context cancel to
// exit on contect close
func (u *UDP) protectChanWrite(t Packet) {
	defer chantools.RecoverFromClosedChan()
	select {
	case u.out <- t:
	case <-u.ctx.Done():
	}
}

// serverConn sets up the socket as a server
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

// clientConn sets up the socket as a client
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

// processInUDP will listen incomming UDP and put on output channel
//
// Notes
//   For now, due to the ReadFromUDP blocking
//   we are going to call wg.Done so things dont
//   wait for us until we get a packet.  This
//   should be a defer wg.Done()
func (u *UDP) processInUDP(wg *sync.WaitGroup) {
	wg.Done()

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

// processInChan will handle the receiving on the input channel and
// output via the UDP connection
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

// mainloop will setup to receive UDP and input channel processing
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

// ------------------------------------------------------------------------------------
// Public Methods
// ------------------------------------------------------------------------------------

// OutputChan returns read only output channel that the incomming UDP packets will
// be placed onto
func (u *UDP) OuputChan() <-chan Packet {
	return u.out
}

// Close will shutdown the output channel and cancel the context for the listen
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

	switch ct {
	case SERVER:
		if err := udp.serverConn(); err != nil {
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
