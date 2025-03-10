package game

import (
	"Monty-Hall-Problem-Simulator/internal/utils"
)

func montyHallSimulation(change bool, numDoors int) bool {
	carPosition := utils.SecureRandomInt(int64(numDoors))
	playerChoice := utils.SecureRandomInt(int64(numDoors))

	if change {
		var montyChoice int
		for montyChoice = 0; montyChoice < numDoors; montyChoice++ {
			if montyChoice != playerChoice && montyChoice != carPosition {
				break
			}
		}
		playerChoice = utils.GetRemainingDoor(playerChoice, montyChoice, numDoors)
	}
	return playerChoice == carPosition
}

func RunSimulations(testCount int, numDoors int) (float64, float64) {
	stayWins, switchWins := 0, 0
	for i := 0; i < testCount; i++ {
		if montyHallSimulation(false, numDoors) { stayWins++ }
		if montyHallSimulation(true, numDoors)  { switchWins++ }
	}
	return float64(stayWins) * 100 / float64(testCount), float64(switchWins) * 100 / float64(testCount)
}

