/*
Copyright (c) 2026 José María Micoli
Licensed under {'license_type': 'BSL', 'change_date': '2033-02-17', 'convert_to': 'Apache-2.0'}

You may:
✔ Study
✔ Modify
✔ Use for internal security testing

You may NOT:
✘ Offer as a commercial service
✘ Sell derived competing products
*/

package utils

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/enrichment"
	"github.com/pterm/pterm"
)

// Global UI State
var UIMode = "CLI" // "CLI" or "TUI"

// Buffer increased to 1000 to prevent blocking during Mass-BOLA/BOPLA operations
var UI_Log_Chan = make(chan string, 1000)

// UI_Log_Buffer maintains the scrollback limit before sending to UI
// Task 2 Part 2: Enforce Hard Cap of 1,000 lines.
type LogBufferStruct struct {
	mu    sync.Mutex
	lines []string
	cap   int
}

var GlobalLogBuffer = &LogBufferStruct{
	lines: make([]string, 0),
	cap:   1000,
}

// SetLoggerMode defines how outputs are rendered
func SetLoggerMode(mode string) {
	UIMode = mode
}

// EscapeTview sanitizes strings to prevent tview from interpreting brackets as color tags.
// This fixes the UI freezing issue when logging JSON or Arrays.
func EscapeTview(text string) string {
	text = StripANSI(text)
	text = strings.ReplaceAll(text, "[", "[[")
	text = strings.ReplaceAll(text, "]", "]]")
	return text
}

// TacticalLog handles generic system messages
// Modified to use GlobalLogBuffer for scrollback limit enforcement
func TacticalLog(msg string) {
	if UIMode == "TUI" {
		cleanMsg := StripANSI(msg)

		if msg == "___CLEAR_SCREEN_SIGNAL___" {
			select {
			case UI_Log_Chan <- msg:
			default:
			}
			return
		}

		colorTag := "[white]"
		if strings.Contains(strings.ToLower(cleanMsg), "success") || strings.Contains(msg, "[green]") {
			colorTag = "[green]"
		} else if strings.Contains(strings.ToLower(cleanMsg), "error") || strings.Contains(msg, "[red]") {
			colorTag = "[red]"
		} else if strings.Contains(strings.ToLower(cleanMsg), "warn") || strings.Contains(msg, "[yellow]") {
			colorTag = "[yellow]"
		} else if strings.Contains(strings.ToLower(cleanMsg), "phase") || strings.Contains(msg, "[cyan]") {
			colorTag = "[cyan]"
		}

		formatted := fmt.Sprintf("[gray][%s][-] %s%s[-]", timeStamp(), colorTag, cleanMsg)

		// Enforce Scrollback Limit
		GlobalLogBuffer.mu.Lock()
		GlobalLogBuffer.lines = append(GlobalLogBuffer.lines, formatted)
		if len(GlobalLogBuffer.lines) > GlobalLogBuffer.cap {
			// Drop oldest
			GlobalLogBuffer.lines = GlobalLogBuffer.lines[1:]
		}
		GlobalLogBuffer.mu.Unlock()

		select {
		case UI_Log_Chan <- formatted:
		default:
			// Drop log if channel is full to prevent deadlocks, but with 1000 this is rare.
		}
	} else {
		pterm.Info.Println(msg)
	}
}

// RecordFinding persists findings and logs them
func RecordFinding(f db.Finding) {
	enrichment.EnrichFinding(&f)
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

	if UIMode == "TUI" {
		var logLine string
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
		case UI_Log_Chan <- logLine:
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
