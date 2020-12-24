package unscrabble

import (
	"errors"
	"example.com/unscrabble/lexicon"
	"strings"
)

var (
	letterScores = map[rune]int{
		'a': 1,
		'b': 4,
		'c': 4,
		'd': 2,
		'e': 1,
		'f': 4,
		'g': 1,
		'h': 3,
		'i': 1,
		'j': 10,
		'k': 5,
		'l': 2,
		'm': 4,
		'n': 2,
		'o': 1,
		'p': 4,
		'q': 10,
		'r': 1,
		's': 1,
		't': 1,
		'u': 2,
		'v': 5,
		'w': 4,
		'x': 8,
		'y': 3,
		'z': 10,
	}
	wwfLetterMultipliers = [][]int{
		{3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 3, 1, 2, 1, 2, 1, 3, 1, 1},
		{1, 1, 1, 3, 1, 1, 1, 3, 1, 1, 1},
		{1, 1, 2, 1, 1, 1, 1, 1, 2, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 2, 1, 1, 1, 1, 1, 2, 1, 1},
		{1, 1, 1, 3, 1, 1, 1, 3, 1, 1, 1},
		{1, 1, 3, 1, 2, 1, 2, 1, 3, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3},
	}
	wwfWordMultipliers = [][]int{
		{1, 1, 3, 1, 1, 1, 1, 1, 3, 1, 1},
		{1, 2, 1, 1, 1, 2, 1, 1, 1, 2, 1},
		{3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3},
		{1, 2, 1, 1, 1, 2, 1, 1, 1, 2, 1},
		{1, 1, 3, 1, 1, 1, 1, 1, 3, 1, 1},
	}
)

const (
	bingoPremium = 35
	rackSize     = 7
	left         = -1
	right        = 1
	above        = -1
	below        = 1
)

// BoardTile is a data structure which contains information relating
// to a possibly empty tile on the board.
type BoardTile struct {
	Letter           rune // If Letter is 0 the tile is empty
	WordMultiplier   int
	LetterMultiplier int           // If the LetterMultiplier is 0 the tile was a blank RackTile
	CrossCheckSet    map[rune]bool // If the CrossCheckSet is nil, any tile can be placed
	CrossScore       int
	IsAnchor         bool
	BoardPosition    *Position
}

type Move struct {
	StartPosition *Position
	Horizontal    bool // true is horizontal, false is vertical
	Word          string
	BlankTiles    []bool
	Score         int
}

type Rack struct {
	letterTiles []rune
	blanks      int
}

type Position struct {
	Row    int
	Column int
}

type Board [][]*BoardTile

func (position *Position) Transpose() {
	position.Row, position.Column = position.Column, position.Row
}

func NewTile(y, x, wordMultiplier, letterMultiplier int) *BoardTile {
	return &BoardTile{
		WordMultiplier:   wordMultiplier,
		LetterMultiplier: letterMultiplier,
		BoardPosition: &Position{
			Row:    y,
			Column: x,
		},
	}
}

// NewBoard returns a new empty board (a 2D slice of Tiles) from 2D slices
// of word multipliers and letter multipliers.
func NewBoard(wordMultipliers, letterMultipliers [][]int) Board {
	boardSize := len(wordMultipliers)
	board := make(Board, boardSize)
	for y := range board {
		board[y] = make([]*BoardTile, boardSize)
		for x := range board[y] {
			board[y][x] = NewTile(
				y,
				x,
				wordMultipliers[y][x],
				letterMultipliers[y][x],
			)
		}
	}
	return board
}

// TODO: check what we need to do about board modifications
// Transpose transposes the tiles of the board.
// This is achieved using an in-place transformation.
// This works on the assumption that the board is square.
func Transpose(board Board) {
	for i := range board {
		for j := i + 1; j < len(board); j++ {
			board[i][j], board[j][i] = board[j][i], board[i][j]
		}
	}
}

// GetAnchors finds the anchors of the rows.
// aka the candidate anchors of the words.
// These anchors are the empty squares which are adjacent
// (horizontally or vertically) to another square.
// TODO: get these incrementally?
func GetAnchors(board Board) []*BoardTile {
	anchors := make([]*BoardTile, 0, len(board)*len(board))
	for y, row := range board {
		for x, tile := range row {
			if !(tile.Letter == 0) {
				continue
			}

			if y > 0 && board[y-1][x].Letter != 0 {
				anchors = append(anchors, tile)
				continue
			}

			if y < (len(board)-1) && board[y+1][x].Letter != 0 {
				anchors = append(anchors, tile)
				continue
			}

			if x > 0 && board[y][x-1].Letter != 0 {
				anchors = append(anchors, tile)
				continue
			}

			if x < (len(board)-1) && board[y][x+1].Letter != 0 {
				anchors = append(anchors, tile)
				continue
			}
		}

	}
	return anchors
}

