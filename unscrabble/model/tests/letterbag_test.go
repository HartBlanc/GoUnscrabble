package model_test

import (
	"testing"

	"example.com/unscrabble/unscrabble/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLetterBagReturnsExpectedLetterBag(t *testing.T) {
	letterCounts := map[rune]int{'a': 2, 'b': 1, 'c': 3}
	letterBag := model.NewLetterBag(letterCounts)

	expectedContents := []rune{'a', 'a', 'b', 'c', 'c', 'c'}
	assert.ElementsMatch(t, expectedContents, letterBag)
}

func TestPopRandomLetterRemovesAndReturnsLetters(t *testing.T) {
	letterCounts := map[rune]int{'a': 2, 'b': 1, 'c': 3}
	letterBag := model.NewLetterBag(letterCounts)

	var letters []rune
	for i := 0; i < 6; i++ {
		letter, err := letterBag.PopRandomLetter()
		require.NoError(t, err)
		letters = append(letters, letter)
		assert.Len(t, letterBag, 6-(i+1))
	}

	expectedContents := []rune{'a', 'a', 'b', 'c', 'c', 'c'}
	assert.ElementsMatch(t, expectedContents, letters)
	assert.Empty(t, letterBag)
}
