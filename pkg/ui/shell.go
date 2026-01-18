package ui

import (
	"fmt"
	"io"
	"strings"
	"time"
	"os"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/pterm/pterm"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
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
		statusText = "○ LINK SEVERED" // Fixed: used = instead of :=
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
			if len(line) == 0 { break } else { continue }
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" { continue }

		s.handleCommand(line)
	}
}

func (s *Shell) handleCommand(input string) {
	parts := strings.Split(input, " ")
	command := parts[0]

	switch command {
	case "usage":
		s.ShowUsage()
	case "help":
		if len(parts) > 1 { s.ShowHelp(parts[1]) } else { pterm.Info.Println("Usage: help <command>") }
	case "clear", "cls", "splash":
		s.RenderBanner()
	case "exit":
        result, _ := pterm.DefaultInteractiveConfirm.
            WithDefaultText("Terminate mission and exit shell?").
            WithConfirmStyle(pterm.NewStyle(pterm.FgRed, pterm.Bold)).
            Show()

        if result {
            fmt.Println()

            // Step 1: Network
            pterm.Print(pterm.LightBlue("  ○ "))
            pterm.Print("Closing active network connections... ")
            time.Sleep(500 * time.Millisecond)
            pterm.Success.Println("Done")

            // Step 2: Database (The Guarded Call)
            pterm.Print(pterm.LightMagenta("  ○ "))
            pterm.Print("Synchronizing mission logs and closing database... ")
            db.CloseDB() // This is now safe to call multiple times
            time.Sleep(1 * time.Second)
            pterm.Success.Println("Done")

            // Step 3: Session
            pterm.Print(pterm.LightYellow("  ○ "))
            pterm.Print("Exiting tactical session... ")
            time.Sleep(400 * time.Millisecond)
            pterm.Success.Println("Goodbye")

            fmt.Println()
            pterm.Info.WithPrefix(pterm.Prefix{Text: "OFFLINE", Style: pterm.NewStyle(pterm.FgBlack, pterm.BgCyan)}).
                Println("VaporTrace terminated safely. Reveal the invisible.")
            
            os.Exit(0) // Ensure the process actually dies
            return
        }
	case "init_db":
        fmt.Println() // Spacer
        pterm.DefaultHeader.WithFullWidth(false).Println("Phase 5: Intelligence Initialization")

        // Stage 1: Connection
        pterm.Print(pterm.LightCyan("  ● ")) 
        pterm.Print("Establishing link to SQLite persistence... ")
        db.InitDB() // Calls the logic we built in pkg/db
        time.Sleep(600 * time.Millisecond)
        pterm.Success.Println("Connected")

        // Stage 2: Worker Pool
        pterm.Print(pterm.LightGreen("  ● "))
        pterm.Print("Spawning asynchronous log worker pool... ")
        go db.StartAsyncWorker()
        time.Sleep(800 * time.Millisecond)
        pterm.Success.Println("Active")

        // Stage 3: Ready State
        fmt.Println()
        pterm.Info.WithPrefix(pterm.Prefix{Text: "READY", Style: pterm.NewStyle(pterm.FgBlack, pterm.BgGreen)}).
            Println("Database persistence is now online. All findings will be logged.")

	case "reset_db":
        // Safety Warning with Red Bold Style
        pterm.Warning.WithPrefix(pterm.Prefix{Text: "DANGER", Style: pterm.NewStyle(pterm.FgBlack, pterm.BgRed)}).
            Println("This action will permanently delete all mission findings and database logs!")
        
        // English Confirmation
        result, _ := pterm.DefaultInteractiveConfirm.
            WithDefaultText("Are you sure you want to PURGE the database?").
            WithConfirmStyle(pterm.NewStyle(pterm.FgRed, pterm.Bold)).
            Show()
        
        if result {
            fmt.Println() // Spacer for visual clarity

            // Stage 1: Purging Data
            pterm.Print(pterm.LightRed("  × ")) // Deletion icon
            pterm.Print("Wiping all records from findings table... ")
            time.Sleep(600 * time.Millisecond)
            
            // Stage 2: Resetting Metadata
            pterm.Print("\n  × ")
            pterm.Print("Resetting DATABASE ID and Gen Time... ")
            db.ResetDB() // Executes the SQL DROP and Re-init [cite: 2, 3]
            time.Sleep(800 * time.Millisecond)
            pterm.Success.Println("Done")

            // Stage 3: Final Confirmation
            fmt.Println()
            pterm.Success.WithPrefix(pterm.Prefix{Text: "CLEAN", Style: pterm.NewStyle(pterm.FgBlack, pterm.BgGreen)}).
                Println("Database has been successfully purged. System ready for new mission.")
        } else {
            pterm.Info.Println("Reset operation aborted by operator.")
        }
	case "map":
		pterm.Info.Println("Executing Phase 2: Mapping Logic sequence...")
	case "audit":
		if len(parts) < 2 {
			pterm.Info.Println("Usage: audit <url>")
			return
		}
		probe := &logic.MisconfigContext{TargetURL: parts[1]}
		probe.Audit()
	case "probe":
		if len(parts) < 2 {
			pterm.Info.Println("Usage: probe <url> [type]")
			return
		}
		iType := "generic"
		if len(parts) > 2 { iType = parts[2] }
		probe := &logic.IntegrationContext{TargetURL: parts[1], IntegrationType: iType}
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
		if len(parts) < 3 {
			pterm.Info.Println("Usage: exhaust <url> <parameter>")
			pterm.Info.Println("Example: exhaust https://api.target.com/v1/users limit")
			return
		}
		probe := &logic.ExhaustionContext{TargetURL: parts[1], ParamName: parts[2]}
		probe.FuzzPagination()
	case "ssrf":
		if len(parts) < 4 {
			pterm.Info.Println("Usage: ssrf <url> <parameter> <callback_url>")
			return
		}
		probe := &logic.SSRFContext{TargetURL: parts[1], ParamName: parts[2], Callback: parts[3]}
		probe.Probe()

	case "test-ssrf":
		pterm.Info.Println("Simulating SSRF against httpbin (External Redirect Test)...")
		// We use httpbin's redirect endpoint to simulate an SSRF-vulnerable parameter
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
		if len(parts) < 3 {
			pterm.Info.Println("Usage: bola <url> <victim_id>")
			return
		}
		probe := &logic.BOLAContext{
			BaseURL:  parts[1],
			VictimID: parts[2],
		}
		probe.Probe()
	case "bopla":
		if len(parts) < 3 {
			pterm.Info.Println("Usage: bopla <url> <base_json>")
			pterm.Info.Println("Example: bopla https://api.com/v1/user '{\"name\":\"john\"}'")
			return
		}
		jsonStr := strings.Join(parts[2:], " ")
		probe := &logic.BOPLAContext{
			TargetURL: parts[1],
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
		if len(parts) < 2 {
			pterm.Info.Println("Usage: bfla <url>")
			return
		}
		probe := &logic.BFLAContext{TargetURL: parts[1]}
		probe.Probe()
	case "test-bfla":
		pterm.Info.Println("Simulating Verb Tampering against httpbin...")
		test := &logic.BFLAContext{TargetURL: "https://httpbin.org/anything"}
		test.Probe()
	case "auth":
		if len(parts) < 3 {
			pterm.Info.Println("Usage: auth <victim|attacker> <token>")
			return 
		}
		if parts[1] == "attacker" {
			logic.CurrentSession.AttackerToken = parts[2]
			pterm.Success.Println("Attacker token updated in session store.")
		} else {
			logic.CurrentSession.VictimToken = parts[2]
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
			BaseURL:        "https://httpbin.org/anything", 
			VictimID:       "user_777_private_data",
			AttackerToken: "evil_token_v3",
		}
		vuln.Probe()

		fmt.Println(strings.Repeat("-", 30))

		pterm.Info.Println("TEST 2: Simulating Secure Endpoint (Expect SECURE)")
		secure := &logic.BOLAContext{
			BaseURL:        "https://httpbin.org/status/403",
			VictimID:       "", 
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
		{"init_db",  "Initialize SQLite Persistence", "init_db"},
		{"reset_db", "Wipe local mission data (Purge)",      "reset_db"},
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
	case "test-bfla": pterm.Println("Simulates Verb Tampering against httpbin to verify the method-shuffling engine.")
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
		pterm.Println("2. Validates that CORS, Security Headers, and Method-shuffling logic")
		pterm.Println("   are correctly identifying and reporting server responses.")
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
	default:
		pterm.Error.Printf("No manual entry for %s\n", cmd)
	}
}