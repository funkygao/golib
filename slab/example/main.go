package main

import (
	"fmt"
	"github.com/funkygao/golib/slab"
)

func main() {
	arena := slab.NewArena(32, 1<<9, 2., nil)
	fmt.Println("1")
	buf, _ := arena.Alloc(130)
	fmt.Println("2")
	arena.Alloc(131)
	arena.AddRef(buf)
}
