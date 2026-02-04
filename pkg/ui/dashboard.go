package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/engine"
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
	neuroView    *tview.TextView
	reportView   *tview.TextView // F7 Report Tab
	statusFooter *tview.TextView
	cmdInput     *tview.InputField

	cmdHistory   []string
	historyIndex int
	historyFile  = ".vapor_history"

	knownCommands = []string{
		"auth", "sessions", "map", "swagger", "scrape", "mine", "proxy", "proxies", "target", "pipeline",
		"flow", "bola", "bopla", "bfla", "exhaust", "ssrf", "audit", "probe",
		"weaver", "loot", "test-bola", "test-bopla", "test-bfla", "test-exhaust", "test-ssrf", "test-audit", "test-probe",
		"neuro", "test-neuro", "neuro-gen",
		"init_db", "seed_db", "reset_db", "report", "clear", "exit", "usage", "help",
	}

	spinnerIdx    = 0
	spinnerFrames = []string{"▰▱▱▱▱", "▰▰▱▱▱", "▰▰▰▱▱", "▰▰▰▰▱", "▰▰▰▰▰"}
)

// LoadHistory reads the command history from disk
func LoadHistory() {
	file, err := os.Open(historyFile)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			cmdHistory = append(cmdHistory, line)
		}
	}
	historyIndex = len(cmdHistory)
}

// SaveHistory appends the last command to the disk
func SaveHistory(cmd string) {
	f, err := os.OpenFile(historyFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(cmd + "\n")
}

func InitTacticalDashboard() {
	utils.SetLoggerMode("TUI")
	LoadHistory() // Initialize history from persistence

	app = tview.NewApplication()
	pages = tview.NewPages()

	header = tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	statusFooter = tview.NewTextView().SetDynamicColors(true)

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

	// --- PIPELINE QUADRANT SETUP ---
	targetColumn = tview.NewTable().SetBorders(true).SetBordersColor(tcell.ColorBlue)
	targetColumn.SetTitle(" [blue]PIPELINE [white]").SetBorder(true)
	// Initial Headers
	targetColumn.SetCell(0, 0, tview.NewTableCell("[black:blue] PROPERTY "))
	targetColumn.SetCell(0, 1, tview.NewTableCell("[black:blue] VALUE "))
	// Pre-populate with defaults
	updatePipelineQuadrant()

	// --- VIEW SETUP ---
	brainLog = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true)

	brainLog.SetTitle(" [green]VAPOR_LOGS (TACTICAL FEED) [white]").SetBorder(true)

	mapView = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetTextAlign(tview.AlignCenter)
	mapView.SetTitle(" [blue]ATTACK_SURFACE [white]").SetBorder(true)

	lootTable = tview.NewTable().
		SetBorders(true).
		SetBordersColor(tcell.ColorDarkCyan).
		SetSelectable(true, false)
	lootTable.SetTitle(" [magenta]LOOT_VAULT [white]").SetBorder(true)
	lootTable.SetCell(0, 0, tview.NewTableCell("[black:cyan] TYPE "))
	lootTable.SetCell(0, 1, tview.NewTableCell("[black:cyan] VALUE "))
	lootTable.SetCell(0, 2, tview.NewTableCell("[black:cyan] SOURCE "))
	lootTable.SetFixed(1, 0) // Fix header row

	reqView = tview.NewTextView().SetDynamicColors(true).SetWordWrap(true).SetScrollable(true)
	reqView.SetTitle(" [yellow]REQUEST (UPPER) - Ctrl+A to Analyze [white]").SetBorder(true)
	resView = tview.NewTextView().SetDynamicColors(true).SetWordWrap(true).SetScrollable(true)
	resView.SetTitle(" [green]RESPONSE (LOWER) [white]").SetBorder(true)

	trafficSplit := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(reqView, 0, 1, false).
		AddItem(resView, 0, 1, false)

	aiView = tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true)
	aiView.SetTitle(" [white:blue] LOGIC_ANALYZER [white] ").SetBorder(true)

	aiView = tview.NewTextView().SetDynamicColors(true).SetWordWrap(true).SetScrollable(true)
	aiView.SetTitle(" [white:blue] CONTEXT_AGGREGATOR (F5) [white] ").SetBorder(true)
	aiView.SetText("[gray]Initializing Context Aggregator...\n\n[blue]●[-] Intelligence Harvest: [green]ACTIVE[-]\n[blue]●[-] Watching For: [white]JWTs, AWS Keys, Bearer Tokens[-]\n[blue]●[-] Correlation Engine: [white]Cross-referencing Findings[-]\n\n")

	neuroView = tview.NewTextView().SetDynamicColors(true).SetWordWrap(true).SetScrollable(true)
	neuroView.SetTitle(" [magenta:black] NEURAL ENGINE (F6) [white] ").SetBorder(true)

	// F7 Report Tab Initialization
	reportView = InitReportTab()

	// Add Pages
	pages.AddPage("logs", brainLog, true, true)
	pages.AddPage("map", mapView, true, false)
	pages.AddPage("loot", lootTable, true, false)
	pages.AddPage("traffic", trafficSplit, true, false)
	pages.AddPage("ai", aiView, true, false)
	pages.AddPage("neuro", neuroView, true, false)
	pages.AddPage("report", reportView, true, false)

	// Main Layout
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 10, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(targetColumn, 35, 1, false).
			AddItem(pages, 0, 4, false),
			0, 4, false).
		AddItem(statusFooter, 1, 1, false).
		AddItem(cmdInput, 3, 1, true)

	updateTabs("logs")

	// --- GLOBAL KEYBINDS & MOUSE SCROLL ---
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// F1 - F7 Tab Switching
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
		case tcell.KeyF6:
			switchTo("neuro")
		case tcell.KeyF7:
			// Auto-Load Findings on switch
			LoadFindings()
			switchTo("report")

		// Tactical Shortcuts
		case tcell.KeyCtrlI:
			logic.InterceptorActive = !logic.InterceptorActive
			state := "OFF"
			color := "[red]"
			if logic.InterceptorActive {
				state = "ON"
				color = "[green]"
			}
			utils.TacticalLog(fmt.Sprintf("%sINTERCEPTOR TOGGLED: %s (Wait for Packets)[-]", color, state))

		case tcell.KeyCtrlH:
			// Global Help Modal
			ShowHelpModal(app, pages)

		case tcell.KeyPgUp:
			row, col := brainLog.GetScrollOffset()
			if row > 0 {
				brainLog.ScrollTo(row-1, col)
			}
		case tcell.KeyPgDn:
			row, col := brainLog.GetScrollOffset()
			brainLog.ScrollTo(row+1, col)
		case tcell.KeyEsc:
			confirmExit()
		}
		return event
	})

	// Input History Navigation
	cmdInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			if len(cmdHistory) > 0 && historyIndex > 0 {
				historyIndex--
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

	// Command Execution
	cmdInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			text := cmdInput.GetText()
			cmdInput.SetText("")

			if text == "" {
				return
			}

			if strings.TrimSpace(text) == "exit" {
				confirmExit()
				return
			}

			// Persistence and Memory
			cmdHistory = append(cmdHistory, text)
			SaveHistory(text)
			historyIndex = len(cmdHistory)

			switchTo("logs")
			// Must execute in goroutine to prevent blocking input loop
			go engine.ExecuteCommand(text)
		}
	})

	startAsyncEngines()

	if err := app.SetRoot(mainFlex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

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
				app.SetFocus(cmdInput)
			}
		})

	pages.AddPage("modal", modal, false, true)
	app.SetFocus(modal)
}

