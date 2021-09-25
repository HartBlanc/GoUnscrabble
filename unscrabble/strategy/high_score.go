package strategy

import (
	"example.com/unscrabble/unscrabble/model"
)

func NewHighScoreStrategy(MoveGenerator MoveGenerator) *HighScoreStrategy {
	return &HighScoreStrategy{
		moveGenerator: MoveGenerator,
	}
}

type HighScoreStrategy struct {
	moveGenerator MoveGenerator
}

// PickMove returns the highest scoring move out of all the moves generated by the provided board
// and rack. If multiple moves have the highest score, the first one provided by the generator is
// returned. If no moves are generated a nil Move is returned.
func (h *HighScoreStrategy) PickMove(board model.Board, rack model.Rack) *model.Move {
	var highScoreMove model.Move
	highestScore := -1

	for _, move := range h.moveGenerator.GenerateMoves(board, rack) {
		if move.Score > highestScore {
			highestScore = move.Score
			highScoreMove = move
		}
	}
	if highestScore == -1 {
		return nil
	}
	return &highScoreMove
}
