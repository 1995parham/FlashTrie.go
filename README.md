<h1 align="center">
    FlashTrie.go
</h1>

<p align="center">
    Generic, high-speed longest-prefix matching based on hash-based prefix-compressed tries
</p>

<p align="center">
    <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/1995parham/FlashTrie.go/test.yaml?style=for-the-badge" />
    <img alt="Codecov" src="https://img.shields.io/codecov/c/github/1995parham/FlashTrie.go?logo=codecov&style=for-the-badge">
</p>

## Introduction

FlashTrie.go is a Go implementation of the FlashTrie algorithm for longest-prefix matching at speeds beyond 100 Gb/s.
The core data structure is generic over key type: any domain that can be encoded into fixed-length binary strings
can be used with FlashTrie through the `Adapter[K]` interface.

Three adapters are provided out of the box:

| Package   | Key type      | Bits | Use case                    |
|-----------|---------------|------|-----------------------------|
| `ipv4`    | `netip.Addr`  | 32   | IPv4 route lookup           |
| `ipv6`    | `netip.Addr`  | 128  | IPv6 route lookup           |
| `geohash` | `geohash.Coord` | 60 | Geographic region matching  |

## Architecture

```
fltrie.FLTrie[K, V]
    |
    |-- Adapter[K]        (encodes keys to binary strings)
    |-- Config             (KeyBits, Stride, CompSize, TrieDepth)
    |
    |-- trie.Trie[V]      (binary trie, used during build)
    |-- pctrie.PCTrie[V]  (hash-based prefix-compressed trie, used at lookup)
```

The build phase divides the binary trie into stride-sized levels.
The first `TrieDepth` levels are kept as a direct trie for O(1) access;
the remaining levels are compressed into path-compressed tries indexed by hash maps.

## Usage

### IPv4

```go
fl := fltrie.New[netip.Addr, string](ipv4.NewAdapter(), ipv4.DefaultConfig())

route, _ := ipv4.ParseCIDR("192.168.73.0/24")
fl.Add(route, "next-hop-B")
fl.Build()

val, found, _ := fl.Lookup(netip.MustParseAddr("192.168.73.10"))
// found=true, val="next-hop-B"
```

### IPv6

```go
fl := fltrie.New[netip.Addr, string](ipv6.NewAdapter(), ipv6.DefaultConfig())

route, _ := ipv6.ParseCIDR("2001:db8::/32")
fl.Add(route, "next-hop-A")
fl.Build()

val, found, _ := fl.Lookup(netip.MustParseAddr("2001:db8::1"))
```

### Geohash

```go
adapter := geohash.NewAdapter(12) // 12-char precision (60 bits)
fl := fltrie.New[geohash.Coord, string](adapter, geohash.DefaultConfig())

prefix, _ := geohash.ParseGeohash("u4pru") // Aalborg, Denmark region
fl.Add(prefix, "Aalborg")
fl.Build()

val, found, _ := fl.Lookup(geohash.Coord{Lat: 57.64911, Lng: 10.40744})
```

### Custom adapter

Implement `fltrie.Adapter[K]` for any key type:

```go
type Adapter[K any] interface {
    Encode(key K) (string, error) // returns a binary string of length KeyBits
    KeyBits() uint
}
```

## Packages

| Package    | Description                                           |
|------------|-------------------------------------------------------|
| `fltrie`   | Core FlashTrie: `New`, `Add`, `Build`, `Lookup`       |
| `trie`     | Binary trie with divide-by-stride support             |
| `pctrie`   | Hash-based prefix-compressed trie                     |
| `ipv4`     | IPv4 adapter and CIDR parser                          |
| `ipv6`     | IPv6 adapter and CIDR parser (rejects IPv4-mapped)    |
| `geohash`  | Geohash adapter with configurable precision           |

## CLI

A small interactive tool is included for IPv4 route lookup:

```bash
go run ./cmd/flashtrie
```

It reads routes from a YAML file (see `sample.yml`) and provides a REPL for lookups.

## Reference

Bando, M., Lin, Y. and Chao, H. (2012). FlashTrie: Beyond 100-Gb/s IP Route Lookup Using Hash-Based Prefix-Compressed Trie. IEEE/ACM Transactions on Networking, 20(4), pp.1262-1275.
