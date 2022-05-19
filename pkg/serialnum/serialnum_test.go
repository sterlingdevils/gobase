// Package serialnum_test will test the public API of serialnum
package serialnum_test

import (
	"fmt"

	"github.com/sterlingdevils/gobase/pkg/serialnum"
)

func Example() {
	in := []byte("Slice1")
	fmt.Println(len(in))
	fmt.Println(in)

	in, _ = serialnum.AddRandom(in)
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

func ExampleAddInc() {
	in := []byte("Slice1")
	fmt.Println(len(in))
	fmt.Println(in)

	in, _ = serialnum.AddInc(in)
	fmt.Println(in)

	in2 := []byte("Slice2")
	in2, _ = serialnum.AddInc(in2)
	fmt.Println(in2)

	// Output:
	// 6
	// [83 108 105 99 101 49]
	// [0 0 0 0 0 0 0 0 83 108 105 99 101 49]
	// [1 0 0 0 0 0 0 0 83 108 105 99 101 50]
}
