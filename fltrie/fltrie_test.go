package fltrie

import (
	"fmt"
	"testing"

	"github.com/1995parham/FlashTrie.go/net"
)

func TestBasic(t *testing.T) {
	lookups := map[string]string{
		"192.168.73.10": "B",
	}

	fltrie := New()

	r1, _ := net.ParseNet("0.0.0.0/31")
	r2, _ := net.ParseNet("192.168.73.0/24")
	r3, _ := net.ParseNet("192.168.75.0/24")
	r4, _ := net.ParseNet("192.168.72.0/24")
	r5, _ := net.ParseNet("192.0.0.0/8")
	r6, _ := net.ParseNet("172.0.0.0/8")

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
		if fltrie.Lookup(net.ParseIP(route)) != nhi {
			t.Fatalf("Invalid lookup %s != %s", nhi, fltrie.Lookup(net.ParseIP(route)))
		}
	}
}
