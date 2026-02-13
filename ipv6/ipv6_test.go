package ipv6_test

import (
	"net/netip"
	"testing"

	"github.com/1995parham/FlashTrie.go/ipv6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCIDR(t *testing.T) {
	t.Parallel()

	result, err := ipv6.ParseCIDR("2001:db8::/32")
	require.NoError(t, err)
	assert.Len(t, result, 33) // 32 bits + "*"
	assert.Equal(t, "00100000000000010000110110111000*", result)
}

func TestParseCIDRLinkLocal(t *testing.T) {
	t.Parallel()

	result, err := ipv6.ParseCIDR("fe80::1/10")
	require.NoError(t, err)
	assert.Len(t, result, 11) // 10 bits + "*"
	assert.Equal(t, "1111111010*", result)
}

func TestParseCIDRInvalid(t *testing.T) {
	t.Parallel()

	_, err := ipv6.ParseCIDR("not-a-cidr")
	require.Error(t, err)
}

func TestParseCIDRIPv4Rejected(t *testing.T) {
	t.Parallel()

	_, err := ipv6.ParseCIDR("192.168.0.0/16")
	require.Error(t, err)
}

func TestParseCIDRIPv4MappedRejected(t *testing.T) {
	t.Parallel()

	_, err := ipv6.ParseCIDR("::ffff:192.168.0.0/112")
	require.Error(t, err)
}

func TestAdapterEncode(t *testing.T) {
	t.Parallel()

	adapter := ipv6.NewAdapter()
	result, err := adapter.Encode(netip.MustParseAddr("2001:db8::1"))
	require.NoError(t, err)
	assert.Len(t, result, 128)
	// First 32 bits should match 2001:0db8
	assert.Equal(t, "00100000000000010000110110111000", result[:32])
}

func TestAdapterEncodeIPv4Rejected(t *testing.T) {
	t.Parallel()

	adapter := ipv6.NewAdapter()
	_, err := adapter.Encode(netip.MustParseAddr("192.168.1.1"))
	require.Error(t, err)
}

func TestAdapterKeyBits(t *testing.T) {
	t.Parallel()

	adapter := ipv6.NewAdapter()
	assert.Equal(t, uint(128), adapter.KeyBits())
}

func BenchmarkParseCIDR(b *testing.B) {
	for b.Loop() {
		_, _ = ipv6.ParseCIDR("2001:db8::/32")
	}
}

func BenchmarkAdapterEncode(b *testing.B) {
	adapter := ipv6.NewAdapter()
	addr := netip.MustParseAddr("2001:db8::1")

	b.ResetTimer()

	for b.Loop() {
		_, _ = adapter.Encode(addr)
	}
}
