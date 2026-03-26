package display

import (
	"fmt"
	"math"
	"strings"

	"Monty-Hall-Problem-Simulator/internal/stats"
	"Monty-Hall-Problem-Simulator/internal/utils"
)

const (
	car    = "🚗"
	goat   = "🐐"
	closed = "🚪"
)

func PrintBanner() {
	fmt.Println("╔══════════════════════════════════════════════╗")
	fmt.Println("║    Monty Hall Problem Simulator  v2.0        ║")
	fmt.Println("╚══════════════════════════════════════════════╝")
}

func PrintIntroduction(numDoors int) {
	fmt.Println("\n Context:")
	fmt.Print("  Doors: ")
	for i := 0; i < numDoors; i++ {
		fmt.Printf("%s ", closed)
	}
	fmt.Printf("\n  One hides %s. The other %d hide %s\n", car, numDoors-1, goat)
	fmt.Println("  You pick a door. The host reveals a goat behind another.")
	fmt.Println("  Do you stick or switch?\n")
}

func PrintResult(r stats.Result) {
	stayDelta := math.Abs(r.StayPct - r.TheoStay)
	switchDelta := math.Abs(r.SwitchPct - r.TheoSwitch)
	fmt.Printf("\n Results after %s simulations:\n", FormatCount(r.TrialCount))
	fmt.Printf("  Stay:   %6.2f%% (+/-%.2f%%)  | Theory: %.2f%% | Delta: %.2f%%\n",
		r.StayPct, r.StayCI95, r.TheoStay, stayDelta)
	fmt.Printf("  Switch: %6.2f%% (+/-%.2f%%)  | Theory: %.2f%% | Delta: %.2f%%\n",
		r.SwitchPct, r.SwitchCI95, r.TheoSwitch, switchDelta)
	fmt.Printf("  Gain by switching: %.2f%%\n", r.SwitchPct-r.StayPct)
	fmt.Println("  " + strings.Repeat("-", 52))
}

// FormatCount formats an integer with thousands separators.
func FormatCount(n int) string {
	s := fmt.Sprintf("%d", n)
	out := make([]byte, 0, len(s)+4)
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, byte(c))
	}
	return string(out)
}

// PrintVerboseDemo shows a single annotated game for pedagogical purposes.
func PrintVerboseDemo(numDoors int) {
	fmt.Println("\n Step-by-step demo (one game):")
	fmt.Println("  " + strings.Repeat("-", 42))

	carPos := utils.SecureRandomInt(int64(numDoors))
	playerChoice := utils.SecureRandomInt(int64(numDoors))

	fmt.Print("  Doors: ")
	for i := 0; i < numDoors; i++ {
		fmt.Printf("%s ", closed)
	}
	fmt.Printf("\n  You pick door %d.\n", playerChoice+1)

	var montyChoice int
	for montyChoice = 0; montyChoice < numDoors; montyChoice++ {
		if montyChoice != playerChoice && montyChoice != carPos {
			break
		}
	}
	fmt.Printf("  Host opens door %d: %s (goat)\n", montyChoice+1, goat)

	remaining := utils.GetRemainingDoor(playerChoice, montyChoice, numDoors)
	fmt.Printf("  Your door: %d  |  Alternative: %d\n\n", playerChoice+1, remaining+1)

	if playerChoice == carPos {
		fmt.Printf("  STAY   -> door %d = %s  WIN\n", playerChoice+1, car)
		fmt.Printf("  SWITCH -> door %d = %s  LOSE\n", remaining+1, goat)
	} else {
		fmt.Printf("  STAY   -> door %d = %s  LOSE\n", playerChoice+1, goat)
		fmt.Printf("  SWITCH -> door %d = %s  WIN\n", remaining+1, car)
	}
	fmt.Println("  " + strings.Repeat("-", 42) + "\n")
}

