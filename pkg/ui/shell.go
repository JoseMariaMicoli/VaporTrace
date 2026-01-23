package ui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"strconv"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/discovery"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/report"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/pterm/pterm"
)

type Shell struct {
	Active       bool
	RemoteActive bool
}

func NewShell() *Shell {
	return &Shell{
		Active:       true,
		RemoteActive: true,
	}
}

// RenderBanner displays the header-based banner with an improved tactical aesthetic
func (s *Shell) RenderBanner() {
	fmt.Print("\033[H\033[2J") // Clear screen

	// Match the ASCII Art from README.md
	bannerStyle := pterm.NewStyle(pterm.FgCyan, pterm.Bold)
	bannerStyle.Println(`
    __  __                         _____                    
    \ \ / /___  _ __  ___  _ __   |_   _| __ __ _  ___ ___ 
     \ V // _ ` + "`" + `| '_ \/ _ \| '__|    | || '__/ _` + "`" + `|/ __/ _ \
      \  / (_| | |_)  (_) | |       | || | | (_| | (_|  __/
       \/ \__,_| .__/\___/|_|       |_||_|  \__,_|\___\___|
               |_|      [ Surgical API Exploitation Suite]`)

	pterm.Println(pterm.Cyan("────────────────────────────────────────────────────────────"))

	statusColor := pterm.FgLightGreen
	statusText := "● SYSTEM ONLINE"

	if !s.RemoteActive {
		statusColor = pterm.FgLightRed
		statusText = "○ LINK SEVERED"
	}

	// Check if proxy is active for the UI
	gateway := "DIRECT"
	if len(logic.ProxyPool) > 0 {
		gateway = fmt.Sprintf("ROTATING (%d IPs)", len(logic.ProxyPool))
	} else if os.Getenv("HTTP_PROXY") != "" {
		gateway = "http://127.0.0.1:8080 (BURP)"
	}

	// Stylized Table using the logic from your shell.go
	pterm.DefaultTable.WithData(pterm.TableData{
		{"UPSTREAM GATEWAY", "LOGIC ENGINE", "BUILD VERSION"},
		{pterm.LightBlue(gateway), statusColor.Sprintf(statusText), pterm.LightMagenta("v3.1-Flash")},
	}).WithBoxed().Render()

	pterm.Printf("\n%s Use 'usage' for tactics or 'help' for manuals.\n\n",
		pterm.Cyan("»"))
}

// Start launches the interactive tactical loop with Auto-Completion
func (s *Shell) Start() {
	s.RenderBanner()

	completer := readline.NewPrefixCompleter(
		readline.PcItem("init_db"),
		readline.PcItem("reset_db"),
		readline.PcItem("target"),
		readline.PcItem("map"),
		readline.PcItem("mine"),
		readline.PcItem("scrape"),
		readline.PcItem("swagger"),
		readline.PcItem("pipeline"),
		readline.PcItem("proxy"),
		readline.PcItem("proxies",        // <--- Ensure comma here
			readline.PcItem("load"),  // <--- Ensure comma here
			readline.PcItem("reset"), // <--- Ensure comma here
		),
		readline.PcItem("flow",
	        readline.PcItem("add"),
	        readline.PcItem("run"),
	        readline.PcItem("list"),
	        readline.PcItem("step"),
	        readline.PcItem("race"),
	        readline.PcItem("clear"),
	    ),
	    readline.PcItem("loot",
		    readline.PcItem("list"),
		    readline.PcItem("clear"),
		),
		readline.PcItem("bola"),
		readline.PcItem("bopla"),
		readline.PcItem("bfla"),
		readline.PcItem("exhaust"),
		readline.PcItem("ssrf"),
		readline.PcItem("audit"),
		readline.PcItem("probe"),
		readline.PcItem("test-probe"),
		readline.PcItem("test-audit"),
		readline.PcItem("test-ssrf"),
		readline.PcItem("test-exhaust"),
		readline.PcItem("test-bola"),
		readline.PcItem("test-bopla"),
		readline.PcItem("test-bfla"),
		readline.PcItem("auth"),
		readline.PcItem("sessions"),
		readline.PcItem("help"),
		readline.PcItem("usage"),
		readline.PcItem("clear"),
		readline.PcItem("splash"),
		readline.PcItem("report"),
		readline.PcItem("exit"),
	)

	statusStr := pterm.NewStyle(pterm.FgGreen, pterm.Bold).Sprint("ONLINE")
	if !s.RemoteActive {
		statusStr = pterm.NewStyle(pterm.FgRed, pterm.Bold).Sprint("OFFLINE")
	}

	prompt := fmt.Sprintf("[%s] %s%s%s ",
		statusStr,
		color.New(color.FgGreen, color.Bold).Sprint("vapor@trace"),
		color.New(color.FgWhite).Sprint(":"),
		color.New(color.FgBlue, color.Bold).Sprint("~$ "),
	)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          prompt,
		HistoryFile:      "/tmp/vaportrace.tmp",
		AutoComplete:     completer,
		InterruptPrompt: "^C",
		EOFPrompt:        "exit",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for s.Active {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		command := parts[0]
		var args []string
		if len(parts) > 1 {
			args = parts[1:]
		}

		s.handleCommand(command, args)
	}
}