// CrossCheck finds the cross check set and the cross score of a tile.
// The cross check set of a tile is the set of letters
// that will form legal down words when making an across
// move through that square. The cross score is the sum of
// the prefix score and the suffix score. Where the prefix score
// and the suffix scores are the sums of the scores of the letters
// of the prefix above the tile and the suffix below the tile.
func CrossCheck(tile *BoardTile, board Board, trie lexicon.Node) (map[rune]bool, int) {
	prefix, prefixScore := GetPrefixAbove(tile, board)
	suffix, suffixScore := GetSuffixBelow(tile, board)
	if prefix == "" && suffix == "" {
		return nil, 0
	}
	return trie.ValidLettersBetweenPrefixAndSuffix(prefix, suffix), prefixScore + suffixScore
}

// GetPrefixAbove finds the prefix and the score
// associated with the consecutive tiles immediately
// above the provided tile.
func GetPrefixAbove(tile *BoardTile, board Board) (string, int) {

	var sb strings.Builder
	x := tile.BoardPosition.Column
	y := tile.BoardPosition.Row - 1
	score := 0

	for ; (y >= 0) && board[y][x].Letter != 0; y-- {
		placedTile := board[y][x]
		sb.WriteRune(placedTile.Letter)
		score += letterScores[placedTile.Letter]
	}

	return reverse(sb.String()), score
}

// GetSuffixBelow finds the suffix and the score
// associated with the consecutive tiles immediately
// below the provided tile.
func GetSuffixBelow(tile *BoardTile, board Board) (string, int) {

	var sb strings.Builder
	x := tile.BoardPosition.Column
	y := tile.BoardPosition.Row + 1
	score := 0

	for ; (y < len(board)) && board[y][x].Letter != 0; y++ {
		placedTile := board[y][x]
		sb.WriteRune(placedTile.Letter)
		score += letterScores[placedTile.Letter]
	}

	return sb.String(), score
}

func (move *Move) CalculateScore(board Board) (int, error) {
	y := move.StartPosition.Row
	x := move.StartPosition.Column

	if len(move.BlankTiles) != len(move.Word) {
		return 0, errors.New("blanks should be same length as word")
	}

	if len(move.Word) > (len(board) - x) {
		return 0, errors.New("word extends beyond end of board")
	}

	crossScore := 0
	horizontalScore := 0
	horizontalWordMultiplier := 1
	tilesPlaced := 0

	for i, char := range word {
		tile := board[y][x+i]
		if move.BlankTiles[i] {
			letterScore := 0
		} else {
			letterScore := letterScores[char] * tile.LetterMultiplier
		}

		horizontalScore += letterScore

		if tile.Letter == 0 {
			horizontalWordMultiplier *= tile.WordMultiplier
			if tile.CrossCheckSet != nil {
				crossScore += (tile.CrossScore + letterScore) * tile.WordMultiplier
			}
			tilesPlaced += 1
		}
	}
	horizontalScore *= horizontalWordMultiplier
	score := horizontalScore + crossScore
	if tilesPlaced == rackSize {
		score += bingoPremium
	}
	return score, nil
}

func reverse(s string) string {
	rns := []rune(s) // convert to rune
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {

		// swap the letters of the string,
		// like first with last and so on.
		rns[i], rns[j] = rns[j], rns[i]
	}

	// return the reversed string.
	return string(rns)
}

// TODO: allow passing generator function as an argument somehow (abstracting the lexicon)
// TODO: how to handle snchor skipping for GADDAG?
func GetLegalMoves(board Board, rack map[rune]int, anchors []*BoardTile, trie *lexicon.Node) []*Move {

	moves := make([]*Move, 0)

	transposed := false

	for i := 0; i < 2; i++ {
		for anchor := range anchors {
			appendToMoves := func(prefix, word string, blanks []bool) {
				x := anchor.BoardPosition.Column - len(prefix)
				y := anchor.BoardPosition.Row
				move := &Move{
					StartPosition: &Position{Row: y, Column: x},
					Horizontal:    !transposed,
					Word:          word,
					BlankTiles:    blanks,
				}
				move.Score, _ = move.CalculateScore(board)
				if transposed {
					move.StartPosition.Transpose()
				}
				_ = append(moves, move)
			}
			GenerateWordsFromAnchorWithTrie(board, rack, anchor, trie, appendToMoves)
		}
		Transpose(board)
		transposed = true
	}
	return moves
}

