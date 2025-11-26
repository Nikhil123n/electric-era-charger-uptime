package main

import "testing"

func TestSimpleUptime(t *testing.T) {
	reports := map[uint32][]Report{
		1000: {
			{start: 0, end: 100, up: true},
			{start: 100, end: 200, up: false},
		},
	}

	chargers := []uint32{1000}

	uptime := computeStationUptime(0, chargers, reports)
	if uptime != 50 {
		t.Fatalf("expected 50, got %d", uptime)
	}
}
