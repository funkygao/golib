package observer

import (
	"fmt"
	"time"
)

func ExampleObserver() {
	go func() {
		for {
			time.Sleep(time.Duration(1) * time.Second)
			Publish("foo", time.Now().Unix())
		}
	}()

	eventCh1 := make(chan interface{})
	Subscribe("foo", eventCh1)
	go func() {
		for {
			data := <-eventCh1
			fmt.Printf("sub1: %#v\n", data)
		}
	}()

	eventCh2 := make(chan interface{})
	Subscribe("foo", eventCh2)
	go func() {
		for {
			data := <-eventCh2
			fmt.Printf("sub2: %v\n", data)
		}
	}()

	<-time.After(6 * time.Second)

	// Output:
}
