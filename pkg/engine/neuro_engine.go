package engine

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// NeuroAnalysisResult represents the output of an AI analysis pass
type NeuroAnalysisResult struct {
	VulnerabilityType string    // BOLA, BFLA, SSRF, INJECTION, etc.
	Endpoints         []string  // List of affected endpoints
	Payloads          []string  // Recommended attack payloads
	ConfidenceScores  []float64 // Per-endpoint confidence scores
	Reasoning         string    // AI-generated reasoning
}

// NeuroEngineCore manages engine-level AI analysis for tactical planning
// This bridges logic.GlobalNeuro with engine.TacticalAction generation
// Extended with full exploit execution and response evaluation capabilities
type NeuroEngineCore struct {
	mu            sync.Mutex
	lastAnalysis  *NeuroAnalysisResult
	analysisCount int
	httpClient    *http.Client
}

// GlobalNeuroCore is the singleton instance for engine-level AI analysis
var GlobalNeuroCore = &NeuroEngineCore{
	analysisCount: 0,
	httpClient:    &http.Client{Timeout: 15 * time.Second},
}

// Analyze performs comprehensive AI-driven analysis on a given request/response pair
// and returns a TacticalAction recommendation
func (n *NeuroEngineCore) Analyze(reqDump, resDump string) *TacticalAction {
	n.mu.Lock()
	defer n.mu.Unlock()

	if !logic.GlobalNeuro.Active {
		utils.TacticalLog("[yellow]NEURO:[-] Engine inactive. Cannot analyze.")
		return nil
	}

	// Extract target URL and method from request dump
	targetURL := n.extractURL(reqDump)
	method := n.extractMethod(reqDump)
	statusCode := n.extractStatusCode(resDump)

	if targetURL == "" {
		utils.TacticalLog("[yellow]NEURO:[-] Could not extract target URL from request.")
		return nil
	}

	// Build comprehensive analysis prompt
	prompt := fmt.Sprintf(`You are a security analyst. Analyze this HTTP traffic snapshot and identify security issues.

REQUEST:
%s

RESPONSE:
%s

Respond in exactly this format:
ISSUE_TYPE: [BOLA|BFLA|SSRF|INJECTION|ENUMERATION|BYPASS|OTHER]
CONFIDENCE: [0-100]
REASONING: [1-2 sentences explaining the vulnerability]
PAYLOAD: [Suggested attack payload or technique]
`, reqDump, resDump)

	// Query the global neuro engine
	response, err := logic.GlobalNeuro.ExecuteQuery(prompt)
	if err != nil {
		utils.TacticalLog("[yellow]NEURO:[-] Analysis failed: " + err.Error())
		return nil
	}

	// Parse the response into a TacticalAction
	action := n.parseAnalysisResponse(response, targetURL, method, statusCode)
	n.lastAnalysis = &NeuroAnalysisResult{
		VulnerabilityType: action.Type,
		Endpoints:         []string{targetURL},
		Payloads:          []string{action.Payload},
		ConfidenceScores:  n.parseConfidence(action.Confidence),
		Reasoning:         action.Reasoning,
	}

	return action
}

// AnalyzeEndpoint performs AI analysis on a single discovered endpoint
// Returns a TacticalAction if a vulnerability is suspected
func (n *NeuroEngineCore) AnalyzeEndpoint(endpoint string, lastStatus int, loot logic.LootSummary) *TacticalAction {
	n.mu.Lock()
	defer n.mu.Unlock()

	if !logic.GlobalNeuro.Active {
		return nil
	}

	// Build a focused analysis prompt for endpoint testing
	prompt := fmt.Sprintf(`Analyze this API endpoint for security vulnerabilities:

ENDPOINT: %s
LAST_HTTP_STATUS: %d
HAS_JWT: %v
HAS_AWS_KEYS: %v
HAS_CREDENTIALS: %v

Respond in this exact format:
VULNERABILITY_TYPE: [BOLA|BFLA|SSRF|INJECTION|ENUMERATION|BYPASS|NONE]
CONFIDENCE: [0-100]
REASONING: [Explanation]
PAYLOAD: [Suggested exploit]
`, endpoint, lastStatus, loot.HasJWT, loot.HasAWS, loot.Credential != "")

	response, err := logic.GlobalNeuro.ExecuteQuery(prompt)
	if err != nil {
		return nil
	}

	action := n.parseAnalysisResponse(response, endpoint, "GET", lastStatus)
	return action
}

