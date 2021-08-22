package trie_test

import (
	"testing"

	"example.com/unscrabble/lexicon"
	"example.com/unscrabble/unscrabble/model"
	triemovegen "example.com/unscrabble/unscrabble/movegen/trie"

	assert "github.com/stretchr/testify/assert"
)

type MockCrossCheckSetGenerator struct{}

func (m MockCrossCheckSetGenerator) ValidLettersBetweenPrefixAndSuffix(
	prefix, suffix string,
) map[rune]bool {
	return nil
}

func TestTrieMoveGeneratorGeneratesMoves(t *testing.T) {
	// emptyMultipliers required for NewBoard to create board of required size (6x6)
	emptyMultipliers := make([][]int, 7, 7)
	for y := range emptyMultipliers {
		emptyMultipliers[y] = make([]int, 7, 7)
	}
	testBoard := model.NewBoard(MockCrossCheckSetGenerator{}, emptyMultipliers, emptyMultipliers)

	testTrieRoot := lexicon.NewTrieNode()
	testTrieRoot.Insert("inarack")
	testTrieRoot.Insert("notinarack")

	testRack := model.NewRack(10)
	for _, char := range "inarack" {
		testRack.AddRune(char)
	}

	testTrieMoveGen := triemovegen.NewTrieMoveGenertator(testTrieRoot)
	moves := testTrieMoveGen.GenerateMoves(testBoard, *testRack)
	expectedMoves := []model.Move{
		{
			StartPosition: &model.Position{
				Row:    3,
				Column: 0,
			},
			Horizontal: true,
			Word: model.Word{
				Chars:      "inarack",
				BlankTiles: make([]bool, 7),
			},
		},
		{
			StartPosition: &model.Position{
				Row:    0,
				Column: 3,
			},
			Horizontal: false,
			Word: model.Word{
				Chars:      "inarack",
				BlankTiles: make([]bool, 7),
			},
		},
	}
	assert.ElementsMatch(t, expectedMoves, moves)
}
