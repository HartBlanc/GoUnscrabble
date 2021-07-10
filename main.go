package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"example.com/unscrabble/unscrabble/model"
	"gopkg.in/yaml.v2"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// Load in confiugration
	// Create a game with that confifguration
	// Play that game
	// return the winner

	var config model.Configuration
	dataPath := os.Args[1]
	configBytes, err := ioutil.ReadFile(dataPath)
	check(err)

	err = yaml.Unmarshal(configBytes, &config)
	check(err)

	fmt.Println(config)

	_, err = model.NewGame(
		2,
		config.BingoPremium,
		config.RackSize,
		convertStringsToRunes(config.LetterScores),
		convertStringsToRunes(config.LetterCounts),
		config.LetterMultipliers,
		config.WordMultipliers,
		nil,
	)
}

func convertStringsToRunes(in map[string]int) map[rune]int {
	out := map[rune]int{}
	for key, value := range in {
		out[[]rune(key)[0]] = value
	}
	return out
}
