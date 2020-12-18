package unscrabble

import "strings"

var (
	Empty        struct{}
	letterScores = map[rune]int{
		'a': 1,
		'b': 4,
		'c': 4,
		'd': 2,
		'e': 1,
		'f': 4,
		'g': 1,
		'h': 3,
		'i': 1,
		'j': 10,
		'k': 5,
		'l': 2,
		'm': 4,
		'n': 2,
		'o': 1,
		'p': 4,
		'q': 10,
		'r': 1,
		's': 1,
		't': 1,
		'u': 2,
		'v': 5,
		'w': 4,
		'x': 8,
		'y': 3,
		'z': 10,
	}
	wwfLetterMultipliers = [][]int{
		{3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 3, 1, 2, 1, 2, 1, 3, 1, 1},
		{1, 1, 1, 3, 1, 1, 1, 3, 1, 1, 1},
		{1, 1, 2, 1, 1, 1, 1, 1, 2, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 2, 1, 1, 1, 1, 1, 2, 1, 1},
		{1, 1, 1, 3, 1, 1, 1, 3, 1, 1, 1},
		{1, 1, 3, 1, 2, 1, 2, 1, 3, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3},
	}
	wwfWordMultipliers = [][]int{
		{1, 1, 3, 1, 1, 1, 1, 1, 3, 1, 1},
		{1, 2, 1, 1, 1, 2, 1, 1, 1, 2, 1},
		{3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3},
		{1, 2, 1, 1, 1, 2, 1, 1, 1, 2, 1},
		{1, 1, 3, 1, 1, 1, 1, 1, 3, 1, 1},
	}
)

// Tile is a data structure contains information relating to
// the tile on the board.
type Tile struct {
	Letter           rune
	WordMultiplier   int
	LetterMultiplier int
	CrossCheckSets   map[rune]struct{}
	BoardPosition    *Position
}

type Position struct {
	Row    int
	Column int
}

type Lexicon interface {
	Contains(string) bool
	ValidLettersBetweenPrefixAndSuffix(string, string) map[rune]struct{}
}

// NewTile returns a new empty tile, with a full cross check set.
func NewTile(x, y, wordMultiplier, letterMultiplier int) *Tile {
	return &Tile{
		WordMultiplier:   wordMultiplier,
		LetterMultiplier: letterMultiplier,
		BoardPosition: &Position{
			Row:    y,
			Column: x,
		},
	}
}

// NewBoard returns a new empty board (a 2D slice of Tiles) from 2D slices
// of word multipliers and letter multipliers.
func NewBoard(wordMultipliers, letterMultipliers [][]int) [][]*Tile {
	boardSize := len(wordMultipliers)
	tiles := make([][]*Tile, boardSize)
	for i := range tiles {
		tiles[i] = make([]*Tile, boardSize)
		for j := range tiles[i] {
			tiles[i][j] = NewTile(
				i,
				j,
				wordMultipliers[i][j],
				letterMultipliers[i][j],
			)
		}
	}
	return tiles
}

// Transpose transposes the tiles of the board.
// This is achieved using an in-place transformation.
// This works on the assumption that the board is square.
func Transpose(tiles [][]*Tile) {
	for i := range tiles {
		for j := i + 1; j < len(tiles); j++ {
			tiles[i][j], tiles[j][i] = tiles[j][i], tiles[i][j]
		}
	}
}

// GetAnchors finds the anchors of the rows.
// aka the candidate anchors of the words.
// These anchors are the empty squares which are adjacent
// (horizontally or vertically) to another square.
// TODO: get these incrementally?
func GetAnchors(tiles [][]*Tile) []*Tile {
	anchors := make([]*Tile, 0, len(tiles)*len(tiles))
	for i, row := range tiles {
		for j, tile := range row {
			if !(tile.Letter == 0) {
				continue
			}

			if i > 0 && tiles[i-1][j].Letter != 0 {
				anchors = append(anchors, tile)
				continue
			}

			if i < (len(tiles)-1) && tiles[i+1][j].Letter != 0 {
				anchors = append(anchors, tile)
				continue
			}

			if j > 0 && tiles[i][j-1].Letter != 0 {
				anchors = append(anchors, tile)
				continue
			}

			if j < (len(tiles)-1) && tiles[i][j+1].Letter != 0 {
				anchors = append(anchors, tile)
				continue
			}
		}

	}
	return anchors
}

// CrossCheck finds the cross check set and the cross score of a tile.
// The cross check set of a tile is the set of letters
// that will form legal down words when making an across
// move through that square. The cross score is the sum of
// the prefix score and the suffix score. Where the prefix score
// and the suffix scores are the sums of the scores of the letters
// of the prefix above the tile and the suffix below the tile.
func CrossCheck(tile *Tile, tiles [][]*Tile, lexicon Lexicon) (map[rune]struct{}, int) {
	prefix, prefixScore := GetPrefixAbove(tile, tiles)
	suffix, suffixScore := GetSuffixBelow(tile, tiles)
	if prefix == "" && suffix == "" {
		return nil, 0
	}
	return lexicon.ValidLettersBetweenPrefixAndSuffix(prefix, suffix), prefixScore + suffixScore
}

// GetPrefixAbove finds the prefix and the score
// associated with the consecutive tiles immediately
// above the provided tile.
func GetPrefixAbove(tile *Tile, tiles [][]*Tile) (string, int) {

	var sb strings.Builder
	x := tile.BoardPosition.Column
	y := tile.BoardPosition.Row - 1
	score := 0

	for ; (y >= 0) && tiles[y][x].Letter != 0; y-- {
		placedTile := tiles[y][x]
		sb.WriteRune(placedTile.Letter)
		score += placedTile.LetterMultiplier * letterScores[placedTile.Letter]
	}

	return reverse(sb.String()), score
}

// GetSuffixBelow finds the suffix and the score
// associated with the consecutive tiles immediately
// below the provided tile.
func GetSuffixBelow(tile *Tile, tiles [][]*Tile) (string, int) {

	var sb strings.Builder
	x := tile.BoardPosition.Column
	y := tile.BoardPosition.Row + 1
	score := 0

	for ; (y < len(tiles)) && tiles[y][x].Letter != 0; y++ {
		placedTile := tiles[y][x]
		sb.WriteRune(placedTile.Letter)
		score += placedTile.LetterMultiplier * letterScores[placedTile.Letter]
	}

	return sb.String(), score
}

// function, which takes a string as
// argument and return the reverse of string.
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
