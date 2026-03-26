package game

import (
	"runtime"
	"sync"

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

// RunSimulations runs testCount trials sequentially for both strategies.
func RunSimulations(testCount, numDoors int) (float64, float64) {
	stayWins, switchWins := 0, 0
	for i := 0; i < testCount; i++ {
		if montyHallSimulation(false, numDoors) {
			stayWins++
		}
		if montyHallSimulation(true, numDoors) {
			switchWins++
		}
	}
	return pct(stayWins, testCount), pct(switchWins, testCount)
}

// RunSimulationsParallel splits trials across all available CPU cores.
func RunSimulationsParallel(testCount, numDoors int) (float64, float64) {
	numWorkers := runtime.NumCPU()
	if numWorkers > testCount {
		numWorkers = testCount
	}

	type result struct{ stay, switch_ int }
	ch := make(chan result, numWorkers)

	var wg sync.WaitGroup
	base := testCount / numWorkers
	for w := 0; w < numWorkers; w++ {
		n := base
		if w == numWorkers-1 {
			n = testCount - base*(numWorkers-1)
		}
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			s, sw := 0, 0
			for i := 0; i < n; i++ {
				if montyHallSimulation(false, numDoors) {
					s++
				}
				if montyHallSimulation(true, numDoors) {
					sw++
				}
			}
			ch <- result{s, sw}
		}(n)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	totalStay, totalSwitch := 0, 0
	for r := range ch {
		totalStay += r.stay
		totalSwitch += r.switch_
	}
	return pct(totalStay, testCount), pct(totalSwitch, testCount)
}

func pct(wins, total int) float64 {
	return float64(wins) * 100.0 / float64(total)
}