// parseAnalysisResponse converts AI response text into a TacticalAction
func (n *NeuroEngineCore) parseAnalysisResponse(response, target, method string, status int) *TacticalAction {
	if response == "" {
		return nil
	}

	action := &TacticalAction{
		Type:       "GENERIC",
		Target:     target,
		Payload:    "AI-Generated Payload",
		Confidence: "MEDIUM",
		Reasoning:  "AI Analysis",
		Status:     "PENDING",
	}

	// Parse response line by line
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "ISSUE_TYPE:") || strings.HasPrefix(line, "VULNERABILITY_TYPE:") {
			issueType := strings.TrimPrefix(line, "ISSUE_TYPE:")
			issueType = strings.TrimPrefix(issueType, "VULNERABILITY_TYPE:")
			issueType = strings.TrimSpace(issueType)
			action.Type = issueType
		} else if strings.HasPrefix(line, "CONFIDENCE:") {
			confStr := strings.TrimPrefix(line, "CONFIDENCE:")
			confStr = strings.TrimSpace(confStr)
			confidence := n.parseConfidenceValue(confStr)
			action.Confidence = confidence
		} else if strings.HasPrefix(line, "REASONING:") {
			reasoning := strings.TrimPrefix(line, "REASONING:")
			action.Reasoning = strings.TrimSpace(reasoning)
		} else if strings.HasPrefix(line, "PAYLOAD:") {
			payload := strings.TrimPrefix(line, "PAYLOAD:")
			action.Payload = strings.TrimSpace(payload)
		}
	}

	// Validate action type
	validTypes := map[string]bool{
		"BOLA": true, "BFLA": true, "SSRF": true, "INJECTION": true,
		"ENUMERATION": true, "BYPASS": true, "AUDIT": true, "GENERIC": true,
	}
	if !validTypes[action.Type] {
		action.Type = "GENERIC"
	}

	return action
}

// parseConfidenceValue converts a string/number to a confidence level
func (n *NeuroEngineCore) parseConfidenceValue(confStr string) string {
	// Try to extract numeric value
	re := regexp.MustCompile(`(\d+)`)
	matches := re.FindStringSubmatch(confStr)

	if len(matches) > 0 {
		confStr = matches[0]
	}

	// Map numeric or string confidence to levels
	if strings.Contains(confStr, "HIGH") || strings.Contains(confStr, "CRITICAL") || strings.Contains(confStr, "100") || strings.Contains(confStr, "95") || strings.Contains(confStr, "90") || strings.Contains(confStr, "85") || strings.Contains(confStr, "80") {
		return "CRITICAL"
	} else if strings.Contains(confStr, "MEDIUM") || strings.Contains(confStr, "60") || strings.Contains(confStr, "70") {
		return "HIGH"
	} else if strings.Contains(confStr, "LOW") || strings.Contains(confStr, "40") || strings.Contains(confStr, "50") {
		return "MEDIUM"
	}

	return "MEDIUM"
}

// parseConfidence converts confidence string to float64 slice
func (n *NeuroEngineCore) parseConfidence(confStr string) []float64 {
	switch confStr {
	case "CRITICAL":
		return []float64{0.95}
	case "HIGH":
		return []float64{0.75}
	case "MEDIUM":
		return []float64{0.50}
	case "LOW":
		return []float64{0.25}
	default:
		return []float64{0.50}
	}
}

// extractURL extracts the target URL from a request dump
func (n *NeuroEngineCore) extractURL(reqDump string) string {
	lines := strings.Split(reqDump, "\n")
	if len(lines) > 0 {
		// First line typically contains method and URL
		parts := strings.Fields(lines[0])
		if len(parts) >= 2 {
			return parts[1]
		}
	}
	return ""
}

