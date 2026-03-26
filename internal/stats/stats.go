package stats

import "math"

// Result holds the outcome of one simulation batch with statistical metadata.
type Result struct {
	TrialCount int     `json:"trial_count"`
	StayPct    float64 `json:"stay_pct"`
	SwitchPct  float64 `json:"switch_pct"`
	StayCI95   float64 `json:"stay_ci95"`
	SwitchCI95 float64 `json:"switch_ci95"`
	TheoStay   float64 `json:"theo_stay"`
	TheoSwitch float64 `json:"theo_switch"`
}

// TheoreticalStay returns the theoretical win rate (%) when staying.
func TheoreticalStay(numDoors int) float64 {
	return 100.0 / float64(numDoors)
}

// TheoreticalSwitch returns the theoretical win rate (%) when switching.
func TheoreticalSwitch(numDoors int) float64 {
	return 100.0 * float64(numDoors-1) / float64(numDoors)
}

// CI95 returns the half-width of a 95% confidence interval for a proportion.
// Uses the Wilson/normal approximation: 1.96 * sqrt(p(1-p)/n).
func CI95(pct float64, n int) float64 {
	p := pct / 100.0
	return 1.96 * math.Sqrt(p*(1-p)/float64(n)) * 100.0
}

// Build constructs a Result with all computed statistical fields.
func Build(trialCount, numDoors int, stayPct, switchPct float64) Result {
	return Result{
		TrialCount: trialCount,
		StayPct:    stayPct,
		SwitchPct:  switchPct,
		StayCI95:   CI95(stayPct, trialCount),
		SwitchCI95: CI95(switchPct, trialCount),
		TheoStay:   TheoreticalStay(numDoors),
		TheoSwitch: TheoreticalSwitch(numDoors),
	}
}
