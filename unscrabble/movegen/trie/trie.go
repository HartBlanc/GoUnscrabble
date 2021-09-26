package trie

import (
	"example.com/unscrabble/lexicon"
	"example.com/unscrabble/unscrabble/model"
)

func NewTrieMoveGenertator(trieRoot *lexicon.TrieNode) TrieMoveGenerator {
	return TrieMoveGenerator{trieRoot: trieRoot}
}

type TrieMoveGenerator struct {
	rack     model.Rack
	board    model.Board
	trieRoot *lexicon.TrieNode
}

func (t *TrieMoveGenerator) GenerateMoves(board model.Board, rack model.Rack) []model.Move {
	var moves []model.Move
	t.board = board
	t.rack = rack

	for _, transposed := range []bool{false, true} {
		for _, row := range board.Tiles {
			for _, tile := range row {
				if !tile.IsAnchor {
					continue
				}
				for _, prefixResult := range t.generatePrefixResults(tile) {
					for _, extendedPrefix := range t.extendPrefix(prefixResult, tile) {
						startPos := model.Position{
							Row:    tile.BoardPosition.Row - len(prefixResult.prefix.Chars),
							Column: tile.BoardPosition.Column,
						}
						if transposed {
							startPos.Row, startPos.Column = startPos.Column, startPos.Row
						}
						moves = append(
							moves,
							model.Move{
								StartPosition: &startPos,
								Horizontal:    transposed,
								Word:          extendedPrefix,
							},
						)
					}
				}
			}
		}
		model.Transpose(t.board)
	}
	return moves
}

func (t TrieMoveGenerator) generatePrefixResults(anchor *model.Tile) []partialPrefixResult {
	// Return the prefix that is already on the board if it exists
	placedPrefixChars := make([]rune, 0)
	for adjTile := anchor.GetAdjacentTile(t.board, 0, -1); adjTile != nil && !adjTile.Empty(); adjTile = adjTile.GetAdjacentTile(t.board, 0, -1) {
		placedPrefixChars = append(placedPrefixChars, adjTile.Letter)
	}
	placedPrefix := string(placedPrefixChars)
	if len(placedPrefix) > 0 {

		// blank *placed* tiles are not blank for the purpose of moves as we
		// don't need to use a blank tile from the rack
		noBlankTiles := make([]bool, len(placedPrefix))
		return []partialPrefixResult{
			{
				prefix: model.Word{
					Chars:      placedPrefix,
					BlankTiles: noBlankTiles,
				},
				remainingRack: t.rack.Copy(),
			}}
	}

	// otherwise generate all valid prefixes that can be placed from
	// the rack
	prefixGenerator := newPrefixResultGenerator(t.rack, anchor.BoardPosition.Column)
	t.trieRoot.VisitNodesWithPruning(prefixGenerator)

	return prefixGenerator.results
}

func (t TrieMoveGenerator) extendPrefix(prefixResult partialPrefixResult, anchor *model.Tile) []model.Word {
	prefixExtender := newPrefixExtender(
		t.board, prefixResult.remainingRack, prefixResult.prefix.BlankTiles, prefixResult.node, anchor,
	)
	prefixResult.node.VisitNodesWithPruning(prefixExtender)

	return prefixExtender.words
}

func newPrefixResultGenerator(rack model.Rack, maxPrefixLength int) *prefixResultGenerator {
	return &prefixResultGenerator{
		rack:            rack,
		prefixBlanks:    make([]bool, maxPrefixLength),
		maxPrefixLength: maxPrefixLength,
		results:         nil,
	}
}

type prefixResultGenerator struct {
	rack            model.Rack
	prefixBlanks    []bool
	maxPrefixLength int
	results         []partialPrefixResult
}

func (t *prefixResultGenerator) IsValidEdge(edge rune) bool {
	return t.rack.Contains(edge)
}

func (t *prefixResultGenerator) Terminate(node *lexicon.TrieNode) bool {
	return len(node.Label) >= t.maxPrefixLength
}

