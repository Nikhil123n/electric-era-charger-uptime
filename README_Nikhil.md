# README_Nikhil.md
Charger Uptime Coding Challenge – Implementation Notes  
Author: Nikhil Jivraj Arethiya

## 1. Goal

This document explains how my solution is structured and how I approached the problem.  
The main README provided in the repository already describes the problem itself, so this file focuses on design, structure, and implementation choices.

---

## 2. High-Level Approach

1. **Parse once, validate early**  
   I read the entire input file, enforce all structural rules up front (sections, station definitions, charger ownership, report format), and stop immediately with `ERROR` if anything is inconsistent.

2. **Represent reports as time intervals**  
   Each availability record for a charger is stored as a simple `[start, end)` interval with a boolean `up` flag. This keeps the data model minimal and easy to reason about.

3. **Compute station uptime from merged “up” intervals**  
   For each station, I:
   - Determine the global time window covered by all its chargers.  
   - Collect all intervals where any charger is explicitly up.  
   - Sort and merge overlapping intervals.  
   - Sum the merged intervals to get the total “up” time.  
   - Compute `floor(100 * totalUp / totalWindow)` as the uptime percentage.

4. **Fail fast on invalid input**  
   Any malformed input or inconsistent relationships (for example, a charger belonging to two stations) immediately results in `ERROR`, as required by the challenge.

---

## 3. File Layout

The solution is split into four files to keep concerns separated and the code easier to read.

- **main.go**  
  Entry point.  
  - Reads the input path from the command line.  
  - Calls `parseInput` to build in-memory structures.  
  - Calls `computeAllStationUptimes` to compute each station’s uptime.  
  - Prints `<station_id> <uptime_pct>` for each station in ascending order.

- **parser.go**  
  Input parsing and validation.  
  - Handles section transitions: `[Stations]` and `[Charger Availability Reports]`.  
  - Populates three maps:
    - `stationToChargers: stationID → []chargerID`  
    - `chargerToStation: chargerID → stationID`  
    - `chargerReports: chargerID → []Report`  
  - Ensures:
    - No charger belongs to more than one station.  
    - No availability record references an undefined charger.  
    - Stations are defined before reports.

- **types.go**  
  Shared data structures and helpers.  
  - `Report` struct for availability intervals.  
  - `StationResult` struct for final outputs.  
  - `fail()` helper that prints `ERROR` and exits.  
  - Small parsing helpers for integers and booleans.  
  - `computePercentage()` for safe percentage calculation using `math/big` to avoid overflow.

- **uptime.go**  
  Core uptime calculation.  
  - `computeAllStationUptimes`:
    - Iterates over stations in sorted order.  
    - Calls `computeStationUptime` for each one.  
  - `computeStationUptime`:
    - Finds the earliest start and latest end timestamps across all chargers in the station.  
    - Collects all `up == true` intervals into one list.  
    - Sorts and merges intervals to avoid double-counting overlapping time.  
    - Computes the uptime percentage relative to the station’s total window.

(Additionally, a `go.mod` file is created via `go mod init charger_uptime` to treat this folder as a Go module.)

---

## 4. Edge Cases and Validation

Some specific cases I handle explicitly:

- **Station with chargers but no reports**  
  Considered invalid input and results in `ERROR`.

- **Overlapping “up” intervals**  
  When multiple chargers are up at the same time, or a single charger has overlapping records, merging ensures overlapping ranges are not double-counted.

- **Zero-length total window**  
  If all reports start and end at the same time, the total window length is zero. In that case, the station uptime is treated as `0`.

- **Malformed lines**  
  Wrong token counts, invalid numbers, or non-boolean values cause an immediate `ERROR`.

---

## 5. How to Run

From the project directory (where `main.go`, `parser.go`, `types.go`, `uptime.go`, and `go.mod` live):

### One-time setup (already done in this project)
```bash
go mod init charger_uptime
```

### Run with an input file
```bash
go run . input_1.txt
```

This compiles all Go files in the module and passes `input_1.txt` as the argument to the program.

Output format:

```text
<station_id> <uptime_percentage>
```

Stations are always printed in ascending station ID order.

---

## 6. Summary

The solution prioritizes:

- Clear structure (separate files for parsing, types, and uptime logic)  
- Strict adherence to the input specification and error behavior  
- Correct handling of gaps and overlaps in time intervals  
- Simplicity in data structures while remaining efficient and scalable  

This layout is intended to make the logic easy to follow and straightforward to extend or test.
