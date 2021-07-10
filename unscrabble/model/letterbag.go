package model

import (
	"errors"
	"math/rand"
)

// LetterBag is an abstract data structure which allows for efficient random
// sampling without replacement. This is achieved by using stack-like item
// popping and shuffling the underlying array any time a new item is added.
type LetterBag []rune

// NewLetterBag is for constructing a new letterbag from counts of letters
func NewLetterBag(letterCounts map[rune]int) LetterBag {
	numLetters := 0
	for _, count := range letterCounts {
		numLetters += count
	}
	bag := make(LetterBag, 0, numLetters)
	bag.AddLetterCounts(letterCounts)
	return bag
}

// PopRandomLetter is for retrieving a random letter from the LetterBag
func (bag LetterBag) PopRandomLetter() (rune, error) {
	if len(bag) == 0 {
		return 0, errors.New("bag is empty")
	}

	randomLetter := bag[len(bag)-1]
	bag = bag[:len(bag)-1]
	return randomLetter, nil
}

// AddLetterCounts is for adding a batch of letters to the LetterBag
func (bag LetterBag) AddLetterCounts(letterCounts map[rune]int) {
	for letter, count := range letterCounts {
		for i := 0; i < count; i++ {
			bag = append(bag, letter)
		}
	}
	bag.shuffle()
}

func (bag LetterBag) shuffle() {
	rand.Shuffle(len(bag), func(i, j int) {
		bag[i], bag[j] = bag[j], bag[i]
	})
}
