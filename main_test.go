/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 21-01-2018
 * |
 * | File Name:     main_test.go
 * +===============================================
 */
package main

import (
	"testing"

	"github.com/AUTProjects/FlashTrie.go/pctrie"
	"github.com/AUTProjects/FlashTrie.go/trie"
)

var ipAddress = []string{
	"172.0.1.1",
	"192.0.1.1",
	"192.0.0.0",
	"172.73.72.75",
	"194.0.0.0",
}

func TestBasic(t *testing.T) {
	r1 := parseCIDR("192.0.2.1/4")
	r2 := parseCIDR("192.0.2.1/8")
	r3 := parseCIDR("172.0.2.1/8")

	trie := trie.New()
	trie.Add(r1, "A")
	trie.Add(r2, "B")
	trie.Add(r3, "C")

	pctrie := pctrie.New(trie, 4)

	for _, ip := range ipAddress {
		if trie.Lookup(parseIP(ip)) != pctrie.Lookup(parseIP(ip)) {
			t.Fatalf("Invalid route %s: %s != %s", ip, trie.Lookup(parseIP(ip)), pctrie.Lookup(parseIP(ip)))
		}
		t.Logf("%s: %s", ip, trie.Lookup(parseIP(ip)))
	}
}
