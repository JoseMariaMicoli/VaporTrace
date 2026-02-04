package logic

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/ai"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// Global state for Neuro Features (Neuro Inverter Toggle)
var NeuroInverterActive bool = false

// NeuroEngine manages AI interactions for Analysis, Exploit Gen, and Auto-Fuzzing
// Uses a Hybrid Architecture: Primary (Cloud) -> Secondary (Local Fallback)
type NeuroEngine struct {
	Primary   ai.LLMProvider
	Secondary ai.LLMProvider
	Active    bool
	mu        sync.Mutex
	lastCall  time.Time // Rate Limiter Timestamp
}

// Global singleton instance
var GlobalNeuro = &NeuroEngine{
	Active: false,
}

// Configure sets up the AI provider with Hydra Optimization (Prioritize Cloud)
func (n *NeuroEngine) Configure(providerType, apiKey, model, endpoint string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Default Fallback is always Ollama (Mistral)
	// Assumes standard localhost:11434 if not specified otherwise in args
	n.Secondary = &ai.OllamaClient{}
	n.Secondary.Configure("", "mistral", "http://localhost:11434")

	if providerType == "" {
		providerType = "openai" // Default to Cloud if unspecified
	}

	switch strings.ToLower(providerType) {
	case "ollama":
		if model == "" {
			model = "mistral"
		}
		// In pure local mode, Primary is also Ollama
		n.Primary = &ai.OllamaClient{}
		n.Primary.Configure("", model, endpoint)
		utils.TacticalLog("[yellow]NEURO:[-] Warning - Local Inference uses significant RAM. Ensure 8GB+ avail.")

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

	case "hybrid":
		// Explicit Hybrid Mode
		if model == "" {
			model = "gpt-4o"
		}
		n.Primary = &ai.OpenAIClient{}
		n.Primary.Configure(apiKey, model, endpoint)
		utils.TacticalLog("[magenta]NEURO:[-] Hybrid Brain Activated. Primary: OpenAI | Fallback: Ollama.")

	default:
		// Fallback
		if model == "" {
			model = "gpt-4o"
		}
		n.Primary = &ai.OpenAIClient{}
		n.Primary.Configure(apiKey, model, endpoint)
	}

	n.Active = true
	// Initialize Rate Limiter with a past timestamp so the FIRST request works instantly
	n.lastCall = time.Now().Add(-15 * time.Second)
	utils.TacticalLog(fmt.Sprintf("[magenta]NEURO:[-] Engine configured with %s (%s) + Smart Rate Limiting (High Latency Mode)", providerType, model))
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
			if strings.Contains(errStr, "429") || strings.Contains(strings.ToLower(errStr), "quota") || strings.Contains(strings.ToLower(errStr), "exhausted") {
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
	if !n.Active {
		utils.TacticalLog("[yellow]NEURO:[-] AI Engine not active. Run 'neuro config'.")
		return
	}

	utils.TacticalLog("[magenta]NEURO:[-] Intercepted Snapshot. Initiating Full-Spectrum Analysis...")
	utils.LogNeural(fmt.Sprintf("[yellow]>>> AUTONOMOUS SEQUENCE STARTED [%s][-]", time.Now().Format("15:04:05")))

	// Safely truncate dumps to avoid token limit hangs
	safeReq := truncateContext(reqDump, 1000) // Lowered token count further for 429 safety
	safeRes := truncateContext(resDump, 1000)

	// Async Execution
	go func() {
		// 1. Construct the Offensive Prompt
		prompt := fmt.Sprintf(`ACT AS AN OFFENSIVE SECURITY AI.
REQUEST:
%s

RESPONSE:
%s

TASK:
1. THINKING: Briefly explain your reasoning process (Chain of Thought).
2. IDENTIFY: Detect logic flaws (BOLA, BFLA, SQLi, XSS).
3. EXPLOIT: Generate 3 specific, high-probability exploit payloads.
4. COMPLIANCE: Map to MITRE and OWASP 2023.

Format:
CHAIN OF THOUGHT: <Reasoning>
ANALYSIS: <Summary>
---PAYLOADS---
<Payload1>
<Payload2>
...
---COMPLIANCE---
<Framework IDs>
`, safeReq, safeRes)

		// 2. Query LLM (Hybrid Execution)
		response, err := n.ExecuteQuery(prompt)
		if err != nil {
			utils.TacticalLog(fmt.Sprintf("[red]NEURO ERROR:[-] %v", err))
			utils.LogNeural(fmt.Sprintf("[red]ERROR: %v[-]", err))
			return
		}

		// 3. Parse Sections
		analysis, payloads, compliance := n.parseAIOutput(response)

		// 4. Report to UI (F6 Neural Tab)
		report := fmt.Sprintf("\n[cyan]=== TACTICAL ANALYSIS ===[-]\n[white]%s[-]\n\n", analysis)
		if compliance != "" {
			report += fmt.Sprintf("[blue]=== COMPLIANCE ===[-]\n[gray]%s[-]\n", compliance)
		}
		utils.LogNeural(report)

		// 5. AUTO-FUZZING: Execute Generated Exploits
		if len(payloads) > 0 {
			targetURL, method := n.extractTargetInfo(reqDump)
			utils.TacticalLog(fmt.Sprintf("[magenta]NEURO-AUTO:[-] Extracted %d exploits. Engaging Target %s...", len(payloads), targetURL))

			// Engage "Smart Fuzzer" logic
			n.executeSmartAttack(targetURL, method, payloads)
		} else {
			utils.TacticalLog("[yellow]NEURO:[-] No viable exploits generated by AI.")
		}
	}()
}

// PerformNeuroBrute implements the Fuzzing logic triggered by Ctrl+B
// It takes context (body or headers) and generates high-entropy mutations.
func (n *NeuroEngine) PerformNeuroBrute(seedBody string) {
	if !n.Active {
		utils.TacticalLog("[yellow]NEURO:[-] Engine inactive.")
		return
	}

	utils.TacticalLog("[blue]NEURO:[-] Generating intelligent mutations for context...")

	go func() {
		// Reduce context drastically for brute gen to save tokens
		truncBody := truncateContext(seedBody, 400)
		prompt := fmt.Sprintf("Generate 5 fuzzing mutations for this data to test for SQLi and BOLA. Return ONLY raw strings/JSON:\n%s", truncBody)

		resp, err := n.ExecuteQuery(prompt)
		if err != nil {
			utils.TacticalLog(fmt.Sprintf("[red]Brute Gen Failed: %v", err))
			return
		}

		payloads := strings.Split(resp, "\n")
		utils.LogNeural("[blue]BRUTE:[-] Generated Mutations. Check Log.")

		for _, p := range payloads {
			p = strings.TrimSpace(p)
			if len(p) > 2 {
				utils.LogNeural(fmt.Sprintf("[white]MUTATION:[-] %s", p))
			}
		}
	}()
}

// executeSmartAttack handles the "Live-Fire" execution of AI payloads
func (n *NeuroEngine) executeSmartAttack(targetURL, method string, payloads []string) {
	client := GlobalClient
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}

	// 1. ESTABLISH BASELINE LATENCY (Calibration)
	utils.LogNeural("[blue]NEURO-AUTO:[-] Calibrating baseline network latency...")
	baseReq, _ := http.NewRequest(method, targetURL, nil)

	if CurrentSession.AttackerToken != "" {
		baseReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", CurrentSession.AttackerToken))
	}

	startBase := time.Now()
	baseResp, errBase := client.Do(baseReq)
	baselineLatency := time.Since(startBase)

	if errBase == nil {
		baseResp.Body.Close()
	} else {
		baselineLatency = 200 * time.Millisecond // Fallback
	}
	utils.LogNeural(fmt.Sprintf("[blue]NEURO-AUTO:[-] Baseline established: %v", baselineLatency))

	// 2. FIRE PAYLOADS
	for i, payload := range payloads {
		if payload == "" {
			continue
		}

		utils.LogNeural(fmt.Sprintf("[yellow]>>> FIRING VECTOR %d/%d: %s[-]", i+1, len(payloads), shortPayload(payload)))

		var req *http.Request
		// Intelligent Injection Strategy
		if method == "POST" || method == "PUT" || method == "PATCH" {
			req, _ = http.NewRequest(method, targetURL, bytes.NewBufferString(payload))
			req.Header.Set("Content-Type", "application/json")
		} else {
			// Query Parameter Injection
			u := targetURL
			if strings.Contains(u, "?") {
				u += "&fuzz=" + url.QueryEscape(payload)
			} else {
				u += "?fuzz=" + url.QueryEscape(payload)
			}
			req, _ = http.NewRequest(method, u, nil)
		}

		if CurrentSession.AttackerToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", CurrentSession.AttackerToken))
		}

		req.Header.Set("X-Neuro-Engine", "Automated-Exploit")
		req.Header.Set("User-Agent", "VaporTrace-Neuro/1.0")

		startAttack := time.Now()
		resp, err := client.Do(req)
		attackDuration := time.Since(startAttack)

		if err != nil {
			if strings.Contains(err.Error(), "Timeout") || attackDuration > 10*time.Second {
				utils.TacticalLog(fmt.Sprintf("[red]CRITICAL: Request Timed Out (%v). Possible Heavy SQLi.[/]", attackDuration))
				n.recordTimeBasedSQLi(targetURL, payload, attackDuration, baselineLatency)
			}
		} else {
			n.evaluateResponse(resp, payload, targetURL, attackDuration, baselineLatency)
			resp.Body.Close()
		}

		// *** CRITICAL RATE LIMIT FIX (LATAM/FREE TIER) ***
		// Wait 6 seconds between automated attacks to protect Quota
		time.Sleep(6000 * time.Millisecond)
	}
	utils.TacticalLog("[green]NEURO-AUTO:[-] Sequence Complete. Verified results in Logs.")
}

