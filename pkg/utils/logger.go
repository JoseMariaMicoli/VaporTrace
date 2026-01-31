package utils

import (
	"fmt"
	"strings"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/pterm/pterm"
)

// Global UI State
var UIMode = "CLI" // "CLI" or "TUI"

// Buffer increased to 1000 to prevent blocking during Mass-BOLA/BOPLA operations
var UI_Log_Chan = make(chan string, 1000)

// SetLoggerMode defines how outputs are rendered
func SetLoggerMode(mode string) {
	UIMode = mode
}

// EscapeTview sanitizes strings to prevent tview from interpreting brackets as color tags.
// This fixes the UI freezing issue when logging JSON or Arrays.
func EscapeTview(text string) string {
	return strings.ReplaceAll(text, "[", "[[")
}

// TacticalLog handles generic system messages
func TacticalLog(msg string) {
	if UIMode == "TUI" {
		// We do not escape here assuming the caller might want to use colors,
		// BUT if raw data is passed, the caller should sanitize it.
		// For safety in mass-ops, we sanitize key variable inputs in the caller,
		// or we can strictly separate "System Messages" (colored) from "Data" (escaped).
		select {
		case UI_Log_Chan <- msg:
		default:
			// Drop log if channel is full to prevent deadlocks, but with 1000 this is rare.
		}
	} else {
		pterm.Info.Println(msg)
	}
}

// RecordFinding is the Unified Pipeline (Phase 10.2.2)
// It handles Persistence (DB) AND Visualization (CLI/TUI) simultaneously.
func RecordFinding(f db.Finding) {
	// 1. Persistence Layer
	// We send a copy to avoid pointer issues if any
	db.LogQueue <- f

	// 2. Visualization Layer
	if UIMode == "TUI" {
		// Sanitize Content to prevent TUI corruption
		safeDetails := EscapeTview(f.Details)
		safeTarget := EscapeTview(f.Target)
		safeOWASP := EscapeTview(f.OWASP_ID)

		// Format specifically for tview dynamic colors
		// High-density format: [STATUS] (OWASP) DETAILS | TARGET
		var colorMsg string

		switch f.Status {
		case "CRITICAL", "VULNERABLE", "EXPLOITED":
			colorMsg = fmt.Sprintf("[red::b]%s[-] [yellow](%s)[-] [white]%s[-] [blue::b]| %s[-]",
				f.Status, safeOWASP, safeDetails, safeTarget)
		case "WEAK CONFIG", "POTENTIAL CALLBACK":
			colorMsg = fmt.Sprintf("[yellow]%s[-] [white]%s[-] [blue]| %s[-]",
				f.Status, safeDetails, safeTarget)
		case "SUCCESS", "INFO":
			colorMsg = fmt.Sprintf("[blue]%s[-] [white]%s[-] [blue]| %s[-]",
				f.Status, safeDetails, safeTarget)
		default:
			colorMsg = fmt.Sprintf("[white]%s[-] %s [blue]| %s[-]",
				f.Status, safeDetails, safeTarget)
		}

		select {
		case UI_Log_Chan <- colorMsg:
		default:
			// Non-blocking drop if UI is overwhelmed
		}
	} else {
		// CLI Pterm Output
		if f.Status == "VULNERABLE" || f.Status == "CRITICAL" || f.Status == "EXPLOITED" {
			pterm.Warning.Prefix = pterm.Prefix{Text: f.Status, Style: pterm.NewStyle(pterm.BgRed, pterm.FgWhite)}
			pterm.Warning.Printfln("%s (OWASP: %s) -> %s", f.Details, f.OWASP_ID, f.Target)
		} else if f.Status == "WEAK CONFIG" {
			pterm.Warning.Printfln("%s -> %s", f.Details, f.Target)
		} else {
			pterm.Success.Printfln("%s -> %s", f.Details, f.Target)
		}
	}
}
