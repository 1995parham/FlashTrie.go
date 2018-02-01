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
	"io/ioutil"
	"testing"

	yaml "gopkg.in/yaml.v2"

	"github.com/AUTProjects/FlashTrie.go/fltrie"
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

var testRoutes []route = []route{
	route{
		Route:   "1.2.3.4",
		Nexthop: "Kiana",
	},
	route{
		Route:   "10.10.10.194",
		Nexthop: "6.6.6.6",
	},
	route{
		Route:   "10.10.10.2",
		Nexthop: "5.5.5.5",
	},
	route{
		Route:   "218.144.10.10",
		Nexthop: "209.244.2.115",
	},
}

func TestFarkiani(t *testing.T) {
	f := "T1.yml"
	var routes []route

	data, err := ioutil.ReadFile(f)
	if err != nil {
		t.Fatalf("Reading file %s failed with: %s\n", f, err)
	}

	err = yaml.Unmarshal(data, &routes)
	if err != nil {
		t.Fatalf("Parsing file %s failed with: %s\n", f, err)
	}

	// Building flash trie

	fltrie := fltrie.New()

	for _, route := range routes {
		r, err := util.ParseNet(route.Route)
		if err != nil {
			t.Fatalf("Parsing file %s failed with: %s\n", f, err)
		}
		fltrie.Add(r, route.Nexthop)
	}

	if err := fltrie.Build(); err != nil {
		t.Fatalf("Building flash trie failed with: %s\n", err)
	}

	for _, r := range testRoutes {
		if l := fltrie.Lookup(util.ParseIP(r.Route)); l != r.Nexthop {
			t.Fatalf("%s -> %s is not equal to %s", r.Route, l, r.Nexthop)
		}
	}
}
