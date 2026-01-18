package ui

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/pterm/pterm"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
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
		readline.PcItem("map"),
		readline.PcItem("mine"),
		readline.PcItem("bola"),
		readline.PcItem("bopla"),
		readline.PcItem("bfla"),
		readline.PcItem("exhaust"),
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
	case "exit", "quit":
		pterm.NewStyle(pterm.FgRed).Println("\n[!] TERMINATING SESSION...")
		time.Sleep(500 * time.Millisecond)
		s.Active = false
	case "map":
		pterm.Info.Println("Executing Phase 2: Mapping Logic sequence...")
	case "exhaust":
		if len(parts) < 3 {
			pterm.Info.Println("Usage: exhaust <url> <parameter>")
			pterm.Info.Println("Example: exhaust https://api.target.com/v1/users limit")
			return
		}
		probe := &logic.ExhaustionContext{TargetURL: parts[1], ParamName: parts[2]}
		probe.FuzzPagination()

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
		{"auth", "Set identity tokens", "auth attacker <token>"},
		{"sessions", "View active tokens", "sessions"},
		{"bola", "Phase 3 BOLA test", "bola <url> <id>"},
		{"bopla", "BOPLA / API3 Mass Assignment", "bopla <url> '{\"id\":1}'"},
		{"bfla", "BFLA / API5 Method Shuffling", "bfla <url>"},
		{"exhaust", "API4 Pagination Fuzzing", "exhaust <url> limit"},
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
	default:
		pterm.Error.Printf("No manual entry for %s\n", cmd)
	}
}