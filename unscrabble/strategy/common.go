package strategy

import "example.com/unscrabble/unscrabble/model"

type MoveGenerator interface {
	GenerateMoves(board model.Board, rack model.Rack) []model.Move
}