func (s *Shell) handleCommand(command string, args []string) {
	switch command {
	case "loot":
		if len(args) < 1 {
			pterm.Info.Println("Usage: loot list | loot clear")
			return
		}

		subCmd := args[0]
		switch subCmd {
		case "list":
			pterm.DefaultHeader.WithFullWidth(false).
				WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).
				Println("DISCOVERY VAULT")

			if len(logic.Vault) == 0 {
				pterm.Warning.Println("No sensitive data recovered yet.")
			} else {
				for _, f := range logic.Vault {
					pterm.Info.Printfln("[%s] %s (Source: %s)",
						pterm.Yellow(f.Type),
						f.Value,
						pterm.Gray(f.Source))
				}
			}

		case "clear":
			logic.Vault = []logic.Finding{}
			pterm.Success.Println("Vault purged.")

		default:
			pterm.Error.Printfln("Unknown sub-command: %s", subCmd)
		}

	case "flow":
		if len(args) < 1 {
			pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).Println("TACTICAL FLOW ENGINE")
			pterm.Info.Println("Phase 7.1/7.2: State-Machine & Sequence Analysis")
			fmt.Println("\nAvailable Subcommands:")
			fmt.Println("  add   - Record a new step (Interactive)")
			fmt.Println("  run   - Execute the current sequence")
			fmt.Println("  list  - View the recorded sequence and context")
			fmt.Println("  step  - [Phase 7.2] Execute a single step (Out-of-order test)")
			fmt.Println("  clear - Purge the current flow and captured context")
			return
		}

		switch args[0] {
		case "add":
			pterm.DefaultSection.Println("Flow Step Recorder")
			name, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Step Name").Show()
			method, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Method (GET/POST/PUT)").Show()
			
			rawURL, _ := pterm.DefaultInteractiveTextInput.
				WithDefaultText("Target URL").
				WithDefaultValue(logic.CurrentSession.TargetURL).
				Show()
				
			rawBody, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Body (Use {{var}} for variables)").Show()
			extract, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("JSON Path to Extract (optional)").Show()

			// Sanitization for Phase 6.3 ANSI Safety
			cleanURL := utils.StripANSI(rawURL)
			cleanBody := utils.StripANSI(rawBody)
			cleanExtract := utils.StripANSI(extract)

			logic.ActiveFlow = append(logic.ActiveFlow, logic.FlowStep{
				Name:        name,
				Method:      strings.ToUpper(method),
				URL:         strings.TrimSpace(cleanURL),
				Body:        strings.TrimSpace(cleanBody),
				ExtractPath: strings.TrimSpace(cleanExtract),
			})
			pterm.Success.Printfln("Step '%s' added to sequence.", name)

		case "run":
			logic.RunFlow()

		case "list":
			pterm.DefaultSection.Println("Current Tactical Flow")
			if len(logic.ActiveFlow) == 0 {
				pterm.Warning.Println("Queue is empty.")
				return
			}
			for i, s := range logic.ActiveFlow {
				pterm.BulletListPrinter{Items: []pterm.BulletListItem{
					{Level: 0, Text: pterm.Cyan(fmt.Sprintf("[%d] %s", i+1, s.Name))},
					{Level: 1, Text: pterm.Gray(fmt.Sprintf("%s %s", s.Method, s.URL))},
				}}.Render()
			}

		case "step":
			// Phase 7.2: Out-of-order State Machine Test
			if len(args) < 2 {
				pterm.Error.Println("Usage: flow step <id>")
				return
			}
			id, err := strconv.Atoi(args[1])
			if err != nil {
				pterm.Error.Println("Invalid step ID.")
				return
			}
			logic.RunStep(id - 1)
		case "race":
			if len(args) < 3 {
				pterm.Info.Println("Usage: flow race <step_id> <threads>")
				return
			}
			id, _ := strconv.Atoi(args[1])
			threads, _ := strconv.Atoi(args[2])
			logic.RunRace(id-1, threads)

		case "clear":
			logic.ActiveFlow = []logic.FlowStep{}
			logic.FlowContext = make(map[string]string) // Reset State-Machine memory
			pterm.Success.Println("Flow and State Context purged.")

		default:
			pterm.Error.Printfln("Unknown flow command: %s", args[0])
		}

	case "target":
		if len(args) < 1 {
			pterm.Error.Println("Usage: target <url> (e.g., target https://api.target.com)")
			return
		}
		// Normalize and set the target in the global logic store
		targetURL := args[0]
		if !strings.HasPrefix(targetURL, "http") {
			targetURL = "https://" + targetURL
		}
		logic.CurrentSession.TargetURL = targetURL
		pterm.Success.Printfln("Target locked: %s", targetURL)

	case "proxies":
		if len(args) < 1 {
			pterm.Info.Println("Usage: proxies load <file> | proxies reset")
			return
		}
		if args[0] == "load" && len(args) == 2 {
			err := logic.LoadProxiesFromFile(args[1])
			if err == nil {
				logic.InitializeRotaryClient()
			}
		} else if args[0] == "reset" {
			logic.ProxyPool = []string{}
			logic.InitializeRotaryClient()
			pterm.Success.Println("Proxy pool purged. Identity returned to default (Direct/Burp).")
		} else {
			pterm.Warning.Println("Invalid subcommand. Use 'load <file>' or 'reset'.")
		}

	case "proxy":
		if len(args) < 1 {
			pterm.Info.Println("Usage: proxy <url> (e.g., proxy http://127.0.0.1:8080)")
			pterm.Info.Println("Current Status: Traffic is direct")
			return
		}

		proxyAddr := args[0]
		client, err := utils.GetClient(proxyAddr) 
		if err != nil {
			pterm.Error.Printf("Failed to initialize proxy: %v\n", err)
			return
		}

		logic.SetGlobalClient(client) 
		discovery.SetGlobalClient(client)

		pterm.Success.WithPrefix(pterm.Prefix{Text: "BURP", Style: pterm.NewStyle(pterm.FgBlack, pterm.BgCyan)}).
			Printf("Proxy active at %s. All tactical traffic now routed through interceptor.\n", proxyAddr)

	case "report":
		fmt.Println()
		pterm.Info.Println("Accessing persistence layer for debrief construction...")
		time.Sleep(600 * time.Millisecond)

		pterm.Print(pterm.LightBlue("  ● "))
		pterm.Print("Aggregating Phase II - IV finding buffers... ")
		time.Sleep(800 * time.Millisecond)
		pterm.Success.Println("Done")

		pterm.Print(pterm.LightCyan("  ● "))
		pterm.Print("Generating Markdown Debrief... ")
		report.GenerateMissionDebrief()

		fmt.Println()
		pterm.Success.WithPrefix(pterm.Prefix{Text: "REPORT", Style: pterm.NewStyle(pterm.FgBlack, pterm.BgGreen)}).
			Println("Mission Intelligence has been successfully exported.")

	case "usage":
		if len(args) > 0 && args[0] == "loot" {
			pterm.DefaultHeader.WithFullWidth(false).Println("PHASE 8.1: PII SCANNER")
			pterm.Println("Automatically scans all incoming HTTP traffic for secrets.")
			pterm.Println("Currently monitoring: AWS Keys, JWTs, Credit Cards, Emails, and Slack Tokens.")
		} else {
			s.ShowUsage()
		}

	case "help":
		if len(args) > 0 {
			s.ShowHelp(args[0])
		} else {
			pterm.Info.Println("Usage: help <command>")
		}

	case "pipeline":
		concurrency := logic.CurrentSession.Threads
		if len(args) > 1 {
			if c, err := strconv.Atoi(args[1]); err == nil {
				concurrency = c
			}
		}
		
		pterm.DefaultHeader.WithFullWidth(false).WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).Println("PIPELINE: Industrialized Execution")
		logic.AnalyzeDiscovery()
		pterm.Info.Println("Starting automated engine sequence via Pipeline...")
		logic.RunPipeline(concurrency)

	case "clear", "cls", "splash":
		s.RenderBanner()

	case "exit":
		result, _ := pterm.DefaultInteractiveConfirm.
			WithDefaultText("Terminate mission and exit shell?").
			WithConfirmStyle(pterm.NewStyle(pterm.FgRed, pterm.Bold)).
			Show()

		if result {
			fmt.Println()
			pterm.Print(pterm.LightBlue("  ○ "))
			pterm.Print("Closing active network connections... ")
			time.Sleep(500 * time.Millisecond)
			pterm.Success.Println("Done")

			pterm.Print(pterm.LightMagenta("  ○ "))
			pterm.Print("Synchronizing mission logs and closing database... ")
			db.CloseDB()
			time.Sleep(1 * time.Second)
			pterm.Success.Println("Done")

			pterm.Print(pterm.LightYellow("  ○ "))
			pterm.Print("Exiting tactical session... ")
			time.Sleep(400 * time.Millisecond)
			pterm.Success.Println("Goodbye")

			fmt.Println()
			pterm.Info.WithPrefix(pterm.Prefix{Text: "OFFLINE", Style: pterm.NewStyle(pterm.FgBlack, pterm.BgCyan)}).
				Println("VaporTrace terminated safely. Reveal the invisible.")

			os.Exit(0)
		}

	case "init_db":
		fmt.Println()
		pterm.DefaultHeader.WithFullWidth(false).Println("Phase 5: Intelligence Initialization")
		pterm.Print(pterm.LightCyan("  ● "))
		pterm.Print("Establishing link to SQLite persistence... ")
		db.InitDB()
		time.Sleep(600 * time.Millisecond)
		pterm.Success.Println("Connected")

		pterm.Print(pterm.LightGreen("  ● "))
		pterm.Print("Spawning asynchronous log worker pool... ")
		go db.StartAsyncWorker()
		time.Sleep(800 * time.Millisecond)
		pterm.Success.Println("Active")

		fmt.Println()
		pterm.Info.WithPrefix(pterm.Prefix{Text: "READY", Style: pterm.NewStyle(pterm.FgBlack, pterm.BgGreen)}).
			Println("Database persistence is now online. All findings will be logged.")

	case "reset_db":
		pterm.Warning.WithPrefix(pterm.Prefix{Text: "DANGER", Style: pterm.NewStyle(pterm.FgBlack, pterm.BgRed)}).
			Println("This action will permanently delete all mission findings and database logs!")

		result, _ := pterm.DefaultInteractiveConfirm.
			WithDefaultText("Are you sure you want to PURGE the database?").
			WithConfirmStyle(pterm.NewStyle(pterm.FgRed, pterm.Bold)).
			Show()

		if result {
			fmt.Println()
			pterm.Print(pterm.LightRed("  × "))
			pterm.Print("Wiping all records from findings table... ")
			time.Sleep(600 * time.Millisecond)
			pterm.Print("\n  × ")
			pterm.Print("Resetting DATABASE ID and Gen Time... ")
			db.ResetDB()
			time.Sleep(800 * time.Millisecond)
			pterm.Success.Println("Done")
			fmt.Println()
			pterm.Success.WithPrefix(pterm.Prefix{Text: "CLEAN", Style: pterm.NewStyle(pterm.FgBlack, pterm.BgGreen)}).
				Println("Database has been successfully purged. System ready for new mission.")
		} else {
			pterm.Info.Println("Reset operation aborted by operator.")
		}

	case "swagger":
		if len(args) < 1 {
			pterm.Error.Println("Usage: swagger <url>")
			return
		}
		endpoints, err := discovery.ParseSwagger(args[0], "")
		if err != nil {
			pterm.Error.Printf("Parsing failed: %v\n", err)
			return
		}
		pterm.Success.Printf("Extracted %d endpoints. Probing live status...\n", len(endpoints))
		for _, ep := range endpoints {
			status, _ := discovery.ProbeEndpoint(args[0], ep, "")
			if status == 200 {
				pterm.Success.Printf("[LIVE] %s (Status: %d)\n", ep, status)
			}
		}

	case "mine":
		if len(args) < 2 {
			pterm.Error.Println("Usage: mine <url> <endpoint>")
			return
		}
		pterm.Info.Printf("Starting parameter miner on %s%s...\n", args[0], args[1])
		discovery.MineParameters(args[0], args[1], "")
		pterm.Success.Println("Mining operation complete.")

	case "scrape": 
		if len(args) < 1 {
			pterm.Error.Println("Usage: scrape <url>")
			return
		}
		paths, err := discovery.ExtractJSPaths(args[0], "")
		if err != nil {
			pterm.Error.Printf("Scraping failed: %v\n", err)
			return
		}
		if len(paths) == 0 {
			pterm.Warning.Println("No API paths discovered.")
			return
		}
		pterm.Success.Printf("Discovered %d potential API paths:\n", len(paths))
		for _, p := range paths {
			pterm.BulletListPrinter{Items: []pterm.BulletListItem{{Level: 0, Text: p}}}.Render()
		}

	case "map":
		if len(args) < 2 {
			pterm.Error.Println("Usage: map -u <url> OR map -j <js_url>")
			return 
		}

		var sURL, jURL string
		for i, arg := range args {
			if arg == "-u" && i+1 < len(args) { sURL = args[i+1] }
			if arg == "-j" && i+1 < len(args) { jURL = args[i+1] }
		}

		pterm.DefaultSection.Println("Phase 2: Intelligence Mapping")
		var foundEndpoints []string

		if jURL != "" {
			spinner, _ := pterm.DefaultSpinner.Start("Scraping JS Bundle: " + jURL)
			endpoints, err := discovery.ExtractJSPaths(jURL, "") 
			if err != nil {
				spinner.Fail("Scrape failed: " + err.Error())
			} else if len(endpoints) == 0 {
				spinner.Warning("No API patterns found in JS.")
			} else {
				spinner.Success(fmt.Sprintf("Harvested %d routes", len(endpoints)))
				foundEndpoints = append(foundEndpoints, endpoints...)
				tableData := pterm.TableData{{"TYPE", "EXTRACTED PATH"}}
				for _, e := range endpoints {
					tableData = append(tableData, []string{"JS_ROUTE", e})
				}
				pterm.DefaultTable.WithHasHeader().WithData(tableData).WithBoxed().Render()
			}
		}

		if sURL != "" {
			spinner, _ := pterm.DefaultSpinner.Start("Analyzing Swagger Spec...")
			endpoints, err := discovery.ParseSwagger(sURL, "")
			if err != nil {
				spinner.Fail("Swagger parse failed")
			} else {
				spinner.Success(fmt.Sprintf("Found %d documented endpoints", len(endpoints)))
				foundEndpoints = append(foundEndpoints, endpoints...)
			}
		}

		if len(foundEndpoints) > 0 {
			pterm.Success.Printf("Mapping complete. %d total endpoints stored in session.\n", len(foundEndpoints))
		}

	case "audit":
		if len(args) < 1 {
			pterm.Info.Println("Usage: audit <url>")
			return
		}
		probe := &logic.MisconfigContext{TargetURL: args[0]}
		probe.Audit()

	case "probe":
		if len(args) < 1 {
			pterm.Info.Println("Usage: probe <url> [type]")
			return
		}

		iType := "generic"
		if len(args) > 1 {
			iType = args[1]
		}

		// CRITICAL FIX: Append to 'ActiveFlow' so 'flow run' can see it
		logic.ActiveFlow = append(logic.ActiveFlow, logic.FlowStep{
			Name:   fmt.Sprintf("Probe-%s", iType),
			Method: "GET", 
			URL:    args[0],
			Body:   "",
		})
		
		pterm.Success.Printfln("Task added to Tactical Queue (Flow): [%s] %s", pterm.Cyan(iType), args[0])

		probe := &logic.IntegrationContext{TargetURL: args[0], IntegrationType: iType}
		probe.Probe()

	case "test-probe":
		pterm.Info.Println("Simulating Webhook injection against httpbin...")
		test := &logic.IntegrationContext{TargetURL: "https://httpbin.org/post", IntegrationType: "generic"}
		test.Probe()

	case "test-audit":
		pterm.Info.Println("Running diagnostic audit against google.com...")
		test := &logic.MisconfigContext{TargetURL: "https://www.google.com"}
		test.Audit()

	case "exhaust":
		if len(args) < 2 {
			pterm.Info.Println("Usage: exhaust <url> <parameter>")
			return
		}
		probe := &logic.ExhaustionContext{TargetURL: args[0], ParamName: args[1]}
		probe.FuzzPagination()

	case "ssrf":
		if len(args) < 3 {
			pterm.Info.Println("Usage: ssrf <url> <parameter> <callback_url>")
			return
		}
		probe := &logic.SSRFContext{TargetURL: args[0], ParamName: args[1], Callback: args[2]}
		probe.Probe()

	case "test-ssrf":
		pterm.Info.Println("Simulating SSRF against httpbin...")
		test := &logic.SSRFContext{
			TargetURL: "https://httpbin.org/redirect-to",
			ParamName: "url",
			Callback:  "https://google.com",
		}
		test.Probe()

	case "test-exhaust":
		pterm.Info.Println("Simulating Pagination Fuzzing against httpbin...")
		test := &logic.ExhaustionContext{TargetURL: "https://httpbin.org/get", ParamName: "limit"}
		test.FuzzPagination()

	case "bola":
		isPipeline := false
		for _, arg := range args {
			if arg == "--pipeline" || arg == "-p" { isPipeline = true }
		}

		if isPipeline {
			ctx := &logic.BOLAContext{BaseURL: logic.CurrentSession.TargetURL}
			idList := []string{"1", "2", "3", "101", "102"} 
			ctx.MassProbe(idList, logic.CurrentSession.Threads) 
			return
		}

		if len(args) < 2 {
			pterm.Error.Println("Usage: bola -u <url> -v <victim_id> [-a <attacker_id>] OR bola --pipeline")
			return
		}

		ctx := &logic.BOLAContext{}
		for i := 0; i < len(args); i++ {
			switch args[i] {
			case "-u":
				if i+1 < len(args) { ctx.BaseURL = args[i+1] }
			case "-v":
				if i+1 < len(args) { ctx.VictimID = args[i+1] }
			case "-a":
				if i+1 < len(args) { ctx.AttackerID = args[i+1] }
			}
		}
		ctx.ProbeSilent()

	case "scan-bola":
		if len(args) < 4 {
			pterm.Error.Println("Usage: scan-bola -u <url> -r <start-end> -t <threads>")
			break
		}

		var urlStr, idRange string
		threads := logic.CurrentSession.Threads
		for i := 0; i < len(args); i++ {
			switch args[i] {
			case "-u":
				if i+1 < len(args) { urlStr = args[i+1] }
			case "-r":
				if i+1 < len(args) { idRange = args[i+1] }
			case "-t":
				if i+1 < len(args) { fmt.Sscanf(args[i+1], "%d", &threads) }
			}
		}

		var start, end int
		fmt.Sscanf(idRange, "%d-%d", &start, &end)
		var ids []string
		for i := start; i <= end; i++ {
			ids = append(ids, fmt.Sprintf("%d", i))
		}

		ctx := &logic.BOLAContext{BaseURL: urlStr}
		ctx.MassProbe(ids, threads)

	case "bopla":
		if len(args) > 0 && (args[0] == "--pipeline" || args[0] == "-p") {
			logic.ExecuteMassBOPLA(logic.CurrentSession.Threads)
			return
		}
		pterm.Error.Println("Usage: bopla --pipeline")

	case "test-bopla":
		pterm.Info.Println("Simulating Mass Assignment against httpbin...")
		test := &logic.BOPLAContext{
			TargetURL: "https://httpbin.org/patch",
			Method:    "PATCH",
			BaseJSON:  `{"username": "vapor_user"}`,
		}
		test.RunFuzzer(1)

	case "bfla":
		if len(args) > 0 && (args[0] == "--pipeline" || args[0] == "-p") {
			logic.ExecuteMassBFLA(logic.CurrentSession.Threads)
			return
		}
		pterm.Error.Println("Usage: bfla --pipeline")

	case "test-bfla":
		pterm.Info.Println("Simulating Verb Tampering against httpbin...")
		ctx := &logic.BFLAContext{TargetURL: "https://httpbin.org/anything"}
		ctx.MassProbe(1)

	case "auth":
		if len(args) < 2 {
			pterm.Info.Println("Usage: auth <victim|attacker> <token>")
			return
		}
		if args[0] == "attacker" {
			logic.CurrentSession.AttackerToken = args[1]
		} else {
			logic.CurrentSession.VictimToken = args[1]
		}
		pterm.Success.Printf("%s token updated in session store.\n", args[0])

	case "sessions":
		pterm.DefaultTable.WithData(pterm.TableData{
			{"ROLE", "TOKEN SNAPSHOT"},
			{"VICTIM", logic.CurrentSession.VictimToken},
			{"ATTACKER", logic.CurrentSession.AttackerToken},
		}).WithBoxed().Render()

	case "test-bola":
		pterm.Info.Println("Diagnostic BOLA test against httpbin...")
		vuln := &logic.BOLAContext{
			BaseURL:        "https://httpbin.org/anything",
			VictimID:       "private_id_777",
			AttackerToken: logic.CurrentSession.AttackerToken,
		}
		vuln.ProbeSilent()

	default:
		pterm.Error.Printf("Unknown tactical command: %s\n", command)
	}
}

