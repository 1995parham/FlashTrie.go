package fltrie_test

import (
	"maps"
	"net/netip"
	"testing"

	"github.com/1995parham/FlashTrie.go/fltrie"
	"github.com/1995parham/FlashTrie.go/ipv4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildTestFLTrie(t *testing.T) *fltrie.FLTrie[netip.Addr, string] {
	t.Helper()

	fl := fltrie.New[netip.Addr, string](ipv4.NewAdapter(), ipv4.DefaultConfig())

	routes := []struct{ cidr, nh string }{
		{"0.0.0.0/31", "A"},
		{"192.168.73.0/24", "B"},
		{"192.168.75.0/24", "C"},
		{"192.168.72.0/24", "D"},
		{"192.0.0.0/8", "E"},
		{"172.0.0.0/8", "F"},
	}

	for _, r := range routes {
		parsed, err := ipv4.ParseCIDR(r.cidr)
		require.NoError(t, err)
		require.NoError(t, fl.Add(parsed, r.nh))
	}

	require.NoError(t, fl.Build())

	return fl
}

func TestBasic(t *testing.T) {
	t.Parallel()

	fl := buildTestFLTrie(t)

	result, found, err := fl.Lookup(netip.MustParseAddr("192.168.73.10"))
	require.NoError(t, err)
	require.True(t, found)
	assert.Equal(t, "B", result)
}

func TestAddAfterBuild(t *testing.T) {
	t.Parallel()

	fl := buildTestFLTrie(t)

	err := fl.Add("1*", "X")
	require.ErrorIs(t, err, fltrie.ErrAlreadyBuilt)
}

func TestBuildWrongHeight(t *testing.T) {
	t.Parallel()

	fl := fltrie.New[netip.Addr, string](ipv4.NewAdapter(), ipv4.DefaultConfig())
	require.NoError(t, fl.Add("*", "A"))

	err := fl.Build()
	require.ErrorIs(t, err, fltrie.ErrInvalidHeight)
}

func TestBuildTwice(t *testing.T) {
	t.Parallel()

	fl := buildTestFLTrie(t)

	err := fl.Build()
	require.ErrorIs(t, err, fltrie.ErrAlreadyBuilt)
}

func TestAll(t *testing.T) {
	t.Parallel()

	fl := buildTestFLTrie(t)

	got := maps.Collect(fl.All())

	// 6 explicit routes + the trie root inherits a value during construction
	assert.GreaterOrEqual(t, len(got), 6, "should iterate over at least the 6 added routes")
}

func BenchmarkBuild(b *testing.B) {
	routes := []struct{ cidr, nh string }{
		{"0.0.0.0/31", "A"},
		{"192.168.73.0/24", "B"},
		{"192.168.75.0/24", "C"},
		{"192.168.72.0/24", "D"},
		{"192.0.0.0/8", "E"},
		{"172.0.0.0/8", "F"},
	}

	parsedRoutes := make([]struct{ route, nh string }, len(routes))

	for i, r := range routes {
		parsed, err := ipv4.ParseCIDR(r.cidr)
		if err != nil {
			b.Fatal(err)
		}

		parsedRoutes[i] = struct{ route, nh string }{parsed, r.nh}
	}

	b.ResetTimer()

	for b.Loop() {
		fl := fltrie.New[netip.Addr, string](ipv4.NewAdapter(), ipv4.DefaultConfig())
		for _, r := range parsedRoutes {
			_ = fl.Add(r.route, r.nh)
		}

		_ = fl.Build()
	}
}

func BenchmarkLookup(b *testing.B) {
	fl := fltrie.New[netip.Addr, string](ipv4.NewAdapter(), ipv4.DefaultConfig())

	routes := []struct{ cidr, nh string }{
		{"0.0.0.0/31", "A"},
		{"192.168.73.0/24", "B"},
		{"192.168.75.0/24", "C"},
		{"192.168.72.0/24", "D"},
		{"192.0.0.0/8", "E"},
		{"172.0.0.0/8", "F"},
	}

	for _, r := range routes {
		parsed, err := ipv4.ParseCIDR(r.cidr)
		if err != nil {
			b.Fatal(err)
		}

		_ = fl.Add(parsed, r.nh)
	}

	if err := fl.Build(); err != nil {
		b.Fatal(err)
	}

	addr := netip.MustParseAddr("192.168.73.10")

	b.ResetTimer()

	for b.Loop() {
		_, _, _ = fl.Lookup(addr)
	}
}
