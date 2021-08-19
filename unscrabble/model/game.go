package model

import (
	"fmt"
)

type MovePicker interface {
	PickMove()
}

type Lexicon interface{}

// Player represents a single player in a Game
type Player struct {
	rack     *Rack
	turns    []*Move
	score    int
	strategy MovePicker
}

type Configuration struct {
	BingoPremium      int            `yaml:"bingo_premium"`
	RackSize          int            `yaml:"rack_size"`
	LetterScores      map[string]int `yaml:"letter_scores"`
	LetterCounts      map[string]int `yaml:"letter_counts"`
	LetterMultipliers [][]int        `yaml:"letter_multipliers"`
	WordMultipliers   [][]int        `yaml:"word_multipliers"`
}

// Game represents a single game
type Game struct {
	letterBag    LetterBag
	players      []Player
	board        Board
	lexicon      Lexicon
	letterScores map[rune]int
	bingoPremium int
}

// NewGame returns a new game using the provided configuration
func NewGame(
	numPlayers, bingoPremium, rackSize int,
	letterScores, letterCounts map[rune]int,
	letterMultipliers, wordMultipliers [][]int,
	lexicon Lexicon,
) (*Game, error) {

	// TODO: make a crosscheckset generator
	board := NewBoard(nil, wordMultipliers, letterMultipliers)
	letterBag := NewLetterBag(letterCounts)
	if (numPlayers * rackSize) < len(letterBag) {
		return nil, fmt.Errorf(
			"too many players (%v) for the rackSize (%v) and number of "+
				"letters in letter bag (%v)",
			numPlayers,
			rackSize,
			len(letterBag),
		)
	}

	var players []Player
	for i := 0; i < numPlayers; i++ {
		players = append(players, Player{rack: NewRack(rackSize)})
		players[i].rack.FillRack(letterBag)
	}

	game := Game{
		letterBag:    letterBag,
		players:      players,
		board:        board,
		lexicon:      lexicon,
		letterScores: letterScores,
		bingoPremium: bingoPremium,
	}
	return &game, nil
}

// Play a game to completition and return the winners
func (g *Game) Play() (winners []Player) {
	g.playAllTurns()
	return g.selectWinners()
}

func (g *Game) selectWinners() []Player {

	tiebreakWinnerScore := 0
	for _, player := range g.players {
		if player.score > tiebreakWinnerScore {
			tiebreakWinnerScore = player.score
		}
	}

	var tiebreakWinners []Player
	for _, player := range g.players {
		if player.score == tiebreakWinnerScore {
			tiebreakWinners = append(tiebreakWinners, player)
		}
	}

	winnerScore := 0
	for _, player := range g.players {

		for letter, count := range player.rack.letterCounts {
			player.score -= g.letterScores[letter] * count
		}

		if player.rack.tileCount == 0 {
			for _, otherPlayer := range g.players {
				for letter, count := range otherPlayer.rack.letterCounts {
					player.score += g.letterScores[letter] * count
				}
			}
		}

		if player.score > winnerScore {
			winnerScore = player.score
		}
	}

	var winners []Player
	for _, player := range g.players {
		if player.score == winnerScore {
			winners = append(winners, player)
		}
	}

	if len(winners) == 1 {
		return winners
	}

	return tiebreakWinners
}

func (g *Game) playAllTurns() {
	for anyPlayerHasMove := true; anyPlayerHasMove; anyPlayerHasMove = false {
		for _, player := range g.players {
			move := player.SelectMove(g)
			g.PerformMove(&player, move)
			if move == nil {
				player.ReplaceRack(g.letterBag)
				continue
			}
			anyPlayerHasMove = true

			player.ReplaceRack(g.letterBag)
			if player.rack.tileCount == 0 {
				return
			}
		}
	}
	return
}

func (g *Game) PerformMove(player *Player, move *Move) {}

func (p *Player) SelectMove(game *Game) *Move {
	return nil
}

func (p *Player) ReplaceRack(letterbag LetterBag) {}
