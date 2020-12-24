package lexicon

// Node is a node in a lexicon data structure.
type Node struct {
	Label     string
	Terminal  bool
	NextNodes map[rune]*Node
}
