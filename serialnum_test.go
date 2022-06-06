// Package serialnum_test will test the public API of serialnum
package gobase_test

import (
	"fmt"

	"github.com/sterlingdevils/gobase"
)

func ExampleSerialNum() {
	in := []byte("Slice1")
	fmt.Println(len(in))
	fmt.Println(in)

	gbsn := gobase.SerialNum{}

	in = gbsn.AddRandom(in)
	fmt.Println(len(in))

	in, _, _ = gbsn.Remove(in)
	fmt.Println(in)
	fmt.Println(len(in))
	// Output:
	// 6
	// [83 108 105 99 101 49]
	// 14
	// [83 108 105 99 101 49]
	// 6
}

// ExampleNew
func ExampleSerialNum_New() {
	in := []byte("Slice1")
	in2 := []byte("Slice2")

	gbsn := gobase.SerialNum{}

	sncomp := gbsn.New()

	in = sncomp.AddInc(in)
	in2 = sncomp.AddInc(in2)

	s1, _ := gbsn.Uint64(in)
	s2, _ := gbsn.Uint64(in2)

	fmt.Print(s2 - s1)

	gobase.SnUint64.AddSn(in, 4959)
	gobase.SnUint.AddSn(in, 4959)
	// Output: 1
}
