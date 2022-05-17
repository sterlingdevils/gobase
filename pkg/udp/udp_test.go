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

	// Wait for 1 second, then send a packet to our self, and display it, exit after 3 seconds
	for {
		select {
		case <-time.After(time.Second * 1):
			in <- Packet{Addr: &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 9092}, Data: []byte("Hello from Us.")}
		case p := <-udpcomp.OuputChan():
			dispPacket(p) // Display the Packet
			return
		}
	}

	// Output: 127.0.0.1:9092: [72 101 108 108 111 32 102 114 111 109 32 85 115 46]
}
