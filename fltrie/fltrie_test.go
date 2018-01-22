/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 22-01-2018
 * |
 * | File Name:     fltrie/fltrie_test.go
 * +===============================================
 */

package fltrie

import (
	"fmt"
	"testing"

	"github.com/AUTProjects/FlashTrie.go/util"
)

var lookups = map[string]string{
	"192.168.73.10": "B",
}

func TestBasic(t *testing.T) {
	fltrie := New()

	r1, _ := util.ParseNet("0.0.0.0/31")
	r2, _ := util.ParseNet("192.168.73.0/24")
	r3, _ := util.ParseNet("192.168.75.0/24")
	r4, _ := util.ParseNet("192.168.72.0/24")
	r5, _ := util.ParseNet("192.0.0.0/8")
	r6, _ := util.ParseNet("172.0.0.0/8")

	fmt.Println(r2)

	fltrie.Add(r1, "A")
	fltrie.Add(r2, "B")
	fltrie.Add(r3, "C")
	fltrie.Add(r4, "D")
	fltrie.Add(r5, "E")
	fltrie.Add(r6, "F")

	if err := fltrie.Build(); err != nil {
		t.Fatal(err)
	}

	for route, nhi := range lookups {
		if fltrie.Lookup(util.ParseIP(route)) != nhi {
			t.Fatalf("Invalid lookup %s != %s", nhi, fltrie.Lookup(util.ParseIP(route)))
		}
	}
}
