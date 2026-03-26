package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
	"time"

	"Monty-Hall-Problem-Simulator/internal/game"
	"Monty-Hall-Problem-Simulator/internal/stats"
)

//go:embed static
var staticFiles embed.FS

// Batches defines the trial counts streamed progressively to the client.
// Keep in sync with TOTAL_BATCHES in app.js.
var Batches = []int{
	10, 25, 50, 100, 250, 500,
	1_000, 2_500, 5_000, 10_000, 25_000,
	50_000, 100_000, 250_000, 500_000, 1_000_000,
}

type simEvent struct {
	TrialCount int     `json:"trial_count"`
	StayPct    float64 `json:"stay_pct"`
	SwitchPct  float64 `json:"switch_pct"`
	StayCI95   float64 `json:"stay_ci95"`
	SwitchCI95 float64 `json:"switch_ci95"`
	TheoStay   float64 `json:"theo_stay"`
	TheoSwitch float64 `json:"theo_switch"`
	Done       bool    `json:"done"`
}

// Start launches the HTTP server on the given port.
func Start(port int) error {
	sub, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("embed static files: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(sub)))
	mux.HandleFunc("/simulate", handleSimulate)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("  Web interface: http://localhost%s\n\n", addr)
	return http.ListenAndServe(addr, mux)
}

func handleSimulate(w http.ResponseWriter, r *http.Request) {
	doors     := clamp(parseIntOr(r, "doors", 3), 3, 20)
	delay     := clamp(parseIntOr(r, "delay", 80), 0, 3000)
	maxTrials := parseIntOr(r, "max", 1_000_000)
	parallel  := parseIntOr(r, "parallel", 0) == 1

	// Only run batches up to maxTrials
	activeBatches := make([]int, 0, len(Batches))
	for _, tc := range Batches {
		if tc <= maxTrials {
			activeBatches = append(activeBatches, tc)
		}
	}
	if len(activeBatches) == 0 {
		activeBatches = []int{10}
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	for i, tc := range activeBatches {
		select {
		case <-ctx.Done():
			return
		default:
		}

		var stayPct, switchPct float64
		if parallel {
			stayPct, switchPct = game.RunSimulationsParallel(tc, doors)
		} else {
			stayPct, switchPct = game.RunSimulations(tc, doors)
		}
		res := stats.Build(tc, doors, stayPct, switchPct)

		ev := simEvent{
			TrialCount: res.TrialCount,
			StayPct:    res.StayPct,
			SwitchPct:  res.SwitchPct,
			StayCI95:   res.StayCI95,
			SwitchCI95: res.SwitchCI95,
			TheoStay:   res.TheoStay,
			TheoSwitch: res.TheoSwitch,
			Done:       i == len(activeBatches)-1,
		}

		data, _ := json.Marshal(ev)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()

		if delay > 0 {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(delay) * time.Millisecond):
			}
		}
	}
}

func parseIntOr(r *http.Request, key string, def int) int {
	v, err := strconv.Atoi(r.URL.Query().Get(key))
	if err != nil {
		return def
	}
	return v
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
