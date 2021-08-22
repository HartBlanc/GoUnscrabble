package model

import (
	"errors"
	"strings"
)

type CrossCheckSetGenerator interface {
	ValidLettersBetweenPrefixAndSuffix(prefix, suffix string) map[rune]bool
}

// Position contains the coordinates of a board Tile
type Position struct {
	Row    int // Row is the zero-indexed row number (top row is 0)
	Column int // Column is the zero-indexed column number (leftmost row is 0)
}

func (position *Position) transpose() {
	position.Row, position.Column = position.Column, position.Row
}

// Tile contains information relating to a (possibly empty) tile on the Board
type Tile struct {
	Letter                 rune // If Letter is 0 the tile is empty
	WordMultiplier         int
	LetterMultiplier       int           // If LetterMultiplier is 0 the tile was a blank RackTile
	CrossCheckSet          map[rune]bool // If CrossCheckSet is nil, any tile can be placed
	CrossScore             int
	transposeCrossCheckSet map[rune]bool
	transposeCrossScore    int
	IsAnchor               bool
	BoardPosition          *Position
}

// NewTile creates a new empty Tile
func NewTile(y, x, wordMultiplier, letterMultiplier int) *Tile {
	return &Tile{
		WordMultiplier:   wordMultiplier,
		LetterMultiplier: letterMultiplier,
		BoardPosition: &Position{
			Row:    y,
			Column: x,
		},
		CrossCheckSet: nil,
	}
}

func (tile *Tile) transpose() {
	tile.BoardPosition.transpose()
	tile.CrossCheckSet, tile.transposeCrossCheckSet = tile.transposeCrossCheckSet, tile.CrossCheckSet
	tile.CrossScore, tile.transposeCrossScore = tile.transposeCrossScore, tile.CrossScore
}

// TODO: Review should vertical even be a parameter here?
// If it should, shouldn't we be bounds checking?
// If not document why not.

// GetAdjacentTile gets the tile adjacent to tile in board.
func (tile *Tile) GetAdjacentTile(board Board, vertical, horizontal int) *Tile {
	currRow := tile.BoardPosition.Row
	currColumn := tile.BoardPosition.Column
	if currColumn+horizontal < 0 || currColumn+horizontal >= len(board.Tiles) {
		return nil
	}
	return board.Tiles[currRow+vertical][currColumn+horizontal]
}

// GetAdjacentTileOrSentinels gets the tile adjacent to tile in board and returns a
// sentinel tile if the tile is out of bounds
func (tile *Tile) GetAdjacentTileOrSentinel(board Board, vertical, horizontal int) *Tile {
	currRow := tile.BoardPosition.Row
	currColumn := tile.BoardPosition.Column
	if currColumn+horizontal < 0 || currColumn+horizontal >= len(board.Tiles) {
		sentinel := &Tile{
			CrossCheckSet: make(map[rune]bool),
			BoardPosition: &Position{
				Column: currColumn + horizontal,
				Row:    currRow + vertical,
			},
		}
		return sentinel
	}
	return board.Tiles[currRow+vertical][currColumn+horizontal]
}

// SetIsAnchor is used for setting the isAnchor property of a Tile safely
func (tile *Tile) SetIsAnchor(isAnchor bool, board Board) error {
	if isAnchor {
		return nil
	}

	if !tile.IsAnchor {
		return errors.New("tile is not an anchor, should not be resetting IsAnchor to false")
	}
	if tile.Letter == 0 {
		return errors.New(
			"tile is empty, letter should be placed before setting IsAnchor to false",
		)
	}

	// the crossCheckSet is set to the placed character to ensure
	// Lexicon traversals are constrained to the placed character
	// when considering new moves that pass through this board position.
	tile.CrossCheckSet = map[rune]bool{tile.Letter: true}
	tile.transposeCrossCheckSet = tile.CrossCheckSet
	tile.CrossScore = 0
	tile.transposeCrossScore = tile.CrossScore
	return nil
}

// SetIsAnchor is used for setting the isAnchor propert of a Tile safely
func (tile *Tile) crossCheck(board Board, letterScores map[rune]int) (map[rune]bool, int) {
	suffix, suffixScore := tile.getSuffixBelow(board, letterScores)
	prefix, prefixScore := tile.getPrefixAbove(board, letterScores)
	if prefix == "" && suffix == "" {
		return nil, 0
	}
	crossCheckSet := board.crossCheckSetGenerator.ValidLettersBetweenPrefixAndSuffix(prefix, suffix)
	return crossCheckSet, prefixScore + suffixScore
}

