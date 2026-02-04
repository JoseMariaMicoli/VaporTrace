package ui

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/engine"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
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
	mapTable     *tview.Table
	lootTable    *tview.Table
	reqView      *tview.TextView
	resView      *tview.TextView
	aiView       *tview.TextView
	neuroView    *tview.TextView

	// Removed reportFlex/reportView redeclarations.
	// They are defined in report_tab.go but visible here in package 'ui'.

	statusFooter *tview.TextView
	cmdInput     *tview.InputField

	cmdHistory   []string
	historyIndex int
	historyFile  = ".vapor_history"

	knownCommands = []string{
		"ask",
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

// SaveHistory appends a new command to the history file
func SaveHistory(cmd string) {
	f, err := os.OpenFile(historyFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(cmd + "\n")
}

// InitTacticalDashboard is the entry point for the TUI
func InitTacticalDashboard() {
	utils.SetLoggerMode("TUI")
	LoadHistory()

	logic.InitializeRotaryClient()
	logic.StartContextAggregator()

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
	targetColumn.SetTitle(" [blue]PIPELINE & STATUS [white]").SetBorder(true)
	targetColumn.SetCell(0, 0, tview.NewTableCell("[black:blue] PROPERTY "))
	targetColumn.SetCell(0, 1, tview.NewTableCell("[black:blue] VALUE "))
	updatePipelineQuadrant()

	// --- VIEW SETUP ---
	brainLog = tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetWordWrap(true).SetScrollable(true)
	brainLog.SetTitle(" [green]VAPOR_LOGS (TACTICAL FEED) [white]").SetBorder(true)

	mapTable = tview.NewTable().SetBorders(true).SetBordersColor(tcell.ColorDarkCyan).SetSelectable(true, false)
	mapTable.SetTitle(" [blue]ATTACK_SURFACE (MAP) [white]").SetBorder(true)
	mapTable.SetCell(0, 0, tview.NewTableCell("[black:cyan] TIMESTAMP "))
	mapTable.SetCell(0, 1, tview.NewTableCell("[black:cyan] SOURCE "))
	mapTable.SetCell(0, 2, tview.NewTableCell("[black:cyan] ENDPOINT "))
	mapTable.SetCell(0, 3, tview.NewTableCell("[black:cyan] META "))
	mapTable.SetFixed(1, 0)

	lootTable = tview.NewTable().SetBorders(true).SetBordersColor(tcell.ColorDarkCyan).SetSelectable(true, false)
	lootTable.SetTitle(" [magenta]LOOT_VAULT [white]").SetBorder(true)
	lootTable.SetCell(0, 0, tview.NewTableCell("[black:cyan] TYPE "))
	lootTable.SetCell(0, 1, tview.NewTableCell("[black:cyan] VALUE "))
	lootTable.SetCell(0, 2, tview.NewTableCell("[black:cyan] SOURCE "))
	lootTable.SetFixed(1, 0)

	reqView = tview.NewTextView().SetDynamicColors(true).SetWordWrap(true).SetScrollable(true)
	reqView.SetTitle(" [yellow]REQUEST (UPPER) - Ctrl+A to Analyze [white]").SetBorder(true)
	resView = tview.NewTextView().SetDynamicColors(true).SetWordWrap(true).SetScrollable(true)
	resView.SetTitle(" [green]RESPONSE (LOWER) [white]").SetBorder(true)

	trafficSplit := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(reqView, 0, 1, false).
		AddItem(resView, 0, 1, false)

	// F4 Input Capture (AI Trigger) with VISUAL FEEDBACK
	trafficSplit.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlA {
			req := reqView.GetText(true)
			res := resView.GetText(true)
			if req == "" {
				utils.TacticalLog("[yellow]NEURO:[-] No request selected to analyze.")
			} else {
				reqView.SetTitle(" [white:red]>>> ANALYZING SNAPSHOT... PLEASE WAIT <<<[white] ")
				utils.TacticalLog("[magenta]NEURO:[-] Snapshot captured. Transmitting to Neural Engine...")

				go func() {
					logic.GlobalNeuro.AnalyzeTrafficSnapshot(req, res)
				}()
			}
			return nil
		}
		return event
	})

	aiView = tview.NewTextView().SetDynamicColors(true).SetWordWrap(true).SetScrollable(true)
	aiView.SetTitle(" [white:blue] CONTEXT_AGGREGATOR (F5) [white] ").SetBorder(true)

	neuroView = tview.NewTextView().SetDynamicColors(true).SetWordWrap(true).SetScrollable(true)
	neuroView.SetTitle(" [magenta:black] NEURAL ENGINE (F6) [white] ").SetBorder(true)

	// F7 Report Tab Initialization
	// Note: reportFlex is variable in ui package, initialized here
	reportFlex = InitReportTab()

	// Add Pages
	pages.AddPage("logs", brainLog, true, true)
	pages.AddPage("map", mapTable, true, false)
	pages.AddPage("loot", lootTable, true, false)
	pages.AddPage("traffic", trafficSplit, true, false)
	pages.AddPage("ai", aiView, true, false)
	pages.AddPage("neuro", neuroView, true, false)
	pages.AddPage("report", reportFlex, true, false)

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

	// Global Key Bindings
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
		case tcell.KeyF6:
			switchTo("neuro")
		case tcell.KeyF7:
			LoadFindings()
			switchTo("report")

		case tcell.KeyCtrlI:
			logic.InterceptorActive = !logic.InterceptorActive
			utils.TacticalLog(fmt.Sprintf("INTERCEPTOR: %v", logic.InterceptorActive))
			updatePipelineQuadrant()

		case tcell.KeyCtrlH:
			ShowHelpModal(app, pages)
			return nil // PATCH: Prevents "h" from being typed in the input field

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

			cmdHistory = append(cmdHistory, text)
			SaveHistory(text)
			historyIndex = len(cmdHistory)

			if strings.HasPrefix(strings.ToLower(text), "ask ") {
				go engine.ExecuteCommand(text)
			} else {
				switchTo("logs")
				go engine.ExecuteCommand(text)
			}
		}
	})

	startAsyncEngines()

	if err := app.SetRoot(mainFlex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func updatePipelineQuadrant() {
	t := logic.CurrentSession.GetTarget()
	if t == "" {
		t = "[red]NOT SET"
	} else {
		t = "[green]" + t
	}
	targetColumn.SetCell(1, 0, tview.NewTableCell("TARGET"))
	targetColumn.SetCell(1, 1, tview.NewTableCell(t))

	aToken := logic.CurrentSession.AttackerToken
	if aToken == "" {
		aToken = "[gray]None"
	} else {
		aToken = "[green]" + shortString(aToken, 15)
	}
	targetColumn.SetCell(2, 0, tview.NewTableCell("AUTH (ATK)"))
	targetColumn.SetCell(2, 1, tview.NewTableCell(aToken))

	vToken := logic.CurrentSession.VictimToken
	if vToken == "" {
		vToken = "[gray]None"
	} else {
		vToken = "[yellow]" + shortString(vToken, 15)
	}
	targetColumn.SetCell(3, 0, tview.NewTableCell("AUTH (VIC)"))
	targetColumn.SetCell(3, 1, tview.NewTableCell(vToken))

	count := 0
	if db.DB != nil && logic.CurrentSession.GetTarget() != "" {
		_ = db.DB.QueryRow("SELECT COUNT(*) FROM context_store WHERE ? LIKE '%' || scope || '%'", logic.CurrentSession.GetTarget()).Scan(&count)
	}
	ctxColor := "[gray]"
	if count > 0 {
		ctxColor = "[magenta]"
	}
	targetColumn.SetCell(4, 0, tview.NewTableCell("CONTEXTS"))
	targetColumn.SetCell(4, 1, tview.NewTableCell(fmt.Sprintf("%s%d Active", ctxColor, count)))

	staticProxy := logic.GetConfiguredProxy()
	if staticProxy == "" {
		staticProxy = "[gray]Direct"
	} else {
		staticProxy = "[blue]" + staticProxy
	}
	targetColumn.SetCell(5, 0, tview.NewTableCell("PROXY (STAT)"))
	targetColumn.SetCell(5, 1, tview.NewTableCell(staticProxy))

	poolCount := len(logic.ProxyPool)
	poolStatus := fmt.Sprintf("[gray]%d Nodes", poolCount)
	if poolCount > 0 {
		poolStatus = fmt.Sprintf("[green]%d Active", poolCount)
	}
	targetColumn.SetCell(6, 0, tview.NewTableCell("PROXY POOL"))
	targetColumn.SetCell(6, 1, tview.NewTableCell(poolStatus))

	intStatus := "[red]OFF (Ctrl+I)"
	if logic.InterceptorActive {
		intStatus = "[green]ACTIVE"
	}
	targetColumn.SetCell(7, 0, tview.NewTableCell("INTERCEPTOR"))
	targetColumn.SetCell(7, 1, tview.NewTableCell(intStatus))

	neuroStatus := "[red]OFF"
	if logic.GlobalNeuro.Active {
		neuroStatus = "[magenta]ONLINE (HYBRID)"
	}
	targetColumn.SetCell(8, 0, tview.NewTableCell("NEURO BRAIN"))
	targetColumn.SetCell(8, 1, tview.NewTableCell(neuroStatus))
}

func shortString(s string, l int) string {
	if len(s) > l {
		return s[:l] + "..."
	}
	return s
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
	// Task 3: Focus routing logic for Report Tab
	if page == "report" {
		if reportFlex != nil {
			app.SetFocus(reportFlex)
		}
	} else {
		app.SetFocus(cmdInput)
	}
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

func startAsyncEngines() {
	go func() {
		ticker := time.NewTicker(250 * time.Millisecond)
		for range ticker.C {
			app.QueueUpdateDraw(func() {
				spinnerIdx = (spinnerIdx + 1) % len(spinnerFrames)
				statusFooter.SetText(fmt.Sprintf(" [blue]SYSTEM SYNC %s [white]| %s", spinnerFrames[spinnerIdx], time.Now().Format("15:04:05")))
				updatePipelineQuadrant()
			})
		}
	}()

	go func() {
		for msg := range utils.UI_Log_Chan {
			app.QueueUpdateDraw(func() {
				if msg == "___CLEAR_SCREEN_SIGNAL___" {
					brainLog.Clear()
					return
				}
				fmt.Fprintln(brainLog, msg)
				brainLog.ScrollToEnd()
			})
		}
	}()

	regexMap := regexp.MustCompile(`DISCOVERED\[-\] (.*?) \[yellow\]\((.*?)\)`)
	go func() {
		for msg := range utils.MapDataChan {
			app.QueueUpdateDraw(func() {
				cleanMsg := utils.StripANSI(msg)
				matches := regexMap.FindStringSubmatch(msg)
				endpoint := cleanMsg
				meta := "Unknown"
				source := "Recon"
				if len(matches) > 2 {
					endpoint = strings.TrimSpace(matches[1])
					meta = strings.TrimSpace(matches[2])
					if strings.Contains(meta, "JS") {
						source = "Scraper"
					}
					if strings.Contains(meta, "Swagger") || strings.Contains(meta, "OpenAPI") {
						source = "Swagger"
					}
					if strings.Contains(meta, "Mining") {
						source = "Miner"
					}
				}
				row := 1
				mapTable.InsertRow(row)
				mapTable.SetCell(row, 0, tview.NewTableCell(time.Now().Format("15:04:05")).SetTextColor(tcell.ColorGray))
				mapTable.SetCell(row, 1, tview.NewTableCell(source).SetTextColor(tcell.ColorBlue))
				mapTable.SetCell(row, 2, tview.NewTableCell(endpoint).SetTextColor(tcell.ColorGreen))
				mapTable.SetCell(row, 3, tview.NewTableCell(meta).SetTextColor(tcell.ColorYellow))
			})
		}
	}()

	go func() {
		for pkt := range utils.LootDataChan {
			app.QueueUpdateDraw(func() {
				row := 1
				lootTable.InsertRow(row)
				lootTable.SetCell(row, 0, tview.NewTableCell(pkt.Type).SetTextColor(tcell.ColorRed))
				lootTable.SetCell(row, 1, tview.NewTableCell(pkt.Value).SetTextColor(tcell.ColorYellow))
				lootTable.SetCell(row, 2, tview.NewTableCell(pkt.Source).SetTextColor(tcell.ColorBlue))
			})
		}
	}()

	go func() {
		for pkt := range utils.TrafficChan {
			app.QueueUpdateDraw(func() {
				reqView.SetText(fmt.Sprintf("[yellow]%s[-]\n\n[white]%s[-]", pkt.ReqHeader, pkt.ReqBody))
				resView.SetText(fmt.Sprintf("[green]%s[-]\n\n[white]%s[-]", pkt.ResHeader, pkt.ResBody))
				// Reset title if it was stuck on processing
				reqView.SetTitle(" [yellow]REQUEST (UPPER) - Ctrl+A to Analyze [white] ")
			})
		}
	}()

	go func() {
		for msg := range utils.ContextLogChan {
			app.QueueUpdateDraw(func() {
				fmt.Fprintln(aiView, msg)
				aiView.ScrollToEnd()
			})
		}
	}()

	go func() {
		for msg := range utils.NeuroLogChan {
			app.QueueUpdateDraw(func() {
				fmt.Fprintf(neuroView, "%s\n", msg)
				neuroView.ScrollToEnd()
				reqView.SetTitle(" [yellow]REQUEST (UPPER) - Ctrl+A to Analyze [white] ")
			})
		}
	}()

	go func() {
		for payload := range logic.InterceptorChan {
			app.QueueUpdateDraw(func() {
				utils.TacticalLog("[yellow]INTERCEPTOR:[-] Incoming Request Paused... Check Modal.")
				ShowInterceptorModal(app, pages, payload)
			})
		}
	}()
}
