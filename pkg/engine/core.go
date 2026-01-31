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
		utils.TacticalLog("[aqua]Injecting dummy data for report testing...[-]")
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
			if args[0] == "off" {
				utils.UpdateGlobalClient("")
			} else {
				utils.UpdateGlobalClient(args[0])
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
func seedDatabase() {
	time.Sleep(500 * time.Millisecond)

	// Base dataset mapped to your specific VaporTrace Suite
	findings := []db.Finding{
		// --- I. INFIL (Recon & Discovery) ---
		{Phase: "I. INFIL: 2.1 OpenAPI", Target: "/v1/swagger.json", Details: "Shadow API Discovery: Hidden /internal/debug identified.", Status: "INFO", OWASP_ID: "API9:2023", MITRE_ID: "T1595.002", NIST_Tag: "ID.RA", CVE_ID: "-", CVSS_Score: "0.0"},
		{Phase: "I. INFIL: 2.2 JS Mining", Target: "main.bundle.js", Details: "Hidden Route Extraction: Scraped 14 endpoints from minified source.", Status: "VULNERABLE", OWASP_ID: "API9:2023", MITRE_ID: "T1592", NIST_Tag: "ID.RA", CVE_ID: "-", CVSS_Score: "3.5"},
		{Phase: "I. INFIL: 3.1 Brute-force", Target: "/api/v0/auth", Details: "Legacy Version ID: Deprecated auth route accessible via version fuzzing.", Status: "VULNERABLE", OWASP_ID: "-", MITRE_ID: "T1589", NIST_Tag: "ID.AM", CVE_ID: "-", CVSS_Score: "5.0"},
		{Phase: "I. INFIL: 2.2 JS Mining", Target: "vendor.js", Details: "Credential Leak: Found hardcoded Stripe 'pk_test' key.", Status: "CRITICAL", OWASP_ID: "API2:2023", MITRE_ID: "T1592", NIST_Tag: "ID.RA", CVE_ID: "-", CVSS_Score: "9.1"},

		// --- II. EXPLOIT (Broken Auth & Injection) ---
		{Phase: "II. EXPLOIT: 4.1 BOLA", Target: "/api/orders/5001", Details: "Unauthorized Data Access: Accessed Order 5001 (User B) as User A.", Status: "EXPLOITED", OWASP_ID: "API1:2023", MITRE_ID: "T1548", NIST_Tag: "PR.AC", CVE_ID: "CVE-202X-BOLA", CVSS_Score: "8.8"},
		{Phase: "II. EXPLOIT: 5.1 BFLA", Target: "/api/system/reboot", Details: "Administrative Escalation: Standard user triggered restricted system action.", Status: "EXPLOITED", OWASP_ID: "API5:2023", MITRE_ID: "T1548.002", NIST_Tag: "PR.AC", CVE_ID: "CVE-202X-BFLA", CVSS_Score: "9.0"},
		{Phase: "II. EXPLOIT: 5.2 BOPLA", Target: "/api/v2/profile", Details: "Internal State Injection: Injected 'tier: platinum' via mass assignment.", Status: "EXPLOITED", OWASP_ID: "API6:2023", MITRE_ID: "T1496", NIST_Tag: "PR.DS", CVE_ID: "CVE-202X-MASS", CVSS_Score: "6.5"},
		{Phase: "II. EXPLOIT: 6.1 JWT", Target: "X-Auth-Token", Details: "Identity Spoofing: Successfully forged admin token using 'none' algorithm.", Status: "EXPLOITED", OWASP_ID: "API2:2023", MITRE_ID: "T1606", NIST_Tag: "PR.AC", CVE_ID: "CVE-202X-JWT", CVSS_Score: "9.8"},

		// --- III. EXPAND (Lateral & Infrastructure) ---
		{Phase: "III. EXPAND: 7.1 SSRF", Target: "169.254.169.254", Details: "Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service.", Status: "CRITICAL", OWASP_ID: "API7:2023", MITRE_ID: "T1046", NIST_Tag: "DE.CM", CVE_ID: "CVE-202X-SSRF", CVSS_Score: "10.0"},
		{Phase: "III. EXPAND: 8.1 DoS", Target: "/api/reports/all", Details: "Backend Service Crash: Resource exhaustion via nested JSON payload.", Status: "VULNERABLE", OWASP_ID: "API4:2023", MITRE_ID: "T1499", NIST_Tag: "RS.AN", CVE_ID: "-", CVSS_Score: "7.5"},
		{Phase: "III. EXPAND: 9.1 Persist", Target: "Mission Database", Details: "Audit Trail Integrity: Findings persisted with NIST framework tagging.", Status: "INFO", OWASP_ID: "-", MITRE_ID: "T1560", NIST_Tag: "PR.DS", CVE_ID: "-", CVSS_Score: "0.0"},

		// --- IV. OBFUSC (Stealth Ops) ---
		{Phase: "IV. OBFUSC: 11.1 Proxy", Target: "127.0.0.1:8080", Details: "Origin IP Masking: Tactical traffic successfully proxied through Burp.", Status: "ACTIVE", OWASP_ID: "-", MITRE_ID: "T1090", NIST_Tag: "PR.PT", CVE_ID: "-", CVSS_Score: "0.0"},
		{Phase: "IV. OBFUSC: 11.2 Rotation", Target: "ProxyPool-Alpha", Details: "Rate-Limit Bypass: Egress IP rotated 15 times during session.", Status: "ACTIVE", OWASP_ID: "-", MITRE_ID: "T1090.003", NIST_Tag: "PR.PT", CVE_ID: "-", CVSS_Score: "0.0"},
		{Phase: "IV. OBFUSC: 12.1 Evasion", Target: "Cloudflare WAF", Details: "WAF Signature Evasion: Randomized JA3 fingerprints and headers.", Status: "ACTIVE", OWASP_ID: "-", MITRE_ID: "T1562.001", NIST_Tag: "PR.PT", CVE_ID: "-", CVSS_Score: "0.0"},

		// --- V. COMPL (Finalization) ---
		{Phase: "V. COMPL: 13.1 Debrief", Target: "mission_logs.md", Details: "Evidence Packaging: Automated Markdown report generated.", Status: "INFO", OWASP_ID: "-", MITRE_ID: "T1020", NIST_Tag: "PR.DS", CVE_ID: "-", CVSS_Score: "0.0"},
	}

	// TRIPLE-PLUS LOOP: Generates 120+ findings with randomized variations
	targets := []string{"prod-api", "dev-cluster", "stg-nodes", "legacy-v1", "edge-gateway"}
	statusList := []string{"EXPLOITED", "VULNERABLE", "CRITICAL", "ACTIVE", "INFO"}

	utils.TacticalLog("[yellow]Initializing High-Density Data Seeding...[-]")

	for i := 0; i < 8; i++ { // 8 iterations * 15 findings = 120 records
		for _, f := range findings {
			// Randomize data to prevent identical duplicates
			f.Target = fmt.Sprintf("https://%s.target.com%s?id=%d", targets[i%5], f.Target, i*100)
			f.Status = statusList[i%5]

			utils.RecordFinding(f)
			time.Sleep(10 * time.Millisecond) // Faster seeding for large volume
		}
	}

	utils.TacticalLog("[green]Mission Environment Seeded: 120+ Findings mapped to MITRE & OWASP.[-]")
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
		"[yellow]exit[-]         | Secure Shutdown",
	}
	for _, c := range cmds {
		utils.TacticalLog(c)
	}
}

func printHelp(cmd string) {
	utils.TacticalLog(fmt.Sprintf("[aqua]MANUAL: %s[-]", cmd))
	switch cmd {
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
