package unscrabble

import (
	"testing"

	"example.com/unscrabble/set"
	gomock "github.com/golang/mock/gomock"
	assert "github.com/stretchr/testify/assert"
)

func TestTranspose(t *testing.T) {

	t.Run("transpose empty", func(t *testing.T) {
		tiles := Board{}
		Transpose(tiles)
		assert.Equal(t, Board{}, tiles)
	})
	t.Run("transpose single", func(t *testing.T) {
		tiles := Board{
			{&BoardTile{BoardPosition: &Position{Row: 0, Column: 0}}},
		}
		expectedTiles := Board{
			{&BoardTile{BoardPosition: &Position{Row: 0, Column: 0}}},
		}
		Transpose(tiles)
		assert.Equal(t, expectedTiles, tiles)
	})
	t.Run("tiles are moved", func(t *testing.T) {

		initialBoard := Board{
			{
				&BoardTile{
					BoardPosition: &Position{
						Row:    0,
						Column: 0,
					},
				},
				&BoardTile{
					BoardPosition: &Position{
						Row:    0,
						Column: 1,
					},
				},
			},
			{
				&BoardTile{
					BoardPosition: &Position{
						Row:    1,
						Column: 0,
					},
				},
				&BoardTile{
					BoardPosition: &Position{
						Row:    1,
						Column: 1,
					},
				},
			},
		}
		transposedBoard := make(Board, len(initialBoard))
		for y := range initialBoard {
			transposedBoard[y] = make(
				[]*BoardTile,
				len(initialBoard[y]),
			)
			copy(transposedBoard[y], initialBoard[y])
		}
		Transpose(transposedBoard)
		for y, row := range transposedBoard {
			for x := range row {
				assert.Same(
					t,
					initialBoard[y][x],
					transposedBoard[x][y],
				)
			}
		}

		Transpose(transposedBoard)
		for y, row := range transposedBoard {
			for x := range row {
				assert.Same(
					t,
					initialBoard[y][x],
					transposedBoard[y][x],
				)
			}
		}
	})
	t.Run("cross checks and positions swapped", func(t *testing.T) {

		sets := make(map[rune]set.RuneSet)
		chars := [8]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}
		for _, char := range chars {
			mySet := set.New(1)
			mySet.AddRune(char)
			sets[char] = mySet
		}
		transposedBoard := Board{
			{
				&BoardTile{
					CrossCheckSet:          sets['a'],
					transposeCrossCheckSet: sets['b'],
					CrossScore:             1,
					transposeCrossScore:    2,
					BoardPosition: &Position{
						Row:    0,
						Column: 0,
					},
				},
				&BoardTile{
					CrossCheckSet:          sets['c'],
					transposeCrossCheckSet: sets['d'],
					CrossScore:             3,
					transposeCrossScore:    4,
					BoardPosition: &Position{
						Row:    0,
						Column: 1,
					},
				},
			},
			{
				&BoardTile{
					CrossCheckSet:          sets['e'],
					transposeCrossCheckSet: sets['f'],
					CrossScore:             5,
					transposeCrossScore:    6,
					BoardPosition: &Position{
						Row:    1,
						Column: 0,
					},
				},
				&BoardTile{
					CrossCheckSet:          sets['g'],
					transposeCrossCheckSet: sets['h'],
					CrossScore:             7,
					transposeCrossScore:    8,
					BoardPosition: &Position{
						Row:    1,
						Column: 1,
					},
				},
			},
		}
		// identical values to initial board
		expectedDoubleTransposeBoard := Board{
			{
				&BoardTile{
					CrossCheckSet:          sets['a'],
					transposeCrossCheckSet: sets['b'],
					CrossScore:             1,
					transposeCrossScore:    2,
					BoardPosition: &Position{
						Row:    0,
						Column: 0,
					},
				},
				&BoardTile{
					CrossCheckSet:          sets['c'],
					transposeCrossCheckSet: sets['d'],
					CrossScore:             3,
					transposeCrossScore:    4,
					BoardPosition: &Position{
						Row:    0,
						Column: 1,
					},
				},
			},
			{
				&BoardTile{
					CrossCheckSet:          sets['e'],
					transposeCrossCheckSet: sets['f'],
					CrossScore:             5,
					transposeCrossScore:    6,
					BoardPosition: &Position{
						Row:    1,
						Column: 0,
					},
				},
				&BoardTile{
					CrossCheckSet:          sets['g'],
					transposeCrossCheckSet: sets['h'],
					CrossScore:             7,
					transposeCrossScore:    8,
					BoardPosition: &Position{
						Row:    1,
						Column: 1,
					},
				},
			},
		}
		expectedTransposedBoard := Board{
			{
				&BoardTile{
					CrossCheckSet:          sets['b'],
					transposeCrossCheckSet: sets['a'],
					CrossScore:             2,
					transposeCrossScore:    1,
					BoardPosition: &Position{
						Row:    0,
						Column: 0,
					},
				},
				&BoardTile{
					CrossCheckSet:          sets['f'],
					transposeCrossCheckSet: sets['e'],
					CrossScore:             6,
					transposeCrossScore:    5,
					BoardPosition: &Position{
						Row:    0,
						Column: 1,
					},
				},
			},
			{
				&BoardTile{
					CrossCheckSet:          sets['d'],
					transposeCrossCheckSet: sets['c'],
					CrossScore:             4,
					transposeCrossScore:    3,
					BoardPosition: &Position{
						Row:    1,
						Column: 0,
					},
				},
				&BoardTile{
					CrossCheckSet:          sets['h'],
					transposeCrossCheckSet: sets['g'],
					CrossScore:             8,
					transposeCrossScore:    7,
					BoardPosition: &Position{
						Row:    1,
						Column: 1,
					},
				},
			},
		}
		Transpose(transposedBoard)
		assert.Equal(t, expectedTransposedBoard, transposedBoard)
		Transpose(transposedBoard)
		assert.Equal(t, expectedDoubleTransposeBoard, transposedBoard)
	})
}

