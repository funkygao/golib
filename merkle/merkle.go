package merkle

import (
	"errors"
	"hash"
)

// A node in the merkle tree
type Node struct {
	Hash  []byte
	Left  *Node
	Right *Node
}

// Creates a node given a hash function and data to hash
func NewNode(h hash.Hash, block []byte) (Node, error) {
	if h == nil || block == nil {
		return Node{}, nil
	}
	defer h.Reset()
	_, err := h.Write(block[:])
	if err != nil {
		return Node{}, err
	}
	return Node{Hash: h.Sum(nil)}, nil
}

// Contains all nodes
type Tree struct {
	// All nodes, linear
	Nodes []Node
	// Points to each level in the node. The first level contains the root node
	Levels [][]Node
}

func NewTree() Tree {
	return Tree{Nodes: nil, Levels: nil}
}

// Returns a slice of the leaf nodes in the tree, if available, else nil
func (self *Tree) Leaves() []Node {
	if self.Levels == nil {
		return nil
	} else {
		return self.Levels[len(self.Levels)-1]
	}
}

// Returns the root node of the tree, if available, else nil
func (self *Tree) Root() *Node {
	if self.Nodes == nil {
		return nil
	} else {
		return &self.Levels[0][0]
	}
}

// Returns all nodes at a given height, where height 1 returns a 1-element
// slice containing the root node, and a height of tree.Height() returns
// the leaves
func (self *Tree) GetNodesAtHeight(h uint64) []Node {
	if self.Levels == nil || h == 0 || h > uint64(len(self.Levels)) {
		return nil
	} else {
		return self.Levels[h-1]
	}
}

// Returns the height of this tree
func (self *Tree) Height() uint64 {
	return uint64(len(self.Levels))
}

// Generates the tree nodes
func (self *Tree) Generate(blocks [][]byte, hashf hash.Hash) error {
	blockCount := uint64(len(blocks))
	if blockCount == 0 {
		return errors.New("Empty tree")
	}
	height, nodeCount := CalculateHeightAndNodeCount(blockCount)
	levels := make([][]Node, height)
	nodes := make([]Node, nodeCount)

	// Create the leaf nodes
	for i, block := range blocks {
		node, err := NewNode(hashf, block)
		if err != nil {
			return err
		}
		nodes[i] = node
	}
	levels[height-1] = nodes[:len(blocks)]

	// Create each node level
	current := nodes[len(blocks):]
	h := height - 1
	for ; h > 0; h-- {
		below := levels[h]
		wrote, err := self.generateNodeLevel(below, current, hashf)
		if err != nil {
			return err
		}
		levels[h-1] = current[:wrote]
		current = current[wrote:]
	}

	self.Nodes = nodes
	self.Levels = levels
	return nil
}

// Creates all the non-leaf nodes for a certain height. The number of nodes
// is calculated to be 1/2 the number of nodes in the lower rung.  The newly
// created nodes will reference their Left and Right children.
// Returns the number of nodes added to current
func (self *Tree) generateNodeLevel(below []Node, current []Node,
	h hash.Hash) (uint64, error) {
	h.Reset()
	size := h.Size()
	data := make([]byte, size*2)
	end := (len(below) + (len(below) % 2)) / 2
	for i := 0; i < end; i++ {
		// Concatenate the two children hashes and hash them, if both are
		// available, otherwise reuse the hash from the lone left node
		node := Node{}
		ileft := 2 * i
		iright := 2*i + 1
		left := &below[ileft]
		var right *Node = nil
		if len(below) > iright {
			right = &below[iright]
		}
		if right == nil {
			b := data[:size]
			copy(b, left.Hash)
			node = Node{Hash: b}
		} else {
			copy(data[:size], below[ileft].Hash)
			copy(data[size:], below[iright].Hash)
			var err error
			node, err = NewNode(h, data)
			if err != nil {
				return 0, err
			}
		}
		// Point the new node to its children and save
		node.Left = left
		node.Right = right
		current[i] = node

		// Reset the data slice
		data = data[:]
	}
	return uint64(end), nil
}

// Returns the height and number of nodes in an unbalanced binary tree given
// number of leaves
func CalculateHeightAndNodeCount(leaves uint64) (height, nodeCount uint64) {
	height = calculateTreeHeight(leaves)
	nodeCount = calculateNodeCount(height, leaves)
	return
}

// Calculates the number of nodes in a binary tree unbalanced strictly on
// the right side.  Height is assumed to be equal to
// calculateTreeHeight(size)
func calculateNodeCount(height, size uint64) uint64 {
	if isPowerOfTwo(size) {
		return 2*size - 1
	}
	count := size
	prev := size
	i := uint64(1)
	for ; i < height; i++ {
		next := (prev + (prev % 2)) / 2
		count += next
		prev = next
	}
	return count
}

// Returns the height of a full, complete binary tree given nodeCount nodes
func calculateTreeHeight(nodeCount uint64) uint64 {
	if nodeCount == 0 {
		return 0
	} else if nodeCount == 1 {
		return 2
	} else {
		return logBaseTwo(nextPowerOfTwo(nodeCount)) + 1
	}
}

// Returns true if n is a power of 2
func isPowerOfTwo(n uint64) bool {
	// http://graphics.stanford.edu/~seander/bithacks.html#DetermineIfPowerOf2
	return n != 0 && (n&(n-1)) == 0
}

// Returns the next highest power of 2 above n, if n is not already a
// power of 2
func nextPowerOfTwo(n uint64) uint64 {
	if n == 0 {
		return 1
	}
	// http://graphics.stanford.edu/~seander/bithacks.html#RoundUpPowerOf2
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	n++
	return n
}

// Lookup table for integer log2 implementation
var log2lookup []uint64 = []uint64{
	0xFFFFFFFF00000000,
	0x00000000FFFF0000,
	0x000000000000FF00,
	0x00000000000000F0,
	0x000000000000000C,
	0x0000000000000002,
}

// Returns log2(n) assuming n is a power of 2
func logBaseTwo(x uint64) uint64 {
	if x == 0 {
		return 0
	}
	ct := uint64(0)
	for x != 0 {
		x >>= 1
		ct += 1
	}
	return ct - 1
}

// Returns the ceil'd log2 value of n
// See: http://stackoverflow.com/a/15327567
func ceilLogBaseTwo(x uint64) uint64 {
	y := uint64(1)
	if (x & (x - 1)) == 0 {
		y = 0
	}
	j := uint64(32)
	i := uint64(0)

	for ; i < 6; i++ {
		k := j
		if (x & log2lookup[i]) == 0 {
			k = 0
		}
		y += k
		x >>= k
		j >>= 1
	}

	return y
}
