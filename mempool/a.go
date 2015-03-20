package mempool

type Node struct {
	nextFree    *Node
	left, right *Node
}

type Slab struct {
	free  *Node
	Nodes []Node
}

func NewSlab(size int) *Slab {
	s := &Slab{Nodes: make([]Node, size)}
	s.free = &s.Nodes[0]
	prev := s.free
	for i := 1; i < len(s.Nodes); i++ {
		curr := &s.Nodes[i]
		prev.nextFree = curr
		prev = curr
	}
	return s
}

func (s *Slab) Malloc() *Node {
	n := s.free
	if n == nil {
		n = &Node{}
	}
	s.free = n.nextFree
	n.nextFree = n // Slab allocated Node marker
	return n
}

func (s *Slab) Free(n *Node) {
	if n.nextFree == n {
		n.nextFree = s.free
		s.free = n
	}
	// Nodes that were not slab allocated are left tn the garbage collector
}
