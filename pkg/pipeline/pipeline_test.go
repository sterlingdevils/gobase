package pipeline_test

import (
	"log"

	"github.com/sterlingdevils/gobase/pkg/buffer"
	"github.com/sterlingdevils/gobase/pkg/pipeline"
)

type Node struct {
	Name string
}

func Example() {
	b, _ := buffer.NewWithPipeline[Node](
		5,
		pipeline.Checkerror(buffer.New[Node](10)))
	_ = b
	b.Close()
	// Output:
}

func Example_betterpipeline() {
	b1, err := buffer.New[any](10)
	if err != nil {
		log.Fatal(err)
	}

	b2, err := buffer.NewWithPipeline[any](5, b1)
	if err != nil {
		log.Fatal(err)
	}

	b2.Close()
	// Output:
}
