package trie

import (
	"strings"
)

// Node is a node in the Trie data structure.
type Node struct {
	Label     string
	Terminal  bool
	NextNodes map[rune]*Node
}

// Trie is a data structure which contains strings.
// Checking whether a word is present in the trie is an O(k) operation
// where k is the length of the string.
type Trie struct {
	Root *Node
}

// Insert inserts the word into the trie.
// It returns a bool representing whether the word was newly inserted.
func (trie *Trie) Insert(word string) bool {
	if trie.Contains(word) {
		return false
	}
	var stringBuilder strings.Builder
	node := trie.Root
	for _, char := range word {
		stringBuilder.WriteRune(char)
		if nextNode, ok := node.NextNodes[char]; ok {
			node = nextNode
		} else {
			node.NextNodes[char] = &Node{
				Label:     stringBuilder.String(),
				Terminal:  false,
				NextNodes: map[rune]*Node{},
			}
			node = node.NextNodes[char]
		}
	}
	node.Terminal = true
	return true
}

// Contains identifies whether the word is in the trie.
// It returns a bool representing whether the word is in the tree.
func (trie *Trie) Contains(word string) bool {
	node := trie.Root
	for _, char := range word {
		nextNode, ok := node.NextNodes[char]
		if !ok {
			return false
		}
		node = nextNode
	}
	return node.Terminal
}

// Delete removes the word from the trie, if it is present in the tree.
// It returns whether the word was present in the tree.
func (trie *Trie) Delete(word string) bool {
	if !trie.Contains(word) {
		return false
	}

	var prefixWordTerminalNode *Node
	var suffixInitialChar rune 
	currNode := trie.Root

	for i, char := range word {
		if currNode.Terminal && i != len(word) {
			prefixWordTerminalNode = currNode
			suffixInitialChar = char
		}
		currNode = currNode.NextNodes[char]
	}

	if prefixWordTerminalNode == nil {
		if len(currNode.NextNodes) == 0 {
			delete(
				trie.Root.NextNodes,
				[]rune(word)[0],
			)
		} else {
			currNode.Terminal = false
		}
	} else {
		delete(prefixWordTerminalNode.NextNodes, suffixInitialChar)
	}
	return true
}

// New returns a pointer to a new empty Trie.
func New() *Trie {
	return &Trie{
		Root: &Node{
			Label: "",
			Terminal: false,
			NextNodes: map[rune]*Node{},
		},
	}
}
