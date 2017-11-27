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
	Root *Node
}

// New creates new trie
func New() Trie {
	return Trie{
		Root: new(Node),
	}
}

// Add adds new route into trie
// given route must be in binary regex format e.g. *, 11*
func (t Trie) Add(route string, nexthop string) {
	it := t.Root
	for _, b := range route {
		if b == '*' {
			it.Prefix = route
			it.NextHop = nexthop
		} else {
			if b == '1' {
				if it.Left == nil {
					it.Left = new(Node)
				}
				it = it.Left
			}
			if b == '0' {
				if it.Right == nil {
					it.Right = new(Node)
				}
				it = it.Right
			}
		}
	}
}

// Lookup lookups given route in tire and returns finded nexhop or -
// given route must be in binary represenation e.g. 111111..
func (t Trie) Lookup(route string) string {
	it := t.Root
	nexthop := "-"
	for _, b := range route {
		if it.Prefix != "" {
			nexthop = it.NextHop
		}

		if b == '0' {
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
