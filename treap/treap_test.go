package treap

import (
	"fmt"
	"math/rand"
	"testing"
)

func init() {
	// so that every run is the same seq of rand numbers
	rand.Seed(0)
}

func StringLess(p, q interface{}) bool {
	return p.(string) < q.(string)
}

func IntLess(p, q interface{}) bool {
	return p.(int) < q.(int)
}

func TestEmpty(t *testing.T) {
	tree := NewTree(StringLess)
	if tree.Len() != 0 {
		t.Errorf("expected tree len 0")
	}

	x := tree.Get("asdf")
	if x != nil {
		t.Errorf("expected nil for nonexistent key")
	}
}

func TestInsert(t *testing.T) {
	tree := NewTree(StringLess)
	tree.Insert("xyz", "adsf")
	x := tree.Get("xyz")
	if x != "adsf" {
		t.Errorf("expected adsf, got %v", x)
	}
}

func TestFromDoc(t *testing.T) {
	tree := NewTree(IntLess)
	tree.Insert(5, "a")
	tree.Insert(7, "b")
	x := tree.Get(5)
	if x != "a" {
		t.Errorf("expected a, got %v", x)
	}
	x = tree.Get(7)
	if x != "b" {
		t.Errorf("expected b, got %v", x)
	}
	tree.Insert(2, "c")
	x = tree.Get(2)
	if x != "c" {
		t.Errorf("exepcted c, got %v", x)
	}
	tree.Insert(2, "d")
	x = tree.Get(2)
	if x != "d" {
		t.Errorf("exepcted d, got %v", x)
	}
	tree.Delete(5)
	if tree.Exists(5) {
		t.Errorf("expected 5 to be removed from tree")
	}
}

func TestBalance(t *testing.T) {
	tree := NewTree(IntLess)
	for i := 0; i < 1000; i++ {
		tree.Insert(i, false)
	}
	for i := 0; i < 1000; i += 50 {
		fmt.Printf("%d: height = %d\n", i, tree.Height(i))
	}
}

// tests copied from petar's llrb
func TestCases(t *testing.T) {
	tree := NewTree(IntLess)
	tree.Insert(1, true)
	tree.Insert(1, false)
	if tree.Len() != 1 {
		t.Errorf("expecting len 1")
	}
	if !tree.Exists(1) {
		t.Errorf("expecting to find key=1")
	}

	tree.Delete(1)
	if tree.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if tree.Exists(1) {
		t.Errorf("not expecting to find key=1")
	}

	tree.Delete(1)
	if tree.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if tree.Exists(1) {
		t.Errorf("not expecting to find key=1")
	}
}

func TestReverseInsertOrder(t *testing.T) {
	tree := NewTree(IntLess)
	n := 100
	for i := 0; i < n; i++ {
		tree.Insert(n-i, true)
	}
	c := tree.IterKeysAscend()
	for j, item := 1, <-c; item != nil; j, item = j+1, <-c {
		if item.(int) != j {
			t.Fatalf("bad order")
		}
	}
}

func TestIterateOverlapNoFunc(t *testing.T) {
	tree := NewTree(IntLess)
	n := 100
	for i := 0; i < n; i++ {
		tree.Insert(n-i, true)
	}
	for v := range tree.IterateOverlap(50) {
		t.Errorf("didn't expect to have any overlap since fn not defined: %v", v)
	}
}

type BucketKey struct {
	Start    int64
	Duration int64
}

type BucketVal struct {
	Start    int64
	Duration int64
	Value    float64
}

func BucketLess(a, b interface{}) bool {
	aa := a.(*BucketKey)
	bb := b.(*BucketKey)
	if aa.Start < bb.Start {
		return true
	}
	if aa.Start == bb.Start {
		return aa.Duration < bb.Duration
	}
	return false
}

func BucketOverlap(a, b interface{}) bool {
	aa := a.(*BucketKey)
	bb := b.(*BucketKey)
	return aa.Start+aa.Duration < bb.Start
}

