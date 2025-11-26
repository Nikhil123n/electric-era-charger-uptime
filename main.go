package main

import (
	"fmt"
	"os"
)

func main() {
	// Expect exactly one argument: the path to the input file.
	// If not, print ERROR and exit (per challenge spec).
	if len(os.Args) != 2 {
		fail()
	}
	path := os.Args[1]

	// Parse the input into well-structured in-memory data.
	stationToChargers, chargerToStation, chargerReports :=
		parseInput(path)

	// Compute uptime for each station in sorted order.
	results :=
		computeAllStationUptimes(stationToChargers, chargerToStation, chargerReports)

	// Print each result on its own line: "<station_id> <uptime_pct>"
	for _, r := range results {
		fmt.Printf("%d %d\n", r.StationID, r.UptimePct)
	}
}
