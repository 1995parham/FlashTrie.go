package main

import (
	"os"
	"testing"

	"github.com/1995parham/FlashTrie.go/fltrie"
	"github.com/1995parham/FlashTrie.go/net"
	"github.com/1995parham/FlashTrie.go/pctrie"
	"github.com/1995parham/FlashTrie.go/trie"
	yaml "gopkg.in/yaml.v2"
)

func TestBasic(t *testing.T) {
	t.Parallel()

	r1, _ := net.ParseNet("192.0.2.1/4")
	r2, _ := net.ParseNet("192.0.2.1/8")
	r3, _ := net.ParseNet("172.0.2.1/8")

	trie := trie.New()
	trie.Add(r1, "A")
	trie.Add(r2, "B")
	trie.Add(r3, "C")

	pctrie := pctrie.New(trie, 4)

	ipAddress := []string{
		"172.0.1.1",
		"192.0.1.1",
		"192.0.0.0",
		"172.73.72.75",
		"194.0.0.0",
	}

	for _, ip := range ipAddress {
		if trie.Lookup(net.ParseIP(ip)) != pctrie.Lookup(net.ParseIP(ip)) {
			t.Fatalf("Invalid route %s: %s != %s", ip, trie.Lookup(net.ParseIP(ip)), pctrie.Lookup(net.ParseIP(ip)))
		}

		t.Logf("%s: %s", ip, trie.Lookup(net.ParseIP(ip)))
	}
}

func TestFarkiani(t *testing.T) {
	t.Parallel()

	testRoutes := []route{
		{
			Route:   "1.2.3.4",
			Nexthop: "Raha",
		},
		{
			Route:   "10.10.10.194",
			Nexthop: "6.6.6.6",
		},
		{
			Route:   "10.10.10.2",
			Nexthop: "5.5.5.5",
		},
		{
			Route:   "218.144.10.10",
			Nexthop: "209.244.2.115",
		},
	}

	f := "T1.yml"

	var routes []route

	data, err := os.ReadFile(f)
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
		r, err := net.ParseNet(route.Route)
		if err != nil {
			t.Fatalf("Parsing file %s failed with: %s\n", f, err)
		}

		fltrie.Add(r, route.Nexthop)
	}

	if err := fltrie.Build(); err != nil {
		t.Fatalf("Building flash trie failed with: %s\n", err)
	}

	for _, r := range testRoutes {
		if l := fltrie.Lookup(net.ParseIP(r.Route)); l != r.Nexthop {
			t.Fatalf("%s -> %s is not equal to %s", r.Route, l, r.Nexthop)
		}
	}
}
