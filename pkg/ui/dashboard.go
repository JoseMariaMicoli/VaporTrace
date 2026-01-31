package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/engine" // New Unified Engine
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app          *tview.Application
	pages        *tview.Pages
	header       *tview.TextView
	targetColumn *tview.Table
	brainLog     *tview.TextView
	mapView      *tview.TextView
	lootTable    *tview.Table
	reqView      *tview.TextView
	resView      *tview.TextView
	aiView       *tview.TextView
	statusFooter *tview.TextView
	cmdInput     *tview.InputField

	// Command History
	cmdHistory   []string
	historyIndex int

	knownCommands = []string{
		"auth", "sessions", "map", "swagger", "scrape", "mine", "proxy", "proxies", "target", "pipeline",
		"flow", "bola", "bopla", "bfla", "exhaust", "ssrf", "audit", "probe",
		"weaver", "loot", "test-bola", "test-bopla", "test-bfla", "test-exhaust", "test-ssrf", "test-audit", "test-probe",
		"init_db", "seed_db", "reset_db", "report", "clear", "exit", "usage", "help",
	}

	spinnerIdx    = 0
	spinnerFrames = []string{"▰▱▱▱▱", "▰▰▱▱▱", "▰▰▰▱▱", "▰▰▰▰▱", "▰▰▰▰▰"}
)

func InitTacticalDashboard() {
	utils.SetLoggerMode("TUI")

	app = tview.NewApplication()
	pages = tview.NewPages()

	header = tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	statusFooter = tview.NewTextView().SetDynamicColors(true)

	// --- INPUT FIELD SETUP ---
	cmdInput = tview.NewInputField().
		SetLabel("[aqua]VAPOR/INT> [white]").
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorYellow)
	cmdInput.SetBorder(true).SetBorderColor(tcell.ColorBlue)

	cmdInput.SetAutocompleteFunc(func(currentText string) (entries []string) {
		if len(currentText) == 0 {
			return nil
		}
		for _, cmd := range knownCommands {
			if strings.HasPrefix(strings.ToLower(cmd), strings.ToLower(currentText)) {
				entries = append(entries, cmd)
			}
		}
		return
	})

	// --- PIPELINE QUADRANT ---
	targetColumn = tview.NewTable().SetBorders(true).SetBordersColor(tcell.ColorBlue)
	targetColumn.SetTitle(" [blue]PIPELINE [white]").SetBorder(true)
	targetColumn.SetCell(0, 0, tview.NewTableCell("[black:blue] PROPERTY "))
	targetColumn.SetCell(0, 1, tview.NewTableCell("[black:blue] VALUE "))
	targetColumn.SetCell(1, 0, tview.NewTableCell("TARGET"))
	targetColumn.SetCell(1, 1, tview.NewTableCell("[red]NOT SET"))

	// --- PAGES SETUP ---
	brainLog = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).  // Added: Enables internal region tagging for complex logs
		SetWordWrap(true). // Prevents truncation of long attack payloads
		SetChangedFunc(func() {
			brainLog.ScrollToEnd() // Ensures the latest tactical data is always visible
			app.Draw()
		})
	brainLog.SetTitle(" [green]VAPOR_LOGS (TACTICAL FEED) [white]").SetBorder(true)
	brainLog.SetScrollable(true)

	mapView = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).  // Added: Allows interactive highlighting of endpoints
		SetWordWrap(true). // Prevents map breakage on small terminal windows
		SetTextAlign(tview.AlignCenter)
	mapView.SetTitle(" [blue]ATTACK_SURFACE [white]").SetBorder(true)

	lootTable = tview.NewTable().
		SetBorders(true).
		SetBordersColor(tcell.ColorDarkCyan).
		SetSelectable(true, false) // Added: Allows navigating through captured secrets
	lootTable.SetTitle(" [magenta]LOOT_VAULT [white]").SetBorder(true)

	reqView = tview.NewTextView().SetDynamicColors(true).SetWordWrap(true)
	reqView.SetTitle(" [yellow]TRAFFIC_REQ [white]").SetBorder(true)

	resView = tview.NewTextView().SetDynamicColors(true).SetWordWrap(true)
	resView.SetTitle(" [green]TRAFFIC_RES [white]").SetBorder(true)

	trafficSplit := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(reqView, 0, 1, false).
		AddItem(resView, 0, 1, false)

	aiView = tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true) // Prevents AI analysis text from cutting off
	aiView.SetTitle(" [white:blue] LOGIC_ANALYZER [white] ").SetBorder(true)

	pages.AddPage("logs", brainLog, true, true)
	pages.AddPage("map", mapView, true, false)
	pages.AddPage("loot", lootTable, true, false)
	pages.AddPage("traffic", trafficSplit, true, false)
	pages.AddPage("ai", aiView, true, false)

	// --- LAYOUT ---
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 10, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(targetColumn, 35, 1, false).
			AddItem(pages, 0, 4, false),
			0, 4, false).
		AddItem(statusFooter, 1, 1, false).
		AddItem(cmdInput, 3, 1, true)

	updateTabs("logs")

	// --- GLOBAL KEYBINDS ---
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			switchTo("logs")
		case tcell.KeyF2:
			switchTo("map")
		case tcell.KeyF3:
			switchTo("loot")
		case tcell.KeyF4:
			switchTo("traffic")
		case tcell.KeyF5:
			switchTo("ai")
		case tcell.KeyEsc:
			confirmExit()
		}
		return event
	})

	// --- INPUT HISTORY & EXECUTION ---
	cmdInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			if len(cmdHistory) > 0 && historyIndex > 0 {
				historyIndex--
				cmdInput.SetText(cmdHistory[historyIndex])
			} else if len(cmdHistory) > 0 && historyIndex == len(cmdHistory) {
				historyIndex = len(cmdHistory) - 1
				cmdInput.SetText(cmdHistory[historyIndex])
			}
			return nil
		case tcell.KeyDown:
			if len(cmdHistory) > 0 && historyIndex < len(cmdHistory)-1 {
				historyIndex++
				cmdInput.SetText(cmdHistory[historyIndex])
			} else {
				historyIndex = len(cmdHistory)
				cmdInput.SetText("")
			}
			return nil
		}
		return event
	})

	cmdInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			text := cmdInput.GetText()
			cmdInput.SetText("") // Clear immediately so no feedback appears in input

			if text == "" {
				return
			}

			// Special handling for 'exit' to trigger Modal
			if strings.TrimSpace(text) == "exit" {
				confirmExit()
				return
			}

			// Add to history
			cmdHistory = append(cmdHistory, text)
			historyIndex = len(cmdHistory)

			// Switch to Logs Tab so user sees feedback
			switchTo("logs")

			// Execute
			engine.ExecuteCommand(text)
		}
	})

	startAsyncEngines()
	if err := app.SetRoot(mainFlex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

// confirmExit displays the TUI modal and ensures focus
func confirmExit() {
	modal := tview.NewModal().
		SetText("Secure Shutdown Protocol?\n(Terminates all listeners)").
		AddButtons([]string{"Yes", "No"}).
		SetBackgroundColor(tcell.ColorDarkRed).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				pages.RemovePage("modal")
				engine.ExecuteCommand("__internal_shutdown")
			} else {
				pages.RemovePage("modal")
				// Return focus to input
				app.SetFocus(cmdInput)
			}
		})

	pages.AddPage("modal", modal, false, true)
	// Explicitly set focus to modal to capture keyboard events immediately
	app.SetFocus(modal)
}

