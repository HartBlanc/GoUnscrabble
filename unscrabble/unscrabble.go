package unscrabble

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"example.com/unscrabble/set"
)

const (
	left  = -1
	right = 1
	above = -1
	below = 1
)

type Lexicon interface {
	ValidLettersBetweenPrefixAndSuffix(string, string) set.RuneSet
	GenerateNodesWithPruning(
		func(rune) bool,
		func(rune, Lexicon),
		func(rune, Lexicon),
		func(Lexicon) bool,
		func(Lexicon),
	)
	Terminal() bool
	Label() string
	FollowEdges(string) Lexicon
}

type Position struct {
	Row    int
	Column int
}

func (position *Position) Transpose() {
	position.Row, position.Column = position.Column, position.Row
}

// BoardTile is a data structure which contains information relating
// to a possibly empty tile on the board.
type BoardTile struct {
	Letter                 rune // If Letter is 0 the tile is empty
	WordMultiplier         int
	LetterMultiplier       int         // If the LetterMultiplier is 0 the tile was a blank RackTile
	CrossCheckSet          set.RuneSet // If the CrossCheckSet is nil, any tile can be placed
	CrossScore             int
	transposeCrossCheckSet set.RuneSet
	transposeCrossScore    int
	IsAnchor               bool
	BoardPosition          *Position
}

func NewTile(y, x, wordMultiplier, letterMultiplier int) *BoardTile {
	return &BoardTile{
		WordMultiplier:   wordMultiplier,
		LetterMultiplier: letterMultiplier,
		BoardPosition: &Position{
			Row:    y,
			Column: x,
		},
		CrossCheckSet: nil,
	}
}

func (tile *BoardTile) Transpose() {
	tile.BoardPosition.Transpose()
	tile.CrossCheckSet, tile.transposeCrossCheckSet = tile.transposeCrossCheckSet, tile.CrossCheckSet
	tile.CrossScore, tile.transposeCrossScore = tile.transposeCrossScore, tile.CrossScore
}

func (tile *BoardTile) GetAdjacentTile(board Board, vertical, horizontal int) *BoardTile {
	if tile.BoardPosition.Column+horizontal < 0 || tile.BoardPosition.Column+horizontal >= len(board) {
		return nil
	}
	return board[tile.BoardPosition.Row+vertical][tile.BoardPosition.Column+horizontal]
}

func (tile *BoardTile) SetIsAnchor(isAnchor bool, board Board, lexi Lexicon) error {
	if !isAnchor {
		if !tile.IsAnchor {
			return errors.New("tile is not an anchor, should not be resetting IsAnchor to false")
		}
		if tile.Letter == 0 {
			return errors.New("tile is empty, letter should be placed before setting IsAnchor to false")
		}

		// the crossCheckSet is set to the placed character to ensure
		// Lexicon traversals are constrained to the placed character
		// when considering new moves that pass through this board position.
		tile.CrossCheckSet = set.New(len(letterScores))
		tile.CrossCheckSet.AddRune(tile.Letter)
		tile.transposeCrossCheckSet = tile.CrossCheckSet
		tile.CrossScore = 0
		tile.transposeCrossScore = tile.CrossScore
	} else {
		tile.CrossCheckSet, tile.CrossScore = CrossCheck(tile, board, lexi)
	}
	return nil
}

// CrossCheck finds the cross check set and the cross score of a tile.
// The cross check set of a tile is the set of letters
// that will form legal down words when making an across
// move through that square. The cross score is the sum of
// the prefix score and the suffix score. Where the prefix score
// and the suffix scores are the sums of the scores of the letters
// of the prefix above the tile and the suffix below the tile.
func CrossCheck(tile *BoardTile, board Board, lexi Lexicon) (set.RuneSet, int) {
	prefix, prefixScore := GetPrefixAbove(tile, board)
	suffix, suffixScore := GetSuffixBelow(tile, board)
	if prefix == "" && suffix == "" {
		return nil, 0
	}
	return lexi.ValidLettersBetweenPrefixAndSuffix(prefix, suffix), prefixScore + suffixScore
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

type Board [][]*BoardTile

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
	board[boardSize/2][boardSize/2].IsAnchor = true
	return board
}

