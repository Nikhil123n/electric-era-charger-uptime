package main

import (
	"bufio"
	"os"
	"strings"
)

type parseState int

// These represent which section of the file we are currently reading.
const (
	stateInit parseState = iota
	stateStations
	stateReports
)

// parseInput reads the file and extracts:
//
// 1. stationToChargers: stationID → list of chargerIDs
// 2. chargerToStation: chargerID → stationID (used to detect duplicates)
// 3. chargerReports:  chargerID → its availability records
//
// Any structural error → print ERROR and exit.
func parseInput(path string) (
	map[uint32][]uint32,
	map[uint32]uint32,
	map[uint32][]Report,
) {
	f, err := os.Open(path)
	if err != nil {
		fail()
	}
	defer f.Close()

	stationToChargers := make(map[uint32][]uint32)
	chargerToStation := make(map[uint32]uint32)
	chargerReports := make(map[uint32][]Report)

	state := stateInit
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Section headers
		switch line {
		case "[Stations]":
			if state != stateInit {
				fail()
			}
			state = stateStations
			continue
		case "[Charger Availability Reports]":
			if state != stateStations {
				fail()
			}
			state = stateReports
			continue
		}

		// Section-specific parsing
		switch state {
		case stateStations:
			parseStationsLine(line, stationToChargers, chargerToStation)
		case stateReports:
			parseReportLine(line, chargerToStation, chargerReports)
		default:
			// We encountered data before "[Stations]" or other invalid transitions.
			fail()
		}
	}

	if err := scanner.Err(); err != nil {
		fail()
	}

	// Must have at least one station defined.
	if len(stationToChargers) == 0 {
		fail()
	}

	return stationToChargers, chargerToStation, chargerReports
}

// parseStationsLine processes one line under [Stations], e.g.:
//    0 1000 1001 1002
// stationID must be unique, and each chargerID can belong to exactly one station.
func parseStationsLine(
	line string,
	stationToChargers map[uint32][]uint32,
	chargerToStation map[uint32]uint32,
) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		fail() // station with no chargers is invalid
	}

	sid := parseUint32(parts[0])
	if _, exists := stationToChargers[sid]; exists {
		// duplicate station definition
		fail()
	}

	chargers := make([]uint32, 0, len(parts)-1)
	for _, p := range parts[1:] {
		cid := parseUint32(p)

		// Ensure charger isn't already assigned to another station.
		if _, exists := chargerToStation[cid]; exists {
			fail()
		}

		chargers = append(chargers, cid)
		chargerToStation[cid] = sid
	}

	stationToChargers[sid] = chargers
}

// parseReportLine processes one availability record.
// Format: <chargerID> <startNanos> <endNanos> <true|false>
func parseReportLine(
	line string,
	chargerToStation map[uint32]uint32,
	chargerReports map[uint32][]Report,
) {
	parts := strings.Fields(line)
	if len(parts) != 4 {
		fail()
	}

	cid := parseUint32(parts[0])
	start := parseUint64(parts[1])
	end := parseUint64(parts[2])
	if end < start {
		fail()
	}
	up := parseBoolTF(parts[3])

	// Reject reports for chargers never declared under [Stations].
	if _, ok := chargerToStation[cid]; !ok {
		fail()
	}

	chargerReports[cid] = append(chargerReports[cid], Report{
		start: start,
		end:   end,
		up:    up,
	})
}
