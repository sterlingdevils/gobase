// Package udp_test will test the public API of udp
package udp_test

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/sterlingdevils/gobase/pkg/udp"
)

// Example of how to create,receive and send packets
//
// This will create a UDP component and then send a packet,
// receive the udp, then display it, and check the display is
// correct.
func Example() {
	wg := new(sync.WaitGroup)

	in := make(chan udp.Packet, 5)
	defer close(in)

	// Must pass in the input channel as we dont assume we own it
	udpcomp, err := udp.New(wg, in, ":9092", udp.SERVER, 1)
	if err != nil {
		log.Fatalln("error creating UDP")
	}
	defer udpcomp.Close()

	// Wait for 1 second, then send a packet to our self, and display it, exit after 3 seconds
	for {
		select {
		case <-time.After(time.Second * 1):
			in <- udp.Packet{Addr: &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 9092}, Data: []byte("Hello from Us.")}
		case p := <-udpcomp.OuputChan():
			fmt.Printf("%v: %v\n", p.Addr, p.Data)
			return
		}
	}

	// Output: 127.0.0.1:9092: [72 101 108 108 111 32 102 114 111 109 32 85 115 46]
}
