package logic

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url" // Added for robust URL parsing
	"time"

	"github.com/pterm/pterm"
)

type ExhaustionContext struct {
	TargetURL string
	ParamName string 
}

var testLimits = []string{"100", "1000", "10000", "50000", "1000000"}

func (e *ExhaustionContext) FuzzPagination() {
	pterm.DefaultHeader.WithFullWidth(false).Println("API4: Resource Exhaustion - Pagination Fuzzer")

	client := &http.Client{
		Timeout: 20 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	for _, val := range testLimits {
		// Robust URL Construction
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
		
		if CurrentSession.AttackerToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", CurrentSession.AttackerToken))
		}

		resp, err := client.Do(req)
		duration := time.Since(start)

		if err != nil {
			pterm.Warning.Printf("Request timed out or connection dropped at limit %s (%v)\n", val, err)
			break
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if duration > 2*time.Second {
				pterm.Warning.Prefix = pterm.Prefix{Text: "VULN", Style: pterm.NewStyle(pterm.BgRed, pterm.FgWhite)}
				pterm.Warning.Printf("Resource Exhaustion Detected! Limit %s took %v\n", val, duration)
			} else {
				pterm.Success.Printf("Limit %s processed in %v\n", val, duration)
			}
		} else {
			pterm.Info.Printf("Server rejected limit %s (Status: %d)\n", val, resp.StatusCode)
			break 
		}
	}
}