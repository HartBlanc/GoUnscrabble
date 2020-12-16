package unscrabble

// Board is a data structure containing tiles.
// It has methods for placing tiles on the board.
type Board struct {
	Tiles [][]*Tile
}

// Tile is a data structure contains information relating to
// the tile on the board.
type Tile struct {
	Letter rune
	WordMultiplier int
	LetterMultiplier int
	CrossChecks []bool
	BoardPosition Position
}

type Position struct {
	Row int
	Column int
}

// Transpose transposes the tiles of the board.
// This is achieved using an in-place transformation.
// This works on the assumption that the board is square.
func Transpose(tiles [][]*Tile){
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
			if !(tile.Letter == 0){
				continue
			}

			if i > 0 && tiles[i-1][j].Letter != 0 {
				anchors = append(anchors, tile)
				continue
			}

			if i < (len(tiles)-1) && tiles[i+1][j].Letter != 0{
				anchors = append(anchors, tile)
				continue
			}

			if j > 0 && tiles[i][j-1].Letter != 0 {
				anchors = append(anchors, tile)
				continue
			}

			if j < (len(tiles)-1) && tiles[i][j+1].Letter != 0{
				anchors = append(anchors, tile)
				continue
			}
		}

	}
	return anchors
}