// Visit vists a TrieNode by removing a tile from the rack and adding the prefix
// to the set of results
func (t *prefixResultGenerator) Visit(node *lexicon.TrieNode) {
	if node.IsRoot() {
		return
	}
	if t.rack.HasTile(node.IncomingEdge()) {
		t.rack.RemoveLetter(node.IncomingEdge())
	} else if t.rack.HasTile('*') {
		t.rack.RemoveLetter('*')
		t.prefixBlanks[len(node.Label)-1] = true
	}

	prefixBlanks := make([]bool, len(node.Label), len(node.Label))
	for i := 0; i < len(node.Label); i++ {
		prefixBlanks[i] = t.prefixBlanks[i]
	}

	result := partialPrefixResult{
		prefix: model.Word{
			Chars:      node.Label,
			BlankTiles: prefixBlanks,
		},
		remainingRack: t.rack.Copy(),
		node:          node,
	}
	t.results = append(t.results, result)
}

// Exit cleans up after a TrieNode and all of its children have been visited by
// adding the tile used to visit the node back on the rack
func (t *prefixResultGenerator) Exit(node *lexicon.TrieNode) {
	if node.IsRoot() {
		return
	}
	if lastTileWasBlank := t.prefixBlanks[len(node.Label)-1]; lastTileWasBlank {
		t.rack.AddLetter('*')
		t.prefixBlanks[len(node.Label)-1] = false
		return
	}
	t.rack.AddLetter(node.IncomingEdge())
}

type partialPrefixResult struct {
	prefix        model.Word
	remainingRack model.Rack
	node          *lexicon.TrieNode
}

func newPrefixExtender(board model.Board, rack model.Rack, prefixBlanks []bool, prefixRoot *lexicon.TrieNode, anchor *model.Tile) *prefixExtender {
	blanks := make([]bool, len(board.Tiles), len(board.Tiles))
	for i := 0; i < len(prefixBlanks); i++ {
		blanks[i] = prefixBlanks[i]
	}
	return &prefixExtender{
		board:      board,
		rack:       rack,
		blanks:     blanks,
		prefixRoot: prefixRoot,
		currTile:   anchor,
	}
}

type prefixExtender struct {
	board      model.Board
	rack       model.Rack
	currTile   *model.Tile
	prefixRoot *lexicon.TrieNode
	blanks     []bool
	words      []model.Word
}

func (t *prefixExtender) IsValidEdge(edge rune) bool {
	if t.currTile.Empty() {
		return t.rack.Contains(edge) && (t.currTile.CrossCheckSet == nil || t.currTile.CrossCheckSet[edge])
	}
	return edge == t.currTile.Letter
}

func (t *prefixExtender) Terminate(node *lexicon.TrieNode) bool {
	return false // termination achieved via sentinels empty cross-set
}

// Visit vists a TrieNode by removing a tile from the rack and adding the prefix
// to the set of results
func (t *prefixExtender) Visit(node *lexicon.TrieNode) {
	if node == t.prefixRoot {
		return
	}

	nextTile := t.currTile.GetAdjacentTileOrSentinel(t.board, 0, 1)

	// note that the sentinel square is 'empty'
	if node.Terminal && t.currTile.Empty() {
		blanks := make([]bool, len(node.Label))
		for i := 0; i < len(node.Label); i++ {
			blanks[i] = t.blanks[i]
		}
		t.words = append(
			t.words,
			model.Word{
				Chars:      node.Label,
				BlankTiles: blanks,
			},
		)
	}

	if t.currTile.Empty() {
		if t.rack.HasTile(node.IncomingEdge()) {
			t.rack.RemoveLetter(node.IncomingEdge())
		} else if t.rack.HasTile('*') {
			t.rack.RemoveLetter('*')
			t.blanks[len(node.Label)-1] = true
		}
	}

	t.currTile = nextTile
}

// Exit cleans up after a TrieNode and all of its children have been visited by
// returning the currTile, rack, and blanks to the state they were in before the
// node was visited
func (t *prefixExtender) Exit(node *lexicon.TrieNode) {
	if node == t.prefixRoot {
		return
	}

	currTileHasLetter := t.currTile.Empty()
	t.currTile = t.currTile.GetAdjacentTile(t.board, 0, -1)

	if !currTileHasLetter {
		return
	}

	if lastTileWasBlank := t.blanks[len(node.Label)-1]; lastTileWasBlank {
		t.rack.AddLetter('*')
		t.blanks[len(node.Label)-1] = false
		return
	}

	t.rack.AddLetter(node.IncomingEdge())
}
