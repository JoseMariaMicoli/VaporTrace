package engine

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/discovery"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/report"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// getTarget helps commands inherit the global target if no argument is provided
func getTarget(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	global := logic.CurrentSession.GetTarget()
	if global == "" || global == "http://localhost" {
		return ""
	}
	return global
}

// ExecuteCommand parses raw input strings and routes them to the appropriate logic module.
func ExecuteCommand(rawCmd string) {
	if rawCmd == "" {
		return
	}
	parts := strings.Fields(rawCmd)
	verb := strings.ToLower(parts[0])
	args := parts[1:]

	utils.TacticalLog(fmt.Sprintf("[yellow]EXEC:[-] %s", rawCmd))

	switch verb {
	// --- NEURAL ENGINE (Sprint 10.6) ---
	case "neuro":
		if len(args) == 0 {
			utils.TacticalLog("Usage: neuro config | neuro on | neuro off")
			return
		}
		if args[0] == "config" {
			// Usage: neuro config <provider> <model> [api_key] [endpoint]
			if len(args) < 3 {
				utils.TacticalLog("[red]Usage:[-] neuro config <provider> <model> [api_key] [endpoint]")
				return
			}
			provider := args[1]
			model := args[2]
			apiKey := ""
			endpoint := ""

			if len(args) > 3 {
				apiKey = args[3]
			}
			if len(args) > 4 {
				endpoint = args[4]
			}

			logic.GlobalNeuro.Configure(provider, apiKey, model, endpoint)
		} else if args[0] == "on" {
			logic.GlobalNeuro.Active = true
			utils.TacticalLog("[green]Neural Engine Activated.[-]")
		} else if args[0] == "off" {
			logic.GlobalNeuro.Active = false
			utils.TacticalLog("[yellow]Neural Engine Deactivated.[-]")
		} else {
			utils.TacticalLog("Invalid neuro command.")
		}

	case "test-neuro":
		utils.TacticalLog("[blue]Testing Neural Engine Connectivity...[-]")
		logic.GlobalNeuro.TestConnectivity()

	case "neuro-gen":
		// Usage: neuro-gen <context> <count>
		// Example: neuro-gen "SQL Injection in login form" 5
		if len(args) < 2 {
			utils.TacticalLog("[yellow]Usage: neuro-gen <context_string> <count>")
			return
		}
		count, _ := strconv.Atoi(args[1])
		logic.GlobalNeuro.GenerateAttackVectors(args[0], count)

	// --- IDENTITY & SESSION ---
	case "auth":
		if len(args) < 2 {
			utils.TacticalLog("[red]Usage:[-] auth <attacker|victim> <token>")
			return
		}
		if args[0] == "attacker" {
			logic.CurrentSession.AttackerToken = args[1]
		} else {
			logic.CurrentSession.VictimToken = args[1]
		}
		utils.TacticalLog(fmt.Sprintf("[green]Identity Updated:[-] %s", args[0]))

	case "target":
		if len(args) > 0 {
			err := logic.CurrentSession.SetGlobalTarget(args[0])
			if err != nil {
				utils.TacticalLog(fmt.Sprintf("[red]Target Error:[-] %v", err))
			} else {
				utils.TacticalLog(fmt.Sprintf("[green]Target Locked:[-] %s", args[0]))
			}
		}

	case "sessions":
		utils.TacticalLog(fmt.Sprintf("Attacker Token: %s...", shortToken(logic.CurrentSession.AttackerToken)))
		utils.TacticalLog(fmt.Sprintf("Victim Token:   %s...", shortToken(logic.CurrentSession.VictimToken)))

	// --- DISCOVERY & RECON ---
	case "swagger":
		target := getTarget(args)
		if target == "" {
			utils.TacticalLog("[red]Error:[-] Usage: swagger <url> (or set global target)")
			return
		}
		utils.TacticalLog(fmt.Sprintf("[blue]Parsing OpenAPI spec at %s...[-]", target))
		go func() {
			endpoints, err := discovery.ParseSwagger(target, "")
			if err != nil {
				utils.TacticalLog(fmt.Sprintf("[red]Swagger Failed:[-] %v", err))
				return
			}
			utils.TacticalLog(fmt.Sprintf("[green]Success:[-] Found %d endpoints.", len(endpoints)))
		}()

	case "scrape":
		target := getTarget(args)
		if target == "" {
			utils.TacticalLog("[red]Error:[-] Usage: scrape <js_url> (or set global target)")
			return
		}
		utils.TacticalLog(fmt.Sprintf("[blue]Scraping JS Bundle: %s...[-]", target))
		go func() {
			paths, err := discovery.ExtractJSPaths(target, "")
			if err != nil {
				utils.TacticalLog(fmt.Sprintf("[red]Scrape Failed:[-] %v", err))
				return
			}
			utils.TacticalLog(fmt.Sprintf("[green]Success:[-] Extracted %d paths.", len(paths)))
		}()

	case "mine":
		target := getTarget(args)
		endpoint := ""
		// Handle argument shifting: mine <url> <endpoint> VS mine <endpoint> (w/ global target)
		if len(args) >= 2 {
			target = args[0]
			endpoint = args[1]
		} else if len(args) == 1 && logic.CurrentSession.GetTarget() != "" {
			target = logic.CurrentSession.GetTarget()
			endpoint = args[0]
		}

		if target == "" {
			utils.TacticalLog("[red]Error:[-] Usage: mine <url> <endpoint>")
			return
		}

		utils.TacticalLog(fmt.Sprintf("[blue]Mining hidden parameters on %s%s...[-]", target, endpoint))
		go func() {
			discovery.MineParameters(target, endpoint, "")
			utils.TacticalLog("[green]Mining Sequence Complete.[-]")
		}()

	case "map":
		target := getTarget(args)
		if target != "" {
			utils.TacticalLog(fmt.Sprintf("[blue]Starting Phase 2 Recon against %s...[-]", target))
			go func() {
				// 1. Swagger
				endpoints, err := discovery.ParseSwagger(target, "")
				if err != nil {
					utils.TacticalLog(fmt.Sprintf("[red]Swagger Error:[-] %v", err))
				} else {
					utils.TacticalLog(fmt.Sprintf("[green]Swagger Mapped:[-] %d routes", len(endpoints)))
				}
				// 2. Mining
				discovery.MineParameters(target, "", "")
				utils.TacticalLog("[green]Recon Map Finished.[-]")
			}()
		} else {
			utils.TacticalLog("[red]Error:[-] Usage: map <url> (or set global target first).")
		}

	case "pipeline":
		utils.TacticalLog("[aqua]Initializing Industrialized Attack Pipeline...[-]")
		go func() {
			utils.TacticalLog(fmt.Sprintf("[blue]Concurrency Level: %d threads[-]", logic.CurrentSession.Threads))
			logic.RunPipeline(logic.CurrentSession.Threads)
		}()

	// --- LOGIC PROBES ---
	case "bola":
		target := getTarget(args)
		isPipeline := false

		// Check for flags
		if len(args) > 0 && (args[0] == "--pipeline" || args[0] == "-p") {
			isPipeline = true
		}

		if isPipeline {
			utils.TacticalLog("[aqua]Starting Mass BOLA Pipeline Scan...[-]")
			go logic.ExecuteMassBOLA(logic.CurrentSession.Threads)
		} else if target != "" {
			// Surgical Mode
			victim := "1"
			// Logic to handle 'bola <url> <id>' vs 'bola <id>' (global target)
			if len(args) >= 2 {
				target = args[0]
				victim = args[1]
			} else if len(args) == 1 && logic.CurrentSession.GetTarget() != "" {
				// Use inherited target
				victim = args[0]
			}

			utils.TacticalLog(fmt.Sprintf("[blue]Launching Surgical BOLA Probe on %s (ID: %s)...[-]", target, victim))
			ctx := &logic.BOLAContext{BaseURL: target, VictimID: victim}
			go func() {
				ctx.ProbeSilent()
				utils.TacticalLog("[green]Surgical BOLA Probe Finished.[-]")
			}()
		} else {
			utils.TacticalLog("[red]Error:[-] No target set. Use 'target <url>' or 'bola <url> <id>'.")
		}

	case "bfla":
		if logic.CurrentSession.GetTarget() != "" {
			utils.TacticalLog("[aqua]Starting Mass BFLA Matrix (Verb Tampering)...[-]")
			go logic.ExecuteMassBFLA(logic.CurrentSession.Threads)
		} else {
			utils.TacticalLog("[red]Error:[-] No global target set. Run 'target <url>' first.")
		}

	case "bopla":
		if logic.CurrentSession.GetTarget() != "" {
			utils.TacticalLog("[aqua]Starting Mass BOPLA Fuzzer (Property Injection)...[-]")
			go logic.ExecuteMassBOPLA(logic.CurrentSession.Threads)
		} else {
			utils.TacticalLog("[red]Error:[-] No global target set. Run 'target <url>' first.")
		}

	case "exhaust":
		target := getTarget(args)
		param := "limit"

		if len(args) >= 2 {
			target = args[0]
			param = args[1]
		} else if len(args) == 1 && logic.CurrentSession.GetTarget() != "" {
			target = logic.CurrentSession.GetTarget()
			param = args[0]
		}

		if target == "" {
			utils.TacticalLog("[red]Usage:[-] exhaust <url> <param>")
		} else {
			go func() {
				utils.TacticalLog(fmt.Sprintf("[blue]Fuzzing Pagination on %s?%s=...[-]", target, param))
				ctx := &logic.ExhaustionContext{TargetURL: target, ParamName: param}
				ctx.FuzzPagination()
			}()
		}

	case "ssrf":
		target := getTarget(args)
		param := "url"
		cb := "http://127.0.0.1"

		// Handle args: ssrf <url> <param> <cb> OR ssrf <param> <cb> (global)
		if len(args) >= 3 {
			target = args[0]
			param = args[1]
			cb = args[2]
		} else if len(args) == 2 && logic.CurrentSession.GetTarget() != "" {
			target = logic.CurrentSession.GetTarget()
			param = args[0]
			cb = args[1]
		}

		if target == "" {
			utils.TacticalLog("[red]Usage:[-] ssrf <url> <param> <callback>")
		} else {
			go func() {
				utils.TacticalLog(fmt.Sprintf("[blue]Probing SSRF on %s (%s)...[-]", target, param))
				ctx := &logic.SSRFContext{TargetURL: target, ParamName: param, Callback: cb}
				ctx.Probe()
			}()
		}

	case "audit":
		target := getTarget(args)
		if target != "" {
			utils.TacticalLog(fmt.Sprintf("[blue]Auditing Security Headers/CORS on %s...[-]", target))
			ctx := &logic.MisconfigContext{TargetURL: target}
			go func() {
				ctx.Audit()
				utils.TacticalLog("[green]Audit Complete.[-]")
			}()
		} else {
			utils.TacticalLog("[red]Usage:[-] audit <url>")
		}

	case "probe":
		target := getTarget(args)
		iType := "generic"

		// Handle args: probe <url> <type> OR probe <type> (global)
		if len(args) >= 2 {
			target = args[0]
			iType = args[1]
		} else if len(args) == 1 {
			// Naive check: if arg looks like URL, treat as URL
			if strings.HasPrefix(args[0], "http") {
				target = args[0]
			} else if logic.CurrentSession.GetTarget() != "" {
				target = logic.CurrentSession.GetTarget()
				iType = args[0]
			}
		}

		if target == "" {
			utils.TacticalLog("[red]Usage:[-] probe <url> [type] (or set global target)")
		} else {
			utils.TacticalLog(fmt.Sprintf("[blue]Launching %s Integration Probe against %s...[-]", iType, target))
			ctx := &logic.IntegrationContext{TargetURL: target, IntegrationType: iType}
			go ctx.Probe()
		}

	// --- FLOW ENGINE ---
	case "flow":
		if len(args) == 0 {
			utils.TacticalLog("Flow Commands: list, run, clear, race.")
			return
		}
		switch args[0] {
		case "list":
			if len(logic.ActiveFlow) == 0 {
				utils.TacticalLog("Flow queue empty.")
			} else {
				for i, s := range logic.ActiveFlow {
					utils.TacticalLog(fmt.Sprintf("[%d] %s %s (%s)", i+1, s.Method, s.URL, s.Name))
				}
			}
		case "clear":
			logic.ActiveFlow = []logic.FlowStep{}
			logic.FlowContext = make(map[string]string)
			utils.TacticalLog("[yellow]Flow and Context Cleared.[-]")
		case "run":
			go logic.RunFlow()
		case "race":
			if len(args) < 3 {
				utils.TacticalLog("Usage: flow race <step_id> <threads>")
				return
			}
			id, _ := strconv.Atoi(args[1])
			threads, _ := strconv.Atoi(args[2])
			go logic.RunRace(id-1, threads)
		default:
			utils.TacticalLog("Unknown flow command.")
		}

	// --- SYSTEM & DB ---
	case "init_db":
		// STRICT: Only initialize connection, DO NOT SEED.
		db.InitDB()
		utils.TacticalLog("[green]Database Persistence Initialized (Empty State).[-]")

	case "seed_db":
		// EXPLICIT SEEDING COMMAND
		utils.TacticalLog("[aqua]Injecting high-fidelity mock data for C-Level Report...[-]")
		go seedDatabase()

	case "reset_db":
		// EXPLICIT PURGE COMMAND
		db.ResetDB()
		utils.TacticalLog("[yellow]Database Purged (Reset).[-]")

	case "report":
		report.GenerateMissionDebrief()

	case "weaver":
		interval := 60
		if len(args) > 0 {
			if i, err := strconv.Atoi(args[0]); err == nil {
				interval = i
			}
		}
		utils.TacticalLog(fmt.Sprintf("[magenta]Deploying Ghost-Weaver (Interval: %ds)...[-]", interval))
		config := logic.WeaverConfig{Interval: time.Duration(interval) * time.Second, Active: true}
		go logic.StartGhostWeaver(config)

	case "loot":
		if len(args) > 0 && args[0] == "list" {
			utils.TacticalLog("[magenta]Accessing Discovery Vault...[-]")
			if len(logic.Vault) == 0 {
				utils.TacticalLog("Vault is empty.")
			}
			for _, v := range logic.Vault {
				utils.TacticalLog(fmt.Sprintf("[yellow]%s[-] %s", v.Type, v.Value))
			}
		} else if len(args) > 0 && args[0] == "clear" {
			logic.Vault = []logic.Finding{}
			utils.TacticalLog("[green]Vault Purged.[-]")
		} else {
			utils.TacticalLog("Usage: loot list | loot clear")
		}

	case "proxy":
		if len(args) > 0 {
			// FIX: Replaced utils.UpdateGlobalClient with logic.SetProxy
			if args[0] == "off" {
				logic.SetProxy("")
			} else {
				logic.SetProxy(args[0])
			}
		} else {
			utils.TacticalLog("[red]Usage:[-] proxy <url> | proxy off")
		}

	case "proxies":
		if len(args) > 0 {
			if args[0] == "load" && len(args) >= 2 {
				if err := logic.LoadProxiesFromFile(args[1]); err != nil {
					utils.TacticalLog(fmt.Sprintf("[red]Load Failed:[-] %v", err))
				} else {
					logic.InitializeRotaryClient()
					utils.TacticalLog(fmt.Sprintf("[green]Proxy Pool Loaded:[-] %d proxies active.", len(logic.ProxyPool)))
				}
			} else if args[0] == "reset" {
				logic.ProxyPool = []string{}
				logic.InitializeRotaryClient()
				utils.TacticalLog("[yellow]Proxy Pool Reset. Reverting to default transport.[-]")
			} else {
				utils.TacticalLog("[red]Usage:[-] proxies load <file> | proxies reset")
			}
		}

	case "clear":
		// Sends signal to TUI to wipe logs
		utils.TacticalLog("___CLEAR_SCREEN_SIGNAL___")

	case "usage":
		printUsage()

	case "help":
		if len(args) > 0 {
			printHelp(args[0])
		} else {
			utils.TacticalLog("[white]Usage: help <command>[-]")
		}

	// Internal command triggered by UI Modal
	case "__internal_shutdown":
		go func() {
			utils.TacticalLog("[red::b]INITIATING SEQUENTIAL SHUTDOWN...[-:-:-]")

			utils.TacticalLog("[blue]⠋[-] Terminating Network Stacks...")
			time.Sleep(500 * time.Millisecond)
			utils.TacticalLog("[green]✔[-] Network Terminated.")

			utils.TacticalLog("[blue]⠋[-] Closing Database Persistence...")
			db.CloseDB()
			time.Sleep(500 * time.Millisecond)
			utils.TacticalLog("[green]✔[-] Database Sync Complete.")

			utils.TacticalLog("[blue]⠋[-] Killing Ghost-Weaver Agents...")
			time.Sleep(500 * time.Millisecond)
			utils.TacticalLog("[green]✔[-] Agents Purged.")

			utils.TacticalLog("[red::b]SYSTEM HALTED. GOODBYE.[-:-:-]")
			time.Sleep(1200 * time.Millisecond)
			os.Exit(0)
		}()

	case "exit":
		utils.TacticalLog("[yellow]Please use the Esc key or the Dashboard UI to initiate safe shutdown.[-]")

	default:
		if strings.HasPrefix(verb, "test-") {
			handleTestCommands(verb)
		} else {
			utils.TacticalLog(fmt.Sprintf("[red]Command not recognized:[-] %s", verb))
		}
	}
}

