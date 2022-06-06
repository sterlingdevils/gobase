// Package serialnum adds and removes serial numbers to byte slices
package gobase

import (
	"encoding/binary"
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
)

const (
	// number of bytes in a serial number
	SERIALNUMSIZE = 8
	UINT64SIZE    = 8
	INTSIZE       = 4
)

type SerialNum struct {
	// Holds the current incremented serial number, This is a Package wide Global
	currentSn    uint64
	currentMutex sync.Mutex
}

type snuint struct{}
type snuint64 struct{}

var (
	SnUint   snuint
	SnUint64 snuint64
)

func (snuint64) AddSn(in []byte, sn uint64) []byte {
	data := make([]byte, SERIALNUMSIZE+len(in))
	binary.LittleEndian.PutUint64(data, sn)
	copy(data[SERIALNUMSIZE:], in)

	return data
}

func (snuint) AddSn(in []byte, sn uint32) []byte {
	data := make([]byte, INTSIZE+len(in))
	binary.LittleEndian.PutUint32(data, sn)
	copy(data[INTSIZE:], in)

	return data
}

// addsn will take a uint64 and prepend it to the in slice, returns a new slice based on a new array
func (*SerialNum) addsn(in []byte, sn uint64) []byte {
	data := make([]byte, SERIALNUMSIZE+len(in))
	binary.LittleEndian.PutUint64(data, sn)
	copy(data[SERIALNUMSIZE:], in)

	return data
}

// Next will return the next sn in sequence
func (s *SerialNum) Next() uint64 {
	s.currentMutex.Lock()
	sn := atomic.LoadUint64(&s.currentSn)
	atomic.AddUint64(&s.currentSn, 1)
	s.currentMutex.Unlock()
	return sn
}

// AddInc adds an incrementing serial number to byte slice
func (s *SerialNum) AddInc(in []byte) []byte {
	return s.addsn(in, s.Next())
}

// AddSn adds a passed in serial number to byte slice
func (s *SerialNum) AddSn(in []byte, sn uint64) []byte {
	return s.addsn(in, sn)
}

// AddRandom adds a random serial number to byte slice
func (s *SerialNum) AddRandom(in []byte) []byte {
	sn := uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
	return s.addsn(in, sn)
}

// Remove will split the serial number from byte slice
// returns the data that is after the sn, the sn, and an error it there was a problem
func (*SerialNum) Remove(in []byte) (data []byte, sn []byte, err error) {
	if len(in) < SERIALNUMSIZE {
		return nil, nil, errors.New("passed in slice is smaller than a serialnumber")
	}

	return in[SERIALNUMSIZE:], in[:SERIALNUMSIZE], nil
}

// Sn will return a slice with the serial number
// note: The returned slice has the same underlying array
func (*SerialNum) Sn(in []byte) (sn []byte, err error) {
	if len(in) < SERIALNUMSIZE {
		return nil, errors.New("passed in slice is smaller than a serialnumber")
	}

	return in[:SERIALNUMSIZE], nil
}

// Uint64 will return a slice with the serial number
// note: The returned slice has the same underlying array
func (*SerialNum) Uint64(in []byte) (sn uint64, err error) {
	if len(in) < SERIALNUMSIZE {
		return 0, errors.New("passed in slice is smaller than a serialnumber")
	}

	return binary.LittleEndian.Uint64(in), nil
}

// New returns a SerialNum component, each instant has its own sn counter
func (*SerialNum) New() *SerialNum {
	sn := SerialNum{currentSn: 0}
	return &sn
}
