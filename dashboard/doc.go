// Package dashboard implements a visualization feature that
// listen on a specified http port and draws stacked trend lines over time.
//
// Usage:
//    d := dashboard.New("gc", 10)
//    if err := d.Validate(); err != nil {
//	      panic(err)
//	  }
//	  gcGraph := dashboard.NewGraph("gc")
//    gcGraph.AddLine("Heap", myline{i: 1})
//    d.AddGraph(gcGraph)
//	  d.Launch(":8000")
package dashboard