func (s *Shell) ShowUsage() {
	table := pterm.DefaultTable.WithHasHeader().WithBoxed()
	table.Data = pterm.TableData{
		{"COMMAND", "DESCRIPTION", "EXAMPLE"},
		{"init_db", "Initialize SQLite Persistence", "init_db"},
		{"reset_db", "Wipe local mission data (Purge)", "reset_db"},
		{"proxy", "Toggle Burp Suite Proxy (8080)", "proxy on"},
		{"proxies load", "Load IP rotation pool from text file", "proxies load p.txt"},
		{"proxies reset", "Clear pool and return to direct mode", "proxies reset"},
		{"flow", "Record and replay multi-step business logic", "flow add"},
		{"loot", "[Phase 8.1] View/Clear captured PII and Secrets", "loot list"},
		{"swagger", "Parse OpenAPI/Swagger docs for routes", "swagger <url>"},
		{"mine", "Fuzz for hidden query parameters", "mine <url> <endpoint>"},
		{"scrape", "Extract API paths from JS files", "scrape <url>"},
		{"pipeline", "Analyze discovery data for BOLA/BFLA/BOPLA targets"},
		{"auth", "Set identity tokens", "auth attacker <token>"},
		{"sessions", "View active tokens", "sessions"},
		{"bola", "Phase 3 BOLA test", "bola <url> <id>"},
		{"bopla", "BOPLA / API3 Mass Assignment", "bopla <url> '{\"id\":1}'"},
		{"bfla", "BFLA / API5 Method Shuffling", "bfla <url>"},
		{"exhaust", "API4 Pagination Fuzzing", "exhaust <url> limit"},
		{"ssrf", "API7 SSRF/OOB Tracker", "ssrf <url> param <callback>"},
		{"audit", "API8 Security Misconfig Audit", "audit <url>"},
		{"probe", "API10 Integration Fuzzer", "probe <url> stripe"},
		{"test-probe", "Verify Integration logic", "test-probe"},
		{"test-audit", "Verify Audit Logic", "test-audit"},
		{"test-ssrf", "Verify SSRF Logic", "test-ssrf"},
		{"test-exhaust", "Verify Exhaustion logic", "test-exhaust"},
		{"test-bola", "Verify BOLA logic", "test-bola"},
		{"test-bopla", "Verify BOPLA logic", "test-bopla"},
		{"test-bfla", "Verify BFLA logic", "test-bfla"},
		{"map", "Full Phase 2 API Recon", "map -u <url>"},
		{"report", "Generate the Report (.md)", "report"},
		{"help", "Show manual", "help map"},
		{"exit", "Shutdown suite", "exit"},
	}
	table.Render()
}

