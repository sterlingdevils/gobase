package obj_test

import (
	"fmt"
	"time"

	"github.com/sterlingdevils/gobase/obj"
)

// ExampleNew will show how to create one of us
func ExampleNew() {
	o, _ := obj.New(time.Millisecond * 100)
	fmt.Println(o.Context().Err())

	time.Sleep(time.Millisecond * 200)
	fmt.Println(o.Context().Err())
	// Output:
	// <nil>
	// context deadline exceeded
}