package lexicon

import (
	"bufio"
	"os"
	"strings"
)

// NewTrieNode returns a pointer to a new empty root TrieNode with an initialised map for NextNodes
func NewTrieNode() *TrieNode {
	return &TrieNode{
		Label:     "",
		Terminal:  false,
		NextNodes: make(map[rune]*TrieNode),
	}
}

// TrieNode is used for efficient prefix searches on a collection of strings
// The zero value may be used as a root node, however the caller is responsible
// for initalising the NextNodes map.
type TrieNode struct {
	Label     string
	Terminal  bool
	NextNodes map[rune]*TrieNode
}

// InsertWordsFromFile inserts words from a file which has a single word on each line. It is
// intended to be called on the root node.
func (t *TrieNode) InsertWordsFromFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t.Insert(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

// Insert inserts the provided word into the trie. It is intended to be called on the root node.
func (t *TrieNode) Insert(word string) {
	if t.Contains(word) {
		return
	}

	var stringBuilder strings.Builder
	currNode := t
	for _, char := range word {
		stringBuilder.WriteRune(char)
		if _, ok := currNode.NextNodes[char]; !ok {
			currNode.NextNodes[char] = &TrieNode{
				Label:     stringBuilder.String(),
				Terminal:  false,
				NextNodes: make(map[rune]*TrieNode),
			}
		}
		currNode = currNode.NextNodes[char]
	}
	currNode.Terminal = true
}

// Contains is used for identifying whether the provided word is in the trie rooted at t
func (t *TrieNode) Contains(word string) bool {
	currNode := t
	for _, char := range word {
		nextNode, ok := currNode.NextNodes[char]
		if !ok {
			return false
		}
		currNode = nextNode
	}
	return currNode.Terminal
}

// Delete removes the word from the trie. It is intended to be called on a root node.
func (t *TrieNode) Delete(word string) {
	if !t.Contains(word) {
		return
	}

	chars := []rune(word)
	currNode := t
	for i := 0; i < len(chars); currNode, i = currNode.NextNodes[chars[i]], i+1 {
		if currNode.Terminal {
			delete(currNode.NextNodes, chars[i])
			return
		}
	}

	if len(currNode.NextNodes) == 0 {
		delete(t.NextNodes, chars[0])
		return
	}

	currNode.Terminal = false
}

// ValidLettersBetweenPrefixAndSuffix returns the set of all letters '?'
// for which there is a word in the trie that looks like: '{prefix}?{suffix}'.
// It is inteded to be called on the root node.
func (t *TrieNode) ValidLettersBetweenPrefixAndSuffix(prefix, suffix string) map[rune]bool {

	validLetters := make(map[rune]bool)
	currNode := t
	prefixInTrie := true

	for _, prefixChar := range prefix {
		currNode, prefixInTrie = currNode.NextNodes[prefixChar]
		if !prefixInTrie {
			return validLetters
		}
	}

	middleNode := currNode

	for middleLetter, currNode := range middleNode.NextNodes {
		wordInTrie := true
		for _, suffixChar := range suffix {
			currNode, wordInTrie = currNode.NextNodes[suffixChar]
			if !wordInTrie {
				break
			}
		}
		if wordInTrie && currNode.Terminal {
			validLetters[middleLetter] = true
		}
	}
	return validLetters
}

// IsRoot returns true if the reciever is the root node of the trie
func (t *TrieNode) IsRoot() bool {
	return t.Label == ""
}

// IncomingEdge returns the edge which links this node with its parent. The zero value is returned
// if this node is the root.
func (t *TrieNode) IncomingEdge() rune {
	if t.IsRoot() {
		return 0
	}
	return []rune(t.Label)[len(t.Label)-1]
}

// EdgePruner is used for indicating which edges should be followed in a pruned traversal
type EdgePruner interface {
	IsValidEdge(edge rune) bool
	Terminate(node *TrieNode) bool
}

// Visitor is used for visiting elements in a collection
type Visitor interface {
	Visit(*TrieNode)
	Exit(*TrieNode)
}

// Visitor is used for visiting elements in a collection using a pruned traversal
type PrunerVisitor interface {
	EdgePruner
	Visitor
}

// VisitNodesWithPruning performs a pruned depth first traversal of the trie rooted at t.
func (t *TrieNode) VisitNodesWithPruning(prunerVisitor PrunerVisitor) {
	prunerVisitor.Visit(t)

	if !prunerVisitor.Terminate(t) {
		for edge, nextNode := range t.NextNodes {
			if prunerVisitor.IsValidEdge(edge) {
				nextNode.VisitNodesWithPruning(prunerVisitor)
			}
		}
	}

	prunerVisitor.Exit(t)
}