func TestGetAnchors(t *testing.T) {

	t.Run("no tiles", func(t *testing.T) {
		tiles := Board{{}}
		expectedAnchors := []*BoardTile{}
		assert.Equal(t, expectedAnchors, GetAnchors(tiles))
	})
	t.Run("all tiles empty", func(t *testing.T) {
		tiles := Board{
			{&BoardTile{}, &BoardTile{}},
			{&BoardTile{}, &BoardTile{}},
		}
		expectedAnchors := []*BoardTile{}
		assert.Equal(t, expectedAnchors, GetAnchors(tiles))
	})
	t.Run("no empty tiles", func(t *testing.T) {
		tiles := Board{
			{&BoardTile{Letter: 'a'}, &BoardTile{Letter: 'a'}},
			{&BoardTile{Letter: 'a'}, &BoardTile{Letter: 'a'}},
		}
		expectedAnchors := []*BoardTile{}
		assert.Equal(t, expectedAnchors, GetAnchors(tiles))
	})
	t.Run("adjacent tiles are anchors", func(t *testing.T) {

		above := &BoardTile{BoardPosition: &Position{Row: 0, Column: 1}}
		left := &BoardTile{BoardPosition: &Position{Row: 1, Column: 0}}
		right := &BoardTile{BoardPosition: &Position{Row: 1, Column: 2}}
		below := &BoardTile{BoardPosition: &Position{Row: 2, Column: 1}}
		tiles := Board{
			{&BoardTile{}, above, &BoardTile{}},
			{left, &BoardTile{Letter: 'a'}, right},
			{&BoardTile{}, below, &BoardTile{}},
		}
		expectedAnchors := []*BoardTile{above, left, right, below}
		assert.ElementsMatch(t, expectedAnchors, GetAnchors(tiles))
	})
}

func TestGetPrefixAbove(t *testing.T) {
	t.Run("test empty string and zero score if no prefix above", func(t *testing.T) {
		tiles := newBoard()
		prefix, score := GetPrefixAbove(tiles[1][1], tiles)
		assert.Equal(t, "", prefix)
		assert.Equal(t, 0, score)
	})
	t.Run("test does not go out of bounds if tile on top row", func(t *testing.T) {
		tiles := newBoard()
		prefix, score := GetPrefixAbove(tiles[0][1], tiles)
		assert.Equal(t, "", prefix)
		assert.Equal(t, 0, score)
	})
	t.Run("test does not go out of bounds if prefix stops at top", func(t *testing.T) {
		tiles := newBoard()
		tiles[0][1].Letter = 'a'
		prefix, score := GetPrefixAbove(tiles[1][1], tiles)
		assert.Equal(t, "a", prefix)
		assert.Equal(t, 1, score)
	})
	t.Run("test stops at empty square", func(t *testing.T) {
		tiles := newBoard()
		tiles[0][1].Letter = 'b' // 1,1 is empty, should not reach here.
		tiles[2][1].Letter = 'a'
		prefix, score := GetPrefixAbove(tiles[3][1], tiles)
		assert.Equal(t, "a", prefix)
		assert.Equal(t, 1, score)
	})
	t.Run("prefix reverses correctly", func(t *testing.T) {
		tiles := newBoard()
		tiles[0][1].Letter = 'a'
		tiles[1][1].Letter = 'b'
		prefix, _ := GetPrefixAbove(tiles[2][1], tiles)
		assert.Equal(t, "ab", prefix)
	})
	t.Run("test multipliers ignored", func(t *testing.T) {
		tiles := newBoard()
		tiles[0][1].Letter = 'a'
		tiles[1][1].Letter = 'b'
		tiles[0][1].LetterMultiplier = 3
		tiles[1][1].WordMultiplier = 3
		_, score := GetPrefixAbove(tiles[2][1], tiles)
		assert.Equal(t, 1+4, score)
	})
}

