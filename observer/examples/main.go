package main

import (
    ob "github.com/funkygao/golib/observer"
	"fmt"
	"time"
)

func main() {
	go func() {
		for {
			time.Sleep(time.Duration(1) * time.Second)
			ob.Publish("foo", time.Now().Unix())
		}
	}()

	eventCh1 := make(chan interface{})
	ob.Subscribe("foo", eventCh1)
	go func() {
		for {
			data := <-eventCh1
			fmt.Printf("sub1: %#v\n", data)
		}
	}()

	eventCh2 := make(chan interface{})
	ob.Subscribe("foo", eventCh2)
	go func() {
		for {
			data := <-eventCh2
			fmt.Printf("sub2: %v\n", data)
		}
	}()

	<-time.After(6 * time.Second)

	// Output:
}
