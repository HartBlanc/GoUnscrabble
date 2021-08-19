package lexicon

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	expectedTrie := &Trie{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*Trie{},
	}

	assert.Equal(t, expectedTrie, New())
}

func TestInsertEmpty(t *testing.T) {

	trie := &Trie{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*Trie{},
	}
	assert.True(t, trie.Insert("abc"))

	cNode := &Trie{
		Label:     "abc",
		Terminal:  true,
		NextNodes: map[rune]*Trie{},
	}
	bNode := &Trie{
		Label:     "ab",
		Terminal:  false,
		NextNodes: map[rune]*Trie{'c': cNode},
	}
	aNode := &Trie{
		Label:     "a",
		Terminal:  false,
		NextNodes: map[rune]*Trie{'b': bNode},
	}

	expectedTrie := &Trie{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*Trie{'a': aNode},
	}

	assert.Equal(t, expectedTrie, trie)
}

func TestInsertDisjoint(t *testing.T) {

	trie := &Trie{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*Trie{},
	}
	assert.True(t, trie.Insert("abc"))
	assert.True(t, trie.Insert("def"))

	cNode := &Trie{
		Label:     "abc",
		Terminal:  true,
		NextNodes: map[rune]*Trie{},
	}
	bNode := &Trie{
		Label:     "ab",
		Terminal:  false,
		NextNodes: map[rune]*Trie{'c': cNode},
	}
	aNode := &Trie{
		Label:     "a",
		Terminal:  false,
		NextNodes: map[rune]*Trie{'b': bNode},
	}

	fNode := &Trie{
		Label:     "def",
		Terminal:  true,
		NextNodes: map[rune]*Trie{},
	}
	eNode := &Trie{
		Label:     "de",
		Terminal:  false,
		NextNodes: map[rune]*Trie{'f': fNode},
	}
	dNode := &Trie{
		Label:     "d",
		Terminal:  false,
		NextNodes: map[rune]*Trie{'e': eNode},
	}

	expectedTrie := &Trie{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*Trie{'d': dNode, 'a': aNode},
	}

	assert.Equal(t, expectedTrie, trie)
}

func TestInsertSharedPrefix(t *testing.T) {

	trie := &Trie{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*Trie{},
	}
	assert.True(t, trie.Insert("abce"))
	assert.True(t, trie.Insert("abcd"))

	eNode := &Trie{
		Label:     "abce",
		Terminal:  true,
		NextNodes: map[rune]*Trie{},
	}

	dNode := &Trie{
		Label:     "abcd",
		Terminal:  true,
		NextNodes: map[rune]*Trie{},
	}

	cNode := &Trie{
		Label:     "abc",
		Terminal:  false,
		NextNodes: map[rune]*Trie{'d': dNode, 'e': eNode},
	}
	bNode := &Trie{
		Label:     "ab",
		Terminal:  false,
		NextNodes: map[rune]*Trie{'c': cNode},
	}
	aNode := &Trie{
		Label:     "a",
		Terminal:  false,
		NextNodes: map[rune]*Trie{'b': bNode},
	}

	expectedTrie := &Trie{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*Trie{'a': aNode},
	}

	assert.Equal(t, expectedTrie, trie)
}

func TestInsertSameWordTwice(t *testing.T) {

	trie := &Trie{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*Trie{},
	}
	assert.True(t, trie.Insert("a"))

	aNode := &Trie{
		Label:     "a",
		Terminal:  true,
		NextNodes: map[rune]*Trie{},
	}
	expectedTrie := &Trie{
		Label:     "",
		Terminal:  false,
		NextNodes: map[rune]*Trie{'a': aNode},
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

func createTrie() *Trie {
	carsNode := &Trie{
		Label:     "cars",
		Terminal:  true,
		NextNodes: map[rune]*Trie{},
	}
	catsNode := &Trie{
		Label:     "cats",
		Terminal:  true,
		NextNodes: map[rune]*Trie{},
	}
	carNode := &Trie{
		Label:     "car",
		Terminal:  true,
		NextNodes: map[rune]*Trie{'s': carsNode},
	}
	catNode := &Trie{
		Label:     "cat",
		Terminal:  true,
		NextNodes: map[rune]*Trie{'s': catsNode},
	}
	caNode := &Trie{
		Label:     "ca",
		Terminal:  false,
		NextNodes: map[rune]*Trie{'r': carNode, 't': catNode},
	}
	cNode := &Trie{
		Label:     "c",
		Terminal:  false,
		NextNodes: map[rune]*Trie{'a': caNode},
	}
	dogsNode := &Trie{
		Label:     "dogs",
		Terminal:  true,
		NextNodes: map[rune]*Trie{},
	}
	dogNode := &Trie{
		Label:     "dog",
		Terminal:  true,
		NextNodes: map[rune]*Trie{'s': dogsNode},
	}
	doneNode := &Trie{
		Label:     "done",
		Terminal:  true,
		NextNodes: map[rune]*Trie{},
	}
	donNode := &Trie{
		Label:     "don",
		Terminal:  false,
		NextNodes: map[rune]*Trie{'e': doneNode},
	}
	doNode := &Trie{
		Label:     "do",
		Terminal:  true,
		NextNodes: map[rune]*Trie{'n': donNode, 'g': dogNode},
	}
	dNode := &Trie{
		Label:     "d",
		Terminal:  false,
		NextNodes: map[rune]*Trie{'o': doNode},
	}
	earsNode := &Trie{
		Label:     "ears",
		Terminal:  true,
		NextNodes: map[rune]*Trie{},
	}
	earNode := &Trie{
		Label:     "ear",
		Terminal:  true,
		NextNodes: map[rune]*Trie{'s': earsNode},
	}
	eatsNode := &Trie{
		Label:     "eats",
		Terminal:  true,
		NextNodes: map[rune]*Trie{},
	}
	eatNode := &Trie{
		Label:     "eat",
		Terminal:  true,
		NextNodes: map[rune]*Trie{'s': eatsNode},
	}
	eaNode := &Trie{
		Label:     "ea",
		Terminal:  false,
		NextNodes: map[rune]*Trie{'r': earNode, 't': eatNode},
	}
	eNode := &Trie{
		Label:     "e",
		Terminal:  false,
		NextNodes: map[rune]*Trie{'a': eaNode},
	}
	beNode := &Trie{
		Label:     "be",
		Terminal:  true,
		NextNodes: map[rune]*Trie{},
	}
	bNode := &Trie{
		Label:     "b",
		Terminal:  false,
		NextNodes: map[rune]*Trie{'e': beNode},
	}
	aNode := &Trie{
		Label:     "a",
		Terminal:  true,
		NextNodes: map[rune]*Trie{},
	}
	trie := &Trie{
		Label:    "",
		Terminal: false,
		NextNodes: map[rune]*Trie{
			'a': aNode,
			'b': bNode,
			'c': cNode,
			'd': dNode,
			'e': eNode,
		},
	}
	return trie
}
