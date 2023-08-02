package fltrie

import (
	"fmt"

	"github.com/1995parham/FlashTrie.go/pctrie"
	"github.com/1995parham/FlashTrie.go/trie"
)

type hashElement struct {
	pctrie  *pctrie.PCTrie
	nextHop string
}

// FLTrie represents flash trie structure
type FLTrie struct {
	trie    *trie.Trie
	pctries []map[string]*hashElement
	build   bool
}

// New creates empty and un-build flash trie
func New() *FLTrie {
	return &FLTrie{
		trie:    trie.New(),
		pctries: make([]map[string]*hashElement, 2),
	}
}

// Add adds new route into unbuild flash trie
// given route must be in binary regex format e.g. *, 11*
func (fl *FLTrie) Add(route string, nexthop string) {
	if !fl.build {
		fl.trie.Add(route, nexthop)
	}
}

// Build builds flash trie three level hierarchy
func (fl *FLTrie) Build() error {
	if fl.trie.Height != 32 {
		return fmt.Errorf("FLTrie height must be 32 that are greater than %d", fl.trie.Height)
	}

	if !fl.build {
		fl.build = true

		tries := fl.trie.Divide(8)

		// Level 2 pctries
		fl.pctries[0] = make(map[string]*hashElement)
		for _, trie := range tries[2] {
			fl.pctries[0][trie.Prefix] = &hashElement{
				pctrie:  pctrie.New(trie, 2),
				nextHop: trie.Root.NextHop,
			}
		}

		// Level 3 pctries
		fl.pctries[1] = make(map[string]*hashElement)
		for _, trie := range tries[3] {
			fl.pctries[1][trie.Prefix] = &hashElement{
				pctrie:  pctrie.New(trie, 2),
				nextHop: trie.Root.NextHop,
			}
		}

		return nil
	}

	return fmt.Errorf("FLTrie is build already")
}

// Lookup lookups given route and returns found nexhop or -
// given route must be in binary representation e.g. 111111..
func (fl *FLTrie) Lookup(route string) string {
	nh := "-"

	if len(route) == 32 && fl.build {
		// Level 1 (16 bit trie)
		if nhi := fl.trie.Lookup(route[:16]); nhi != "-" {
			nh = nhi
		}
		// Level 2 (25 bit - 8 bit pctrie)
		if he, ok := fl.pctries[0][route[:16]]; ok {
			if nhi := he.pctrie.Lookup(route[16:24]); nhi != "-" {
				nh = nhi
			} else {
				nh = he.nextHop
			}
		}
		// Level 3 (32 bit - 8 bit pctrie)
		if he, ok := fl.pctries[1][route[:24]]; ok {
			if nhi := he.pctrie.Lookup(route[24:32]); nhi != "-" {
				nh = nhi
			} else {
				nh = he.nextHop
			}
		}
	}

	return nh
}
