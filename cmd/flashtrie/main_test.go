package main

import (
	"net/netip"
	"os"
	"testing"

	"github.com/1995parham/FlashTrie.go/fltrie"
	"github.com/1995parham/FlashTrie.go/ipv4"
	"github.com/1995parham/FlashTrie.go/pctrie"
	"github.com/1995parham/FlashTrie.go/trie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

func TestBasic(t *testing.T) {
	t.Parallel()

	r1, err := ipv4.ParseCIDR("192.0.2.1/4")
	require.NoError(t, err)

	r2, err := ipv4.ParseCIDR("192.0.2.1/8")
	require.NoError(t, err)

	r3, err := ipv4.ParseCIDR("172.0.2.1/8")
	require.NoError(t, err)

	tr := trie.New[string]()
	tr.Add(r1, "A")
	tr.Add(r2, "B")
	tr.Add(r3, "C")

	pc := pctrie.New[string](tr, 4)

	adapter := ipv4.NewAdapter()

	ipAddresses := []string{
		"172.0.1.1",
		"192.0.1.1",
		"192.0.0.0",
		"172.73.72.75",
		"194.0.0.0",
	}

	for _, ip := range ipAddresses {
		addr := netip.MustParseAddr(ip)

		parsed, err := adapter.Encode(addr)
		require.NoError(t, err)

		trieLookup, trieFound := tr.Lookup(parsed)

		pcLookup, pcFound, err := pc.Lookup(parsed)
		require.NoError(t, err)

		assert.Equal(t, trieFound, pcFound, "Found mismatch for %s", ip)

		if trieFound && pcFound {
			assert.Equal(t, trieLookup, pcLookup, "Mismatch for %s", ip)
		}

		t.Logf("%s: found=%v val=%s", ip, trieFound, trieLookup)
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

	data, err := os.ReadFile("../../testdata/T1.yml")
	require.NoError(t, err)

	var routes []route

	err = yaml.Unmarshal(data, &routes)
	require.NoError(t, err)

	// Building flash trie
	fl := fltrie.New[netip.Addr, string](ipv4.NewAdapter(), ipv4.DefaultConfig())

	for _, route := range routes {
		r, err := ipv4.ParseCIDR(route.Route)
		require.NoError(t, err)
		require.NoError(t, fl.Add(r, route.Nexthop))
	}

	require.NoError(t, fl.Build())

	for _, r := range testRoutes {
		addr := netip.MustParseAddr(r.Route)

		result, found, err := fl.Lookup(addr)
		require.NoError(t, err)
		require.True(t, found, "Expected to find %s", r.Route)
		assert.Equal(t, r.Nexthop, result, "%s -> %s is not equal to %s", r.Route, result, r.Nexthop)
	}
}
