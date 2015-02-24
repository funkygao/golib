// Copyright 2011 Numrotron Inc.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.
//
// Developed at www.stathat.com by Patrick Crosby
// Contact us on twitter with any questions:  twitter.com/stat_hat

// The treap package provides a balanced binary tree datastructure, expected
// to have logarithmic height.
package treap

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// A Tree is the data structure that stores everything
type Tree struct {
	less    LessFunc
	overlap OverlapFunc
	count   int
	root    *Node
}

// LessFunc returns true if a < b
type LessFunc func(a, b interface{}) bool

// OverlapFunc return true if a overlaps b.  Optional.  Only used by overlap trees.
type OverlapFunc func(a, b interface{}) bool

// Key can be anything.  It will use LessFunc to compare keys.
type Key interface{}

// Item can be anything.
type Item interface{}

// A Node in the Tree.
type Node struct {
	key      Key
	item     Item
	priority int
	left     *Node
	right    *Node
}

func newNode(key Key, item Item, priority int) *Node {
	result := new(Node)
	result.key = key
	result.item = item
	result.priority = priority
	return result
}

// To create a Tree, you need to supply a LessFunc that can compare the
// keys in the Node.
func NewTree(lessfn LessFunc) *Tree {
	t := new(Tree)
	t.less = lessfn
	return t
}

// To create a tree that also lets you iterate by key overlap, supply a LessFunc
// and an OverlapFunc
func NewOverlapTree(lessfn LessFunc, overfn OverlapFunc) *Tree {
	t := new(Tree)
	t.less = lessfn
	t.overlap = overfn
	return t
}

// Remove everything from the tree.
func (t *Tree) Reset() {
	t.root = nil
	t.count = 0
}

// The length of the tree (number of nodes).
func (t *Tree) Len() int {
	return t.count
}

// Get an Item in the tree.
func (t *Tree) Get(key Key) Item {
	return t.get(t.root, key)
}

func (t *Tree) get(node *Node, key Key) Item {
	if node == nil {
		return nil
	}
	if t.less(key, node.key) {
		return t.get(node.left, key)
	}
	if t.less(node.key, key) {
		return t.get(node.right, key)
	}
	return node.item
}

// Returns true if there is an item in the tree with this key.
func (t *Tree) Exists(key Key) bool {
	return t.exists(t.root, key)
}

func (t *Tree) exists(node *Node, key Key) bool {
	if node == nil {
		return false
	}
	if t.less(key, node.key) {
		return t.exists(node.left, key)
	}
	if t.less(node.key, key) {
		return t.exists(node.right, key)
	}
	return true
}

// Insert an item into the tree.
func (t *Tree) Insert(key Key, item Item) {
	priority := rand.Int()
	t.root = t.insert(t.root, key, item, priority)
}

func (t *Tree) insert(node *Node, key Key, item Item, priority int) *Node {
	if node == nil {
		t.count++
		return newNode(key, item, priority)
	}
	if t.less(key, node.key) {
		node.left = t.insert(node.left, key, item, priority)
		if node.left.priority < node.priority {
			return t.leftRotate(node)
		}
		return node
	}
	if t.less(node.key, key) {
		node.right = t.insert(node.right, key, item, priority)
		if node.right.priority < node.priority {
			return t.rightRotate(node)
		}
		return node
	}

	// equal: replace the value
	node.item = item
	return node
}

func (t *Tree) leftRotate(node *Node) *Node {
	result := node.left
	x := result.right
	result.right = node
	node.left = x
	return result
}

func (t *Tree) rightRotate(node *Node) *Node {
	result := node.right
	x := result.left
	result.left = node
	node.right = x
	return result
}

// Split the tree by creating a tree with a node of priority -1 so it will be the root
func (t *Tree) Split(key Key) (*Node, *Node) {
	inserted := t.insert(t.root, key, nil, -1)
	return inserted.left, inserted.right
}

// Merge two trees together by supplying the root node of each tree.
func (t *Tree) Merge(left, right *Node) *Node {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}
	if left.priority < right.priority {
		result := left
		x := left.right
		result.right = t.Merge(x, right)
		return result
	}

	result := right
	x := right.left
	result.left = t.Merge(x, left)
	return result
}

// Delete the item from the tree that has this key.
func (t *Tree) Delete(key Key) {
	if t.Exists(key) == false {
		return
	}
	t.root = t.delete(t.root, key)
}

func (t *Tree) delete(node *Node, key Key) *Node {
	if node == nil {
		panic("key not found")
	}

	if t.less(key, node.key) {
		result := node
		x := node.left
		result.left = t.delete(x, key)
		return result
	}
	if t.less(node.key, key) {
		result := node
		x := node.right
		result.right = t.delete(x, key)
		return result
	}
	t.count--
	return t.Merge(node.left, node.right)
}

// Returns the height (depth) of the tree.
func (t *Tree) Height(key Key) int {
	return t.height(t.root, key)
}

func (t *Tree) height(node *Node, key Key) int {
	if node == nil {
		return 0
	}
	if t.less(key, node.key) {
		depth := t.height(node.left, key)
		return depth + 1
	}
	if t.less(node.key, key) {
		depth := t.height(node.right, key)
		return depth + 1
	}
	return 0
}

// Returns a channel of Items whose keys are in ascending order.
func (t *Tree) IterAscend() <-chan Item {
	c := make(chan Item)
	go func() {
		iterateInOrder(t.root, c)
		close(c)
	}()
	return c
}

func iterateInOrder(h *Node, c chan<- Item) {
	if h == nil {
		return
	}
	iterateInOrder(h.left, c)
	c <- h.item
	iterateInOrder(h.right, c)
}

// Returns a channel of keys in ascending order.
func (t *Tree) IterKeysAscend() <-chan Key {
	c := make(chan Key)
	go func() {
		iterateKeysInOrder(t.root, c)
		close(c)
	}()
	return c
}

func iterateKeysInOrder(h *Node, c chan<- Key) {
	if h == nil {
		return
	}
	iterateKeysInOrder(h.left, c)
	c <- h.key
	iterateKeysInOrder(h.right, c)
}

// Returns a channel of items that overlap key.
func (t *Tree) IterateOverlap(key Key) <-chan Item {
	c := make(chan Item)
	go func() {
		if t.overlap != nil {
			t.iterateOverlap(t.root, key, c)
		}
		close(c)
	}()
	return c
}

func (t *Tree) iterateOverlap(h *Node, key Key, c chan<- Item) {
	if h == nil {
		return
	}
	lessThanLower := t.overlap(h.key, key)
	greaterThanUpper := t.overlap(key, h.key)

	if !lessThanLower {
		t.iterateOverlap(h.left, key, c)
	}
	if !lessThanLower && !greaterThanUpper {
		c <- h.item
	}
	if !greaterThanUpper {
		t.iterateOverlap(h.right, key, c)
	}
}

// Returns the minimum item in the tree.
func (t *Tree) Min() Item {
	return min(t.root)
}

func min(h *Node) Item {
	if h == nil {
		return nil
	}
	if h.left == nil {
		return h.item
	}
	return min(h.left)
}

// Returns the maximum item in the tree.
func (t *Tree) Max() Item {
	return max(t.root)
}

func max(h *Node) Item {
	if h == nil {
		return nil
	}
	if h.right == nil {
		return h.item
	}
	return max(h.right)
}
