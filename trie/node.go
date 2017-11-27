/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 27-11-2017
 * |
 * | File Name:     node.go
 * +===============================================
 */

package trie

// Node represents trie node wich may contains a prefix
// if node has prefix it will have non-nil prefix value and
// otherwise it will have nil prefix value
type Node struct {
	Prefix  string
	NextHop string
	Right   *Node
	Left    *Node
}
