package main

import (
	"fmt"
	"strconv"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const outputGuessesPerLine = 5

const debugGuessLimit = -1 // Stop execution after 'x' guesses (-1 = unlimited)
const debugWord = ""       // Stop execution after processing defined word

// Test Case (solving for "shall")
// Guesses so far:
// | caper
// | logan
// | blast
// INPUT:
// var (
// 	Grays   = "cperognbt"
// 	Yellows = map[rune]string{
// 		'a': "24",	// 'a' has yellows in 2/4
// 		'l': "12",	// 'l' has yellows in 1/2
// 		's': "4",	// 's' has yellows in 4
// 	}
// 	Greens = []rune("--a--")
// )

var (
	Grays   = ""
	Yellows = map[rune]string{}
	Greens  = []rune("-----")
)

func main() {
	// Grays Utility
	graysMap := map[rune]struct{}{}
	for _, letter := range Grays {
		graysMap[letter] = struct{}{}
	}
	isGray := func(letter rune) bool {
		_, ok := graysMap[letter]
		return ok
	}

	// Yellows Utility
	yellowsMap := map[rune]map[rune]struct{}{}
	for letter, positions := range Yellows {
		for _, pos := range positions {
			yellowsMap[letter] = map[rune]struct{}{pos: {}}
		}
	}
	containsAllYellows := func(guess string) bool {
		for yellow := range Yellows {
			isPresent := false
			for _, letter := range guess {
				if letter == yellow {
					isPresent = true
				}
			}
			if isPresent == false {
				return isPresent
			}
		}
		return true
	}

	// Greens Utility
	greensMap := map[int]rune{}
	greensList := map[rune]struct{}{}
	LAST_INDEX := 0
	for i, letter := range Greens {
		if letter == '-' {
			LAST_INDEX = i
		} else {
			// Mark a green
			greensMap[i+1] = letter
			greensList[letter] = struct{}{}
		}
	}

	containsAllGreens := func(guess string) bool {
		for green := range greensList {
			isPresent := false
			for _, letter := range guess {
				if letter == green {
					isPresent = true
				}
			}
			if isPresent == false {
				return false
			}
		}
		return true
	}

	itor := func(num int) rune {
		return rune(strconv.Itoa(num)[0])
	}

	// Permutation engine
	guesses := []string{}

	tempGuess := []rune{}
	for _, letter := range Greens {
		tempGuess = append(tempGuess, letter)
	}

	debugAbort := false

	var StartGuessFrom func(c int)
	StartGuessFrom = func(guessIndex int) {
		if debugAbort || guessIndex > LAST_INDEX {
			return
		}

		if debugGuessLimit >= 0 && len(guesses) >= debugGuessLimit {
			return
		}

		// If green, skip current index
		if Greens[guessIndex] != '-' {
			StartGuessFrom(guessIndex + 1)
			return
		}

		for _, currentAlpha := range alphabet {
			// Check if current alpha is gray
			if isGray(currentAlpha) {
				// Skip impossible guess
				continue
			}

			// Check if current alpha is a yellow
			if impossibles, ok := yellowsMap[currentAlpha]; ok {
				// Check if current alpha is impossible using yellows info
				if _, impossible := impossibles[itor(guessIndex+1)]; impossible {
					// Skip impossible guess
					continue
				}
			}
			tempGuess[guessIndex] = currentAlpha

			// Recurse to next index
			StartGuessFrom(guessIndex + 1)

			// Abort post-recursion if marked
			if debugAbort {
				return
			}

			// Add guess if it matches criteria for greens/yellows
			if guessIndex == LAST_INDEX {
				guess := string(tempGuess)
				if containsAllGreens(guess) && containsAllYellows(guess) {
					guesses = append(guesses, guess)
				}
				if guess == debugWord {
					debugAbort = true
				}
				if debugGuessLimit >= 0 && len(guesses) >= debugGuessLimit {
					return
				}
			}
		}
	}

	// EXECUTE
	fmt.Println("START")
	StartGuessFrom(0)

	// PRINT
	i := 0
	for _, guess := range guesses {
		if i >= outputGuessesPerLine {
			fmt.Println()
			i = 0
		}
		fmt.Printf("\t%s", guess)
		i++
	}
}
