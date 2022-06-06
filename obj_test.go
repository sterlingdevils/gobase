package gobase_test

import (
	"fmt"
	"time"

	"github.com/sterlingdevils/gobase"
)

// ExampleNew will show how to create one of us
func ExampleObj_New() {
	o, _ := gobase.Obj{}.New(time.Millisecond * 100)
	fmt.Println(o.Context().Err())

	time.Sleep(time.Millisecond * 200)
	fmt.Println(o.Context().Err())
	// Output:
	// <nil>
	// context deadline exceeded
}
