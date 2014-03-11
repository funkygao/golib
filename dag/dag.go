// Directed Acyclic Graph implementation in golang
package dag

import (
	"fmt"
	"github.com/funkygao/golib/str"
	"os"
)

type Dag struct {
	nodes map[string]*Node
}

type Node struct {
	name         string
	indegree     int
	val          interface{}
	dependencies []string
	children     []*Node
}

func New() *Dag {
	this := new(Dag)
	this.nodes = make(map[string]*Node)
	return this
}

func (this *Dag) AddVertex(name string) *Node {
	node := &Node{name: name}
	this.nodes[name] = node
	return node
}

func (this *Dag) AddEdge(from, to string) {
	fromNode := this.nodes[from]
	toNode := this.nodes[to]
	fromNode.children = append(fromNode.children, toNode)
	toNode.indegree++
}

func (this *Dag) MakeDotGraph(fn string) {
	file, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	sb := str.NewStringBuilder()
	sb.WriteString("digraph depgraph {\n\trankdir=LR;\n")
	for _, node := range this.nodes {
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

func (this *Node) Children() []*Node {
	return this.children
}
