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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

type ExhaustionContext struct {
	TargetURL string
	ParamName string
}

var testLimits = []string{"100", "1000", "10000", "50000", "1000000"}

func estimateResponseSize(resp *http.Response, body []byte) int64 {
	if resp != nil && resp.ContentLength > 0 {
		return resp.ContentLength
	}
	return int64(len(body))
}

func (e *ExhaustionContext) FuzzPagination() {
	// FIXED: Now measures actual resource consumption (response size, latency, server behavior)
	utils.TacticalLog("[cyan]API4: Resource Exhaustion - Pagination Fuzzer Started[-]")

	var prevSize int64 = 0
	var baselineLatency time.Duration = 0
	exhaustionDetected := false
	exhaustionPattern := ""
	tested := 0
	failures := 0

	for idx, val := range testLimits {
		u, err := url.Parse(e.TargetURL)
		if err != nil {
			utils.TacticalLog(fmt.Sprintf("[red]EXHAUSTION ERROR:[-] invalid target URL: %v", err))
			return
		}

		q := u.Query()
		q.Set(e.ParamName, val)
		u.RawQuery = q.Encode()
		fuzzedURL := u.String()

		utils.TacticalLog(fmt.Sprintf("[blue]Testing %s=%s...", e.ParamName, val))

		start := time.Now()
		req, _ := http.NewRequest("GET", fuzzedURL, nil)

		if CurrentSession.AttackerToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", CurrentSession.AttackerToken))
		}

		resp, err := SafeDo(req, false, "EXHAUSTION")
		duration := time.Since(start)
		tested++

		if err != nil {
			failures++
			utils.TacticalLog(fmt.Sprintf("[yellow]Request error at %s=%s: %v (possible exhaustion)", e.ParamName, val, err))
			if failures >= 2 {
				exhaustionPattern = fmt.Sprintf("Repeated connection failures from %s=%s onward", e.ParamName, val)
				exhaustionDetected = true
				break
			}
			continue
		}

		// FIXED: Measure actual response body size
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
		respSize := estimateResponseSize(resp, bodyBytes)
		resp.Body.Close()

		utils.TacticalLog(fmt.Sprintf("[cyan]Response: %d bytes in %dms (Status: %d)", respSize, duration.Milliseconds(), resp.StatusCode))

		// Record baseline on first request
		if idx == 0 {
			baselineLatency = duration
			prevSize = respSize
			continue
		}

		// FIXED: Detect exhaustion patterns
		latencyGrowth := 1.0
		if baselineLatency > 0 {
			latencyGrowth = float64(duration) / float64(baselineLatency)
		}

		sizeGrowth := 1.0
		if prevSize > 0 {
			sizeGrowth = float64(respSize) / float64(prevSize)
		}

		// Pattern 1: Exponential latency increase (server struggling)
		if latencyGrowth > 5.0 && duration > 1500*time.Millisecond {
			exhaustionDetected = true
			exhaustionPattern = fmt.Sprintf("Exponential latency spike: %.1fx slowdown at %s=%s", latencyGrowth, e.ParamName, val)
			utils.TacticalLog(fmt.Sprintf("[red]EXHAUSTION DETECTED:[-] %s", exhaustionPattern))
			break
		}

		// Pattern 2: Server returning large datasets (memory pressure)
		if sizeGrowth > 100.0 && respSize > 10*1024*1024 { // > 10MB
			exhaustionDetected = true
			exhaustionPattern = fmt.Sprintf("Memory exhaustion: %.0fx size growth (%d bytes) at %s=%s", sizeGrowth, respSize, e.ParamName, val)
			utils.TacticalLog(fmt.Sprintf("[red]EXHAUSTION DETECTED:[-] %s", exhaustionPattern))
			break
		}

		// Pattern 3: Server timeout/503 at high values
		if resp.StatusCode >= 500 {
			exhaustionDetected = true
			exhaustionPattern = fmt.Sprintf("Server error (HTTP %d) at %s=%s - resource exhaustion triggered", resp.StatusCode, e.ParamName, val)
			utils.TacticalLog(fmt.Sprintf("[red]EXHAUSTION DETECTED:[-] %s", exhaustionPattern))
			break
		}

		prevSize = respSize
	}

	// Record finding if exhaustion was detected
	if exhaustionDetected {
		utils.RecordFinding(db.Finding{
			Phase:   "PHASE 9.9: EXHAUSTION",
			Command: "exhaust",
			Target:  e.TargetURL,
			Details: fmt.Sprintf("API4:2023 Resource Exhaustion: %s", exhaustionPattern),
			Status:  "VULNERABLE",
		})
		utils.TacticalLog("[green]✔ Finding recorded[-]")
		utils.TacticalLog(fmt.Sprintf("[green]EXHAUST COMPLETE:[-] tested=%d failures=%d result=VULNERABLE", tested, failures))
	} else {
		utils.TacticalLog(fmt.Sprintf("[yellow]EXHAUST COMPLETE:[-] tested=%d failures=%d result=NO-EVIDENCE", tested, failures))
	}

	utils.TacticalLog(fmt.Sprintf("[green]✔[-] Exhaustion probe complete on %s", e.ParamName))
}
