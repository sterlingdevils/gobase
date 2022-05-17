package udp

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func dispPacket(p Packet) {
	fmt.Printf("%v: %v\n", p.Addr, p.Data)
}

// Example of how to create,receive and send packets
func Example() {
	wg := new(sync.WaitGroup)

	in := make(chan Packet, 5)
	defer close(in)

	// Must pass in the input channel as we dont assume we own it
	udpcomp, err := New(wg, in, ":9092", SERVER, 1)
	if err != nil {
		log.Fatalln("error creating UDP")
	}
	defer udpcomp.Close()

	// Example Write packet to UDP
	in <- Packet{Addr: &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 9999}, Data: []byte("Hello from Us.")}

	for {
		select {
		case <-time.After(time.Second * 20):
			return
		case p := <-udpcomp.OuputChan():
			dispPacket(p) // Display the Packet
		}
	}
}
