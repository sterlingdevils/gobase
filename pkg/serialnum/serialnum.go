// Package serialnum adds and removes serial numbers to byte slices
package serialnum

import "crypto/rand"

const (
	SERIALNUMSIZE = 8
)

// Adds serial number to byte slice
func Add(in []byte) []byte {
	sn := make([]byte, SERIALNUMSIZE)
	rand.Read(sn)
	return append(sn, in...)
}

// Removes serial number from byte slice or returns empty array if the size is too small
func Remove(in []byte) []byte {
	if len(in) < SERIALNUMSIZE {
		return make([]byte, 0)
	}
	return in[SERIALNUMSIZE:]
}
