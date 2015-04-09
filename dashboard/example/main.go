package main

import (
	"github.com/funkygao/golib/dashboard"
)

type mydata struct {
	i int
}

func (this *mydata) Data() int {
	this.i++
	return this.i
}

func main() {
	d := dashboard.New("test of dashboard", 5)
	g := d.AddGraph("graph1")
	g.AddLine("Heap", &mydata{i: 1})
	d.Launch(":8000")
}
