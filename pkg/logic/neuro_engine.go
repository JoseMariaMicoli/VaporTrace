package logic

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/ai"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// NeuroEngine manages AI interactions
type NeuroEngine struct {
	Provider ai.LLMProvider
	Active   bool
	mu       sync.Mutex
}

var GlobalNeuro = &NeuroEngine{
	Active: false,
}

// Configure sets up the AI provider
func (n *NeuroEngine) Configure(providerType, apiKey, model, endpoint string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	switch strings.ToLower(providerType) {
	case "ollama":
		n.Provider = &ai.OllamaClient{}
	case "openai":
		n.Provider = &ai.OpenAIClient{}
	default:
		n.Provider = &ai.OllamaClient{} // Default to local
	}

	n.Provider.Configure(apiKey, model, endpoint)
	n.Active = true
	utils.TacticalLog(fmt.Sprintf("[magenta]NEURO:[-] Engine configured with %s (%s)", providerType, model))
}

// AnalyzeTrafficSnapshot takes raw HTTP dumps and queries the AI
func (n *NeuroEngine) AnalyzeTrafficSnapshot(reqDump, resDump string) {
	if !n.Active || n.Provider == nil {
		utils.TacticalLog("[yellow]NEURO:[-] AI Engine not active. Run 'neuro config'.")
		return
	}

	utils.TacticalLog("[magenta]NEURO:[-] Sending snapshot to Neural Engine...")
	utils.LogNeural(fmt.Sprintf("[yellow]>>> ANALYSIS REQUEST STARTED [%s][-]", time.Now().Format("15:04:05")))

	// Run Async
	go func() {
		prompt := fmt.Sprintf(ai.TrafficAnalysisPrompt, reqDump, resDump)

		analysis, err := n.Provider.Analyze(prompt)
		if err != nil {
			utils.TacticalLog(fmt.Sprintf("[red]NEURO ERROR:[-] %v", err))
			utils.LogNeural(fmt.Sprintf("[red]ERROR: %v[-]", err))
			return
		}

		// Format output for TUI
		formatted := fmt.Sprintf("\n[cyan]--- AI ANALYSIS REPORT ---\n[white]%s\n[cyan]------------------------[-]\n", analysis)

		// Send to F7 Tab
		utils.LogNeural(formatted)

		// Notification in Main Log
		utils.TacticalLog("[green]NEURO:[-] Analysis complete. Check F7 Tab.")
	}()
}

// GenerateAttackVectors asks the AI for specific payloads (Dry-run)
func (n *NeuroEngine) GenerateAttackVectors(context string, count int) {
	if !n.Active || n.Provider == nil {
		utils.TacticalLog("[yellow]NEURO:[-] AI Engine not active.")
		return
	}

	go func() {
		utils.TacticalLog(fmt.Sprintf("[magenta]NEURO:[-] Generating %d payloads for '%s'...", count, context))
		payloads, err := n.Provider.GeneratePayloads(context, count)
		if err != nil {
			utils.TacticalLog(fmt.Sprintf("[red]NEURO GEN FAIL:[-] %v", err))
			return
		}

		output := fmt.Sprintf("\n[cyan]--- NEURO PAYLOADS (%s) ---\n[white]", context)
		for _, p := range payloads {
			output += fmt.Sprintf("- %s\n", p)
		}
		output += "[cyan]----------------------------[-]\n"
		utils.LogNeural(output)
		utils.TacticalLog("[green]NEURO:[-] Payloads generated. Check F7 Tab.")
	}()
}

// AutonomousFuzz asks AI for payloads and executes them against the target (Live-Fire)
func (n *NeuroEngine) AutonomousFuzz(targetURL, method, context string, count int) {
	if !n.Active || n.Provider == nil {
		utils.TacticalLog("[yellow]NEURO:[-] AI Engine not active.")
		return
	}

	go func() {
		utils.TacticalLog(fmt.Sprintf("[magenta]NEURO-FUZZ:[-] Priming %d autonomous vectors for [%s]...", count, context))

		// 1. Generate
		payloads, err := n.Provider.GeneratePayloads(context, count)
		if err != nil {
			utils.TacticalLog(fmt.Sprintf("[red]NEURO GEN FAIL:[-] %v", err))
			return
		}

		utils.TacticalLog(fmt.Sprintf("[magenta]NEURO-FUZZ:[-] Engaging target %s...", targetURL))

		// 2. Execute
		for _, payload := range payloads {
			if strings.TrimSpace(payload) == "" {
				continue
			}

			// We send the payload as the request body. (Adjustable based on logic)
			req, _ := http.NewRequest(method, targetURL, bytes.NewBufferString(payload))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-VaporTrace-AI", "Autonomous-Fuzzer")

			// Inject Identity if available
			if CurrentSession.AttackerToken != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", CurrentSession.AttackerToken))
			}

			// This triggers the TacticalTransport -> Evasion -> Logger (F4) -> Interceptor
			resp, err := SafeDo(req, false, "NEURO-FUZZ")
			if err != nil {
				utils.TacticalLog(fmt.Sprintf("[red]Fuzz Error:[-] %v", err))
				continue
			}

			// 3. Heuristic Check
			statusColor := "[green]"
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				statusColor = "[yellow]"
			} else if resp.StatusCode >= 500 {
				statusColor = "[red]"
			}

			utils.TacticalLog(fmt.Sprintf("%sFUZZ [%d][-] | Payload: %s", statusColor, resp.StatusCode, payload))
			resp.Body.Close()
			time.Sleep(200 * time.Millisecond) // Built-in jitter
		}

		utils.TacticalLog("[green]NEURO-FUZZ:[-] Autonomous campaign complete. Check F4 for full traffic logs.")
	}()
}

// TestConnectivity runs a dummy prompt
func (n *NeuroEngine) TestConnectivity() {
	if !n.Active {
		utils.TacticalLog("[yellow]NEURO:[-] Engine is toggled OFF. Run 'neuro on'.")
		return
	}

	if n.Provider == nil {
		utils.TacticalLog("[red]NEURO CRITICAL:[-] Provider not configured. Run 'neuro config <provider> <model>' first.")
		n.Active = false // Safety toggle to prevent further calls
		return
	}

	go func() {
		utils.TacticalLog("[blue]NEURO:[-] Sending heartbeat packet to LLM...")
		resp, err := n.Provider.Analyze("Explain 'BOLA' vulnerability in one sentence.")
		if err != nil {
			utils.TacticalLog(fmt.Sprintf("[red]NEURO FAIL:[-] %v", err))
		} else {
			utils.TacticalLog("[green]NEURO ONLINE:[-] " + resp)
			utils.LogNeural("[green]CONNECTIVITY CHECK PASSED:[-] " + resp)
		}
	}()
}
