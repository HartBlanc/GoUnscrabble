package lexicon

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	expectedTrie := &TrieNode{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{},
	}

	assert.Equal(t, expectedTrie, New())
}

func TestInsertEmpty(t *testing.T) {

	trie := &TrieNode{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{},
	}
	assert.True(t, trie.Insert("abc"))

	cNode := &TrieNode{
		Label:     "abc",
		Terminal:  true,
		NextNodes: map[rune]*TrieNode{},
	}
	bNode := &TrieNode{
		Label:     "ab",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{'c': cNode},
	}
	aNode := &TrieNode{
		Label:     "a",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{'b': bNode},
	}

	expectedTrie := &TrieNode{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{'a': aNode},
	}

	assert.Equal(t, expectedTrie, trie)
}

func TestInsertDisjoint(t *testing.T) {

	trie := &TrieNode{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{},
	}
	assert.True(t, trie.Insert("abc"))
	assert.True(t, trie.Insert("def"))

	cNode := &TrieNode{
		Label:     "abc",
		Terminal:  true,
		NextNodes: map[rune]*TrieNode{},
	}
	bNode := &TrieNode{
		Label:     "ab",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{'c': cNode},
	}
	aNode := &TrieNode{
		Label:     "a",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{'b': bNode},
	}

	fNode := &TrieNode{
		Label:     "def",
		Terminal:  true,
		NextNodes: map[rune]*TrieNode{},
	}
	eNode := &TrieNode{
		Label:     "de",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{'f': fNode},
	}
	dNode := &TrieNode{
		Label:     "d",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{'e': eNode},
	}

	expectedTrie := &TrieNode{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{'d': dNode, 'a': aNode},
	}

	assert.Equal(t, expectedTrie, trie)
}

func TestInsertSharedPrefix(t *testing.T) {

	trie := &TrieNode{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{},
	}
	assert.True(t, trie.Insert("abce"))
	assert.True(t, trie.Insert("abcd"))

	eNode := &TrieNode{
		Label:     "abce",
		Terminal:  true,
		NextNodes: map[rune]*TrieNode{},
	}

	dNode := &TrieNode{
		Label:     "abcd",
		Terminal:  true,
		NextNodes: map[rune]*TrieNode{},
	}

	cNode := &TrieNode{
		Label:     "abc",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{'d': dNode, 'e': eNode},
	}
	bNode := &TrieNode{
		Label:     "ab",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{'c': cNode},
	}
	aNode := &TrieNode{
		Label:     "a",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{'b': bNode},
	}

	expectedTrie := &TrieNode{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{'a': aNode},
	}

	assert.Equal(t, expectedTrie, trie)
}

func TestInsertSameWordTwice(t *testing.T) {

	trie := &TrieNode{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{},
	}
	assert.True(t, trie.Insert("a"))

	aNode := &TrieNode{
		Label:     "a",
		Terminal:  true,
		NextNodes: map[rune]*TrieNode{},
	}
	expectedTrie := &TrieNode{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{'a': aNode},
	}

	assert.Equal(t, expectedTrie, trie)
	assert.False(t, trie.Insert("a"))
	assert.Equal(t, expectedTrie, trie)
}

func TestContains(t *testing.T) {

	trie := createTrie()
	testCases := []struct {
		Name             string
		Word             string
		ExpectedContains bool
	}{
		{
			Name:             "present and terminal",
			Word:             "dog",
			ExpectedContains: true,
		},
		{
			Name:             "present but not terminal",
			Word:             "ea",
			ExpectedContains: false,
		},
		{
			Name:             "not present",
			Word:             "missing",
			ExpectedContains: false,
		},
		{
			Name:             "empty string",
			Word:             "",
			ExpectedContains: false,
		},
	}

	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.Name, func(t *testing.T) {

			assert.Equal(
				t,
				testCase.ExpectedContains,
				trie.Contains(testCase.Word),
			)
		})
	}
}

func TestDelete(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {

		trie := createTrie()
		assert.False(t, trie.Delete(""))
		assert.Equal(t, createTrie(), trie)
	})

	t.Run("word not present", func(t *testing.T) {

		trie := createTrie()
		assert.False(t, trie.Delete("missing"))
		assert.Equal(t, createTrie(), trie)
	})

	t.Run("word has prefix", func(t *testing.T) {

		trie := createTrie()
		assert.True(t, trie.Delete("cars"))
		expectedTrie := createTrie()
		delete(expectedTrie.NextNodes['c'].NextNodes['a'].NextNodes['r'].NextNodes, 's')
		assert.Equal(
			t,
			expectedTrie,
			trie,
		)
	})

	t.Run("word is prefix", func(t *testing.T) {
		trie := createTrie()
		assert.True(t, trie.Delete("car"))

		expectedTrie := createTrie()
		expectedTrie.NextNodes['c'].NextNodes['a'].NextNodes['r'].Terminal = false
		assert.Equal(
			t,
			expectedTrie,
			trie,
		)
	})

	t.Run("no prefixes", func(t *testing.T) {
		trie := createTrie()
		trie.Delete("be")

		expectedTrie := createTrie()
		delete(expectedTrie.NextNodes, 'b')

		assert.Equal(
			t,
			expectedTrie,
			trie,
		)
	})
}

