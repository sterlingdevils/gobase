// Package serialnum adds and removes serial numbers to byte slices
package serialnum

import (
	"encoding/binary"
	"errors"
	"math/rand"
	"sync/atomic"
	"time"
)

const (
	// number of bytes in a serial number
	SERIALNUMSIZE = 8
)

var (
	// Holds the current incremented serial number
	currentSn uint64
)

// init is called when this package is loaded
func init() {
	currentSn = 0
	rand.Seed(time.Now().UnixNano())
}

// Adds an incrementing serial number to byte slice
func AddInc(in []byte) ([]byte, error) {
	sn := make([]byte, SERIALNUMSIZE)
	binary.LittleEndian.PutUint64(sn, currentSn)

	// Add one to currentSn safely
	atomic.AddUint64(&currentSn, 1)

	return append(sn, in...), nil
}

// Adds serial number to byte slice
func AddRandom(in []byte) ([]byte, error) {
	sn := make([]byte, SERIALNUMSIZE)
	_, err := rand.Read(sn)
	if err != nil {
		return nil, err
	}

	return append(sn, in...), nil
}

// Removes serial number from byte slice or returns empty array if the size is too small
// returns the data that is after the sn, the sn, and an error it there was a problem
func Remove(in []byte) (data []byte, sn []byte, err error) {
	if len(in) < SERIALNUMSIZE {
		return nil, nil, errors.New("passed in slice is smaller than a serialnumber")
	}

	return in[SERIALNUMSIZE:], in[:SERIALNUMSIZE], nil
}

// Sn will return a slice with the serial number
// note: The returned slice has the same underlying array
func Sn(in []byte) (sn []byte, err error) {
	if len(in) < SERIALNUMSIZE {
		return nil, errors.New("passed in slice is smaller than a serialnumber")
	}

	return in[:SERIALNUMSIZE], nil
}
