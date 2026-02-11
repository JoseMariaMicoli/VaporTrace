package attack

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// AttackMode defines the strategy used for payload injection
type AttackMode int

const (
	Sniper AttackMode = iota // One payload set, one position
	// Future modes: Pitchfork, ClusterBomb, etc.
)

// IntruderConfig holds the configuration for the attack session
type IntruderConfig struct {
	TargetURL    string
	Param        string // The specific query parameter to fuzz
	WordlistPath string
	Concurrency  int
	Mode         AttackMode
}

// IntruderResult captures the outcome of a specific payload attempt
type IntruderResult struct {
	Payload       string
	StatusCode    int
	ContentLength int64
	Duration      time.Duration
	IsAnomaly     bool
	AnomalyReason string
}

// RunSniper executes a single-position fuzzing attack using a worker pool
func RunSniper(config IntruderConfig) {
	utils.TacticalLog(fmt.Sprintf("[magenta::b]INTRUDER SNIPER:[-] Targeting %s param '%s' with wordlist %s",
		config.TargetURL, config.Param, config.WordlistPath))

	// 1. Load Payloads
	payloads, err := loadWordlist(config.WordlistPath)
	if err != nil {
		utils.TacticalLog(fmt.Sprintf("[red]ERROR:[-] Failed to load wordlist: %v", err))
		return
	}
	utils.TacticalLog(fmt.Sprintf("[blue]LOADED:[-] %d payloads ready.", len(payloads)))

	// 2. Establish Baseline
	baseline, err := getBaseline(config.TargetURL)
	if err != nil {
		utils.TacticalLog(fmt.Sprintf("[red]ERROR:[-] Baseline request failed: %v", err))
		return
	}
	utils.TacticalLog(fmt.Sprintf("[cyan]BASELINE:[-] Status: %d | Length: %d", baseline.StatusCode, baseline.ContentLength))

	// 3. Initialize Worker Pool
	jobs := make(chan string, len(payloads))
	results := make(chan IntruderResult, len(payloads))
	var wg sync.WaitGroup

	// Cap concurrency to session defaults if not specified
	workers := config.Concurrency
	if workers <= 0 {
		workers = 10
	}

	utils.TacticalLog(fmt.Sprintf("[blue]WORKERS:[-] Spawning %d threads...", workers))

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for payload := range jobs {
				res := executeProbe(config.TargetURL, config.Param, payload, baseline)
				results <- res
			}
		}()
	}

	// 4. Feed Jobs
	go func() {
		for _, p := range payloads {
			jobs <- p
		}
		close(jobs)
	}()

	// 5. Monitor Results (Async)
	go func() {
		wg.Wait()
		close(results)
	}()

	// 6. Process Outcomes
	anomalies := 0
	for res := range results {
		if res.IsAnomaly {
			anomalies++
			utils.TacticalLog(fmt.Sprintf("[green]HIT:[-] %s | Status: %d | Len: %d (%s)",
				res.Payload, res.StatusCode, res.ContentLength, res.AnomalyReason))

			// Record to Database
			utils.RecordFinding(db.Finding{
				Phase:        "PHASE III: INTRUDER",
				Command:      "intruder",
				Target:       config.TargetURL,
				Details:      fmt.Sprintf("Intruder Anomaly: %s (Payload: %s)", res.AnomalyReason, res.Payload),
				Status:       "VULNERABLE",
				CVSS_Numeric: 5.0, // Default medium, needs manual triage
			})
		}
	}

	utils.TacticalLog(fmt.Sprintf("[green]✓ INTRUDER COMPLETE:[-] Scan finished. %d anomalies detected.", anomalies))
}

// executeProbe constructs the request, sends it, and analyzes the response
func executeProbe(baseURL, param, payload string, baseline IntruderResult) IntruderResult {
	// Construct URL with injected payload
	u, _ := url.Parse(baseURL)
	q := u.Query()
	q.Set(param, payload)
	u.RawQuery = q.Encode()
	target := u.String()

	req, _ := http.NewRequest("GET", target, nil)

	// Apply existing logic evasion (User-Agent rotation, etc)
	logic.ApplyEvasion(req)

	start := time.Now()
	// Use GlobalClient directly to avoid double-logging from SafeDo,
	// but ensure we use the configured transport (proxy, etc)
	resp, err := logic.GlobalClient.Do(req)
	duration := time.Since(start)

	result := IntruderResult{
		Payload:  payload,
		Duration: duration,
	}

	if err != nil {
		return result // Failed request
	}
	defer resp.Body.Close()

	// Calculate Content Length
	bodyBytes, _ := io.ReadAll(resp.Body)
	result.StatusCode = resp.StatusCode
	result.ContentLength = int64(len(bodyBytes))

	// --- ANALYSIS LOGIC ---
	// 1. Status Code Change
	if result.StatusCode != baseline.StatusCode {
		// Ignore 429s as anomalies (usually rate limits)
		if result.StatusCode != 429 && result.StatusCode != 503 {
			result.IsAnomaly = true
			result.AnomalyReason = fmt.Sprintf("Status %d != %d", result.StatusCode, baseline.StatusCode)
		}
	}

	// 2. Significant Length Change (>10% difference)
	diff := result.ContentLength - baseline.ContentLength
	if diff < 0 {
		diff = -diff
	}

	// Avoid division by zero
	baseLen := baseline.ContentLength
	if baseLen == 0 {
		baseLen = 1
	}

	percentDiff := (float64(diff) / float64(baseLen)) * 100
	if percentDiff > 10.0 && !result.IsAnomaly { // Only flag if not already flagged by status
		result.IsAnomaly = true
		result.AnomalyReason = fmt.Sprintf("Length Delta %.0f%% (%d bytes)", percentDiff, diff)
	}

	return result
}

// getBaseline captures the behavior of the target without fuzzing
func getBaseline(target string) (IntruderResult, error) {
	req, _ := http.NewRequest("GET", target, nil)
	logic.ApplyEvasion(req)

	start := time.Now()
	resp, err := logic.GlobalClient.Do(req)
	if err != nil {
		return IntruderResult{}, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	return IntruderResult{
		Payload:       "BASELINE",
		StatusCode:    resp.StatusCode,
		ContentLength: int64(len(bodyBytes)),
		Duration:      time.Since(start),
	}, nil
}

// loadWordlist reads a file line by line
func loadWordlist(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}
