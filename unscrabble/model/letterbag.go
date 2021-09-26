package model

import (
	"errors"
	"math/rand"
)

// RandomLetterBag is an abstract data structure which allows for efficient random
// sampling without replacement. This is achieved by using stack-like item
// popping and shuffling the underlying array any time a new item is added.
type RandomLetterBag []rune

// NewRandomLetterBag is for constructing a new letterbag from counts of letters
func NewRandomLetterBag(letterCounts map[rune]int) RandomLetterBag {
	numLetters := 0
	for _, count := range letterCounts {
		numLetters += count
	}
	bag := make(RandomLetterBag, 0, numLetters)
	bag.addLetterCounts(letterCounts)
	return bag
}

func (bag *RandomLetterBag) addLetterCounts(letterCounts map[rune]int) {
	for letter, count := range letterCounts {
		for i := 0; i < count; i++ {
			*bag = append(*bag, letter)
		}
	}
	bag.shuffle()
}

func (bag *RandomLetterBag) shuffle() {
	rand.Shuffle(len(*bag), func(i, j int) {
		(*bag)[i], (*bag)[j] = (*bag)[j], (*bag)[i]
	})
}

// GetLetter is for retrieving a random letter from the RandomLetterBag
func (bag *RandomLetterBag) GetLetter() (rune, error) {
	if !bag.HasLetter() {
		return 0, errors.New("bag is empty")
	}

	randomLetter := (*bag)[len(*bag)-1]
	*bag = (*bag)[:len(*bag)-1]
	return randomLetter, nil
}

// HasLetter is for checking whether the bag is empty or not
func (bag *RandomLetterBag) HasLetter() bool {
	return len(*bag) != 0
}