func switchTo(page string) {
	updateTabs(page)
	pages.SwitchToPage(page)
}

func updateTabs(active string) {
	tabs := []string{"LOGS (F1)", "MAP (F2)", "LOOT (F3)", "TRAFFIC (F4)", "CTX (F5)", "NEURAL (F6)", "REPORT (F7)"}
	descs := []string{"System", "Recon", "Exfil", "Sniffer", "Intel", "AI-Ops", "Debrief"}
	ids := []string{"logs", "map", "loot", "traffic", "ai", "neuro", "report"}

	var topRow, bottomRow []string

	for i, t := range tabs {
		style := "[black:blue]"
		if active == ids[i] {
			style = "[black:aqua]"
		}
		topRow = append(topRow, fmt.Sprintf("%s┠ %s ┨[-]", style, t))
		padLen := (len(t) + 4 - len(descs[i])) / 2
		if padLen < 0 {
			padLen = 0
		}
		pad := strings.Repeat(" ", padLen)
		bottomRow = append(bottomRow, fmt.Sprintf("[gray]%s%s%s[-]", pad, descs[i], pad))
	}

	headerText := fmt.Sprintf(`[aqua:black:b]
██╗   ██╗ █████╗ ██████╗  ██████╗ ██████╗ ████████╗██████╗  █████╗  ██████╗███████╗
██║   ██║██╔══██╗██╔══██╗██╔═══██╗██╔══██╗╚══██╔══╝██╔══██╗██╔══██╗██╔════╝██╔════╝
╚██╗ ██╔╝███████║██████╔╝██║   ██║██████╔╝   ██║   ██████╔╝███████║██║     █████╗  
 ╚████╔╝ ██╔══██║██╔═══╝ ██║   ██║██╔══██╗   ██║   ██╔══██╗██╔══██║██║     ██╔══╝  
  ╚██╔╝  ██║  ██║██║     ╚██████╔╝██║  ██║   ██║   ██║  ██║██║  ██║╚██████╗███████╗
   ╚═╝   ╚═╝  ╚═╝╚═╝      ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚══════╝[-]
%s
%s`, strings.Join(topRow, " "), strings.Join(bottomRow, " "))

	header.SetText(headerText)
}