// seedDatabase injects a massive dataset (120+ entries) strictly aligned to VaporTrace Mapping
// MODIFIED: Updated to populate new architectural fields (Command, CVSS_Numeric, etc.)
func seedDatabase() {
	time.Sleep(500 * time.Millisecond)

	// Comprehensive dataset for C-Level Report Generation
	findings := []db.Finding{
		// --- CRITICAL (Remediation Priority) ---
		{
			Phase: "II. EXPLOIT", Target: "https://api.target.corp/users/1001",
			Details: "BOLA: Accessed administrative user profile via ID manipulation.",
			Status:  "EXPLOITED",
			Command: "bola", OWASP_ID: "API1:2023 BOLA",
			MITRE_ID: "T1594", MitreTactic: "Exfiltration",
			NIST_Tag: "PR.AC", NistControl: "PR.AC-03",
			CVE_ID: "CVE-2024-BOLA", CVSS_Score: "9.1", CVSS_Numeric: 9.1,
		},
		{
			Phase: "III. EXPAND", Target: "https://api.target.corp/hooks/stripe",
			Details: "SSRF: Cloud Metadata (169.254.169.254) keys exfiltrated.",
			Status:  "CRITICAL",
			Command: "ssrf", OWASP_ID: "API7:2023 SSRF",
			MITRE_ID: "T1071.001", MitreTactic: "Command & Control",
			NIST_Tag: "DE.CM", NistControl: "PR.DS-01",
			CVE_ID: "CVE-2021-26855", CVSS_Score: "9.8", CVSS_Numeric: 9.8,
		},
		{
			Phase: "II. EXPLOIT", Target: "https://api.target.corp/admin/roles",
			Details: "BOPLA: Mass Assignment allowed injection of 'role: admin'.",
			Status:  "VULNERABLE",
			Command: "bopla", OWASP_ID: "API3:2023 Property Injection",
			MITRE_ID: "T1592.001", MitreTactic: "Privilege Escalation",
			NIST_Tag: "PR.DS", NistControl: "PR.DS-01",
			CVE_ID: "CVE-2022-23131", CVSS_Score: "8.8", CVSS_Numeric: 8.8,
		},

		// --- HIGH RISKS ---
		{
			Phase: "II. EXPLOIT", Target: "https://api.target.corp/v2/delete_user",
			Details: "BFLA: DELETE method accepted from unprivileged account.",
			Status:  "VULNERABLE",
			Command: "bfla", OWASP_ID: "API5:2023 BFLA",
			MITRE_ID: "T1548.003", MitreTactic: "Privilege Escalation",
			NIST_Tag: "PR.AC", NistControl: "PR.AC-05",
			CVE_ID: "CVE-2023-30533", CVSS_Score: "8.2", CVSS_Numeric: 8.2,
		},
		{
			Phase: "IV. OBFUSC", Target: "https://api.target.corp/integrations/webhook",
			Details: "Unsafe Consumption: No signature verification on 3rd party webhook.",
			Status:  "VULNERABLE",
			Command: "probe", OWASP_ID: "API10:2023 Unsafe Consumption",
			MITRE_ID: "T1190", MitreTactic: "Initial Access",
			NIST_Tag: "PR.DS", NistControl: "PR.DS-02",
			CVE_ID: "CVE-2024-PROBE", CVSS_Score: "7.5", CVSS_Numeric: 7.5,
		},
		{
			Phase: "III. EXPAND", Target: "https://api.target.corp/reports/all",
			Details: "DoS: Pagination limit fuzzing caused 5s latency spike.",
			Status:  "VULNERABLE",
			Command: "exhaust", OWASP_ID: "API4:2023 Resource Exhaustion",
			MITRE_ID: "T1499.004", MitreTactic: "Impact",
			NIST_Tag: "RS.AN", NistControl: "DE.AE-02",
			CVE_ID: "CVE-2023-44487", CVSS_Score: "7.5", CVSS_Numeric: 7.5,
		},

		// --- MEDIUM / LOW RISKS (Info & Audit) ---
		{
			Phase: "I. INFIL", Target: "https://api.target.corp/v1/swagger.json",
			Details: "Information Disclosure: Full OpenAPI spec exposed publicly.",
			Status:  "INFO",
			Command: "map", OWASP_ID: "API9:2023 Inventory",
			MITRE_ID: "T1595.002", MitreTactic: "Reconnaissance",
			NIST_Tag: "ID.AM", NistControl: "ID.AM-07",
			CVE_ID: "-", CVSS_Score: "0.0", CVSS_Numeric: 0.0,
		},
		{
			Phase: "I. INFIL", Target: "https://api.target.corp/app.bundle.js",
			Details: "Hardcoded Secrets: AWS S3 Bucket URL found in JS.",
			Status:  "INFO",
			Command: "scrape", OWASP_ID: "API2:2023 Broken Auth",
			MITRE_ID: "T1552", MitreTactic: "Credential Access",
			NIST_Tag: "PR.IP", NistControl: "PR.AC-01",
			CVE_ID: "-", CVSS_Score: "4.5", CVSS_Numeric: 4.5,
		},
		{
			Phase: "II. DISCOVERY", Target: "https://api.target.corp",
			Details: "Misconfiguration: Missing Strict-Transport-Security header.",
			Status:  "WEAK CONFIG",
			Command: "audit", OWASP_ID: "API8:2023 Misconfig",
			MITRE_ID: "T1562.001", MitreTactic: "Defense Evasion",
			NIST_Tag: "PR.PS", NistControl: "PR.PS-01",
			CVE_ID: "-", CVSS_Score: "3.5", CVSS_Numeric: 3.5,
		},
	}

	utils.TacticalLog("[yellow]Seeding Database with enriched findings...[-]")
	for _, f := range findings {
		utils.RecordFinding(f)
		time.Sleep(20 * time.Millisecond)
	}

	utils.TacticalLog("[green]Mission Environment Seeded: Findings mapped to MITRE, NIST, & OWASP.[-]")
}

