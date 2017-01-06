// Package top provides unix top like UI framework.
//
//
/*
Usage example:

	package main

	import (
		"fmt"
		"time"

		"github.com/funkygao/golib/top"
	)

	func main() {
		t := top.New("Id|Value", "%2s %8s")
		go func() {
			for i := 0; i < 1000; i++ {
				rows := make([]string, 0)

				if i%2 == 0 {
					rows = append(rows, fmt.Sprintf("%d|%d", 1, 100))
				} else {
					rows = append(rows, fmt.Sprintf("%d|%d", 2, 200))
				}

				t.Refresh(rows)
				time.Sleep(time.Second)
			}
		}()

		if err := t.Start(); err != nil {
			panic(err)
		}
	}

*/
package top
