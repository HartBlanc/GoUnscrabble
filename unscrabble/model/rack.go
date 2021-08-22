package model

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

func NewRack(rackSize int) *Rack {
	return &Rack{
		letterCounts: map[rune]int{},
		letterSet:    map[rune]bool{},
		tileCount:    0,
		capacity:     rackSize,
	}
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

func (rack *Rack) AddRune(letter rune) {
	rack.letterCounts[letter]++
	rack.tileCount++
	if letter != '*' {
		rack.letterSet[letter] = true
	}
}

func (rack *Rack) RemoveRune(letter rune) {
	if rack.letterCounts[letter] == 0 {
		return
	}
	rack.letterCounts[letter]--
	rack.tileCount--
	if rack.letterCounts[letter] == 0 && letter != '*' {
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
func (rack *Rack) HasTile(tileLetter rune) bool {
	return rack.letterCounts[tileLetter] > 0
}

func (rack *Rack) FillRack(letterBag LetterBag) error {
	for rack.tileCount < rack.capacity && len(letterBag) > 0 {
		randomLetter, err := letterBag.PopRandomLetter()
		if err != nil {
			return err
		}
		rack.AddRune(randomLetter)
	}
	return nil
}