func (s *Shell) ShowHelp(cmd string) {
	pterm.DefaultHeader.WithFullWidth(false).Printf("Manual: %s\n", cmd)

	switch cmd {
	case "proxy":
		pterm.Bold.Println("DESCRIPTION:")
		pterm.Println("Enables or disables global traffic routing through a proxy.")
		pterm.Println("Defaults to Burp Suite at http://127.0.0.1:8080.")
	case "init_db":
		pterm.Bold.Println("DESCRIPTION:")
		pterm.Println("Initializes the Phase 5 SQLite engine and starts the async worker.")
		pterm.Println("This creates a persistent link between your tactical actions and the")
		pterm.Println("final 'Mission Debrief' report.")
		pterm.Bold.Println("\nTECHNICAL DETAILS:")
		pterm.BulletListPrinter{Items: []pterm.BulletListItem{
			{Level: 0, Text: "Engine: SQLite3 (Local Persistence)"},
			{Level: 0, Text: "I/O Mode: Asynchronous Non-blocking (Goroutine Worker)"},
			{Level: 0, Text: "Default File: ./vaportrace.db"},
		}}.Render()
	case "proxies":
		pterm.Bold.Println("PHASE 6.2: IP ROTATION & EVASION")
		pterm.Println("Distributes requests across a pool of HTTP/SOCKS5 proxies to bypass rate-limits.")
		pterm.Println("\nCOMMANDS:")
		pterm.BulletListPrinter{Items: []pterm.BulletListItem{
			{Level: 0, Text: "load <file> : Ingests a line-separated list of proxy URLs."},
			{Level: 0, Text: "reset        : Wipes the pool. VaporTrace will fall back to Burp or Direct."},
		}}.Render()
		pterm.Println("\nFILE FORMAT:")
		pterm.Cyan("http://user:pass@host:port\nsocks5://host:port")
	case "reset_db":
		pterm.Bold.Println("DESCRIPTION:")
		pterm.Println("Safely purges the mission database and metadata.")
		pterm.Println("Use this command to 'clean' the environment before a new operation")
		pterm.Println("or to reset the DATABASE ID to 1.")
		pterm.Bold.Println("\nWARNING:")
		pterm.NewStyle(pterm.FgRed, pterm.Bold).Println("This action is irreversible. All logged findings will be lost.")

	case "help":
		pterm.DefaultHeader.WithFullWidth(false).Println("VaporTrace Tactical Help")
		helpItems := [][]string{
			{"swagger", "Map API surface via OpenAPI/Swagger docs", "swagger <url>"},
			{"mine", "Fuzz for hidden query parameters", "mine <url> <endpoint>"},
			{"scrape", "Extract API paths from JS source", "scrape <url>"},
			{"flow", "Business Logic sequencing engine", "flow <add|run|list|clear>"},
		}
		for _, item := range helpItems {
			pterm.DefaultBulletList.WithItems([]pterm.BulletListItem{
				{Level: 0, Text: pterm.Cyan(item[0]) + ": " + item[1]},
				{Level: 1, Text: pterm.Gray("Usage: ") + item[2]},
			}).Render()
		}
	case "loot":
	    pterm.DefaultHeader.WithFullWidth(false).Println("PHASE 8.3: DISCOVERY VAULT & CLOUD PIVOT")
	    pterm.Bold.Println("DESCRIPTION:")
	    pterm.Println("Manages the storage of detected secrets and handles tactical exfiltration.")
	    pterm.Println("High-value findings are encrypted using AES-256-GCM (Ghost-Pipeline Standard).")
	    
	    pterm.Bold.Println("\nAUTOMATED CLOUD PIVOT:")
	    pterm.Println("Upon detection of 169.254.169.254, the engine spawns an IMDSv2 prober in the background.")

	    pterm.Bold.Println("\nCOMMANDS:")
	    pterm.BulletListPrinter{Items: []pterm.BulletListItem{
	        {Level: 0, Text: pterm.Cyan("list") + "  : Displays the table of captured secrets, tokens, and PII."},
	        {Level: 0, Text: pterm.Cyan("clear") + " : Purges the Vault from current session memory."},
	    }}.Render()

	    pterm.Bold.Println("\nSTEALTH SIGNATURE:")
	    pterm.Warning.Println("Deprecated dependency 'net/v1.0.4' (Camouflaged AES Payload)")
	case "flow":
		pterm.DefaultHeader.WithFullWidth(false).Println("USAGE: TACTICAL FLOWS")
		pterm.Println("VaporTrace mimics complex user journeys to find Business Logic flaws.")

		fmt.Println(pterm.Bold.Sprint("\nVariable Chaining (Phase 7.1):"))
		pterm.Println("Capture values using GJSON paths. Example: 'data.user.id'")
		pterm.Println("Inject them in later steps using: {{data.user.id}}")

		fmt.Println(pterm.Bold.Sprint("\nState-Machine Mapping (Phase 7.2):"))
		pterm.Println("Use 'flow step <id>' to execute a sensitive action (like /download)")
		pterm.Println("without the prerequisite steps (like /pay).")

		fmt.Println(pterm.Bold.Sprint("\nRace Condition Engine (Phase 7.3):"))
		pterm.Println("Use 'flow race <id> <threads>' to fire synchronized requests.")
		pterm.Println("Attempts to exploit TOCTOU flaws (e.g., double-spending).")

		fmt.Println(pterm.Bold.Sprint("\nCommands:"))
		pterm.BulletListPrinter{Items: []pterm.BulletListItem{
			{Level: 0, Text: "flow add   : Interactive step recording"},
			{Level: 0, Text: "flow run   : Full sequence execution"},
			{Level: 0, Text: "flow step  : Targeted out-of-order execution"},
			{Level: 0, Text: "flow race  : Synchronized high-concurrency probe"},
			{Level: 0, Text: "flow list  : View sequence and variables"},
			{Level: 0, Text: "flow clear : Clear flow and state memory"},
		}}.Render()

    // DELETED: case "usage" block (It was causing the error and is not needed here)

	case "auth":
		pterm.Println("Configures identity contexts (JWT/Cookies) for cross-account authorization testing.")
	case "sessions":
		pterm.Println("Displays the currently loaded authentication tokens for the Attacker and Victim roles.")
	case "map":
		pterm.Println("Parses Swagger/OpenAPI specs and probes for hidden shadow versions (API9).")
	case "mine":
		pterm.Println("Fuzzes discovered endpoints for hidden administrative or debug parameters.")
	case "pipeline":
		pterm.Bold.Println("COMMAND: pipeline")
		pterm.Println("DESCRIPTION:")
		pterm.Println("The Pipeline engine analyzes all endpoints stored in the Global Tactical Store")
		pterm.Println("(populated by 'map' or 'swagger' commands). It uses heuristic regex to")
		pterm.Println("categorize routes as potential BOLA, BFLA, or BOPLA targets.")
		pterm.Println("\nOPERATION:")
		pterm.BulletListPrinter{Items: []pterm.BulletListItem{
			{Level: 0, Text: "ID Detection: Finds {id} or UUIDs for BOLA testing."},
			{Level: 0, Text: "Verb Mapping: Prepares all routes for BFLA method shuffling."},
			{Level: 0, Text: "Write Detection: Flags POST/PUT routes for BOPLA/Mass-Assignment."},
		}}.Render()
		pterm.Println("\nUSAGE:")
		pterm.Cyan("pipeline")
	case "bola":
		pterm.Println("Attempts Broken Object Level Authorization (API1) by swapping identity tokens across resource IDs.")
	case "test-bola":
		pterm.Println("Runs a diagnostic BOLA sequence against httpbin to verify identity-swap and detection logic.")
	case "bopla":
		pterm.Println("Attempts Mass Assignment (API3) by injecting administrative properties into JSON payloads.")
	case "test-bopla":
		pterm.Println("Simulates a Mass Assignment attack against a reflection endpoint to verify injection logic.")
	case "bfla":
		pterm.Println("Tests for Broken Function Level Authorization by attempting administrative HTTP verbs (DELETE, POST, etc.) using a low-privilege session.")
	case "test-bfla":
		pterm.Println("Simulates Verb Tampering against httpbin to verify the method-shuffling engine.")
	case "exhaust":
		pterm.Bold.Println("DESCRIPTION:")
		pterm.Println("Tests for API4:2023 Resource Exhaustion by fuzzing pagination parameters.")
		pterm.Println("It attempts to force the server to process massive datasets by exponentially")
		pterm.Println("increasing parameters like 'limit', 'size', or 'per_page'.")
		pterm.Println("\nSTRATEGY:")
		pterm.Println("1. Identifies if the server lacks a maximum cap on requested resources.")
		pterm.Println("2. Monitors response latency to detect database/memory stress.")
		pterm.Println("\nUSAGE:")
		pterm.Cyan("exhaust <url> <parameter>")
	case "test-exhaust":
		pterm.Println("Diagnostic tool that runs the exhaustion logic against httpbin.org.")
		pterm.Println("Useful for verifying the network stack and latency timing engine.")
	case "ssrf":
		pterm.Bold.Println("DESCRIPTION:")
		pterm.Println("Tests for API7:2023 Server-Side Request Forgery by injecting")
		pterm.Println("internal and external URLs into target parameters.")
		pterm.Bold.Println("\nSTRATEGY:")
		pterm.BulletListPrinter{Items: []pterm.BulletListItem{
			{Level: 0, Text: "Injects cloud metadata IPs (169.254.169.254)"},
			{Level: 0, Text: "Injects local loopback (127.0.0.1) to find internal services"},
			{Level: 0, Text: "Uses a 'Callback' URL for Out-of-Band (OOB) detection"},
		}}.Render()
		pterm.Println("\nUSAGE:")
		pterm.Cyan("ssrf <url> <parameter> <callback_url>")
	case "test-ssrf":
		pterm.Println("Diagnostic tool that simulates an SSRF injection against httpbin.")
		pterm.Println("Verifies if the engine correctly identifies redirects and successful injections.")
	case "audit":
		pterm.Bold.Println("DESCRIPTION:")
		pterm.Println("Performs a passive and active audit of API8:2023 Security Misconfigurations.")
		pterm.Println("\nCHECKS:")
		pterm.BulletListPrinter{Items: []pterm.BulletListItem{
			{Level: 0, Text: "CORS Reflection: Checks if the API reflects arbitrary origins."},
			{Level: 0, Text: "Security Headers: Scans for missing HSTS, CSP, and X-Frame-Options."},
			{Level: 0, Text: "Verbose Errors: Attempts to trigger stack traces via invalid HTTP methods."},
		}}.Render()
	case "test-audit":
		pterm.Bold.Println("DESCRIPTION:")
		pterm.Println("A safe diagnostic command to verify the Misconfiguration Engine.")
		pterm.Println("\nOPERATION:")
		pterm.Println("1. Targets a well-known stable endpoint (google.com).")
		pterm.Println("2. Validates that CORS, Security Headers, and Method-shuffling logic are correctly identifying and reporting server responses.")
		pterm.Println("\nUSAGE:")
		pterm.Cyan("test-audit")
	case "probe":
		pterm.Bold.Println("DESCRIPTION:")
		pterm.Println("Tests for API10:2023 Unsafe Consumption of APIs.")
		pterm.Println("This module targets endpoints that process data from third-party services")
		pterm.Println("(e.g., GitHub webhooks, Stripe events, or Cloud storage callbacks).")
		pterm.Bold.Println("\nSTRATEGY:")
		pterm.BulletListPrinter{Items: []pterm.BulletListItem{
			{Level: 0, Text: "Signature Bypass: Tests if the API validates HMAC/Webhook signatures."},
			{Level: 0, Text: "Injection via Trust: Checks if unsanitized 3rd-party data enters DBs/Shells."},
			{Level: 0, Text: "SSRF via Integration: Injects malicious URLs into integration metadata fields."},
		}}.Render()
		pterm.Println("\nUSAGE:")
		pterm.Cyan("probe <url> [type]")
		pterm.Println("Valid types: generic, stripe, github")
	case "test-probe":
		pterm.Bold.Println("DESCRIPTION:")
		pterm.Println("Safe diagnostic for the Integration Probe logic.")
		pterm.Println("\nOPERATION:")
		pterm.Println("Sends mock GitHub/Stripe payloads to httpbin.org/post to verify")
		pterm.Println("header formation and payload delivery.")
		pterm.Println("\nUSAGE:")
		pterm.Cyan("test-probe")
	case "report":
		pterm.DefaultHeader.WithFullWidth(false).Println("COMMAND: report")
		pterm.Bold.Println("DESCRIPTION:")
		pterm.Println("Compiles all findings stored in the local SQLite database into a professional Markdown report.")
		pterm.Bold.Println("\nREPORT STRUCTURE:")
		pterm.BulletListPrinter{Items: []pterm.BulletListItem{
			{Level: 0, Text: "Mission Metadata (Database ID, Gen Time)"},
			{Level: 0, Text: "Phase-specific Attack Vectors"},
			{Level: 0, Text: "Vulnerability Status (VULNERABLE, EXPLOITED, MITIGATED)"},
			{Level: 0, Text: "OWASP API Top 10 Mapping"},
		}}.Render()
		pterm.Bold.Println("\nOUTPUT FILE:")
		pterm.Info.Println("VAPOR_DEBRIEF_[YYYY-MM-DD].md")
	default:
		pterm.Error.Printf("No manual entry for %s\n", cmd)
	}
}