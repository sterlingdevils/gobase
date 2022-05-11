package main

import (
	"fmt"
	"log"
	"sdbase/pkg/chantools"
	"sdbase/pkg/dirscan"
	"sdbase/pkg/metrics"
	"sdbase/pkg/udprx"
	"sdbase/pkg/udptx"
	"strconv"
	"sync"
	"time"
)

func mainloop(wg *sync.WaitGroup) {
	defer wg.Done()

	u, _ := udprx.New(wg, 9999)
	udptx.New(wg, "localhost:9998", u.OutputChan())
	defer u.Close()

	mex := metrics.New()

	ds, err := dirscan.New(wg, ".", time.Second*2, 0)
	if err != nil {
		log.Fatalln(err)
	}
	defer ds.Close()
	ds.SetMetric(mex)

	c2 := make(chan string)
	wg.Add(1)
	go func() {
		defer chantools.RecoverFromClosedChan()
		defer wg.Done()
		for i := 0; i < 10; i++ {
			c2 <- strconv.Itoa(i)
			time.Sleep(time.Second)
		}
	}()
	defer close(c2)

	mux := chantools.New(wg, ds.OutputChan(), c2, 0)
	defer mux.Close()

	inchan := mux.OutputChan()

	//	fnchan := mux.GetChan()
	//	achan := chantools.AsyncSkip(wg, fnchan, 0)

	wg.Add(1)
	go chantools.DisplayChan(wg, inchan)

	time.Sleep(time.Second * 55)

	if v, err := mex.GetValue("scanDir.Size"); err == nil {
		fmt.Println(v)
	}
}

func main() {
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go mainloop(wg)

	wg.Wait()
}
