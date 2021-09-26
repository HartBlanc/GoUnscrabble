package model

// LetterGetter is for getting letters to fill the rack with
type LetterGetter interface {
	// GetLetter gets a letter, and should return an error if the getter does not have any letters
	GetLetter() (rune, error)
	HasLetter() bool
}

// NewRack is for creating a new empty rack
func NewRack(rackSize int) *Rack {
	return &Rack{
		letterCounts: map[rune]int{},
		letterSet:    map[rune]bool{},
		tileCount:    0,
		capacity:     rackSize,
	}
}

// The rack is an abstract data type which is essentially a multi-set.
// However, the rack also maintains a traditional set based on the presence
// of at least one count of the letter being present in the multiset.
// This allows for efficient set operations to be calculated with other sets of
// interest (e.g. cross-check sets and trie edge sets).
type Rack struct {
	letterCounts map[rune]int
	letterSet    map[rune]bool
	tileCount    int
	capacity     int
}

// Copy copies a rack
func (r Rack) Copy() Rack {

	copyLetterCounts := make(map[rune]int, len(r.letterCounts))
	for key, value := range r.letterCounts {
		copyLetterCounts[key] = value
	}
	copyLetterSet := make(map[rune]bool, len(r.letterSet))
	for key, value := range r.letterSet {
		copyLetterSet[key] = value
	}

	return Rack{
		letterCounts: copyLetterCounts,
		letterSet:    copyLetterSet,
		tileCount:    r.tileCount,
		capacity:     r.capacity,
	}
}

// AddLetter adds a new letter to the rack
func (rack *Rack) AddLetter(letter rune) {
	rack.letterCounts[letter]++
	rack.tileCount++
	rack.letterSet[letter] = true
}

// RemoveLetter removes an existing letter from the rack. It panics if the letter is not in the rack.
func (rack *Rack) RemoveLetter(letter rune) {
	if rack.letterCounts[letter] == 0 {
		panic("tried to remove a letter but it's not actually in the rack!")
	}
	rack.letterCounts[letter]--
	rack.tileCount--
	if rack.letterCounts[letter] == 0 {
		rack.letterSet[letter] = false
	}
}

// Contains tells us whether the rack contains a letter.
// If the rack contains a blank tile, it contains all letters.
func (rack *Rack) Contains(letter rune) bool {
	return rack.letterSet[letter] || rack.letterCounts['*'] > 0
}

// Has tile asks if the rack actually contains the tile.
// Wildcards are treated identically to all other tiles.
func (rack *Rack) HasTile(letter rune) bool {
	return rack.letterSet[letter]
}

// Fill fills the rack with tiles from a letterGetter
func (rack *Rack) Fill(letterGetter LetterGetter) {
	for rack.tileCount < rack.capacity && letterGetter.HasLetter() {
		letter, _ := letterGetter.GetLetter()
		rack.AddLetter(letter)
	}
}
