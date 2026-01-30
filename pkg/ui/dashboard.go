package ui

import (
	"fmt"
	"strings"
	"time"

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

	knownCommands = []string{
		"auth", "sessions", "map", "swagger", "scrape", "mine", "proxy", "proxies", "target", "pipeline",
		"flow", "bola", "bopla", "bfla", "exhaust", "ssrf", "audit", "probe",
		"weaver", "loot", "test-bola", "test-bopla", "test-bfla", "test-exhaust", "test-ssrf", "test-audit", "test-probe",
		"init_db", "reset_db", "report", "clear", "exit",
	}

	spinnerIdx    = 0
	spinnerFrames = []string{"▰▱▱▱▱", "▰▰▱▱▱", "▰▰▰▱▱", "▰▰▰▰▱", "▰▰▰▰▰"}
)

func InitTacticalDashboard() {
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
		if len(currentText) == 0 { return nil }
		for _, cmd := range knownCommands {
			if strings.HasPrefix(strings.ToLower(cmd), strings.ToLower(currentText)) {
				entries = append(entries, cmd)
			}
		}
		return
	})

	targetColumn = tview.NewTable().SetBorders(true).SetBordersColor(tcell.ColorBlue)
	targetColumn.SetTitle(" [blue]PIPELINE [white]").SetBorder(true)
	targetColumn.SetCell(0, 0, tview.NewTableCell("[black:blue] ENDPOINT "))
	targetColumn.SetCell(0, 1, tview.NewTableCell("[black:blue] RISK "))

	// --- TYPE ASSERTION FIXES BELOW ---
	brainLog = tview.NewTextView().SetDynamicColors(true).SetWordWrap(true).SetChangedFunc(func() {
		app.Draw()
		brainLog.ScrollToEnd()
	})
	brainLog.SetTitle(" [green]VAPOR_LOGS [white]").SetBorder(true)

	mapView = tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	mapView.SetTitle(" [blue]ATTACK_SURFACE [white]").SetBorder(true)

	lootTable = tview.NewTable().SetBorders(true).SetBordersColor(tcell.ColorDarkCyan)
	lootTable.SetTitle(" [magenta]LOOT_VAULT [white]").SetBorder(true)

	// Fix: Asserting (*tview.Box) back to (*tview.TextView)
	reqView = tview.NewTextView().SetDynamicColors(true)
	reqView.SetTitle(" [yellow]TRAFFIC_REQ [white]").SetBorder(true)
	
	resView = tview.NewTextView().SetDynamicColors(true)
	resView.SetTitle(" [green]TRAFFIC_RES [white]").SetBorder(true)
	
	trafficSplit := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(reqView, 0, 1, false).
		AddItem(resView, 0, 1, false)

	aiView = tview.NewTextView().SetDynamicColors(true)
	aiView.SetTitle(" [white:blue] LOGIC_ANALYZER [white] ").SetBorder(true)

	pages.AddPage("logs", brainLog, true, true)
	pages.AddPage("map", mapView, true, false)
	pages.AddPage("loot", lootTable, true, false)
	pages.AddPage("traffic", trafficSplit, true, false)
	pages.AddPage("ai", aiView, true, false)

	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 10, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(targetColumn, 35, 1, false).
			AddItem(pages, 0, 4, false),
		0, 4, false).
		AddItem(statusFooter, 1, 1, false).
		AddItem(cmdInput, 3, 1, true)

	updateTabs("logs")

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1: switchTo("logs")
		case tcell.KeyF2: switchTo("map")
		case tcell.KeyF3: switchTo("loot")
		case tcell.KeyF4: switchTo("traffic")
		case tcell.KeyF5: switchTo("ai")
		case tcell.KeyEsc: initiateScorchedEarth()
		}
		return event
	})

	cmdInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			handleCommand(cmdInput.GetText())
			cmdInput.SetText("")
		}
	})

	startAsyncEngines()
	if err := app.SetRoot(mainFlex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func handleCommand(cmd string) {
	if cmd == "" { return }
	fields := strings.Fields(cmd)
	verb := strings.ToLower(fields[0])

	switch verb {
	case "auth":
		utils.TacticalLog("[yellow]AUTH:[-] Setting identity tokens.")
	case "sessions":
		utils.TacticalLog("[blue]SESSION:[-] Viewing active tokens.")
	case "map":
		switchTo("map")
		utils.TacticalLog("[blue]RECON:[-] Executing Phase 2 Endpoint Mapping.")
	case "swagger", "scrape", "mine":
		utils.TacticalLog(fmt.Sprintf("[blue]DISCOVERY:[-] Running %s recon module.", verb))
	case "proxy", "proxies":
		utils.TacticalLog("[orange]TRAFFIC:[-] Modifying proxy state.")
	case "target":
		targetColumn.SetCell(1, 0, tview.NewTableCell("[yellow]LOCKED"))
		utils.TacticalLog("[aqua]PIPELINE:[-] Base URL Locked.")
	case "pipeline":
		utils.TacticalLog("[aqua]PIPELINE:[-] Categorizing targets.")
	case "flow":
		utils.TacticalLog("[red]LOGIC:[-] Managing sequences.")
	case "bola", "bopla", "bfla", "exhaust", "ssrf", "audit", "probe":
		utils.TacticalLog(fmt.Sprintf("[red]EXPLOIT:[-] Launching %s logic probe.", verb))
	case "weaver":
		utils.TacticalLog("[magenta]AGENT:[-] Deploying Ghost-Weaver.")
	case "loot":
		switchTo("loot")
		utils.TacticalLog("[magenta]VAULT:[-] Listing secrets.")
	case "test-bola", "test-bopla", "test-bfla", "test-exhaust", "test-ssrf", "test-audit", "test-probe":
		utils.TacticalLog("[green]VERIFY:[-] Testing verification engine.")
	case "init_db", "reset_db", "report":
		utils.TacticalLog("[white]SYS:[-] Persistence operation.")
	case "clear":
		brainLog.Clear()
	case "exit":
		initiateScorchedEarth()
	case "help", "usage":
		printTacticalManual()
	default:
		utils.TacticalLog("EXEC> " + cmd)
	}
}

func printTacticalManual() {
	switchTo("logs")
	manual := `
 [aqua:black:b] TACTICAL COMMAND MANUAL [-:-:-]
 [yellow]auth[-]     | Set identity tokens
 [yellow]map[-]      | Recon Mapping (F2)
 [yellow]target[-]   | Lock base URL
 [yellow]flow[-]     | Sequence management
 [yellow]bola/ssrf[-]| Launch logic probes
 [yellow]loot[-]     | View secrets (F3)
 [yellow]exit[-]     | Secure shutdown
`
	fmt.Fprintf(brainLog, "\n%s\n", manual)
}

func initiateScorchedEarth() {
	switchTo("logs")
	brainLog.Clear()
	go func() {
		utils.TacticalLog("[red]SHUTDOWN INITIATED...[-]")
		time.Sleep(400 * time.Millisecond)
		app.Stop()
	}()
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
	go func() {
		for {
			time.Sleep(250 * time.Millisecond)
			app.QueueUpdateDraw(func() {
				spinnerIdx = (spinnerIdx + 1) % len(spinnerFrames)
				statusFooter.SetText(fmt.Sprintf(" [blue]SYNC %s [white]| %s", spinnerFrames[spinnerIdx], time.Now().Format("15:04:05")))
			})
		}
	}()
	go func() {
		for msg := range utils.UI_Log_Chan {
			app.QueueUpdateDraw(func() {
				fmt.Fprintf(brainLog, "[%s] [green]>[white] %s\n", time.Now().Format("15:04:05"), msg)
			})
		}
	}()
}