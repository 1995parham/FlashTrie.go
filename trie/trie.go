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

// Trie represents binary trie for route lookup
type Trie struct {
	Root   *Node
	Stride uint
}

// New creates new trie
func New() *Trie {
	return &Trie{
		Root:   new(Node),
		Stride: 1,
	}
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
					if t.Stride < it.Left.height+1 {
						t.Stride = it.Left.height + 1
					}
				}
				it = it.Left
			}
			if b == '1' {
				if it.Right == nil {
					it.Right = new(Node)
					it.Right.height = it.height + 1
					if t.Stride < it.Right.height+1 {
						t.Stride = it.Right.height + 1
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
func (t *Trie) ToArray() []Node {
	nodes := make([]Node, 1<<t.Stride)

	nodes[1] = *t.Root

	for i := 1; i < 1<<t.Stride; i++ {
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
