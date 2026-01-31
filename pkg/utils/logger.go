package utils

import (
	"fmt"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/pterm/pterm"
)

// Global UI State
var UIMode = "CLI" // "CLI" or "TUI"

var UI_Log_Chan = make(chan string, 100)

// SetLoggerMode defines how outputs are rendered
func SetLoggerMode(mode string) {
	UIMode = mode
}

// TacticalLog handles generic system messages
func TacticalLog(msg string) {
	if UIMode == "TUI" {
		select {
		case UI_Log_Chan <- msg:
		default:
			// Drop log if channel is full to prevent deadlocks
		}
	} else {
		// Strip ANSI codes if needed for raw log, but pterm handles them well
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
		// Format specifically for tview dynamic colors
		var colorMsg string

		switch f.Status {
		case "CRITICAL", "VULNERABLE", "EXPLOITED":
			colorMsg = fmt.Sprintf("[red::b]%s[-:-:-] [yellow](%s)[-] [white]%s[-]", f.Status, f.OWASP_ID, f.Details)
		case "WEAK CONFIG", "POTENTIAL CALLBACK":
			colorMsg = fmt.Sprintf("[yellow]%s[-] [white]%s[-]", f.Status, f.Details)
		case "SUCCESS", "INFO":
			colorMsg = fmt.Sprintf("[blue]%s[-] [white]%s[-]", f.Status, f.Details)
		default:
			colorMsg = fmt.Sprintf("[white]%s[-] %s", f.Status, f.Details)
		}

		select {
		case UI_Log_Chan <- colorMsg:
		default:
		}
	} else {
		// CLI Pterm Output
		if f.Status == "VULNERABLE" || f.Status == "CRITICAL" || f.Status == "EXPLOITED" {
			pterm.Warning.Prefix = pterm.Prefix{Text: f.Status, Style: pterm.NewStyle(pterm.BgRed, pterm.FgWhite)}
			pterm.Warning.Printfln("%s (OWASP: %s)", f.Details, f.OWASP_ID)
		} else if f.Status == "WEAK CONFIG" {
			pterm.Warning.Printfln("%s", f.Details)
		} else {
			pterm.Success.Printfln("%s", f.Details)
		}
	}
}
