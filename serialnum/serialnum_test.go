// Package serialnum_test will test the public API of serialnum
package serialnum_test

import (
	"fmt"

	"github.com/sterlingdevils/gobase/serialnum"
)

func Example() {
	in := []byte("Slice1")
	fmt.Println(len(in))
	fmt.Println(in)

	in = serialnum.AddRandom(in)
	fmt.Println(len(in))

	in, _, _ = serialnum.Remove(in)
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
func ExampleNew() {
	in := []byte("Slice1")
	in2 := []byte("Slice2")

	sncomp := serialnum.New()

	in = sncomp.AddInc(in)
	in2 = sncomp.AddInc(in2)

	s1, _ := serialnum.Uint64(in)
	s2, _ := serialnum.Uint64(in2)

	fmt.Print(s2 - s1)

	serialnum.SnUint64.AddSn(in, 4959)
	serialnum.SnUint.AddSn(in, 4959)

	// Output: 1
}
