// Directed Acyclic Graph implementation in golang
package dag

import (
	"fmt"
	"github.com/funkygao/golib/str"
	"os"
)

type Dag map[string]*Node

type Node struct {
	name         string
	indegree     int
	val          interface{}
	dependencies []string
	children     []*Node
}

func New() Dag {
	return make(map[string]*Node)
}

func (this Dag) AddEdge(from, to string) {
	fromNode := this[from]
	toNode := this[to]
	fromNode.children = append(fromNode.children, toNode)
	toNode.indegree++
}

func (this Dag) MakeDotGraph(fn string) {
	file, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	sb := str.NewStringBuilder()
	sb.WriteString("digraph depgraph {\n\trankdir=LR;\n")
	for _, node := range this {
		node.dotGraph(sb)
	}
	sb.WriteString("}\n")
	file.WriteString(sb.String())
}

func (this *Dag) HasPathTo(that string) bool {
	return false
}

func (this *Node) dotGraph(sb *str.StringBuilder) {
	if len(this.dependencies) == 0 {
		sb.WriteString(fmt.Sprintf("\t\"%s\";\n", this.name))
	}
}
