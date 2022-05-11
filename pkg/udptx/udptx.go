package udptx

import (
	"log"
	"net"
	"sync"
)

const (
	maxDSz = 65507
)

type UDPTx struct {
	addr string
	in1  <-chan []byte
	conn *net.UDPConn
}

func (tx *UDPTx) newUDPConn() error {
	a, err := net.ResolveUDPAddr("udp4", tx.addr)
	if err != nil {
		return err
	}

	tx.conn, err = net.DialUDP("udp4", nil, a)
	if err != nil {
		return err
	}

	return nil
}

func (tx *UDPTx) mainloop(wg *sync.WaitGroup) {
	defer wg.Done()

	err := tx.newUDPConn()
	if err != nil {
		return
	}
	defer tx.conn.Close()

	for b := range tx.in1 {
		_, err = tx.conn.Write(b)
		if err != nil {
			log.Println("udp write failed")
		}
	}
}

// ------- Public Methods -------

// Close
func (tx *UDPTx) Close() {
}

// New
func New(wg *sync.WaitGroup, addr string, in1 <-chan []byte) (*UDPTx, error) {
	tx := UDPTx{addr: addr, in1: in1}

	wg.Add(1)
	go tx.mainloop(wg)

	return &tx, nil
}