// PrintConvergenceGraph renders a 2D ASCII chart showing how both strategies
// converge toward their theoretical win rates as trial count increases.
func PrintConvergenceGraph(results []stats.Result, numDoors int) {
	if len(results) == 0 {
		return
	}

	const (
		graphWidth  = 56
		graphHeight = 21
	)

	grid := make([][]byte, graphHeight)
	for i := range grid {
		grid[i] = make([]byte, graphWidth)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	pctToRow := func(pct float64) int {
		row := graphHeight - 1 - int(math.Round(pct/100.0*float64(graphHeight-1)))
		if row < 0 {
			row = 0
		}
		if row >= graphHeight {
			row = graphHeight - 1
		}
		return row
	}

	idxToCol := func(i int) int {
		if len(results) <= 1 {
			return graphWidth / 2
		}
		return i * (graphWidth - 1) / (len(results) - 1)
	}

	// Theoretical lines
	theoSwitchRow := pctToRow(results[0].TheoSwitch)
	theoStayRow := pctToRow(results[0].TheoStay)
	for x := 0; x < graphWidth; x++ {
		grid[theoSwitchRow][x] = '-'
		grid[theoStayRow][x] = '-'
	}

	// Data points (drawn after lines so they take priority)
	for i, r := range results {
		col := idxToCol(i)
		grid[pctToRow(r.SwitchPct)][col] = 'S'
		grid[pctToRow(r.StayPct)][col] = 'T'
	}

	fmt.Println("\n Convergence Graph")
	fmt.Printf(" Doors: %d  |  P(switch)=%.2f%%  |  P(stay)=%.2f%%\n",
		numDoors, results[0].TheoSwitch, results[0].TheoStay)
	fmt.Println(" S = Switch result  |  T = Stay result  |  --- = Theoretical\n")

	for row, line := range grid {
		pct := int(math.Round(float64(graphHeight-1-row) / float64(graphHeight-1) * 100))
		fmt.Printf("  %3d%% |%s\n", pct, string(line))
	}
	fmt.Printf("       +%s\n", strings.Repeat("-", graphWidth))

	// X-axis labels
	labelLine := make([]byte, graphWidth)
	for i := range labelLine {
		labelLine[i] = ' '
	}
	for i, r := range results {
		col := idxToCol(i)
		label := []byte(formatTrials(r.TrialCount))
		// Shift label left if it would overflow the right edge
		start := col
		if start+len(label) > graphWidth {
			start = graphWidth - len(label)
		}
		for j, c := range label {
			if start+j >= 0 && start+j < graphWidth {
				labelLine[start+j] = c
			}
		}
	}
	fmt.Printf("        %s\n", string(labelLine))
}

func formatTrials(n int) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%dM", n/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%dK", n/1_000)
	default:
		return fmt.Sprintf("%d", n)
	}
}

// PrintMathExplanation outputs a complete probability derivation including
// Bayesian proof and generalization to N doors.
func PrintMathExplanation(numDoors int) {
	theoStay := 1.0 / float64(numDoors) * 100
	theoSwitch := float64(numDoors-1) / float64(numDoors) * 100

	fmt.Println("\n" + strings.Repeat("=", 56))
	fmt.Println(" Mathematical Proof")
	fmt.Println(strings.Repeat("=", 56))

	fmt.Printf("\n P(win | stay)   = 1/%d = %.4f%%\n", numDoors, theoStay)
	fmt.Printf(" P(win | switch) = %d/%d = %.4f%%\n\n", numDoors-1, numDoors, theoSwitch)

	fmt.Println(" Core intuition:")
	fmt.Println("   You win by switching iff your initial pick was WRONG.")
	fmt.Printf("   P(initial pick = goat) = %d/%d\n", numDoors-1, numDoors)
	fmt.Println("   The host reveals all other goats, so the remaining door")
	fmt.Println("   absorbs the full probability mass of the initial group.\n")

	if numDoors == 3 {
		fmt.Println(" Bayesian derivation (N=3):")
		fmt.Println("   You pick door 1. Host opens door 3 (goat).")
		fmt.Println("   What is P(car = door 2 | host opened door 3)?\n")
		fmt.Println("   P(car=2 | H=3) = P(H=3 | car=2) * P(car=2) / P(H=3)\n")
		fmt.Println("   P(car=2)              = 1/3")
		fmt.Println("   P(H=3 | car=2)        = 1    (only option for host)")
		fmt.Println("   P(H=3 | car=1)        = 1/2  (host can open 2 or 3)")
		fmt.Println("   P(H=3 | car=3)        = 0    (host cannot reveal the car)\n")
		fmt.Println("   P(H=3) = (1/2)(1/3) + (1)(1/3) + (0)(1/3) = 1/6 + 2/6 = 1/2\n")
		fmt.Println("   P(car=2 | H=3) = (1 * 1/3) / (1/2) = 2/3\n")
		fmt.Println("   Switching wins with probability 2/3.")
		fmt.Println("   Staying wins with probability 1/3 = P(car=1 | H=3).\n")
	}

	fmt.Println(" Generalization to N doors:")
	fmt.Println("   P(win | stay)   = 1/N")
	fmt.Println("   P(win | switch) = (N-1)/N")
	fmt.Println("   Switching is (N-1)x more likely to win.\n")
	fmt.Println("   N=3:  2x better   (66.67% vs 33.33%)")
	fmt.Println("   N=4:  3x better   (75.00% vs 25.00%)")
	fmt.Println("   N=10: 9x better   (90.00% vs 10.00%)")
	fmt.Println("   N=100: 99x better (99.00% vs  1.00%)")
	fmt.Println(strings.Repeat("=", 56))
}
