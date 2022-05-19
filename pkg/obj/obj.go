package obj

import (
	"context"
	"time"
)

type Obj struct {
	Sn  uint64
	Ctx context.Context
	Can context.CancelFunc

	Data []byte
}

func New(timeout time.Duration) (*Obj, error) {
	c, cancel := context.WithTimeout(context.Background(), timeout)
	o := Obj{Ctx: c, Can: cancel}

	return &o, nil
}