func TestValidLettersBetweenPrefixAndSuffix(t *testing.T) {
	trie := createTrie()

	t.Run("empty prefix", func(t *testing.T) {
		crossSet := trie.ValidLettersBetweenPrefixAndSuffix("", "o")
		assert.Equal(
			t,
			map[rune]bool{
				'd': true,
			},
			crossSet,
		)
	})
	t.Run("empty suffix", func(t *testing.T) {
		crossSet := trie.ValidLettersBetweenPrefixAndSuffix("do", "")
		assert.Equal(
			t,
			map[rune]bool{
				'g': true,
			},
			crossSet,
		)
	})
	t.Run("empty prefix and empty suffix returns single letter words", func(t *testing.T) {
		crossSet := trie.ValidLettersBetweenPrefixAndSuffix("", "")
		assert.Equal(
			t,
			map[rune]bool{'a': true},
			crossSet,
		)
	})
	t.Run("prefix and suffix", func(t *testing.T) {
		crossSet := trie.ValidLettersBetweenPrefixAndSuffix("ca", "s")
		assert.Equal(
			t,
			map[rune]bool{
				'r': true,
				't': true,
			},
			crossSet,
		)
	})
	t.Run("break in suffix", func(t *testing.T) {
		crossSet := trie.ValidLettersBetweenPrefixAndSuffix("", "z")
		assert.Equal(
			t,
			map[rune]bool{},
			crossSet,
		)
	})
	t.Run("break in prefix", func(t *testing.T) {
		crossSet := trie.ValidLettersBetweenPrefixAndSuffix("z", "")
		assert.Equal(
			t,
			map[rune]bool{},
			crossSet,
		)
	})
	t.Run("no cross chars ", func(t *testing.T) {
		crossSet := trie.ValidLettersBetweenPrefixAndSuffix("a", "")
		assert.Equal(
			t,
			map[rune]bool{},
			crossSet,
		)
	})
	t.Run("final suffix node is not terminal", func(t *testing.T) {
		crossSet := trie.ValidLettersBetweenPrefixAndSuffix("d", "n")
		assert.Equal(
			t,
			map[rune]bool{},
			crossSet,
		)
	})
}

func createTrie() *TrieNode {
	carsNode := &TrieNode{
		Label:     "cars",
		Terminal:  true,
		NextNodes: map[rune]*TrieNode{},
	}
	catsNode := &TrieNode{
		Label:     "cats",
		Terminal:  true,
		NextNodes: map[rune]*TrieNode{},
	}
	carNode := &TrieNode{
		Label:     "car",
		Terminal:  true,
		NextNodes: map[rune]*TrieNode{'s': carsNode},
	}
	catNode := &TrieNode{
		Label:     "cat",
		Terminal:  true,
		NextNodes: map[rune]*TrieNode{'s': catsNode},
	}
	caNode := &TrieNode{
		Label:     "ca",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{'r': carNode, 't': catNode},
	}
	cNode := &TrieNode{
		Label:     "c",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{'a': caNode},
	}
	dogsNode := &TrieNode{
		Label:     "dogs",
		Terminal:  true,
		NextNodes: map[rune]*TrieNode{},
	}
	dogNode := &TrieNode{
		Label:     "dog",
		Terminal:  true,
		NextNodes: map[rune]*TrieNode{'s': dogsNode},
	}
	doneNode := &TrieNode{
		Label:     "done",
		Terminal:  true,
		NextNodes: map[rune]*TrieNode{},
	}
	donNode := &TrieNode{
		Label:     "don",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{'e': doneNode},
	}
	doNode := &TrieNode{
		Label:     "do",
		Terminal:  true,
		NextNodes: map[rune]*TrieNode{'n': donNode, 'g': dogNode},
	}
	dNode := &TrieNode{
		Label:     "d",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{'o': doNode},
	}
	earsNode := &TrieNode{
		Label:     "ears",
		Terminal:  true,
		NextNodes: map[rune]*TrieNode{},
	}
	earNode := &TrieNode{
		Label:     "ear",
		Terminal:  true,
		NextNodes: map[rune]*TrieNode{'s': earsNode},
	}
	eatsNode := &TrieNode{
		Label:     "eats",
		Terminal:  true,
		NextNodes: map[rune]*TrieNode{},
	}
	eatNode := &TrieNode{
		Label:     "eat",
		Terminal:  true,
		NextNodes: map[rune]*TrieNode{'s': eatsNode},
	}
	eaNode := &TrieNode{
		Label:     "ea",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{'r': earNode, 't': eatNode},
	}
	eNode := &TrieNode{
		Label:     "e",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{'a': eaNode},
	}
	beNode := &TrieNode{
		Label:     "be",
		Terminal:  true,
		NextNodes: map[rune]*TrieNode{},
	}
	bNode := &TrieNode{
		Label:     "b",
		Terminal:  false,
		NextNodes: map[rune]*TrieNode{'e': beNode},
	}
	aNode := &TrieNode{
		Label:     "a",
		Terminal:  true,
		NextNodes: map[rune]*TrieNode{},
	}
	trie := &TrieNode{
		Label:    "",
		Terminal: false,
		NextNodes: map[rune]*TrieNode{
			'a': aNode,
			'b': bNode,
			'c': cNode,
			'd': dNode,
			'e': eNode,
		},
	}
	return trie
}
