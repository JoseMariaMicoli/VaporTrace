package logic

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/pterm/pterm"
)

type MisconfigContext struct {
	TargetURL string
}

func (m *MisconfigContext) Audit() {
	pterm.DefaultHeader.WithFullWidth(false).Println("API8: Security Misconfiguration Audit")

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// 1. CORS Audit
	req, _ := http.NewRequest("GET", m.TargetURL, nil)
	req.Header.Set("Origin", "https://evil-attacker.com")
	
	resp, err := client.Do(req)
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
	} else {
		pterm.Success.Println("CORS policy appears restrictive.")
	}

	// 2. Security Headers Check
	headers := []string{"Strict-Transport-Security", "Content-Security-Policy", "X-Content-Type-Options"}
	for _, h := range headers {
		if resp.Header.Get(h) == "" {
			pterm.Info.Printf("Missing Security Header: %s\n", h)
		} else {
			pterm.Success.Printf("Header Found: %s\n", h)
		}
	}

	pterm.Info.Println("\nAudit Complete. Triggering verbose error test...")
	m.TriggerVerboseError(client)
}

func (m *MisconfigContext) TriggerVerboseError(client *http.Client) {
	// Attempting to trigger an error by sending a malformed Method/Payload
	req, _ := http.NewRequest("TRACE", m.TargetURL, nil)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		pterm.Warning.Println("Server returned 5xx error. Check response body for stack traces or debug info.")
	}
}