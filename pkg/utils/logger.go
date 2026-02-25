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

var UIMode = "CLI" // CLI | TUI

var UI_Log_Chan = make(chan string, 1000)
var ContextLogChan = make(chan string, 500)
var NeuroLogChan = make(chan string, 500)

type TrafficPacket struct {
	ReqHeader string
	ReqBody   string
	ResHeader string
	ResBody   string
}

var TrafficChan = make(chan TrafficPacket, 500)

type LogBufferStruct struct {
	mu    sync.Mutex
	lines []string
	cap   int
}

var GlobalLogBuffer = &LogBufferStruct{
	lines: make([]string, 0),
	cap:   1000,
}

func SetLoggerMode(mode string) {
	UIMode = mode
}

func EscapeTview(text string) string {
	text = StripANSI(text)
	text = strings.ReplaceAll(text, "[", "[[")
	text = strings.ReplaceAll(text, "]", "]]")
	return text
}

func timeStamp() string {
	return time.Now().Format("15:04:05")
}

func TacticalLog(msg string) {
	if UIMode == "TUI" {
		if msg == "___CLEAR_SCREEN_SIGNAL___" {
			select {
			case UI_Log_Chan <- msg:
			default:
			}
			return
		}

		cleanMsg := StripANSI(msg)
		colorTag := "[white]"
		lower := strings.ToLower(cleanMsg)
		switch {
		case strings.Contains(lower, "error") || strings.Contains(msg, "[red]"):
			colorTag = "[red]"
		case strings.Contains(lower, "warn") || strings.Contains(msg, "[yellow]"):
			colorTag = "[yellow]"
		case strings.Contains(lower, "success") || strings.Contains(msg, "[green]"):
			colorTag = "[green]"
		case strings.Contains(lower, "phase") || strings.Contains(msg, "[cyan]"):
			colorTag = "[cyan]"
		}

		formatted := fmt.Sprintf("[gray][%s][-] %s%s[-]", timeStamp(), colorTag, cleanMsg)

		GlobalLogBuffer.mu.Lock()
		GlobalLogBuffer.lines = append(GlobalLogBuffer.lines, formatted)
		if len(GlobalLogBuffer.lines) > GlobalLogBuffer.cap {
			GlobalLogBuffer.lines = GlobalLogBuffer.lines[1:]
		}
		GlobalLogBuffer.mu.Unlock()

		select {
		case UI_Log_Chan <- formatted:
		default:
		}
		return
	}

	pterm.Info.Println(msg)
}

func LogContext(msg string) {
	if UIMode == "TUI" {
		select {
		case ContextLogChan <- msg:
		default:
		}
		return
	}
	TacticalLog(msg)
}

func LogNeural(msg string) {
	if UIMode == "TUI" {
		select {
		case NeuroLogChan <- msg:
		default:
		}
		return
	}
	TacticalLog(msg)
}

func LogTraffic(reqHeader, reqBody, resHeader, resBody string) {
	if UIMode != "TUI" {
		return
	}
	pkt := TrafficPacket{
		ReqHeader: reqHeader,
		ReqBody:   reqBody,
		ResHeader: resHeader,
		ResBody:   resBody,
	}
	select {
	case TrafficChan <- pkt:
	default:
	}
}

func RecordFinding(f db.Finding) {
	enrichment.EnrichFinding(&f)
	db.LogQueue <- f

	if UIMode == "TUI" {
		safeDetails := EscapeTview(f.Details)
		safeTarget := EscapeTview(f.Target)
		safeOWASP := EscapeTview(f.OWASP_ID)

		var line string
		switch f.Status {
		case "CRITICAL", "VULNERABLE", "EXPLOITED":
			line = fmt.Sprintf("[red::b]%s[-] [yellow](%s)[-] [white]%s[-] [blue::b]| %s[-]", f.Status, safeOWASP, safeDetails, safeTarget)
		case "WEAK CONFIG", "POTENTIAL CALLBACK":
			line = fmt.Sprintf("[yellow]%s[-] [white]%s[-] [blue]| %s[-]", f.Status, safeDetails, safeTarget)
		case "SUCCESS", "INFO":
			line = fmt.Sprintf("[blue]%s[-] [white]%s[-] [blue]| %s[-]", f.Status, safeDetails, safeTarget)
		default:
			line = fmt.Sprintf("[white]%s[-] %s [blue]| %s[-]", f.Status, safeDetails, safeTarget)
		}

		select {
		case UI_Log_Chan <- line:
		default:
		}
		return
	}

	if f.Status == "VULNERABLE" || f.Status == "CRITICAL" || f.Status == "EXPLOITED" {
		pterm.Warning.Prefix = pterm.Prefix{Text: f.Status, Style: pterm.NewStyle(pterm.BgRed, pterm.FgWhite)}
		pterm.Warning.Printfln("%s (OWASP: %s) -> %s", f.Details, f.OWASP_ID, f.Target)
	} else if f.Status == "WEAK CONFIG" {
		pterm.Warning.Printfln("%s -> %s", f.Details, f.Target)
	} else {
		pterm.Success.Printfln("%s -> %s", f.Details, f.Target)
	}
}