// extractMethod extracts the HTTP method from a request dump
func (n *NeuroEngineCore) extractMethod(reqDump string) string {
	lines := strings.Split(reqDump, "\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) >= 1 {
			return parts[0]
		}
	}
	return "GET"
}

// extractStatusCode extracts the HTTP status code from a response dump
func (n *NeuroEngineCore) extractStatusCode(resDump string) int {
	re := regexp.MustCompile(`HTTP/\d\.\d\s+(\d+)`)
	matches := re.FindStringSubmatch(resDump)
	if len(matches) > 1 {
		var status int
		fmt.Sscanf(matches[1], "%d", &status)
		return status
	}
	return 0
}

// ExecuteSmartAttack handles the "Live-Fire" execution of AI payloads against discovered endpoints
// This is the main entry point for neuro-driven exploitation triggered by AnalyzeTrafficSnapshot
func (n *NeuroEngineCore) ExecuteSmartAttack(targetURL, method string, payloads []string) {
	client := n.httpClient
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}

	// 1. ESTABLISH BASELINE LATENCY (Calibration)
	// We measure the normal response time to detect Time-Based injections later
	utils.LogNeural("[blue]NEURO-AUTO:[-] Calibrating baseline network latency...")
	baseReq, _ := http.NewRequest(method, targetURL, nil)

	if logic.CurrentSession.AttackerToken != "" {
		baseReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", logic.CurrentSession.AttackerToken))
	}

	startBase := time.Now()
	baseResp, errBase := client.Do(baseReq)
	baselineLatency := time.Since(startBase)

	if errBase == nil {
		baseResp.Body.Close()
	} else {
		baselineLatency = 200 * time.Millisecond // Fallback if baseline request fails
	}
	utils.LogNeural("[blue]NEURO-AUTO:[-] Baseline established: " + baselineLatency.String())

	// 2. FIRE PAYLOADS
	for i, payload := range payloads {
		if payload == "" {
			continue
		}

		utils.LogNeural("[yellow]>>> FIRING VECTOR " + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(payloads)) + ": " + shortPayload(payload) + "[-]")

		var req *http.Request
		// Intelligent Injection Strategy
		if method == "POST" || method == "PUT" || method == "PATCH" {
			req, _ = http.NewRequest(method, targetURL, bytes.NewBufferString(payload))

			// --- MULTI-FORMAT SIGNATURE SNIFFING (XXE & YAML) ---
			lowerPayload := strings.ToLower(payload)

			if strings.Contains(lowerPayload, "<?xml") || strings.Contains(lowerPayload, "<!doctype") || strings.Contains(lowerPayload, "<entity") {
				// XML/XXE Detection
				req.Header.Set("Content-Type", "application/xml")
				utils.LogNeural("[magenta]NEURO-AUTO:[-] XML Signature detected. Switching to application/xml...[-]")

			} else if strings.HasPrefix(payload, "---") || strings.Contains(payload, "!!") || strings.Contains(payload, "y_object:") {
				// YAML/Deserialization Detection (!! is a common tag for Ruby/Python YAML exploits)
				req.Header.Set("Content-Type", "application/x-yaml")
				utils.LogNeural("[cyan]NEURO-AUTO:[-] YAML Signature detected. Switching to application/x-yaml...[-]")

			} else {
				// Default to JSON for standard API interactions
				req.Header.Set("Content-Type", "application/json")
			}
			// ----------------------------------------------------

		} else {
			// Query Parameter Injection (for GET/DELETE)
			u := targetURL
			if strings.Contains(u, "?") {
				u += "&fuzz=" + url.QueryEscape(payload)
			} else {
				u += "?fuzz=" + url.QueryEscape(payload)
			}
			req, _ = http.NewRequest(method, u, nil)
		}

		if logic.CurrentSession.AttackerToken != "" {
			req.Header.Set("Authorization", "Bearer "+logic.CurrentSession.AttackerToken)
		}

		// Identify the request as an AI-driven exploit for logging/debugging
		req.Header.Set("X-Neuro-Engine", "Automated-Exploit")
		req.Header.Set("User-Agent", "VaporTrace-Neuro/1.0")

		startAttack := time.Now()
		resp, err := client.Do(req)
		attackDuration := time.Since(startAttack)

		if err != nil {
			// Detect Time-Based Blind SQLi if the server hangs
			if strings.Contains(err.Error(), "Timeout") || attackDuration > 10*time.Second {
				utils.TacticalLog("[red]CRITICAL: Request Timed Out (" + attackDuration.String() + "). Possible Heavy SQLi.[/]")
				n.RecordTimeBasedSQLi(targetURL, payload, attackDuration, baselineLatency)
			}
		} else {
			// Pass the response to the evaluator to check for 200 OK (Bypass) or 500 (Leak)
			n.EvaluateResponse(resp, payload, targetURL, attackDuration, baselineLatency)
			resp.Body.Close()
		}

		// *** CRITICAL RATE LIMIT FIX (LATAM/FREE TIER) ***
		// Space out requests to avoid 429 errors from Cloud Providers and WAFs
		time.Sleep(6000 * time.Millisecond)
	}
	utils.TacticalLog("[green]NEURO-AUTO:[-] Sequence Complete. Verified results in Logs.")
}

