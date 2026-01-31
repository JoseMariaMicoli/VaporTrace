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
				// Log Success via Channel (Fixes UI Leak)
				utils.TacticalLog(fmt.Sprintf("[green]Target Locked:[-] %s", args[0]))
			}
		}

	case "sessions":
		utils.TacticalLog(fmt.Sprintf("Attacker Token: %s...", shortToken(logic.CurrentSession.AttackerToken)))
		utils.TacticalLog(fmt.Sprintf("Victim Token:   %s...", shortToken(logic.CurrentSession.VictimToken)))

	// --- DISCOVERY & RECON ---
	case "swagger":
		target := logic.CurrentSession.GetTarget()
		if len(args) > 0 {
			target = args[0]
		}
		if target == "" || target == "http://localhost" {
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
		target := logic.CurrentSession.GetTarget()
		if len(args) > 0 {
			target = args[0]
		}
		if target == "" {
			utils.TacticalLog("[red]Error:[-] Usage: scrape <js_url>")
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
		target := logic.CurrentSession.GetTarget()
		endpoint := ""
		if len(args) >= 2 {
			target = args[0]
			endpoint = args[1]
		} else if len(args) == 1 && target != "" {
			endpoint = args[0]
		}

		if target == "" {
			utils.TacticalLog("[red]Error:[-] Usage: mine <url> <endpoint>")
			return
		}
		utils.TacticalLog(fmt.Sprintf("[blue]Mining hidden parameters on %s%s...[-]", target, endpoint))
		go discovery.MineParameters(target, endpoint, "")

	case "map":
		utils.TacticalLog("[blue]Starting Phase 2 Recon...[-]")
		target := logic.CurrentSession.GetTarget()
		if target != "" && target != "http://localhost" {
			go func() {
				_, err := discovery.ParseSwagger(target, "")
				if err != nil {
					utils.TacticalLog(fmt.Sprintf("[red]Swagger Error:[-] %v", err))
				}
				discovery.MineParameters(target, "", "")
			}()
		} else {
			utils.TacticalLog("[red]Error:[-] Set global target first.")
		}

	case "pipeline":
		utils.TacticalLog("[aqua]Running full tactical pipeline...[-]")
		go logic.RunPipeline(logic.CurrentSession.Threads)

	// --- LOGIC PROBES ---
	case "bola":
		target := logic.CurrentSession.GetTarget()
		if len(args) > 0 && (args[0] == "--pipeline" || args[0] == "-p") {
			utils.TacticalLog("[aqua]Starting Pipeline BOLA scan...[-]")
			go logic.ExecuteMassBOLA(logic.CurrentSession.Threads)
		} else if target != "" && target != "http://localhost" {
			ctx := &logic.BOLAContext{BaseURL: target, VictimID: "1"}
			if len(args) >= 1 {
				ctx.VictimID = args[0]
			}
			go ctx.ProbeSilent()
		} else {
			utils.TacticalLog("[red]Error:[-] No target set.")
		}

	case "bfla":
		if logic.CurrentSession.GetTarget() != "" {
			utils.TacticalLog("[aqua]Starting Mass BFLA Matrix...[-]")
			go logic.ExecuteMassBFLA(logic.CurrentSession.Threads)
		} else {
			utils.TacticalLog("[red]Error:[-] No target set.")
		}

	case "bopla":
		if logic.CurrentSession.GetTarget() != "" {
			utils.TacticalLog("[aqua]Starting Mass BOPLA Fuzzer...[-]")
			go logic.ExecuteMassBOPLA(logic.CurrentSession.Threads)
		} else {
			utils.TacticalLog("[red]Error:[-] No target set.")
		}

	case "exhaust":
		if len(args) >= 2 {
			ctx := &logic.ExhaustionContext{TargetURL: args[0], ParamName: args[1]}
			go ctx.FuzzPagination()
		} else {
			utils.TacticalLog("[red]Usage:[-] exhaust <url> <param>")
		}

	case "ssrf":
		if len(args) >= 3 {
			ctx := &logic.SSRFContext{TargetURL: args[0], ParamName: args[1], Callback: args[2]}
			go ctx.Probe()
		} else {
			utils.TacticalLog("[red]Usage:[-] ssrf <url> <param> <callback>")
		}

	case "audit":
		if len(args) > 0 {
			ctx := &logic.MisconfigContext{TargetURL: args[0]}
			go ctx.Audit()
		} else if t := logic.CurrentSession.GetTarget(); t != "" {
			ctx := &logic.MisconfigContext{TargetURL: t}
			go ctx.Audit()
		} else {
			utils.TacticalLog("[red]Usage:[-] audit <url>")
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
		db.InitDB()
		utils.TacticalLog("[green]Database Persistence Initialized.[-]")

	case "seed_db":
		utils.TacticalLog("[aqua]Injecting dummy data for report testing...[-]")
		go seedDatabase()

	case "reset_db":
		db.ResetDB()
		utils.TacticalLog("[yellow]Database Purged.[-]")

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
			utils.UpdateGlobalClient(args[0])
		}

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

// seedDatabase injects 20 dummy findings
func seedDatabase() {
	time.Sleep(500 * time.Millisecond)

	findings := []db.Finding{
		{Phase: "PHASE III: AUTH LOGIC", Target: "https://api.target.com/users/102", Details: "BOLA ID Swap: Accessed User 102", Status: "EXPLOITED", OWASP_ID: "API1:2023", MITRE_ID: "T1548", NIST_Tag: "DE.AE"},
		{Phase: "PHASE III: AUTH LOGIC", Target: "https://api.target.com/users/103", Details: "BOLA ID Swap: Accessed User 103", Status: "EXPLOITED", OWASP_ID: "API1:2023", MITRE_ID: "T1548", NIST_Tag: "DE.AE"},
		{Phase: "PHASE III: AUTH LOGIC", Target: "https://api.target.com/users/admin", Details: "BOLA ID Swap: Failed", Status: "INFO", OWASP_ID: "API1:2023", MITRE_ID: "T1548", NIST_Tag: "DE.AE"},
		{Phase: "PHASE IV: INJECTION", Target: "https://api.target.com/hooks", Details: "SSRF Internal: 169.254.169.254", Status: "CRITICAL", OWASP_ID: "API7:2023", MITRE_ID: "T1071", NIST_Tag: "DE.CM"},
		{Phase: "PHASE IV: INJECTION", Target: "https://api.target.com/callback", Details: "SSRF Internal: 127.0.0.1", Status: "CRITICAL", OWASP_ID: "API7:2023", MITRE_ID: "T1071", NIST_Tag: "DE.CM"},
		{Phase: "PHASE IV: INJECTION", Target: "https://api.target.com/img", Details: "SSRF: Open Redirect to evil.com", Status: "VULNERABLE", OWASP_ID: "API7:2023", MITRE_ID: "T1071", NIST_Tag: "DE.CM"},
		{Phase: "PHASE II: DISCOVERY", Target: "https://api.target.com/v1/swagger.json", Details: "Swagger Documentation Exposed", Status: "INFO", OWASP_ID: "API9:2023", MITRE_ID: "T1595", NIST_Tag: "ID.AM"},
		{Phase: "PHASE II: DISCOVERY", Target: "https://api.target.com/v2/api-docs", Details: "OpenAPI v3 Spec Found", Status: "INFO", OWASP_ID: "API9:2023", MITRE_ID: "T1595", NIST_Tag: "ID.AM"},
		{Phase: "PHASE II: DISCOVERY", Target: "https://api.target.com/.env", Details: "Environment File (403 Forbidden)", Status: "INFO", OWASP_ID: "API9:2023", MITRE_ID: "T1595", NIST_Tag: "ID.AM"},
		{Phase: "PHASE VIII: EXFILTRATION", Target: "https://api.target.com/debug", Details: "Leaked AWS_KEY: AKIA........", Status: "VULNERABLE", OWASP_ID: "API2:2023", MITRE_ID: "T1552", NIST_Tag: "PR.AC"},
		{Phase: "PHASE VIII: EXFILTRATION", Target: "https://api.target.com/logs", Details: "Leaked JWT Token in Body", Status: "VULNERABLE", OWASP_ID: "API2:2023", MITRE_ID: "T1552", NIST_Tag: "PR.AC"},
		{Phase: "PHASE VIII: EXFILTRATION", Target: "https://api.target.com/slack", Details: "Leaked Slack Webhook URL", Status: "VULNERABLE", OWASP_ID: "API2:2023", MITRE_ID: "T1552", NIST_Tag: "PR.AC"},
		{Phase: "PHASE III: AUTH LOGIC", Target: "https://api.target.com/admin/user", Details: "BFLA: DELETE Method Allowed", Status: "VULNERABLE", OWASP_ID: "API5:2023", MITRE_ID: "T1548.002", NIST_Tag: "DE.CM"},
		{Phase: "PHASE III: AUTH LOGIC", Target: "https://api.target.com/admin/settings", Details: "BFLA: POST Method Allowed", Status: "VULNERABLE", OWASP_ID: "API5:2023", MITRE_ID: "T1548.002", NIST_Tag: "DE.CM"},
		{Phase: "PHASE IV: INJECTION", Target: "https://api.target.com/profile", Details: "BOPLA: 'is_admin' Injected", Status: "EXPLOITED", OWASP_ID: "API3:2023", MITRE_ID: "T1538", NIST_Tag: "PR.AC"},
		{Phase: "PHASE IV: INJECTION", Target: "https://api.target.com/order", Details: "BOPLA: 'discount' Injected", Status: "EXPLOITED", OWASP_ID: "API3:2023", MITRE_ID: "T1538", NIST_Tag: "PR.AC"},
		{Phase: "PHASE 9.9: EXHAUSTION", Target: "https://api.target.com/feed", Details: "DoS: limit=1000000 (Latency: 5s)", Status: "VULNERABLE", OWASP_ID: "API4:2023", MITRE_ID: "T1499", NIST_Tag: "RS.AN"},
		{Phase: "PHASE 9.9: EXHAUSTION", Target: "https://api.target.com/search", Details: "DoS: Deep Pagination (Offset 50k)", Status: "VULNERABLE", OWASP_ID: "API4:2023", MITRE_ID: "T1499", NIST_Tag: "RS.AN"},
		{Phase: "PHASE II: DISCOVERY", Target: "https://api.target.com", Details: "Missing Header: HSTS", Status: "WEAK CONFIG", OWASP_ID: "API8:2023", MITRE_ID: "T1592", NIST_Tag: "PR.IP"},
		{Phase: "PHASE II: DISCOVERY", Target: "https://api.target.com", Details: "CORS: * (Wildcard)", Status: "WEAK CONFIG", OWASP_ID: "API8:2023", MITRE_ID: "T1592", NIST_Tag: "PR.IP"},
	}

	for _, f := range findings {
		utils.RecordFinding(f)
		time.Sleep(50 * time.Millisecond)
	}
	utils.TacticalLog("[green]Seeding Complete. 20 Findings Injected.[-]")
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
