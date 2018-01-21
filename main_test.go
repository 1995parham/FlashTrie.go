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
	"github.com/AUTProjects/FlashTrie.go/util"
)

var ipAddress = []string{
	"172.0.1.1",
	"192.0.1.1",
	"192.0.0.0",
	"172.73.72.75",
	"194.0.0.0",
}

func TestBasic(t *testing.T) {
	r1, _ := util.ParseNet("192.0.2.1/4")
	r2, _ := util.ParseNet("192.0.2.1/8")
	r3, _ := util.ParseNet("172.0.2.1/8")

	trie := trie.New()
	trie.Add(r1, "A")
	trie.Add(r2, "B")
	trie.Add(r3, "C")

	pctrie := pctrie.New(trie, 4)

	for _, ip := range ipAddress {
		if trie.Lookup(util.ParseIP(ip)) != pctrie.Lookup(util.ParseIP(ip)) {
			t.Fatalf("Invalid route %s: %s != %s", ip, trie.Lookup(util.ParseIP(ip)), pctrie.Lookup(util.ParseIP(ip)))
		}
		t.Logf("%s: %s", ip, trie.Lookup(util.ParseIP(ip)))
	}
}