func TestGetSuffixBelow(t *testing.T) {
	t.Run("test empty string and zero score if no suffix below", func(t *testing.T) {
		tiles := newBoard()
		suffix, score := GetSuffixBelow(tiles[1][1], tiles)
		assert.Equal(t, "", suffix)
		assert.Equal(t, 0, score)
	})
	t.Run("test does not go out of bounds if tile on bottom row", func(t *testing.T) {
		tiles := newBoard()
		suffix, score := GetSuffixBelow(tiles[len(tiles)-1][1], tiles)
		assert.Equal(t, "", suffix)
		assert.Equal(t, 0, score)
	})
	t.Run("test does not go out of bounds if suffix stops at bottom", func(t *testing.T) {
		tiles := newBoard()
		tiles[len(tiles)-1][1].Letter = 'a'
		suffix, score := GetSuffixBelow(tiles[len(tiles)-2][1], tiles)
		assert.Equal(t, "a", suffix)
		assert.Equal(t, 1, score)
	})
	t.Run("test stops at empty square", func(t *testing.T) {
		tiles := newBoard()
		tiles[1][1].Letter = 'a'
		tiles[3][1].Letter = 'b' // 2,1 is empty, should not reach here.
		suffix, score := GetSuffixBelow(tiles[0][1], tiles)
		assert.Equal(t, "a", suffix)
		assert.Equal(t, 1, score)
	})
	t.Run("test suffix is not reversed", func(t *testing.T) {
		tiles := newBoard()
		tiles[1][1].Letter = 'a'
		tiles[2][1].Letter = 'b'
		suffix, _ := GetSuffixBelow(tiles[0][1], tiles)
		assert.Equal(t, "ab", suffix)
	})
	t.Run("test multipliers ignored", func(t *testing.T) {
		tiles := newBoard()
		tiles[1][1].Letter = 'a'
		tiles[2][1].Letter = 'b'
		tiles[1][1].LetterMultiplier = 3
		tiles[2][1].WordMultiplier = 3
		_, score := GetPrefixAbove(tiles[3][1], tiles)
		assert.Equal(t, 1+4, score)
	})
}

func TestCrossCheck(t *testing.T) {
	t.Run("test no prefix or suffix returns nil set and zero score, without consulting lexicon", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tiles := newBoard()
		tile := tiles[0][0]
		mockLexicon := NewMockLexicon(ctrl)
		mockLexicon.EXPECT().ValidLettersBetweenPrefixAndSuffix(
			gomock.Any(),
			gomock.Any(),
		).Times(0)
		crossCheckSet, score := CrossCheck(tile, tiles, mockLexicon)
		assert.Nil(t, crossCheckSet)
		assert.Equal(t, 0, score)
	})
	t.Run("test prefix", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tiles := newBoard()
		tiles[0][0].Letter = 'c'
		tiles[1][0].Letter = 'a'
		tile := tiles[2][0]
		expectedCrossCheckSet := set.RuneMap{'t': true, 'r': true}
		mockLexicon := NewMockLexicon(ctrl)
		mockLexicon.EXPECT().ValidLettersBetweenPrefixAndSuffix(
			"ca",
			"",
		).Times(1).Return(expectedCrossCheckSet)
		crossCheckSet, score := CrossCheck(tile, tiles, mockLexicon)
		assert.Equal(t, expectedCrossCheckSet, crossCheckSet)
		assert.Equal(t, 5, score)
	})
	t.Run("test suffix", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tiles := newBoard()
		tiles[1][0].Letter = 'a'
		tiles[2][0].Letter = 't'
		tile := tiles[0][0]
		expectedCrossCheckSet := set.RuneMap{'c': true, 'b': true}
		mockLexicon := NewMockLexicon(ctrl)
		mockLexicon.EXPECT().ValidLettersBetweenPrefixAndSuffix(
			"",
			"at",
		).Times(1).Return(expectedCrossCheckSet)
		crossCheckSet, score := CrossCheck(tile, tiles, mockLexicon)
		assert.Equal(t, expectedCrossCheckSet, crossCheckSet)
		assert.Equal(t, 2, score)
	})
	t.Run("test prefix and suffix", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tiles := newBoard()
		tiles[0][0].Letter = 'c'
		tiles[2][0].Letter = 't'
		tile := tiles[1][0]
		expectedCrossCheckSet := set.RuneMap{'a': true}
		mockLexicon := NewMockLexicon(ctrl)
		mockLexicon.EXPECT().ValidLettersBetweenPrefixAndSuffix(
			"c",
			"t",
		).Times(1).Return(expectedCrossCheckSet)
		crossCheckSet, score := CrossCheck(tile, tiles, mockLexicon)
		assert.Equal(t, expectedCrossCheckSet, crossCheckSet)
		assert.Equal(t, 5, score)
	})
}

func newBoard() Board {
	return Board{
		{NewTile(0, 0, 1, 1), NewTile(0, 1, 1, 1), NewTile(0, 2, 1, 1), NewTile(0, 3, 1, 1)},
		{NewTile(1, 0, 1, 1), NewTile(1, 1, 1, 1), NewTile(1, 2, 1, 1), NewTile(1, 3, 1, 1)},
		{NewTile(2, 0, 1, 1), NewTile(2, 1, 1, 1), NewTile(2, 2, 1, 1), NewTile(2, 3, 1, 1)},
		{NewTile(3, 0, 1, 1), NewTile(3, 1, 1, 1), NewTile(3, 2, 1, 1), NewTile(3, 3, 1, 1)},
	}
}