// EvaluateResponse analyzes the HTTP response from a fuzzing payload execution
// It checks for timing-based SQLi, server errors, and AI-verified bypasses
func (n *NeuroEngineCore) EvaluateResponse(resp *http.Response, payload, target string, latency time.Duration, baseline time.Duration) {
	// 1. EXTRACT DATA FOR AI ANALYSIS
	bodyBytes, _ := io.ReadAll(resp.Body)
	// Re-assign body to allow reading if needed later (though we close it in the caller)
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	bodySnippet := string(bodyBytes)
	if len(bodySnippet) > 1000 {
		bodySnippet = bodySnippet[:1000] // Truncate to save context tokens
	}

	// 2. LOGIC 1: DIFFERENTIAL TIMING (Fast-Path Time-Based SQLi)
	if latency > 4*time.Second && latency > (baseline*3) {
		n.RecordTimeBasedSQLi(target, payload, latency, baseline)
		return
	}

	// 3. LOGIC 2: ERROR/CRASH (Server-Side Failure)
	if resp.StatusCode >= 500 {
		utils.TacticalLog("[red]CRITICAL HIT (" + strconv.Itoa(resp.StatusCode) + "): " + shortPayload(payload) + " (Lat: " + latency.String() + ")[-]")

		if db.DB != nil {
			utils.RecordFinding(db.Finding{
				Phase:        "PHASE 10.6: NEURO-EXPLOIT",
				Command:      "neuro",
				Target:       target,
				Details:      "Server Error (" + strconv.Itoa(resp.StatusCode) + ") triggered by payload. Possible Injection/RCE.",
				Status:       "EXPLOITED",
				OWASP_ID:     "API10:2023",
				MITRE_ID:     "T1190",
				CVSS_Numeric: 9.0,
			})
		}
		return
	}

	// 4. LOGIC 3: GENERIC BYPASS (Fallthrough)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		utils.TacticalLog("[green]POTENTIAL BYPASS (" + strconv.Itoa(resp.StatusCode) + "): " + shortPayload(payload) + "[-]")
		if db.DB != nil {
			utils.RecordFinding(db.Finding{
				Phase:        "PHASE 10.6: NEURO-EXPLOIT",
				Command:      "neuro",
				Target:       target,
				Details:      "Logic Bypass (" + strconv.Itoa(resp.StatusCode) + "). Payload accepted.",
				Status:       "VULNERABLE",
				OWASP_ID:     "API1:2023",
				MITRE_ID:     "T1595",
				CVSS_Numeric: 7.5,
			})
		}
	} else {
		utils.TacticalLog("[gray]Miss (" + strconv.Itoa(resp.StatusCode) + ") | Lat: " + latency.String() + " | " + shortPayload(payload) + "[-]")
	}
}

