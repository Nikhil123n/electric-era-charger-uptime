package main

import (
	"sort"
)

// computeAllStationUptimes walks each station, computes its uptime,
// and returns a sorted list of results by stationID.
func computeAllStationUptimes(
	stationToChargers map[uint32][]uint32,
	chargerToStation map[uint32]uint32, // unused but kept for future extensibility
	chargerReports map[uint32][]Report,
) []StationResult {

	// We must print results in ascending station ID order.
	stationIDs := make([]uint32, 0, len(stationToChargers))
	for sid := range stationToChargers {
		stationIDs = append(stationIDs, sid)
	}
	sort.Slice(stationIDs, func(i, j int) bool { return stationIDs[i] < stationIDs[j] })

	results := make([]StationResult, 0, len(stationIDs))

	for _, sid := range stationIDs {
		uptime := computeStationUptime(sid, stationToChargers[sid], chargerReports)
		results = append(results, StationResult{
			StationID: sid,
			UptimePct: uptime,
		})
	}

	return results
}

// computeStationUptime calculates the uptime % for a single station.
//
// Steps:
// 1. Find the minimum start time and maximum end time across *all* chargers.
//    This defines the total "station window".
// 2. Collect all intervals where a charger was UP into one list.
// 3. Sort and merge overlapping UP intervals.
// 4. Calculate totalUp = sum of merged intervals.
// 5. uptime = floor(100 * totalUp / totalWindow).
func computeStationUptime(
	stationID uint32,
	chargers []uint32,
	chargerReports map[uint32][]Report,
) uint64 {
	var (
		minStart    uint64
		maxEnd      uint64
		minStartSet bool
		maxEndSet   bool
	)

	upIntervals := make([]Report, 0)
	hasAnyReports := false

	// Walk all chargers belonging to this station.
	for _, cid := range chargers {
		reps := chargerReports[cid]

		for _, r := range reps {
			hasAnyReports = true

			if !minStartSet || r.start < minStart {
				minStart = r.start
				minStartSet = true
			}
			if !maxEndSet || r.end > maxEnd {
				maxEnd = r.end
				maxEndSet = true
			}

			if r.up {
				upIntervals = append(upIntervals, r)
			}
		}
	}

	// A station with chargers but no reports is invalid.
	if !hasAnyReports {
		fail()
	}

	totalWindow := maxEnd - minStart
	if totalWindow == 0 {
		// All reports happened at the same nanosecond.
		return 0
	}

	if len(upIntervals) == 0 {
		// No UP time at all.
		return 0
	}

	// Sort by start time, then by end time.
	sort.Slice(upIntervals, func(i, j int) bool {
		if upIntervals[i].start == upIntervals[j].start {
			return upIntervals[i].end < upIntervals[j].end
		}
		return upIntervals[i].start < upIntervals[j].start
	})

	// Merge overlapping intervals.
	currentStart := upIntervals[0].start
	currentEnd := upIntervals[0].end
	var totalUp uint64

	for _, r := range upIntervals[1:] {
		// Overlap or touch?
		if r.start <= currentEnd {
			if r.end > currentEnd {
				currentEnd = r.end
			}
		} else {
			// Disjoint interval -> close previous one.
			totalUp += currentEnd - currentStart
			currentStart = r.start
			currentEnd = r.end
		}
	}
	// Add the final interval.
	totalUp += currentEnd - currentStart

	uptimePct := computePercentage(totalUp, totalWindow)
	if uptimePct > 100 {
		uptimePct = 100
	}

	return uptimePct
}
