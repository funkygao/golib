/*
Package dashboard helps for visualization that listen on a
http port and draws stacked trend lines over time.

Usage:
	d := dashboard.New("test of dashboard", 5)
	g := d.AddGraph("graph1")
	g.AddLine("data1", &mydata{i: 1})
	g.AddLine("data2", &mydata2{i: 0})
	d.Launch(":8000")

*/
package dashboard
