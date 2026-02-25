package ui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/discovery"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/report"
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

	// Using a "Deep Sea/Cyan" theme for a more sophisticated look
	pterm.DefaultHeader.
		WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).
		WithTextStyle(pterm.NewStyle(pterm.FgBlack, pterm.Bold)).
		WithMargin(10).
		Println("VaporTrace | Surgical API Exploitation Suite")

	statusColor := pterm.FgLightGreen
	statusText := "● SYSTEM ONLINE"

	if !s.RemoteActive {
		statusColor = pterm.FgLightRed
		statusText = "○ LINK SEVERED"
	}

	// Stylized Table using a sleeker box style
	pterm.DefaultTable.WithData(pterm.TableData{
		{"UPSTREAM GATEWAY", "LOGIC ENGINE", "BUILD VERSION"},
		{"http://127.0.0.1:8080", statusColor.Sprintf(statusText), "v2.0.1-stable"},
	}).WithBoxed().Render()

	pterm.Printf("\n%s Use 'usage' for tactics or 'help' for manuals.\n\n",
		pterm.LightBlue("»"))
}

// Start launches the interactive tactical loop with Auto-Completion
func (s *Shell) Start() {
	s.RenderBanner()

	completer := readline.NewPrefixCompleter(
		readline.PcItem("init_db"),
		readline.PcItem("reset_db"),
		readline.PcItem("map"),
		readline.PcItem("mine"),
		readline.PcItem("scrape"),
		readline.PcItem("swagger"),
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
		HistoryFile:     "/tmp/vaportrace.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
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

		// FIXED: Parsing the line into command and args to pass to handleCommand
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
		s.ShowUsage()

	case "help":
		if len(args) > 0 {
			s.ShowHelp(args[0])
		} else {
			pterm.Info.Println("Usage: help <command>")
		}

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
		pterm.Info.Println("Executing Phase 2: Mapping Logic sequence...")

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
			pterm.Info.Println("Example: exhaust https://api.target.com/v1/users limit")
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
		pterm.Info.Println("Simulating SSRF against httpbin (External Redirect Test)...")
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
		if len(args) < 2 {
			pterm.Info.Println("Usage: bola <url> <victim_id>")
			return
		}
		probe := &logic.BOLAContext{
			BaseURL:  args[0],
			VictimID: args[1],
		}
		probe.Probe()

	case "bopla":
		if len(args) < 2 {
			pterm.Info.Println("Usage: bopla <url> <base_json>")
			pterm.Info.Println("Example: bopla https://api.com/v1/user '{\"name\":\"john\"}'")
			return
		}
		jsonStr := strings.Join(args[1:], " ")
		probe := &logic.BOPLAContext{
			TargetURL: args[0],
			Method:    "PATCH",
			BaseJSON:  jsonStr,
		}
		probe.Fuzz()

	case "test-bopla":
		pterm.DefaultHeader.WithFullWidth(false).Println("BOPLA Logic Test Sequence")
		pterm.Info.Println("Simulating Mass Assignment against httpbin reflection...")
		test := &logic.BOPLAContext{
			TargetURL: "https://httpbin.org/patch",
			Method:    "PATCH",
			BaseJSON:  `{"username": "vapor_user", "email": "vapor@trace.local"}`,
		}
		test.Fuzz()

	case "bfla":
		if len(args) < 1 {
			pterm.Info.Println("Usage: bfla <url>")
			return
		}
		probe := &logic.BFLAContext{TargetURL: args[0]}
		probe.Probe()

	case "test-bfla":
		pterm.Info.Println("Simulating Verb Tampering against httpbin...")
		test := &logic.BFLAContext{TargetURL: "https://httpbin.org/anything"}
		test.Probe()

	case "auth":
		if len(args) < 2 {
			pterm.Info.Println("Usage: auth <victim|attacker> <token>")
			return
		}
		if args[0] == "attacker" {
			logic.CurrentSession.AttackerToken = args[1]
			pterm.Success.Println("Attacker token updated in session store.")
		} else {
			logic.CurrentSession.VictimToken = args[1]
			pterm.Success.Println("Victim token updated in session store.")
		}

	case "sessions":
		pterm.DefaultTable.WithData(pterm.TableData{
			{"ROLE", "TOKEN SNAPSHOT"},
			{"VICTIM (User A)", logic.CurrentSession.VictimToken},
			{"ATTACKER (User B)", logic.CurrentSession.AttackerToken},
		}).WithBoxed().Render()

	case "test-bola":
		pterm.DefaultHeader.WithFullWidth(false).Println("BOLA Logic Test Sequence")
		pterm.Info.Println("TEST 1: Simulating Vulnerable Endpoint (Expect VULN)")
		vuln := &logic.BOLAContext{
			BaseURL:       "https://httpbin.org/anything",
			VictimID:      "user_777_private_data",
			AttackerToken: "evil_token_v3",
		}
		vuln.Probe()
		fmt.Println(strings.Repeat("-", 30))
		pterm.Info.Println("TEST 2: Simulating Secure Endpoint (Expect SECURE)")
		secure := &logic.BOLAContext{
			BaseURL:       "https://httpbin.org/status/403",
			VictimID:      "",
			AttackerToken: "evil_token_v3",
		}
		secure.Probe()

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
		{"swagger", "Parse OpenAPI/Swagger docs for routes", "swagger <url>"},
		{"mine", "Fuzz for hidden query parameters", "mine <url> <endpoint>"},
		{"scrape", "Extract API paths from JS files", "scrape <url>"},
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
		}
		for _, item := range helpItems {
			pterm.DefaultBulletList.WithItems([]pterm.BulletListItem{
				{Level: 0, Text: pterm.Cyan(item[0]) + ": " + item[1]},
				{Level: 1, Text: pterm.Gray("Usage: ") + item[2]},
			}).Render()
		}

	case "auth":
		pterm.Println("Configures identity contexts (JWT/Cookies) for cross-account authorization testing.")
	case "sessions":
		pterm.Println("Displays the currently loaded authentication tokens for the Attacker and Victim roles.")
	case "map":
		pterm.Println("Parses Swagger/OpenAPI specs and probes for hidden shadow versions (API9).")
	case "mine":
		pterm.Println("Fuzzes discovered endpoints for hidden administrative or debug parameters.")
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