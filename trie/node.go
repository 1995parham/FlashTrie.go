package trie

// Node represents trie node which may contains a prefix
// if node has prefix it will have non-nil prefix value and
// otherwise it will have nil prefix value.
type Node struct {
	prefix  string
	NextHop string
	Right   *Node
	Left    *Node
	height  uint
}
