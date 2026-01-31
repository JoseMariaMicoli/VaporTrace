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
	statusFooter *tview.TextView
	cmdInput     *tview.InputField

	// Command History
	cmdHistory   []string
	historyIndex int
	historyFile  = ".vapor_history"

	knownCommands = []string{
		"auth", "sessions", "map", "swagger", "scrape", "mine", "proxy", "proxies", "target", "pipeline",
		"flow", "bola", "bopla", "bfla", "exhaust", "ssrf", "audit", "probe",
		"weaver", "loot", "test-bola", "test-bopla", "test-bfla", "test-exhaust", "test-ssrf", "test-audit", "test-probe",
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

	reqView = tview.NewTextView().SetDynamicColors(true).SetWordWrap(true)
	reqView.SetTitle(" [yellow]TRAFFIC_REQ [white]").SetBorder(true)

	resView = tview.NewTextView().SetDynamicColors(true).SetWordWrap(true)
	resView.SetTitle(" [green]TRAFFIC_RES [white]").SetBorder(true)

	trafficSplit := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(reqView, 0, 1, false).
		AddItem(resView, 0, 1, false)

	aiView = tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true)
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

	// --- GLOBAL KEYBINDS & MOUSE SCROLL ---
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

	// --- INPUT HISTORY & EXECUTION ---
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
	// Status Bar Spinner
	go func() {
		ticker := time.NewTicker(250 * time.Millisecond)
		for range ticker.C {
			app.QueueUpdateDraw(func() {
				spinnerIdx = (spinnerIdx + 1) % len(spinnerFrames)
				statusFooter.SetText(fmt.Sprintf(" [blue]SYNC %s [white]| %s", spinnerFrames[spinnerIdx], time.Now().Format("15:04:05")))
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
}