func GenerateWordsFromAnchorWithTrie(board Board, rack map[rune]int, anchor *BoardTile, trie *lexicon.Node, processPrefixWordAndBlanks func(string, string, []bool)) {

	currDirection := left
	currBoardTile := anchor.GetAdjacentTile(board, 0, currDirection)
	inRack := func(edgeChar rune) bool {
		return rack[edgeChar] == 0 && rack[0] == 0
	}
	inRackAndCrossSet := func(edgeChar rune) bool {
		return inRack(edgeChar) && (currBoardTile.CrossCheckSet == nil || currBoardTile.CrossCheckSet[edgeChar])
	}

	blanks := make([]bool, len(board))
	fromRackShiftTile := func(edgeChar rune, nextNode *lexicon.Node) {
		if currBoardTile.Letter == 0 {
			if rack[edgeChar] == 0 {
				blanks[len(nextNode.Label)-1] = true
				rack[0]--
			} else {
				rack[edgeChar]--
			}
		}
		currBoardTile = anchor.GetAdjacentTile(board, 0, currDirection)
	}
	toRackShiftTileBack := func(edgeChar rune, nextNode *lexicon.Node) {
		if currBoardTile.Letter == 0 {
			if blanks[len(nextNode.Label)-1] {
				rack[0]++
				blanks[len(nextNode.Label)-1] = false
			} else {
				rack[edgeChar]++
			}
		}
		currBoardTile = anchor.GetAdjacentTile(board, 0, -currDirection)
	}

	untilEdge := func(node *lexicon.Node) bool {
		return currBoardTile == nil
	}
	untilAnchorOrEdge := func(node *lexicon.Node) bool {
		return untilEdge(node) || currBoardTile.IsAnchor
	}

	prefix := ""
	processWord := func(wordNode *lexicon.Node) {
		if !wordNode.Terminal {
			return
		}
		word := wordNode.Label
		resultBlanks := make([]bool, len(word))
		copy(blanks, resultBlanks)
		processPrefixWordAndBlanks(prefix, word, resultBlanks)
	}
	extendPrefix := func(prefixNode *lexicon.Node) {

		prefixBoardTile := currBoardTile
		currBoardTile = anchor
		currDirection = right

		prefix = prefixNode.Label
		prefixNode.GenerateNodesWithPruning(inRackAndCrossSet, fromRackShiftTile, toRackShiftTileBack, untilEdge, processWord)

		currBoardTile = prefixBoardTile
		currDirection = left
	}
	trie.GenerateNodesWithPruning(inRack, fromRackShiftTile, toRackShiftTileBack, untilAnchorOrEdge, extendPrefix)
}

func (tile *BoardTile) GetAdjacentTile(board Board, vertical, horizontal int) *BoardTile {
	if tile.BoardPosition.Column+horizontal < 0 || tile.BoardPosition.Column+horizontal >= len(board) {
		return nil
	}
	return board[tile.BoardPosition.Row+vertical][tile.BoardPosition.Column+horizontal]
}

// PerformMove places tiles on the board from the rack.
// Update cross sets / anchors
func PerformMove(move *Move, board Board, rack map[rune]int, trie lexicon.Node) {

	if !move.Horizontal {
		Transpose(board)
		move.StartPosition.Transpose()
	}
	y := move.StartPosition.Row
	x := move.StartPosition.Column

	if x > 0 {
		leftTile := BoardTile[y][x-1]
		if leftTile.Letter == 0 {
			leftTile.IsAnchor = true
		}
	}

	for i, char := range move.Word {
		currTile := board[y][x+i]
		if currTile.Letter != 0 {
			continue
		}
		if move.BlankTiles[i] {
			currTile.LetterMultiplier = 0
		} else {
			currTile.LetterMultiplier = 1
		}

		// the crossCheckSet is set to the placed character to ensure
		// lexicon traversals are constrained to the placed character
		// when considering new moves that pass through this board position.
		currTile.CrossCheckSet = map[rune]bool{char: true}
		currTile.Letter = char
		currTile.IsAnchor = false

		for i, direction := range []int{above, below} {
			adjTile := currTile.GetAdjacentTile(board, direction, 0)
			if adjTile != nil && adjTile.Letter == 0 {
				adjTile.IsAnchor = true
			}
		}
		UpdateAdjacentCrossCheckSets(tile, board, trie)
	}

	if (x + len(word)) < len(board) {
		rightTile := BoardTile[y][x+len(word)]
		if rightTile.Letter == 0 {
			rightTile.IsAnchor = true
		}
	}
	if !move.Horizontal {
		Transpose(board)
		move.StartPosition.Transpose()
	}
}

func UpdateAdjacentCrossCheckSets(tile *BoardTile, board Board, trie lexicon.Node) {
	currTile := tile
	for ; currTile != nil && currTile.Letter != 0; currTile = currTile.GetAdjacentTile(board, above, 0) {
	}
	if currTile != nil {
		currTile.CrossCheckSet, currTile.CrossScore = CrossCheck(currTile, board, trie)
	}

	currTile = tile
	for ; currTile != nil && currTile.Letter != 0; currTile = currTile.GetAdjacentTile(board, below, 0) {
	}

	if currTile != nil {
		currTile.CrossCheckSet, currTile.CrossScore = CrossCheck(currTile, board, trie)
	}
}

func PlayGame() {
	// make an empty board (with center square as initial anchor)
	// make trie from source lexicon
	// generate n racks from the letter bag
	// generate moves for the first player
	// pick the move with the highest score
	// play that move
	// iterate to next player
	// continue until game is complete
	// identify winner
}

// TODO: Implement play game
// TODO: Test transposing
// TODO: Test score calculation
// TODO: Test move generation
// TODO: Test move placing
// TODO: Test game playing

// TODO: Implement bitset representation for crossCheckSets and trie edges.
// TODO: Try out different representations for bitset (uint64, bitset, roaring)
// TODO: Abstract out set representation as an interface
// TODO: Compare performance between representations

// TODO: Optimise edge following (vs current pruning) in trie
// TODO: Implement DAWG
// TODO: Implement GADDAG
