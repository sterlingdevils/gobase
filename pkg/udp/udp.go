/*
  Package udp implements a UDP socket component that uses
  go channels for sending and receiving UDP packets.

  The main interface is two channels, one input channel
  and one output channel.  Any Packet put onto the input channel
  will be sent out a UDP network socket. Any received UDP packets
  from the network socket will be placed onto the output channel.

  The input channel is passed into this component so the caller can control the
  life time of the channel.  It should be closed to cause the input channel
  processing routine to finish.

  For the SERVER mode, the component will listen on the passed in port, the
  underling socket does not contain a destination address so it needs to be
  set in the Packet that is put onto the input channel.

  For the CLIENT mode, the component will Dial the address passed in, the
  underling socket contains that address and the send will ignore any address
  set in the Packet.  All outgoing Packets will be sent the the address passed
  in the New call.

  New(port) will handle creating the waitgroup and input channel
  NewwithParams(...) can be give the caller more options
*/
package udp

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/sterlingdevils/gobase/pkg/chantools"
)

// Packet holds a UDP address and Data from the UDP
// The channel types for the input and output channel are of this type
type Packet struct {
	// Addr holds a UDP address (with port) for the packet
	// will be ignored if UDP is created in CLIENT mode
	Addr *net.UDPAddr
	// Data contains the data
	Data []byte
}

// udp constants for the protocol
const (
	// Max size of a Packet Data slice
	MaxPacketSize = 65507
)

// ConnType are constants for UDP socket type
type ConnType int

// Socket Connection type
const (
	// SERVER is used to create a listen socket
	SERVER = ConnType(1)
	// CLIENT is used to create a connected socket (using Dial)
	CLIENT = ConnType(2)
)

// UDP holds our private data for the component
type UDP struct {
	addr string
	in   chan Packet
	out  chan Packet

	conn *net.UDPConn

	ctx  context.Context
	can  context.CancelFunc
	once sync.Once

	ct ConnType
}

// protectChanWrite sends to a channel with a context cancel to
// exit on contect close even if the write to channel is blocked
func (u *UDP) protectChanWrite(t Packet) {
	defer chantools.RecoverFromClosedChan()
	select {
	case u.out <- t:
	case <-u.ctx.Done():
	}
}

// startConn sets up the socket as a server
func (u *UDP) startConn() error {
	addr, err := net.ResolveUDPAddr("udp4", u.addr)
	if err != nil {
		return err
	}

	switch u.ct {
	case SERVER:
		u.conn, err = net.ListenUDP("udp4", addr)
		if err != nil {
			return err
		}
	case CLIENT:
		u.conn, err = net.DialUDP("udp4", nil, addr)
		if err != nil {
			return err
		}
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
	defer wg.Done()

	for {
		// Check if the context is cancled
		if u.ctx.Err() != nil {
			return
		}

		buf := make([]byte, MaxPacketSize)
		u.conn.SetReadDeadline(time.Now().Add(2 * time.Second))

		n, a, err := u.conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		p := Packet{Addr: a, Data: buf[:n]}
		u.protectChanWrite(p)
	}
}

// processInChan will handle the receiving on the input channel and
// output via the UDP connection
func (u *UDP) processInChan(wg *sync.WaitGroup) {
	defer wg.Done()

	send := func(p Packet) {
		if len(p.Data) > MaxPacketSize {
			log.Printf("packet size exceeds max: %v\n", len(p.Data))
			return
		}
		switch u.ct {
		case SERVER:
			_, err := u.conn.WriteToUDP(p.Data, p.Addr)
			if err != nil {
				log.Println("udp write failed")
			}
		case CLIENT:
			_, err := u.conn.Write(p.Data)
			if err != nil {
				log.Println("udp write failed")
			}
		}
	}

	// wait for packets on the input channel or the context to close
	for {
		select {
		case b, more := <-u.in:
			if !more { // if the channel is closed, then we are done
				return
			}
			send(b)
		case <-u.ctx.Done():
			return
		}
	}
}

// ------------------------------------------------------------------------------------
// Public Methods
// ------------------------------------------------------------------------------------

// InChan returns a write only channel that the incomming packets will be read from
func (u *UDP) InChan() chan<- Packet {
	return u.in
}

// OutputChan returns a read only output channel that the incomming UDP packets will
// be placed onto
func (u *UDP) OuputChan() <-chan Packet {
	return u.out
}

// Close will shutdown the output channel and cancel the context for the listen
func (u *UDP) Close() {
	u.can()
	u.once.Do(func() {
		u.conn.Close()
		close(u.out)
	})
}

// ------------------------------------------------------------------------------------
// New Functions to create a UDP
// ------------------------------------------------------------------------------------

// NewWithParams will return a UDP connection component,  it can be setup with as a Server to listen
// for incomming connections, or a client to connect out to a server.  After that client and
// server mode work the same.
// Either way it will read from in channel and then send the packet, and it will listen
// for incomming packets on the socket and put them onto the output channel
//
// This code uses the waitgoup and will add 1 for each routine it starts.  The Close method
// needs to be called so we stop all our routines
//
//  NOTE:
//    The input channel we will not close, we assume we do not own it
func NewWithParams(wg *sync.WaitGroup, in1 chan Packet, addr string, ct ConnType, outChanSize int) (*UDP, error) {
	c, cancel := context.WithCancel(context.Background())
	udp := UDP{out: make(chan Packet, outChanSize), addr: addr, ctx: c, can: cancel, in: in1, ct: ct}

	if err := udp.startConn(); err != nil {
		return nil, err
	}

	wg.Add(1)
	go udp.processInUDP(wg)

	wg.Add(1)
	go udp.processInChan(wg)

	return &udp, nil
}

// NewwithChan will create a UDP component with little fuss for the caller
// it takes just a port and input channel.  It will always setup a SERVER mode component
func NewWithChan(port int, in chan Packet) (*UDP, error) {
	addr := fmt.Sprintf(":%v", port)
	udpc, err := NewWithParams(new(sync.WaitGroup), in, addr, SERVER, 1)
	if err != nil {
		return nil, err
	}
	return udpc, nil
}

// New will create a UDP component with little fuss for the caller
// it takes just a port.  It will always setup a SERVER mode component
func New(port int) (*UDP, error) {
	udpc, err := NewWithChan(port, make(chan Packet, 1))
	if err != nil {
		return nil, err
	}
	return udpc, nil
}
