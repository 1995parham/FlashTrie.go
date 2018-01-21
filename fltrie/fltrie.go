/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 22-01-2018
 * |
 * | File Name:     fltrie/fltrie.go
 * +===============================================
 */

package fltrie

import (
	"fmt"

	"github.com/AUTProjects/FlashTrie.go/pctrie"
	"github.com/AUTProjects/FlashTrie.go/trie"
)

// FLTrie represents flash trie structure
type FLTrie struct {
	trie    *trie.Trie
	pctries []map[string]*pctrie.PCTrie
	build   bool
}

// New creates empty and unbuild flash trie
func New() *FLTrie {
	return &FLTrie{
		trie:    trie.New(),
		pctries: make([]map[string]*pctrie.PCTrie, 2),
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
		fl.pctries[0] = make(map[string]*pctrie.PCTrie)
		for _, trie := range tries[2] {
			fl.pctries[0][trie.Prefix] = pctrie.New(trie, 2)
		}
		// Level 3 pctries
		fl.pctries[1] = make(map[string]*pctrie.PCTrie)
		for _, trie := range tries[3] {
			fl.pctries[1][trie.Prefix] = pctrie.New(trie, 2)
		}
		return nil
	}
	return fmt.Errorf("FLTrie is build already")
}
