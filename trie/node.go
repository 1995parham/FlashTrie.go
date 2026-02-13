package trie

// Node represents trie node which may contains a prefix
// if node has prefix it will have non-nil prefix value and
// otherwise it will have nil prefix value.
type Node[V any] struct {
	prefix string
	Value  *V
	Right  *Node[V]
	Left   *Node[V]
	height uint
}
