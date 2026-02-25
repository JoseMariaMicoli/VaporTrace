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
		if model == "" {
			model = "gpt-4o"
		}
		n.Primary = &ai.OpenAIClient{}
		n.Primary.Configure(apiKey, model, endpoint)
		utils.TacticalLog("[green]NEURO:[-] OpenAI Cloud Selected.")

	case "google", "gemini":
		if model == "" {
			// Using the stable alias to avoid beta quota issues in LATAM regions
			model = "gemini-1.5-flash"
		}
		n.Primary = &ai.GeminiClient{}
		n.Primary.Configure(apiKey, model, endpoint)
		utils.TacticalLog(fmt.Sprintf("[cyan]NEURO:[-] Google Gemini Selected (%s).", model))

	case "groq":
		if model == "" {
			model = "llama-3.1-8b-instant"
		}
		if endpoint == "" {
			endpoint = "https://api.groq.com/openai/v1/chat/completions"
		}
		n.Primary = &ai.OpenAIClient{} // Groq is OpenAI Compatible
		n.Primary.Configure(apiKey, model, endpoint)
		utils.TacticalLog("[cyan]NEURO:[-] Groq LPU Cloud Selected.")

	case "hybrid":
		// Explicit Hybrid Mode
		if model == "" {
			model = "gpt-4o"
		}
		n.Primary = &ai.OpenAIClient{}
		n.Primary.Configure(apiKey, model, endpoint)
		utils.TacticalLog("[magenta]NEURO:[-] Hybrid Brain Activated. Primary: OpenAI | Fallback: Ollama.")

	default:
		n.Provider = &ai.OllamaClient{} // Default to local
	}

	n.Provider.Configure(apiKey, model, endpoint)
	n.Active = true
	utils.TacticalLog(fmt.Sprintf("[magenta]NEURO:[-] Engine configured with %s (%s)", providerType, model))
}

// enforceRateLimit ensures we don't hit 429s by spacing requests
func (n *NeuroEngine) enforceRateLimit() {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Strict 6 seconds between calls for Free Tier safety in high-latency regions
	elapsed := time.Since(n.lastCall)
	if elapsed < 6*time.Second {
		wait := 6*time.Second - elapsed
		time.Sleep(wait)
	}
	n.lastCall = time.Now()
}

// ExecuteQuery tries primary provider, and immediately falls back to secondary on 429
func (n *NeuroEngine) ExecuteQuery(prompt string) (string, error) {
	var primaryErr error

	if n.Primary != nil {
		// 1. Rate Limit Check
		n.enforceRateLimit()

		// 2. Primary Attempt
		res, err := n.Primary.Analyze(prompt)

		// 3. Smart Error Handling
		if err != nil {
			primaryErr = err
			errStr := err.Error()

			// Detect 429 / Quota issues explicitly
			if strings.Contains(errStr, "429") || strings.Contains(strings.ToLower(errStr), "quota") || strings.Contains(strings.ToLower(errStr), "exhausted") || strings.Contains(strings.ToLower(errStr), "rate limit") {
				utils.TacticalLog("[red]NEURO:[-] Cloud Brain Quota/Rate-Limit (429) Hit.")
				utils.TacticalLog("[yellow]NEURO:[-] BYPASSING RETRY -> Engaging Local Brain (Ollama) Immediately.")
				// We do NOT retry primary here.
				// Fallthrough directly to secondary.
			} else {
				// Other errors (Network, Auth) log and fallthrough
				utils.TacticalLog(fmt.Sprintf("[red]NEURO:[-] Primary Brain Error: %v. Switching to Fallback...", err))
			}
		} else if res != "" {
			return res, nil
		}
	}

	// 4. Fallback Execution (Local / Ollama)
	if n.Secondary != nil {
		utils.TacticalLog("[blue]NEURO:[-] Using Local Mistral (Ollama)...")
		res, err := n.Secondary.Analyze(prompt)
		if err != nil {
			// If both fail, return a combined error message
			finalErr := fmt.Errorf("Hybrid Failure - Cloud: %v | Local: %v", primaryErr, err)
			return "", finalErr
		}
		return res, nil
	}

	return "", fmt.Errorf("all neural paths failed (Primary: %v, No Secondary configured)", primaryErr)
}

// AnalyzeTrafficSnapshot is the Core Trigger (Ctrl+A).
// It safely executes the analysis in a background thread to keep the UI responsive.
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
