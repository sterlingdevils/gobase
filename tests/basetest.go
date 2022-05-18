package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/sterlingdevils/gobase/pkg/chantools"
	"github.com/sterlingdevils/gobase/pkg/dirscan"
	"github.com/sterlingdevils/gobase/pkg/metrics"
	"github.com/sterlingdevils/gobase/pkg/udp"
)

func dispPacket(p udp.Packet) {
	fmt.Printf("%v: %v\n", p.Addr, p.Data)
}

// test UDP component
func testUDP(wg *sync.WaitGroup) {
	defer wg.Done()

	// Create a UDP component
	//  First step is to create an input channel
	//  Next step call New
	//  Then defer the close for the UDP component
	//  Then get the output Channel
	//  Use it
	in := make(chan udp.Packet, 5)
	rx, err := udp.NewwithParams(wg, in, "localhost:9999", udp.CLIENT, 1)
	if err != nil {
		log.Fatalln("error creating UDP")
	}
	defer rx.Close()

	// Test sending a packet
	//	in <- udp.Packet{Addr: &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 9999}, Data: []byte("Hello from Us.")}
	in <- udp.Packet{Data: []byte("Hello from Number Two")}

	// Receive any packets on the channel that Received Packets goto
	for p := range rx.OuputChan() {
		dispPacket(p) // Display the Packet
		in <- p       // Send them Back out to the sender
	}

	close(in)
}

func mainloop(wg *sync.WaitGroup) {
	defer wg.Done()

	// Test out the UDP Component
	wg.Add(1)
	go testUDP(wg)

	// Testout Metricss
	mex := metrics.New()

	ds, err := dirscan.New(wg, ".", time.Second*2, 0)
	if err != nil {
		log.Fatalln(err)
	}
	defer ds.Close()
	ds.SetMetric(mex)

	c2 := make(chan string)
	wg.Add(1)
	go func() {
		defer chantools.RecoverFromClosedChan()
		defer wg.Done()
		for i := 0; i < 10; i++ {
			c2 <- strconv.Itoa(i)
			time.Sleep(time.Second)
		}
	}()
	defer close(c2)

	mux := chantools.New(wg, ds.OutputChan(), c2, 0)
	defer mux.Close()

	inchan := mux.OutputChan()

	//	fnchan := mux.GetChan()
	//	achan := chantools.AsyncSkip(wg, fnchan, 0)

	wg.Add(1)
	go chantools.DisplayChan(wg, inchan)

	time.Sleep(time.Second * 55)

	if v, err := mex.GetValue("scanDir.Size"); err == nil {
		fmt.Println(v)
	}
}

func main() {
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go mainloop(wg)

	wg.Wait()
}
