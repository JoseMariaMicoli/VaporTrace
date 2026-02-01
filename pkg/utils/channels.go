package utils

import (
	"time"
)

// LootPacket defines the structure for F3 (Loot) Table Data
type LootPacket struct {
	Type      string
	Value     string
	Source    string
	Timestamp time.Time
}

// MapDataChan handles Recon/Discovery updates (F2 Tab)
// Buffer size 1000 to handle high-concurrency scraping without blocking
var MapDataChan = make(chan string, 1000)

// LootDataChan handles Secret/PII updates (F3 Tab)
var LootDataChan = make(chan LootPacket, 500)

// LogMap sends discovery intelligence to the F2 Tab
func LogMap(endpoint, method, status string) {
	// Only log if in TUI mode to avoid blocking CLI
	if UIMode == "TUI" {
		msg := "[green]DISCOVERED[-] " + endpoint
		if status != "" {
			msg += " [yellow](" + status + ")[-]"
		}
		select {
		case MapDataChan <- msg:
		default:
			// Drop if buffer full to prevent pipeline stall
		}
	}
}

// LogLoot sends captured secrets to the F3 Tab
func LogLoot(lType, value, source string) {
	if UIMode == "TUI" {
		pkt := LootPacket{
			Type:      lType,
			Value:     value,
			Source:    source,
			Timestamp: time.Now(),
		}
		select {
		case LootDataChan <- pkt:
		default:
			// Drop if buffer full
		}
	}
}
