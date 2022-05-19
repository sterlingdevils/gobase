package retrycontainer_test

import (
	"github.com/sterlingdevils/gobase/pkg/obj"
	"github.com/sterlingdevils/gobase/retrycontainer"
)

func ExampleRetryContainer() {
	r, _ := retrycontainer.New[uint64, obj.Obj]()
	_ = r
	// Output:
}
