package trie

import "fmt"

// Trie represents binary trie for prefix lookup.
type Trie[V any] struct {
	Root   *Node[V]
	Height uint
	Prefix string // when trie is subtrie this field represents old prefix of root
}

// New creates new trie.
func New[V any]() *Trie[V] {
	return &Trie[V]{
		Root:   new(Node[V]),
		Height: 1,
		Prefix: "",
	}
}

// Divide divides the binary trie into different levels
// based on these k-bit subtries. If k = 4 thus, level 0 contains
// from prefix length 0 to prefix length 7, and so on.
// Each level contains one or more subtries.
// nolint: funlen, gocognit, cyclop
func (t *Trie[V]) Divide(stride uint) [][]*Trie[V] {
	// How many levels we need?
	levels := t.Height / stride
	if t.Height%stride != 0 {
		levels++
	}

	tries := make([][]*Trie[V], levels)

	// Creates subtries of each level with builder
	// Builder is called on each root node of next level
	// and returns trie of that node

	var builder func(root *Node[V], level uint) *Trie[V]

	builder = func(root *Node[V], level uint) *Trie[V] {
		// Corrects root value
		if root.Value == nil {
			if val, found := t.Lookup(root.prefix + "*"); found {
				root.Value = &val
			}
		}

		// Trie of given root node
		sub := New[V]()
		sub.Prefix = root.prefix

		// Roots array of next level tries
		if level+1 < uint(len(tries)) && tries[level+1] == nil {
			tries[level+1] = make([]*Trie[V], 0)
		}

		// BFS queue
		queue := []*Node[V]{root}

		// BFS loop
		for len(queue) > 0 {
			n := queue[0]
			queue = queue[1:]

			if n.Value != nil {
				// Adds existing prefix into new trie
				sub.Add(n.prefix[len(sub.Prefix):]+"*", *n.Value)
			}

			if n.Right != nil {
				if n.Right.height >= stride*(level+1) {
					tries[level+1] = append(tries[level+1], builder(n.Right, level+1))
				} else {
					queue = append(queue, n.Right)
				}
			}

			if n.Left != nil {
				if n.Left.height >= stride*(level+1) {
					tries[level+1] = append(tries[level+1], builder(n.Left, level+1))
				} else {
					queue = append(queue, n.Left)
				}
			}
		}

		sub.Height = stride

		return sub
	}
	tries[0] = []*Trie[V]{
		builder(t.Root, 0),
	}

	return tries
}

// Add adds new route into trie.
// Given route must be in binary regex format e.g. *, 11*.
func (t *Trie[V]) Add(route string, value V) {
	it := t.Root

	for _, b := range route {
		switch b {
		case '*':
			it.Value = &value
		case '0':
			if it.Left == nil {
				it.Left = new(Node[V])
				it.Left.prefix = it.prefix + "0"
				it.Left.height = it.height + 1

				if t.Height < it.Left.height+1 {
					t.Height = it.Left.height + 1
				}
			}

			it = it.Left
		case '1':
			if it.Right == nil {
				it.Right = new(Node[V])
				it.Right.prefix = it.prefix + "1"
				it.Right.height = it.height + 1

				if t.Height < it.Right.height+1 {
					t.Height = it.Right.height + 1
				}
			}

			it = it.Right
		}
	}
}

// Lookup looks up given route in trie and returns found value.
// Given route must be in binary representation e.g. 111111..
func (t *Trie[V]) Lookup(route string) (V, bool) {
	it := t.Root

	var result V

	found := false

	for _, b := range route {
		if it.Value != nil {
			result = *it.Value
			found = true
		}

		// nolint: nestif
		if b == '1' {
			if it.Right != nil {
				it = it.Right
			} else {
				return result, found
			}
		} else {
			if it.Left != nil {
				it = it.Left
			} else {
				return result, found
			}
		}
	}

	if it.Value != nil {
		result = *it.Value
		found = true
	}

	return result, found
}

// maxToArrayHeight is the maximum trie height allowed for array conversion.
// 2^20 = ~1M nodes is a reasonable upper bound.
const maxToArrayHeight = 20

// ToArray returns node array of trie
// with following structure:
//
//	   i
//	  / \
//	2i  2i+1
//
// e.g.
//
//	      1
//	    /  \
//	   2    3
//	 /  \  / \
//	4   5 6   7
func (t *Trie[V]) ToArray() []Node[V] {
	if t.Height > maxToArrayHeight {
		panic(fmt.Sprintf("trie height %d exceeds maximum of %d for array conversion", t.Height, maxToArrayHeight))
	}

	nodes := make([]Node[V], 1<<t.Height)

	nodes[1] = *t.Root

	for i := 1; i < 1<<t.Height; i++ {
		c := nodes[i]
		if c.Left != nil {
			nodes[2*i] = *c.Left
		}

		if c.Right != nil {
			nodes[2*i+1] = *c.Right
		}
	}

	return nodes
}
