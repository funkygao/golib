package main

import (
	"github.com/funkygao/golib/dashboard"
)

type myline struct {
	i int
}

func (this *myline) Data() int {
	this.i++
	return this.i
}

func main() {
	d := dashboard.New("test", 10)
	gcGraph := d.AddGraph("gc")
	gcGraph.AddLine("Heap", &myline{i: 1})
	d.Launch(":8000")
}
