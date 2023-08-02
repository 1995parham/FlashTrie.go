package trie

import (
	dll "github.com/emirpasic/gods/lists/doublylinkedlist"
)

// Trie represents binary trie for IP route lookup.
type Trie struct {
	Root   *Node
	Height uint
	Prefix string // when trie is subtrie this field represents old prefix of root
}

// New creates new trie.
func New() *Trie {
	return &Trie{
		Root:   new(Node),
		Height: 1,
		Prefix: "",
	}
}

// NewFromArray creates new trie based on given node array
// with following structure.
//
//	   i
//	  / \
//	2i  2i+1
func NewFromArray(_ []Node) *Trie {
	panic("implement me!")
}

// Divide divides the binary trie into different levels
// based on these k-bit subtries. If k = 4 thus, level 0 contains
// from prefix length 0 to prefix length 7, and so on.
// Each level contains one or more subtries.
// nolint: funlen, gocognit, cyclop
func (t *Trie) Divide(stride uint) [][]*Trie {
	// How many levels we need?
	levels := t.Height / stride
	if t.Height%stride != 0 {
		levels++
	}

	tries := make([][]*Trie, levels)

	// Creates subtries of each level with builder
	// Builder is called on each root node of next level
	// and returns trie of that node

	var builder func(root *Node, level uint) *Trie
	builder = func(root *Node, level uint) *Trie {
		// Corrects root next hop
		if root.NextHop == "" {
			root.NextHop = t.Lookup(root.prefix + "*")
		}

		// Trie of given root node
		t := New()
		t.Prefix = root.prefix

		// Roots array of next level tries
		if int(level+1) < len(tries) && tries[level+1] == nil {
			tries[level+1] = make([]*Trie, 0)
		}

		// BFS queue
		q := dll.New()

		// BFS initiation
		q.Add(root)

		// BFS loop
		for !q.Empty() {
			i, _ := q.Get(0)

			n, ok := i.(*Node)
			if !ok {
				panic("")
			}

			if n.NextHop != "" {
				// Adds existing prefix into new trie
				t.Add(n.prefix[len(t.Prefix):]+"*", n.NextHop)
			}

			if n.Right != nil {
				if n.Right.height >= stride*(level+1) {
					tries[level+1] = append(tries[level+1], builder(n.Right, level+1))
				} else {
					q.Add(n.Right)
				}
			}

			if n.Left != nil {
				if n.Left.height >= stride*(level+1) {
					tries[level+1] = append(tries[level+1], builder(n.Left, level+1))
				} else {
					q.Add(n.Left)
				}
			}

			q.Remove(0)
		}

		t.Height = stride

		return t
	}
	tries[0] = []*Trie{
		builder(t.Root, 0),
	}

	return tries
}

// Add adds new route into trie
// given route must be in binary regex format e.g. *, 11*.
func (t *Trie) Add(route string, nexthop string) {
	it := t.Root

	// nolint: nestif
	for _, b := range route {
		if b == '*' {
			it.NextHop = nexthop
		} else {
			if b == '0' {
				if it.Left == nil {
					it.Left = new(Node)
					it.Left.prefix = it.prefix + "0"
					it.Left.height = it.height + 1
					if t.Height < it.Left.height+1 {
						t.Height = it.Left.height + 1
					}
				}
				it = it.Left
			}
			if b == '1' {
				if it.Right == nil {
					it.Right = new(Node)
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
}

// Lookup lookups given route in tire and returns finded nexhop or -
// given route must be in binary representation e.g. 111111..
func (t *Trie) Lookup(route string) string {
	it := t.Root
	nexthop := "-"

	for _, b := range route {
		if it.NextHop != "" {
			nexthop = it.NextHop
		}

		// nolint: nestif
		if b == '1' {
			if it.Right != nil {
				it = it.Right
			} else {
				return nexthop
			}
		} else {
			if it.Left != nil {
				it = it.Left
			} else {
				return nexthop
			}
		}
	}

	if it.NextHop != "" {
		nexthop = it.NextHop
	}

	return nexthop
}

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
func (t *Trie) ToArray() []Node {
	nodes := make([]Node, 1<<t.Height)

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
