package model

import (
	"errors"
)

// Move contains a single candidate word, and a position for that word, that a
// player can play.
type Move struct {
	StartPosition *Position
	Horizontal    bool // true is horizontal, false is vertical
	Word          Word
	Score         int
}

type Word struct {
	Chars      string
	BlankTiles []bool
}

func (move *Move) CalculateScore(board Board, letterScores map[rune]int, rackSize, bingoPremium int) (int, error) {
	y := move.StartPosition.Row
	x := move.StartPosition.Column

	if len(move.Word.BlankTiles) != len(move.Word.Chars) {
		return 0, errors.New("blanks should be same length as word")
	}

	if len(move.Word.Chars) > (len(board.Tiles) - x) {
		return 0, errors.New("word extends beyond end of board.tiles")
	}

	crossScore := 0
	horizontalScore := 0
	horizontalWordMultiplier := 1
	tilesPlaced := 0

	for i, char := range move.Word.Chars {
		tile := board.Tiles[y][x+i]
		letterScore := letterScores[char] * tile.LetterMultiplier
		if move.Word.BlankTiles[i] {
			letterScore = 0
		}
		horizontalScore += letterScore

		if tile.Letter == 0 {
			horizontalWordMultiplier *= tile.WordMultiplier
			if tile.CrossCheckSet != nil {
				crossScore += (tile.CrossScore + letterScore) * tile.WordMultiplier
			}
			tilesPlaced += 1
		}
	}
	horizontalScore *= horizontalWordMultiplier
	score := horizontalScore + crossScore
	if tilesPlaced == rackSize {
		score += bingoPremium
	}
	return score, nil
}
