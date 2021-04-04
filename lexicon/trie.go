package lexicon

import (
	"bufio"
	"os"
	"strings"

	"example.com/unscrabble/set"
)

// CreateTrieFromFile builds a trie from a file which has a single word on each
// line.
func CreateTrieFromFile(filePath string) *Node {
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
func (node *Node) Insert(word string) bool {
	if node.Contains(word) {
		return false
	}
	var stringBuilder strings.Builder
	currNode := node
	for _, char := range word {
		stringBuilder.WriteRune(char)
		if _, ok := currNode.NextNodes[char]; !ok {
			currNode.NextNodes[char] = &Node{
				Label:     stringBuilder.String(),
				Terminal:  false,
				NextNodes: map[rune]*Node{},
			}
		}
		currNode = currNode.NextNodes[char]
	}

	currNode.Terminal = true
	return true
}

// Contains identifies whether the word is in the trie.  It returns a bool
// representing whether the word is in the trie.
func (node *Node) Contains(word string) bool {
	currNode := node
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
func (node *Node) Delete(word string) bool {
	if !node.Contains(word) {
		return false
	}

	var prefixWordTerminalNode *Node
	var suffixInitialChar rune
	currNode := node

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
				node.NextNodes,
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
func (node *Node) ValidLettersBetweenPrefixAndSuffix(
	prefix, suffix string,
) set.RuneMap {

	validLetters := make(set.RuneMap, 0)
	currNode := node
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
			validLetters.AddRune(middleLetter)
		}
	}
	return validLetters
}

// GenerateNodesWithPruning calls the provided hook while traversing the nodes
// of the trie rooted at the receiver node using a pruned depth first traversal.
func (node *Node) GenerateNodesWithPruning(
	validEdge func(rune) bool,
	preExpandHook func(rune, *Node),
	postExpandHook func(rune, *Node),
	terminate func(*Node) bool,
	processNode func(*Node),
) {
	processNode(node)
	if terminate(node) {
		return
	}
	for edge, nextNode := range node.NextNodes {
		if !validEdge(edge) {
			continue
		}
		preExpandHook(edge, nextNode)
		nextNode.GenerateNodesWithPruning(
			validEdge,
			preExpandHook,
			postExpandHook,
			terminate,
			processNode,
		)
		postExpandHook(edge, nextNode)
	}
}

// FollowEdges is used for following the edges in a trie and returns the node at
// the end of the path
func (node *Node) FollowEdges(word string) *Node {
	currNode := node
	for _, char := range word {
		nextNode, ok := currNode.NextNodes[char]
		if !ok {
			return nil
		}
		currNode = nextNode
	}
	return currNode
}

// New returns a pointer to a new empty Node.
func New() *Node {
	return &Node{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*Node{},
	}
}
