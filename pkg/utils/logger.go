package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/pterm/pterm"
)

// Global UI State
var UIMode = "CLI" // "CLI" or "TUI"
// Increased buffer size to prevent blocking during mass scans
var UI_Log_Chan = make(chan string, 5000)

// SetLoggerMode defines how outputs are rendered
func SetLoggerMode(mode string) {
	UIMode = mode
}

// EscapeTview sanitizes strings to prevent tview from interpreting brackets as tags
func EscapeTview(text string) string {
	text = StripANSI(text)
	return strings.ReplaceAll(text, "[", "[[")
}

func timeStamp() string {
	return time.Now().Format("15:04:05")
}

// TacticalLog handles generic system messages with enforced formatting
func TacticalLog(msg string) {
	// 1. TUI MODE: Send to channel (Thread-Safe)
	if UIMode == "TUI" {
		cleanMsg := StripANSI(msg)

		if msg == "___CLEAR_SCREEN_SIGNAL___" {
			select {
			case UI_Log_Chan <- msg:
			default:
			}
			return
		}

		// Color coding for generic logs
		colorTag := "[white]"
		if strings.Contains(strings.ToLower(cleanMsg), "success") || strings.Contains(msg, "[green]") {
			colorTag = "[green]"
		} else if strings.Contains(strings.ToLower(cleanMsg), "error") || strings.Contains(strings.ToLower(cleanMsg), "fail") || strings.Contains(msg, "[red]") {
			colorTag = "[red]"
		} else if strings.Contains(strings.ToLower(cleanMsg), "warn") || strings.Contains(msg, "[yellow]") {
			colorTag = "[yellow]"
		} else if strings.Contains(strings.ToLower(cleanMsg), "phase") || strings.Contains(msg, "[cyan]") {
			colorTag = "[cyan]"
		}

		// Format: [TIME] MESSAGE
		formatted := fmt.Sprintf("[gray][%s][-] %s%s[-]", timeStamp(), colorTag, cleanMsg)

		select {
		case UI_Log_Chan <- formatted:
		default:
			// Drop message if buffer full to prevent UI freeze
		}
	} else {
		// 2. CLI MODE: Print directly via Pterm
		pterm.Info.Println(msg)
	}
}

// RecordFinding is the Unified Pipeline for vulnerability reports
func RecordFinding(f db.Finding) {
	// 1. Persistence Layer
	db.LogQueue <- f

	// 2. Visualization Layer
	safeDetails := EscapeTview(f.Details)
	safeTarget := EscapeTview(f.Target)
	safeOWASP := EscapeTview(f.OWASP_ID)
	ts := timeStamp()

	if UIMode == "TUI" {
		var logLine string

		// Strict Format: [TIME] [STATUS] (OWASP) DETAILS | TARGET
		switch f.Status {
		case "CRITICAL", "EXPLOITED":
			logLine = fmt.Sprintf("[gray][%s][-] [red::b][%s][-] [yellow](%s)[-] [white]%s[-] [blue]| %s[-]",
				ts, f.Status, safeOWASP, safeDetails, safeTarget)
		case "VULNERABLE":
			logLine = fmt.Sprintf("[gray][%s][-] [red][%s][-] [yellow](%s)[-] [white]%s[-] [blue]| %s[-]",
				ts, f.Status, safeOWASP, safeDetails, safeTarget)
		case "WEAK CONFIG", "POTENTIAL CALLBACK":
			logLine = fmt.Sprintf("[gray][%s][-] [yellow][%s][-] [white]%s[-] [blue]| %s[-]",
				ts, f.Status, safeDetails, safeTarget)
		case "SUCCESS", "ACTIVE":
			logLine = fmt.Sprintf("[gray][%s][-] [green][%s][-] [white]%s[-] [blue]| %s[-]",
				ts, f.Status, safeDetails, safeTarget)
		case "INFO":
			logLine = fmt.Sprintf("[gray][%s][-] [blue][%s][-] [white]%s[-] [blue]| %s[-]",
				ts, "INFO", safeDetails, safeTarget)
		default:
			logLine = fmt.Sprintf("[gray][%s][-] [white][%s][-] %s [blue]| %s[-]",
				ts, f.Status, safeDetails, safeTarget)
		}

		select {
		case UI_Log_Chan <- logLine:
		default:
		}
	} else {
		// CLI Pterm Output
		if f.Status == "VULNERABLE" || f.Status == "CRITICAL" {
			pterm.Warning.Printfln("[%s] %s -> %s", f.Status, f.Details, f.Target)
		} else {
			pterm.Info.Printfln("[%s] %s -> %s", f.Status, f.Details, f.Target)
		}
	}
}
