package lexicon

import (
	"bufio"
	"os"
	"strings"
)

// Trie is a data structure used for efficient prefix searches
type Trie struct {
	Label     string
	Terminal  bool
	NextNodes map[rune]*Trie
}

// CreateTrieFromFile builds a trie from a file which has a single word on each
// line.
func CreateTrieFromFile(filePath string) *Trie {
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
func (n *Trie) Insert(word string) bool {
	if n.Contains(word) {
		return false
	}
	var stringBuilder strings.Builder
	currNode := n
	for _, char := range word {
		stringBuilder.WriteRune(char)
		if _, ok := currNode.NextNodes[char]; !ok {
			currNode.NextNodes[char] = &Trie{
				Label:     stringBuilder.String(),
				Terminal:  false,
				NextNodes: map[rune]*Trie{},
			}
		}
		currNode = currNode.NextNodes[char]
	}

	currNode.Terminal = true
	return true
}

// Contains identifies whether the word is in the trie.  It returns a bool
// representing whether the word is in the trie.
func (n *Trie) Contains(word string) bool {
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
func (n *Trie) Delete(word string) bool {
	if !n.Contains(word) {
		return false
	}

	var prefixWordTerminalNode *Trie
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
func (n *Trie) ValidLettersBetweenPrefixAndSuffix(prefix, suffix string) map[rune]bool {

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

// GenerateNodesWithPruning calls the provided hook while traversing the nodes
// of the trie rooted at the receiver node using a pruned depth first traversal.
func (n *Trie) GenerateNodesWithPruning(
	validEdge func(rune) bool,
	preExpandHook func(rune, *Trie),
	postExpandHook func(rune, *Trie),
	terminate func(*Trie) bool,
	processNode func(*Trie),
) {
	processNode(n)
	if terminate(n) {
		return
	}
	for edge, nextNode := range n.NextNodes {
		if !validEdge(edge) {
			continue
		}
		preExpandHook(edge, n)
		nextNode.GenerateNodesWithPruning(
			validEdge,
			preExpandHook,
			postExpandHook,
			terminate,
			processNode,
		)
		postExpandHook(edge, n)
	}
}

// FollowEdges is used for following the edges in a trie and returns the node at
// the end of the path
func (n *Trie) FollowEdges(word string) *Trie {
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

// New returns a pointer to a new empty Node.
func New() *Trie {
	return &Trie{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*Trie{},
	}
}
