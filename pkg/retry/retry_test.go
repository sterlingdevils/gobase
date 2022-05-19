package retry_test

import (
	"fmt"
	"log"
	"time"

	"github.com/sterlingdevils/gobase/pkg/obj"
	"github.com/sterlingdevils/gobase/pkg/retry"
)

func Example() {
	retry, err := retry.New()
	if err != nil {
		return
	}

	retry.Close()
	// Output:
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
		fmt.Println(o.Sn)
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