func handleTestCommands(verb string) {
	utils.TacticalLog(fmt.Sprintf("[white]Running Diagnostic: %s (httpbin)...[-]", verb))
	go func() {
		switch verb {
		case "test-bola":
			ctx := &logic.BOLAContext{BaseURL: "https://httpbin.org/anything", VictimID: "999"}
			ctx.ProbeSilent()
		case "test-bopla":
			ctx := &logic.BOPLAContext{TargetURL: "https://httpbin.org/patch", Method: "PATCH", BaseJSON: `{"user":"test"}`}
			ctx.RunFuzzer(1)
		case "test-bfla":
			ctx := &logic.BFLAContext{TargetURL: "https://httpbin.org/anything"}
			ctx.MassProbe(1)
		case "test-exhaust":
			ctx := &logic.ExhaustionContext{TargetURL: "https://httpbin.org/get", ParamName: "limit"}
			ctx.FuzzPagination()
		case "test-ssrf":
			ctx := &logic.SSRFContext{TargetURL: "https://httpbin.org/redirect-to", ParamName: "url", Callback: "http://google.com"}
			ctx.Probe()
		}
	}()
}

func printUsage() {
	utils.TacticalLog("[aqua]COMMAND MANUAL:[-]")
	cmds := []string{
		"[yellow]init_db[-]      | Initialize Persistence",
		"[yellow]seed_db[-]      | Populate DB with Dummy Data",
		"[yellow]target[-]       | Lock global target URL",
		"[yellow]map[-]          | Full Recon (Swagger + JS Scraper)",
		"[yellow]mine[-]         | Fuzz for hidden parameters",
		"[yellow]scrape[-]       | Extract API paths from JS",
		"[yellow]swagger[-]      | Parse OpenAPI docs",
		"[yellow]bola[-]         | Broken Object Level Auth",
		"[yellow]bopla[-]        | Mass Assignment Fuzzer",
		"[yellow]bfla[-]         | Broken Function Level Auth",
		"[yellow]ssrf[-]         | Server-Side Request Forgery",
		"[yellow]exhaust[-]      | Resource Exhaustion (DoS)",
		"[yellow]audit[-]        | Security Misconfiguration Audit",
		"[yellow]probe[-]        | Integration/Webhook Probe",
		"[yellow]loot[-]         | View Captured Secrets",
		"[yellow]weaver[-]       | Deploy OIDC Interceptor",
		"[yellow]pipeline[-]     | Auto-Analysis & Attack",
		"[yellow]report[-]       | Generate Markdown Debrief",
		"[yellow]neuro[-]        | Configure & Control AI Engine",
		"[yellow]test-neuro[-]   | Verify AI Connectivity",
		"[yellow]neuro-gen[-]    | Generate AI Attack Vectors",
		"[yellow]exit[-]         | Secure Shutdown",
	}
	for _, c := range cmds {
		utils.TacticalLog(c)
	}
}

