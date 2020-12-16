package unscrabble

import (
	"testing"
	assert "github.com/stretchr/testify/assert"
)

func TestTranspose(t *testing.T){

	t.Run("transpose empty", func(t *testing.T){
		tiles := [][]*Tile{{}}
		Transpose(tiles)
		assert.Equal(t, [][]*Tile{{}}, tiles)
	})

	t.Run("transpose single", func(t *testing.T){
		tiles := [][]*Tile{
			{&Tile{}},
		}
		expectedTiles := [][]*Tile{
			{&Tile{}},
		}
		Transpose(tiles)
		assert.Equal(t, expectedTiles, tiles)
	})

	t.Run("transpose square", func(t *testing.T){
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

func TestGetAnchors(t *testing.T){

	t.Run("no tiles", func(t *testing.T){
		tiles := [][]*Tile{{}}
		expectedAnchors := []*Tile{}
		assert.Equal(t, expectedAnchors, GetAnchors(tiles))
	})
	t.Run("all tiles empty", func(t *testing.T){
		tiles := [][]*Tile{
			{&Tile{}, &Tile{}},
			{&Tile{}, &Tile{}},
		}
		expectedAnchors := []*Tile{}
		assert.Equal(t, expectedAnchors, GetAnchors(tiles))
	})
	t.Run("no empty tiles", func(t *testing.T){
		tiles := [][]*Tile{
			{&Tile{Letter: 'a'}, &Tile{Letter: 'a'}},
			{&Tile{Letter: 'a'}, &Tile{Letter: 'a'}},
		}
		expectedAnchors := []*Tile{}
		assert.Equal(t, expectedAnchors, GetAnchors(tiles))
	})
	t.Run("adjacent tiles are anchors", func(t *testing.T){

		above := &Tile{BoardPosition: Position{Row:0, Column:1}}
		left := &Tile{BoardPosition: Position{Row:1, Column:0}}
		right := &Tile{BoardPosition: Position{Row:1, Column:2}}
		below := &Tile{BoardPosition: Position{Row:2, Column:1}}
		tiles := [][]*Tile{
			{&Tile{}, above, &Tile{}},
			{left, &Tile{Letter: 'a'}, right},
			{&Tile{}, below, &Tile{}},
		}
		expectedAnchors := []*Tile{above, left, right, below}
		assert.ElementsMatch(t, expectedAnchors, GetAnchors(tiles))
	})
}
