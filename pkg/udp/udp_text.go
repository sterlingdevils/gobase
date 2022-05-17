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

func Example() {
	in := make(chan Packet, 5)
	defer close(in)
	wg := new(sync.WaitGroup)
	rx, err := New(wg, in, ":9092", SERVER, 1)
	if err != nil {
		log.Fatalln("error creating UDP")
	}
	defer rx.Close()

	in <- Packet{Addr: &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 9999}, Data: []byte("Hello from Us.")}
	out := rx.OuputChan()
	for {
		select {
		case <-time.After(time.Second * 20):
			return
		case p := <-out:
			dispPacket(p) // Display the Packet
		}
	}
}
