package dag

import (
	"testing"
)

func TestMakeDotFile(t *testing.T) {
	d := New()
	d.AddVertex("s1", 1)
	d.AddVertex("s2", 2)
	d.AddVertex("t1", 12)
	d.AddVertex("t2", 13)
	d.AddEdge("s1", "t1")
	d.AddEdge("s1", "t2")
	d.AddEdge("s2", "t2")
	d.MakeDotGraph("test.dot")
	t.Logf("dot -o test.png -T png test.dot")
}
