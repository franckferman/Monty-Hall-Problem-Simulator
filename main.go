package main

import (
	"fmt"
	"Monty-Hall-Problem-Simulator/internal/game"
	"Monty-Hall-Problem-Simulator/internal/utils"
)

func main() {
	fmt.Println("🎉 Welcome to the Monty Hall Problem Simulator! 🎉")
	utils.PrintIntroduction()

tests := []int{10, 100, 1000, 10000, 100000, 1000000}

	for _, testCount := range tests {
		stayPct, switchPct := game.RunSimulations(testCount, 3)
		utils.PrintResults(testCount, stayPct, switchPct)
	}

	utils.PrintExplanation(3)
}

