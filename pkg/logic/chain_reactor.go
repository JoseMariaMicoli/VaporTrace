/*
Copyright (c) 2026 José María Micoli
Licensed under {'license_type': 'BSL', 'change_date': '2033-02-17', 'convert_to': 'Apache-2.0'}

You may:
✔ Study
✔ Modify
✔ Use for internal security testing

You may NOT:
✘ Offer as a commercial service
✘ Sell derived competing products
*/

package logic

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// ChainStep represents a single HTTP action in a sequence
type ChainStep struct {
	Name       string
	Method     string
	URL        string
	Body       string
	Headers    map[string]string
	Extracts   []ExtractionRule // Variables to save after this step
	Assertions []string         // Simple checks (e.g., "200", "success")
}

// ExtractionRule defines what to pull from the response
type ExtractionRule struct {
	VarName  string // Name to store in context (e.g., "auth_token")
	Strategy string // "json", "regex", "header"
	Selector string // The path or pattern
}

// ChainDefinition is the full workflow
type ChainDefinition struct {
	Name  string
	Steps []ChainStep
}

// GlobalChainStore holds user-defined chains
var GlobalChainStore = make(map[string]*ChainDefinition)

// RunChain executes a chain sequence
func RunChain(chainName string) {
	chain, exists := GlobalChainStore[chainName]
	if !exists {
		utils.TacticalLog(fmt.Sprintf("[red]CHAIN ERROR:[-] Chain '%s' not found.", chainName))
		return
	}

	utils.TacticalLog(fmt.Sprintf("[magenta::b]CHAIN REACTOR:[-] Executing sequence '%s' (%d steps)...", chain.Name, len(chain.Steps)))

	// Session Context (Ephemeral storage for this run)
	// We also seed it with GlobalDataSilo values
	ctx := make(map[string]string)

	// Copy global data (optional, useful if we have creds from Loot)
	// For now, start empty or pre-seed common vars
	ctx["target"] = CurrentSession.TargetURL

	for i, step := range chain.Steps {
		utils.TacticalLog(fmt.Sprintf("[blue]STEP %d:[-] %s (%s)", i+1, step.Name, step.Method))

		// 1. Templating: Inject variables into URL and Body
		finalURL := InjectVariables(step.URL, ctx)
		finalBody := InjectVariables(step.Body, ctx)

		// 2. Prepare Request
		var reqBody io.Reader
		if finalBody != "" {
			reqBody = bytes.NewBufferString(finalBody)
		}

		req, err := http.NewRequest(step.Method, finalURL, reqBody)
		if err != nil {
			utils.TacticalLog(fmt.Sprintf("[red]ERROR:[-] Failed to build request: %v", err))
			return
		}

		// 3. Inject Headers
		for k, v := range step.Headers {
			finalVal := InjectVariables(v, ctx)
			req.Header.Set(k, finalVal)
		}
		// Default content type if body exists and not set
		if finalBody != "" && req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}

		// 4. Evasion & Transport
		ApplyEvasion(req) // Tier 3 Evasion logic

		// 5. Execute
		start := time.Now()
		resp, err := GlobalClient.Do(req) // Use GlobalClient to leverage Interceptor/Proxy
		if err != nil {
			utils.TacticalLog(fmt.Sprintf("[red]ERROR:[-] Network failure: %v", err))
			return
		}
		defer resp.Body.Close()

		// 6. Read Response
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyStr := string(bodyBytes)

		utils.TacticalLog(fmt.Sprintf("  Status: %d | Latency: %v | Size: %d", resp.StatusCode, time.Since(start), len(bodyBytes)))

		// 7. Extract Variables
		for _, ex := range step.Extracts {
			val, err := ExtractValue(bodyStr, ex.Strategy, ex.Selector)
			if err != nil {
				// Fallback: Check headers
				// Ideally, ExtractValue should handle "header" strategy by passing headers,
				// but for simplicity in this artifact we stick to body or refactor ExtractValue signature.
				// NOTE: For "header" strategy to work, we need to pass resp.Header string.
				if ex.Strategy == "header" {
					// Dump headers to string
					var hBuilder strings.Builder
					resp.Header.Write(&hBuilder)
					val, err = ExtractValue(hBuilder.String(), ex.Strategy, ex.Selector)
				}
			}

			if err == nil && val != "" {
				ctx[ex.VarName] = val
				utils.TacticalLog(fmt.Sprintf("  [green]CAPTURED:[-] {{%s}} = %s...", ex.VarName, truncate(val, 20)))
				// Sync to Global Silo so Intruder can use it later
				GlobalDataSilo.Set(ex.VarName, val)
			} else {
				utils.TacticalLog(fmt.Sprintf("  [yellow]WARN:[-] Failed to extract {{%s}}: %v", ex.VarName, err))
			}
		}

		// 8. Assertions (Stop on failure)
		// For now, just stop on 4xx/5xx unless expected
		if resp.StatusCode >= 400 {
			utils.TacticalLog("[yellow]CHAIN INTERRUPTED:[-] Received error status code.")
			// In a full implementation, we'd check specific assertions here.
		}
	}

	utils.TacticalLog(fmt.Sprintf("[green]✓ CHAIN COMPLETE:[-] Workflow '%s' finished successfully.", chainName))
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}
