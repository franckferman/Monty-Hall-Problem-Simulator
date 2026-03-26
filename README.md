<div id="top" align="center">

[![CI][ci-shield]](https://github.com/franckferman/Monty-Hall-Problem-Simulator/actions/workflows/ci.yml)
[![Contributors][contributors-shield]](https://github.com/franckferman/Monty-Hall-Problem-Simulator/graphs/contributors)
[![Forks][forks-shield]](https://github.com/franckferman/Monty-Hall-Problem-Simulator/network/members)
[![Stargazers][stars-shield]](https://github.com/franckferman/Monty-Hall-Problem-Simulator/stargazers)
[![License][license-shield]](https://github.com/franckferman/Monty-Hall-Problem-Simulator/blob/stable/LICENSE)
[![Go Version][go-shield]](https://go.dev/)
[![Platform][platform-shield]](https://github.com/franckferman/Monty-Hall-Problem-Simulator/releases)

<a href="https://github.com/franckferman/Monty-Hall-Problem-Simulator">
  <img src="https://raw.githubusercontent.com/franckferman/Monty-Hall-Problem-Simulator/refs/heads/stable/docs/github/graphical_resources/Logo-Monty-Hall-Problem-Simulator.png" alt="Monty Hall Problem Simulator" width="auto" height="auto">
</a>

<h3 align="center">Monty Hall Problem Simulator</h3>
<p align="center">
  <em>Should you switch or stay? One million simulations later, the answer is still the same.</em>
  <br><br>
  A probability simulator built around one of the most counterintuitive results in mathematics.<br>
  Interactive web interface and full-featured CLI, with the complete mathematical derivation.
</p>

</div>

---

## Table of Contents

<details open>
  <summary><strong>Click to collapse/expand</strong></summary>
  <ol>
    <li><a href="#the-problem">The Problem</a></li>
    <li><a href="#the-mathematics">The Mathematics</a></li>
    <li><a href="#quick-start">Quick Start</a></li>
    <li><a href="#web-interface">Web Interface</a></li>
    <li><a href="#command-line-interface">Command-Line Interface</a></li>
    <li><a href="#build-from-source">Build from Source</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#star-evolution">Star Evolution</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>

---

## The Problem

In 1963, a game show called *Let's Make a Deal* began airing on American television. Contestants faced three closed doors. Behind one was a car. Behind the other two, goats. The host — Monty Hall — always knew where the car was.

The game went like this: you pick a door. Monty then opens one of the other doors, always revealing a goat. He asks: do you want to switch to the remaining door, or stick with your first choice?

Most people say it doesn't matter. 50/50, right? Two doors left, one car. In 1990, Marilyn vos Savant published the correct answer in *Parade* magazine — **always switch, you win twice as often** — and received nearly 10,000 letters from mathematicians and academics telling her she was wrong.

She wasn't.

This simulator exists to make that result undeniable. Run it once with 10 trials and the noise might fool you. Run it a million times and the numbers converge exactly where the math says they should.

---

## The Mathematics

### Why intuition fails

When Monty opens a door, most people unconsciously apply what psychologists call the *equal likelihood* heuristic: two doors remain, so each must have a 50% chance. This reasoning is wrong because it ignores the nature of Monty's action.

Monty does not open a door at random. He has perfect knowledge and always reveals a goat. This constraint makes his action informative — he is not generating new uncertainty, he is removing one specific possibility in a way that is directly conditioned on where the car is. Once you account for that, the symmetry breaks.

### Conditional probability

You pick door 1. Monty opens door 3 (a goat). Two cases:

- **Staying** wins if and only if your initial pick was correct. That probability is 1/3. Nothing that happens afterwards changes that prior.
- **Switching** wins if and only if your initial pick was *wrong*. That probability is 2/3.

This is the entire argument, and it requires no formula. If you picked wrong (probability 2/3), the car is behind one of the other two doors. Monty eliminates the one without the car. The remaining door therefore has the car. Switching wins.

### Bayesian derivation

Define C_i as the event "car is behind door i". You have picked door 1. Monty opens door 3. We want P(C_2 | host opens door 3).

By Bayes' theorem:

```
P(C_2 | host opens 3) = P(host opens 3 | C_2) × P(C_2) / P(host opens 3)
```

Computing each term:

```
P(C_2)                = 1/3

P(host opens 3 | C_1) = 1/2  (car at 1: host can open 2 or 3 freely)
P(host opens 3 | C_2) = 1    (car at 2: host cannot open 1 or 2, door 3 is the only option)
P(host opens 3 | C_3) = 0    (car at 3: host cannot reveal it)

P(host opens 3) = (1/2)(1/3) + (1)(1/3) + (0)(1/3) = 1/6 + 1/3 = 1/2
```

Substituting back:

```
P(C_2 | host opens 3) = (1 × 1/3) / (1/2) = 2/3   →  switch
P(C_1 | host opens 3) = (1/2 × 1/3) / (1/2) = 1/3  →  stay
```

The key term is `P(host opens 3 | C_2) = 1`. When the car is at door 2, Monty has no choice but to open door 3. This forced action is what shifts the probability mass. The host's lack of freedom is exactly what makes switching advantageous.

### Generalization to N doors

With N doors, you pick one, the host opens N−2 of the remaining doors (all goats), leaving exactly one alternative.

- **P(win | stay)** = 1/N — your initial pick carries only the original prior.
- **P(win | switch)** = (N−1)/N — you win by switching whenever your initial pick was wrong.

| N doors | P(win \| stay) | P(win \| switch) | Advantage |
|:-------:|:--------------:|:----------------:|:---------:|
| 3       | 33.33%         | 66.67%           | 2×        |
| 4       | 25.00%         | 75.00%           | 3×        |
| 5       | 20.00%         | 80.00%           | 4×        |
| 10      | 10.00%         | 90.00%           | 9×        |
| 100     | 1.00%          | 99.00%           | 99×       |

The N=100 case is instructive: you pick one door out of a hundred, the host opens 98 to reveal goats, leaving yours and one other. Nobody argues 50/50 at N=100. The three-door version is the same problem — only the numbers are less extreme, which is exactly why the intuition fails there.

### The role of host knowledge

The result depends entirely on the host knowing where the car is and always opening a goat door. If the host opened a door at random and happened to reveal a goat, the conditional probability calculation changes completely and the 50/50 intuition would actually be correct. The information asymmetry is the whole mechanism — Monty's constrained action communicates information about the car's location, and switching is the rational response to that information.

---

## Quick Start

Download a pre-built binary from the [releases page](https://github.com/franckferman/Monty-Hall-Problem-Simulator/releases).

| Platform | File |
|----------|------|
| Linux (x86_64) | `Monty-Hall-Problem-Simulator-Linux` |
| Windows (x86_64) | `Monty-Hall-Problem-Simulator-Windows.exe` |
| macOS (Intel) | `Monty-Hall-Problem-Simulator-MacOS` |
| macOS (Apple Silicon) | `Monty-Hall-Problem-Simulator-MacOS-ARM64` |

**Linux / macOS:**
```bash
chmod +x Monty-Hall-Problem-Simulator-Linux
./Monty-Hall-Problem-Simulator-Linux
```

**Windows:**
```
Monty-Hall-Problem-Simulator-Windows.exe
```

---

## Web Interface

```bash
./Monty-Hall-Problem-Simulator-Linux --web
```

Open **http://localhost:8080** in your browser.

The web interface runs simulations server-side and streams results to the browser in real time using Server-Sent Events. No page reloads, no polling.

**Live door demo** — Three doors animate continuously. The host reveals a goat, switch and stay outcomes are shown, then the cycle repeats. Switch to Play mode to make the decisions yourself and track your personal win rate over time.

**Real-time convergence chart** — Start a simulation and watch the switch and stay win rates draw themselves toward the theoretical predictions. With the instant speed setting, a million trials complete in under a second.

**N-door slider** — Drag from 3 to 10 doors. The chart, statistics, and mathematical explanation all update immediately.

**Math accordion** — The complete Bayesian derivation and N-door comparison table, expandable at the bottom of the page.

**CSV export** — Once a simulation completes, download the full dataset for use in Python, R, or Excel.

### Custom port

```bash
./Monty-Hall-Problem-Simulator-Linux --web --port 9000
```

---

## Command-Line Interface

Running the binary without flags executes the full suite of simulations (10 to 1,000,000 trials) and prints results with 95% confidence intervals and delta from the theoretical value.

```bash
./Monty-Hall-Problem-Simulator-Linux
```

```
 Results after 1,000,000 simulations:
  Stay:    33.27% (±0.09%)  | Theory: 33.33% | Delta: 0.06%
  Switch:  66.72% (±0.09%)  | Theory: 66.67% | Delta: 0.05%
  Gain by switching: 33.45%
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--web` | `false` | Launch the web interface |
| `--port` | `8080` | Port for the web interface |
| `--doors` | `3` | Number of doors (minimum 3) |
| `--trials` | `0` | Custom trial count (0 = full default suite) |
| `--graph` | `false` | Show ASCII convergence graph |
| `--output` | `text` | Output format: `text`, `json`, `csv` |
| `--parallel` | `false` | Use all CPU cores |
| `--verbose` | `false` | Step-by-step annotated game trace |
| `--no-math` | `false` | Skip the mathematical explanation |

### Examples

```bash
# Step-by-step game trace followed by the ASCII convergence graph
./Monty-Hall-Problem-Simulator-Linux --verbose --graph

# Simulate with 10 doors instead of 3
./Monty-Hall-Problem-Simulator-Linux --doors 10

# Run five million trials in parallel across all CPU cores
./Monty-Hall-Problem-Simulator-Linux --trials 5000000 --parallel

# Export results to CSV for Python or Excel
./Monty-Hall-Problem-Simulator-Linux --output csv > results.csv

# Export to JSON
./Monty-Hall-Problem-Simulator-Linux --output json | python3 -m json.tool

# Simulation only, no math section
./Monty-Hall-Problem-Simulator-Linux --no-math
```

### ASCII convergence graph

```bash
./Monty-Hall-Problem-Simulator-Linux --graph --no-math
```

```
 Convergence Graph
 Doors: 3  |  P(switch)=66.67%  |  P(stay)=33.33%
 S = Switch result  |  T = Stay result  |  --- = Theoretical

  100% |
   ...
   65% |--S-----------------------S----------S----------S-----S
   ...
   35% |--T-----------------------T----------T----------T-----T
   ...
    0% |
       +--------------------------------------------------------
        10         100        1K         10K        100K      1M
```

### CSV output

```csv
trials,stay_pct,switch_pct,stay_ci95,switch_ci95,theo_stay,theo_switch
10,30.0000,70.0000,28.4000,28.4000,33.3333,66.6667
100,35.0000,66.0000,9.3500,9.2800,33.3333,66.6667
1000000,33.2700,66.7200,0.0920,0.0920,33.3333,66.6667
```

---

## Build from Source

Requires [Go 1.21+](https://go.dev/dl/).

```bash
git clone https://github.com/franckferman/Monty-Hall-Problem-Simulator.git
cd Monty-Hall-Problem-Simulator
go build -o simulator .
./simulator
```

The web interface static files are embedded directly into the binary at build time using Go's `embed` package — no separate assets to ship or manage.

### Cross-platform build

```bash
chmod +x build.sh
./build.sh
```

Produces stripped, statically linked binaries in `bin/` for Linux, Windows, macOS Intel, and macOS ARM64. Release binaries are also built and published automatically via GitHub Actions on each version tag.

---

## Contributing

Issues and pull requests are welcome. The codebase is small and deliberately self-contained — the simulation logic, statistics, display, and HTTP server each live in their own package under `internal/`.

<p align="right">(<a href="#top">Back to top</a>)</p>

---

## Star Evolution

<a href="https://star-history.com/#franckferman/Monty-Hall-Problem-Simulator&Timeline">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=franckferman/Monty-Hall-Problem-Simulator&type=Timeline&theme=dark" />
    <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=franckferman/Monty-Hall-Problem-Simulator&type=Timeline" />
  </picture>
</a>

<p align="right">(<a href="#top">Back to top</a>)</p>

---

## License

GNU Affero General Public License v3.0. See [LICENSE](https://github.com/franckferman/Monty-Hall-Problem-Simulator/blob/stable/LICENSE) for details.

<p align="right">(<a href="#top">Back to top</a>)</p>

---

## Contact

[![ProtonMail][protonmail-shield]](mailto:contact@franckferman.fr)
[![LinkedIn][linkedin-shield]](https://www.linkedin.com/in/franckferman)
[![Twitter][twitter-shield]](https://www.twitter.com/franckferman)

<p align="right">(<a href="#top">Back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
[ci-shield]: https://github.com/franckferman/Monty-Hall-Problem-Simulator/actions/workflows/ci.yml/badge.svg?branch=stable
[contributors-shield]: https://img.shields.io/github/contributors/franckferman/Monty-Hall-Problem-Simulator.svg?style=for-the-badge
[forks-shield]: https://img.shields.io/github/forks/franckferman/Monty-Hall-Problem-Simulator.svg?style=for-the-badge
[stars-shield]: https://img.shields.io/github/stars/franckferman/Monty-Hall-Problem-Simulator.svg?style=for-the-badge
[license-shield]: https://img.shields.io/github/license/franckferman/Monty-Hall-Problem-Simulator.svg?style=for-the-badge
[go-shield]: https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white
[platform-shield]: https://img.shields.io/badge/Platform-Linux%20%7C%20Windows%20%7C%20macOS-lightgrey?style=for-the-badge
[protonmail-shield]: https://img.shields.io/badge/ProtonMail-8B89CC?style=for-the-badge&logo=protonmail&logoColor=blueviolet
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=blue
[twitter-shield]: https://img.shields.io/badge/-Twitter-black.svg?style=for-the-badge&logo=twitter&colorB=blue
