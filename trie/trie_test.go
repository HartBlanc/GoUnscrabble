package trie

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {

	expectedTrie := &Trie{
		Root: &Node{
			Label:     "",
			Terminal:  false,
			NextNodes: map[rune]*Node{},
		},
	}

	assert.Equal(t, expectedTrie, New())
}

func TestInsertEmpty(t *testing.T) {

	trie := &Trie{
		Root: &Node{
			Label:     "",
			Terminal:  false,
			NextNodes: map[rune]*Node{},
		},
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

	expectedTrie := &Trie{
		Root: &Node{
			Label:     "",
			Terminal:  false,
			NextNodes: map[rune]*Node{'a': aNode},
		},
	}

	assert.Equal(t, expectedTrie, trie)
}

func TestInsertDisjoint(t *testing.T) {

	trie := &Trie{
		Root: &Node{
			Label:     "",
			Terminal:  false,
			NextNodes: map[rune]*Node{},
		},
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

	expectedTrie := &Trie{
		Root: &Node{
			Label:     "",
			Terminal:  false,
			NextNodes: map[rune]*Node{'d': dNode, 'a': aNode},
		},
	}

	assert.Equal(t, expectedTrie, trie)
}

func TestInsertSharedPrefix(t *testing.T) {

	trie := &Trie{
		Root: &Node{
			Label:     "",
			Terminal:  false,
			NextNodes: map[rune]*Node{},
		},
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

	expectedTrie := &Trie{
		Root: &Node{
			Label:     "",
			Terminal:  false,
			NextNodes: map[rune]*Node{'a': aNode},
		},
	}

	assert.Equal(t, expectedTrie, trie)
}

func TestInsertSameWordTwice(t *testing.T) {

	trie := &Trie{
		Root: &Node{
			Label:     "",
			Terminal:  false,
			NextNodes: map[rune]*Node{},
		},
	}
	assert.True(t, trie.Insert("a"))

	aNode := &Node{
		Label:     "a",
		Terminal:  true,
		NextNodes: map[rune]*Node{},
	}
	expectedTrie := &Trie{
		Root: &Node{
			Label:     "",
			Terminal:  false,
			NextNodes: map[rune]*Node{'a': aNode},
		},
	}

	assert.Equal(t, expectedTrie, trie)
	assert.False(t, trie.Insert("a"))
	assert.Equal(t, expectedTrie, trie)
}

func TestContains(t *testing.T) {

	dNode := &Node{
		Label:     "abcd",
		Terminal:  true,
		NextNodes: map[rune]*Node{},
	}
	cNode := &Node{
		Label:     "abc",
		Terminal:  false,
		NextNodes: map[rune]*Node{'d': dNode},
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

	trie := &Trie{
		Root: &Node{
			Label:     "",
			Terminal:  false,
			NextNodes: map[rune]*Node{'a': aNode},
		},
	}

	testCases := []struct {
		Name             string
		Word             string
		ExpectedContains bool
	}{
		{
			Name:             "present and terminal",
			Word:             "abcd",
			ExpectedContains: true,
		},
		{
			Name:             "present but not terminal",
			Word:             "abc",
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
		assert.True(t, trie.Delete("abcd"))
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
			Label:     "ef",
			Terminal:  true,
			NextNodes: map[rune]*Node{},
		}
		eNode := &Node{
			Label:     "e",
			Terminal:  false,
			NextNodes: map[rune]*Node{'f': fNode},
		}
		assert.Equal(
			t,
			&Trie{
				Root: &Node{
					Label:     "",
					Terminal:  false,
					NextNodes: map[rune]*Node{'a': aNode, 'e': eNode},
				},
			},
			trie,
		)
	})

	t.Run("word is prefix", func(t *testing.T) {

		dNode := &Node{
			Label:     "abcd",
			Terminal:  true,
			NextNodes: map[rune]*Node{},
		}
		cNode := &Node{
			Label:     "abc",
			Terminal:  false,
			NextNodes: map[rune]*Node{'d': dNode},
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
			Label:     "ef",
			Terminal:  true,
			NextNodes: map[rune]*Node{},
		}
		eNode := &Node{
			Label:     "e",
			Terminal:  false,
			NextNodes: map[rune]*Node{'f': fNode},
		}

		trie := createTrie()
		assert.True(t, trie.Delete("abc"))
		assert.Equal(
			t,
			&Trie{
				Root: &Node{
					Label:     "",
					Terminal:  false,
					NextNodes: map[rune]*Node{'a': aNode, 'e': eNode},
				},
			},
			trie,
		)
	})

	t.Run("no prefixes", func(t *testing.T) {

		dNode := &Node{
			Label:     "abcd",
			Terminal:  true,
			NextNodes: map[rune]*Node{},
		}
		cNode := &Node{
			Label:     "abc",
			Terminal:  true,
			NextNodes: map[rune]*Node{'d': dNode},
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

		trie := createTrie()
		trie.Delete("ef")
		assert.Equal(
			t,
			&Trie{
				Root: &Node{
					Label:     "",
					Terminal:  false,
					NextNodes: map[rune]*Node{'a': aNode},
				},
			},
			trie,
		)
	})
}

func createTrie() *Trie {
	dNode := &Node{
		Label:     "abcd",
		Terminal:  true,
		NextNodes: map[rune]*Node{},
	}

	cNode := &Node{
		Label:     "abc",
		Terminal:  true,
		NextNodes: map[rune]*Node{'d': dNode},
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
		Label:     "ef",
		Terminal:  true,
		NextNodes: map[rune]*Node{},
	}
	eNode := &Node{
		Label:     "e",
		Terminal:  false,
		NextNodes: map[rune]*Node{'f': fNode},
	}
	trie := &Trie{
		Root: &Node{
			Label:     "",
			Terminal:  false,
			NextNodes: map[rune]*Node{'a': aNode, 'e': eNode},
		},
	}

	return trie

}
