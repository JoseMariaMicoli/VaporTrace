package logic

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/pterm/pterm"
)

// ExhaustionContext defines the parameters for API4:2023 Resource Exhaustion
type ExhaustionContext struct {
	TargetURL string
	ParamName string 
}

// testLimits defines common thresholds to test for lack of resources limiting
var testLimits = []string{"100", "1000", "10000", "50000", "1000000"}

// FuzzPagination probes the endpoint with increasing limit values to detect DoS or memory pressure
func (e *ExhaustionContext) FuzzPagination() {
	pterm.DefaultHeader.WithFullWidth(false).Println("API4: Resource Exhaustion - Pagination Fuzzer")

	for _, val := range testLimits {
		// Construct the target URL with the fuzzed parameter
		u, err := url.Parse(e.TargetURL)
		if err != nil {
			pterm.Error.Printf("Invalid Target URL: %v\n", err)
			return
		}

		q := u.Query()
		q.Set(e.ParamName, val)
		u.RawQuery = q.Encode()
		fuzzedURL := u.String()

		pterm.Info.Printf("Probing Limit: %s | URL: %s\n", val, fuzzedURL)

		start := time.Now()
		req, _ := http.NewRequest("GET", fuzzedURL, nil)
		
		// Use the synchronized session token
		if CurrentSession.AttackerToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", CurrentSession.AttackerToken))
		}

		// Execute via GlobalClient (managed in store.go)
		resp, err := GlobalClient.Do(req)
		duration := time.Since(start)

		if err != nil {
			pterm.Warning.Printf("Request timed out or connection dropped at limit %s (%v)\n", val, err)
			break
		}
		defer resp.Body.Close()

		// Analysis: If the response time spikes significantly, it indicates a vulnerability
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if duration > 2*time.Second {
				pterm.Warning.Prefix = pterm.Prefix{Text: "VULN", Style: pterm.NewStyle(pterm.BgRed, pterm.FgWhite)}
				pterm.Warning.Printf("Resource Exhaustion Detected! Limit %s took %v\n", val, duration)

				db.LogQueue <- db.Finding{
					Phase:   "PHASE 9.9: EXHAUSTION",
					Target:  e.TargetURL,
					Details: fmt.Sprintf("Resource Exhaustion via %s=%s (Latency: %v)", e.ParamName, val, duration),
					Status:  "VULNERABLE",
				}
			}
		}
	}
}