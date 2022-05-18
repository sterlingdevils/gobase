// Package serialnum_test will test the public API of serialnum
package serialnum_test

import (
	"fmt"

	"github.com/sterlingdevils/gobase/pkg/serialnum"
)

func Example() {
	in := []byte("Slice1")
	fmt.Println(in)
	in = serialnum.Add(in)
	fmt.Println(in)
	in = serialnum.Remove(in)
	fmt.Println(in)
}
