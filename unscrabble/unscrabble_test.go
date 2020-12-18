package unscrabble

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

var testTiles = [][]*Tile{}

func TestTranspose(t *testing.T) {

	t.Run("transpose empty", func(t *testing.T) {
		tiles := [][]*Tile{{}}
		Transpose(tiles)
		assert.Equal(t, [][]*Tile{{}}, tiles)
	})
	t.Run("transpose single", func(t *testing.T) {
		tiles := [][]*Tile{
			{&Tile{}},
		}
		expectedTiles := [][]*Tile{
			{&Tile{}},
		}
		Transpose(tiles)
		assert.Equal(t, expectedTiles, tiles)
	})
	t.Run("transpose square", func(t *testing.T) {
		tiles := [][]*Tile{
			{&Tile{Letter: 'a'}, &Tile{Letter: 'b'}},
			{&Tile{Letter: 'c'}, &Tile{Letter: 'd'}},
		}
		expectedTiles := [][]*Tile{
			{&Tile{Letter: 'a'}, &Tile{Letter: 'c'}},
			{&Tile{Letter: 'b'}, &Tile{Letter: 'd'}},
		}
		Transpose(tiles)
		assert.Equal(t, expectedTiles, tiles)
		Transpose(tiles)

		expectedTiles = [][]*Tile{
			{&Tile{Letter: 'a'}, &Tile{Letter: 'b'}},
			{&Tile{Letter: 'c'}, &Tile{Letter: 'd'}},
		}
		assert.Equal(t, expectedTiles, tiles)
	})
}

func TestGetAnchors(t *testing.T) {

	t.Run("no tiles", func(t *testing.T) {
		tiles := [][]*Tile{{}}
		expectedAnchors := []*Tile{}
		assert.Equal(t, expectedAnchors, GetAnchors(tiles))
	})
	t.Run("all tiles empty", func(t *testing.T) {
		tiles := [][]*Tile{
			{&Tile{}, &Tile{}},
			{&Tile{}, &Tile{}},
		}
		expectedAnchors := []*Tile{}
		assert.Equal(t, expectedAnchors, GetAnchors(tiles))
	})
	t.Run("no empty tiles", func(t *testing.T) {
		tiles := [][]*Tile{
			{&Tile{Letter: 'a'}, &Tile{Letter: 'a'}},
			{&Tile{Letter: 'a'}, &Tile{Letter: 'a'}},
		}
		expectedAnchors := []*Tile{}
		assert.Equal(t, expectedAnchors, GetAnchors(tiles))
	})
	t.Run("adjacent tiles are anchors", func(t *testing.T) {

		above := &Tile{BoardPosition: &Position{Row: 0, Column: 1}}
		left := &Tile{BoardPosition: &Position{Row: 1, Column: 0}}
		right := &Tile{BoardPosition: &Position{Row: 1, Column: 2}}
		below := &Tile{BoardPosition: &Position{Row: 2, Column: 1}}
		tiles := [][]*Tile{
			{&Tile{}, above, &Tile{}},
			{left, &Tile{Letter: 'a'}, right},
			{&Tile{}, below, &Tile{}},
		}
		expectedAnchors := []*Tile{above, left, right, below}
		assert.ElementsMatch(t, expectedAnchors, GetAnchors(tiles))
	})
}

func TestGetPrefixAbove(t *testing.T) {
	newBoard := func() [][]*Tile {
		return [][]*Tile{
			{NewTile(0, 0, 1, 1), NewTile(1, 0, 1, 1), NewTile(2, 0, 1, 1), NewTile(3, 0, 1, 1)},
			{NewTile(0, 1, 1, 1), NewTile(1, 1, 1, 1), NewTile(2, 1, 1, 1), NewTile(3, 1, 1, 1)},
			{NewTile(0, 2, 1, 1), NewTile(1, 2, 1, 1), NewTile(2, 2, 1, 1), NewTile(3, 2, 1, 1)},
			{NewTile(0, 3, 1, 1), NewTile(1, 3, 1, 1), NewTile(2, 3, 1, 1), NewTile(3, 3, 1, 1)},
		}
	}
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

// TODO: here, and in prefix, check multipliers / score calculations.
func TestGetSuffixBelow(t *testing.T) {
	newBoard := func() [][]*Tile {
		return [][]*Tile{
			{NewTile(0, 0, 1, 1), NewTile(1, 0, 1, 1), NewTile(2, 0, 1, 1), NewTile(3, 0, 1, 1)},
			{NewTile(0, 1, 1, 1), NewTile(1, 1, 1, 1), NewTile(2, 1, 1, 1), NewTile(3, 1, 1, 1)},
			{NewTile(0, 2, 1, 1), NewTile(1, 2, 1, 1), NewTile(2, 2, 1, 1), NewTile(3, 2, 1, 1)},
			{NewTile(0, 3, 1, 1), NewTile(1, 3, 1, 1), NewTile(2, 3, 1, 1), NewTile(3, 3, 1, 1)},
		}
	}
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
		_, score := GetPrefixAbove(tiles[0][1], tiles)
		assert.Equal(t, 1+4, score)
	})
}
