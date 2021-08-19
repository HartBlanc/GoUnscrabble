package movegen

import (
	"example.com/unscrabble/lexicon"
	"example.com/unscrabble/unscrabble/model"
)

type triePrefixGenerator struct {
	trieMoveGenerator *TrieMoveGenerator
}

type triePrefixExtender struct {
	trie *TrieMoveGenerator
}

type TrieMoveGenerator struct {
	trie          *lexicon.Trie
	board         model.Board
	rack          model.Rack
	candidateTile *model.Tile
	moveBlanks    []bool
}

func (t *TrieMoveGenerator) GenerateMoves(board model.Board) []*model.Move {
	t.board = board

	t.moveBlanks = nil
	for i := 0; i < len(board.Tiles); i++ {
		t.moveBlanks = append(t.moveBlanks, false)
	}

	var moves []*model.Move
	for _, row := range board.Tiles {
		for _, tile := range row {
			if tile.IsAnchor {
				moves = append(moves, t.generateMovesForAnchor(tile)...)
			}
		}
	}
	return moves
}

func (t *TrieMoveGenerator) generateMovesForAnchor(anchor *model.Tile) []*model.Move {
	moves := []*model.Move{}
	for _, prefix := range t.generatePrefixes(anchor) {
		moves = append(moves, t.extendPrefix(prefix)...)
	}
	return moves
}

func (t TrieMoveGenerator) generatePrefixes(anchor *model.Tile) []string {
	placedPrefixChars := make([]rune, 0)
	for adjTile := t.candidateTile; adjTile != nil && adjTile.Letter != 0; adjTile = adjTile.GetAdjacentTile(t.board, 0, -1) {
		placedPrefixChars = append(placedPrefixChars, t.candidateTile.Letter)
	}

	// Extend the placed prefix if it exists
	if len(placedPrefixChars) > 0 {
		return []string{string(placedPrefixChars)}
	}

	var prefixes []string
	appendToPrefixes := func(node *lexicon.Trie) {
		prefixes = append(prefixes, node.Label)
	}

	// otherwise generate and extend all valid prefixes
	t.trie.GenerateNodesWithPruning(
		t.rack.Contains,
		t.addLetterOrBlankToMove,
		t.removeLetterOrBlankFromMove,
		t.prefixExtendedToAnchorOrEdge,
		appendToPrefixes,
	)

	return prefixes
}

func (t TrieMoveGenerator) addLetterOrBlankToMove(letter rune, node *lexicon.Trie) {
	if t.rack.HasTile('*') {
		t.rack.RemoveRune('*')
		t.moveBlanks[len(node.Label)-1] = true
		return
	}
	t.rack.RemoveRune(letter)
}

func (t TrieMoveGenerator) removeLetterOrBlankFromMove(letter rune, node *lexicon.Trie) {
	if lastTileWasBlank := t.moveBlanks[len(node.Label)-1]; lastTileWasBlank {
		t.rack.AddRune('*')
		t.moveBlanks[len(node.Label)-1] = false
		return
	}
	t.rack.AddRune(letter)
}

func (t TrieMoveGenerator) prefixExtendedToAnchorOrEdge(node *lexicon.Trie) bool {
	currTile := t.candidateTile.GetAdjacentTile(t.board, 0, len(node.Label))
	if currTile == nil {
		return true
	}
	return currTile.IsAnchor && node.Label != ""
}