func (n *NeuroEngine) evaluateResponse(resp *http.Response, payload, target string, latency time.Duration, baseline time.Duration) {
	// LOGIC 1: DIFFERENTIAL TIMING (Time-Based SQLi)
	if latency > 4*time.Second && latency > (baseline*3) {
		n.recordTimeBasedSQLi(target, payload, latency, baseline)
		return
	}

	// LOGIC 2: ERROR/CRASH
	if resp.StatusCode >= 500 {
		utils.TacticalLog(fmt.Sprintf("[red]CRITICAL HIT (%d): %s (Lat: %v)[-]", resp.StatusCode, shortPayload(payload), latency))

		if db.DB != nil {
			utils.RecordFinding(db.Finding{
				Phase:        "PHASE 10.6: NEURO-EXPLOIT",
				Command:      "neuro",
				Target:       target,
				Details:      fmt.Sprintf("Server Error (%d) triggered by payload. Possible Injection/RCE.", resp.StatusCode),
				Status:       "EXPLOITED",
				OWASP_ID:     "API10:2023",
				MITRE_ID:     "T1190",
				CVSS_Numeric: 9.0,
			})
		}
		return
	}

	// LOGIC 3: BYPASS
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		utils.TacticalLog(fmt.Sprintf("[green]POTENTIAL BYPASS (%d): %s[-]", resp.StatusCode, shortPayload(payload)))
		if db.DB != nil {
			utils.RecordFinding(db.Finding{
				Phase:        "PHASE 10.6: NEURO-EXPLOIT",
				Command:      "neuro",
				Target:       target,
				Details:      fmt.Sprintf("Logic Bypass (%d). Payload accepted.", resp.StatusCode),
				Status:       "VULNERABLE",
				OWASP_ID:     "API1:2023",
				MITRE_ID:     "T1595",
				CVSS_Numeric: 7.5,
			})
		}
	} else {
		utils.TacticalLog(fmt.Sprintf("[gray]Miss (%d) | Lat: %v | %s[-]", resp.StatusCode, latency, shortPayload(payload)))
	}
}