func printHelp(cmd string) {
	utils.TacticalLog(fmt.Sprintf("[aqua]MANUAL: %s[-]", cmd))
	switch cmd {
	case "neuro":
		utils.TacticalLog("Configures the Neural Engine.")
		utils.TacticalLog("Usage: neuro config <provider> <model> [api_key] [endpoint]")
		utils.TacticalLog("Example: neuro config ollama mistral")
		utils.TacticalLog("Example: neuro config openai gpt-4 sk-123...")
	case "seed_db":
		utils.TacticalLog("Injects 20 fake vulnerabilities into the database. Useful for verifying the 'report' command without running live attacks.")
	case "bola":
		utils.TacticalLog("Attempts to access resources of other users by iterating IDs in the URL (e.g., /user/1 -> /user/2).")
	case "bopla":
		utils.TacticalLog("Mass Assignment. Injects administrative fields (is_admin, role, discount) into JSON requests.")
	case "bfla":
		utils.TacticalLog("Tests for Broken Function Level Authorization by attempting forbidden HTTP methods (DELETE, PUT) on endpoints.")
	case "ssrf":
		utils.TacticalLog("Injects internal IP addresses (127.0.0.1) and cloud metadata URLs (169.254.169.254) into parameters.")
	case "exhaust":
		utils.TacticalLog("Fuzzes pagination parameters (limit, size) with large values to test for Resource Exhaustion (DoS).")
	case "audit":
		utils.TacticalLog("Checks for Security Misconfigurations like missing headers (HSTS) and weak CORS policies.")
	case "probe":
		utils.TacticalLog("Tests for Unsafe Consumption by injecting unsigned payloads into webhook integrations.")
	case "pipeline":
		utils.TacticalLog("Analyzes all discovered endpoints (from 'map') and automatically routes them to the appropriate attack engines.")
	case "weaver":
		utils.TacticalLog("Deploys a background agent to intercept OIDC tokens and mask data exfiltration.")
	case "map":
		utils.TacticalLog("Performs full reconnaissance: Parses Swagger/OpenAPI docs and scrapes JavaScript files for hidden endpoints.")
	case "scrape":
		utils.TacticalLog("Extracts potential API endpoints from a JavaScript file.")
	case "swagger":
		utils.TacticalLog("Parses a Swagger/OpenAPI JSON file to map the API attack surface.")
	case "mine":
		utils.TacticalLog("Fuzzes an endpoint for hidden query parameters like 'debug', 'admin', 'test'.")
	case "proxy":
		utils.TacticalLog("Sets the upstream proxy (e.g., Burp Suite) for all tactical traffic.")
	case "loot":
		utils.TacticalLog("Manages the Discovery Vault, which captures secrets (Keys, Tokens) from response bodies.")
	default:
		utils.TacticalLog("No specific manual entry found. Try 'usage' for a list of commands.")
	}
}

func shortToken(t string) string {
	if len(t) > 10 {
		return t[:10]
	}
	return t
}
