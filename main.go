package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"

	"Monty-Hall-Problem-Simulator/internal/display"
	"Monty-Hall-Problem-Simulator/internal/game"
	"Monty-Hall-Problem-Simulator/internal/server"
	"Monty-Hall-Problem-Simulator/internal/stats"
)

type jsonOutput struct {
	Config  jsonConfig     `json:"config"`
	Results []stats.Result `json:"results"`
}

type jsonConfig struct {
	Doors    int  `json:"doors"`
	Parallel bool `json:"parallel"`
}

func main() {
	web       := flag.Bool("web", false, "Launch the web GUI in your browser")
	port      := flag.Int("port", 8080, "Port for the web GUI (used with --web)")
	doors     := flag.Int("doors", 3, "Number of doors (minimum 3)")
	trials    := flag.Int("trials", 0, "Custom trial count (0 = default suite: 10 to 1M)")
	showGraph := flag.Bool("graph", false, "Show ASCII convergence graph after simulations")
	output    := flag.String("output", "text", "Output format: text | json | csv")
	parallel  := flag.Bool("parallel", false, "Use goroutines to parallelize simulations")
	verbose   := flag.Bool("verbose", false, "Show a step-by-step annotated game before simulating")
	noMath    := flag.Bool("no-math", false, "Skip the mathematical explanation at the end")
	flag.Parse()

	if *web {
		display.PrintBanner()
		if err := server.Start(*port); err != nil {
			fmt.Fprintln(os.Stderr, "Error starting web server:", err)
			os.Exit(1)
		}
		return
	}

	if *doors < 3 {
		fmt.Fprintln(os.Stderr, "Error: --doors must be >= 3")
		os.Exit(1)
	}
	if *output != "text" && *output != "json" && *output != "csv" {
		fmt.Fprintln(os.Stderr, "Error: --output must be one of: text, json, csv")
		os.Exit(1)
	}

	isText := *output == "text"

	if isText {
		display.PrintBanner()
		display.PrintIntroduction(*doors)
		if *parallel {
			fmt.Printf("  [Parallel mode: %d CPU cores]\n\n", runtime.NumCPU())
		}
	}

	if isText && *verbose {
		display.PrintVerboseDemo(*doors)
	}

	testCounts := []int{10, 100, 1_000, 10_000, 100_000, 1_000_000}
	if *trials > 0 {
		testCounts = []int{*trials}
	}

	results := make([]stats.Result, 0, len(testCounts))
	for _, tc := range testCounts {
		var stayPct, switchPct float64
		if *parallel {
			stayPct, switchPct = game.RunSimulationsParallel(tc, *doors)
		} else {
			stayPct, switchPct = game.RunSimulations(tc, *doors)
		}
		r := stats.Build(tc, *doors, stayPct, switchPct)
		results = append(results, r)
		if isText {
			display.PrintResult(r)
		}
	}

	switch *output {
	case "json":
		out := jsonOutput{
			Config:  jsonConfig{Doors: *doors, Parallel: *parallel},
			Results: results,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(out); err != nil {
			fmt.Fprintln(os.Stderr, "Error encoding JSON:", err)
			os.Exit(1)
		}
	case "csv":
		w := csv.NewWriter(os.Stdout)
		_ = w.Write([]string{"trials", "stay_pct", "switch_pct", "stay_ci95", "switch_ci95", "theo_stay", "theo_switch"})
		for _, r := range results {
			_ = w.Write([]string{
				fmt.Sprintf("%d", r.TrialCount),
				fmt.Sprintf("%.4f", r.StayPct),
				fmt.Sprintf("%.4f", r.SwitchPct),
				fmt.Sprintf("%.4f", r.StayCI95),
				fmt.Sprintf("%.4f", r.SwitchCI95),
				fmt.Sprintf("%.4f", r.TheoStay),
				fmt.Sprintf("%.4f", r.TheoSwitch),
			})
		}
		w.Flush()
		if err := w.Error(); err != nil {
			fmt.Fprintln(os.Stderr, "Error writing CSV:", err)
			os.Exit(1)
		}
	}

	if isText {
		if *showGraph {
			display.PrintConvergenceGraph(results, *doors)
		}
		if !*noMath {
			display.PrintMathExplanation(*doors)
		}
	}
}