func (n *NeuroEngine) recordTimeBasedSQLi(target, payload string, latency, baseline time.Duration) {
	msg := fmt.Sprintf("[red]!!! TIME-BASED SQLI CONFIRMED !!! Latency: %v (Base: %v) | Vector: %s[-]", latency, baseline, payload)
	utils.LogNeural(msg)
	utils.TacticalLog(msg)

	if db.DB != nil {
		utils.RecordFinding(db.Finding{
			Phase:        "PHASE 10.6: NEURO-EXPLOIT",
			Command:      "neuro",
			Target:       target,
			Details:      fmt.Sprintf("High-Confidence Time-Based Blind SQLi. Response delayed by %v (Baseline: %v).", latency, baseline),
			Status:       "CRITICAL",
			OWASP_ID:     "API8:2023 Injection",
			MITRE_ID:     "T1190",
			MitreTactic:  "Initial Access",
			CVSS_Numeric: 9.8,
		})
	}
}

func (n *NeuroEngine) parseAIOutput(raw string) (analysis string, payloads []string, compliance string) {
	lines := strings.Split(raw, "\n")
	section := "ANALYSIS"

	for _, line := range lines {
		cleanLine := strings.TrimSpace(line)

		if cleanLine == "---PAYLOADS---" {
			section = "PAYLOADS"
			continue
		} else if cleanLine == "---COMPLIANCE---" {
			section = "COMPLIANCE"
			continue
		} else if strings.HasPrefix(cleanLine, "ANALYSIS:") || strings.HasPrefix(cleanLine, "CHAIN OF THOUGHT:") {
			section = "ANALYSIS"
		}

		switch section {
		case "ANALYSIS":
			analysis += cleanLine + "\n"
		case "PAYLOADS":
			if cleanLine != "" && !strings.HasPrefix(cleanLine, "`") && !strings.HasPrefix(cleanLine, "Analysis") {
				cleanLine = strings.TrimPrefix(cleanLine, "- ")
				cleanLine = strings.TrimPrefix(cleanLine, "* ")
				cleanLine = strings.TrimPrefix(cleanLine, "1. ")
				cleanLine = strings.Trim(cleanLine, "\"")
				cleanLine = strings.Trim(cleanLine, "'")
				if len(cleanLine) > 2 {
					payloads = append(payloads, cleanLine)
				}
			}
		case "COMPLIANCE":
			compliance += line + "\n"
		}
	}
	return
}

