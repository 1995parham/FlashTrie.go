/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 27-11-2017
 * |
 * | File Name:     trie/trie.go
 * +===============================================
 */

package trie

// Trie represents binary trie for IP route lookup
type Trie struct {
	Root   *Node
	Height uint
}

// New creates new trie
func New() *Trie {
	return &Trie{
		Root:   new(Node),
		Height: 1,
	}
}

// NewFromArray creates new trie based on given node array
// with following structure
//       i
//      / \
//    2i  2i+1
// TODO
func NewFromArray(nodes []Node) *Trie {
	for i := 0; i < len(nodes); i++ {
	}

	return nil
}

// Divide divides the binary trie into different levels
// based on these k-bit subtries. If k = 4 thus, level 0 contains
// from prefix length 0 to prefix length 7, and so on.
// Each level contains one or more subtries.
func (t *Trie) Divide(stride uint) [][]*Trie {
	// How many levels we need?
	levels := t.Height / stride
	if t.Height%stride != 0 {
		levels++
	}

	tries := make([][]*Trie, levels)

	// Creates subtries of each level with buidler

	var builder func(root *Node, level uint)
	builder = func(root *Node, level uint) {
		it := root
		for it.height >= stride*level && it.height <= stride*(level+1) {
		}
	}
	builder(t.Root, 1)

	return tries
}

// Add adds new route into trie
// given route must be in binary regex format e.g. *, 11*
func (t *Trie) Add(route string, nexthop string) {
	it := t.Root
	for _, b := range route {
		if b == '*' {
			it.Prefix = route
			it.NextHop = nexthop
		} else {
			if b == '0' {
				if it.Left == nil {
					it.Left = new(Node)
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
// given route must be in binary represenation e.g. 111111..
func (t *Trie) Lookup(route string) string {
	it := t.Root
	nexthop := "-"
	for _, b := range route {
		if it.Prefix != "" {
			nexthop = it.NextHop
		}

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

	return nexthop
}

// ToArray returns node array of trie
// with following structure:
//       i
//      / \
//    2i  2i+1
//
// e.g.
//          1
//        /  \
//       2    3
//     /  \  / \
//    4   5 6   7
//
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
