package set

type RuneSet interface {
	AddRune(rune)
	RemoveRune(rune)
	Contains(rune) bool
	Intersection(RuneSet) RuneSet
}

type RuneMap map[rune]bool

func (runeMap RuneMap) AddRune(r rune) {
	runeMap[r] = true
}

func (runeMap RuneMap) RemoveRune(r rune) {
	delete(runeMap, r)
}

func (runeMap RuneMap) Contains(r rune) bool {
	return runeMap[r]
}

//TODO: Consider using shorter set to iterate over if possible
func (self RuneMap) Intersection(otherRuneSet RuneSet) RuneSet {
	var intersection RuneSet
	intersection = make(RuneMap, 0)

	for r := range self {
		if otherRuneSet.Contains(r) {
			intersection.AddRune(r)
		}
	}
	return intersection
}

func New(size int) RuneSet {
	return make(RuneMap, size)
}
