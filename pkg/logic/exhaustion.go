package logic

import (
	"fmt"
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

func (e *ExhaustionContext) FuzzPagination() {
	// FIX: Removed pterm.DefaultHeader
	utils.TacticalLog("[cyan]API4: Resource Exhaustion - Pagination Fuzzer Started[-]")

	for _, val := range testLimits {
		u, err := url.Parse(e.TargetURL)
		if err != nil {
			return
		}

		q := u.Query()
		q.Set(e.ParamName, val)
		u.RawQuery = q.Encode()
		fuzzedURL := u.String()

		start := time.Now()
		req, _ := http.NewRequest("GET", fuzzedURL, nil)

		if CurrentSession.AttackerToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", CurrentSession.AttackerToken))
		}

		resp, err := GlobalClient.Do(req)
		duration := time.Since(start)

		if err != nil {
			break
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if duration > 2*time.Second {
				utils.RecordFinding(db.Finding{
					Phase:      "PHASE 9.9: EXHAUSTION",
					Target:     e.TargetURL,
					Details:    fmt.Sprintf("Resource Exhaustion via %s=%s (Latency: %v)", e.ParamName, val, duration),
					Status:     "VULNERABLE",
					OWASP_ID:   "API4:2023",
					MITRE_ID:   "T1499",
					NIST_Tag:   "RS.AN",
					CVE_ID:     "CVE-202X-GENERIC-DOS", // Generic Classification
					CVSS_Score: "5.3",                  // CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:L
				})
			}
		}
	}
	utils.TacticalLog(fmt.Sprintf("[green]✔[-] Exhaustion probe complete on %s", e.ParamName))
}