// startAsyncEngines consolidates all UI listeners
func startAsyncEngines() {
	// Status Bar Spinner
	go func() {
		ticker := time.NewTicker(250 * time.Millisecond)
		for range ticker.C {
			app.QueueUpdateDraw(func() {
				// Spinner
				spinnerIdx = (spinnerIdx + 1) % len(spinnerFrames)

				// TASK 1: Dynamic Status Indicators
				intStatus := "[gray][Ctrl+I: INT-OFF]"
				if logic.InterceptorActive {
					intStatus = "[black:green] Ctrl+I: INTERCEPTING (ACTIVE) [-:-]"
				}

				aiStatus := "[gray][F6: AI-IDLE]"
				if logic.GlobalNeuro.Active {
					aiStatus = "[black:magenta] NEURO ENGINE: ACTIVE (HYBRID) [-:-]"
				} else {
					aiStatus = "[black:red] NEURO ENGINE: OFF [-:-]"
				}

				statusFooter.SetText(fmt.Sprintf(" %s %s [blue]SYNC %s [white]| %s", intStatus, aiStatus, spinnerFrames[spinnerIdx], time.Now().Format("15:04:05")))

				// Pipeline Quadrant Update
				updatePipelineQuadrant()
			})
		}
	}()

	// Tactical Log Consumer (The BrainLog)
	go func() {
		for msg := range utils.UI_Log_Chan {
			app.QueueUpdateDraw(func() {
				// Special handling for Target Updates in sidebar
				if strings.Contains(msg, "Target Locked") {
					parts := strings.Split(msg, "Target Locked:[-] ")
					if len(parts) > 1 {
						// Clean escapes for the table cell
						url := strings.ReplaceAll(strings.TrimSpace(parts[1]), "[[", "[")
						targetColumn.SetCell(1, 1, tview.NewTableCell("[green]"+url))
					}
				}

				// Print to BrainLog
				fmt.Fprintf(brainLog, "[%s] %s\n", time.Now().Format("15:04:05"), msg)

				// Force scroll to end to ensure visibility of latest findings
				brainLog.ScrollToEnd()
			})
		}
	}()

	// 3. F2 (Map) Consumer
	go func() {
		for msg := range utils.MapDataChan {
			app.QueueUpdateDraw(func() {
				fmt.Fprintln(mapView, msg)
				mapView.ScrollToEnd()
			})
		}
	}()

	// 4. F3 (Loot) Consumer
	go func() {
		for pkt := range utils.LootDataChan {
			app.QueueUpdateDraw(func() {
				row := 1 // Header is 0, insert below header
				lootTable.InsertRow(row)
				lootTable.SetCell(row, 0, tview.NewTableCell(pkt.Type).SetTextColor(tcell.ColorRed))
				lootTable.SetCell(row, 1, tview.NewTableCell(pkt.Value).SetTextColor(tcell.ColorYellow))
				lootTable.SetCell(row, 2, tview.NewTableCell(pkt.Source).SetTextColor(tcell.ColorBlue))
			})
		}
	}()

	// 5. F4 (Traffic) Consumer
	go func() {
		for pkt := range utils.TrafficChan {
			app.QueueUpdateDraw(func() {
				reqView.SetText(fmt.Sprintf("[yellow]%s[-]\n\n[white]%s[-]", pkt.ReqHeader, pkt.ReqBody))
				resView.SetText(fmt.Sprintf("[green]%s[-]\n\n[white]%s[-]", pkt.ResHeader, pkt.ResBody))
			})
		}
	}()

	// 6. F5 (Context) Consumer
	go func() {
		for msg := range utils.ContextLogChan {
			app.QueueUpdateDraw(func() {
				fmt.Fprintln(aiView, msg)
				aiView.ScrollToEnd()
			})
		}
	}()

	// 7. F6 (Neural) Consumer
	go func() {
		for msg := range utils.NeuroLogChan {
			app.QueueUpdateDraw(func() {
				fmt.Fprintf(neuroView, "%s\n", msg)
				neuroView.ScrollToEnd()
			})
		}
	}()

	// 8. Interceptor Modal Listener
	go func() {
		for payload := range logic.InterceptorChan {
			app.QueueUpdateDraw(func() {
				utils.TacticalLog("[yellow]INTERCEPTOR:[-] Incoming Request Paused... Check Modal.")
				ShowInterceptorModal(app, pages, payload)
			})
		}
	}()
}
