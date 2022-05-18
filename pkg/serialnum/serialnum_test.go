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
	in = serialnum.Add(in)
	fmt.Println(len(in))
	in = serialnum.Remove(in)
	fmt.Println(in)
	fmt.Println(len(in))
	// Output:
	// 6
	// [83 108 105 99 101 49]
	// 14
	// [83 108 105 99 101 49]
	// 6
}