func switchTo(page string) {
	updateTabs(page)
	pages.SwitchToPage(page)
}

func updateTabs(active string) {
	tabs := []string{"LOGS", "MAP", "LOOT", "TRAFFIC", "ANALYSIS"}
	ids := []string{"logs", "map", "loot", "traffic", "ai"}
	var formatted []string
	for i, t := range tabs {
		if active == ids[i] {
			formatted = append(formatted, "[black:aqua]┢ "+t+" ┪[-]")
		} else {
			formatted = append(formatted, "[black:blue]┠ "+t+" ┨[-]")
		}
	}
	header.SetText(fmt.Sprintf(`[aqua:black:b]
██╗   ██╗ █████╗ ██████╗  ██████╗ ██████╗ ████████╗██████╗  █████╗  ██████╗███████╗
██║   ██║██╔══██╗██╔══██╗██╔═══██╗██╔══██╗╚══██╔══╝██╔══██╗██╔══██╗██╔════╝██╔════╝
╚██╗ ██╔╝███████║██████╔╝██║   ██║██████╔╝   ██║   ██████╔╝███████║██║     █████╗  
 ╚████╔╝ ██╔══██║██╔═══╝ ██║   ██║██╔══██╗   ██║   ██╔══██╗██╔══██║██║     ██╔══╝  
  ╚██╔╝  ██║  ██║██║     ╚██████╔╝██║  ██║   ██║   ██║  ██║██║  ██║╚██████╗███████╗
   ╚═╝   ╚═╝  ╚═╝╚═╝      ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚══════╝[-]
%s`, strings.Join(formatted, "  ")))
}

func startAsyncEngines() {
	// 1. Spinner Logic
	go func() {
		for {
			time.Sleep(250 * time.Millisecond)
			app.QueueUpdateDraw(func() {
				spinnerIdx = (spinnerIdx + 1) % len(spinnerFrames)
				statusFooter.SetText(fmt.Sprintf(" [blue]SYNC %s [white]| %s", spinnerFrames[spinnerIdx], time.Now().Format("15:04:05")))
			})
		}
	}()

	// 2. Log & Target Update Logic
	go func() {
		for msg := range utils.UI_Log_Chan {
			app.QueueUpdateDraw(func() {
				if strings.Contains(msg, "Target Locked:") {
					parts := strings.Split(msg, "Target Locked:[-] ")
					if len(parts) > 1 {
						url := strings.TrimSpace(parts[1])
						targetColumn.SetCell(1, 1, tview.NewTableCell("[green]"+url))
					}
				}
				fmt.Fprintf(brainLog, "[%s] %s\n", time.Now().Format("15:04:05"), msg)
				brainLog.ScrollToEnd()
			})
		}
	}()
}