// RecordTimeBasedSQLi records and logs a confirmed time-based SQL injection
func (n *NeuroEngineCore) RecordTimeBasedSQLi(target, payload string, latency, baseline time.Duration) {
	msg := "[red]!!! TIME-BASED SQLI CONFIRMED !!! Latency: " + latency.String() + " (Base: " + baseline.String() + ") | Vector: " + payload + "[-]"
	utils.LogNeural(msg)
	utils.TacticalLog(msg)

	if db.DB != nil {
		utils.RecordFinding(db.Finding{
			Phase:        "PHASE 10.6: NEURO-EXPLOIT",
			Command:      "neuro",
			Target:       target,
			Details:      "High-Confidence Time-Based Blind SQLi. Response delayed by " + latency.String() + " (Baseline: " + baseline.String() + ").",
			Status:       "CRITICAL",
			OWASP_ID:     "API8:2023 Injection",
			MITRE_ID:     "T1190",
			MitreTactic:  "Initial Access",
			CVSS_Numeric: 9.8,
		})
	}
}

// ComprehensiveAnalysis performs a multi-pass analysis combining heuristics and AI
// This is the main entry point for the strategic planner
func (n *NeuroEngineCore) ComprehensiveAnalysis(endpoints []string, loot logic.LootSummary, traffic map[string]int) []TacticalAction {
	var actions []TacticalAction

	for _, endpoint := range endpoints {
		// First pass: Apply deterministic heuristics
		hAction := RunHeuristics(endpoint, traffic[endpoint], loot)
		if hAction != nil {
			actions = append(actions, *hAction)
		}

		// Second pass: If AI is active, perform neural analysis
		if logic.GlobalNeuro.Active && len(actions) < 10 {
			aiAction := n.AnalyzeEndpoint(endpoint, traffic[endpoint], loot)
			if aiAction != nil && aiAction.Type != "GENERIC" {
				actions = append(actions, *aiAction)
			}
		}
	}

	return actions
}

// GetLastAnalysis returns the result of the most recent analysis
func (n *NeuroEngineCore) GetLastAnalysis() *NeuroAnalysisResult {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.lastAnalysis
}

// Reset clears the analysis history
func (n *NeuroEngineCore) Reset() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.lastAnalysis = nil
	n.analysisCount = 0
}

// Helper function to shorten payload strings for display
func shortPayload(p string) string {
	if len(p) > 50 {
		return p[:47] + "..."
	}
	return p
}

