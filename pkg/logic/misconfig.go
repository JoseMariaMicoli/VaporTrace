package logic

import (
	"fmt" // Added for formatting
	"net/http"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db" // Added Persistence
	"github.com/pterm/pterm"
)

type MisconfigContext struct {
	TargetURL string
}

func (m *MisconfigContext) Audit() {
	pterm.DefaultHeader.WithFullWidth(false).Println("API8: Security Misconfiguration Audit")

	// PATCH: Removed local client definition

	// 1. CORS Audit
	req, _ := http.NewRequest("GET", m.TargetURL, nil)
	req.Header.Set("Origin", "https://evil-attacker.com")
	
	// PATCH: Using GlobalClient
	resp, err := GlobalClient.Do(req)
	if err != nil {
		pterm.Error.Printf("Audit failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	pterm.DefaultSection.Println("Header & CORS Analysis")
	
	cors := resp.Header.Get("Access-Control-Allow-Origin")
	if cors == "*" || cors == "https://evil-attacker.com" {
		pterm.Warning.Prefix = pterm.Prefix{Text: "VULN", Style: pterm.NewStyle(pterm.BgRed, pterm.FgWhite)}
		pterm.Warning.Printf("Weak CORS Policy: %s\n", cors)

		// PERSISTENCE HOOK
		db.LogQueue <- db.Finding{
			Phase:   "PHASE II: DISCOVERY",
			Target:  m.TargetURL,
			Details: fmt.Sprintf("Weak CORS Policy: %s", cors),
			Status:  "VULNERABLE",
		}
	} else {
		pterm.Success.Println("CORS policy appears restrictive.")
	}

	// 2. Security Headers Check
	headers := []string{"Strict-Transport-Security", "Content-Security-Policy", "X-Content-Type-Options"}
	for _, h := range headers {
		if resp.Header.Get(h) == "" {
			pterm.Info.Printf("Missing Security Header: %s\n", h)
			
			// Optional Persistence for Missing Headers (Info level)
			db.LogQueue <- db.Finding{
				Phase:   "PHASE II: DISCOVERY",
				Target:  m.TargetURL,
				Details: fmt.Sprintf("Missing Header: %s", h),
				Status:  "WEAK CONFIG",
			}
		} else {
			pterm.Success.Printf("Header Found: %s\n", h)
		}
	}

	pterm.Info.Println("\nAudit Complete. Triggering verbose error test...")
	m.TriggerVerboseError()
}

// PATCH: Removed client argument. Uses global client logic.
func (m *MisconfigContext) TriggerVerboseError() {
	// Attempting to trigger an error by sending a malformed Method/Payload
	req, _ := http.NewRequest("TRACE", m.TargetURL, nil)
	
	// PATCH: Using GlobalClient
	resp, err := GlobalClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		pterm.Warning.Println("Server returned 5xx error. Check response body for stack traces or debug info.")
		
		db.LogQueue <- db.Finding{
			Phase:   "PHASE II: DISCOVERY",
			Target:  m.TargetURL,
			Details: "Verbose Error / Stack Trace",
			Status:  "INFO LEAK",
		}
	}
}