func TestIterateOverlap(t *testing.T) {
	tree := NewOverlapTree(BucketLess, BucketOverlap)

	tree.Insert(&BucketKey{100, 10}, &BucketVal{100, 10, 5.0})
	tree.Insert(&BucketKey{110, 10}, &BucketVal{110, 10, 6.0})
	tree.Insert(&BucketKey{120, 10}, &BucketVal{120, 10, 7.0})
	tree.Insert(&BucketKey{130, 10}, &BucketVal{130, 10, 8.0})

	for v := range tree.IterateOverlap(&BucketKey{105, 7}) {
		fmt.Printf("val: %v\n", v)
	}
}

/*
before:                            after:

        A                             B
       / \                           / \
      B   C                         D   A
     / \ / \                           / \
    D  E F G                          E   C
                                         / \
                                        F   G

*/
func TestLeftRotate(t *testing.T) {
	// create a tree by hand...
	a := newNode("a", "a", 1)
	b := newNode("b", "b", 2)
	c := newNode("c", "c", 3)
	d := newNode("d", "d", 4)
	e := newNode("e", "e", 5)
	f := newNode("f", "g", 5)
	g := newNode("g", "g", 5)
	a.left = b
	a.right = c
	b.left = d
	b.right = e
	c.left = f
	c.right = g

	x := new(Tree)
	root := x.leftRotate(a)
	if root != b {
		t.Errorf("expected root to be b")
	}
	if root.left != d {
		t.Errorf("expected root.left to be d")
	}
	if root.right != a {
		t.Errorf("expected root.right to be a")
	}
	if a.left != e {
		t.Errorf("expected a.left to be e")
	}
	if a.right != c {
		t.Errorf("expected a.right to be c")
	}
	if c.left != f {
		t.Errorf("expected c.left to be f")
	}
	if c.right != g {
		t.Errorf("expected c.right to be g")
	}
}

/*
before:                            after:

        A                             C
       / \                           / \
      B   C                         A   G
     / \ / \                       / \
    D  E F G                      B   F
                                 / \
                                D  E

*/
func TestRightRotate(t *testing.T) {
	// create a tree by hand...
	a := newNode("a", "a", 1)
	b := newNode("b", "b", 2)
	c := newNode("c", "c", 3)
	d := newNode("d", "d", 4)
	e := newNode("e", "e", 5)
	f := newNode("f", "g", 5)
	g := newNode("g", "g", 5)
	a.left = b
	a.right = c
	b.left = d
	b.right = e
	c.left = f
	c.right = g

	x := new(Tree)
	root := x.rightRotate(a)
	if root != c {
		t.Errorf("expected root to be c")
	}
	if root.left != a {
		t.Errorf("expected root.left to be a")
	}
	if root.right != g {
		t.Errorf("expected root.right to be g")
	}
	if a.left != b {
		t.Errorf("expected a.left to be b")
	}
	if a.right != f {
		t.Errorf("expected a.right to be f")
	}
	if b.left != d {
		t.Errorf("expected b.left to be d")
	}
	if b.right != e {
		t.Errorf("expected b.right to be e")
	}
}

func treeOfInts(ints []int) (tree *Tree) {
	tree = NewTree(IntLess)
	for _, i := range ints {
		tree.Insert(i, i)
	}
	return
}

func BenchmarkInsert(b *testing.B) {
	b.StopTimer()
	ints := rand.Perm(b.N)
	b.StartTimer()
	_ = treeOfInts(ints)
}

func BenchmarkDelete(b *testing.B) {
	b.StopTimer()
	ints := rand.Perm(b.N)
	tree := treeOfInts(ints)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.Delete(i)
	}
}

func BenchmarkLookup(b *testing.B) {
	b.StopTimer()
	ints := rand.Perm(b.N)
	tree := treeOfInts(ints)
	b.StartTimer()
	for j := 0; j < 10; j++ {
		for i := 0; i < len(ints)/10; i++ {
			_ = tree.Exists(ints[i])
		}
	}
}