func (n *NeuroEngine) extractTargetInfo(reqDump string) (string, string) {
	if reqDump == "" {
		return "", "GET"
	}
	reader := bufio.NewReader(strings.NewReader(reqDump))
	requestLine, _ := reader.ReadString('\n')
	parts := strings.Fields(requestLine)
	if len(parts) >= 2 {
		method := parts[0]
		path := parts[1]
		host := ""
		scanner := bufio.NewScanner(strings.NewReader(reqDump))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Host:") {
				host = strings.TrimSpace(strings.TrimPrefix(line, "Host:"))
				break
			}
		}
		scheme := "http"
		if strings.Contains(reqDump, "443") {
			scheme = "https"
		}
		return fmt.Sprintf("%s://%s%s", scheme, host, path), method
	}
	return "", "GET"
}

func truncateContext(s string, limit int) string {
	if len(s) > limit {
		return s[:limit] + "...[TRUNCATED]"
	}
	return s
}

func shortPayload(p string) string {
	if len(p) > 50 {
		return p[:47] + "..."
	}
	return p
}

// Legacy wrappers to satisfy interface if needed, or legacy calls
func (n *NeuroEngine) GenerateAttackVectors(context string, count int) {
	if !n.Active {
		utils.TacticalLog("[yellow]NEURO:[-] AI Engine not active.")
		return
	}
	go func() {
		// Use ExecuteQuery logic to access primary/secondary via unified interface,
		// but since GeneratePayloads is an interface method, we access fields directly
		// with checks to ensure we aren't calling nil.
		var payloads []string
		var err error

		if n.Primary != nil {
			n.enforceRateLimit()
			payloads, err = n.Primary.GeneratePayloads(context, count)
		}

		// Fallback if Primary failed or nil
		if (err != nil || n.Primary == nil) && n.Secondary != nil {
			payloads, err = n.Secondary.GeneratePayloads(context, count)
		}

		if err == nil {
			output := fmt.Sprintf("\n[cyan]--- NEURO PAYLOADS (%s) ---\n[white]", context)
			for _, p := range payloads {
				output += fmt.Sprintf("- %s\n", p)
			}
			utils.LogNeural(output)
		}
	}()
}

func (n *NeuroEngine) AutonomousFuzz(targetURL, method, context string, count int) {
	if !n.Active {
		return
	}
	go func() {
		var payloads []string
		var err error

		// Try Primary with Rate Limit
		if n.Primary != nil {
			n.enforceRateLimit()
			payloads, err = n.Primary.GeneratePayloads(context, count)
		}

		// Try Secondary if Primary failed
		if (err != nil || n.Primary == nil) && n.Secondary != nil {
			payloads, err = n.Secondary.GeneratePayloads(context, count)
		}

		if err == nil {
			var clean []string
			for _, p := range payloads {
				clean = append(clean, strings.TrimSpace(p))
			}
			n.executeSmartAttack(targetURL, method, clean)
		}
	}()
}

func (n *NeuroEngine) TestConnectivity() {
	if !n.Active {
		utils.TacticalLog("[yellow]NEURO:[-] Engine is toggled OFF. Run 'neuro on'.")
		return
	}
	go func() {
		utils.TacticalLog("[blue]NEURO:[-] Sending heartbeat packet to Brain...")
		// Use ExecuteQuery to test the whole fallback chain
		resp, err := n.ExecuteQuery("Ping")
		if err != nil {
			utils.TacticalLog(fmt.Sprintf("[red]NEURO FAIL:[-] %v", err))
		} else {
			utils.LogNeural("[green]CONNECTIVITY CHECK:[-] " + resp)
			utils.TacticalLog("[green]NEURO ONLINE:[-] Check Neural Tab.")
		}
	}()
}
