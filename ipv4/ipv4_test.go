package ipv4_test

import (
	"net/netip"
	"testing"

	"github.com/1995parham/FlashTrie.go/ipv4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCIDR(t *testing.T) {
	t.Parallel()

	result, err := ipv4.ParseCIDR("192.0.0.0/4")
	require.NoError(t, err)
	assert.Equal(t, "1100*", result)
}

func TestParseCIDRFull(t *testing.T) {
	t.Parallel()

	result, err := ipv4.ParseCIDR("10.10.10.0/24")
	require.NoError(t, err)
	assert.Equal(t, "000010100000101000001010*", result)
}

func TestParseCIDRInvalid(t *testing.T) {
	t.Parallel()

	_, err := ipv4.ParseCIDR("not-a-cidr")
	require.Error(t, err)
}

func TestAdapterEncode(t *testing.T) {
	t.Parallel()

	adapter := ipv4.NewAdapter()
	result, err := adapter.Encode(netip.MustParseAddr("192.168.1.1"))
	require.NoError(t, err)
	assert.Len(t, result, 32)
}

func TestAdapterEncodeIPv6Rejected(t *testing.T) {
	t.Parallel()

	adapter := ipv4.NewAdapter()
	_, err := adapter.Encode(netip.MustParseAddr("::1"))
	require.Error(t, err)
}

func TestAdapterKeyBits(t *testing.T) {
	t.Parallel()

	adapter := ipv4.NewAdapter()
	assert.Equal(t, uint(32), adapter.KeyBits())
}

func BenchmarkParseCIDR(b *testing.B) {
	for b.Loop() {
		_, _ = ipv4.ParseCIDR("192.168.73.0/24")
	}
}

func BenchmarkAdapterEncode(b *testing.B) {
	adapter := ipv4.NewAdapter()
	addr := netip.MustParseAddr("192.168.73.10")

	b.ResetTimer()

	for b.Loop() {
		_, _ = adapter.Encode(addr)
	}
}
