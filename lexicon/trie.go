package lexicon

import (
	"bufio"
	"os"
	"strings"
)

// TrieNode is a data structure used for efficient prefix searches
type TrieNode struct {
	Label     string
	Terminal  bool
	NextNodes map[rune]*TrieNode
}

// CreateTrieFromFile builds a trie from a file which has a single word on each
// line.
func CreateTrieFromFile(filePath string) *TrieNode {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	trie := New()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		trie.Insert(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return trie
}

// Insert inserts the word into the trie. It returns a bool representing whether
// the word was newly inserted.
func (n *TrieNode) Insert(word string) bool {
	if n.Contains(word) {
		return false
	}
	var stringBuilder strings.Builder
	currNode := n
	for _, char := range word {
		stringBuilder.WriteRune(char)
		if _, ok := currNode.NextNodes[char]; !ok {
			currNode.NextNodes[char] = &TrieNode{
				Label:     stringBuilder.String(),
				Terminal:  false,
				NextNodes: map[rune]*TrieNode{},
			}
		}
		currNode = currNode.NextNodes[char]
	}

	currNode.Terminal = true
	return true
}

// Contains identifies whether the word is in the trie.  It returns a bool
// representing whether the word is in the trie.
func (n *TrieNode) Contains(word string) bool {
	currNode := n
	for _, char := range word {
		nextNode, ok := currNode.NextNodes[char]
		if !ok {
			return false
		}
		currNode = nextNode
	}
	return currNode.Terminal
}

// Delete removes the word from the node, if it is present in the trie. It
// returns whether the word was present in the trie.
func (n *TrieNode) Delete(word string) bool {
	if !n.Contains(word) {
		return false
	}

	var prefixWordTerminalNode *TrieNode
	var suffixInitialChar rune
	currNode := n

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
				n.NextNodes,
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

// ValidLettersBetweenPrefixAndSuffix returns the set of all letters '?'
// for which there is a word in the node that looks like: '{prefix}?{suffix}'.
func (n *TrieNode) ValidLettersBetweenPrefixAndSuffix(prefix, suffix string) map[rune]bool {

	validLetters := map[rune]bool{}
	currNode := n
	prefixOkay := true

	for _, prefixChar := range prefix {
		currNode, prefixOkay = currNode.NextNodes[prefixChar]

		// All placed prefixes of length > 1 should be valid words,
		// and therefore contained in lexicon.
		// However, in theory at least, a single character
		// could be placed as part of an across word, and there
		// could be no valid words in the lexicon that start with this
		// word. Here we take a precautionary approach to allow for
		// unanticipated use cases.
		if !prefixOkay {
			return validLetters
		}
	}

	middleNode := currNode

	for middleLetter, currNode := range middleNode.NextNodes {
		suffixOkay := true
		for _, suffixChar := range suffix {
			currNode, suffixOkay = currNode.NextNodes[suffixChar]
			if !suffixOkay {
				break
			}
		}
		if suffixOkay && currNode.Terminal {
			validLetters[middleLetter] = true
		}
	}
	return validLetters
}

func (n *TrieNode) IsRoot() bool {
	return n.Label == ""
}

func (n *TrieNode) IncomingEdge() rune {
	if n.IsRoot() {
		return 0
	}
	return rune(n.Label[len(n.Label)-1])
}

type EdgePruner interface {
	IsValidEdge(edge rune) bool
	Terminate(node *TrieNode) bool
}

type Visitor interface {
	Visit(*TrieNode)
	Exit(*TrieNode)
}

type PrunerVisitor interface {
	EdgePruner
	Visitor
}

// VisitNodesWithPruning calls the provided hook while traversing the nodes
// of the trie rooted at the receiver node using a pruned depth first traversal.
func (n *TrieNode) VisitNodesWithPruning(prunerVisitor PrunerVisitor) {
	prunerVisitor.Visit(n)

	if !prunerVisitor.Terminate(n) {
		for edge, nextNode := range n.NextNodes {
			if prunerVisitor.IsValidEdge(edge) {
				nextNode.VisitNodesWithPruning(prunerVisitor)
			}
		}
	}

	prunerVisitor.Exit(n)
}

// FollowEdges is used for following the edges in a trie and returns the node at
// the end of the path
func (n *TrieNode) FollowEdges(word string) *TrieNode {
	currNode := n
	for _, char := range word {
		nextNode, ok := currNode.NextNodes[char]
		if !ok {
			return nil
		}
		currNode = nextNode
	}
	return currNode
}

// New returns a pointer to a new empty TrieNode.
func New() *TrieNode {
	return &TrieNode{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{},
	}
}
