package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

// Report represents one availability record for a charger.
// Each record says: from [start,end] this charger was up or down.
type Report struct {
	start uint64
	end   uint64
	up    bool
}

// StationResult represents the final output we print per station.
type StationResult struct {
	StationID uint32
	UptimePct uint64
}

// fail prints "ERROR" and exits immediately.
// The challenge requires this behavior for *any* malformed input.
func fail() {
	fmt.Println("ERROR")
	os.Exit(1)
}

// parseUint32 parses a string as an unsigned 32-bit integer.
// On failure, follow challenge rules and print ERROR.
func parseUint32(s string) uint32 {
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		fail()
	}
	return uint32(v)
}

// parseUint64 parses a string as an unsigned 64-bit integer.
func parseUint64(s string) uint64 {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		fail()
	}
	return v
}

// parseBoolTF parses "true"/"false" (case-insensitive).
// Anything else is invalid.
func parseBoolTF(s string) bool {
	l := strings.ToLower(s)
	if l == "true" {
		return true
	}
	if l == "false" {
		return false
	}
	fail()
	return false
}

// computePercentage returns floor(100 * totalUp / totalWindow).
// We use math/big to avoid overflow when multiplying large nanos.
func computePercentage(totalUp, totalWindow uint64) uint64 {
	if totalWindow == 0 {
		return 0
	}
	up := new(big.Int).SetUint64(totalUp)
	win := new(big.Int).SetUint64(totalWindow)

	// up = (up * 100) / win
	up.Mul(up, big.NewInt(100))
	up.Div(up, win)

	return up.Uint64()
}
