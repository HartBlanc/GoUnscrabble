package lexicon

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {

	expectedTrie := &Node{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*Node{},
	}

	assert.Equal(t, expectedTrie, New())
}

func TestInsertEmpty(t *testing.T) {

	trie := &Node{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*Node{},
	}
	assert.True(t, trie.Insert("abc"))

	cNode := &Node{
		Label:     "abc",
		Terminal:  true,
		NextNodes: map[rune]*Node{},
	}
	bNode := &Node{
		Label:     "ab",
		Terminal:  false,
		NextNodes: map[rune]*Node{'c': cNode},
	}
	aNode := &Node{
		Label:     "a",
		Terminal:  false,
		NextNodes: map[rune]*Node{'b': bNode},
	}

	expectedTrie := &Node{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*Node{'a': aNode},
	}

	assert.Equal(t, expectedTrie, trie)
}

func TestInsertDisjoint(t *testing.T) {

	trie := &Node{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*Node{},
	}
	assert.True(t, trie.Insert("abc"))
	assert.True(t, trie.Insert("def"))

	cNode := &Node{
		Label:     "abc",
		Terminal:  true,
		NextNodes: map[rune]*Node{},
	}
	bNode := &Node{
		Label:     "ab",
		Terminal:  false,
		NextNodes: map[rune]*Node{'c': cNode},
	}
	aNode := &Node{
		Label:     "a",
		Terminal:  false,
		NextNodes: map[rune]*Node{'b': bNode},
	}

	fNode := &Node{
		Label:     "def",
		Terminal:  true,
		NextNodes: map[rune]*Node{},
	}
	eNode := &Node{
		Label:     "de",
		Terminal:  false,
		NextNodes: map[rune]*Node{'f': fNode},
	}
	dNode := &Node{
		Label:     "d",
		Terminal:  false,
		NextNodes: map[rune]*Node{'e': eNode},
	}

	expectedTrie := &Node{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*Node{'d': dNode, 'a': aNode},
	}

	assert.Equal(t, expectedTrie, trie)
}

func TestInsertSharedPrefix(t *testing.T) {

	trie := &Node{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*Node{},
	}
	assert.True(t, trie.Insert("abce"))
	assert.True(t, trie.Insert("abcd"))

	eNode := &Node{
		Label:     "abce",
		Terminal:  true,
		NextNodes: map[rune]*Node{},
	}

	dNode := &Node{
		Label:     "abcd",
		Terminal:  true,
		NextNodes: map[rune]*Node{},
	}

	cNode := &Node{
		Label:     "abc",
		Terminal:  false,
		NextNodes: map[rune]*Node{'d': dNode, 'e': eNode},
	}
	bNode := &Node{
		Label:     "ab",
		Terminal:  false,
		NextNodes: map[rune]*Node{'c': cNode},
	}
	aNode := &Node{
		Label:     "a",
		Terminal:  false,
		NextNodes: map[rune]*Node{'b': bNode},
	}

	expectedTrie := &Node{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*Node{'a': aNode},
	}

	assert.Equal(t, expectedTrie, trie)
}

func TestInsertSameWordTwice(t *testing.T) {

	trie := &Node{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*Node{},
	}
	assert.True(t, trie.Insert("a"))

	aNode := &Node{
		Label:     "a",
		Terminal:  true,
		NextNodes: map[rune]*Node{},
	}
	expectedTrie := &Node{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*Node{'a': aNode},
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

func createTrie() *Node {
	carsNode := &Node{
		Label:     "cars",
		Terminal:  true,
		NextNodes: map[rune]*Node{},
	}
	catsNode := &Node{
		Label:     "cats",
		Terminal:  true,
		NextNodes: map[rune]*Node{},
	}
	carNode := &Node{
		Label:     "car",
		Terminal:  true,
		NextNodes: map[rune]*Node{'s': carsNode},
	}
	catNode := &Node{
		Label:     "cat",
		Terminal:  true,
		NextNodes: map[rune]*Node{'s': catsNode},
	}
	caNode := &Node{
		Label:     "ca",
		Terminal:  false,
		NextNodes: map[rune]*Node{'r': carNode, 't': catNode},
	}
	cNode := &Node{
		Label:     "c",
		Terminal:  false,
		NextNodes: map[rune]*Node{'a': caNode},
	}
	dogsNode := &Node{
		Label:     "dogs",
		Terminal:  true,
		NextNodes: map[rune]*Node{},
	}
	dogNode := &Node{
		Label:     "dog",
		Terminal:  true,
		NextNodes: map[rune]*Node{'s': dogsNode},
	}
	doneNode := &Node{
		Label:     "done",
		Terminal:  true,
		NextNodes: map[rune]*Node{},
	}
	donNode := &Node{
		Label:     "don",
		Terminal:  false,
		NextNodes: map[rune]*Node{'e': doneNode},
	}
	doNode := &Node{
		Label:     "do",
		Terminal:  true,
		NextNodes: map[rune]*Node{'n': donNode, 'g': dogNode},
	}
	dNode := &Node{
		Label:     "d",
		Terminal:  false,
		NextNodes: map[rune]*Node{'o': doNode},
	}
	earsNode := &Node{
		Label:     "ears",
		Terminal:  true,
		NextNodes: map[rune]*Node{},
	}
	earNode := &Node{
		Label:     "ear",
		Terminal:  true,
		NextNodes: map[rune]*Node{'s': earsNode},
	}
	eatsNode := &Node{
		Label:     "eats",
		Terminal:  true,
		NextNodes: map[rune]*Node{},
	}
	eatNode := &Node{
		Label:     "eat",
		Terminal:  true,
		NextNodes: map[rune]*Node{'s': eatsNode},
	}
	eaNode := &Node{
		Label:     "ea",
		Terminal:  false,
		NextNodes: map[rune]*Node{'r': earNode, 't': eatNode},
	}
	eNode := &Node{
		Label:     "e",
		Terminal:  false,
		NextNodes: map[rune]*Node{'a': eaNode},
	}
	beNode := &Node{
		Label:     "be",
		Terminal:  true,
		NextNodes: map[rune]*Node{},
	}
	bNode := &Node{
		Label:     "b",
		Terminal:  false,
		NextNodes: map[rune]*Node{'e': beNode},
	}
	aNode := &Node{
		Label:     "a",
		Terminal:  true,
		NextNodes: map[rune]*Node{},
	}
	trie := &Node{
		Label:    "",
		Terminal: false,
		NextNodes: map[rune]*Node{
			'a': aNode,
			'b': bNode,
			'c': cNode,
			'd': dNode,
			'e': eNode,
		},
	}
	return trie
}