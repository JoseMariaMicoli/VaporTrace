package logic

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db" // Added Persistence
	"github.com/pterm/pterm"
)

// BOLAContext defines the parameters for an ID-swap attack
type BOLAContext struct {
	BaseURL       string
	VictimID      string
	AttackerToken string
}

func (b *BOLAContext) Probe() {
	pterm.DefaultHeader.WithFullWidth(false).Println("BOLA Probe Engine (API1:2023)")

	// 1. Token Priority Logic: Context-specific -> Global Store -> Error
	activeToken := b.AttackerToken
	if activeToken == "" {
		activeToken = CurrentSession.AttackerToken
	}

	if activeToken == "" {
		pterm.Error.Println("No Attacker Token configured. Use 'auth attacker <token>' first.")
		return
	}

	// 2. Build Target URL safely
	u, err := url.Parse(b.BaseURL)
	if err != nil {
		pterm.Error.Printf("Invalid Base URL: %v\n", err)
		return
	}

	if b.VictimID != "" {
		u.Path = path.Join(u.Path, b.VictimID)
	}
	target := u.String()

	pterm.Info.Printf("Targeting: %s\n", target)
	pterm.Info.Printf("Using Token Snapshot: %s...\n", activeToken[:8])

	// PATCH: Use GlobalClient instead of creating a new one
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		pterm.Error.Printf("Failed to create request: %v\n", err)
		return
	}

	// Inject the resolved token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", activeToken))
	req.Header.Set("User-Agent", "VaporTrace/2.0.1 (API Recon Suite)")

	spinner, _ := pterm.DefaultSpinner.Start("Performing ID-Swap cross-validation...")

	resp, err := GlobalClient.Do(req)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Connection failed: %v", err))
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Phase 3 Logic Analysis: 200 OK vs 403 Forbidden
	if resp.StatusCode == http.StatusOK {
		spinner.Success("BOLA VULNERABILITY DETECTED")
		pterm.Warning.Prefix = pterm.Prefix{Text: "VULN", Style: pterm.NewStyle(pterm.BgRed, pterm.FgWhite)}
		pterm.Warning.Println("Attacker context accessed Victim resource successfully.")

		pterm.DefaultTable.WithData(pterm.TableData{
			{"METRIC", "VALUE"},
			{"STATUS CODE", "200 OK"},
			{"RESOURCE ID", b.VictimID},
			{"LEAK SIZE", fmt.Sprintf("%d bytes", len(body))},
		}).WithBoxed().Render()

		// PERSISTENCE HOOK: Log Success
		db.LogQueue <- db.Finding{
			Phase:   "PHASE III: AUTH LOGIC",
			Target:  target,
			Details: fmt.Sprintf("BOLA ID-Swap on %s", b.VictimID),
			Status:  "EXPLOITED",
		}

	} else if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		spinner.Success("Target Secure")
		pterm.Success.Println("Access denied: Authorization logic correctly enforced.")

		// PERSISTENCE HOOK: Log Mitigation
		db.LogQueue <- db.Finding{
			Phase:   "PHASE III: AUTH LOGIC",
			Target:  target,
			Details: "BOLA Attempt",
			Status:  "MITIGATED",
		}
	} else {
		spinner.Warning(fmt.Sprintf("Inconclusive: Server returned %d", resp.StatusCode))
	}
}