func (tile *Tile) getSuffixBelow(board Board, letterScores map[rune]int) (string, int) {

	var sb strings.Builder
	x := tile.BoardPosition.Column
	y := tile.BoardPosition.Row + 1
	score := 0

	for ; (y < len(board.Tiles)) && board.Tiles[y][x].Letter != 0; y++ {
		placedTile := board.Tiles[y][x]
		sb.WriteRune(placedTile.Letter)
		score += letterScores[placedTile.Letter]
	}

	return sb.String(), score
}

func (tile *Tile) getPrefixAbove(board Board, letterScores map[rune]int) (string, int) {

	var sb strings.Builder
	x := tile.BoardPosition.Column
	y := tile.BoardPosition.Row - 1
	score := 0

	for ; (y >= 0) && board.Tiles[y][x].Letter != 0; y-- {
		placedTile := board.Tiles[y][x]
		sb.WriteRune(placedTile.Letter)
		score += letterScores[placedTile.Letter]
	}

	return reverse(sb.String()), score
}

func reverse(s string) string {
	rns := []rune(s) // convert to rune
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {

		// swap the letters of the string,
		// like first with last and so on.
		rns[i], rns[j] = rns[j], rns[i]
	}

	// return the reversed string.
	return string(rns)
}

func (t *Tile) Empty() bool {
	return t.Letter == 0
}

// Board is a collection of Tiles
type Board struct {
	Tiles                  [][]*Tile
	crossCheckSetGenerator CrossCheckSetGenerator
}

// NewBoard returns a new empty board (a 2D slice of Tiles) from 2D slices
// of word multipliers and letter multipliers.
func NewBoard(crossCheckSetGenerator CrossCheckSetGenerator, wordMultipliers, letterMultipliers [][]int) Board {
	boardSize := len(wordMultipliers)
	tiles := make([][]*Tile, boardSize)
	for y := range tiles {
		tiles[y] = make([]*Tile, boardSize)
		for x := range tiles[y] {
			tiles[y][x] = NewTile(
				y,
				x,
				wordMultipliers[y][x],
				letterMultipliers[y][x],
			)
		}
	}
	board := Board{
		Tiles:                  tiles,
		crossCheckSetGenerator: crossCheckSetGenerator,
	}
	board.Tiles[boardSize/2][boardSize/2].IsAnchor = true
	return board
}

// Transpose transposes the tiles of the board.
// This is achieved using an in-place transformation.
// This works on the assumption that the board is square.
func Transpose(board Board) {
	for y := range board.Tiles {
		for x := y + 1; x < len(board.Tiles); x++ {
			board.Tiles[y][x].transpose()
			board.Tiles[x][y].transpose()
			board.Tiles[y][x], board.Tiles[x][y] = board.Tiles[x][y], board.Tiles[y][x]
		}
	}
	// Transpose also flips the cross sets so we need to do the diagonal too
	for y := range board.Tiles {
		board.Tiles[y][y].transpose()
	}
}

// GetAnchors is for finding the anchors of the rows. Anchors are the empty
// Tiles which are adjacent (horizontally or vertically) to a non-empty
// Tile.
func GetAnchors(board Board) []*Tile {
	anchors := make([]*Tile, 0, len(board.Tiles)*len(board.Tiles))
	for y, row := range board.Tiles {
		for x, tile := range row {
			if !(tile.Letter == 0) {
				continue
			}

			if y > 0 && board.Tiles[y-1][x].Letter != 0 {
				anchors = append(anchors, tile)
				continue
			}

			if y < (len(board.Tiles)-1) && board.Tiles[y+1][x].Letter != 0 {
				anchors = append(anchors, tile)
				continue
			}

			if x > 0 && board.Tiles[y][x-1].Letter != 0 {
				anchors = append(anchors, tile)
				continue
			}

			if x < (len(board.Tiles)-1) && board.Tiles[y][x+1].Letter != 0 {
				anchors = append(anchors, tile)
				continue
			}
		}

	}
	return anchors
}
