package strategy_test

import (
	"testing"

	"example.com/unscrabble/unscrabble/model"
	"example.com/unscrabble/unscrabble/strategy"
	"example.com/unscrabble/unscrabble/strategy/mock_strategy"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHighScoreStrategyPickMoveReturnsMoveWithHighestScore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMoveGenerator := mock_strategy.NewMockMoveGenerator(ctrl)

	generatedMoves := []model.Move{
		{
			StartPosition: &model.Position{Row: 0, Column: 1},
			Horizontal:    true,
			Word: model.Word{
				Chars:      "abc",
				BlankTiles: []bool{true, false, false},
			},
			Score: 9,
		},
		{
			StartPosition: &model.Position{Row: 0, Column: 1},
			Horizontal:    true,
			Word: model.Word{
				Chars:      "abcd",
				BlankTiles: []bool{true, false, false, false},
			},
			Score: 11,
		},
		{
			StartPosition: &model.Position{Row: 0, Column: 0},
			Horizontal:    false,
			Word: model.Word{
				Chars:      "ab",
				BlankTiles: []bool{false, false},
			},
			Score: 5,
		},
	}

	mockMoveGenerator.
		EXPECT().
		GenerateMoves(gomock.Any(), gomock.Any()).
		Return(generatedMoves)

	highScoreStrategy := strategy.NewHighScoreStrategy(mockMoveGenerator)
	expectedMove := generatedMoves[1]
	assert.Equal(t, &expectedMove, highScoreStrategy.PickMove(model.Board{}, model.Rack{}))
}

func TestHighScoreStrategyPickMoveReturnsFirstMoveWithHighestScore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMoveGenerator := mock_strategy.NewMockMoveGenerator(ctrl)

	generatedMoves := []model.Move{
		{
			StartPosition: &model.Position{Row: 0, Column: 1},
			Horizontal:    true,
			Word: model.Word{
				Chars:      "abc",
				BlankTiles: []bool{true, false, false},
			},
			Score: 11,
		},
		{
			StartPosition: &model.Position{Row: 0, Column: 1},
			Horizontal:    true,
			Word: model.Word{
				Chars:      "abcd",
				BlankTiles: []bool{true, false, false, false},
			},
			Score: 11,
		},
	}

	mockMoveGenerator.
		EXPECT().
		GenerateMoves(gomock.Any(), gomock.Any()).
		Return(generatedMoves)

	highScoreStrategy := strategy.NewHighScoreStrategy(mockMoveGenerator)
	expectedMove := generatedMoves[0]
	assert.Equal(t, &expectedMove, highScoreStrategy.PickMove(model.Board{}, model.Rack{}))
}
