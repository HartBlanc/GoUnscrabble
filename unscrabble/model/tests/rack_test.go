package model_test

import (
	"errors"
	"testing"

	"example.com/unscrabble/unscrabble/model"
	"example.com/unscrabble/unscrabble/model/mock_model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFillFillsRackFromLetterGetter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLetterGetter := mock_model.NewMockLetterGetter(ctrl)

	mockLetterGetter.EXPECT().HasLetter().Return(true).Times(2)
	mockLetterGetter.EXPECT().GetLetter().Return('a', nil)
	mockLetterGetter.EXPECT().GetLetter().Return('b', nil)

	rack := model.NewRack(2)
	rack.Fill(mockLetterGetter)
	assert.True(t, rack.HasTile('a'))
	assert.True(t, rack.HasTile('b'))
}

func TestFillStopsEarlyIfLetterGetterDoesNotHaveAnyLetters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLetterGetter := mock_model.NewMockLetterGetter(ctrl)
	mockLetterGetter.EXPECT().HasLetter().Return(false)

	// no expected calls top GetLetter

	rack := model.NewRack(2)
	rack.Fill(mockLetterGetter)
}

func TestFillIgnoresErrorFromLetterGetter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLetterGetter := mock_model.NewMockLetterGetter(ctrl)

	mockLetterGetter.EXPECT().HasLetter().Return(true).Times(2)
	mockLetterGetter.EXPECT().GetLetter().Return('a', errors.New("an error!"))
	mockLetterGetter.EXPECT().GetLetter().Return('b', errors.New("another error!"))

	rack := model.NewRack(2)
	rack.Fill(mockLetterGetter)
	assert.True(t, rack.HasTile('a'))
	assert.True(t, rack.HasTile('b'))
}

func TestHasTileReturnsTrueIfRackContainsTile(t *testing.T) {
	rack := model.NewRack(1)
	rack.AddLetter('*')
	assert.True(t, rack.HasTile('*'))
}

func TestContainsReturnsTrueIfRackContainsBlankTile(t *testing.T) {
	rack := model.NewRack(1)
	rack.AddLetter('*')
	for _, letter := range []rune{'a', 'b', 'c', 'd', 'e'} {
		assert.True(t, rack.Contains(letter))
	}
}

func TestCopyCopiesRack(t *testing.T) {
	rack := model.NewRack(1)
	rack.AddLetter('a')

	copyRack := rack.Copy()
	assert.Equal(t, rack, &copyRack)
	assert.NotSame(t, rack, &copyRack)
}