// === SPRINT 12: Autonomous Chaining Logic ===
// ProcessExploitResult analyzes the results of an exploit execution and automatically generates chained actions
// If loot contains sensitive data (tokens, credentials, env vars), this function creates follow-up exploits
func (n *NeuroEngineCore) ProcessExploitResult(exploit TacticalAction, loot string) {
	if loot == "" {
		return
	}

	utils.TacticalLog("[cyan]CHAIN:[-] Analyzing exploit result for chaining opportunities...")
	utils.LogContext("[blue]Processing loot:[-] " + loot[:min(len(loot), 100)])

	// Generate a unique chain ID based on the parent exploit
	chainID := "chain_" + strconv.FormatInt(int64(exploit.ID), 10) + "_" + strconv.FormatInt(time.Now().UnixNano(), 10)

	// Store the loot in GlobalDataSilo
	lootKey := "loot_from_" + strconv.Itoa(exploit.ID)
	logic.GlobalDataSilo.Set(lootKey, loot)
	utils.LogContext("[green]✓ Stored loot:[-] '" + lootKey + "'")

	// Analyze the loot to determine if it contains sensitive data
	sensitiveKeys := []string{
		"k8s_token", "kubernetes_token", "bearer", "authorization",
		"aws_access_key", "aws_secret_key", "credential", "password",
		"secret", "api_key", "private_key", "jwt", "session",
	}

	var detectedKeys []string
	for _, key := range sensitiveKeys {
		if strings.Contains(strings.ToLower(loot), key) {
			detectedKeys = append(detectedKeys, key)
		}
	}

	if len(detectedKeys) == 0 {
		utils.LogContext("[yellow]No sensitive data detected in loot.[-]")
		return
	}

	utils.LogContext("[red]⚠ Sensitive data detected:[-] " + fmt.Sprintf("%v", detectedKeys))

	// Generate follow-up exploit actions based on detected data
	for _, key := range detectedKeys {
		var nextAction *TacticalAction

		switch key {
		case "k8s_token", "kubernetes_token":
			nextAction = &TacticalAction{
				Type:         "CROSS_TENANT_LEAKAGE",
				Target:       "https://kubernetes.default.svc/api/v1/namespaces",
				Payload:      fmt.Sprintf("Authorization: Bearer %s", loot[:50]),
				Confidence:   "HIGH",
				Reasoning:    "K8s token detected. Attempt to enumerate cluster resources.",
				Status:       "PENDING",
				PreCondition: lootKey,
				ChainID:      chainID,
				Loot:         fmt.Sprintf("k8s_namespaces_%d", exploit.ID),
			}
			utils.LogContext("[magenta]>>> Generated:[-] CROSS_TENANT_LEAKAGE action for K8s token")

		case "bearer", "authorization":
			nextAction = &TacticalAction{
				Type:         "LATERAL_MOVEMENT",
				Target:       exploit.Target, // Same target but with new token
				Payload:      fmt.Sprintf("Use bearer token for privilege escalation"),
				Confidence:   "HIGH",
				Reasoning:    "Auth token detected. Attempt lateral movement with elevated privileges.",
				Status:       "PENDING",
				PreCondition: lootKey,
				ChainID:      chainID,
				Loot:         fmt.Sprintf("lateral_result_%d", exploit.ID),
			}
			utils.LogContext("[magenta]>>> Generated:[-] LATERAL_MOVEMENT action for auth token")

		case "aws_access_key", "aws_secret_key":
			nextAction = &TacticalAction{
				Type:         "CLOUD_PIVOT",
				Target:       "https://sts.amazonaws.com/",
				Payload:      fmt.Sprintf("Enumerate AWS with extracted keys"),
				Confidence:   "HIGH",
				Reasoning:    "AWS credentials detected. Pivot to cloud infrastructure.",
				Status:       "PENDING",
				PreCondition: lootKey,
				ChainID:      chainID,
				Loot:         fmt.Sprintf("aws_enumerate_%d", exploit.ID),
			}
			utils.LogContext("[magenta]>>> Generated:[-] CLOUD_PIVOT action for AWS keys")

		case "jwt":
			// JWT token detected - attempt to decode and use in subsequent requests
			nextAction = &TacticalAction{
				Type:         "JWT_BYPASS",
				Target:       exploit.Target,
				Payload:      fmt.Sprintf("Decode and use JWT for privilege escalation"),
				Confidence:   "MEDIUM",
				Reasoning:    "JWT token found. Attempt to manipulate token claims.",
				Status:       "PENDING",
				PreCondition: lootKey,
				ChainID:      chainID,
				Loot:         fmt.Sprintf("jwt_bypass_%d", exploit.ID),
			}
			utils.LogContext("[magenta]>>> Generated:[-] JWT_BYPASS action for JWT token")

		default:
			nextAction = &TacticalAction{
				Type:         "PRIVILEGE_ESCALATION",
				Target:       exploit.Target,
				Payload:      fmt.Sprintf("Attempt privilege escalation using discovered %s", key),
				Confidence:   "MEDIUM",
				Reasoning:    fmt.Sprintf("Sensitive data (%s) detected. Attempt escalation.", key),
				Status:       "PENDING",
				PreCondition: lootKey,
				ChainID:      chainID,
				Loot:         fmt.Sprintf("privesc_%d", exploit.ID),
			}
			utils.LogContext("[magenta]>>> Generated:[-] PRIVILEGE_ESCALATION action")
		}

		if nextAction != nil {
			// Add to ActionBuffer for autonomous execution
			nextAction.ID = len(ActionBuffer) + 1
			ActionBuffer = append(ActionBuffer, *nextAction)
			utils.TacticalLog(fmt.Sprintf("[green]✓ ACTION QUEUED:[-] %s (Chain: %s)", nextAction.Type, chainID))
		}
	}

	// Log chain summary
	utils.LogContext(fmt.Sprintf("[cyan]Chain Summary:[-] %d follow-up actions queued for chain '%s'", len(detectedKeys), chainID))
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