// Transpose transposes the tiles of the board.
// This is achieved using an in-place transformation.
// This works on the assumption that the board is square.
func Transpose(board Board) {
	for y := range board {
		for x := y + 1; x < len(board); x++ {
			board[y][x].Transpose()
			board[x][y].Transpose()
			board[y][x], board[x][y] = board[x][y], board[y][x]
		}
	}
	// Transpose also flips the cross sets so we need to do the diagonal too
	for y := range board {
		board[y][y].Transpose()
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

// TODO: allow passing generator function as an argument somehow (abstracting the Lexicon)
// TODO: how to handle anchor skipping for GADDAG?
func GetLegalMoves(board Board, rack *Rack, lexi Lexicon) []*Move {

	moves := make([]*Move, 0)
	transposed := false

	for i := 0; i < 2; i++ {
		for _, row := range board {
			for _, anchor := range row {
				if !anchor.IsAnchor {
					continue
				}
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
					moves = append(moves, move)
				}
				GenerateWordsFromAnchorWithTrie(board, rack, anchor, lexi, appendToMoves)
			}
		}
		Transpose(board)
		transposed = true
	}
	return moves
}

func GenerateWordsFromAnchorWithTrie(
	board Board,
	rack *Rack,
	anchor *BoardTile,
	lexi Lexicon,
	processPrefixWordAndBlanks func(string, string, []bool),
) {

	currDirection := left
	currBoardTile := anchor.GetAdjacentTile(board, 0, currDirection)
	inRackAndCrossSet := func(edgeChar rune) bool {
		return (rack.Contains(edgeChar) &&
			(currBoardTile.CrossCheckSet == nil ||
				currBoardTile.CrossCheckSet.Contains(edgeChar)))
	}

	blanks := make([]bool, len(board))
	fromRackShiftTile := func(edgeChar rune, nextNode Lexicon) {
		if currBoardTile.Letter == 0 {
			if rack.HasTile('*') {
				blanks[len(nextNode.Label())-1] = true
				rack.RemoveRune('*')
			} else {
				rack.RemoveRune(edgeChar)
			}
		}
		currBoardTile = currBoardTile.GetAdjacentTile(board, 0, currDirection)
	}
	toRackShiftTileBack := func(edgeChar rune, nextNode Lexicon) {
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

	untilEdge := func(node Lexicon) bool {
		return currBoardTile == nil
	}
	untilAnchorOrEdge := func(node Lexicon) bool {
		return untilEdge(node) || currBoardTile.IsAnchor
	}

	prefix := ""
	processWord := func(wordNode Lexicon) {
		if !wordNode.Terminal() {
			return
		}
		word := wordNode.Label()
		resultBlanks := make([]bool, len(word))
		copy(blanks, resultBlanks)
		processPrefixWordAndBlanks(prefix, word, resultBlanks)
	}
	extendPrefix := func(prefixNode Lexicon) {

		// Save the currBoardTile so that once we have extended the current
		// prefix to the right we can resume generating all valid prefixes from
		// the same position.
		prefixBoardTile := currBoardTile
		currBoardTile = anchor
		currDirection = right

		prefix = prefixNode.Label()
		prefixNode.GenerateNodesWithPruning(
			inRackAndCrossSet,
			fromRackShiftTile,
			toRackShiftTileBack,
			untilEdge,
			processWord,
		)

		currBoardTile = prefixBoardTile
		currDirection = left
	}

	placedPrefixChars := make([]rune, 0)
	for adjTile := currBoardTile; adjTile != nil && adjTile.Letter != 0; adjTile = adjTile.GetAdjacentTile(board, left, 0) {
		currBoardTile = adjTile
		placedPrefixChars = append(placedPrefixChars, currBoardTile.Letter)
	}

	// Extend the placed prefix if it exists
	if len(placedPrefixChars) > 0 {
		placedPrefixNode := lexi.FollowEdges(string(placedPrefixChars))
		extendPrefix(placedPrefixNode)
		return
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

// PerformMove places tiles on the board from the rack.
// Updates cross sets and anchors
// TODO: should return the new list of anchors
func PerformMove(move *Move, board Board, rack *Rack, lexi Lexicon) {

	if !move.Horizontal {
		Transpose(board)
		move.StartPosition.Transpose()
	}
	y := move.StartPosition.Row
	x := move.StartPosition.Column

	for i, char := range move.Word {
		currTile := board[y][x+i]
		if currTile.Letter != 0 {
			continue
		}
		if move.BlankTiles[i] {
			currTile.LetterMultiplier = 0
			rack.RemoveRune('*')
		} else {
			currTile.LetterMultiplier = 1
			rack.RemoveRune(char)
		}

		currTile.Letter = char
		err := currTile.SetIsAnchor(false, board, lexi)
		if err != nil {
			panic(err)
		}

		for _, direction := range [2]int{above, below} {
			adjEmptyTile := currTile
			for ; adjEmptyTile != nil && adjEmptyTile.Letter != 0; adjEmptyTile = adjEmptyTile.GetAdjacentTile(board, direction, 0) {
			}

			if adjEmptyTile != nil {
				adjEmptyTile.SetIsAnchor(true, board, lexi)
			}
		}
	}

	// The board is transposed here to ensure that the transposed
	// cross check is completed, using the placed prefix to the left and
	// suffix to the right
	Transpose(board)
	if x > 0 {
		leftTile := board[y][x-1]
		if leftTile.Letter == 0 {
			leftTile.SetIsAnchor(true, board, lexi)
		}
	}
	if (x + len(move.Word)) < len(board) {
		rightTile := board[y][x+len(move.Word)]
		rightTile.SetIsAnchor(true, board, lexi)
	}
	Transpose(board)

	if !move.Horizontal {
		Transpose(board)
		move.StartPosition.Transpose()
	}
}

type Move struct {
	StartPosition *Position
	Horizontal    bool // true is horizontal, false is vertical
	Word          string
	BlankTiles    []bool
	Score         int
}

// The letter bag is an abstract data structure which allows for
// efficient random sampling without replacement. This is achieved by using
// stack-like item popping and shuffling the underlying array
// any time a new item is added.
type LetterBag []rune

func newLetterBag(letterCounts map[rune]int) LetterBag {
	numLetters := 0
	for _, count := range letterCounts {
		numLetters += count
	}
	bag := make(LetterBag, 0, numLetters)
	bag.AddLetterCounts(letterCounts)
	return bag
}

func (bag LetterBag) shuffle() {
	rand.Shuffle(len(bag), func(i, j int) {
		bag[i], bag[j] = bag[j], bag[i]
	})
}

func (bag LetterBag) PopRandomLetter() (rune, error) {
	if len(bag) == 0 {
		return -1, errors.New("bag is empty!")
	}
	randomLetter := bag[len(bag)-1]
	bag = bag[:len(bag)-1]
	return randomLetter, nil
}

func (bag LetterBag) AddLetterCounts(letterCounts map[rune]int) {
	for letter, count := range letterCounts {
		for i := 0; i < count; i++ {
			bag = append(bag, letter)
		}
	}
	bag.shuffle()
}

// Plays a Game with n players and returns the winner
func PlayGame(numberOfPlayers int, lexi Lexicon) (int, error) {

	letterBag := newLetterBag(letterCounts)
	if (numberOfPlayers * rackSize) < len(letterBag) {
		return -1, fmt.Errorf(
			"too many players (%v) for the rackSize (%v) and number of letters in letter bag (%v)",
			numberOfPlayers,
			rackSize,
			len(letterBag),
		)
	}
	racks := make([]*Rack, numberOfPlayers)
	for player := range racks {
		racks[player] = newRack()
		err := racks[player].FillRack(letterBag)
		if err != nil {
			return -1, err
		}
	}

	anyPlayerHasMove := true
	scores := make([]int, numberOfPlayers)
	board := NewBoard(wwfWordMultipliers, wwfLetterMultipliers)

	for anyPlayerHasMove {
		anyPlayerHasMove = false

		for player, rack := range racks {
			playerMoves := GetLegalMoves(board, rack, lexi)
			if len(playerMoves) == 0 {
				newRack := newRack()
				err := newRack.FillRack(letterBag)
				if err != nil {
					return -1, err
				}
				letterBag.AddLetterCounts(racks[player].letterCounts)
				racks[player] = newRack
				continue
			}

			anyPlayerHasMove = true
			bestMove := playerMoves[0]
			for _, move := range playerMoves[1:] {
				if move.Score > bestMove.Score {
					bestMove = move
				}
			}
			PerformMove(bestMove, board, racks[player], lexi)
			scores[player] += bestMove.Score
			err := racks[player].FillRack(letterBag)
			if err != nil {
				return -1, err
			}
			if racks[player].tileCount == 0 && len(letterBag) == 0 {
				return player, nil // winner!
			}
		}
	}

	winner := 0
	highestScore := -1
	for player, score := range scores {
		if score > highestScore {
			highestScore = score
			winner = player
		}
	}
	return winner, nil
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
