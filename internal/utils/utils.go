package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
)

const (
	CAR    = "🚗"
	GOAT   = "🐐"
	CLOSED = "🚪"
)

func SecureRandomInt(limit int64) int {
	n, err := rand.Int(rand.Reader, big.NewInt(limit))
	if err != nil {
		log.Fatalf("[ERROR] Failed to generate secure random int: %v", err)
	}
	return int(n.Int64())
}

func GetRemainingDoor(choice, montyChoice, numDoors int) int {
	for i := 0; i < numDoors; i++ {
		if i != choice && i != montyChoice {
			return i
		}
	}
	return choice // fallback (should never hit)
}

func PrintIntroduction() {
	fmt.Println("\n📖 Context:")
	fmt.Println("Imagine you're on a game show. In front of you are three doors:", CLOSED, CLOSED, CLOSED)
	fmt.Println("Behind one is a shiny car", CAR, "and behind the others are goats", GOAT)
	fmt.Println("You pick a door. The host reveals a goat behind another door.")
	fmt.Println("Will you stick or switch? Let's explore the odds!\n")
}

func PrintResults(testCount int, stayPct, switchPct float64) {
	fmt.Printf("\n🔍 Results after %d simulations:\n", testCount)
	fmt.Printf("Staying: %.2f%% | Switching: %.2f%% | Gain by Switching: %.2f%%\n", stayPct, switchPct, switchPct-stayPct)
	fmt.Println("---------------------------------------------------")
}

func PrintExplanation(numDoors int) {
	fmt.Println("\n🔢 The Math Unraveled:")
	fmt.Printf("Initially, you have a 1/%d chance of picking the car %s and a %d/%d chance of picking a goat %s.\n", numDoors, CAR, numDoors-1, numDoors, GOAT)
	fmt.Println("Switching effectively gives you the better odds!")
}
