package retry_test

import (
	"fmt"
	"log"
	"time"

	"github.com/sterlingdevils/gobase/pkg/obj"
	"github.com/sterlingdevils/gobase/pkg/retry"
	"github.com/sterlingdevils/gobase/pkg/serialnum"
)

func Example() {
	retry, err := retry.New()
	if err != nil {
		return
	}

	retry.Close()
	// Output:
}

func ExampleRetry_inout() {
	sn := serialnum.New()
	retry, err := retry.New()
	if err != nil {
		log.Fatal("error on create")
	}

	for i := 0; i < 10; i++ {
		o, _ := obj.New(2 * time.Second)
		o.Sn = sn.Next()
		retry.ObjIn() <- o
	}

	go func() {
		time.Sleep(3 * time.Second)
		retry.Close()
	}()

	for o := range retry.ObjOut() {
		fmt.Println(o.Key())
	}

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
}

func ExampleRetry_pointercheck() {
	retry, err := retry.New()
	if err != nil {
		log.Fatal("error on create")
	}

	// Check that we are passing pointer
	o, _ := obj.New(5 * time.Second)
	retry.ObjIn() <- o

	o.Sn = 5
	retry.ObjIn() <- o

	go func() {
		time.Sleep(2 * time.Second)
		retry.Close()
	}()

	for o := range retry.ObjOut() {
		fmt.Println(o.Key())
	}

	// Output:
	// 5
	// 5
}

func ExampleRetry_Close() {
	retry, err := retry.New()
	if err != nil {
		log.Fatal("error on create")
	}
	retry.Close()
	// Output:
}
