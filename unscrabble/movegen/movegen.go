package movegen

import "example.com/unscrabble/unscrabble/model"

type triePrefixGenerator struct {
	trieMoveGenerator *TrieMoveGenerator
}

type triePrefixExtender struct {
	trie *TrieMoveGenerator
}

type TrieMoveGenerator struct {
	trie          *Trie
	rack          model.Rack
	candidateTile model.Tile
	moveWord      []rune
	moveBlanks    []bool
}

func (t *TrieMoveGenerator) GenerateMoves(model.Board) []*model.Move {
	moves := []*Move{}
	for y, row := range board {
		for x, tile := range row {
			if tile.IsAnchor {
				moves = append(moves, t.generateMovesForAnchor(tile)...)
			}
		}
	}
}

func (t *TrieMoveGenerator) generateMovesForAnchor(anchor model.Tile) []*model.Move {
	moves := []*Move{}
	for _, prefix := range t.generatePrefixes(anchor) {
		moves = append(moves, t.extendPrefix(prefix)...)
	}
	return moves
}

func (t TrieMoveGenerator) generatePrefixes(anchor model.Tile) []string {
	placedPrefixChars := make([]rune, 0)
	for adjTile := currBoardTile; adjTile != nil && adjTile.Letter != 0; adjTile = adjTile.GetAdjacentTile(board, left, 0) {
		placedPrefixChars = append(placedPrefixChars, currBoardTile.Letter)
	}

	// Extend the placed prefix if it exists
	if len(placedPrefixChars) > 0 {
		return []string{string(placedPrefixChars)}
	}

	// otherwise generate and extend all valid prefixes
	lexi.GenerateNodesWithPruning(
		rack.Contains,
		fromRackShiftTile,
		toRackShiftTileBack,
		untilAnchorOrEdge,
		extendPrefix,
	)
}

func (t TrieMoveGenerator) takeBlankOrLetterFromRack(letter rune) rune {}

func (t TrieMoveGenerator) nextCandidateTile(direction int) {
	t.candidateTile = t.candidateTile.GetAdjacentTile(t.board, 0, direction)
}

func (t TrieMoveGenerator) addLetterOrBlankToMove(letter rune) {
	if t.rack.HasTile('*') {
		t.rack.RemoveRune('*')
		t.moveBlanks[len(t.moveWord)-1] = true
		return
	}
	t.rack.RemoveRune(letter)
}

func (t TrieMoveGenerator) toRackShiftTileBack(edgeChar rune, nextNode Lexicon) {
	if currBoardTile.Letter == 0 {
		if blanks[len(nextNode.Label())-1] {
			rack.AddRune('*')
			blanks[len(nextNode.Label())-1] = false
		} else {
			rack.AddRune(edgeChar)
		}
	}
	currBoardTile = currBoardTile.GetAdjacentTile(board, 0, -currDirection)
}

func (t triePrefixGenerator) IsEdgeValid(edge rune) bool {
	return t.trieMoveGenerator.Rack[edge]
}
func (t triePrefixGenerator) PreExpandHook(edge rune) {
	t.trieMoveGenerator.addBlankOrLetterToMove(edge)
	t.trieMoveGenerator.nextCandidateTile(left)
}
func (t triePrefixGenerator) PostExpandHook(edge rune) {

}
func (t triePrefixGenerator) IsFinished(node *Node) {}

func (t triePrefixGenerator) IsFinished(node *Node) {}
