package strategy

import "example.com/unscrabble/unscrabble/model"

type HighScoreStrategy struct{}

func (h *HighScoreStrategy) PickMove(board model.Board, rack model.Rack) *model.Move {
	return &model.Move{}
}
