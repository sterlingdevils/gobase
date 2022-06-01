package signalutils

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func SignalHandler(wg *sync.WaitGroup, ctx context.Context, can context.CancelFunc) {
	defer wg.Done()
	defer can()

	ch := make(chan os.Signal, 10)
	signal.Notify(ch, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case s := <-ch:
			log.Println("received sig: ", s.String())
			switch s {
			case syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM:
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
